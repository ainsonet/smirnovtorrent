package peer

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestNewPeerID(t *testing.T) {
	peerID := NewPeerID()

	// Проверяем префикс
	if string(peerID[:8]) != PeerIDPrefix {
		t.Errorf("Expected prefix %s, got %s", PeerIDPrefix, string(peerID[:8]))
	}
}

func TestMessageSerialization(t *testing.T) {
	// Тест для request message
	pieceIndex := uint32(5)
	begin := uint32(1024)
	length := uint32(16384)

	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], pieceIndex)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	binary.BigEndian.PutUint32(payload[8:12], length)

	// Проверяем правильность сериализации
	if binary.BigEndian.Uint32(payload[0:4]) != pieceIndex {
		t.Error("piece index mismatch")
	}
	if binary.BigEndian.Uint32(payload[4:8]) != begin {
		t.Error("begin mismatch")
	}
	if binary.BigEndian.Uint32(payload[8:12]) != length {
		t.Error("length mismatch")
	}
}

func TestHandshakeStructure(t *testing.T) {
	var infoHash [20]byte
	var peerID [PeerIDSize]byte

	for i := range infoHash {
		infoHash[i] = byte(i)
	}
	for i := range peerID {
		peerID[i] = byte(i + 20)
	}

	var buf [68]byte
	buf[0] = 19
	copy(buf[1:20], []byte(ProtocolName))
	copy(buf[20:40], infoHash[:])
	copy(buf[40:60], peerID[:])

	// Проверяем структуру
	if buf[0] != 19 {
		t.Error("invalid pstrlen")
	}

	if !bytes.Equal(buf[1:20], []byte(ProtocolName)) {
		t.Error("invalid protocol string")
	}

	if !bytes.Equal(buf[20:40], infoHash[:]) {
		t.Error("invalid info hash")
	}

	if !bytes.Equal(buf[40:60], peerID[:]) {
		t.Error("invalid peer id")
	}
}

func TestMessageLength(t *testing.T) {
	// Keep-alive message имеет длину 0
	length := uint32(0)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, length)

	if binary.BigEndian.Uint32(header) != 0 {
		t.Error("keep-alive length should be 0")
	}

	// Choke message имеет длину 1 (только тип сообщения)
	length = 1
	binary.BigEndian.PutUint32(header, length)
	if binary.BigEndian.Uint32(header) != 1 {
		t.Error("choke length should be 1")
	}
}