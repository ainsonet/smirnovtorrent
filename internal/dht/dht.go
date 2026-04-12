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

	client := &DHTClient{
		nodeID:     nodeID,
		udpConn:    conn,
		bootstrap:  bootstrap,
		kademlia:   &KademliaTable{nodes: make(map[string]*DHTNode)},
		ctx:        ctx,
		cancel:     cancel,
		peersFound: make(chan []string, 10),
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

	// В реальной реализации здесь был бы запрос Kademlia DHT
	// Пока возвращаем заглушку
	log.Printf("Finding peers for info hash: %s", infoHash)
	
	// Эмулируем поиск (в реальности это async запрос)
	go func() {
		time.Sleep(1 * time.Second)
		// Возвращаем несколько фейковых пиров для демонстрации
		peers := []string{
			"192.168.1.100:6881",
			"10.0.0.50:6882",
		}
		select {
		case d.peersFound <- peers:
		default:
		}
	}()

	// Ждём результаты
	select {
	case peers := <-d.peersFound:
		return peers, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout finding peers")
	}
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

	// Извлекаем узлы из ответа
	if nodes, ok := response["values"]; ok {
		if list, ok := nodes.(bencode.List); ok {
			log.Printf("Found %d peers in DHT response", len(list))
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