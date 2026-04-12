package engine

import (
	"crypto/sha1"
	"testing"
)

func TestNewPieceManager(t *testing.T) {
	pieceLength := 16384
	totalSize := 1024 * 1024 // 1MB
	numPieces := 64

	// Создаём тестовые хеши
	pieceHashes := make([]byte, numPieces*20)
	for i := 0; i < numPieces; i++ {
		hash := sha1.Sum([]byte(fmt.Sprintf("piece%d", i)))
		copy(pieceHashes[i*20:(i+1)*20], hash[:])
	}

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)

	if pm.TotalPieces() != numPieces {
		t.Errorf("Expected %d pieces, got %d", numPieces, pm.TotalPieces())
	}

	if pm.CompletePieces() != 0 {
		t.Errorf("Expected 0 complete pieces, got %d", pm.CompletePieces())
	}
}

func TestMarkPieceComplete(t *testing.T) {
	pieceLength := 16384
	totalSize := 16384 * 2 // 2 pieces
	numPieces := 2

	pieceHashes := make([]byte, numPieces*20)
	data := []byte("test data for piece 0")
	hash := sha1.Sum(data)
	copy(pieceHashes[0:20], hash[:])

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)

	err := pm.MarkPieceComplete(0, data)
	if err != nil {
		t.Fatalf("Failed to mark piece complete: %v", err)
	}

	if pm.CompletePieces() != 1 {
		t.Errorf("Expected 1 complete piece, got %d", pm.CompletePieces())
	}

	if !pm.IsComplete() {
		t.Error("Expected torrent to be incomplete")
	}
}

func TestMarkPieceComplete_HashMismatch(t *testing.T) {
	pieceLength := 16384
	totalSize := 16384
	numPieces := 1

	pieceHashes := make([]byte, numPieces*20)
	// Создаём неправильный хеш
	hash := sha1.Sum([]byte("wrong data"))
	copy(pieceHashes[0:20], hash[:])

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)

	// Пытаемся добавить кусок с неправильным хешем
	data := []byte("test data")
	err := pm.MarkPieceComplete(0, data)
	if err == nil {
		t.Fatal("Expected error for hash mismatch, got nil")
	}
}

func TestProgress(t *testing.T) {
	pieceLength := 16384
	totalSize := 16384 * 4 // 4 pieces
	numPieces := 4

	pieceHashes := make([]byte, numPieces*20)
	for i := 0; i < numPieces; i++ {
		data := []byte(fmt.Sprintf("piece%d", i))
		hash := sha1.Sum(data)
		copy(pieceHashes[i*20:(i+1)*20], hash[:])
	}

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)

	// Загружаем 2 из 4 кусков
	for i := 0; i < 2; i++ {
		data := []byte(fmt.Sprintf("piece%d", i))
		pm.MarkPieceComplete(i, data)
	}

	progress := pm.Progress()
	if progress != 50.0 {
		t.Errorf("Expected 50%% progress, got %.2f%%", progress)
	}
}

func TestGetNextPiece(t *testing.T) {
	pieceLength := 16384
	totalSize := 16384 * 3
	numPieces := 3

	pieceHashes := make([]byte, numPieces*20)
	for i := 0; i < numPieces; i++ {
		data := []byte(fmt.Sprintf("piece%d", i))
		hash := sha1.Sum(data)
		copy(pieceHashes[i*20:(i+1)*20], hash[:])
	}

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)

	// Получаем кусочки
	piece1 := pm.GetNextPiece()
	if piece1 == nil || piece1.Index != 0 {
		t.Error("Expected first piece")
	}

	piece2 := pm.GetNextPiece()
	if piece2 == nil || piece2.Index != 1 {
		t.Error("Expected second piece")
	}

	// Отмечаем первый кусок как завершённый
	pm.MarkPieceComplete(0, []byte("piece0"))

	// Следующий кусок должен быть 2 (1 уже запрошен)
	piece3 := pm.GetNextPiece()
	if piece3 == nil || piece3.Index != 2 {
		t.Error("Expected third piece")
	}
}

func TestIsComplete(t *testing.T) {
	pieceLength := 16384
	totalSize := 16384 * 2
	numPieces := 2

	pieceHashes := make([]byte, numPieces*20)
	for i := 0; i < numPieces; i++ {
		data := []byte(fmt.Sprintf("piece%d", i))
		hash := sha1.Sum(data)
		copy(pieceHashes[i*20:(i+1)*20], hash[:])
	}

	pm := NewPieceManager(pieceLength, totalSize, pieceHashes)

	if pm.IsComplete() {
		t.Error("Expected incomplete torrent")
	}

	// Загружаем все кусочки
	for i := 0; i < numPieces; i++ {
		data := []byte(fmt.Sprintf("piece%d", i))
		pm.MarkPieceComplete(i, data)
	}

	if !pm.IsComplete() {
		t.Error("Expected complete torrent")
	}
}