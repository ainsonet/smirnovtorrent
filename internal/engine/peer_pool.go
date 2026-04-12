package engine

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"smirnovtorrent/internal/peer"
	"smirnovtorrent/internal/tracker"
)

// PeerPool пул подключений к пирам
type PeerPool struct {
	peers       map[string]*PeerSlot
	mu          sync.RWMutex
	infoHash    string
	peerID      string
	port        uint16
	maxPeers    int
	connected   int
	connectedMu sync.Mutex
	pexClient   *peer.PEXClient
	usePEX      bool
}

// PeerSlot состояние одного пирa в пуле
type PeerSlot struct {
	PeerID       string
	Conn         *peer.PeerConnection
	HasPieces    []bool // какие куски есть у пирa
	Downloading  bool   // качаем ли сейчас у этого пирa
	LastActive   time.Time
	Choked       bool
	Interested   bool
}

// NewPeerPool создаёт новый пул пиров
func NewPeerPool(infoHash string, peerID string, port uint16, maxPeers int) *PeerPool {
	pool := &PeerPool{
		peers:    make(map[string]*PeerSlot),
		infoHash: infoHash,
		peerID:   peerID,
		port:     port,
		maxPeers: maxPeers,
	}
	
	// Инициализируем PEX клиент
	infoHashBytes, _ := hex.DecodeString(infoHash)
	peerIDBytes := []byte(peerID)
	var infoHashArr, peerIDArr [20]byte
	copy(infoHashArr[:], infoHashBytes)
	copy(peerIDArr[:], peerIDBytes)
	pool.pexClient = peer.NewPEXClient(peerIDArr, infoHashArr)

	return pool
}

// EnablePEX включает Peer Exchange
func (pp *PeerPool) EnablePEX() {
	pp.usePEX = true
}

// SendPEX отправляет PEX сообщения всем пирам
func (pp *PeerPool) SendPEX() {
	if !pp.usePEX {
		return
	}

	pp.mu.RLock()
	defer pp.mu.RUnlock()

	pexData, err := pp.pexClient.CreatePEXMessage()
	if err != nil {
		log.Printf("PEX: failed to create message: %v", err)
		return
	}

	// Отправляем extended message с PEX данным
	for _, slot := range pp.peers {
		if slot.Conn != nil {
			// Extended message type 20, ut_pex = 1
			pp.sendExtendedMessage(slot.Conn, 1, pexData)
		}
	}

	log.Printf("PEX: sent updates to %d peers", len(pp.peers))
}

// AddDiscoveredPeer добавляет пирa обнаруженного через PEX
func (pp *PeerPool) AddDiscoveredPeer(ip string, port uint16) {
	if !pp.usePEX {
		return
	}
	
	pp.pexClient.AddPeer(ip, port)
	
	// Пробуем подключиться к новому пиру
	peerInfo := tracker.PeerInfo{
		IP:   ip,
		Port: port,
	}
	pp.AddPeer(peerInfo)
}

// RemoveDisconnectedPeer удаляет отключенного пирa из PEX
func (pp *PeerPool) RemoveDisconnectedPeer(ip string, port uint16) {
	if !pp.usePEX {
		return
	}
	
	pp.pexClient.RemovePeer(ip, port)
}

// sendExtendedMessage отправляет extended сообщение
func (pp *PeerPool) sendExtendedMessage(conn *peer.PeerConnection, extMsgID byte, data []byte) error {
	// Extended message format: length (4) + msg_type (1=20) + ext_msgid (1) + payload
	payload := append([]byte{extMsgID}, data...)
	
	// Длина: 1 (msg_type=20) + len(payload)
	totalLen := uint32(1 + len(payload))
	
	// Заголовок: длина (4 байта big-endian)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, totalLen)
	
	if _, err := conn.Conn.Write(header); err != nil {
		return err
	}
	
	// Тип сообщения extended (20)
	if _, err := conn.Conn.Write([]byte{20}); err != nil {
		return err
	}

	// Payload
	_, err := conn.Conn.Write(payload)
	return err
}

