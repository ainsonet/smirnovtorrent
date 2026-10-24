package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"smirnovtorrent/config"
	"smirnovtorrent/engine"
	"smirnovtorrent/sysintegration"
)

//go:embed assets/*
var assetsFS embed.FS

func init() {
	_ = mime.AddExtensionType(".woff2", "font/woff2")
}

// GetLogoDark returns embedded dark logo asset bytes.
func GetLogoDark() []byte {
	b, _ := assetsFS.ReadFile("assets/smirnovtorrent.jpg")
	return b
}

// GetLogoLight returns embedded light logo asset bytes.
func GetLogoLight() []byte {
	b, _ := assetsFS.ReadFile("assets/smirnovtorrent-light.jpg")
	return b
}

type Server struct {
	mu            sync.RWMutex
	eng           *engine.Engine
	cfgManager    *config.ConfigManager
	listener      net.Listener
	port          int
	onWindowShow  func()
	onThemeChange func(string)
}

func NewServer(eng *engine.Engine, cfgManager *config.ConfigManager) (*Server, error) {
	// Listen on loopback on random available port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to bind loopback listener: %w", err)
	}

	port := ln.Addr().(*net.TCPAddr).Port

	s := &Server{
		eng:        eng,
		cfgManager: cfgManager,
		listener:   ln,
		port:       port,
	}

	return s, nil
}

func (s *Server) SetOnWindowShow(fn func()) {
	s.mu.Lock()
	s.onWindowShow = fn
	s.mu.Unlock()
}

func (s *Server) SetOnThemeChange(fn func(string)) {
	s.mu.Lock()
	s.onThemeChange = fn
	s.mu.Unlock()
}

func (s *Server) GetURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d", s.port)
}

func (s *Server) GetPort() int {
	return s.port
}

func (s *Server) Close() error {
	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		return err
	}
	return nil
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API Handlers
	mux.HandleFunc("/api/torrents", s.handleGetTorrents)
	mux.HandleFunc("/api/torrent", s.handleGetTorrentDetail)
	mux.HandleFunc("/api/stats", s.handleGetStats)
	mux.HandleFunc("/api/add-magnet", s.handleAddMagnet)
	mux.HandleFunc("/api/add-file-dialog", s.handleAddFileDialog)
	mux.HandleFunc("/api/upload-torrent", s.handleUploadTorrent)
	mux.HandleFunc("/api/torrent/pause", s.handlePause)
	mux.HandleFunc("/api/torrent/resume", s.handleResume)
	mux.HandleFunc("/api/torrent/delete", s.handleDelete)
	mux.HandleFunc("/api/torrent/sequential", s.handleSequential)
	mux.HandleFunc("/api/torrent/verify", s.handleVerify)
	mux.HandleFunc("/api/torrent/redownload", s.handleRedownload)
	mux.HandleFunc("/api/torrent/file-priority", s.handleFilePriority)
	mux.HandleFunc("/api/torrent/open-file", s.handleOpenFile)
	mux.HandleFunc("/api/torrent/show-file-in-folder", s.handleShowFileInFolder)
	mux.HandleFunc("/api/open-folder", s.handleOpenFolder)
	mux.HandleFunc("/api/browse-dir", s.handleBrowseDir)
	mux.HandleFunc("/api/settings", s.handleSettings)
	mux.HandleFunc("/api/system/integration", s.handleSystemIntegration)
	mux.HandleFunc("/api/add-torrent-path", s.handleAddTorrentPath)
	mux.HandleFunc("/api/window/show", s.handleWindowShow)
	mux.HandleFunc("/api/window/theme", s.handleWindowTheme)

	// Embedded Static Assets
	strippedFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(strippedFS))
	mux.Handle("/", fileServer)

	server := &http.Server{
		Handler: mux,
	}

	return server.Serve(s.listener)
}

func (s *Server) handleGetTorrents(w http.ResponseWriter, r *http.Request) {
	torrents := s.eng.GetTorrents()
	writeJSON(w, torrents)
}

func (s *Server) handleGetTorrentDetail(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "missing hash", http.StatusBadRequest)
		return
	}

	detail, err := s.eng.GetTorrentDetail(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, detail)
}

