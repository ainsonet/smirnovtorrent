package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"smirnovtorrent/config"
)

type SavedTorrent struct {
	InfoHash       string      `json:"info_hash"`
	Name           string      `json:"name"`
	MagnetURI      string      `json:"magnet_uri,omitempty"`
	TorrentFile    string      `json:"torrent_file,omitempty"`
	AddedAt        time.Time   `json:"added_at"`
	SavePath       string      `json:"save_path"`
	IsPaused       bool        `json:"is_paused"`
	Sequential     bool        `json:"sequential"`
	FilePriorities map[int]int `json:"file_priorities,omitempty"`
}

func cloneSavedTorrent(st *SavedTorrent) *SavedTorrent {
	if st == nil {
		return nil
	}
	cp := *st
	if st.FilePriorities != nil {
		cp.FilePriorities = make(map[int]int, len(st.FilePriorities))
		for k, v := range st.FilePriorities {
			cp.FilePriorities[k] = v
		}
	}
	return &cp
}

type SessionStore struct {
	mu         sync.RWMutex
	fileMu     sync.Mutex
	storePath  string
	torrentDir string
	items      map[string]*SavedTorrent
}

func NewSessionStore() *SessionStore {
	baseDir := config.GetConfigDir()
	torrentDir := filepath.Join(baseDir, "torrents")
	_ = os.MkdirAll(torrentDir, 0755)

	storePath := filepath.Join(baseDir, "session.json")

	ss := &SessionStore{
		storePath:  storePath,
		torrentDir: torrentDir,
		items:      make(map[string]*SavedTorrent),
	}

	ss.load()
	return ss
}

func (ss *SessionStore) load() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	data, err := os.ReadFile(ss.storePath)
	if err != nil {
		return
	}

	var list []*SavedTorrent
	if err := json.Unmarshal(data, &list); err == nil {
		for _, item := range list {
			if item != nil && item.InfoHash != "" {
				ih := strings.ToLower(item.InfoHash)
				item.InfoHash = ih
				ss.items[ih] = cloneSavedTorrent(item)
			}
		}
	}
}

func (ss *SessionStore) Save() error {
	ss.fileMu.Lock()
	defer ss.fileMu.Unlock()

	ss.mu.RLock()
	list := make([]*SavedTorrent, 0, len(ss.items))
	for _, item := range ss.items {
		list = append(list, cloneSavedTorrent(item))
	}
	ss.mu.RUnlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return atomicWriteSessionFile(ss.storePath, data, 0644)
}

func (ss *SessionStore) Put(item *SavedTorrent) error {
	if item == nil {
		return nil
	}
	cp := cloneSavedTorrent(item)
	cp.InfoHash = strings.ToLower(cp.InfoHash)

	ss.mu.Lock()
	ss.items[cp.InfoHash] = cp
	ss.mu.Unlock()

	return ss.Save()
}

func (ss *SessionStore) Get(infoHash string) (*SavedTorrent, bool) {
	ih := strings.ToLower(infoHash)
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	item, ok := ss.items[ih]
	if !ok {
		return nil, false
	}
	return cloneSavedTorrent(item), true
}

func (ss *SessionStore) Remove(infoHash string) error {
	ih := strings.ToLower(infoHash)
	ss.mu.Lock()
	item, ok := ss.items[ih]
	delete(ss.items, ih)
	ss.mu.Unlock()

	// Also delete saved .torrent file if any, but only if it resides within the session torrent directory
	if ok && item != nil && item.TorrentFile != "" {
		cleanTorrentFile := filepath.Clean(item.TorrentFile)
		cleanTorrentDir := filepath.Clean(ss.torrentDir)
		rel, err := filepath.Rel(cleanTorrentDir, cleanTorrentFile)
		if err == nil && !strings.HasPrefix(rel, "..") && rel != "." {
			_ = os.Remove(cleanTorrentFile)
		}
	}
	torrentFile := filepath.Join(ss.torrentDir, ih+".torrent")
	_ = os.Remove(torrentFile)

	return ss.Save()
}

func (ss *SessionStore) All() []*SavedTorrent {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	res := make([]*SavedTorrent, 0, len(ss.items))
	for _, item := range ss.items {
		res = append(res, cloneSavedTorrent(item))
	}
	return res
}

func (ss *SessionStore) GetTorrentDir() string {
	return ss.torrentDir
}

func atomicWriteSessionFile(filename string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	_ = os.MkdirAll(dir, 0755)

	tmpFile, err := os.CreateTemp(dir, "session-tmp-*.json")
	if err != nil {
		return os.WriteFile(filename, data, perm)
	}
	tmpName := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	_ = os.Chmod(tmpName, perm)

	var renameErr error
	for i := 0; i < 5; i++ {
		renameErr = os.Rename(tmpName, filename)
		if renameErr == nil {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	if writeErr := os.WriteFile(filename, data, perm); writeErr == nil {
		return nil
	}
	return renameErr
}
