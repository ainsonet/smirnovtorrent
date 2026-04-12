package tracker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnnounce_Success(t *testing.T) {
	// Создаём тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"interval": 300,
			"peers": [
				{"ip": "192.168.1.1", "port": 6881, "peer id": "-TR3000-abcdef123456", "downloaded": 0, "left": 1024, "uploaded": 0}
			]
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer server.Close()

	tracker := NewTracker(server.URL)

	params := AnnounceParams{
		InfoHash:   "1234567890123456789012345678901234567890",
		PeerID:     "-TR3000-abcdef123456",
		Port:       6881,
		Downloaded: 0,
		Left:       1024,
		Uploaded:   0,
		Event:      "started",
	}

	resp, err := tracker.Announce(params)
	if err != nil {
		t.Fatalf("Announce failed: %v", err)
	}

	if resp.Interval != 300 {
		t.Errorf("Expected interval 300, got %d", resp.Interval)
	}

	if len(resp.Peers) != 1 {
		t.Fatalf("Expected 1 peer, got %d", len(resp.Peers))
	}

	if resp.Peers[0].IP != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", resp.Peers[0].IP)
	}
}

func TestAnnounce_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tracker := NewTracker(server.URL)

	params := AnnounceParams{
		InfoHash: "1234567890123456789012345678901234567890",
		PeerID:   "test",
		Port:     6881,
	}

	_, err := tracker.Announce(params)
	if err == nil {
		t.Fatal("Expected error for 500 status, got nil")
	}
}

func TestAnnounce_FailureReason(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"failure reason": "invalid info_hash"}`))
	}))
	defer server.Close()

	tracker := NewTracker(server.URL)

	params := AnnounceParams{
		InfoHash: "invalid",
		PeerID:   "test",
		Port:     6881,
	}

	_, err := tracker.Announce(params)
	if err == nil {
		t.Fatal("Expected error for failure reason, got nil")
	}

	if err.Error() != "tracker error: invalid info_hash" {
		t.Errorf("Expected tracker error message, got: %s", err.Error())
	}
}

func TestParsePeerURL(t *testing.T) {
	url := "http://tracker.example.com/announce"
	tracker, err := ParsePeerURL(url)
	if err != nil {
		t.Fatalf("ParsePeerURL failed: %v", err)
	}

	if tracker.announceURL != url {
		t.Errorf("Expected announce URL %s, got %s", url, tracker.announceURL)
	}
}

func TestParsePeerURL_Empty(t *testing.T) {
	_, err := ParsePeerURL("")
	if err == nil {
		t.Fatal("Expected error for empty URL, got nil")
	}
}