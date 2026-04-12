package parser

import (
	"os"
	"testing"

	"smirnovtorrent/pkg/bencode"
)

// createTestTorrent создаёт минимальный валидный .torrent файл
// Формат bencode:
// d = dict start
// 8:announce = key "announce" (8 байт)
// 15:http://test.com = value (15 байт)
// 4:info = key "info" (4 байта)
// d = dict start (info)
// 6:length = key "length" (6 байт)
// i1024e = int 1024
// 4:name = key "name" (4 байта)
// 4:test = value "test" (4 байта)
// 12:piece length = key "piece length" (12 байт)
// i16384e = int 16384
// 6:pieces = key "pieces" (6 байт)
// 20:xxxxxxxxxxxxxxxxxxxx = 20 bytes pieces
// e = dict end (info)
// e = dict end (root)
func createTestTorrent() []byte {
	return []byte("d8:announce15:http://test.com4:infod6:lengthi1024e4:name4:test12:piece lengthi16384e6:pieces20:01234567890123456789ee")
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
	// Создаём полный torrent через Marshal
	root := bencode.Dict{
		"announce": bencode.String("http://test.com/announce"),
		"info": bencode.Dict{
			"files": bencode.List{
				bencode.Dict{
					"length": bencode.Int(512),
					"path":   bencode.List{bencode.String("file1.txt")},
				},
				bencode.Dict{
					"length": bencode.Int(1024),
					"path":   bencode.List{bencode.String("file2.txt")},
				},
			},
			"piece length": bencode.Int(16384),
			"pieces":       bencode.String(make([]byte, 20)),
			"name":         bencode.String("myfolder"),
		},
	}

	data, err := bencode.Marshal(root)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	torrent, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(torrent.Info.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(torrent.Info.Files))
	}
}

func TestInfoHash(t *testing.T) {
	data := createTestTorrent()
	torrent, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Хеш должен быть 40 символов hex
	if len(torrent.Info.InfoHash) != 40 {
		t.Errorf("Expected 40 char hash, got %d", len(torrent.Info.InfoHash))
	}
}
