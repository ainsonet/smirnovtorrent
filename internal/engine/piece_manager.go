package engine

import (
	"crypto/sha1"
	"fmt"
	"sync"

	"smirnovtorrent/internal/parser"
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

// GetFileRange возвращает данные для конкретного файла
func (pm *PieceManager) GetFileRange(fileIndex int, files []parser.FileInfo) ([]byte, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if fileIndex < 0 || fileIndex >= len(files) {
		return nil, fmt.Errorf("invalid file index: %d", fileIndex)
	}

	file := files[fileIndex]
	
	// Вычисляем смещение файла в потоке данных
	var offset int64 = 0
	for i := 0; i < fileIndex; i++ {
		offset += files[i].Size
	}

	// Находим куски которые содержат этот файл
	startPiece := int(offset) / pm.pieceLength
	endPiece := int(offset + file.Size - 1) / pm.pieceLength

	var data []byte
	for i := startPiece; i <= endPiece; i++ {
		if i < len(pm.pieces) && pm.pieces[i].Complete {
			pieceData := pm.pieces[i].Data
			
			// Вычисляем смещение внутри куска
			pieceStart := i * pm.pieceLength
			fileStartInPiece := offset - int64(pieceStart)
			
			if fileStartInPiece < 0 {
				fileStartInPiece = 0
			}
			
			// Вычисляем сколько данных нужно взять из этого куска
			remaining := file.Size - int64(len(data))
			available := int64(len(pieceData)) - fileStartInPiece
			
			take := remaining
			if available < take {
				take = available
			}
			
			data = append(data, pieceData[fileStartInPiece:fileStartInPiece+take]...)
		}
	}

	return data, nil
}

// MarkPieceCompleteDirect отмечает кусок как завершённый без проверки хэша (для resume)
func (pm *PieceManager) MarkPieceCompleteDirect(index int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if index < 0 || index >= len(pm.pieces) {
		return
	}

	piece := pm.pieces[index]
	piece.Complete = true
	piece.Requested = false
}