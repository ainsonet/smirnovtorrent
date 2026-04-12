package engine

import (
	"math/rand"
	"sync"
	"time"
)

// RarestFirstManager выбирает куски по алгоритму Rarest-first
type RarestFirstManager struct {
	pieceManager *PieceManager
	peerPool     *PeerPool
	pieceCounts  []int // сколько пиров имеет каждый кусок
	mu           sync.RWMutex
	lastRefresh  time.Time
	refreshInterval time.Duration
}

// NewRarestFirstManager создаёт менеджер с Rarest-first
func NewRarestFirstManager(pm *PieceManager, pool *PeerPool) *RarestFirstManager {
	return &RarestFirstManager{
		pieceManager:    pm,
		peerPool:        pool,
		pieceCounts:     make([]int, pm.TotalPieces()),
		refreshInterval: 5 * time.Second,
	}
}

// GetNextPieceRarest возвращает следующий кусок по Rarest-first
func (rf *RarestFirstManager) GetNextPieceRarest() *Piece {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// Обновляем статистику если прошло много времени
	if time.Since(rf.lastRefresh) > rf.refreshInterval {
		rf.refreshPieceCounts()
	}

	// Находим самые редкие куски
	missing := rf.pieceManager.GetMissingPieces()
	if len(missing) == 0 {
		return nil
	}

	// Считаем количество пиров для каждого куска
	minCount := len(rf.pieceCounts) + 1
	for _, idx := range missing {
		if idx < len(rf.pieceCounts) {
			count := rf.pieceCounts[idx]
			if count < minCount {
				minCount = count
			}
		}
	}

	// Собираем все куски с минимальным количеством
	var candidates []int
	for _, idx := range missing {
		if idx < len(rf.pieceCounts) && rf.pieceCounts[idx] == minCount {
			candidates = append(candidates, idx)
		}
	}

	// Выбираем случайный из кандидатов
	if len(candidates) == 0 {
		return rf.pieceManager.GetNextPiece()
	}

	chosenIndex := candidates[rand.Intn(len(candidates))]
	
	// Проверяем что у кого-то этот кусок есть
	peersWithPiece := rf.peerPool.GetPeersWithPiece(chosenIndex)
	if len(peersWithPiece) == 0 {
		// Никто не имеет этот кусок, берём следующий доступный
		return rf.pieceManager.GetNextPiece()
	}

	return rf.pieceManager.pieces[chosenIndex]
}

// refreshPieceCounts обновляет статистику наличия кусков
func (rf *RarestFirstManager) refreshPieceCounts() {
	rf.pieceCounts = make([]int, rf.pieceManager.TotalPieces())

	poolSize := rf.peerPool.GetActivePeerCount()
	if poolSize == 0 {
		return
	}

	// Проходим по всем пирам и считаем куски
	// Упрощённая версия - в реальности нужно хранить информацию о каждом пире
	rf.peerPool.mu.RLock()
	for _, slot := range rf.peerPool.peers {
		for i, has := range slot.HasPieces {
			if has && i < len(rf.pieceCounts) {
				rf.pieceCounts[i]++
			}
		}
	}
	rf.peerPool.mu.RUnlock()

	rf.lastRefresh = time.Now()
}

// UpdatePieceCount обновляет счётчик для конкретного куска
func (rf *RarestFirstManager) UpdatePieceCount(pieceIndex int, delta int) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if pieceIndex < len(rf.pieceCounts) {
		rf.pieceCounts[pieceIndex] += delta
		if rf.pieceCounts[pieceIndex] < 0 {
			rf.pieceCounts[pieceIndex] = 0
		}
	}
}

// GetPieceRarity возвращает редкость куска (сколько пиров его имеют)
func (rf *RarestFirstManager) GetPieceRarity(pieceIndex int) int {
	rf.mu.RLock()
	defer rf.mu.RUnlock()

	if pieceIndex < len(rf.pieceCounts) {
		return rf.pieceCounts[pieceIndex]
	}
	return 0
}

// EmergencyMode возвращает любой доступный кусок когда пиров мало
func (rf *RarestFirstManager) EmergencyMode() *Piece {
	return rf.pieceManager.GetNextPiece()
}