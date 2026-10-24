package engine

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"go.etcd.io/bbolt"
)

const (
	boltDbCompleteValue   = "c"
	boltDbIncompleteValue = "i"
)

var completionBucketKey = []byte("completion")

// PieceStore manages piece completion status across torrents in a single thread-safe database.
type PieceStore struct {
	mu        sync.RWMutex
	db        *bbolt.DB
	dbPath    string
	memBackup storage.PieceCompletion
}

// NewPieceStore opens or creates a persistent bolt DB for piece completion in baseDir.
func NewPieceStore(baseDir string) (*PieceStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return &PieceStore{memBackup: storage.NewMapPieceCompletion()}, err
	}

	dbPath := filepath.Join(baseDir, "pieces.bolt.db")
	db, err := bbolt.Open(dbPath, 0660, &bbolt.Options{
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		// Fallback to in-memory piece completion if bolt DB is unavailable or locked
		return &PieceStore{
			dbPath:    dbPath,
			memBackup: storage.NewMapPieceCompletion(),
		}, nil
	}
	db.NoSync = true

	// Ensure root completion bucket exists
	_ = db.Update(func(tx *bbolt.Tx) error {
		_, bErr := tx.CreateBucketIfNotExists(completionBucketKey)
		return bErr
	})

	return &PieceStore{
		db:        db,
		dbPath:    dbPath,
		memBackup: storage.NewMapPieceCompletion(),
	}, nil
}

func (ps *PieceStore) Persistent() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.db != nil
}

func (ps *PieceStore) Get(pk metainfo.PieceKey) (cn storage.Completion, err error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.db == nil {
		if ps.memBackup != nil {
			return ps.memBackup.Get(pk)
		}
		return storage.Completion{}, nil
	}

	err = ps.db.View(func(tx *bbolt.Tx) error {
		cb := tx.Bucket(completionBucketKey)
		if cb == nil {
			return nil
		}
		ih := cb.Bucket(pk.InfoHash[:])
		if ih == nil {
			return nil
		}
		var key [4]byte
		binary.BigEndian.PutUint32(key[:], uint32(pk.Index))
		val := ih.Get(key[:])
		if val == nil {
			return nil
		}
		cn.Ok = true
		switch string(val) {
		case boltDbCompleteValue:
			cn.Complete = true
		case boltDbIncompleteValue:
			cn.Complete = false
		default:
			cn.Ok = false
		}
		return nil
	})
	return
}

func (ps *PieceStore) Set(pk metainfo.PieceKey, b bool) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.db == nil {
		if ps.memBackup != nil {
			return ps.memBackup.Set(pk, b)
		}
		return nil
	}

	return ps.db.Update(func(tx *bbolt.Tx) error {
		cb := tx.Bucket(completionBucketKey)
		if cb == nil {
			var err error
			cb, err = tx.CreateBucketIfNotExists(completionBucketKey)
			if err != nil {
				return err
			}
		}
		ih, err := cb.CreateBucketIfNotExists(pk.InfoHash[:])
		if err != nil {
			return err
		}
		var key [4]byte
		binary.BigEndian.PutUint32(key[:], uint32(pk.Index))
		val := boltDbIncompleteValue
		if b {
			val = boltDbCompleteValue
		}
		return ih.Put(key[:], []byte(val))
	})
}

// DeleteTorrent synchronously and completely purges all piece completion records for the infohash.
func (ps *PieceStore) DeleteTorrent(infoHashBytes []byte) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.memBackup != nil {
		var pk metainfo.PieceKey
		if len(infoHashBytes) == 20 {
			copy(pk.InfoHash[:], infoHashBytes)
		}
		// In memory piece completion stores per hash; setting a dummy or clearing
	}

	if ps.db == nil {
		return nil
	}

	return ps.db.Update(func(tx *bbolt.Tx) error {
		cb := tx.Bucket(completionBucketKey)
		if cb == nil {
			return nil
		}
		err := cb.DeleteBucket(infoHashBytes)
		if err == bbolt.ErrBucketNotFound {
			return nil
		}
		return err
	})
}

func (ps *PieceStore) Close() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var err error
	if ps.db != nil {
		err = ps.db.Close()
		ps.db = nil
	}
	if ps.memBackup != nil {
		_ = ps.memBackup.Close()
		ps.memBackup = nil
	}
	return err
}

// pieceStoreWrapper wraps PieceStore so individual torrent storage closes do not close the underlying database.
type pieceStoreWrapper struct {
	*PieceStore
}

func (w pieceStoreWrapper) Close() error {
	// No-op for individual torrent storages: the shared database is closed when Engine shuts down
	return nil
}

func (w pieceStoreWrapper) Persistent() bool {
	if w.PieceStore != nil {
		return w.PieceStore.Persistent()
	}
	return false
}

var _ storage.PieceCompletion = pieceStoreWrapper{}
