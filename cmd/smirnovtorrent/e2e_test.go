//go:build e2e
// +build e2e

package main

import (
	"os"
	"testing"
	"time"

	"smirnovtorrent/internal/engine"
	"smirnovtorrent/internal/parser"
)

// TestE2EDownload тест полного цикла загрузки торрента
// Запуск: go test -tags e2e -v -timeout 5m
func TestE2EDownload(t *testing.T) {
	// Пропускаем если нет тестового файла
	torrentFile := os.Getenv("TORRENT_FILE")
	if torrentFile == "" {
		t.Skip("TORRENT_FILE not set, skipping E2E test")
	}

	// Парсим торрент
	torrent, err := parser.ParseFile(torrentFile)
	if err != nil {
		t.Fatalf("Failed to parse torrent: %v", err)
	}

	t.Logf("Testing download: %s (%d bytes)", torrent.Info.Name, torrent.TotalSize())

	// Создаём временную директорию для загрузки
	tmpDir := "test_download"
	defer os.RemoveAll(tmpDir)

	// Создаём движок
	eng := engine.NewDownloadEngine(torrent, tmpDir)

	// Устанавливаем callback для прогресса
	var lastProgress float64
	eng.SetProgressCallback(func(progress float64, current, total, peers int, speed float64) {
		if progress > lastProgress {
			t.Logf("Progress: %.1f%% (%d/%d pieces), Peers: %d, Speed: %.1f KB/s",
				progress, current, total, peers, speed/1024)
			lastProgress = progress
		}
	})

	// Запускаем загрузку с timeout
	done := make(chan error)
	go func() {
		done <- eng.Start()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Download failed: %v", err)
		}
		t.Log("Download completed successfully!")
	case <-time.After(5 * time.Minute):
		t.Fatal("Download timeout after 5 minutes")
	}

	// Проверяем что файлы созданы
	if len(torrent.Info.Files) > 0 {
		for _, file := range torrent.Info.Files {
			filePath := tmpDir + "/" + file.Path
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected file not found: %s", filePath)
			}
		}
	}
}
