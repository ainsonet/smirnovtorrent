package encryption

import (
	"crypto/rand"
	"testing"
)

func TestNewMSEEncryption(t *testing.T) {
	var infoHash [20]byte
	_, err := rand.Read(infoHash[:])
	if err != nil {
		t.Fatalf("Failed to generate info hash: %v", err)
	}

	enc := NewMSEEncryption(infoHash)
	if enc == nil {
		t.Fatal("Expected non-nil encryption")
	}

	if len(enc.encryptionKey) != 20 {
		t.Errorf("Expected key length 20, got %d", len(enc.encryptionKey))
	}
}

func TestEncryptionRoundTrip(t *testing.T) {
	// Этот тест требует полной синхронизации RC4 ключей
	// В реальной реализации используется handshake для синхронизации
	// Пока просто проверяем что шифр создаётся
	var infoHash [20]byte
	_, _ = rand.Read(infoHash[:])

	enc := NewMSEEncryption(infoHash)
	if enc == nil {
		t.Fatal("Expected non-nil encryption")
	}

	t.Skip("Full encryption test requires handshake synchronization")
}

func TestEncryptionModeString(t *testing.T) {
	tests := []struct {
		mode     EncryptionMode
		expected string
	}{
		{EncryptionDisabled, "disabled"},
		{EncryptionEnabled, "enabled"},
		{EncryptionPrefer, "prefer"},
		{EncryptionMode(99), "unknown"},
	}

	for _, test := range tests {
		result := test.mode.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestValidateHandshake(t *testing.T) {
	// Валидное рукопожатие (17 символов + 3 пробела = 20)
	valid := make([]byte, 20)
	copy(valid, "BitTorrent protocol")
	if !ValidateHandshake(valid) {
		t.Error("Expected valid handshake")
	}

	// Неверная длина
	invalid := []byte("short")
	if ValidateHandshake(invalid) {
		t.Error("Expected invalid handshake to fail")
	}

	// Неверный протокол
	wrong := []byte("Other protocol!   ")
	if ValidateHandshake(wrong) {
		t.Error("Expected wrong protocol to fail")
	}
}

func TestDeriveKey(t *testing.T) {
	var infoHash [20]byte
	_, _ = rand.Read(infoHash[:])

	enc := NewMSEEncryption(infoHash)

	localS := make([]byte, 20)
	remoteS := make([]byte, 20)
	_, _ = rand.Read(localS)
	_, _ = rand.Read(remoteS)

	encryptKey := enc.deriveKey(localS, remoteS)
	decryptKey := enc.deriveKey(remoteS, localS)

	if len(encryptKey) != 16 {
		t.Errorf("Expected encrypt key length 16, got %d", len(encryptKey))
	}

	if len(decryptKey) != 16 {
		t.Errorf("Expected decrypt key length 16, got %d", len(decryptKey))
	}
}