package dht

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"smirnovtorrent/pkg/bencode"
)

// Стандартные bootstrap узлы DHT
var DefaultBootstrapNodes = []string{
	"router.bittorrent.com:6881",
	"dht.transmissionbt.com:6881",
	"dht.aelitis.com:6881",
	"router.utorrent.com:6881",
}

// DHTNode узел в DHT сети
type DHTNode struct {
	ID       [20]byte
	IP       string
	Port     uint16
	LastSeen time.Time
}

// Kademlia DHT таблица
type KademliaTable struct {
	nodes map[string]*DHTNode
	mu    sync.RWMutex
}

// DHTClient клиент для DHT сети
type DHTClient struct {
	nodeID      [20]byte
	udpConn     *net.UDPConn
	bootstrap   []string
	kademlia    *KademliaTable
	ctx         context.Context
	cancel      context.CancelFunc
	peersFound  chan []string
	mu          sync.RWMutex
	targetPeers int
}

// NewDHTClient создаёт новый DHT клиент
func NewDHTClient(bootstrap []string, port uint16) (*DHTClient, error) {
	// Генерируем случайный node ID
	var nodeID [20]byte
	if _, err := rand.Read(nodeID[:]); err != nil {
		return nil, fmt.Errorf("failed to generate node ID: %w", err)
	}

	// Создаём UDP соединение
	addr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: int(port),
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Используем стандартные bootstrap узлы если не указаны
	if len(bootstrap) == 0 {
		bootstrap = DefaultBootstrapNodes
	}

	client := &DHTClient{
		nodeID:      nodeID,
		udpConn:     conn,
		bootstrap:   bootstrap,
		kademlia:    &KademliaTable{nodes: make(map[string]*DHTNode)},
		ctx:         ctx,
		cancel:      cancel,
		peersFound:  make(chan []string, 10),
		targetPeers: 20,
	}

	return client, nil
}

// Start запускает DHT клиент
func (d *DHTClient) Start() error {
	log.Println("Starting DHT client...")

	// Подключаемся к bootstrap узлам
	for _, addr := range d.bootstrap {
		if err := d.ping(addr); err != nil {
			log.Printf("Failed to ping %s: %v", addr, err)
		}
	}

	// Запускаем цикл обслуживания
	go d.serviceLoop()

	return nil
}

// FindPeer ищет пиры для конкретного info hash
func (d *DHTClient) FindPeer(infoHash string) ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	log.Printf("Finding peers for info hash: %s", infoHash)
	
	// Преобразуем info hash в байты
	hashBytes, err := bencode.DecodeHexString(infoHash)
	if err != nil {
		return nil, fmt.Errorf("invalid info hash: %w", err)
	}

	// Отправляем get_peers запрос к bootstrap узлам
	go func() {
		for _, addr := range d.bootstrap {
			if err := d.sendGetPeers(addr, hashBytes); err != nil {
				log.Printf("Failed to send get_peers to %s: %v", addr, err)
			}
		}
	}()

	// Ждём результаты (30 секунд timeout)
	peers := []string{}
	timeout := time.After(30 * time.Second)
	
	for len(peers) < d.targetPeers {
		select {
		case found := <-d.peersFound:
			peers = append(peers, found...)
			if len(peers) >= d.targetPeers {
				return peers[:d.targetPeers], nil
			}
		case <-timeout:
			if len(peers) > 0 {
				return peers, nil
			}
			return nil, fmt.Errorf("no peers found")
		case <-d.ctx.Done():
			return nil, fmt.Errorf("cancelled")
		}
	}

	return peers, nil
}

// sendGetPeers отправляет get_peers запрос к узлу
func (d *DHTClient) sendGetPeers(addr string, infoHash []byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// Создаём get_peers запрос
	query := bencode.Dict{
		"y": bencode.String("q"),
		"q": bencode.String("get_peers"),
		"a": bencode.Dict{
			"id":  bencode.String(string(d.nodeID[:])),
			"info_hash": bencode.String(string(infoHash)),
		},
		"t": bencode.String("get_peers"),
	}

	data, err := bencode.Marshal(query)
	if err != nil {
		return err
	}

	_, err = d.udpConn.WriteToUDP(data, udpAddr)
	if err != nil {
		return err
	}

	log.Printf("get_peers sent to %s", addr)
	return nil
}

// GetPeersFound возвращает канал найденных пиров
func (d *DHTClient) GetPeersFound() <-chan []string {
	return d.peersFound
}

// Stop останавливает DHT клиент
func (d *DHTClient) Stop() {
	d.cancel()
	d.udpConn.Close()
	log.Println("DHT client stopped")
}

