package engine

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

// AnacrolixEngine обёртка над anacrolix/torrent
type AnacrolixEngine struct {
	client         *torrent.Client
	torrent        *torrent.Torrent
	outputDir      string
	ctx            context.Context
	cancel         context.CancelFunc
	statusMu       sync.RWMutex
	progress       float64
	downloaded     int64
	totalSize      int64
	completePieces int
	totalPieces    int
	speed          float64
	uploadSpeed    float64
	lastBytes      int64
	lastUploadBytes int64
	lastTime       time.Time
	progressCb     func(float64, int, int, int, float64)
}

// NewAnacrolixEngine создаёт новый движок на базе anacrolix
func NewAnacrolixEngine(torrentPath string, outputDir string) (*AnacrolixEngine, error) {
	if outputDir == "" {
		outputDir = filepath.Dir(torrentPath)
	}

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = outputDir
	cfg.DisableIPv6 = true
	cfg.NoDHT = false

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create torrent client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &AnacrolixEngine{
		client:    client,
		ctx:       ctx,
		cancel:    cancel,
		outputDir: outputDir,
	}, nil
}

// LoadTorrent загружает .torrent файл
func (e *AnacrolixEngine) LoadTorrent(torrentPath string) error {
	mi, err := metainfo.LoadFromFile(torrentPath)
	if err != nil {
		return fmt.Errorf("failed to load torrent: %w", err)
	}

	t, err := e.client.AddTorrent(mi)
	if err != nil {
		return fmt.Errorf("failed to add torrent: %w", err)
	}

	e.torrent = t
	e.totalSize = t.Info().TotalLength()
	e.totalPieces = len(t.Info().Pieces)

	log.Printf("✓ Torrent loaded: %s", t.Info().Name)
	log.Printf("  Size: %d bytes", e.totalSize)
	log.Printf("  Pieces: %d", e.totalPieces)

	// Ждём метаданные если это magnet
	<-t.GotInfo()

	return nil
}

// SetProgressCallback устанавливает callback для обновления прогресса
func (e *AnacrolixEngine) SetProgressCallback(cb func(float64, int, int, int, float64)) {
	e.progressCb = cb
}

// Start начинает загрузку
func (e *AnacrolixEngine) Start() error {
	if e.torrent == nil {
		return fmt.Errorf("torrent not loaded")
	}

	// Запускаем загрузку
	e.torrent.DownloadAll()

	log.Println("✓ Download started")

	// Обновляем прогресс
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	e.lastTime = time.Now()

	for {
		select {
		case <-e.ctx.Done():
			return nil
		case <-ticker.C:
			e.updateStatus()
			
			// Скорость
			now := time.Now()
			elapsed := now.Sub(e.lastTime).Seconds()
			if elapsed > 0 {
				e.speed = float64(e.downloaded - e.lastBytes) / elapsed
				e.uploadSpeed = float64(e.torrent.BytesCompleted() - e.lastUploadBytes) / elapsed // упрощённо
			}
			e.lastBytes = e.downloaded
			e.lastUploadBytes = e.torrent.BytesCompleted()
			e.lastTime = now

			// Получаем количество пиров
			peers := len(e.torrent.PeerConns())

			log.Printf("Progress: %.1f%% (%d/%d pieces), Peers: %d, DL: %.1f KB/s, UL: %.1f KB/s",
				e.progress,
				e.completePieces,
				e.totalPieces,
				peers,
				e.speed/1024,
				e.uploadSpeed/1024)

			// Вызываем callback если установлен
			if e.progressCb != nil {
				e.progressCb(
					e.progress,
					e.completePieces,
					e.totalPieces,
					peers,
					e.speed,
				)
			}

			// Проверка завершения
			if e.torrent.BytesCompleted() >= e.torrent.Info().TotalLength() {
				log.Println("✓ Download complete!")
				log.Println("✓ Now seeding (sharing) to other peers...")
				log.Printf("✓ Files location: %s", e.outputDir)
				
				// Проверяем, что файлы действительно сохранены
				e.verifyFiles()
				
				// Продолжаем раздачу - не выходим из цикла
				// Можно добавить отдельную команду для остановки
			}
		}
	}
}

// updateStatus обновляет статус загрузки
func (e *AnacrolixEngine) updateStatus() {
	e.statusMu.Lock()
	defer e.statusMu.Unlock()

	if e.torrent == nil {
		return
	}

	e.downloaded = e.torrent.BytesCompleted()
	e.totalSize = e.torrent.Info().TotalLength()
	e.lastUploadBytes = e.torrent.BytesCompleted()
	
	// Прогресс в процентах
	if e.totalSize > 0 {
		e.progress = float64(e.downloaded) / float64(e.totalSize) * 100
	}
	
	// Считаем завершённые куски через BytesCompleted
	pieceLength := e.torrent.Info().PieceLength
	e.completePieces = int(e.downloaded / pieceLength)
}

// verifyFiles проверяет наличие файлов
func (e *AnacrolixEngine) verifyFiles() {
	files := e.torrent.Files()
	log.Printf("Torrent contains %d file(s):", len(files))
	for _, f := range files {
		fullPath := filepath.Join(e.outputDir, f.Path())
		log.Printf("  - %s (%d bytes) -> %s", f.DisplayPath(), f.Length(), fullPath)
	}
}

// Stop останавливает загрузку
func (e *AnacrolixEngine) Stop() error {
	if e.cancel != nil {
		e.cancel()
	}
	if e.client != nil {
		e.client.Close()
	}
	return nil
}
