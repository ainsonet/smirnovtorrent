package dht

import (
	"testing"
)

func TestNewDHTClient(t *testing.T) {
	bootstrap := []string{"router.bittorrent.com:6881", "dht.transmissionbt.com:6881"}

	client, err := NewDHTClient(bootstrap, 0)
	if err != nil {
		t.Fatalf("NewDHTClient failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.GetNodeCount() != 0 {
		t.Errorf("Expected 0 nodes, got %d", client.GetNodeCount())
	}
}

func TestPeerAddressParsing(t *testing.T) {
	ip, port, err := ParsePeerAddress("192.168.1.100:6881")
	if err != nil {
		t.Fatalf("ParsePeerAddress failed: %v", err)
	}

	if ip != "192.168.1.100" {
		t.Errorf("Expected IP '192.168.1.100', got '%s'", ip)
	}

	if port != 6881 {
		t.Errorf("Expected port 6881, got %d", port)
	}
}

func TestPeerToString(t *testing.T) {
	str := PeerToString("10.0.0.1", 6882)

	if str != "10.0.0.1:6882" {
		t.Errorf("Expected '10.0.0.1:6882', got '%s'", str)
	}
}

func TestEncodeDecodePeerInfo(t *testing.T) {
	ip := "192.168.1.1"
	port := uint16(6881)

	encoded, err := EncodePeerInfo(ip, port)
	if err != nil {
		t.Fatalf("EncodePeerInfo failed: %v", err)
	}

	decodedIP, decodedPort, err := DecodePeerInfo(encoded)
	if err != nil {
		t.Fatalf("DecodePeerInfo failed: %v", err)
	}

	if decodedIP != ip {
		t.Errorf("Expected IP '%s', got '%s'", ip, decodedIP)
	}

	if decodedPort != port {
		t.Errorf("Expected port %d, got %d", port, decodedPort)
	}
}

func TestEncodePeerInfo_InvalidIP(t *testing.T) {
	_, err := EncodePeerInfo("invalid.ip.address", 6881)
	if err == nil {
		t.Error("Expected error for invalid IP, got nil")
	}
}

func TestDecodePeerInfo_InvalidLength(t *testing.T) {
	data := []byte{0, 1, 2} // Менее 6 байт
	_, _, err := DecodePeerInfo(data)
	if err == nil {
		t.Error("Expected error for invalid length, got nil")
	}
}

func TestPeerInfoRoundTrip(t *testing.T) {
	// Тестируем кодирование/декодирование нескольких пиров
	peers := []struct {
		ip   string
		port uint16
	}{
		{"192.168.1.1", 6881},
		{"10.0.0.1", 6882},
		{"172.16.0.1", 6883},
	}

	for _, p := range peers {
		encoded, err := EncodePeerInfo(p.ip, p.port)
		if err != nil {
			t.Fatalf("Encode failed for %s:%d: %v", p.ip, p.port, err)
		}

		decodedIP, decodedPort, err := DecodePeerInfo(encoded)
		if err != nil {
			t.Fatalf("Decode failed for %s:%d: %v", p.ip, p.port, err)
		}

		if decodedIP != p.ip || decodedPort != p.port {
			t.Errorf("Round trip failed for %s:%d: got %s:%d",
				p.ip, p.port, decodedIP, decodedPort)
		}
	}
}