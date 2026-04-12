package peer

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"smirnovtorrent/pkg/bencode"
)

// PEXClient клиент для Peer Exchange (BEP 11)
type PEXClient struct {
	peerID     [20]byte
	infoHash   [20]byte
	peers      map[string]PeerInfo
	added      []PeerInfo
	dropped    []PeerInfo
	mu         sync.RWMutex
	lastUpdate time.Time
}

// PeerInfo информация о пире для PEX
type PeerInfo struct {
	IP   string
	Port uint16
}

// NewPEXClient создаёт новый PEX клиент
func NewPEXClient(peerID, infoHash [20]byte) *PEXClient {
	return &PEXClient{
		peerID:   peerID,
		infoHash: infoHash,
		peers:    make(map[string]PeerInfo),
		added:    make([]PeerInfo, 0),
		dropped:  make([]PeerInfo, 0),
	}
}

// AddPeer добавляет пирa в список
func (p *PEXClient) AddPeer(ip string, port uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := fmt.Sprintf("%s:%d", ip, port)
	if _, exists := p.peers[key]; !exists {
		p.peers[key] = PeerInfo{IP: ip, Port: port}
		p.added = append(p.added, PeerInfo{IP: ip, Port: port})
	}
}

// RemovePeer удаляет пирa из списка
func (p *PEXClient) RemovePeer(ip string, port uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := fmt.Sprintf("%s:%d", ip, port)
	if _, exists := p.peers[key]; exists {
		delete(p.peers, key)
		p.dropped = append(p.dropped, PeerInfo{IP: ip, Port: port})
	}
}

// CreatePEXMessage создаёт PEX сообщение для отправки
func (p *PEXClient) CreatePEXMessage() ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Формируем compact peer list для добавленных
	added4 := p.encodePeersIPv4(p.added)
	
	// Формируем compact peer list для удалённых
	dropped4 := p.encodePeersIPv4(p.dropped)

	// Создаём bencode сообщение
	msg := bencode.Dict{
		"added":  bencode.String(added4),
		"dropped": bencode.String(dropped4),
	}

	// Добавляем IPv6 если есть (пока пусто)
	msg["added6"] = bencode.String("")
	msg["dropped6"] = bencode.String("")

	data, err := bencode.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Очищаем списки после отправки
	p.added = make([]PeerInfo, 0)
	p.dropped = make([]PeerInfo, 0)
	p.lastUpdate = time.Now()

	return data, nil
}

// encodePeersIPv4 кодирует список пиров в compact формат
func (p *PEXClient) encodePeersIPv4(peers []PeerInfo) []byte {
	// Каждый пир: 4 байта IP + 2 байта порт = 6 байт
	data := make([]byte, 0, len(peers)*6)
	
	for _, peer := range peers {
		ip := net.ParseIP(peer.IP)
		if ip == nil {
			continue
		}
		
		ip4 := ip.To4()
		if ip4 == nil {
			continue
		}
		
		data = append(data, ip4...)
		
		portBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(portBytes, peer.Port)
		data = append(data, portBytes...)
	}
	
	return data
}

// ParsePEXMessage парсит полученное PEX сообщение
func (p *PEXClient) ParsePEXMessage(data []byte) ([]PeerInfo, []PeerInfo, error) {
	val, err := bencode.Unmarshal(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse PEX message: %w", err)
	}

	dict, ok := val.(bencode.Dict)
	if !ok {
		return nil, nil, fmt.Errorf("invalid PEX message format")
	}

	// Парсим добавленных пиров (IPv4)
	added := p.parsePeerList(dict, "added")
	
	// Парсим удалённых пиров (IPv4)
	dropped := p.parsePeerList(dict, "dropped")

	// Парсим IPv6 пиров (если есть)
	added6 := p.parsePeerList(dict, "added6")
	dropped6 := p.parsePeerList(dict, "dropped6")

	added = append(added, added6...)
	dropped = append(dropped, dropped6...)

	return added, dropped, nil
}

// parsePeerList парсит список пиров из bencode dict
func (p *PEXClient) parsePeerList(dict bencode.Dict, key string) []PeerInfo {
	val, exists := dict[key]
	if !exists {
		return nil
	}

	peerBytes, ok := val.(bencode.String)
	if !ok {
		return nil
	}

	return p.decodePeers([]byte(peerBytes))
}

// decodePeers декодирует compact peer list
func (p *PEXClient) decodePeers(data []byte) []PeerInfo {
	peers := make([]PeerInfo, 0)
	
	// IPv4: 6 байт на пирa
	for i := 0; i+6 <= len(data); i += 6 {
		ip := net.IP(data[i : i+4]).String()
		port := binary.BigEndian.Uint16(data[i+4 : i+6])
		peers = append(peers, PeerInfo{IP: ip, Port: port})
	}
	
	// IPv6: 18 байт на пирa (16 IP + 2 порт)
	for i := 0; i+18 <= len(data); i += 18 {
		ip := net.IP(data[i : i+16]).String()
		port := binary.BigEndian.Uint16(data[i+16 : i+18])
		peers = append(peers, PeerInfo{IP: ip, Port: port})
	}
	
	return peers
}

// GetPeers возвращает текущий список пиров
func (p *PEXClient) GetPeers() []PeerInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	peers := make([]PeerInfo, 0, len(p.peers))
	for _, peer := range p.peers {
		peers = append(peers, peer)
	}
	return peers
}

// GetPeerCount возвращает количество пиров
func (p *PEXClient) GetPeerCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.peers)
}

// ShouldUpdate проверяет нужно ли отправлять обновление
func (p *PEXClient) ShouldUpdate(interval time.Duration) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return time.Since(p.lastUpdate) > interval || len(p.added) > 0 || len(p.dropped) > 0
}

// Reset сбрасывает PEX состояние
func (p *PEXClient) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.added = make([]PeerInfo, 0)
	p.dropped = make([]PeerInfo, 0)
	p.lastUpdate = time.Now()
}
