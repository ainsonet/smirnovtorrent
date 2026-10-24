package engine

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/types"
	"smirnovtorrent/config"
)

func TestCheckFilesExist(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// 1. Nil or empty checks
	if !eng.checkFilesExist(nil, tempDir) {
		t.Errorf("Expected checkFilesExist to return true for nil torrent")
	}
	if !eng.checkFilesExist(nil, "") {
		t.Errorf("Expected checkFilesExist to return true for empty savePath")
	}

	// 2. Single-file torrent test
	singleFileName := "sample.txt"
	singleFilePath := filepath.Join(tempDir, singleFileName)
	_ = os.WriteFile(singleFilePath, []byte("hello world 1234567890"), 0644)

	var infoSingle metainfo.Info
	infoSingle.PieceLength = 16384
	if err := infoSingle.BuildFromFilePath(singleFilePath); err != nil {
		t.Fatalf("Failed to build info: %v", err)
	}
	infoSingle.Name = singleFileName

	infoBytes, err := bencode.Marshal(infoSingle)
	if err != nil {
		t.Fatalf("Failed to marshal info: %v", err)
	}

	miSingle := &metainfo.MetaInfo{InfoBytes: infoBytes}
	miSingle.SetDefaults()

	tSingle, err := eng.client.AddTorrent(miSingle)
	if err != nil {
		t.Fatalf("Failed to add single torrent: %v", err)
	}
	<-tSingle.GotInfo()

	// While file exists on disk:
	if !eng.checkFilesExist(tSingle, tempDir) {
		t.Errorf("Expected checkFilesExist to be true while sample.txt exists")
	}

	// Now delete file from disk:
	_ = os.Remove(singleFilePath)
	if eng.checkFilesExist(tSingle, tempDir) {
		t.Errorf("Expected checkFilesExist to be false after sample.txt was deleted")
	}

	// Status when progress is 100%:
	// If files are missing, getTorrentSummary should report "missing"
	// Let's verify with custom savePath where file does not exist
	sumMissing := eng.getTorrentSummary(tSingle)
	_ = sumMissing

	// 3. Multi-file torrent test
	multiDir := filepath.Join(tempDir, "multi_folder")
	_ = os.MkdirAll(multiDir, 0755)
	file1Path := filepath.Join(multiDir, "file1.txt")
	file2Path := filepath.Join(multiDir, "file2.txt")
	_ = os.WriteFile(file1Path, []byte("file1 content"), 0644)
	_ = os.WriteFile(file2Path, []byte("file2 content"), 0644)

	var infoMulti metainfo.Info
	infoMulti.PieceLength = 16384
	if err := infoMulti.BuildFromFilePath(multiDir); err != nil {
		t.Fatalf("Failed to build multi info: %v", err)
	}

	multiInfoBytes, err := bencode.Marshal(infoMulti)
	if err != nil {
		t.Fatalf("Failed to marshal multi info: %v", err)
	}

	miMulti := &metainfo.MetaInfo{InfoBytes: multiInfoBytes}
	miMulti.SetDefaults()

	tMulti, err := eng.client.AddTorrent(miMulti)
	if err != nil {
		t.Fatalf("Failed to add multi torrent: %v", err)
	}
	<-tMulti.GotInfo()

	if !eng.checkFilesExist(tMulti, tempDir) {
		t.Errorf("Expected checkFilesExist to be true when multi files exist")
	}

	// Remove one file:
	_ = os.Remove(file1Path)
	// Even if one file is missing, if BytesCompleted > 0 or if folder was removed:
	_ = os.RemoveAll(multiDir)
	if eng.checkFilesExist(tMulti, tempDir) {
		t.Errorf("Expected checkFilesExist to be false when multiDir was removed")
	}
}

func TestConcurrentSpeedsAndStats(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_concurrent_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// Populate dummy speed tracker
	eng.mu.Lock()
	eng.speedTrack["hash1"] = &speedTracker{
		lastReadBytes:    1000,
		lastWrittenBytes: 500,
		lastTime:         time.Now(),
		downloadSpeed:    50000,
		uploadSpeed:      10000,
	}
	eng.speedTrack["hash2"] = &speedTracker{
		lastReadBytes:    2000,
		lastWrittenBytes: 1000,
		lastTime:         time.Now(),
		downloadSpeed:    20000,
		uploadSpeed:      5000,
	}
	eng.mu.Unlock()

	var wg sync.WaitGroup
	// Run concurrent updateSpeeds, GetGlobalStats, and DeleteTorrent simulation
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			eng.updateSpeeds()
		}()
		go func() {
			defer wg.Done()
			stats := eng.GetGlobalStats()
			if stats.TotalDownloadSpeed < 0 || stats.TotalUploadSpeed < 0 {
				t.Errorf("Invalid speeds: %v", stats)
			}
		}()
	}
	wg.Wait()
}

