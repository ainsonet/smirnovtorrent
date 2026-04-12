package engine

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"smirnovtorrent/internal/encryption"
	"smirnovtorrent/internal/parser"
	"smirnovtorrent/internal/peer"
	"smirnovtorrent/internal/tracker"
)

// DownloadEngine основной движок загрузки
type DownloadEngine struct {
	torrent      *parser.Torrent
	outputDir    string
	pieceManager *PieceManager
	peerPool     *PeerPool
	rarestMgr    *RarestFirstManager
	seedManager  *SeedManager
	tracker      *tracker.Tracker
	cancel       context.CancelFunc
	ctx          context.Context
	status       DownloadStatus
	statusMu     sync.RWMutex
	numWorkers   int
	seedMode     bool
	progressCallback func(float64, int, int, int, float64)
	
	// Новые функции
	encryption     *encryption.MSEMessageStreamEncryption
	limiter        *RateLimiter
	resumeManager  *ResumeManager
}

// DownloadStatus состояние загрузки
type DownloadStatus struct {
	Progress     float64
	Downloaded   int64
	TotalSize    int64
	ActivePeers  int
	DownloadSpeed float64 // байт/сек
	StartTime    time.Time
}

// NewDownloadEngine создаёт новый движок загрузки
func NewDownloadEngine(torrent *parser.Torrent, outputDir string) *DownloadEngine {
	if outputDir == "" {
		outputDir = torrent.Info.Name
	}

	return &DownloadEngine{
		torrent:    torrent,
		outputDir:  outputDir,
		ctx:        nil,
		cancel:     nil,
		numWorkers: 4, // 4 параллельных загрузчика
	}
}

// SetProgressCallback устанавливает callback для обновления прогресса
func (e *DownloadEngine) SetProgressCallback(cb func(float64, int, int, int, float64)) {
	e.progressCallback = cb
}

// SetRateLimits устанавливает ограничения скорости
func (e *DownloadEngine) SetRateLimits(downloadRate, uploadRate int64) {
	if e.limiter == nil {
		e.limiter = NewRateLimiter(downloadRate, uploadRate)
	} else {
		e.limiter.SetMaxDownloadRate(downloadRate)
		e.limiter.SetMaxUploadRate(uploadRate)
	}
}

// SetEncryptionKey включает шифрование
func (e *DownloadEngine) SetEncryptionKey(key []byte) {
	e.encryption = encryption.NewMSEEncryption(key)
}

// EnableResume включает продолжение загрузки
func (e *DownloadEngine) EnableResume() {
	e.resumeManager = NewResumeManager(e.torrent.Info.InfoHash, e.outputDir)
	
	// Загружаем сохранённый прогресс
	if err := e.resumeManager.Load(); err != nil {
		log.Printf("Failed to load resume data: %v", err)
		return
	}
	
	// Восстанавливаем завершённые куски
	// В реальной реализации здесь нужно будет восстановить данные кусков
	
	// Запускаем авто-сохранение
	e.resumeManager.StartAutoSave(30 * time.Second)
}

