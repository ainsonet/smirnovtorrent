package engine

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	tracker      *tracker.Tracker
	cancel       context.CancelFunc
	ctx          context.Context
	status       DownloadStatus
	statusMu     sync.RWMutex
	numWorkers   int
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

	// Генерируем peer ID
	peerID := peer.NewPeerID()
	peerIDStr := string(peerID[:])

	return &DownloadEngine{
		torrent:    torrent,
		outputDir:  outputDir,
		ctx:        nil,
		cancel:     nil,
		numWorkers: 4, // 4 параллельных загрузчика
	}
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
		
		e.cancel()
	}

	return nil
}

// downloadLoop основной цикл загрузки
func (e *DownloadEngine) downloadLoop() error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return nil
		case <-ticker.C:
			e.updateStatus()
			log.Printf("Progress: %.1f%% (%d/%d pieces), Peers: %d",
				e.pieceManager.Progress(),
				e.pieceManager.CompletePieces(),
				e.pieceManager.TotalPieces(),
				e.peerPool.GetActivePeerCount())
		}
	}
}

// downloadLoop основной цикл загрузки
func (e *DownloadEngine) downloadLoop() error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return nil
		case <-ticker.C:
			e.updateStatus()
			log.Printf("Progress: %.1f%% (%d/%d pieces)",
				e.pieceManager.Progress(),
				e.pieceManager.CompletePieces(),
				e.pieceManager.TotalPieces())
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