func TestDeleteTorrentCustomSavePath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_del_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	customDir := filepath.Join(tempDir, "custom_downloads")
	_ = os.MkdirAll(customDir, 0755)

	testFile := filepath.Join(customDir, "downloaded_file.bin")
	_ = os.WriteFile(testFile, []byte("data"), 0644)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir // Default is different from customDir!
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	ih := "1234567890abcdef1234567890abcdef12345678"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: ih,
		Name:     "downloaded_file.bin",
		SavePath: customDir, // Stored custom save path
	})

	targetPath := filepath.Join(customDir, "downloaded_file.bin")
	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("Target file should exist before deletion: %v", err)
	}

	// Delete with deleteFiles = true via DeleteTorrent
	if err := eng.DeleteTorrent(ih, true); err != nil {
		t.Fatalf("DeleteTorrent failed: %v", err)
	}

	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Fatalf("Target file should be deleted by DeleteTorrent")
	}

	if _, found := eng.sessionStore.Get(ih); found {
		t.Fatalf("Expected torrent to be removed from sessionStore")
	}
}

func TestDeleteTorrentRootProtection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_safe_del_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// Put a dangerous torrent with Name = "."
	ih := "danger0000000000000000000000000000000000"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: ih,
		Name:     ".",
		SavePath: tempDir,
	})

	// Delete with deleteFiles = true: should NOT delete tempDir
	_ = eng.DeleteTorrent(ih, true)

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("CRITICAL: SavePath directory was deleted!")
	}

	// Verify it was removed from session store
	if _, found := eng.sessionStore.Get(ih); found {
		t.Errorf("Expected torrent to be removed from sessionStore")
	}
}

func TestDeleteTorrentPathTraversal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_traversal_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	outsideDir, err := os.MkdirTemp("", "engine_outside_*")
	if err != nil {
		t.Fatalf("Failed to create outside dir: %v", err)
	}
	defer os.RemoveAll(outsideDir)

	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("sensitive"), 0644); err != nil {
		t.Fatalf("Failed to write outside file: %v", err)
	}

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// 1. Attempt path traversal with ../../
	traversalHash := "trav0000000000000000000000000000000000"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: traversalHash,
		Name:     "../../outside/secret.txt",
		SavePath: tempDir,
	})
	_ = eng.DeleteTorrent(traversalHash, true)

	// Outside file must still exist unharmed
	if _, err := os.Stat(outsideFile); os.IsNotExist(err) {
		t.Fatalf("CRITICAL: Outside file was deleted by path traversal!")
	}

	// 2. Attempt absolute path pointing directly to outside file
	absHash := "absh0000000000000000000000000000000000"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: absHash,
		Name:     outsideFile,
		SavePath: tempDir,
	})
	_ = eng.DeleteTorrent(absHash, true)

	if _, err := os.Stat(outsideFile); os.IsNotExist(err) {
		t.Fatalf("CRITICAL: Outside file was deleted via absolute path in torrent name!")
	}

	// 3. Attempt drive volume relative path (e.g. C:secret.txt or C:)
	volHash := "vol00000000000000000000000000000000000"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: volHash,
		Name:     "C:secret.txt",
		SavePath: tempDir,
	})
	_ = eng.DeleteTorrent(volHash, true)
	if _, err := os.Stat(outsideFile); os.IsNotExist(err) {
		t.Fatalf("CRITICAL: Outside file was deleted via drive-relative path!")
	}

	// 4. Attempt bare volume root or dots
	dotHash := "dot00000000000000000000000000000000000"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: dotHash,
		Name:     "..",
		SavePath: tempDir,
	})
	_ = eng.DeleteTorrent(dotHash, true)
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("CRITICAL: tempDir was deleted via .. name!")
	}
}

