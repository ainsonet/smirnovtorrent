package engine

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anacrolix/dht/v2"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/anacrolix/torrent/types"
	"golang.org/x/time/rate"
	"smirnovtorrent/config"
)

type speedTracker struct {
	lastReadBytes    int64
	lastWrittenBytes int64
	lastTime         time.Time
	downloadSpeed    int64
	uploadSpeed      int64
}

type Engine struct {
	mu              sync.RWMutex
	cfgManager      *config.ConfigManager
	sessionStore    *SessionStore
	pieceStore      *PieceStore
	client          *torrent.Client
	downloadLimiter *rate.Limiter
	uploadLimiter   *rate.Limiter
	speedTrack      map[string]*speedTracker
	verifying       map[string]struct{}
	stopCh          chan struct{}
	closed          bool
}

func NewEngine(cfgManager *config.ConfigManager) (*Engine, error) {
	cfg := cfgManager.Get()

	var cl *torrent.Client
	var err error

	downloadLimiter := rate.NewLimiter(rate.Inf, 1<<20)
	uploadLimiter := rate.NewLimiter(rate.Inf, 1<<20)
	if cfg.DownloadLimitBps > 0 {
		downloadLimiter.SetLimit(rate.Limit(cfg.DownloadLimitBps))
		downloadLimiter.SetBurst(max(int(cfg.DownloadLimitBps), 1<<20))
	}
	if cfg.UploadLimitBps > 0 {
		uploadLimiter.SetLimit(rate.Limit(cfg.UploadLimitBps))
		uploadLimiter.SetBurst(max(int(cfg.UploadLimitBps), 1<<20))
	}

	pieceStore, psErr := NewPieceStore(config.GetConfigDir())
	if psErr != nil {
		log.Printf("[Engine] Warning: failed to initialize persistent piece store: %v", psErr)
	}

	var dummyListeners []io.Closer
	defer func() {
		for _, l := range dummyListeners {
			_ = l.Close()
		}
	}()

	listenPort := cfg.ListenPort
	maxRetries := 5
	if listenPort > 0 {
		maxRetries = 2
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		tcfg := torrent.NewDefaultClientConfig()
		tcfg.DataDir = cfg.DownloadDir
		tcfg.DefaultStorage = storage.NewFileOpts(storage.NewFileClientOpts{
			ClientBaseDir:   cfg.DownloadDir,
			PieceCompletion: pieceStoreWrapper{pieceStore},
		})
		tcfg.NoDHT = !cfg.EnableDHT
		tcfg.DisablePEX = !cfg.EnablePEX
		tcfg.NoDefaultPortForwarding = !cfg.EnableUPnP
		tcfg.DisableUTP = false
		tcfg.NoUpload = false
		tcfg.EstablishedConnsPerTorrent = cfg.MaxEstablished
		tcfg.ListenPort = listenPort
		tcfg.DownloadRateLimiter = downloadLimiter
		tcfg.UploadRateLimiter = uploadLimiter

		cl, err = torrent.NewClient(tcfg)
		if err == nil {
			break
		}
		if cl != nil {
			cl.Close()
			cl = nil
		}
		if attempt < maxRetries {
			if listenPort == 0 {
				if l, lErr := net.Listen("tcp4", ":0"); lErr == nil {
					dummyListeners = append(dummyListeners, l)
				}
				if pc, pErr := net.ListenPacket("udp4", ":0"); pErr == nil {
					dummyListeners = append(dummyListeners, pc)
				}
			}
			time.Sleep(time.Duration(attempt+1) * 75 * time.Millisecond)
		} else if listenPort > 0 {
			// Fallback: If configured static port could not be bound, switch to dynamic port 0
			log.Printf("[Engine] Warning: configured ListenPort %d could not be bound (%v); falling back to random available port", listenPort, err)
			listenPort = 0
			attempt = -1
			maxRetries = 5
		}
	}

	if err != nil {
		if pieceStore != nil {
			_ = pieceStore.Close()
		}
		return nil, fmt.Errorf("failed to create torrent client: %w", err)
	}

	eng := &Engine{
		cfgManager:      cfgManager,
		sessionStore:    NewSessionStore(),
		pieceStore:      pieceStore,
		client:          cl,
		downloadLimiter: downloadLimiter,
		uploadLimiter:   uploadLimiter,
		speedTrack:      make(map[string]*speedTracker),
		verifying:       make(map[string]struct{}),
		stopCh:          make(chan struct{}),
	}

	// Restore saved torrents in background
	go eng.restoreSession()

	// Start 1-second speed & stats ticker
	go eng.statsLoop()

	return eng, nil
}

func (e *Engine) newTorrentStorage(savePath string) storage.ClientImpl {
	return storage.NewFileOpts(storage.NewFileClientOpts{
		ClientBaseDir:   savePath,
		PieceCompletion: pieceStoreWrapper{e.pieceStore},
	})
}

