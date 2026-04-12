package parser

import (
	"os"
	"testing"
)

// Создаем тестовый .torrent файл в памяти
func createTestTorrent() []byte {
	// Минимальный валидный .torrent (single file)
	return []byte("d8:announce15:http://test.com4:infod6:lengthi1024e4:name4:test12:piece lengthi16384e6:pieces20:00000000000000000000e")
}

func TestParseMinimalTorrent(t *testing.T) {
	data := createTestTorrent()
	torrent, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if torrent.Announce != "http://test.com" {
		t.Errorf("Expected announce 'http://test.com', got '%s'", torrent.Announce)
	}

	if torrent.Info.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", torrent.Info.Name)
	}

	if torrent.Info.PieceLength != 16384 {
		t.Errorf("Expected piece length 16384, got %d", torrent.Info.PieceLength)
	}

	if len(torrent.Info.Pieces) != 20 {
		t.Errorf("Expected 20 pieces bytes, got %d", len(torrent.Info.Pieces))
	}
}

func TestParseTorrentFile(t *testing.T) {
	// Создаем временный файл
	data := createTestTorrent()
	tmpFile := "test.torrent"
	err := os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)

	torrent, err := ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if torrent.Info.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", torrent.Info.Name)
	}
}

func TestTorrentIsValid(t *testing.T) {
	data := createTestTorrent()
	torrent, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	err = torrent.IsValid()
	if err != nil {
		t.Errorf("Torrent should be valid: %v", err)
	}
}

func TestTorrentTotalSize(t *testing.T) {
	data := createTestTorrent()
	torrent, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if torrent.TotalSize() != 1024 {
		t.Errorf("Expected total size 1024, got %d", torrent.TotalSize())
	}
}

func TestParseMultiFileTorrent(t *testing.T) {
	// Multi-file torrent
	data := []byte("d8:announce15:http://test.com4:infod4:filesld6:lengthi512e4:pa" +
		"thl3:fil3:txteee6:lengthi1024e4:pa" +
		"thl3:fil2:txteee12:piece lengthi16384e6:pieces20:00000000000000000000e4:name10:myfolder12:announce-listlli15:http://test2.comeeee")

	torrent, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(torrent.Info.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(torrent.Info.Files))
	}

	if len(torrent.AnnounceList) > 0 {
		if len(torrent.AnnounceList[0]) != 1 {
			t.Errorf("Expected 1 backup tracker, got %d", len(torrent.AnnounceList[0]))
		}
	}
}

func TestInfoHash(t *testing.T) {
	data := createTestTorrent()
	torrent, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Хеш должен быть 40 символов hex
	if len(torrent.InfoHash) != 40 {
		t.Errorf("Expected 40 char hash, got %d", len(torrent.InfoHash))
	}
}