// Start начинает загрузку
func (e *DownloadEngine) Start() error {
	e.ctx, e.cancel = context.WithCancel(context.Background())
	defer e.cancel()

	// Инициализируем PieceManager
	e.pieceManager = NewPieceManager(
		e.torrent.Info.PieceLength,
		e.torrent.TotalSize(),
		e.torrent.Info.Pieces,
	)

	// Создаём пул пиров
	peerID := peer.NewPeerID()
	peerIDStr := string(peerID[:])
	e.peerPool = NewPeerPool(e.torrent.Info.InfoHash, peerIDStr, 6881, 50)

	// Создаём менеджер Rarest-first
	e.rarestMgr = NewRarestFirstManager(e.pieceManager, e.peerPool)

	// Создаём трекер
	if e.torrent.Announce == "" {
		return fmt.Errorf("no tracker URL available")
	}

	var err error
	e.tracker, err = tracker.ParsePeerURL(e.torrent.Announce)
	if err != nil {
		return fmt.Errorf("failed to create tracker: %w", err)
	}

	log.Printf("Starting download: %s", e.torrent.Info.Name)
	log.Printf("Total size: %d bytes", e.torrent.TotalSize())
	log.Printf("Pieces: %d", e.pieceManager.TotalPieces())
	log.Printf("Workers: %d", e.numWorkers)

	// Создаём директорию для загрузки
	if err := os.MkdirAll(e.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Получаем список пиров от трекера
	peers, err := e.tracker.GetPeers(
		e.torrent.Info.InfoHash,
		peerIDStr,
		6881,
	)
	if err != nil {
		log.Printf("Warning: failed to get peers from tracker: %v", err)
	} else {
		log.Printf("Got %d peers from tracker", len(peers))
	}

	// Подключаемся к пирам (до 20 сразу)
	maxInitialPeers := 20
	if len(peers) < maxInitialPeers {
		maxInitialPeers = len(peers)
	}
	
	for i := 0; i < maxInitialPeers; i++ {
		go e.peerPool.AddPeer(peers[i])
	}

	// Ждём немного чтобы пиры подключились
	time.Sleep(2 * time.Second)

	// Запускаем воркеров для параллельной загрузки
	e.startWorkers()

	// Основной цикл загрузки
	return e.downloadLoop()
}

// startWorkers запускает воркеров для параллельной загрузки
func (e *DownloadEngine) startWorkers() {
	for i := 0; i < e.numWorkers; i++ {
		go e.worker(i)
	}
}

// worker воркер который загружает куски
func (e *DownloadEngine) worker(id int) {
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-e.ctx.Done():
			return
		default:
			// Получаем следующий кусок по Rarest-first
			piece := e.rarestMgr.GetNextPieceRarest()
			if piece == nil {
				// Нет доступных кусков, ждём
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Пробуем загрузить кусок
			if err := e.downloadPiece(id, piece); err != nil {
				log.Printf("Worker %d: failed to download piece %d: %v", id, piece.Index, err)
				// Возвращаем кусок как незапрошенный
				piece.Requested = false
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// downloadPiece загружает конкретный кусок
func (e *DownloadEngine) downloadPiece(workerID int, piece *Piece) error {
	// Ищем пирa у которого есть этот кусок
	peers := e.peerPool.GetPeersWithPiece(piece.Index)
	if len(peers) == 0 {
		return fmt.Errorf("no peers with piece %d", piece.Index)
	}

	// Выбираем первого свободного пирa
	var chosenPeer *PeerSlot
	for _, p := range peers {
		if !p.Downloading && !p.Choked {
			chosenPeer = p
			break
		}
	}

	if chosenPeer == nil {
		return fmt.Errorf("all peers with piece %d are busy or choked", piece.Index)
	}

	// Помечаем что качаем
	e.peerPool.MarkPieceDownloading(chosenPeer)
	defer e.peerPool.MarkPieceDone(chosenPeer)

	// Запрашиваем кусок
	pieceLength := e.torrent.Info.PieceLength
	if err := chosenPeer.Conn.SendRequest(uint32(piece.Index), 0, uint32(pieceLength)); err != nil {
		return err
	}

	// Читаем ответ
	msgType, payload, err := chosenPeer.Conn.ReadMessage()
	if err != nil {
		return err
	}

	if msgType != peer.MsgPiece {
		return fmt.Errorf("expected piece message, got %d", msgType)
	}

	// Проверяем и сохраняем кусок
	if err := e.pieceManager.MarkPieceComplete(piece.Index, payload); err != nil {
		return err
	}

	log.Printf("Worker %d: downloaded piece %d/%d (%.1f%%)",
		workerID,
		e.pieceManager.CompletePieces(),
		e.pieceManager.TotalPieces(),
		e.pieceManager.Progress())

	// Проверяем завершена ли загрузка
	if e.pieceManager.IsComplete() {
		log.Println("Download complete!")
		
		// Ждём немного чтобы все куски были сохранены
		time.Sleep(1 * time.Second)
		
		// Собираем файлы
		if err := e.assembleFiles(); err != nil {
			log.Printf("Failed to assemble files: %v", err)
		} else {
			log.Printf("Files saved to: %s", e.outputDir)
		}
		
		// Переходим в seed режим
		go e.startSeedMode()
		
		e.cancel()
		return nil
	}

	return nil
}

// startSeedMode запускает режим раздачи
func (e *DownloadEngine) startSeedMode() {
	e.seedMode = true
	log.Println("Switching to seed mode...")

	// Создаём SeedManager
	e.seedManager = NewSeedManager(e.pieceManager, e.peerPool, 20)

	// Запускаем seed режим
	if err := e.seedManager.StartSeed(); err != nil {
		log.Printf("Seed mode error: %v", err)
	}
}

// downloadLoop основной цикл загрузки
func (e *DownloadEngine) downloadLoop() error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastBytes int64
	var lastUpdate time.Time

	for {
		select {
		case <-e.ctx.Done():
			// Если мы в seed режиме, продолжаем работать
			if e.seedMode {
				log.Println("Download complete, running in seed mode...")
				// Seed режим работает в отдельной горутине
				return nil
			}
			return nil
		case <-ticker.C:
			e.updateStatus()
			
			// Вычисляем скорость
			now := time.Now()
			currentBytes := e.status.Downloaded
			elapsed := now.Sub(lastUpdate).Seconds()
			
			var downloadSpeed float64
			if elapsed > 0 {
				downloadSpeed = float64(currentBytes-lastBytes) / elapsed
			}
			
			lastBytes = currentBytes
			lastUpdate = now

			// Вызываем callback если установлен
			if e.progressCallback != nil {
				e.progressCallback(
					e.status.Progress,
					e.pieceManager.CompletePieces(),
					e.pieceManager.TotalPieces(),
					e.peerPool.GetActivePeerCount(),
					downloadSpeed,
				)
			}

			// Логирование
			log.Printf("Progress: %.1f%% (%d/%d pieces), Peers: %d, Speed: %.1f KB/s",
				e.status.Progress,
				e.pieceManager.CompletePieces(),
				e.pieceManager.TotalPieces(),
				e.peerPool.GetActivePeerCount(),
				downloadSpeed/1024)
		}
	}
}

// updateStatus обновляет статус загрузки
func (e *DownloadEngine) updateStatus() {
	e.statusMu.Lock()
	defer e.statusMu.Unlock()

	e.status.Progress = e.pieceManager.Progress()
	e.status.Downloaded = int64(e.pieceManager.CompletePieces()) * int64(e.torrent.Info.PieceLength)
	e.status.TotalSize = e.torrent.TotalSize()
	e.status.ActivePeers = e.peerPool.GetActivePeerCount()
	e.status.StartTime = time.Now()
}

// GetStatus возвращает текущий статус
func (e *DownloadEngine) GetStatus() DownloadStatus {
	e.statusMu.RLock()
	defer e.statusMu.RUnlock()
	return e.status
}

// Stop останавливает загрузку
func (e *DownloadEngine) Stop() {
	if e.cancel != nil {
		e.cancel()
	}

	if e.peerPool != nil {
		e.peerPool.CloseAll()
	}

	if e.seedManager != nil {
		e.seedManager.Stop()
	}
}

// assembleFiles собирает файлы из кусков
func (e *DownloadEngine) assembleFiles() error {
	if len(e.torrent.Info.Files) == 1 {
		// Single file mode
		data := e.pieceManager.AssembleFile()
		filePath := filepath.Join(e.outputDir, e.torrent.Info.Files[0].Path)
		return os.WriteFile(filePath, data, 0644)
	}

	// Multi-file mode
	for i := range e.torrent.Info.Files {
		data, err := e.pieceManager.GetFileRange(i, e.torrent.Info.Files)
		if err != nil {
			return fmt.Errorf("failed to get file %d: %w", i, err)
		}

		filePath := filepath.Join(e.outputDir, e.torrent.Info.Files[i].Path)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}