func (e *Engine) restoreSession() {
	saved := e.sessionStore.All()
	defaultDir := e.cfgManager.Get().DownloadDir
	for _, item := range saved {
		var t *torrent.Torrent
		var err error
		savePath := item.SavePath
		if savePath == "" {
			savePath = defaultDir
		}

		if item.TorrentFile != "" && fileExists(item.TorrentFile) {
			mi, loadErr := metainfo.LoadFromFile(item.TorrentFile)
			if loadErr == nil {
				spec := torrent.TorrentSpecFromMetaInfo(mi)
				spec.Storage = e.newTorrentStorage(savePath)
				if item.MagnetURI != "" {
					if magSpec, parseErr := torrent.TorrentSpecFromMagnetUri(item.MagnetURI); parseErr == nil {
						if len(spec.Trackers) == 0 && len(magSpec.Trackers) > 0 {
							spec.Trackers = magSpec.Trackers
						}
					}
				}
				t, _, err = e.client.AddTorrentSpec(spec)
			} else {
				err = loadErr
			}
		} else if item.MagnetURI != "" {
			spec, parseErr := torrent.TorrentSpecFromMagnetUri(item.MagnetURI)
			if parseErr == nil {
				spec.Storage = e.newTorrentStorage(savePath)
				t, _, err = e.client.AddTorrentSpec(spec)
			} else {
				err = parseErr
			}
		}

		if err != nil {
			log.Printf("[Engine] Failed to restore torrent %s: %v", item.InfoHash, err)
			continue
		}

		if t != nil {
			cfg := e.cfgManager.Get()
			if cfg.MaxEstablished > 0 {
				t.SetMaxEstablishedConns(cfg.MaxEstablished)
			}

			if item.IsPaused {
				t.DisallowDataDownload()
			} else {
				t.DownloadAll()
			}

			if t.Info() != nil {
				if len(item.FilePriorities) > 0 {
					e.applyFilePriorities(t, item.FilePriorities)
				}
				if item.Sequential && !item.IsPaused {
					e.applySequential(t)
				}
			}

			// Wait for info asynchronously to save .torrent metadata if needed
			go e.watchTorrentGotInfo(t, item)
		}
	}
}

func (e *Engine) watchTorrentGotInfo(t *torrent.Torrent, item *SavedTorrent) {
	select {
	case <-t.GotInfo():
		ih := strings.ToLower(t.InfoHash().HexString())

		// Verify that the torrent was not dropped/deleted
		if _, ok := e.findTorrent(ih); !ok {
			return
		}

		saved, found := e.sessionStore.Get(ih)
		if !found {
			return
		}

		// Restore saved file priorities if any
		if len(saved.FilePriorities) > 0 {
			e.applyFilePriorities(t, saved.FilePriorities)
		}

		if t.Name() != "" {
			saved.Name = t.Name()
		}

		if saved.Sequential && !saved.IsPaused {
			e.applySequential(t)
		}

		// Save .torrent file copy if not present
		if saved.TorrentFile == "" || !fileExists(saved.TorrentFile) {
			torPath := filepath.Join(e.sessionStore.GetTorrentDir(), ih+".torrent")
			f, err := os.Create(torPath)
			if err == nil {
				mi := t.Metainfo()
				writeErr := mi.Write(f)
				closeErr := f.Close()
				if writeErr == nil && closeErr == nil {
					saved.TorrentFile = torPath
				} else {
					_ = os.Remove(torPath)
				}
			}
		}

		// Check once more that torrent wasn't deleted while writing file
		if _, ok := e.findTorrent(ih); !ok {
			if saved.TorrentFile != "" {
				_ = os.Remove(saved.TorrentFile)
			}
			return
		}

		if latest, ok := e.sessionStore.Get(ih); ok {
			if saved.Name != "" {
				latest.Name = saved.Name
			}
			if saved.TorrentFile != "" {
				latest.TorrentFile = saved.TorrentFile
			}
			_ = e.sessionStore.Put(latest)
		} else {
			_ = e.sessionStore.Put(saved)
		}

		// Bugfix: If content files do not exist on disk, ensure we don't have stale piece completion
		// that causes the torrent to falsely think it is complete and get stuck in "missing"
		savePath := saved.SavePath
		if savePath == "" {
			savePath = e.cfgManager.Get().DownloadDir
		}
		if !e.checkFilesExist(t, savePath) {
			if e.pieceStore != nil {
				_ = e.pieceStore.DeleteTorrent(t.InfoHash().Bytes())
			}
			if t.BytesCompleted() > 0 || t.Seeding() {
				go func() {
					_ = e.VerifyData(ih)
					if !saved.IsPaused {
						t.AllowDataDownload()
						t.DownloadAll()
					}
				}()
			}
		}
	case <-t.Closed():
		return
	case <-e.stopCh:
		return
	}
}

