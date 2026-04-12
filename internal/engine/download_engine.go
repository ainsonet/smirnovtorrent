package engine

import (
	"context"
	"fmt"
	"log"
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
	torrent     *parser.Torrent
	outputDir   string
	pieceManager *PieceManager
	peers       map[string]*peer.PeerConnection
	peerMutex   sync.RWMutex
	tracker     *tracker.Tracker
	cancel      context.CancelFunc
	ctx         context.Context
	status      DownloadStatus
	statusMu    sync.RWMutex
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
		torrent:      torrent,
		outputDir:    outputDir,
		peers:        make(map[string]*peer.PeerConnection),
		StartTime:    time.Now(),
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

	// Создаём директорию для загрузки
	if err := os.MkdirAll(e.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Получаем список пиров от трекера
	peerID := peer.NewPeerID()
	infoHash := e.torrent.InfoHashBytes()

	peers, err := e.tracker.GetPeers(
		e.torrent.Info.InfoHash,
		string(peerID[:]),
		6881, // стандартный порт
	)
	if err != nil {
		log.Printf("Warning: failed to get peers from tracker: %v", err)
	} else {
		log.Printf("Got %d peers from tracker", len(peers))
	}

	// Подключаемся к пирам
	go e.connectToPeers(peers)

	// Основной цикл загрузки
	return e.downloadLoop()
}

// connectToPeers подключается к пирам
func (e *DownloadEngine) connectToPeers(peers []tracker.PeerInfo) {
	for _, p := range peers {
		select {
		case <-e.ctx.Done():
			return
		default:
			go e.connectToPeer(p)
		}
	}
}

// connectToPeer подключается к одному пиру
func (e *DownloadEngine) connectToPeer(p tracker.PeerInfo) {
	peerObj := &peer.Peer{
		IP:     p.IP,
		Port:   p.Port,
		PeerID: peer.NewPeerID(),
	}

	conn, err := peerObj.Connect()
	if err != nil {
		log.Printf("Failed to connect to %s:%d: %v", p.IP, p.Port, err)
		return
	}

	defer conn.Close()

	// Отправляем handshake
	infoHash := e.torrent.InfoHashBytes()
	peerID := peer.NewPeerID()
	if err := conn.SendHandshake(infoHash, peerID); err != nil {
		log.Printf("Failed to send handshake: %v", err)
		return
	}

	// Читаем handshake от пирa
	remoteInfoHash, remotePeerID, err := conn.ReadHandshake()
	if err != nil {
		log.Printf("Failed to read handshake: %v", err)
		return
	}

	// Проверяем info hash
	if remoteInfoHash != infoHash {
		log.Printf("Info hash mismatch, disconnecting")
		return
	}

	log.Printf("Connected to peer %s", remotePeerID)

	// Сохраняем соединение
	e.peerMutex.Lock()
	e.peers[fmt.Sprintf("%s:%d", p.IP, p.Port)] = conn
	e.peerMutex.Unlock()

	// Читаем bitfield
	bitfield, err := conn.ReadBitfield()
	if err != nil {
		log.Printf("Failed to read bitfield: %v", err)
		return
	}

	log.Printf("Peer has %d pieces", len(bitfield)*8)

	// Сообщаем что мы interested
	if err := conn.SendInterested(); err != nil {
		log.Printf("Failed to send interested: %v", err)
		return
	}

	// Начинаем запрашивать куски
	e.requestPieces(conn)
}

// requestPieces запрашивает куски у пирa
func (e *DownloadEngine) requestPieces(conn *peer.PeerConnection) {
	for {
		select {
		case <-e.ctx.Done():
			return
		default:
			piece := e.pieceManager.GetNextPiece()
			if piece == nil {
				// Все куски запрошены или загружены
				time.Sleep(1 * time.Second)
				continue
			}

			// Запрашиваем кусок
			if err := conn.SendRequest(uint32(piece.Index), 0, uint32(e.torrent.PieceLength)); err != nil {
				log.Printf("Failed to request piece: %v", err)
				return
			}

			// Читаем ответ
			msgType, payload, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Failed to read response: %v", err)
				return
			}

			if msgType == peer.MsgPiece {
				// Проверяем и сохраняем кусок
				if err := e.pieceManager.MarkPieceComplete(piece.Index, payload); err != nil {
					log.Printf("Piece validation failed: %v", err)
					continue
				}

				log.Printf("Downloaded piece %d/%d (%.1f%%)",
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
					return
				}
			}
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

	e.peerMutex.RLock()
	e.status.ActivePeers = len(e.peers)
	e.peerMutex.RUnlock()

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

	e.peerMutex.Lock()
	defer e.peerMutex.Unlock()

	for _, conn := range e.peers {
		conn.Close()
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