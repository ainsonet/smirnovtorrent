package encryption

import (
	"crypto/rc4"
	"crypto/rand"
	"testing"
)

func TestNewMSEEncryption(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	enc := NewMSEEncryption(key)
	if enc == nil {
		t.Fatal("Expected non-nil encryption")
	}

	if len(enc.encryptionKey) != 32 {
		t.Errorf("Expected key length 32, got %d", len(enc.encryptionKey))
	}
}

func TestComputeHash(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	enc := NewMSEEncryption(key)
	
	data := []byte("test data")
	var infoHash [20]byte
	_, _ = rand.Read(infoHash[:])

	hash := enc.computeHash(data, infoHash, key)
	
	if len(hash) != 20 {
		t.Errorf("Expected hash length 20, got %d", len(hash))
	}
}

func TestEncryptionRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	enc := NewMSEEncryption(key)

	// Устанавливаем шифры
	encryptKey := enc.deriveKey(key, key, true)
	decryptKey := enc.deriveKey(key, key, false)
	enc.encryptCipher, _ = rc4.NewCipher(encryptKey)
	enc.decryptCipher, _ = rc4.NewCipher(decryptKey)

	// Создаём тестовые данные
	original := []byte("Hello, BitTorrent!")

	// Шифруем
	encrypted, err := enc.EncryptData(original)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Расшифровываем
	decrypted, err := enc.DecryptData(encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if string(decrypted) != string(original) {
		t.Errorf("Expected '%s', got '%s'", string(original), string(decrypted))
	}
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
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	enc := NewMSEEncryption(key)

	localS := make([]byte, 20)
	remoteS := make([]byte, 20)
	_, _ = rand.Read(localS)
	_, _ = rand.Read(remoteS)

	encryptKey := enc.deriveKey(localS, remoteS, true)
	decryptKey := enc.deriveKey(remoteS, localS, false)

	if len(encryptKey) != 16 {
		t.Errorf("Expected encrypt key length 16, got %d", len(encryptKey))
	}

	if len(decryptKey) != 16 {
		t.Errorf("Expected decrypt key length 16, got %d", len(decryptKey))
	}
}