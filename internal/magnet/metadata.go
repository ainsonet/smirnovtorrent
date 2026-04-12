package magnet

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"smirnovtorrent/internal/dht"
	"smirnovtorrent/internal/peer"
	"smirnovtorrent/pkg/bencode"
)

// MetadataDownloader загрузчик метаданных через DHT (BEP 9)
type MetadataDownloader struct {
	infoHash   [20]byte
	dhtClient  *dht.DHTClient
	ctx        context.Context
	cancel     context.CancelFunc
	metadata   []byte
	mu         sync.RWMutex
	pieceSize  int // 16KB для metadata transfer
	numPieces  int
	received   map[int][]byte
}

// NewMetadataDownloader создаёт загрузчик метаданных
func NewMetadataDownloader(infoHash string, dhtClient *dht.DHTClient) (*MetadataDownloader, error) {
	hashBytes, err := hex.DecodeString(infoHash)
	if err != nil {
		return nil, fmt.Errorf("invalid info hash: %w", err)
	}

	var hashArray [20]byte
	copy(hashArray[:], hashBytes)

	ctx, cancel := context.WithCancel(context.Background())

	return &MetadataDownloader{
		infoHash:  hashArray,
		dhtClient: dhtClient,
		ctx:       ctx,
		cancel:    cancel,
		pieceSize: 16 * 1024, // 16KB
		received:  make(map[int][]byte),
	}, nil
}

// Download загружает метаданные .torrent файла
func (m *MetadataDownloader) Download() ([]byte, error) {
	log.Printf("Downloading metadata for info hash: %x", m.infoHash)

	// Находим пиры через DHT
	peers, err := m.dhtClient.FindPeer(hex.EncodeToString(m.infoHash[:]))
	if err != nil {
		return nil, fmt.Errorf("failed to find peers: %w", err)
	}

	log.Printf("Found %d peers for metadata download", len(peers))

	if len(peers) == 0 {
		return nil, fmt.Errorf("no peers available")
	}

	// Подключаемся к пирам и запрашиваем metadata
	return m.downloadFromPeers(peers)
}

// downloadFromPeers загружает метаданные от пиров
func (m *MetadataDownloader) downloadFromPeers(peers []string) ([]byte, error) {
	resultCh := make(chan []byte, 1)
	errCh := make(chan error, len(peers))

	// Пытаемся подключиться к нескольким пирам параллельно
	maxConcurrent := 3
	if len(peers) < maxConcurrent {
		maxConcurrent = len(peers)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

	for _, peerAddr := range peers {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			metadata, err := m.tryPeer(addr)
			if err != nil {
				errCh <- err
				return
			}

			// Успешно получили метаданные
			select {
			case resultCh <- metadata:
			default:
			}
		}(peerAddr)
	}

	// Ждём результат или таймаут
	timeout := time.After(60 * time.Second)
	select {
	case metadata := <-resultCh:
		m.cancel()
		log.Printf("Metadata downloaded successfully: %d bytes", len(metadata))
		return metadata, nil
	case err := <-errCh:
		log.Printf("Peer error: %v", err)
		// Продолжаем ждать другие пиры
	case <-timeout:
		m.cancel()
		return nil, fmt.Errorf("metadata download timeout")
	case <-m.ctx.Done():
		return nil, fmt.Errorf("cancelled")
	}

	// Ждём пока остальные попытки завершатся
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case metadata := <-resultCh:
		return metadata, nil
	case <-done:
		return nil, fmt.Errorf("all peers failed")
	case <-timeout:
		return nil, fmt.Errorf("metadata download timeout")
	}
}

// tryPeer пытается получить метаданные от одного пира
func (m *MetadataDownloader) tryPeer(peerAddr string) ([]byte, error) {
	// Парсим адрес
	host, portStr, err := net.SplitHostPort(peerAddr)
	if err != nil {
		return nil, err
	}
	
	var port uint16
	fmt.Sscanf(portStr, "%d", &port)

	// Подключаемся к пиру
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	// Отправляем handshake
	peerID := peer.NewPeerID()
	handshake := createHandshake(m.infoHash[:], peerID[:])
	
	if _, err := conn.Write(handshake); err != nil {
		return nil, fmt.Errorf("handshake failed: %w", err)
	}

	// Читаем ответный handshake
	response := make([]byte, 68)
	if _, err := conn.Read(response); err != nil {
		return nil, fmt.Errorf("handshake response failed: %w", err)
	}

	// Проверяем handshake
	if response[0] != 19 || string(response[1:20]) != "BitTorrent protocol" {
		return nil, fmt.Errorf("invalid handshake response")
	}

	// Проверяем поддержку extension protocol (BEP 9)
	supportsExtensions := false
	for _, b := range response[20:28] {
		if b != 0 {
			supportsExtensions = true
			break
		}
	}

	if !supportsExtensions {
		return nil, fmt.Errorf("peer doesn't support extensions")
	}

	// Отправляем extended handshake для получения metadata size
	extHandshake := m.createExtendedHandshake()
	if err := m.sendMessage(conn, extHandshake); err != nil {
		return nil, err
	}

	// Читаем extended handshake от пира
	metadataSize, utMetadata, err := m.readExtendedHandshake(conn)
	if err != nil {
		return nil, err
	}

	if metadataSize == 0 {
		return nil, fmt.Errorf("invalid metadata size")
	}

	log.Printf("Metadata size: %d bytes", metadataSize)

	// Вычисляем количество кусков
	m.numPieces = (metadataSize + m.pieceSize - 1) / m.pieceSize
	m.received = make(map[int][]byte)

	// Запрашиваем все куски метаданных
	for i := 0; i < m.numPieces; i++ {
		if err := m.requestMetadataPiece(conn, utMetadata, i); err != nil {
			return nil, fmt.Errorf("failed to request piece %d: %w", i, err)
		}
	}

	// Собираем метаданные
	return m.assembleMetadata(metadataSize)
}