func (e *Engine) applyFilePriorities(t *torrent.Torrent, priorities map[int]int) {
	if t == nil || t.Info() == nil || len(priorities) == 0 {
		return
	}
	files := t.Files()
	for idx, prio := range priorities {
		if idx >= 0 && idx < len(files) {
			f := files[idx]
			switch prio {
			case 0:
				f.SetPriority(types.PiecePriorityNone)
				t.CancelPieces(f.BeginPieceIndex(), f.EndPieceIndex())
			case 1:
				f.SetPriority(types.PiecePriorityNormal)
				t.DownloadPieces(f.BeginPieceIndex(), f.EndPieceIndex())
			case 2:
				f.SetPriority(types.PiecePriorityHigh)
				t.DownloadPieces(f.BeginPieceIndex(), f.EndPieceIndex())
			default:
				f.SetPriority(types.PiecePriorityNormal)
				t.DownloadPieces(f.BeginPieceIndex(), f.EndPieceIndex())
			}
		}
	}
}

func (e *Engine) statsLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.updateSpeeds()
			e.updateSequential()
		case <-e.stopCh:
			return
		}
	}
}

func (e *Engine) updateSequential() {
	torrents := e.client.Torrents()
	for _, t := range torrents {
		ih := strings.ToLower(t.InfoHash().HexString())
		saved, ok := e.sessionStore.Get(ih)
		if ok && saved.Sequential && !saved.IsPaused {
			e.applySequential(t)
		}
	}
}

func (e *Engine) applySequential(t *torrent.Torrent) {
	if t == nil || t.Info() == nil {
		return
	}
	numPieces := t.NumPieces()
	if numPieces <= 0 {
		return
	}

	windowSize := 10
	highCount := 0

	for i := 0; i < numPieces; i++ {
		ps := t.PieceState(i)
		// Never touch pieces that are already complete, and never touch pieces
		// belonging to deselected files (priority None).
		if ps.Complete || ps.Priority == types.PiecePriorityNone {
			continue
		}
		if highCount < windowSize {
			t.Piece(i).SetPriority(types.PiecePriorityHigh)
			highCount++
		} else {
			t.Piece(i).SetPriority(types.PiecePriorityNormal)
		}
	}
}

func (e *Engine) updateSpeeds() {
	torrents := e.client.Torrents()
	now := time.Now()

	type speedSample struct {
		ih           string
		readBytes    int64
		writtenBytes int64
	}
	samples := make([]speedSample, 0, len(torrents))
	for _, t := range torrents {
		stats := t.Stats()
		samples = append(samples, speedSample{
			ih:           strings.ToLower(t.InfoHash().HexString()),
			readBytes:    stats.BytesReadUsefulData.Int64(),
			writtenBytes: stats.BytesWrittenData.Int64(),
		})
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	activeHashes := make(map[string]struct{}, len(samples))
	for _, s := range samples {
		activeHashes[s.ih] = struct{}{}
		st, exists := e.speedTrack[s.ih]
		if !exists {
			e.speedTrack[s.ih] = &speedTracker{
				lastReadBytes:    s.readBytes,
				lastWrittenBytes: s.writtenBytes,
				lastTime:         now,
				downloadSpeed:    0,
				uploadSpeed:      0,
			}
			continue
		}

		elapsed := now.Sub(st.lastTime).Seconds()
		if elapsed >= 0.5 {
			dBytes := s.readBytes - st.lastReadBytes
			wBytes := s.writtenBytes - st.lastWrittenBytes
			if dBytes < 0 {
				dBytes = 0
			}
			if wBytes < 0 {
				wBytes = 0
			}

			st.downloadSpeed = int64(float64(dBytes) / elapsed)
			st.uploadSpeed = int64(float64(wBytes) / elapsed)
			st.lastReadBytes = s.readBytes
			st.lastWrittenBytes = s.writtenBytes
			st.lastTime = now
		}
	}

	for ih := range e.speedTrack {
		if _, ok := activeHashes[ih]; !ok {
			delete(e.speedTrack, ih)
		}
	}
}

func (e *Engine) AddMagnet(magnetURI string) (*TorrentSummary, error) {
	magnetURI = strings.TrimSpace(magnetURI)
	if magnetURI == "" {
		return nil, fmt.Errorf("empty magnet URI")
	}

	saveDir := e.cfgManager.Get().DownloadDir
	spec, err := torrent.TorrentSpecFromMagnetUri(magnetURI)
	if err != nil {
		return nil, fmt.Errorf("invalid magnet link: %w", err)
	}
	spec.Storage = e.newTorrentStorage(saveDir)

	t, _, err := e.client.AddTorrentSpec(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to add torrent: %w", err)
	}

	cfg := e.cfgManager.Get()
	if cfg.MaxEstablished > 0 {
		t.SetMaxEstablishedConns(cfg.MaxEstablished)
	}

	t.DownloadAll()

	ih := strings.ToLower(t.InfoHash().HexString())
	name := t.Name()
	if name == "" {
		name = "Downloading metadata..."
	}

	saved, found := e.sessionStore.Get(ih)
	if !found {
		saved = &SavedTorrent{
			InfoHash:   ih,
			Name:       name,
			MagnetURI:  magnetURI,
			AddedAt:    time.Now(),
			SavePath:   saveDir,
			IsPaused:   false,
			Sequential: false,
		}
	} else if saved.MagnetURI == "" {
		saved.MagnetURI = magnetURI
	}

	_ = e.sessionStore.Put(saved)
	go e.watchTorrentGotInfo(t, saved)

	summary := e.getTorrentSummary(t)
	return &summary, nil
}

func (e *Engine) AddTorrentFile(filePath string) (*TorrentSummary, error) {
	mi, err := metainfo.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse torrent file: %w", err)
	}

	saveDir := e.cfgManager.Get().DownloadDir
	spec := torrent.TorrentSpecFromMetaInfo(mi)
	spec.Storage = e.newTorrentStorage(saveDir)

	t, _, err := e.client.AddTorrentSpec(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to add torrent: %w", err)
	}

	cfg := e.cfgManager.Get()
	if cfg.MaxEstablished > 0 {
		t.SetMaxEstablishedConns(cfg.MaxEstablished)
	}

	t.DownloadAll()

	ih := strings.ToLower(t.InfoHash().HexString())
	name := t.Name()
	if name == "" {
		name = filepath.Base(filePath)
	}

	// Make a copy in session torrentDir
	destPath := filepath.Join(e.sessionStore.GetTorrentDir(), ih+".torrent")
	_ = copyFile(filePath, destPath)

	saved, found := e.sessionStore.Get(ih)
	if !found {
		saved = &SavedTorrent{
			InfoHash:    ih,
			Name:        name,
			TorrentFile: destPath,
			AddedAt:     time.Now(),
			SavePath:    saveDir,
			IsPaused:    false,
			Sequential:  false,
		}
	} else {
		saved.TorrentFile = destPath
	}

	_ = e.sessionStore.Put(saved)
	go e.watchTorrentGotInfo(t, saved)

	// If content files do not exist on disk, clean stale piece completion and verify
	if !e.checkFilesExist(t, saveDir) {
		if e.pieceStore != nil {
			_ = e.pieceStore.DeleteTorrent(t.InfoHash().Bytes())
		}
		if t.BytesCompleted() > 0 || t.Seeding() {
			go func() {
				_ = e.VerifyData(ih)
				t.AllowDataDownload()
				t.DownloadAll()
			}()
		}
	}

	summary := e.getTorrentSummary(t)
	return &summary, nil
}

