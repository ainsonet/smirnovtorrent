package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSessionStoreFilePriorities(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storePath := filepath.Join(tempDir, "session.json")
	torrentDir := filepath.Join(tempDir, "torrents")
	_ = os.MkdirAll(torrentDir, 0755)

	store := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}

	st := &SavedTorrent{
		InfoHash:       "0123456789abcdef",
		Name:           "Test Torrent",
		AddedAt:        time.Now(),
		SavePath:       tempDir,
		IsPaused:       false,
		FilePriorities: map[int]int{0: 0, 1: 2, 2: 1},
	}

	if err := store.Put(st); err != nil {
		t.Fatalf("Failed to put saved torrent: %v", err)
	}

	// Create new session store pointing to same file and test load
	loadedStore := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}
	loadedStore.load()

	item, ok := loadedStore.Get("0123456789abcdef")
	if !ok {
		t.Fatalf("Expected torrent to be found after reload")
	}

	if len(item.FilePriorities) != 3 {
		t.Fatalf("Expected 3 file priorities, got %d", len(item.FilePriorities))
	}

	if item.FilePriorities[0] != 0 || item.FilePriorities[1] != 2 || item.FilePriorities[2] != 1 {
		t.Fatalf("File priorities did not match: got %v", item.FilePriorities)
	}
}

func TestSessionStoreRemove(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_remove_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storePath := filepath.Join(tempDir, "session.json")
	torrentDir := filepath.Join(tempDir, "torrents")
	_ = os.MkdirAll(torrentDir, 0755)

	store := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}

	ih := "aabbccddeeff"
	dummyTor := filepath.Join(torrentDir, ih+".torrent")
	_ = os.WriteFile(dummyTor, []byte("dummy"), 0644)

	store.items[ih] = &SavedTorrent{
		InfoHash: ih,
		Name:     "To Remove",
	}
	_ = store.Save()

	if err := store.Remove(ih); err != nil {
		t.Fatalf("Failed to remove torrent: %v", err)
	}

	if _, ok := store.Get(ih); ok {
		t.Fatalf("Expected torrent to be removed from store")
	}

	if _, err := os.Stat(dummyTor); !os.IsNotExist(err) {
		t.Fatalf("Expected saved .torrent file to be deleted")
	}
}

func TestSessionStoreCaseInsensitivity(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_case_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storePath := filepath.Join(tempDir, "session.json")
	torrentDir := filepath.Join(tempDir, "torrents")
	_ = os.MkdirAll(torrentDir, 0755)

	store := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}

	upperHash := "4A1B2C3D4E5F6A7B8C9D"
	store.Put(&SavedTorrent{
		InfoHash: upperHash,
		Name:     "Case Test",
	})

	// Get with lowercase
	lowerHash := "4a1b2c3d4e5f6a7b8c9d"
	item, ok := store.Get(lowerHash)
	if !ok || item == nil {
		t.Fatalf("Expected torrent to be found using lowercase hash")
	}
	if item.InfoHash != lowerHash {
		t.Errorf("Expected normalized lowercase InfoHash, got %s", item.InfoHash)
	}

	// Remove with upperHash
	if err := store.Remove(upperHash); err != nil {
		t.Fatalf("Failed to remove using uppercase hash: %v", err)
	}
	if _, ok := store.Get(lowerHash); ok {
		t.Fatalf("Expected torrent to be removed")
	}
}

func TestSessionStoreDeepCloneIsolation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_clone_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storePath := filepath.Join(tempDir, "session.json")
	torrentDir := filepath.Join(tempDir, "torrents")
	_ = os.MkdirAll(torrentDir, 0755)

	store := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}

	ih := "1122334455"
	orig := &SavedTorrent{
		InfoHash:       ih,
		Name:           "Original",
		FilePriorities: map[int]int{0: 1, 1: 2},
	}
	_ = store.Put(orig)

	// Mutating the original struct or map should not affect the store
	orig.Name = "Mutated Externally"
	orig.FilePriorities[0] = 0

	stored, ok := store.Get(ih)
	if !ok {
		t.Fatalf("Expected stored item to exist")
	}
	if stored.Name != "Original" {
		t.Errorf("Store item name was mutated externally: %s", stored.Name)
	}
	if stored.FilePriorities[0] != 1 {
		t.Errorf("Store item file priorities were mutated externally: %v", stored.FilePriorities)
	}

	// Mutating the retrieved struct should not affect store either
	stored.FilePriorities[0] = 99
	stored2, _ := store.Get(ih)
	if stored2.FilePriorities[0] != 1 {
		t.Errorf("Store item file priorities mutated via Get: %v", stored2.FilePriorities)
	}
}

func TestSessionStoreConcurrentUpdates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_conc_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storePath := filepath.Join(tempDir, "session.json")
	torrentDir := filepath.Join(tempDir, "torrents")
	_ = os.MkdirAll(torrentDir, 0755)

	store := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ih := fmt.Sprintf("hash%04d", idx)
			_ = store.Put(&SavedTorrent{
				InfoHash: ih,
				Name:     fmt.Sprintf("Torrent %d", idx),
			})
			_, _ = store.Get(ih)
			_ = store.All()
		}(i)
	}
	wg.Wait()

	items := store.All()
	if len(items) != 20 {
		t.Errorf("Expected 20 items in store, got %d", len(items))
	}

	// Verify reload from disk
	reloaded := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}
	reloaded.load()
	if len(reloaded.All()) != 20 {
		t.Errorf("Expected 20 items loaded from disk, got %d", len(reloaded.All()))
	}
}

func TestSessionStoreRemoveDoesNotDeleteExternalFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_ext_del_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outsideDir, err := os.MkdirTemp("", "session_user_dir_*")
	if err != nil {
		t.Fatalf("Failed to create outside dir: %v", err)
	}
	defer os.RemoveAll(outsideDir)

	userTorrentFile := filepath.Join(outsideDir, "important_user_torrent.torrent")
	_ = os.WriteFile(userTorrentFile, []byte("user data"), 0644)

	storePath := filepath.Join(tempDir, "session.json")
	torrentDir := filepath.Join(tempDir, "torrents")
	_ = os.MkdirAll(torrentDir, 0755)

	store := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}

	ih := "safe998877665544332211"
	store.items[ih] = &SavedTorrent{
		InfoHash:    ih,
		Name:        "Safe Test",
		TorrentFile: userTorrentFile, // User's file outside torrentDir
	}
	_ = store.Save()

	// Removing torrent from session store must NOT delete user's external file!
	if err := store.Remove(ih); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if _, err := os.Stat(userTorrentFile); os.IsNotExist(err) {
		t.Fatalf("CRITICAL: User's external .torrent file was deleted by SessionStore.Remove!")
	}
}