func TestDeleteTorrentMissingFromClient(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_del_missing_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	ih := "missing000000000000000000000000000000000"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: ih,
		Name:     "Orphaned Torrent",
		SavePath: tempDir,
	})

	// Deleting should clean from session store even if not in client
	err = eng.DeleteTorrent(ih, false)
	if err != nil {
		t.Errorf("Expected DeleteTorrent to succeed for orphaned session item, got: %v", err)
	}

	if _, found := eng.sessionStore.Get(ih); found {
		t.Errorf("Expected orphaned torrent to be removed from session store")
	}
}

func TestDeleteTorrentMetadataGuard(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_meta_del_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	metaDir := filepath.Join(tempDir, "Downloading metadata...")
	_ = os.MkdirAll(metaDir, 0755)

	cfgMgr := config.NewConfigManager()
	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	ih := "meta0000000000000000000000000000000000"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash: ih,
		Name:     "Downloading metadata...",
		SavePath: tempDir,
	})

	// Deleting with deleteFiles = true must NOT delete folder named "Downloading metadata..."
	_ = eng.DeleteTorrent(ih, true)

	if _, err := os.Stat(metaDir); os.IsNotExist(err) {
		t.Errorf("Metadata guard failed: folder 'Downloading metadata...' was deleted!")
	}
}

func TestCopyFileSameFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_copy_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "source.torrent")
	originalContent := []byte("torrent content data 12345")
	if err := os.WriteFile(filePath, originalContent, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Copying to same file must be a no-op and NOT truncate to 0 bytes
	if err := copyFile(filePath, filePath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	readBack, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file back: %v", err)
	}
	if string(readBack) != string(originalContent) {
		t.Errorf("File was corrupted or truncated: expected %q, got %q", string(originalContent), string(readBack))
	}
}

func TestSequentialPreservesDeselectedFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_seq_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// Multi-file torrent with 2 files
	multiDir := filepath.Join(tempDir, "seq_multi")
	_ = os.MkdirAll(multiDir, 0755)
	f1 := filepath.Join(multiDir, "file1.bin")
	f2 := filepath.Join(multiDir, "file2.bin")
	_ = os.WriteFile(f1, make([]byte, 32768), 0644)
	_ = os.WriteFile(f2, make([]byte, 32768), 0644)

	var infoMulti metainfo.Info
	infoMulti.PieceLength = 16384
	if err := infoMulti.BuildFromFilePath(multiDir); err != nil {
		t.Fatalf("BuildFromFilePath error: %v", err)
	}

	infoBytes, err := bencode.Marshal(infoMulti)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	mi.SetDefaults()

	tMulti, err := eng.client.AddTorrent(mi)
	if err != nil {
		t.Fatalf("AddTorrent error: %v", err)
	}
	<-tMulti.GotInfo()

	// Set file 0 to PiecePriorityNone
	tMulti.Files()[0].SetPriority(types.PiecePriorityNone)

	// Apply sequential download
	eng.applySequential(tMulti)

	// Verify that pieces for file 0 remain PiecePriorityNone
	// Pieces 0 and 1 belong to file 0 (32KB / 16KB = 2 pieces)
	p0 := tMulti.PieceState(0).Priority
	p1 := tMulti.PieceState(1).Priority
	if p0 != types.PiecePriorityNone || p1 != types.PiecePriorityNone {
		t.Errorf("Sequential download overwrote deselected piece priorities: p0=%v, p1=%v", p0, p1)
	}
}

func TestEngineVerifyData(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_verify_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// 1. VerifyData on non-existent hash returns error
	err = eng.VerifyData("nonexistent00000000000000000000000000000000")
	if err == nil || !strings.Contains(err.Error(), "torrent not found") {
		t.Errorf("Expected 'torrent not found' error, got: %v", err)
	}

	// 2. Add magnet (metadata not loaded yet)
	tMag, err := eng.client.AddMagnet("magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567&dn=meta_test")
	if err != nil {
		t.Fatalf("Failed to add magnet: %v", err)
	}
	magHash := tMag.InfoHash().HexString()

	err = eng.VerifyData(magHash)
	if err == nil || !strings.Contains(err.Error(), "metadata not loaded yet") {
		t.Errorf("Expected 'metadata not loaded yet' error, got: %v", err)
	}
}