func (s *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
	stats := s.eng.GetGlobalStats()
	writeJSON(w, stats)
}

type AddMagnetRequest struct {
	MagnetURI string `json:"magnet_uri"`
}

func (s *Server) handleAddMagnet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddMagnetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "invalid json"})
		return
	}

	summary, err := s.eng.AddMagnet(req.MagnetURI)
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"success": true, "torrent": summary})
}

func (s *Server) handleAddFileDialog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath, err := PickTorrentFile()
	if err != nil || filePath == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "cancelled"})
		return
	}

	summary, err := s.eng.AddTorrentFile(filePath)
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"success": true, "torrent": summary})
}

func (s *Server) handleUploadTorrent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 MB max
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "file too large or invalid multipart"})
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing file in request"})
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "uploaded-*.torrent")
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	if _, err := io.Copy(tempFile, file); err != nil {
		_ = tempFile.Close()
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	_ = tempFile.Close()

	summary, err := s.eng.AddTorrentFile(tempPath)
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"success": true, "torrent": summary})
}

func (s *Server) handlePause(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash parameter"})
		return
	}
	err := s.eng.PauseTorrent(hash)
	writeResult(w, err)
}

func (s *Server) handleResume(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash parameter"})
		return
	}
	err := s.eng.ResumeTorrent(hash)
	writeResult(w, err)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash parameter"})
		return
	}
	deleteFiles := r.URL.Query().Get("deleteFiles") == "true"
	err := s.eng.DeleteTorrent(hash, deleteFiles)
	writeResult(w, err)
}

func (s *Server) handleSequential(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash parameter"})
		return
	}
	enabled := r.URL.Query().Get("enabled") == "true"
	err := s.eng.SetSequential(hash, enabled)
	writeResult(w, err)
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash parameter"})
		return
	}
	err := s.eng.VerifyData(hash)
	writeResult(w, err)
}

type FilePriorityReq struct {
	Hash      string `json:"hash"`
	FileIndex int    `json:"file_index"`
	Priority  int    `json:"priority"`
}

func (s *Server) handleFilePriority(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FilePriorityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "invalid json"})
		return
	}

	err := s.eng.SetFilePriority(req.Hash, req.FileIndex, req.Priority)
	writeResult(w, err)
}

func (s *Server) handleOpenFolder(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash parameter"})
		return
	}

	detail, err := s.eng.GetTorrentDetail(hash)
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "torrent not found"})
		return
	}

	savePath := detail.Summary.SavePath
	_ = os.MkdirAll(savePath, 0755)

	target := savePath
	name := detail.Summary.Name
	if name != "" && name != "Downloading metadata..." {
		tPath := filepath.Join(savePath, name)
		if _, err := os.Stat(tPath); err == nil {
			target = tPath
		}
	}

	if err := OpenInExplorer(target); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"success": true})
}

func (s *Server) handleRedownload(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash parameter"})
		return
	}
	err := s.eng.RedownloadTorrent(hash)
	writeResult(w, err)
}

func (s *Server) handleOpenFile(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	indexStr := r.URL.Query().Get("file_index")
	if hash == "" || indexStr == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash or file_index parameter"})
		return
	}
	var index int
	if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "invalid file_index"})
		return
	}

	filePath, err := s.eng.GetTorrentFilePath(hash, index)
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	if _, err := os.Stat(filePath); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "file not found on disk"})
		return
	}

	if err := OpenFile(filePath); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"success": true})
}

func (s *Server) handleShowFileInFolder(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	indexStr := r.URL.Query().Get("file_index")
	if hash == "" || indexStr == "" {
		writeJSON(w, map[string]interface{}{"success": false, "error": "missing hash or file_index parameter"})
		return
	}
	var index int
	if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "invalid file_index"})
		return
	}

	filePath, err := s.eng.GetTorrentFilePath(hash, index)
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	target := filePath
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		target = filepath.Dir(filePath)
		_ = os.MkdirAll(target, 0755)
	}

	if err := OpenInExplorer(target); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"success": true})
}