func (e *Engine) PauseTorrent(infoHash string) error {
	t, ok := e.findTorrent(infoHash)
	if !ok {
		return fmt.Errorf("torrent not found")
	}

	t.DisallowDataDownload()

	if saved, found := e.sessionStore.Get(infoHash); found {
		saved.IsPaused = true
		_ = e.sessionStore.Put(saved)
	}

	return nil
}

func (e *Engine) ResumeTorrent(infoHash string) error {
	ih := strings.ToLower(infoHash)
	t, ok := e.findTorrent(ih)
	if !ok {
		return fmt.Errorf("torrent not found")
	}

	savePath := e.cfgManager.Get().DownloadDir
	if saved, found := e.sessionStore.Get(ih); found && saved.SavePath != "" {
		savePath = saved.SavePath
	}

	// If files are missing from disk, heal by re-verifying so missing pieces download
	if !e.checkFilesExist(t, savePath) && (t.BytesCompleted() > 0 || t.Seeding()) {
		if e.pieceStore != nil {
			_ = e.pieceStore.DeleteTorrent(t.InfoHash().Bytes())
		}
		_ = e.VerifyData(ih)
	}

	t.AllowDataDownload()
	t.DownloadAll()

	if saved, found := e.sessionStore.Get(ih); found {
		saved.IsPaused = false
		_ = e.sessionStore.Put(saved)

		if t.Info() != nil && len(saved.FilePriorities) > 0 {
			e.applyFilePriorities(t, saved.FilePriorities)
		}
		if saved.Sequential {
			e.applySequential(t)
		}
	}

	return nil
}

func (e *Engine) RedownloadTorrent(infoHash string) error {
	ih := strings.ToLower(infoHash)
	t, ok := e.findTorrent(ih)
	if !ok {
		return fmt.Errorf("torrent not found")
	}

	if e.pieceStore != nil {
		_ = e.pieceStore.DeleteTorrent(t.InfoHash().Bytes())
	}
	_ = e.VerifyData(ih)
	t.AllowDataDownload()
	t.DownloadAll()

	if saved, found := e.sessionStore.Get(ih); found {
		saved.IsPaused = false
		_ = e.sessionStore.Put(saved)
	}

	return nil
}