// AddPeer добавляет пирa в пул
func (pp *PeerPool) AddPeer(p tracker.PeerInfo) error {
	pp.connectedMu.Lock()
	if pp.connected >= pp.maxPeers {
		pp.connectedMu.Unlock()
		return nil // Достигнут лимит пиров
	}
	pp.connectedMu.Unlock()

	peerObj := &peer.Peer{
		IP:     p.IP,
		Port:   p.Port,
		PeerID: peer.NewPeerID(),
	}

	conn, err := peerObj.Connect()
	if err != nil {
		log.Printf("Failed to connect to %s:%d: %v", p.IP, p.Port, err)
		return err
	}

	// Отправляем handshake
	infoHashBytes, _ := hex.DecodeString(pp.infoHash)
	var infoHashArr [20]byte
	copy(infoHashArr[:], infoHashBytes)

	peerIDBytes := []byte(pp.peerID)
	var peerIDArr [20]byte
	copy(peerIDArr[:], peerIDBytes)

	if err := conn.SendHandshake(infoHashArr, peerIDArr); err != nil {
		conn.Close()
		return err
	}

	// Читаем handshake от пира
	remoteInfoHash, remotePeerID, err := conn.ReadHandshake()
	if err != nil {
		conn.Close()
		return err
	}

	// Проверяем info hash
	if remoteInfoHash != infoHashArr {
		conn.Close()
		return fmt.Errorf("info hash mismatch")
	}

	// Читаем bitfield
	var bitfield []byte
	msgType, payload, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return err
	}

	if msgType == peer.MsgBitfield {
		bitfield = payload
	}

	// Создаём массив наличие кусков
	numPieces := (len(bitfield) * 8)
	hasPieces := make([]bool, numPieces)
	for i := 0; i < numPieces; i++ {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8)
		if byteIdx < len(bitfield) {
			hasPieces[i] = (bitfield[byteIdx] & (1 << bitIdx)) != 0
		}
	}

	peerSlot := &PeerSlot{
		PeerID:    string(remotePeerID[:]),
		Conn:      conn,
		HasPieces: hasPieces,
		LastActive: time.Now(),
	}

	pp.mu.Lock()
	key := net.JoinHostPort(p.IP, fmt.Sprintf("%d", p.Port))
	pp.peers[key] = peerSlot
	pp.mu.Unlock()

	pp.connectedMu.Lock()
	pp.connected++
	pp.connectedMu.Unlock()

	log.Printf("Connected to peer %s (%d pieces)", remotePeerID[:8], countBits(hasPieces))

	// Отправляем interested
	conn.SendInterested()

	return nil
}

// GetPeer возвращает активного пирa который не choked
func (pp *PeerPool) GetPeer() *PeerSlot {
	pp.mu.RLock()
	defer pp.mu.RUnlock()

	for _, slot := range pp.peers {
		if !slot.Choked && !slot.Downloading && slot.Interested {
			return slot
		}
	}

	return nil
}

// GetPeersWithPiece возвращает пиров у которых есть конкретный кусок
func (pp *PeerPool) GetPeersWithPiece(pieceIndex int) []*PeerSlot {
	pp.mu.RLock()
	defer pp.mu.RUnlock()

	var result []*PeerSlot
	for _, slot := range pp.peers {
		if pieceIndex < len(slot.HasPieces) && slot.HasPieces[pieceIndex] {
			result = append(result, slot)
		}
	}

	return result
}

// MarkPieceDownloading помечает что кусок качается у пирa
func (pp *PeerPool) MarkPieceDownloading(slot *PeerSlot) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if slot != nil {
		slot.Downloading = true
	}
}

// MarkPieceDone помечает что кусок загружен у пирa
func (pp *PeerPool) MarkPieceDone(slot *PeerSlot) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if slot != nil {
		slot.Downloading = false
		slot.LastActive = time.Now()
	}
}

// UpdatePeerBitfield обновляет bitfield пирa (после have сообщения)
func (pp *PeerPool) UpdatePeerBitfield(peerKey string, pieceIndex int, hasPiece bool) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if slot, ok := pp.peers[peerKey]; ok {
		if pieceIndex < len(slot.HasPieces) {
			slot.HasPieces[pieceIndex] = hasPiece
		}
		slot.LastActive = time.Now()
	}
}

// GetActivePeerCount возвращает количество активных пиров
func (pp *PeerPool) GetActivePeerCount() int {
	pp.connectedMu.Lock()
	defer pp.connectedMu.Unlock()
	return pp.connected
}

// CloseAll закрывает все соединения
func (pp *PeerPool) CloseAll() {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	for key, slot := range pp.peers {
		slot.Conn.Close()
		delete(pp.peers, key)
	}

	pp.connected = 0
}

// Helper functions

func countBits(bits []bool) int {
	count := 0
	for _, b := range bits {
		if b {
			count++
		}
	}
	return count
}