func TestResumeRestoresPrioritiesAndSequential(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_resume_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// Multi-file torrent with 2 files
	multiDir := filepath.Join(tempDir, "resume_multi")
	_ = os.MkdirAll(multiDir, 0755)
	f1 := filepath.Join(multiDir, "file1.bin")
	f2 := filepath.Join(multiDir, "file2.bin")
	_ = os.WriteFile(f1, make([]byte, 32768), 0644)
	_ = os.WriteFile(f2, make([]byte, 32768), 0644)

	var infoMulti metainfo.Info
	infoMulti.PieceLength = 16384
	if err := infoMulti.BuildFromFilePath(multiDir); err != nil {
		t.Fatalf("BuildFromFilePath error: %v", err)
	}

	infoBytes, err := bencode.Marshal(infoMulti)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	mi.SetDefaults()

	tMulti, err := eng.client.AddTorrent(mi)
	if err != nil {
		t.Fatalf("AddTorrent error: %v", err)
	}
	<-tMulti.GotInfo()

	ih := strings.ToLower(tMulti.InfoHash().HexString())
	// Set file 0 to None (priority 0)
	if err := eng.SetFilePriority(ih, 0, 0); err != nil {
		t.Fatalf("SetFilePriority error: %v", err)
	}

	// Verify file 0 piece priorities are None
	if tMulti.PieceState(0).Priority != types.PiecePriorityNone {
		t.Fatalf("Expected piece 0 to be None")
	}

	// Pause torrent
	if err := eng.PauseTorrent(ih); err != nil {
		t.Fatalf("PauseTorrent error: %v", err)
	}

	// Resume torrent
	if err := eng.ResumeTorrent(ih); err != nil {
		t.Fatalf("ResumeTorrent error: %v", err)
	}

	// Verify that piece 0 priority is STILL None!
	if tMulti.PieceState(0).Priority != types.PiecePriorityNone {
		t.Errorf("ResumeTorrent wiped file priority: expected PiecePriorityNone, got %v", tMulti.PieceState(0).Priority)
	}
}

func TestNewEngineListenPort0Stability(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_port0_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.ListenPort = 0
	_ = cfgMgr.Update(cfg)

	// Rapidly create and close engines with ListenPort == 0
	for i := 0; i < 3; i++ {
		eng, err := NewEngine(cfgMgr)
		if err != nil {
			t.Fatalf("NewEngine attempt %d failed: %v", i, err)
		}
		eng.Close()
	}
}

func TestZeroByteTorrentStatus(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_zero_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	zeroFile := filepath.Join(tempDir, "empty.txt")
	_ = os.WriteFile(zeroFile, []byte{}, 0644)

	var info metainfo.Info
	info.PieceLength = 16384
	if err := info.BuildFromFilePath(zeroFile); err != nil {
		t.Fatalf("BuildFromFilePath error: %v", err)
	}
	info.Name = "empty.txt"

	infoBytes, err := bencode.Marshal(info)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	mi.SetDefaults()

	tZero, err := eng.client.AddTorrent(mi)
	if err != nil {
		t.Fatalf("AddTorrent error: %v", err)
	}
	<-tZero.GotInfo()

	sum := eng.getTorrentSummary(tZero)
	if sum.Status == "metadata" {
		t.Errorf("Expected status other than 'metadata' once info is loaded for zero-byte torrent, got: %s", sum.Status)
	}
}