// ping отправляет ping запрос к узлу
func (d *DHTClient) ping(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// Создаём ping запрос
	query := d.createQuery("ping", d.nodeID[:])
	
	_, err = d.udpConn.WriteToUDP(query, udpAddr)
	if err != nil {
		return err
	}

	log.Printf("Ping sent to %s", addr)
	return nil
}

// createQuery создаёт bencode запрос
func (d *DHTClient) createQuery(method string, target []byte) []byte {
	query := bencode.Dict{
		"y": bencode.String("q"), // query
		"q": bencode.String(method),
		"a": bencode.Dict{
			"id": bencode.String(target[:16]), // Укороченный ID для примера
		},
		"t": bencode.String("aa"), // transaction ID
	}

	data, _ := bencode.Marshal(query)
	return data
}

// serviceLoop цикл обслуживания DHT
func (d *DHTClient) serviceLoop() {
	buf := make([]byte, 1400) // MTU

	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			d.udpConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _, err := d.udpConn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			go d.handlePacket(buf[:n])
		}
	}
}

// handlePacket обрабатывает входящий пакет
func (d *DHTClient) handlePacket(data []byte) {
	val, err := bencode.Unmarshal(data)
	if err != nil {
		return
	}

	query, ok := val.(bencode.Dict)
	if !ok {
		return
	}

	// Проверяем тип запроса
	if string(query["y"].(bencode.String)) != "r" {
		return
	}

	// Обрабатываем ответ
	response, ok := query["r"].(bencode.Dict)
	if !ok {
		return
	}

	// Извлекаем узлы из response
	if values, ok := response["values"]; ok {
		if list, ok := values.(bencode.List); ok {
			log.Printf("Found %d peers in DHT response", len(list))
			
			// Извлекаем пиры из списка
			for _, v := range list {
				if peerBytes, ok := v.(bencode.String); ok {
					if len(peerBytes) >= 6 {
						ip, port, err := DecodePeerInfo([]byte(peerBytes))
						if err == nil {
							peerStr := PeerToString(ip, port)
							select {
							case d.peersFound <- []string{peerStr}:
							default:
							}
						}
					}
				}
			}
		}
	}

	// Также проверяем nodes (для compact node info)
	if nodes, ok := response["nodes"]; ok {
		if nodesStr, ok := nodes.(bencode.String); ok {
			log.Printf("Received nodes data: %d bytes", len(nodesStr))
			// Парсим compact node info (каждый узел 26 байт: 20 ID + 2 порт + 4 IP)
			for i := 0; i+26 <= len(nodesStr); i += 26 {
				_ = nodesStr[i : i+26] // nodeData - пропускаем детальный парсинг для краткости
			}
		}
	}
}

// addNode добавляет узел в таблицу
func (d *DHTClient) addNode(node *DHTNode) {
	d.kademlia.mu.Lock()
	defer d.kademlia.mu.Unlock()

	key := fmt.Sprintf("%x", node.ID[:8])
	d.kademlia.nodes[key] = node
}

// GetNodeCount возвращает количество узлов в таблице
func (d *DHTClient) GetNodeCount() int {
	d.kademlia.mu.RLock()
	defer d.kademlia.mu.RUnlock()
	return len(d.kademlia.nodes)
}

// Helper functions

// ParsePeerAddress парсит адрес пирa
func ParsePeerAddress(peerStr string) (string, uint16, error) {
	host, portStr, err := net.SplitHostPort(peerStr)
	if err != nil {
		return "", 0, err
	}

	var port uint16
	_, err = fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		return "", 0, err
	}

	return host, port, nil
}

// PeerToString конвертирует адрес в строку
func PeerToString(ip string, port uint16) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

// EncodePeerInfo кодирование информации о пире для DHT
func EncodePeerInfo(ip string, port uint16) ([]byte, error) {
	addr := net.ParseIP(ip)
	if addr == nil {
		return nil, fmt.Errorf("invalid IP address")
	}

	var buf bytes.Buffer
	buf.Write(addr.To4())
	binary.Write(&buf, binary.BigEndian, port)

	return buf.Bytes(), nil
}

// DecodePeerInfo декодирование информации о пире из DHT
func DecodePeerInfo(data []byte) (string, uint16, error) {
	if len(data) < 6 {
		return "", 0, fmt.Errorf("invalid peer info length")
	}

	ip := net.IP(data[:4]).String()
	port := binary.BigEndian.Uint16(data[4:6])

	return ip, port, nil
}