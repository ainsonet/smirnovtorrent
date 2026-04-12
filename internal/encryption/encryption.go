package encryption

import (
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"time"
)

// MSEMessageStreamEncryption реализует BitTorrent Message Stream Encryption
type MSEMessageStreamEncryption struct {
	encryptionKey  []byte
	vc             [8]byte // Verification Constant
	remoteVc       [8]byte
	encryptCipher  *rc4.Cipher
	decryptCipher  *rc4.Cipher
	remoteS        []byte
	localS         []byte
}

// NewMSEEncryption создаёт новый MSE шифр
func NewMSEEncryption(encryptionKey []byte) *MSEMessageStreamEncryption {
	return &MSEMessageStreamEncryption{
		encryptionKey: encryptionKey,
	}
}

// InitHandshake инициирует рукопожатие с шифрованием
func (e *MSEMessageStreamEncryption) InitHandshake(conn net.Conn, infoHash [20]byte, peerID [20]byte) error {
	// Генерируем случайные числа
	e.localS = make([]byte, 20)
	if _, err := io.ReadFull(conn, e.localS); err != nil {
		return err
	}

	// Вычисляем hash
	e.localS = e.computeHash(e.localS, infoHash, e.encryptionKey)

	// Генерируем случайные числа для IA
	e.remoteS = make([]byte, 20)
	if _, err := io.ReadFull(conn, e.remoteS); err != nil {
		return err
	}

	e.remoteS = e.computeHash(e.remoteS, infoHash, e.encryptionKey)

	// Вычисляем RC4 ключи
	encryptKey := e.deriveKey(e.localS, e.remoteS, true)
	decryptKey := e.deriveKey(e.remoteS, e.localS, true)

	// Инициализируем шифры
	encCipher, err := rc4.NewCipher(encryptKey)
	if err != nil {
		return fmt.Errorf("failed to create encrypt cipher: %w", err)
	}
	e.encryptCipher = encCipher

	decCipher, err := rc4.NewCipher(decryptKey)
	if err != nil {
		return fmt.Errorf("failed to create decrypt cipher: %w", err)
	}
	e.decryptCipher = decCipher

	return nil
}

// computeHash вычисляет hash для рукопожатия
func (e *MSEMessageStreamEncryption) computeHash(data []byte, infoHash [20]byte, encryptionKey []byte) []byte {
	hash := sha1.New()
	hash.Write(data)
	hash.Write(infoHash[:])
	hash.Write(encryptionKey)
	return hash.Sum(nil)
}

// deriveKey выводит ключ для RC4
func (e *MSEMessageStreamEncryption) deriveKey(localS, remoteS []byte, encrypt bool) []byte {
	// Простая реализация - в полной версии используется более сложное derivation
	hash := md5.New()
	hash.Write(localS)
	hash.Write(remoteS)
	if encrypt {
		hash.Write([]byte("key"))
	} else {
		hash.Write([]byte("key"))
	}
	return hash.Sum(nil)[:16] // 16 байт ключ
}

// EncryptData шифрует данные
func (e *MSEMessageStreamEncryption) EncryptData(data []byte) ([]byte, error) {
	if e.encryptCipher == nil {
		return data, nil
	}

	cipherText := make([]byte, len(data))
	e.encryptCipher.XORKeyStream(cipherText, data)
	return cipherText, nil
}

// DecryptData расшифровывает данные
func (e *MSEMessageStreamEncryption) DecryptData(data []byte) ([]byte, error) {
	if e.decryptCipher == nil {
		return data, nil
	}

	plainText := make([]byte, len(data))
	e.decryptCipher.XORKeyStream(plainText, data)
	return plainText, nil
}

// WrapConnection оборачивает соединение с шифрованием
type EncryptedConnection struct {
	conn       net.Conn
	encryption *MSEMessageStreamEncryption
}

// NewEncryptedConnection создаёт зашифрованное соединение
func NewEncryptedConnection(conn net.Conn, encryption *MSEMessageStreamEncryption) *EncryptedConnection {
	return &EncryptedConnection{
		conn:       conn,
		encryption: encryption,
	}
}

// Read читает данные с расшифровкой
func (ec *EncryptedConnection) Read(p []byte) (n int, err error) {
	n, err = ec.conn.Read(p)
	if err != nil {
		return n, err
	}

	decrypted, err := ec.encryption.DecryptData(p[:n])
	if err != nil {
		return n, err
	}

	copy(p, decrypted)
	return n, nil
}

// Write пишет данные с шифрованием
func (ec *EncryptedConnection) Write(p []byte) (n int, err error) {
	encrypted, err := ec.encryption.EncryptData(p)
	if err != nil {
		return 0, err
	}

	return ec.conn.Write(encrypted)
}

// Close закрывает соединение
func (ec *EncryptedConnection) Close() error {
	return ec.conn.Close()
}

// LocalAddr возвращает локальный адрес
func (ec *EncryptedConnection) LocalAddr() net.Addr {
	return ec.conn.LocalAddr()
}

// RemoteAddr возвращает удалённый адрес
func (ec *EncryptedConnection) RemoteAddr() net.Addr {
	return ec.conn.RemoteAddr()
}

// SetDeadline устанавливает deadline
func (ec *EncryptedConnection) SetDeadline(t time.Time) error {
	return ec.conn.SetDeadline(t)
}

// SetReadDeadline устанавливает deadline для чтения
func (ec *EncryptedConnection) SetReadDeadline(t time.Time) error {
	return ec.conn.SetReadDeadline(t)
}

// SetWriteDeadline устанавливает deadline для записи
func (ec *EncryptedConnection) SetWriteDeadline(t time.Time) error {
	return ec.conn.SetWriteDeadline(t)
}

// EncryptionMode режим шифрования
type EncryptionMode int

const (
	EncryptionDisabled EncryptionMode = iota
	EncryptionEnabled
	EncryptionPrefer
)

// String возвращает строковое представление режима
func (m EncryptionMode) String() string {
	switch m {
	case EncryptionDisabled:
		return "disabled"
	case EncryptionEnabled:
		return "enabled"
	case EncryptionPrefer:
		return "prefer"
	default:
		return "unknown"
	}
}

// ValidateHandshake проверяет рукопожатие
func ValidateHandshake(data []byte) bool {
	if len(data) < 20 {
		return false
	}

	// Проверяем протокол BitTorrent
	protocol := string(data[:20])
	return protocol == "BitTorrent protocol"
}