// createHandshake создаёт BitTorrent handshake
func createHandshake(infoHash, peerID []byte) []byte {
	var buf [68]byte
	buf[0] = 19 // длина имени протокола
	copy(buf[1:20], []byte("BitTorrent protocol"))
	copy(buf[20:40], infoHash[:20])
	copy(buf[40:60], peerID[:20])
	return buf[:]
}

// createExtendedHandshake создаёт extended handshake
func (m *MetadataDownloader) createExtendedHandshake() []byte {
	// BEP 9 extended handshake
	metadata := bencode.Dict{
		"m": bencode.Dict{
			"ut_metadata": bencode.Int(1),
		},
		"metadata_size": bencode.Int(0),
		"p":             bencode.Int(6881),
		"v":             bencode.String("SmirnovTorrent 0.9.0"),
	}

	data, _ := bencode.Marshal(metadata)

	// Extended message: msg_type=0 (handshake), ext_header=bencode
	msg := make([]byte, 2)
	msg[0] = 0 // extended message type
	msg[1] = 0 // extended handshake

	msg = append(msg, data...)

	// Prefix with message length
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(msg)))
	
	return append(length, msg...)
}

// readExtendedHandshake читает extended handshake от пира
func (m *MetadataDownloader) readExtendedHandshake(conn net.Conn) (int, int, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// Читаем длину сообщения
	lengthBuf := make([]byte, 4)
	if _, err := conn.Read(lengthBuf); err != nil {
		return 0, 0, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	if length > 1024*1024 { // 1MB max
		return 0, 0, fmt.Errorf("extended handshake too large")
	}

	// Читаем сообщение
	data := make([]byte, length)
	if _, err := conn.Read(data); err != nil {
		return 0, 0, err
	}

	if len(data) < 2 || data[0] != 0 || data[1] != 0 {
		return 0, 0, fmt.Errorf("invalid extended message")
	}

	// Парсим bencode
	val, err := bencode.Unmarshal(data[2:])
	if err != nil {
		return 0, 0, err
	}

	dict, ok := val.(bencode.Dict)
	if !ok {
		return 0, 0, fmt.Errorf("invalid extended handshake format")
	}

	// Извлекаем metadata_size
	var metadataSize int
	var utMetadata int

	if size, ok := dict["metadata_size"]; ok {
		if sizeInt, ok := size.(bencode.Int); ok {
			metadataSize = int(sizeInt)
		}
	}

	if m, ok := dict["m"]; ok {
		if mDict, ok := m.(bencode.Dict); ok {
			if utm, ok := mDict["ut_metadata"]; ok {
				if utmInt, ok := utm.(bencode.Int); ok {
					utMetadata = int(utmInt)
				}
			}
		}
	}

	return metadataSize, utMetadata, nil
}

// requestMetadataPiece запрашивает кусок метаданных
func (m *MetadataDownloader) requestMetadataPiece(conn net.Conn, utMetadata, piece int) error {
	// BEP 9 metadata request message
	msg := bencode.Dict{
		"msg_type": bencode.Int(0), // request
		"piece":    bencode.Int(piece),
	}

	data, _ := bencode.Marshal(msg)

	// Extended message: msg_type=ext, ext_header=bencode
	extMsg := make([]byte, 2)
	extMsg[0] = 20 // extended message type
	extMsg[1] = byte(utMetadata)

	extMsg = append(extMsg, data...)

	// Prefix with message length
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(extMsg)))

	fullMsg := append(length, extMsg...)

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := conn.Write(fullMsg)
	return err
}

// assembleMetadata собирает метаданные из кусков
func (m *MetadataDownloader) assembleMetadata(totalSize int) ([]byte, error) {
	timeout := time.After(30 * time.Second)
	
	for len(m.received) < m.numPieces {
		select {
		case <-timeout:
			return nil, fmt.Errorf("metadata assembly timeout")
		case <-m.ctx.Done():
			return nil, fmt.Errorf("cancelled")
		default:
			// Читаем ответ от пира
			// (в реальной реализации здесь нужно читать из соединения)
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Собираем куски
	metadata := make([]byte, 0, totalSize)
	for i := 0; i < m.numPieces; i++ {
		if piece, ok := m.received[i]; ok {
			metadata = append(metadata, piece...)
		}
	}

	// Проверяем hash
	hash := sha1.Sum(metadata)
	if hash != m.infoHash {
		return nil, fmt.Errorf("metadata hash mismatch")
	}

	return metadata, nil
}

// sendMessage отправляет сообщение пиру
func (m *MetadataDownloader) sendMessage(conn net.Conn, msg []byte) error {
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := conn.Write(msg)
	return err
}