func (e *Engine) GetTorrentFilePath(infoHash string, fileIndex int) (string, error) {
	ih := strings.ToLower(infoHash)
	t, ok := e.findTorrent(ih)
	if !ok {
		return "", fmt.Errorf("torrent not found")
	}
	if t.Info() == nil {
		return "", fmt.Errorf("metadata not available yet")
	}
	files := t.Files()
	if fileIndex < 0 || fileIndex >= len(files) {
		return "", fmt.Errorf("invalid file index")
	}

	savePath := e.cfgManager.Get().DownloadDir
	if saved, found := e.sessionStore.Get(ih); found && saved.SavePath != "" {
		savePath = saved.SavePath
	}

	return filepath.Join(savePath, filepath.FromSlash(files[fileIndex].Path())), nil
}

func (e *Engine) DeleteTorrent(infoHash string, deleteFiles bool) error {
	ih := strings.ToLower(infoHash)
	t, ok := e.findTorrent(ih)

	var torrentName string
	savePath := e.cfgManager.Get().DownloadDir

	if saved, found := e.sessionStore.Get(ih); found {
		if saved.SavePath != "" {
			savePath = saved.SavePath
		}
		if saved.Name != "" {
			torrentName = saved.Name
		}
	}

	if ok {
		if t.Name() != "" {
			torrentName = t.Name()
		}
		t.Drop()
	}

	e.mu.Lock()
	delete(e.speedTrack, ih)
	delete(e.verifying, ih)
	e.mu.Unlock()

	_ = e.sessionStore.Remove(ih)

	// Clean piece completion synchronously so re-adding torrent never has stale completed pieces
	if e.pieceStore != nil {
		ihBytes, _ := hex.DecodeString(ih)
		if len(ihBytes) == 20 {
			_ = e.pieceStore.DeleteTorrent(ihBytes)
		}
	}

	if !ok && torrentName == "" {
		return fmt.Errorf("torrent not found")
	}

	// Guard against deleting when metadata is still downloading or invalid names
	if torrentName == "Downloading metadata..." {
		torrentName = ""
	}

	// Ensure torrentName is not an absolute path, has no volume name, and contains no path separators
	if filepath.IsAbs(torrentName) || filepath.VolumeName(torrentName) != "" || strings.ContainsAny(torrentName, `/\:`) {
		torrentName = filepath.Base(filepath.Clean(torrentName))
		if idx := strings.LastIndexAny(torrentName, `/\:`); idx != -1 {
			torrentName = torrentName[idx+1:]
		}
	}

	if torrentName == "." || torrentName == ".." || torrentName == "/" || torrentName == "\\" || torrentName == ":" || torrentName == "" {
		torrentName = ""
	}

	if deleteFiles && torrentName != "" {
		cleanSavePath := filepath.Clean(savePath)
		if cleanSavePath != "" && cleanSavePath != "." {
			targetPath := filepath.Clean(filepath.Join(cleanSavePath, torrentName))
			rel, err := filepath.Rel(cleanSavePath, targetPath)
			// Safety check: ensure targetPath is inside cleanSavePath and not the root save directory itself
			if err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				for r := 0; r < 5; r++ {
					rmErr := os.RemoveAll(targetPath)
					if rmErr == nil || os.IsNotExist(rmErr) {
						break
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
		}
	}

	return nil
}

func (e *Engine) SetSequential(infoHash string, seq bool) error {
	t, ok := e.findTorrent(infoHash)
	if !ok {
		return fmt.Errorf("torrent not found")
	}

	saved, found := e.sessionStore.Get(infoHash)
	if found {
		saved.Sequential = seq
		_ = e.sessionStore.Put(saved)
	}

	if seq {
		e.applySequential(t)
	} else {
		// When disabling sequential, reapply user's file priorities if any, else DownloadAll
		t.DownloadAll()
		if saved != nil && len(saved.FilePriorities) > 0 && t.Info() != nil {
			e.applyFilePriorities(t, saved.FilePriorities)
		}
	}

	return nil
}

func (e *Engine) VerifyData(infoHash string) error {
	ih := strings.ToLower(infoHash)
	t, ok := e.findTorrent(ih)
	if !ok {
		return fmt.Errorf("torrent not found")
	}

	if t.Info() == nil {
		return fmt.Errorf("torrent metadata not loaded yet")
	}

	e.mu.Lock()
	if e.verifying == nil {
		e.verifying = make(map[string]struct{})
	}
	if _, already := e.verifying[ih]; already {
		e.mu.Unlock()
		return fmt.Errorf("torrent data verification already in progress")
	}
	e.verifying[ih] = struct{}{}
	e.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-e.stopCh:
			cancel()
		case <-t.Closed():
			cancel()
		case <-ctx.Done():
		}
	}()

	go func() {
		defer cancel()
		defer func() {
			e.mu.Lock()
			delete(e.verifying, ih)
			e.mu.Unlock()
		}()
		_ = t.VerifyDataContext(ctx)
	}()

	return nil
}