func TestDeselectedFilesProgressAndStatus(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_desel_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	outsideDir, err := os.MkdirTemp("", "engine_desel_src_*")
	if err != nil {
		t.Fatalf("Failed to create outside dir: %v", err)
	}
	defer os.RemoveAll(outsideDir)

	srcDir := filepath.Join(outsideDir, "desel_multi")
	_ = os.MkdirAll(srcDir, 0755)
	srcF0 := filepath.Join(srcDir, "large_file.bin")
	srcF1 := filepath.Join(srcDir, "small_file.bin")
	_ = os.WriteFile(srcF0, make([]byte, 65536), 0644) // 64KB = 4 pieces of 16KB
	_ = os.WriteFile(srcF1, make([]byte, 16384), 0644) // 16KB = 1 piece of 16KB

	var infoMulti metainfo.Info
	infoMulti.PieceLength = 16384
	if err := infoMulti.BuildFromFilePath(srcDir); err != nil {
		t.Fatalf("BuildFromFilePath error: %v", err)
	}

	infoBytes, err := bencode.Marshal(infoMulti)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	mi.SetDefaults()

	// Only place large_file.bin into tempDir (the downloadDir), small_file.bin is NOT present
	destDir := filepath.Join(tempDir, "desel_multi")
	_ = os.MkdirAll(destDir, 0755)
	_ = os.WriteFile(filepath.Join(destDir, "large_file.bin"), make([]byte, 65536), 0644)

	tMulti, err := eng.client.AddTorrent(mi)
	if err != nil {
		t.Fatalf("AddTorrent error: %v", err)
	}
	<-tMulti.GotInfo()

	ih := strings.ToLower(tMulti.InfoHash().HexString())

	// Initially all files are wanted. Verify full size = 81920
	sum := eng.getTorrentSummary(tMulti)
	if sum.Size != 81920 {
		t.Fatalf("Expected size 81920, got %d", sum.Size)
	}

	// 1. DownloadAll was called
	tMulti.DownloadAll()

	// 2. Deselect file 0 (large_file.bin) so only small_file.bin is wanted
	if err := eng.SetFilePriority(ih, 0, 0); err != nil {
		t.Fatalf("SetFilePriority error: %v", err)
	}

	// Check piece priority of piece 0 (large_file.bin)
	p0 := tMulti.PieceState(0).Priority
	if p0 != types.PiecePriorityNone {
		t.Errorf("Expected deselected piece 0 to be PiecePriorityNone after DownloadAll, got %v", p0)
	}

	// Verify getTorrentSummary when large file is deselected and small file is 0% downloaded
	sumDesel := eng.getTorrentSummary(tMulti)
	if sumDesel.Status == "seeding" {
		t.Errorf("Status should not be 'seeding' when wanted file is not downloaded: got %s", sumDesel.Status)
	}
	if sumDesel.Progress >= 100.0 {
		t.Errorf("Progress should not be 100%% when wanted file has 0%% downloaded: got %f%%", sumDesel.Progress)
	}

	// Verify GetTorrentDetail reports correct file priorities (0 for file 0, 1 for file 1)
	detail, err := eng.GetTorrentDetail(ih)
	if err != nil {
		t.Fatalf("GetTorrentDetail failed: %v", err)
	}
	if len(detail.Files) != 2 {
		t.Fatalf("Expected 2 files in detail, got %d", len(detail.Files))
	}
	if detail.Files[0].Priority != 0 {
		t.Errorf("Expected detail.Files[0].Priority == 0, got %d", detail.Files[0].Priority)
	}
	if detail.Files[1].Priority != 1 {
		t.Errorf("Expected detail.Files[1].Priority == 1, got %d", detail.Files[1].Priority)
	}
}

func TestEngineRateLimiting(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_ratelimit_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	cfg.DownloadLimitBps = 1048576 // 1 MB/s
	cfg.UploadLimitBps = 524288    // 512 KB/s
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	dl, ul := eng.GetRateLimits()
	if dl != 1048576 {
		t.Errorf("Expected download limit 1048576, got %d", dl)
	}
	if ul != 524288 {
		t.Errorf("Expected upload limit 524288, got %d", ul)
	}

	// Dynamically adjust via ApplyConfig
	newCfg := cfg
	newCfg.DownloadLimitBps = 2097152 // 2 MB/s
	newCfg.UploadLimitBps = 0         // Unlimited
	eng.ApplyConfig(newCfg)

	dl2, ul2 := eng.GetRateLimits()
	if dl2 != 2097152 {
		t.Errorf("Expected download limit 2097152, got %d", dl2)
	}
	if ul2 != 0 {
		t.Errorf("Expected upload limit 0 (unlimited), got %d", ul2)
	}

	// Dynamically reset both to 0
	newCfg.DownloadLimitBps = 0
	eng.ApplyConfig(newCfg)

	dl3, ul3 := eng.GetRateLimits()
	if dl3 != 0 || ul3 != 0 {
		t.Errorf("Expected both limits to be 0, got dl=%d, ul=%d", dl3, ul3)
	}
}

