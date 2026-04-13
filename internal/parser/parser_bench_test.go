package parser

import (
	"testing"
)

// BenchmarkParseMinimalTorrent бенчмарк для парсинга минимального торрента
func BenchmarkParseMinimalTorrent(b *testing.B) {
	data := createMinimalTorrent()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// BenchmarkParseLargeTorrent бенчмарк для парсинга большого торрента
func BenchmarkParseLargeTorrent(b *testing.B) {
	data := createLargeTorrent(100) // 100 файлов
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// createMinimalTorrent создаёт минимальный торрент для тестов
func createMinimalTorrent() []byte {
	// Простой bencode торрент
	return []byte("d8:announce27:http://tracker.example.com4:infod6:lengthi1024e4:name9:test.txt12:piece lengthi16384e6:pieces20:01234567890123456789ee")
}

// createLargeTorrent создаёт большой торрент с множеством файлов
func createLargeTorrent(numFiles int) []byte {
	// Генерируем bencode с множеством файлов
	result := "d8:announce27:http://tracker.example.com4:infod5:filesl"
	
	for i := 0; i < numFiles; i++ {
		result += "d6:lengthi1024e4:pathl7:file"
		result += string(rune(i%10 + '0'))
		result += "ee"
	}
	
	result += "e12:piece lengthi16384e6:pieces20:01234567890123456789eee"
	
	return []byte(result)
}