func (s *Server) handleBrowseDir(w http.ResponseWriter, r *http.Request) {
	current := s.cfgManager.Get().DownloadDir
	folder, err := PickFolder(current)
	if err != nil || folder == "" {
		writeJSON(w, map[string]interface{}{"path": current})
		return
	}
	writeJSON(w, map[string]interface{}{"path": folder})
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cfg := s.cfgManager.Get()
		m := map[string]interface{}{
			"download_dir":         cfg.DownloadDir,
			"default_download_dir": config.GetDefaultDownloadDir(),
			"listen_port":          cfg.ListenPort,
			"enable_dht":           cfg.EnableDHT,
			"enable_pex":           cfg.EnablePEX,
			"enable_upnp":          cfg.EnableUPnP,
			"max_established":      cfg.MaxEstablished,
			"download_limit_bps":   cfg.DownloadLimitBps,
			"upload_limit_bps":     cfg.UploadLimitBps,
			"associate_torrents":   cfg.AssociateTorrents,
			"associate_magnets":    cfg.AssociateMagnets,
			"autostart":            cfg.Autostart,
			"minimize_to_tray":     cfg.MinimizeToTray,
			"close_to_tray":        cfg.CloseToTray,
			"start_minimized":      cfg.StartMinimized,
			"theme":                cfg.Theme,
			"language":             cfg.Language,
		}
		writeJSON(w, m)
		return
	}

	if r.Method == http.MethodPost {
		newCfg := s.cfgManager.Get()
		if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
			writeJSON(w, map[string]interface{}{"success": false, "error": "invalid json"})
			return
		}
		if err := s.cfgManager.Update(newCfg); err != nil {
			writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		s.eng.ApplyConfig(newCfg)
		_ = sysintegration.ApplyIntegration(newCfg)
		s.mu.RLock()
		themeFn := s.onThemeChange
		s.mu.RUnlock()
		if themeFn != nil {
			themeFn(newCfg.Theme)
		}
		writeJSON(w, map[string]interface{}{"success": true})
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleSystemIntegration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		status, err := sysintegration.GetStatus()
		if err != nil {
			writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		writeJSON(w, status)
		return
	}

	if r.Method == http.MethodPost {
		var req sysintegration.IntegrationStatus
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, map[string]interface{}{"success": false, "error": "invalid json"})
			return
		}

		if err := sysintegration.SetTorrentAssociation(req.TorrentAssociated); err != nil {
			writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		if err := sysintegration.SetMagnetAssociation(req.MagnetAssociated); err != nil {
			writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
		if err := sysintegration.SetAutostart(req.AutostartEnabled); err != nil {
			writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}

		cfg := s.cfgManager.Get()
		cfg.AssociateTorrents = req.TorrentAssociated
		cfg.AssociateMagnets = req.MagnetAssociated
		cfg.Autostart = req.AutostartEnabled
		_ = s.cfgManager.Update(cfg)

		writeJSON(w, map[string]interface{}{"success": true})
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

type AddTorrentPathReq struct {
	FilePath string `json:"file_path"`
}

func (s *Server) handleAddTorrentPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddTorrentPathReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "invalid json"})
		return
	}

	summary, err := s.eng.AddTorrentFile(req.FilePath)
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"success": true, "torrent": summary})
}

func (s *Server) handleWindowShow(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	fn := s.onWindowShow
	s.mu.RUnlock()
	if fn != nil {
		fn()
	}
	writeJSON(w, map[string]interface{}{"success": true})
}

type WindowThemeReq struct {
	Theme string `json:"theme"`
}

func (s *Server) handleWindowTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req WindowThemeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": "invalid json"})
		return
	}
	if req.Theme != "light" && req.Theme != "dark" {
		req.Theme = "dark"
	}
	s.mu.RLock()
	fn := s.onThemeChange
	s.mu.RUnlock()
	if fn != nil {
		fn(req.Theme)
	}
	writeJSON(w, map[string]interface{}{"success": true})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeResult(w http.ResponseWriter, err error) {
	if err != nil {
		writeJSON(w, map[string]interface{}{"success": false, "error": err.Error()})
	} else {
		writeJSON(w, map[string]interface{}{"success": true})
	}
}
