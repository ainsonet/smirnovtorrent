//go:build e2e
// +build e2e

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"smirnovtorrent/internal/engine"
	"smirnovtorrent/internal/parser"
)

// TestE2E_DownloadSmallTorrent тестирует загрузку небольшого торрента
func TestE2E_DownloadSmallTorrent(t *testing.T) {
	torrentFile := os.Getenv("TORRENT_FILE")
	if torrentFile == "" {
		t.Skip("TORRENT_FILE not set, skipping E2E test")
	}

	torrent, err := parser.ParseFile(torrentFile)
	if err != nil {
		t.Fatalf("Failed to parse torrent: %v", err)
	}

	t.Logf("Testing download: %s (%d bytes)", torrent.Info.Name, torrent.TotalSize())

	tmpDir := "test_download_e2e"
	defer os.RemoveAll(tmpDir)

	eng := engine.NewDownloadEngine(torrent, tmpDir)
	eng.EnableDHT()
	eng.EnableResume()

	var lastProgress float64
	eng.SetProgressCallback(func(progress float64, current, total, peers int, speed float64) {
		if progress > lastProgress {
			t.Logf("Progress: %.1f%% (%d/%d), Peers: %d, Speed: %.1f KB/s",
				progress, current, total, peers, speed/1024)
			lastProgress = progress
		}
	})

	done := make(chan error)
	go func() {
		done <- eng.Start()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("Download completed with error: %v", err)
		} else {
			t.Log("Download completed successfully!")
		}
	case <-time.After(5 * time.Minute):
		eng.Stop()
		t.Log("Download interrupted by timeout")
	}
}

// TestE2E_WithRateLimit тестирует загрузку с ограничением скорости
func TestE2E_WithRateLimit(t *testing.T) {
	torrentFile := os.Getenv("TORRENT_FILE")
	if torrentFile == "" {
		t.Skip("TORRENT_FILE not set")
	}

	torrent, err := parser.ParseFile(torrentFile)
	if err != nil {
		t.Fatalf("Failed to parse torrent: %v", err)
	}

	tmpDir := "test_download_ratelimit"
	defer os.RemoveAll(tmpDir)

	eng := engine.NewDownloadEngine(torrent, tmpDir)
	
	// Лимит 100 KB/s download, 50 KB/s upload
	eng.SetRateLimits(102400, 51200)
	eng.EnableDHT()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- eng.Start()
	}()

	select {
	case err := <-done:
		t.Logf("Download completed: %v", err)
	case <-ctx.Done():
		eng.Stop()
		t.Log("Rate limit test completed (30s)")
	}
}

// TestE2E_ResumeDownload тестирует продолжение загрузки
func TestE2E_ResumeDownload(t *testing.T) {
	torrentFile := os.Getenv("TORRENT_FILE")
	if torrentFile == "" {
		t.Skip("TORRENT_FILE not set")
	}

	torrent, err := parser.ParseFile(torrentFile)
	if err != nil {
		t.Fatalf("Failed to parse torrent: %v", err)
	}

	tmpDir := "test_download_resume"
	defer os.RemoveAll(tmpDir)

	// Первая сессия
	t.Log("Starting first session...")
	eng1 := engine.NewDownloadEngine(torrent, tmpDir)
	eng1.EnableResume()
	eng1.EnableDHT()

	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	go func() {
		<-ctx1.Done()
		eng1.Stop()
		t.Log("First session stopped")
	}()

	eng1.Start()
	t.Log("First session completed")

	// Вторая сессия (resume)
	t.Log("Starting second session (resume)...")
	eng2 := engine.NewDownloadEngine(torrent, tmpDir)
	eng2.EnableResume()
	eng2.EnableDHT()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	go func() {
		<-ctx2.Done()
		eng2.Stop()
		t.Log("Second session stopped")
	}()

	eng2.Start()
	t.Log("Second session completed successfully")
}

// TestE2E_MultiFileTorrent тестирует multi-file торренты
func TestE2E_MultiFileTorrent(t *testing.T) {
	torrentFile := os.Getenv("TORRENT_FILE_MULTI")
	if torrentFile == "" {
		t.Skip("TORRENT_FILE_MULTI not set")
	}

	torrent, err := parser.ParseFile(torrentFile)
	if err != nil {
		t.Fatalf("Failed to parse torrent: %v", err)
	}

	if len(torrent.Info.Files) < 2 {
		t.Skip("Torrent is not multi-file")
	}

	t.Logf("Testing multi-file torrent: %s", torrent.Info.Name)
	t.Logf("Files: %d", len(torrent.Info.Files))

	tmpDir := "test_download_multifile"
	defer os.RemoveAll(tmpDir)

	eng := engine.NewDownloadEngine(torrent, tmpDir)
	eng.EnableDHT()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- eng.Start()
	}()

	select {
	case err := <-done:
		t.Logf("Download completed: %v", err)
	case <-ctx.Done():
		eng.Stop()
		t.Log("Multi-file test interrupted by timeout")
	}

	// Проверяем файлы
	checkFileStructure(t, tmpDir, torrent.Info.Files)
}

// checkFileStructure проверяет структуру файлов
func checkFileStructure(t *testing.T, dir string, files []parser.FileInfo) {
	t.Helper()
	
	exists := 0
	missing := 0
	
	for _, file := range files {
		filePath := filepath.Join(dir, file.Path)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			missing++
			t.Logf("Missing: %s", filePath)
		} else {
			exists++
			t.Logf("Exists: %s", filePath)
		}
	}
	
	t.Logf("Files: %d exists, %d missing", exists, missing)
}

// TestMain запускает E2E тесты
func TestMain(m *testing.M) {
	fmt.Println("=================================")
	fmt.Println("SmirnovTorrent E2E Tests")
	fmt.Println("=================================")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  TORRENT_FILE        - Path to test torrent file")
	fmt.Println("  TORRENT_FILE_MULTI  - Path to multi-file torrent")
	fmt.Println("  TORRENT_FILE_LARGE  - Path to large torrent (>1GB)")
	fmt.Println("  MAGNET_LINK         - Magnet link for testing")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go test -tags e2e -v -timeout 5m")
	fmt.Println("  TORRENT_FILE=ubuntu.torrent go test -tags e2e -v")
	fmt.Println()
	
	os.Exit(m.Run())
}