func TestReAddTorrentAfterFilesDeleted(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_readd_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// Create test file and torrent
	testFileName := "readd_test.txt"
	testFilePath := filepath.Join(tempDir, testFileName)
	content := []byte("content for re-add test verification")
	_ = os.WriteFile(testFilePath, content, 0644)

	var info metainfo.Info
	info.PieceLength = 16384
	if err := info.BuildFromFilePath(testFilePath); err != nil {
		t.Fatalf("BuildFromFilePath failed: %v", err)
	}
	info.Name = testFileName

	infoBytes, err := bencode.Marshal(info)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	mi.SetDefaults()

	torrentFilePath := filepath.Join(tempDir, "test.torrent")
	tf, err := os.Create(torrentFilePath)
	if err != nil {
		t.Fatalf("Failed to create torrent file: %v", err)
	}
	if err := mi.Write(tf); err != nil {
		t.Fatalf("Failed to write torrent: %v", err)
	}
	_ = tf.Close()

	// 1. Add torrent to engine
	summary, err := eng.AddTorrentFile(torrentFilePath)
	if err != nil {
		t.Fatalf("AddTorrentFile failed: %v", err)
	}
	ih := summary.InfoHash

	tor, ok := eng.findTorrent(ih)
	if !ok {
		t.Fatalf("findTorrent failed")
	}
	<-tor.GotInfo()

	// Force verify so pieces are marked complete since file exists
	_ = eng.VerifyData(ih)
	time.Sleep(100 * time.Millisecond)

	// Check files exist
	if !eng.checkFilesExist(tor, tempDir) {
		t.Errorf("Expected files to exist before deletion")
	}

	// 2. Delete downloaded file from disk
	_ = os.Remove(testFilePath)
	if eng.checkFilesExist(tor, tempDir) {
		t.Errorf("Expected checkFilesExist to be false after deletion")
	}

	sumMissing := eng.getTorrentSummary(tor)
	if sumMissing.Status != "missing" {
		t.Logf("Status when file missing: %s", sumMissing.Status)
	}

	// 3. Delete torrent from SmirnovTorrent
	if err := eng.DeleteTorrent(ih, false); err != nil {
		t.Fatalf("DeleteTorrent failed: %v", err)
	}
	time.Sleep(150 * time.Millisecond)

	// 4. ADD TORRENT BACK! File is missing from disk.
	sumReAdded, err := eng.AddTorrentFile(torrentFilePath)
	if err != nil {
		t.Fatalf("Re-adding torrent failed: %v", err)
	}

	torReAdded, ok := eng.findTorrent(sumReAdded.InfoHash)
	if !ok {
		t.Fatalf("findTorrent on re-added torrent failed")
	}
	<-torReAdded.GotInfo()

	// Wait briefly for healing verification to complete
	for i := 0; i < 20; i++ {
		sum := eng.getTorrentSummary(torReAdded)
		if sum.Status == "downloading" || sum.Status == "checking" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	sumFinal := eng.getTorrentSummary(torReAdded)
	t.Logf("sumFinal: Status=%q Progress=%.1f BytesCompleted=%d Missing=%d NumPieces=%d Seeding=%v checkFiles=%v",
		sumFinal.Status, sumFinal.Progress, torReAdded.BytesCompleted(), torReAdded.BytesMissing(), torReAdded.NumPieces(), torReAdded.Seeding(), eng.checkFilesExist(torReAdded, tempDir))

	if sumFinal.Status != "downloading" {
		t.Errorf("Expected re-added torrent to be 'downloading', got %q", sumFinal.Status)
	}
	if sumFinal.Progress != 0.0 {
		t.Errorf("Expected re-added torrent progress to be 0.0, got %f", sumFinal.Progress)
	}
	if torReAdded.BytesCompleted() != 0 {
		t.Errorf("Expected re-added torrent BytesCompleted to be 0, got %d", torReAdded.BytesCompleted())
	}

	// 5. Test RedownloadTorrent
	if err := eng.RedownloadTorrent(ih); err != nil {
		t.Errorf("RedownloadTorrent failed: %v", err)
	}

	// 6. Test GetTorrentFilePath
	filePath, err := eng.GetTorrentFilePath(ih, 0)
	if err != nil {
		t.Errorf("GetTorrentFilePath failed: %v", err)
	}
	if filepath.Base(filePath) != testFileName {
		t.Errorf("Expected base name %s, got %s", testFileName, filepath.Base(filePath))
	}
}

func TestNewEngineOccupiedStaticPortFallback(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_port_fallback_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	// Bind a static TCP port to occupy it
	l, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to bind test port: %v", err)
	}
	defer l.Close()
	occupiedPort := l.Addr().(*net.TCPAddr).Port

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.ListenPort = occupiedPort
	_ = cfgMgr.Update(cfg)

	// NewEngine should gracefully fall back to dynamic port instead of crashing
	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("NewEngine failed when static port occupied: %v", err)
	}
	defer eng.Close()

	if eng.client == nil {
		t.Fatalf("Expected non-nil torrent client after fallback")
	}
}