func (e *Engine) SetFilePriority(infoHash string, fileIndex int, priority int) error {
	t, ok := e.findTorrent(infoHash)
	if !ok {
		return fmt.Errorf("torrent not found")
	}

	if t.Info() == nil {
		return fmt.Errorf("torrent metadata not loaded yet")
	}

	files := t.Files()
	if fileIndex < 0 || fileIndex >= len(files) {
		return fmt.Errorf("file index out of range")
	}

	saved, found := e.sessionStore.Get(infoHash)
	if !found {
		saved = &SavedTorrent{
			InfoHash: strings.ToLower(infoHash),
			Name:     t.Name(),
			SavePath: e.cfgManager.Get().DownloadDir,
			AddedAt:  time.Now(),
		}
	}
	if saved.FilePriorities == nil {
		saved.FilePriorities = make(map[int]int)
	}
	saved.FilePriorities[fileIndex] = priority

	f := files[fileIndex]
	switch priority {
	case 0:
		f.SetPriority(types.PiecePriorityNone)
		t.CancelPieces(f.BeginPieceIndex(), f.EndPieceIndex())
		// Ensure any other wanted files don't lose shared boundary pieces
		for idx, otherFile := range files {
			if idx != fileIndex {
				otherPrio := 1
				if p, ok := saved.FilePriorities[idx]; ok {
					otherPrio = p
				}
				if otherPrio > 0 {
					t.DownloadPieces(otherFile.BeginPieceIndex(), otherFile.EndPieceIndex())
				}
			}
		}
	case 1:
		f.SetPriority(types.PiecePriorityNormal)
		t.DownloadPieces(f.BeginPieceIndex(), f.EndPieceIndex())
	case 2:
		f.SetPriority(types.PiecePriorityHigh)
		t.DownloadPieces(f.BeginPieceIndex(), f.EndPieceIndex())
	default:
		f.SetPriority(types.PiecePriorityNormal)
		t.DownloadPieces(f.BeginPieceIndex(), f.EndPieceIndex())
	}

	_ = e.sessionStore.Put(saved)
	if saved.Sequential {
		e.applySequential(t)
	}

	return nil
}

func (e *Engine) ApplyConfig(cfg config.AppConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.client != nil {
		for _, t := range e.client.Torrents() {
			if cfg.MaxEstablished > 0 {
				t.SetMaxEstablishedConns(cfg.MaxEstablished)
			}
		}
	}

	if e.downloadLimiter != nil {
		if cfg.DownloadLimitBps > 0 {
			e.downloadLimiter.SetLimit(rate.Limit(cfg.DownloadLimitBps))
			e.downloadLimiter.SetBurst(max(int(cfg.DownloadLimitBps), 1<<20))
		} else {
			e.downloadLimiter.SetLimit(rate.Inf)
		}
	}

	if e.uploadLimiter != nil {
		if cfg.UploadLimitBps > 0 {
			e.uploadLimiter.SetLimit(rate.Limit(cfg.UploadLimitBps))
			e.uploadLimiter.SetBurst(max(int(cfg.UploadLimitBps), 1<<20))
		} else {
			e.uploadLimiter.SetLimit(rate.Inf)
		}
	}
}

func (e *Engine) GetRateLimits() (downloadBps int64, uploadBps int64) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.downloadLimiter != nil {
		lim := e.downloadLimiter.Limit()
		if lim != rate.Inf {
			downloadBps = int64(lim)
		}
	}
	if e.uploadLimiter != nil {
		lim := e.uploadLimiter.Limit()
		if lim != rate.Inf {
			uploadBps = int64(lim)
		}
	}
	return downloadBps, uploadBps
}

func (e *Engine) GetTorrents() []TorrentSummary {
	torrents := e.client.Torrents()
	list := make([]TorrentSummary, 0, len(torrents))
	for _, t := range torrents {
		list = append(list, e.getTorrentSummary(t))
	}
	return list
}

