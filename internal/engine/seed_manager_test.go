package engine

import (
	"crypto/sha1"
	"fmt"
	"testing"
	"time"
)

func TestNewSeedManager(t *testing.T) {
	pieceLength := 16384
	totalSize := int64(16384 * 10)
	numPieces := 10

	pieceHashes := make([]byte, numPieces*20)
	for i := 0; i < numPieces; i++ {
		data := []byte(fmt.Sprintf("piece%d", i))
		hash := sha1.Sum(data)
		copy(pieceHashes[i*20:(i+1)*20], hash[:])
	}

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)
	pool := NewPeerPool("test", "test", 6881, 10)

	sm := NewSeedManager(pm, pool, 5)

	if sm.targetPeers != 5 {
		t.Errorf("Expected target peers 5, got %d", sm.targetPeers)
	}

	if sm.GetActivePeerCount() != 0 {
		t.Errorf("Expected 0 active peers, got %d", sm.GetActivePeerCount())
	}
}

func TestSeedStats(t *testing.T) {
	pieceLength := 16384
	totalSize := int64(16384 * 5)
	numPieces := 5

	pieceHashes := make([]byte, numPieces*20)
	for i := 0; i < numPieces; i++ {
		data := []byte(fmt.Sprintf("piece%d", i))
		hash := sha1.Sum(data)
		copy(pieceHashes[i*20:(i+1)*20], hash[:])
	}

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)
	pool := NewPeerPool("test", "test", 6881, 10)

	sm := NewSeedManager(pm, pool, 5)

	stats := sm.GetStats()

	if stats.ActivePeers != 0 {
		t.Errorf("Expected 0 active peers, got %d", stats.ActivePeers)
	}

	if stats.TargetPeers != 5 {
		t.Errorf("Expected target 5, got %d", stats.TargetPeers)
	}

	if stats.TotalSent != 0 {
		t.Errorf("Expected 0 bytes sent, got %d", stats.TotalSent)
	}
}

func TestHandleRequest_InvalidPayload(t *testing.T) {
	pieceLength := 16384
	totalSize := int64(16384)
	numPieces := 1

	pieceHashes := make([]byte, numPieces*20)
	data := []byte("test data")
	hash := sha1.Sum(data)
	copy(pieceHashes[0:20], hash[:])

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)
	pm.MarkPieceComplete(0, data)

	pool := NewPeerPool("test", "test", 6881, 10)
	sm := NewSeedManager(pm, pool, 5)

	// Неверный размер payload
	payload := []byte{0, 0, 0} // Менее 12 байт

	err := sm.handleRequest(&SeedPeerSlot{}, payload)
	if err == nil {
		t.Error("Expected error for invalid payload, got nil")
	}
}

func TestHandleRequest_PieceNotComplete(t *testing.T) {
	pieceLength := 16384
	totalSize := int64(16384)
	numPieces := 1

	pieceHashes := make([]byte, numPieces*20)
	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)

	pool := NewPeerPool("test", "test", 6881, 10)
	sm := NewSeedManager(pm, pool, 5)

	// Payload для запроса куска 0
	payload := make([]byte, 12)
	// pieceIndex = 0, begin = 0, length = 16384

	err := sm.handleRequest(&SeedPeerSlot{}, payload)
	if err == nil {
		t.Error("Expected error for incomplete piece, got nil")
	}
}

func TestBitfieldGeneration(t *testing.T) {
	numPieces := 20
	bitfield := make([]byte, (numPieces+7)/8)

	// Устанавливаем все биты
	for i := 0; i < numPieces; i++ {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8)
		bitfield[byteIdx] |= (1 << bitIdx)
	}

	// Проверяем что все биты установлены
	expectedBytes := (numPieces + 7) / 8
	if len(bitfield) != expectedBytes {
		t.Errorf("Expected %d bytes, got %d", expectedBytes, len(bitfield))
	}

	// Все байты должны быть 0xFF
	for i, b := range bitfield {
		if i < expectedBytes-1 {
			if b != 0xFF {
				t.Errorf("Byte %d expected 0xFF, got 0x%02X", i, b)
			}
		}
	}
}

func TestCleanupInactivePeers(t *testing.T) {
	pieceLength := 16384
	totalSize := int64(16384)
	numPieces := 1

	pieceHashes := make([]byte, numPieces*20)
	data := []byte("test")
	hash := sha1.Sum(data)
	copy(pieceHashes[0:20], hash[:])

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)
	pool := NewPeerPool("test", "test", 6881, 10)

	sm := NewSeedManager(pm, pool, 5)

	// Добавляем фейкового пирa
	sm.peers["test-peer"] = &SeedPeerSlot{
		PeerID:     "test",
		LastActive: time.Now().Add(-10 * time.Minute), // Неактивен 10 минут
	}

	// Запускаем очистку
	sm.cleanupInactivePeers()

	// Пир должен быть удалён
	if len(sm.peers) != 0 {
		t.Errorf("Expected 0 peers after cleanup, got %d", len(sm.peers))
	}
}
