package parser

import (
	"testing"
	
	"smirnovtorrent/pkg/bencode"
)

// BenchmarkParseMinimalTorrent бенчмарк для парсинга минимального торрента
func BenchmarkParseMinimalTorrent(b *testing.B) {
	data := createTestTorrent()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// BenchmarkParseMultiFileTorrent бенчмарк для парсинга multi-file торрента
func BenchmarkParseMultiFileTorrent(b *testing.B) {
	root := createMultiFileTorrentData(50) // 50 файлов
	
	data, err := bencode.Marshal(root)
	if err != nil {
		b.Fatalf("Marshal failed: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// createMultiFileTorrentData создаёт данные для multi-file торрента
func createMultiFileTorrentData(numFiles int) bencode.Dict {
	files := bencode.List{}
	for i := 0; i < numFiles; i++ {
		files = append(files, bencode.Dict{
			"length": bencode.Int(1024),
			"path":   bencode.List{bencode.String("file" + string(rune(i+'0')) + ".txt")},
		})
	}
	
	return bencode.Dict{
		"announce": bencode.String("http://tracker.example.com"),
		"info": bencode.Dict{
			"files":        files,
			"piece length": bencode.Int(16384),
			"pieces":       bencode.String(make([]byte, 20)),
			"name":         bencode.String("testfolder"),
		},
	}
}