func (e *Engine) getTorrentSummary(t *torrent.Torrent) TorrentSummary {
	ih := strings.ToLower(t.InfoHash().HexString())
	name := t.Name()
	saved, _ := e.sessionStore.Get(ih)
	isPaused := false
	seq := false
	addedAt := time.Now()
	savePath := e.cfgManager.Get().DownloadDir

	if saved != nil {
		isPaused = saved.IsPaused
		seq = saved.Sequential
		addedAt = saved.AddedAt
		if saved.SavePath != "" {
			savePath = saved.SavePath
		}
		if name == "" && saved.Name != "" {
			name = saved.Name
		}
	}

	if name == "" {
		name = "Downloading metadata..."
	}

	totalLength := t.Length()
	bytesCompleted := t.BytesCompleted()

	var wantedLength int64
	var wantedCompleted int64
	hasUnwanted := false
	if t.Info() != nil {
		var priorities map[int]int
		if saved != nil {
			priorities = saved.FilePriorities
		}
		for idx, f := range t.Files() {
			prio := 1
			if priorities != nil {
				if p, ok := priorities[idx]; ok {
					prio = p
				}
			}
			if prio != 0 {
				wantedLength += f.Length()
				wantedCompleted += f.BytesCompleted()
			} else {
				hasUnwanted = true
			}
		}
	} else {
		wantedLength = totalLength
		wantedCompleted = bytesCompleted
	}

	effectiveTotal := totalLength
	effectiveCompleted := bytesCompleted
	if hasUnwanted && wantedLength >= 0 {
		effectiveTotal = wantedLength
		effectiveCompleted = wantedCompleted
	}

	var progress float64
	if effectiveTotal > 0 {
		progress = (float64(effectiveCompleted) / float64(effectiveTotal)) * 100.0
		if progress > 100.0 {
			progress = 100.0
		}
	} else if t.Info() != nil && totalLength == 0 {
		progress = 100.0
	}

	stats := t.Stats()
	uploaded := stats.BytesWrittenData.Int64()

	var ratio float64
	if bytesCompleted > 0 {
		ratio = float64(uploaded) / float64(bytesCompleted)
	}

	e.mu.RLock()
	st := e.speedTrack[ih]
	var dlSpeed, ulSpeed int64
	if st != nil {
		dlSpeed = st.downloadSpeed
		ulSpeed = st.uploadSpeed
	}
	e.mu.RUnlock()

	var eta int64 = -1
	var missing int64
	if t.Info() != nil {
		if hasUnwanted {
			missing = wantedLength - wantedCompleted
			if missing < 0 {
				missing = 0
			}
		} else {
			missing = t.BytesMissing()
			if missing < 0 {
				missing = 0
			}
		}
	} else {
		missing = totalLength - bytesCompleted
		if missing < 0 {
			missing = 0
		}
	}
	if dlSpeed > 0 && missing > 0 {
		eta = missing / dlSpeed
	}

	e.mu.RLock()
	isVerifying := false
	if e.verifying != nil {
		_, isVerifying = e.verifying[ih]
	}
	e.mu.RUnlock()

	status := "downloading"
	isDone := (t.Info() != nil && totalLength == 0) ||
		(totalLength > 0 && !hasUnwanted && (bytesCompleted >= totalLength || (t.Info() != nil && missing == 0))) ||
		(hasUnwanted && wantedLength > 0 && wantedCompleted >= wantedLength) ||
		(t.Seeding() && bytesCompleted > 0)

	if isVerifying {
		status = "checking"
	} else if isPaused {
		status = "paused"
	} else if isDone {
		status = "seeding"
		eta = 0
		progress = 100.0
		if bytesCompleted > 0 && !e.checkFilesExist(t, savePath) {
			status = "missing"
		}
	} else if t.Info() == nil {
		status = "metadata"
	}

	numFiles := 0
	if t.Info() != nil {
		numFiles = len(t.Files())
	}

	return TorrentSummary{
		InfoHash:       ih,
		Name:           name,
		Size:           totalLength,
		Downloaded:     bytesCompleted,
		Uploaded:       uploaded,
		Progress:       progress,
		Status:         status,
		DownloadSpeed:  dlSpeed,
		UploadSpeed:    ulSpeed,
		ETA:            eta,
		PeersConnected: stats.ActivePeers,
		PeersTotal:     stats.TotalPeers,
		SeedsConnected: stats.ConnectedSeeders,
		Ratio:          ratio,
		AddedAt:        addedAt,
		SavePath:       savePath,
		Sequential:     seq,
		NumFiles:       numFiles,
	}
}

