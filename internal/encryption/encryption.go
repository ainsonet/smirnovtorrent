package encryption

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"time"
)

// MSE支持的模式
const (
	// Plaintext handshake
	MSEPlainText = iota
	// RC4 encrypted
	MSEEncrypted
	// RC4 with header obfuscation
	MSEObfuscated
)

// MSEMessageStreamEncryption реализует BitTorrent Message Stream Encryption
type MSEMessageStreamEncryption struct {
	encryptionKey  []byte // S = hash(info_hash, key)
	vc             [8]byte // Verification Constant (8 zero bytes)
	remoteVc       [8]byte
	encryptCipher  *rc4.Cipher
	decryptCipher  *rc4.Cipher
	remoteS        []byte
	localS         []byte
	mode           int
	infoHash       [20]byte
	peerID         [20]byte
}

// NewMSEEncryption создаёт новый MSE шифр
func NewMSEEncryption(infoHash [20]byte) *MSEMessageStreamEncryption {
	// S = hash('key', S) где S = info_hash
	key := []byte("key")
	hash := sha1.New()
	hash.Write(key)
	hash.Write(infoHash[:])
	encryptionKey := hash.Sum(nil)
	
	return &MSEMessageStreamEncryption{
		encryptionKey: encryptionKey,
		infoHash:      infoHash,
		mode:          MSEEncrypted,
	}
}

// SetMode устанавливает режим шифрования
func (e *MSEMessageStreamEncryption) SetMode(mode int) {
	e.mode = mode
}

// InitHandshake инициирует рукопожатие с шифрованием (outgoing connection)
func (e *MSEMessageStreamEncryption) InitHandshake(conn net.Conn) error {
	// Генерируем случайный S (20 bytes)
	e.localS = make([]byte, 20)
	if _, err := rand.Read(e.localS); err != nil {
		return err
	}

	// Отправляем S
	if _, err := conn.Write(e.localS); err != nil {
		return err
	}

	// Читаем S от пира
	e.remoteS = make([]byte, 20)
	if _, err := io.ReadFull(conn, e.remoteS); err != nil {
		return err
	}

	// Вычисляем RC4 ключи
	encryptKey := e.deriveKey(e.localS, e.remoteS)
	decryptKey := e.deriveKey(e.remoteS, e.localS)

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

	// Синхронизируем RC4 с VC (8 zero bytes)
	e.vc = [8]byte{}
	e.remoteVc = [8]byte{}
	
	// Discard initial RC4 keystream (VC synchronization)
	dummy := make([]byte, 8)
	e.encryptCipher.XORKeyStream(dummy, dummy)
	e.decryptCipher.XORKeyStream(dummy, dummy)

	return nil
}

// AcceptHandshake принимает рукопожатие (incoming connection)
func (e *MSEMessageStreamEncryption) AcceptHandshake(conn net.Conn) error {
	// Читаем S от пира
	e.remoteS = make([]byte, 20)
	if _, err := io.ReadFull(conn, e.remoteS); err != nil {
		return err
	}

	// Генерируем наш S
	e.localS = make([]byte, 20)
	if _, err := rand.Read(e.localS); err != nil {
		return err
	}

	// Отправляем наш S
	if _, err := conn.Write(e.localS); err != nil {
		return err
	}

	// Вычисляем RC4 ключи
	encryptKey := e.deriveKey(e.remoteS, e.localS)
	decryptKey := e.deriveKey(e.localS, e.remoteS)

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

	// Синхронизируем RC4 с VC
	dummy := make([]byte, 8)
	e.encryptCipher.XORKeyStream(dummy, dummy)
	e.decryptCipher.XORKeyStream(dummy, dummy)

	return nil
}

// deriveKey выводит ключ для RC4 (16 bytes)
func (e *MSEMessageStreamEncryption) deriveKey(localS, remoteS []byte) []byte {
	// K = MD5(S, S_other)
	hash := md5.New()
	hash.Write(localS)
	hash.Write(remoteS)
	return hash.Sum(nil)[:16]
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

	if n > 0 && ec.encryption.decryptCipher != nil {
		decrypted := make([]byte, n)
		ec.encryption.decryptCipher.XORKeyStream(decrypted, p[:n])
		copy(p, decrypted)
	}

	return n, nil
}

// Write пишет данные с шифрованием
func (ec *EncryptedConnection) Write(p []byte) (n int, err error) {
	if ec.encryption.encryptCipher != nil {
		encrypted := make([]byte, len(p))
		ec.encryption.encryptCipher.XORKeyStream(encrypted, p)
		return ec.conn.Write(encrypted)
	}
	return ec.conn.Write(p)
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

	// Проверяем что данные начинаются с "BitTorrent protocol" (17 символов)
	protocol := string(data[:19])
	return protocol == "BitTorrent protocol"
}

// TryEncryption пытается установить зашифрованное соединение
func TryEncryption(conn net.Conn, infoHash [20]byte) (net.Conn, *MSEMessageStreamEncryption, error) {
	mse := NewMSEEncryption(infoHash)
	
	// Пробуем encrypted handshake
	if err := mse.InitHandshake(conn); err != nil {
		// Не удалось установить шифрование, возвращаем обычное соединение
		return conn, nil, nil
	}
	
	// Создаём зашифрованное соединение
	encConn := NewEncryptedConnection(conn, mse)
	return encConn, mse, nil
}