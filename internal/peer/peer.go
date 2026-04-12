package peer

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// PeerID размер peer ID в байтах
const PeerIDSize = 20

// ProtocolName имя протокола BitTorrent
const ProtocolName = "BitTorrent protocol"

// PeerIDPrefix префикс для нашего peer ID
const PeerIDPrefix = "-SMRV0100-"

// Message types
const (
	MsgChoke        = 0
	MsgUnchoke      = 1
	MsgInterested   = 2
	MsgNotInterested = 3
	MsgHave         = 4
	MsgBitfield     = 5
	MsgRequest      = 6
	MsgPiece        = 7
	MsgCancel       = 8
)

// Peer представляет удалённый пир
type Peer struct {
	IP       string
	Port     uint16
	PeerID   [PeerIDSize]byte
	Bitfield []byte
	Choked   bool
	Interested bool
}

// PeerConnection активное соединение с пиром
type PeerConnection struct {
	Peer
	Conn net.Conn
}

// NewPeerID создаёт уникальный peer ID
func NewPeerID() [PeerIDSize]byte {
	var peerID [PeerIDSize]byte
	copy(peerID[:], PeerIDPrefix)
	// Оставшуюся часть заполняем случайными данными
	for i := len(PeerIDPrefix); i < PeerIDSize; i++ {
		peerID[i] = byte(i)
	}
	return peerID
}

// Connect устанавливает соединение с пиром
func (p *Peer) Connect() (*PeerConnection, error) {
	addr := fmt.Sprintf("%s:%d", p.IP, p.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*1000000000) // 5 секунд
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &PeerConnection{
		Peer: *p,
		Conn: conn,
	}, nil
}

// SendHandshake отправляет handshake пиру
func (pc *PeerConnection) SendHandshake(infoHash [20]byte, peerID [PeerIDSize]byte) error {
	var buf [68]byte
	buf[0] = 19 // длина имени протокола
	copy(buf[1:20], []byte(ProtocolName))
	copy(buf[20:40], infoHash[:])
	copy(buf[40:60], peerID[:])

	_, err := pc.Conn.Write(buf[:])
	return err
}

// ReadHandshake читает handshake от пирa
func (pc *PeerConnection) ReadHandshake() ([20]byte, [PeerIDSize]byte, error) {
	var buf [68]byte
	_, err := io.ReadFull(pc.Conn, buf[:])
	if err != nil {
		return [20]byte{}, [PeerIDSize]byte{}, fmt.Errorf("failed to read handshake: %w", err)
	}

	// Проверяем длину протокола
	if buf[0] != 19 {
		return [20]byte{}, [PeerIDSize]byte{}, fmt.Errorf("invalid protocol length: %d", buf[0])
	}

	// Проверяем имя протокола
	protocol := string(buf[1:20])
	if protocol != ProtocolName {
		return [20]byte{}, [PeerIDSize]byte{}, fmt.Errorf("invalid protocol: %s", protocol)
	}

	var infoHash [20]byte
	var peerID [PeerIDSize]byte
	copy(infoHash[:], buf[20:40])
	copy(peerID[:], buf[40:60])

	return infoHash, peerID, nil
}

// SendMessage отправляет сообщение пиру
func (pc *PeerConnection) SendMessage(msgType byte, payload []byte) error {
	length := uint32(len(payload)) + 1 // +1 для типа сообщения

	// Заголовок: длина (4 байта) + тип сообщения
	header := make([]byte, 5)
	binary.BigEndian.PutUint32(header[:4], length)
	header[4] = msgType

	_, err := pc.Conn.Write(header)
	if err != nil {
		return err
	}

	if len(payload) > 0 {
		_, err = pc.Conn.Write(payload)
	}

	return err
}

// ReadMessage читает сообщение от пирa
func (pc *PeerConnection) ReadMessage() (byte, []byte, error) {
	// Читаем длину (4 байта)
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(pc.Conn, lengthBuf)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read length: %w", err)
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		// Keep-alive message
		return 0, nil, nil
	}

	// Читаем тип сообщения и payload
	msgBuf := make([]byte, length)
	_, err = io.ReadFull(pc.Conn, msgBuf)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read message: %w", err)
	}

	msgType := msgBuf[0]
	payload := msgBuf[1:]

	return msgType, payload, nil
}

// SendChoke отправляет choke сообщение
func (pc *PeerConnection) SendChoke() error {
	return pc.SendMessage(MsgChoke, nil)
}

// SendUnchoke отправляет unchoke сообщение
func (pc *PeerConnection) SendUnchoke() error {
	return pc.SendMessage(MsgUnchoke, nil)
}

// SendInterested отправляет interested сообщение
func (pc *PeerConnection) SendInterested() error {
	return pc.SendMessage(MsgInterested, nil)
}

// SendNotInterested отправляет not interested сообщение
func (pc *PeerConnection) SendNotInterested() error {
	return pc.SendMessage(MsgNotInterested, nil)
}

// SendHave отправляет have сообщение
func (pc *PeerConnection) SendHave(pieceIndex uint32) error {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, pieceIndex)
	return pc.SendMessage(MsgHave, payload)
}

// SendRequest отправляет request сообщение
func (pc *PeerConnection) SendRequest(pieceIndex, begin, length uint32) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], pieceIndex)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	binary.BigEndian.PutUint32(payload[8:12], length)
	return pc.SendMessage(MsgRequest, payload)
}

// SendPiece отправляет piece сообщение
func (pc *PeerConnection) SendPiece(pieceIndex, begin uint32, data []byte) error {
	payload := make([]byte, 8+len(data))
	binary.BigEndian.PutUint32(payload[0:4], pieceIndex)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	copy(payload[8:], data)
	return pc.SendMessage(MsgPiece, payload)
}

// SendCancel отправляет cancel сообщение
func (pc *PeerConnection) SendCancel(pieceIndex, begin, length uint32) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], pieceIndex)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	binary.BigEndian.PutUint32(payload[8:12], length)
	return pc.SendMessage(MsgCancel, payload)
}

// ReadBitfield читает bitfield от пирa
func (pc *PeerConnection) ReadBitfield() ([]byte, error) {
	msgType, payload, err := pc.ReadMessage()
	if err != nil {
		return nil, err
	}

	if msgType != MsgBitfield {
		return nil, fmt.Errorf("expected bitfield message, got %d", msgType)
	}

	return payload, nil
}

// Close закрывает соединение
func (pc *PeerConnection) Close() error {
	return pc.Conn.Close()
}