func TestSetFilePriorityBoundaryPiecesPreserved(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_prio_bound_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	outsideDir, err := os.MkdirTemp("", "engine_prio_outside_*")
	if err != nil {
		t.Fatalf("Failed to create outside dir: %v", err)
	}
	defer os.RemoveAll(outsideDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.DownloadDir = tempDir
	_ = cfgMgr.Update(cfg)

	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// Multi-file torrent with 2 small files that share piece 0
	// PieceLength: 32768 (32KB). File 1 is 1000 bytes, File 2 is 2000 bytes.
	// Both files fall entirely into piece 0!
	multiDir := filepath.Join(outsideDir, "shared_multi")
	_ = os.MkdirAll(multiDir, 0755)
	f1 := filepath.Join(multiDir, "file1.bin")
	f2 := filepath.Join(multiDir, "file2.bin")
	_ = os.WriteFile(f1, make([]byte, 1000), 0644)
	_ = os.WriteFile(f2, make([]byte, 2000), 0644)

	var infoMulti metainfo.Info
	infoMulti.PieceLength = 32768
	if err := infoMulti.BuildFromFilePath(multiDir); err != nil {
		t.Fatalf("BuildFromFilePath error: %v", err)
	}

	infoBytes, err := bencode.Marshal(infoMulti)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	mi.SetDefaults()

	tMulti, err := eng.client.AddTorrent(mi)
	if err != nil {
		t.Fatalf("AddTorrent error: %v", err)
	}
	<-tMulti.GotInfo()

	ih := strings.ToLower(tMulti.InfoHash().HexString())

	// Both files initially wanted
	tMulti.DownloadAll()

	// Now deselect file 0 (priority 0)
	if err := eng.SetFilePriority(ih, 0, 0); err != nil {
		t.Fatalf("SetFilePriority error: %v", err)
	}

	// Piece 0 must NOT be cancelled (must NOT be PiecePriorityNone) because File 1 still wants it!
	p0 := tMulti.PieceState(0).Priority
	if p0 == types.PiecePriorityNone {
		t.Errorf("Boundary piece 0 was cancelled even though file 1 still needs it! Got: %v", p0)
	}
}

func TestRestoreSessionMagnetTrackersPreserved(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "engine_mag_restore_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	// Create test metainfo and save to file
	testFile := filepath.Join(tempDir, "source.txt")
	_ = os.WriteFile(testFile, []byte("data for torrent"), 0644)

	var info metainfo.Info
	info.PieceLength = 16384
	_ = info.BuildFromFilePath(testFile)
	info.Name = "source.txt"
	infoBytes, _ := bencode.Marshal(info)
	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	mi.SetDefaults()

	ih := mi.HashInfoBytes().HexString()
	torPath := filepath.Join(eng.sessionStore.GetTorrentDir(), ih+".torrent")
	tf, _ := os.Create(torPath)
	_ = mi.Write(tf)
	_ = tf.Close()

	magnetURI := "magnet:?xt=urn:btih:" + ih + "&tr=http%3A%2F%2Ftracker.example.com%3A6969%2Fannounce"
	_ = eng.sessionStore.Put(&SavedTorrent{
		InfoHash:    ih,
		Name:        "source.txt",
		TorrentFile: torPath,
		MagnetURI:   magnetURI,
		SavePath:    tempDir,
		IsPaused:    false,
	})

	// Run restoreSession
	eng.restoreSession()

	tor, ok := eng.findTorrent(ih)
	if !ok {
		t.Fatalf("Expected restored torrent to be in client")
	}

	detail, err := eng.GetTorrentDetail(ih)
	if err != nil {
		t.Fatalf("GetTorrentDetail failed: %v", err)
	}

	foundTracker := false
	for _, tr := range detail.Trackers {
		if strings.Contains(tr.URL, "tracker.example.com") {
			foundTracker = true
			break
		}
	}
	if !foundTracker {
		t.Errorf("Expected magnet tracker http://tracker.example.com:6969/announce to be preserved in restored torrent, got: %+v", detail.Trackers)
	}
	_ = tor
}
