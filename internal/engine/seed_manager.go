package engine

import (
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"time"

	"smirnovtorrent/internal/peer"
)

// SeedManager управляет режимом раздачи (seed)
type SeedManager struct {
	pieceManager  *PieceManager
	peerPool      *PeerPool
	peers         map[string]*SeedPeerSlot
	mu            sync.RWMutex
	seedStartTime time.Time
	targetPeers   int
}

// SeedPeerSlot информация о пире при раздаче
type SeedPeerSlot struct {
	Conn        *peer.PeerConnection
	Interested  bool
	Choked      bool
	Bitfield    []byte
	PeerID      string
	LastActive  time.Time
	BytesSent   int64
}

// NewSeedManager создаёт менеджер раздачи
func NewSeedManager(pm *PieceManager, pool *PeerPool, targetPeers int) *SeedManager {
	return &SeedManager{
		pieceManager: pm,
		peerPool:     pool,
		peers:        make(map[string]*SeedPeerSlot),
		targetPeers:  targetPeers,
	}
}

// StartSeed начинает режим раздачи
func (sm *SeedManager) StartSeed() error {
	log.Println("Entering seed mode - sharing pieces...")
	sm.seedStartTime = time.Now()

	// Обновляем bitfield чтобы все пиры знали что у нас есть все куски
	sm.updateBitfield()

	// Ждём подключения новых пиров
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.cleanupInactivePeers()
			sm.announceInterested()
		}
	}
}

// HandlePeerConnection обрабатывает новое соединение от пирa
func (sm *SeedManager) HandlePeerConnection(conn *peer.PeerConnection, peerID [20]byte, bitfield []byte) error {
	peerKey := fmt.Sprintf("%p", conn)

	slot := &SeedPeerSlot{
		Conn:       conn,
		PeerID:     string(peerID[:]),
		Bitfield:   bitfield,
		LastActive: time.Now(),
	}

	sm.mu.Lock()
	sm.peers[peerKey] = slot
	sm.mu.Unlock()

	log.Printf("New peer connected: %s (%d pieces)", slot.PeerID[:8], len(bitfield)*8)

	// Отправляем наш bitfield (у нас есть все куски)
	sm.sendBitfield(conn)

	// Если пир interested, отправляем unchoke
	return nil
}

// HandleMessage обрабатывает сообщение от пирa
func (sm *SeedManager) HandleMessage(peerKey string, msgType byte, payload []byte) error {
	sm.mu.Lock()
	slot, exists := sm.peers[peerKey]
	if !exists {
		sm.mu.Unlock()
		return fmt.Errorf("unknown peer")
	}
	slot.LastActive = time.Now()
	sm.mu.Unlock()

	switch msgType {
	case peer.MsgInterested:
		sm.mu.Lock()
		slot.Interested = true
		sm.mu.Unlock()
		log.Printf("Peer %s is interested", slot.PeerID[:8])
		
		// Отправляем unchoke
		if err := slot.Conn.SendUnchoke(); err != nil {
			return err
		}

	case peer.MsgNotInterested:
		sm.mu.Lock()
		slot.Interested = false
		sm.mu.Unlock()

	case peer.MsgRequest:
		if slot.Interested && !slot.Choked {
			return sm.handleRequest(slot, payload)
		}

	case peer.MsgChoke:
		sm.mu.Lock()
		slot.Choked = true
		sm.mu.Unlock()

	case peer.MsgUnchoke:
		sm.mu.Lock()
		slot.Choked = false
		sm.mu.Unlock()
	}

	return nil
}

// handleRequest обрабатывает запрос куска от пирa
func (sm *SeedManager) handleRequest(slot *SeedPeerSlot, payload []byte) error {
	if len(payload) < 12 {
		return fmt.Errorf("invalid request payload")
	}

	pieceIndex := binary.BigEndian.Uint32(payload[0:4])
	begin := binary.BigEndian.Uint32(payload[4:8])
	length := binary.BigEndian.Uint32(payload[8:12])

	// Получаем кусок из pieceManager
	piece, err := sm.pieceManager.GetPiece(int(pieceIndex))
	if err != nil {
		return err
	}

	if !piece.Complete {
		return fmt.Errorf("piece %d not complete", pieceIndex)
	}

	// Извлекаем нужный диапазон
	end := begin + length
	if int(end) > len(piece.Data) {
		return fmt.Errorf("invalid range: begin=%d, length=%d, dataLen=%d", begin, length, len(piece.Data))
	}

	data := piece.Data[begin:end]

	// Отправляем кусок
	if err := slot.Conn.SendPiece(pieceIndex, begin, data); err != nil {
		return err
	}

	slot.BytesSent += int64(len(data))
	log.Printf("Sent piece %d:%d-%d to %s (%.1f KB)",
		pieceIndex, begin, end, slot.PeerID[:8], float64(len(data))/1024)

	return nil
}

// sendBitfield отправляет bitfield пиру
func (sm *SeedManager) sendBitfield(conn *peer.PeerConnection) error {
	numPieces := sm.pieceManager.TotalPieces()
	bitfield := make([]byte, (numPieces+7)/8)

	// У нас есть все куски, поэтому устанавливаем все биты
	for i := 0; i < numPieces; i++ {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8)
		bitfield[byteIdx] |= (1 << bitIdx)
	}

	return conn.SendMessage(peer.MsgBitfield, bitfield)
}

// updateBitfield обновляет bitfield для всех пиров
func (sm *SeedManager) updateBitfield() {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for key, slot := range sm.peers {
		if err := sm.sendBitfield(slot.Conn); err != nil {
			log.Printf("Failed to update bitfield for %s: %v", key, err)
		}
	}
}

// announceInterested сообщает всем пирам что мы interested
func (sm *SeedManager) announceInterested() {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, slot := range sm.peers {
		if !slot.Interested {
			slot.Conn.SendInterested()
		}
	}
}

// cleanupInactivePeers закрывает неактивные соединения
func (sm *SeedManager) cleanupInactivePeers() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for key, slot := range sm.peers {
		if now.Sub(slot.LastActive) > 5*time.Minute {
			log.Printf("Closing inactive peer: %s", slot.PeerID[:8])
			if slot.Conn != nil {
				slot.Conn.Close()
			}
			delete(sm.peers, key)
		}
	}
}

// GetActivePeerCount возвращает количество активных пиров
func (sm *SeedManager) GetActivePeerCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.peers)
}

// GetStats возвращает статистику раздачи
func (sm *SeedManager) GetStats() SeedStats {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := SeedStats{
		StartTime:   sm.seedStartTime,
		ActivePeers: len(sm.peers),
		TargetPeers: sm.targetPeers,
	}

	var totalSent int64
	for _, slot := range sm.peers {
		totalSent += slot.BytesSent
	}
	stats.TotalSent = totalSent

	return stats
}

// SeedStats статистика раздачи
type SeedStats struct {
	StartTime   time.Time
	ActivePeers int
	TargetPeers int
	TotalSent   int64
}

// Stop завершает раздачу
func (sm *SeedManager) Stop() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for key, slot := range sm.peers {
		slot.Conn.Close()
		delete(sm.peers, key)
	}

	log.Printf("Seed mode stopped. Total sent: %.1f MB", float64(sm.GetStats().TotalSent)/1024/1024)
}