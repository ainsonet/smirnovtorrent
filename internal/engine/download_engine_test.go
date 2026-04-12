package engine

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"smirnovtorrent/internal/parser"
)

func TestNewDownloadEngine(t *testing.T) {
	torrent := &parser.Torrent{
		Info: parser.TorrentInfo{
			Name:        "test torrent",
			PieceLength: 16384,
			Pieces:      make([]byte, 20*10), // 10 pieces
			Files: []parser.FileInfo{
				{Path: "test.txt", Size: 163840},
			},
		},
		Announce: "http://test.com/announce",
	}

	eng := NewDownloadEngine(torrent, "")

	if eng.torrent != torrent {
		t.Error("Torrent not set correctly")
	}

	if eng.outputDir != torrent.Info.Name {
		t.Errorf("Expected output dir %s, got %s", torrent.Info.Name, eng.outputDir)
	}

	if eng.pieceManager != nil {
		t.Error("Piece manager should be nil initially")
	}
}

func TestNewDownloadEngine_CustomOutputDir(t *testing.T) {
	torrent := &parser.Torrent{
		Info: parser.TorrentInfo{
			Name: "test torrent",
		},
	}

	eng := NewDownloadEngine(torrent, "/custom/path")

	if eng.outputDir != "/custom/path" {
		t.Errorf("Expected custom output dir, got %s", eng.outputDir)
	}
}

func TestGetNextPiece(t *testing.T) {
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

	// Получаем все кусочки
	for i := 0; i < numPieces; i++ {
		piece := pm.GetNextPiece()
		if piece == nil {
			t.Fatalf("Expected piece %d, got nil", i)
		}
		if piece.Index != i {
			t.Errorf("Expected piece index %d, got %d", i, piece.Index)
		}
	}

	// Следующий кусок должен быть nil (все запрошены)
	piece := pm.GetNextPiece()
	if piece != nil {
		t.Error("Expected nil when all pieces requested")
	}
}

func TestAssembleFiles_SingleFile(t *testing.T) {
	// TODO: реализовать тест для single file mode
	// Требуется создать тестовый торрент и проверить сборку файла
}

func TestDownloadStatus(t *testing.T) {
	torrent := &parser.Torrent{
		Info: parser.TorrentInfo{
			Name:        "test",
			PieceLength: 16384,
			Pieces:      make([]byte, 20*10),
			Files: []parser.FileInfo{
				{Path: "test.txt", Size: 163840},
			},
		},
		Announce: "http://test.com/announce",
	}

	eng := NewDownloadEngine(torrent, "")
	eng.pieceManager = NewPieceManager(16384, 163840, torrent.Info.Pieces)

	// Инициализируем статус
	eng.updateStatus()

	status := eng.GetStatus()

	if status.Progress != 0 {
		t.Errorf("Expected 0 progress, got %f", status.Progress)
	}

	if status.Downloaded != 0 {
		t.Errorf("Expected 0 downloaded, got %d", status.Downloaded)
	}

	if status.TotalSize != 163840 {
		t.Errorf("Expected total size 163840, got %d", status.TotalSize)
	}
}