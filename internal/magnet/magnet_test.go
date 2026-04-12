package magnet

import (
	"strings"
	"testing"
)

func TestParse_ValidMagnet(t *testing.T) {
	magnet := "magnet:?xt=urn:btih:73ef7ed9f70e94f1e3a4b8b5c2d1e0f9a8b7c6d5&dn=Test+Torrent&tr=http://tracker.example.com/announce"

	link, err := Parse(magnet)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if link.InfoHash != "73ef7ed9f70e94f1e3a4b8b5c2d1e0f9a8b7c6d5" {
		t.Errorf("Expected info hash '73ef7ed9f70e94f1e3a4b8b5c2d1e0f9a8b7c6d5', got '%s'", link.InfoHash)
	}

	if link.DisplayName != "Test Torrent" {
		t.Errorf("Expected display name 'Test Torrent', got '%s'", link.DisplayName)
	}

	if len(link.Trackers) != 1 {
		t.Errorf("Expected 1 tracker, got %d", len(link.Trackers))
	}

	if link.Trackers[0] != "http://tracker.example.com/announce" {
		t.Errorf("Expected tracker 'http://tracker.example.com/announce', got '%s'", link.Trackers[0])
	}
}

func TestParse_MultipleTrackers(t *testing.T) {
	magnet := "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678&tr=http://tracker1.com&tr=http://tracker2.com"

	link, err := Parse(magnet)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(link.Trackers) != 2 {
		t.Errorf("Expected 2 trackers, got %d", len(link.Trackers))
	}
}

func TestParse_InvalidMagnet(t *testing.T) {
	_, err := Parse("not a magnet link")
	if err == nil {
		t.Error("Expected error for invalid magnet link, got nil")
	}
}

func TestParse_MissingInfoHash(t *testing.T) {
	magnet := "magnet:?dn=Test"
	_, err := Parse(magnet)
	if err == nil {
		t.Error("Expected error for missing info hash, got nil")
	}
}

func TestParse_WithPEX(t *testing.T) {
	magnet := "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678&x.pe=1"

	link, err := Parse(magnet)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !link.PEX {
		t.Error("Expected PEX to be true")
	}
}

func TestParse_WithDHT(t *testing.T) {
	magnet := "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678&dht=on"

	link, err := Parse(magnet)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !link.DHT {
		t.Error("Expected DHT to be true")
	}
}

func TestString(t *testing.T) {
	link := &MagnetLink{
		InfoHash:    "1234567890abcdef1234567890abcdef12345678",
		DisplayName: "Test Torrent",
		Trackers:    []string{"http://tracker.com/announce"},
		DHT:         true,
	}

	str := link.String()

	expected := "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678&dn=Test+Torrent&tr=http%3A%2F%2Ftracker.com%2Fannounce&dht=on"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestExtractInfoHash(t *testing.T) {
	magnet := "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678"

	hash, err := ExtractInfoHash(magnet)
	if err != nil {
		t.Fatalf("ExtractInfoHash failed: %v", err)
	}

	if hash != "1234567890abcdef1234567890abcdef12345678" {
		t.Errorf("Expected '1234567890abcdef1234567890abcdef12345678', got '%s'", hash)
	}
}

func TestExtractDisplayName(t *testing.T) {
	magnet := "magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678&dn=My+Torrent"

	name, err := ExtractDisplayName(magnet)
	if err != nil {
		t.Fatalf("ExtractDisplayName failed: %v", err)
	}

	if name != "My Torrent" {
		t.Errorf("Expected 'My Torrent', got '%s'", name)
	}
}

func TestIsMagnetLink(t *testing.T) {
	if !IsMagnetLink("magnet:?xt=urn:btih:1234567890abcdef1234567890abcdef12345678") {
		t.Error("Expected true for magnet link")
	}

	if IsMagnetLink("http://example.com/file.torrent") {
		t.Error("Expected false for non-magnet link")
	}
}

func TestBuildMagnetLink(t *testing.T) {
	trackers := []string{"http://tracker1.com", "http://tracker2.com"}
	link := BuildMagnetLink("1234567890abcdef1234567890abcdef12345678", "Test", trackers)

	if !strings.HasPrefix(link, "magnet:?") {
		t.Error("Expected magnet link to start with 'magnet:?'")
	}

	if !strings.Contains(link, "1234567890abcdef1234567890abcdef12345678") {
		t.Error("Expected link to contain info hash")
	}
}