func (e *Engine) GetTorrentDetail(infoHash string) (*TorrentDetail, error) {
	t, ok := e.findTorrent(infoHash)
	if !ok {
		return nil, fmt.Errorf("torrent not found")
	}

	summary := e.getTorrentSummary(t)
	saved, _ := e.sessionStore.Get(strings.ToLower(t.InfoHash().HexString()))

	var pieceLen int64
	var pieceCount int
	pieceMap := make([]byte, 0)
	filesList := make([]FileEntry, 0)
	var comment, createdBy string
	var creationDate int64

	if info := t.Info(); info != nil {
		pieceLen = info.PieceLength
		pieceCount = t.NumPieces()

		pieceMap = make([]byte, pieceCount)
		for i := 0; i < pieceCount; i++ {
			ps := t.PieceState(i)
			if ps.Complete {
				pieceMap[i] = 2
			} else if ps.Partial {
				pieceMap[i] = 1
			} else {
				pieceMap[i] = 0
			}
		}

		// Files
		var priorities map[int]int
		if saved != nil {
			priorities = saved.FilePriorities
		}
		for idx, f := range t.Files() {
			var prog float64
			fLen := f.Length()
			fComp := f.BytesCompleted()
			if fLen > 0 {
				prog = (float64(fComp) / float64(fLen)) * 100.0
			}

			prio := 1
			if priorities != nil {
				if p, ok := priorities[idx]; ok {
					prio = p
				}
			}

			filesList = append(filesList, FileEntry{
				Index:          idx,
				Path:           f.DisplayPath(),
				Length:         fLen,
				BytesCompleted: fComp,
				Progress:       prog,
				Priority:       prio,
			})
		}
	}

	mi := t.Metainfo()
	comment = mi.Comment
	createdBy = mi.CreatedBy
	creationDate = mi.CreationDate

	// Connected peers
	peersList := make([]PeerEntry, 0)
	for _, pc := range t.PeerConns() {
		clientName := "Unknown"
		if val := pc.PeerClientName.Load(); val != nil {
			clientName = fmt.Sprintf("%v", val)
		}

		pStats := pc.Stats()
		var peerProgress float64
		if pieceCount > 0 && pStats.RemotePieceCount > 0 {
			peerProgress = (float64(pStats.RemotePieceCount) / float64(pieceCount)) * 100.0
			if peerProgress > 100.0 {
				peerProgress = 100.0
			}
		}

		peersList = append(peersList, PeerEntry{
			Address:      pc.RemoteAddr.String(),
			Client:       clientName,
			DownloadRate: pc.DownloadRate(),
			UploadRate:   pStats.LastWriteUploadRate,
			Progress:     peerProgress,
		})
	}

	// Trackers
	trackersList := make([]TrackerEntry, 0)
	for tierIdx, tier := range mi.UpvertedAnnounceList() {
		for _, u := range tier {
			trackersList = append(trackersList, TrackerEntry{
				URL:     u,
				Status:  "Working",
				Tier:    tierIdx,
				Message: "",
			})
		}
	}

	return &TorrentDetail{
		Summary:      summary,
		PieceLength:  pieceLen,
		PieceCount:   pieceCount,
		PieceMap:     pieceMap,
		Comment:      comment,
		CreatedBy:    createdBy,
		CreationDate: creationDate,
		Files:        filesList,
		Peers:        peersList,
		Trackers:     trackersList,
	}, nil
}

func (e *Engine) GetGlobalStats() GlobalStats {
	var totalDlSpeed, totalUlSpeed, totalDl, totalUl int64
	active := 0

	torrents := e.client.Torrents()

	e.mu.RLock()
	speeds := make(map[string]struct{ dl, ul int64 }, len(e.speedTrack))
	for ih, st := range e.speedTrack {
		if st != nil {
			speeds[ih] = struct{ dl, ul int64 }{st.downloadSpeed, st.uploadSpeed}
		}
	}
	e.mu.RUnlock()

	for _, t := range torrents {
		ih := strings.ToLower(t.InfoHash().HexString())
		sp, ok := speeds[ih]
		if ok {
			totalDlSpeed += sp.dl
			totalUlSpeed += sp.ul
		}
		stats := t.Stats()
		totalDl += t.BytesCompleted()
		totalUl += stats.BytesWrittenData.Int64()
		if stats.ActivePeers > 0 || (ok && (sp.dl > 0 || sp.ul > 0)) {
			active++
		}
	}

	dhtCount := 0
	for _, d := range e.client.DhtServers() {
		if d != nil {
			if st, ok := d.Stats().(dht.ServerStats); ok {
				if st.GoodNodes > 0 {
					dhtCount += st.GoodNodes
				} else if st.Nodes > 0 {
					dhtCount += st.Nodes
				}
			}
		}
	}

	return GlobalStats{
		TotalDownloadSpeed: totalDlSpeed,
		TotalUploadSpeed:   totalUlSpeed,
		TotalDownloaded:    totalDl,
		TotalUploaded:      totalUl,
		ActiveTorrents:     active,
		TotalTorrents:      len(torrents),
		DHTNodes:           dhtCount,
	}
}

func (e *Engine) findTorrent(infoHash string) (*torrent.Torrent, bool) {
	ihLower := strings.ToLower(infoHash)
	for _, t := range e.client.Torrents() {
		if strings.ToLower(t.InfoHash().HexString()) == ihLower {
			return t, true
		}
	}
	return nil, false
}

func (e *Engine) Close() {
	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		return
	}
	e.closed = true
	close(e.stopCh)
	e.mu.Unlock()

	e.client.Close()
	if e.pieceStore != nil {
		_ = e.pieceStore.Close()
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func copyFile(src, dst string) error {
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)
	if strings.EqualFold(cleanSrc, cleanDst) {
		return nil
	}

	in, err := os.Open(cleanSrc)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(cleanDst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func (e *Engine) checkFilesExist(t *torrent.Torrent, savePath string) bool {
	if t == nil || savePath == "" || t.Info() == nil {
		return true
	}
	name := t.Name()
	if name == "" || name == "Downloading metadata..." {
		return true
	}

	target := filepath.Join(savePath, filepath.FromSlash(name))
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}

	if info.IsDir() && t.Info() != nil {
		files := t.Files()
		if len(files) > 0 {
			for _, f := range files {
				if f.BytesCompleted() > 0 {
					filePath := filepath.Join(savePath, filepath.FromSlash(f.Path()))
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						return false
					}
				}
			}
		}
	}
	return true
}
