package engine

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

// Piece представляет один кусок торрента
type Piece struct {
	Index     int
	Hash      []byte
	Data      []byte
	Size      int64
	Complete  bool
	Requested bool
}

// PieceManager управляет кусками торрента
type PieceManager struct {
	pieces     []*Piece
	pieceLength int
	totalSize  int64
	mu         sync.RWMutex
}

// NewPieceManager создаёт новый менеджер кусков
func NewPieceManager(pieceLength int, totalSize int64, pieceHashes []byte) *PieceManager {
	numPieces := len(pieceHashes) / 20 // каждый хеш SHA-1 = 20 байт
	
	pieces := make([]*Piece, numPieces)
	for i := 0; i < numPieces; i++ {
		hash := pieceHashes[i*20 : (i+1)*20]
		pieces[i] = &Piece{
			Index:    i,
			Hash:     hash,
			Complete: false,
		}
	}

	return &PieceManager{
		pieces:      pieces,
		pieceLength: pieceLength,
		totalSize:   totalSize,
	}
}

// GetPiece возвращает кусок по индексу
func (pm *PieceManager) GetPiece(index int) (*Piece, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if index < 0 || index >= len(pm.pieces) {
		return nil, fmt.Errorf("invalid piece index: %d", index)
	}

	return pm.pieces[index], nil
}

// GetNextPiece возвращает следующий неотмеченный кусок
func (pm *PieceManager) GetNextPiece() *Piece {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, piece := range pm.pieces {
		if !piece.Complete && !piece.Requested {
			piece.Requested = true
			return piece
		}
	}

	return nil
}
	}

	return nil
}

// MarkPieceComplete отмечает кусок как завершённый
func (pm *PieceManager) MarkPieceComplete(index int, data []byte) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if index < 0 || index >= len(pm.pieces) {
		return fmt.Errorf("invalid piece index: %d", index)
	}

	piece := pm.pieces[index]
	
	// Проверяем хеш
	hash := sha1.Sum(data)
	if string(hash[:]) != string(piece.Hash) {
		return fmt.Errorf("piece %d hash mismatch", index)
	}

	piece.Data = data
	piece.Complete = true
	piece.Requested = false

	return nil
}

// IsComplete проверяет, все ли куски загружены
func (pm *PieceManager) IsComplete() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, piece := range pm.pieces {
		if !piece.Complete {
			return false
		}
	}

	return true
}

// Progress возвращает прогресс загрузки в процентах
func (pm *PieceManager) Progress() float64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	completed := 0
	for _, piece := range pm.pieces {
		if piece.Complete {
			completed++
		}
	}

	return float64(completed) / float64(len(pm.pieces)) * 100
}

// TotalPieces возвращает общее количество кусков
func (pm *PieceManager) TotalPieces() int {
	return len(pm.pieces)
}

// CompletePieces возвращает количество завершённых кусков
func (pm *PieceManager) CompletePieces() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	count := 0
	for _, piece := range pm.pieces {
		if piece.Complete {
			count++
		}
	}

	return count
}

// GetMissingPieces возвращает индексы незавершённых кусков
func (pm *PieceManager) GetMissingPieces() []int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	missing := []int{}
	for _, piece := range pm.pieces {
		if !piece.Complete {
			missing = append(missing, piece.Index)
		}
	}

	return missing
}

// AssembleFile собирает файл из кусков
func (pm *PieceManager) AssembleFile() []byte {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var data []byte
	for _, piece := range pm.pieces {
		if piece.Complete {
			data = append(data, piece.Data...)
		}
	}

	return data
}