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

// KBucketSize размер K-bucket (параметр Kademlia)
const KBucketSize = 20

// RPC_timeout таймаут для RPC запросов
const RPCTimeout = 5 * time.Second

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

// KBucket K-bucket для хранения узлов
type KBucket struct {
	nodes []*DHTNode
	mu    sync.RWMutex
}

// KademliaTable DHT таблица с K-buckets
type KademliaTable struct {
	nodeID  [20]byte
	buckets [160]*KBucket // 160 бит для SHA-1 ID
	mu      sync.RWMutex
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
	transactionID uint32
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
		kademlia:    NewKademliaTable(nodeID),
		ctx:         ctx,
		cancel:      cancel,
		peersFound:  make(chan []string, 10),
		targetPeers: 20,
		transactionID: 0,
	}

	return client, nil
}

// NewKademliaTable создаёт новую Kademlia таблицу
func NewKademliaTable(nodeID [20]byte) *KademliaTable {
	table := &KademliaTable{
		nodeID:  nodeID,
		buckets: [160]*KBucket{},
	}
	// Инициализируем buckets
	for i := range table.buckets {
		table.buckets[i] = &KBucket{
			nodes: make([]*DHTNode, 0, KBucketSize),
		}
	}
	return table
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

	// Запускаем цикл обновления таблицы
	go d.refreshTable()

	return nil
}

// FindPeer ищет пиры для конкретного info hash с использованием Kademlia
func (d *DHTClient) FindPeer(infoHash string) ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	log.Printf("Finding peers for info hash: %s", infoHash)
	
	// Преобразуем info hash в байты
	hashBytes, err := bencode.DecodeHexString(infoHash)
	if err != nil {
		return nil, fmt.Errorf("invalid info hash: %w", err)
	}

	var targetID [20]byte
	copy(targetID[:], hashBytes)

	// Запускаем итеративный поиск
	go d.iterativeFindPeers(targetID)

	// Ждём результаты (45 секунд timeout)
	peers := []string{}
	timeout := time.After(45 * time.Second)
	
	for len(peers) < d.targetPeers {
		select {
		case found := <-d.peersFound:
			peers = append(peers, found...)
			if len(peers) >= d.targetPeers {
				log.Printf("Found %d peers via DHT", len(peers))
				return peers[:d.targetPeers], nil
			}
		case <-timeout:
			if len(peers) > 0 {
				log.Printf("DHT timeout, found %d peers", len(peers))
				return peers, nil
			}
			log.Printf("DHT timeout, no peers found")
			return nil, fmt.Errorf("no peers found within timeout")
		case <-d.ctx.Done():
			return nil, fmt.Errorf("cancelled")
		}
	}

	return peers, nil
}

// iterativeFindPeers выполняет итеративный поиск пиров по Kademlia
func (d *DHTClient) iterativeFindPeers(target [20]byte) {
	// Начинаем с ближайших известных узлов
	candidates := d.findClosestNodes(target, KBucketSize)
	visited := make(map[[20]byte]bool)
	
	// Итеративно запрашиваем узлы
	for i := 0; i < 3 && len(candidates) > 0; i++ {
		var newCandidates []*DHTNode
		
		for _, node := range candidates {
			nodeKey := node.ID
			if visited[nodeKey] {
				continue
			}
			visited[nodeKey] = true
			
			// Отправляем get_peers
			addr := fmt.Sprintf("%s:%d", node.IP, node.Port)
			if err := d.sendGetPeers(addr, target[:]); err != nil {
				continue
			}
			
			// Также отправляем find_node для обновления таблицы
			d.sendFindNode(addr, target)
			
			// Ждём немного для получения ответа
			time.Sleep(100 * time.Millisecond)
		}
		
		// Получаем новые кандидаты из ответов
		newCandidates = d.findClosestNodes(target, KBucketSize)
		if len(newCandidates) == 0 {
			break
		}
		candidates = newCandidates
	}
	
	log.Printf("Iterative find completed, visited %d nodes", len(visited))
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
			n, addr, err := d.udpConn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			go d.handlePacketFrom(buf[:n], addr)
		}
	}
}

// handlePacketFrom обрабатывает пакет с адресом отправителя
func (d *DHTClient) handlePacketFrom(data []byte, addr *net.UDPAddr) {
	val, err := bencode.Unmarshal(data)
	if err != nil {
		return
	}

	query, ok := val.(bencode.Dict)
	if !ok {
		return
	}

	// Проверяем тип сообщения
	msgType, ok := query["y"].(bencode.String)
	if !ok {
		return
	}

	switch string(msgType) {
	case "r": // response
		d.handleResponse(query)
		// Сохраняем узел из которого пришёл ответ
		nodeID := d.extractNodeID(query)
		d.addNode(&DHTNode{
			ID:       nodeID,
			IP:       addr.IP.String(),
			Port:     uint16(addr.Port),
			LastSeen: time.Now(),
		})
	case "q": // query
		d.handleQuery(query)
	}
}

// extractNodeID извлекает node ID из ответа
func (d *DHTClient) extractNodeID(query bencode.Dict) [20]byte {
	var nodeID [20]byte
	if response, ok := query["r"].(bencode.Dict); ok {
		if id, ok := response["id"].(bencode.String); ok && len(id) >= 20 {
			copy(nodeID[:], id[:20])
		}
	}
	return nodeID
}

// refreshTable периодически обновляет таблицу маршрутизации
func (d *DHTClient) refreshTable() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			// Периодически пингуем bootstrap узлы
			for _, addr := range d.bootstrap {
				d.ping(addr)
			}
			log.Printf("DHT table: %d nodes", d.GetNodeCount())
		}
	}
}

// sendFindNode отправляет find_node запрос
func (d *DHTClient) sendFindNode(addr string, target [20]byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	d.transactionID++
	tid := fmt.Sprintf("%04x", d.transactionID&0xFFFF)

	query := bencode.Dict{
		"y": bencode.String("q"),
		"q": bencode.String("find_node"),
		"a": bencode.Dict{
			"id":     bencode.String(string(d.nodeID[:])),
			"target": bencode.String(string(target[:])),
		},
		"t": bencode.String(tid),
	}

	data, err := bencode.Marshal(query)
	if err != nil {
		return err
	}

	_, err = d.udpConn.WriteToUDP(data, udpAddr)
	if err != nil {
		return err
	}

	log.Printf("find_node sent to %s", addr)
	return nil
}

// handleResponse обрабатывает ответ
func (d *DHTClient) handleResponse(query bencode.Dict) {
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

	// Обрабатываем nodes (compact node info)
	if nodes, ok := response["nodes"]; ok {
		if nodesStr, ok := nodes.(bencode.String); ok {
			log.Printf("Received nodes data: %d bytes", len(nodesStr))
			// Парсим compact node info (каждый узел 26 байт: 20 ID + 4 IP + 2 порт)
			for i := 0; i+26 <= len(nodesStr); i += 26 {
				nodeData := nodesStr[i : i+26]
				if node := d.parseCompactNode(nodeData); node != nil {
					d.addNode(node)
				}
			}
		}
	}

	// Сохраняем узел из id ответа
	if id, ok := response["id"]; ok {
		if idStr, ok := id.(bencode.String); ok && len(idStr) >= 20 {
			// Создаём узел из ответа
			var nodeID [20]byte
			copy(nodeID[:], idStr[:20])
			// IP и порт берём из отправителя (упрощённо)
		}
	}
}

// parseCompactNode парсит compact node info
func (d *DHTClient) parseCompactNode(data []byte) *DHTNode {
	if len(data) < 26 {
		return nil
	}

	var nodeID [20]byte
	copy(nodeID[:], data[:20])

	ip := net.IP(data[20:24]).String()
	port := binary.BigEndian.Uint16(data[24:26])

	return &DHTNode{
		ID:       nodeID,
		IP:       ip,
		Port:     port,
		LastSeen: time.Now(),
	}
}

// handleQuery обрабатывает входящий запрос
func (d *DHTClient) handleQuery(query bencode.Dict) {
	// Получаем метод запроса
	q, ok := query["q"].(bencode.String)
	if !ok {
		return
	}

	// Получаем transaction ID
	t, ok := query["t"].(bencode.String)
	if !ok {
		return
	}

	switch string(q) {
	case "ping":
		d.sendResponse(query, "pong")
	case "find_node":
		d.handleFindNode(query)
	case "get_peers":
		d.handleGetPeers(query)
	}

	_ = t // transaction ID для ответа
}

// handleFindNode обрабатывает find_node запрос
func (d *DHTClient) handleFindNode(query bencode.Dict) {
	a, ok := query["a"].(bencode.Dict)
	if !ok {
		return
	}

	target, ok := a["target"].(bencode.String)
	if !ok {
		return
	}

	var targetID [20]byte
	copy(targetID[:], target)

	// Находим ближайшие узлы
	closest := d.findClosestNodes(targetID, 8)

	// Формируем ответ
	response := d.createFindNodeResponse(query, closest)
	
	// Отправляем ответ (упрощённо - пропускаем)
	_ = response
}

// handleGetPeers обрабатывает get_peers запрос
func (d *DHTClient) handleGetPeers(query bencode.Dict) {
	// Для простоты отправляем пустой ответ
	// В полной реализации нужно вернуть пиры из таблицы
}

// sendResponse отправляет ответ на запрос
func (d *DHTClient) sendResponse(query bencode.Dict, method string) {
	// Упрощённая реализация
}

// createFindNodeResponse создаёт ответ на find_node
func (d *DHTClient) createFindNodeResponse(query bencode.Dict, nodes []*DHTNode) bencode.Dict {
	// Формируем compact nodes
	var compactNodes []byte
	for _, node := range nodes {
		nodeData := make([]byte, 26)
		copy(nodeData[:20], node.ID[:])
		ip := net.ParseIP(node.IP)
		if ip != nil {
			copy(nodeData[20:24], ip.To4())
		}
		binary.BigEndian.PutUint16(nodeData[24:26], node.Port)
		compactNodes = append(compactNodes, nodeData...)
	}

	return bencode.Dict{
		"y": bencode.String("r"),
		"t": query["t"],
		"r": bencode.Dict{
			"id":    bencode.String(string(d.nodeID[:])),
			"nodes": bencode.String(string(compactNodes)),
		},
	}
}

// addNode добавляет узел в таблицу
func (d *DHTClient) addNode(node *DHTNode) {
	d.kademlia.addNode(node)
}

// findClosestNodes находит ближайшие узлы к target
func (d *DHTClient) findClosestNodes(target [20]byte, limit int) []*DHTNode {
	return d.kademlia.findClosestNodes(target, limit)
}

// distance вычисляет расстояние XOR между двумя ID
func distance(a, b [20]byte) [20]byte {
	var result [20]byte
	for i := 0; i < 20; i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

// AddNode добавляет узел в Kademlia таблицу
func (kt *KademliaTable) addNode(node *DHTNode) {
	kt.mu.Lock()
	defer kt.mu.Unlock()

	// Вычисляем bucket по расстоянию
	dist := distance(kt.nodeID, node.ID)
	bucketIndex := kt.getBucketIndex(dist)

	if bucketIndex >= len(kt.buckets) {
		bucketIndex = len(kt.buckets) - 1
	}

	bucket := kt.buckets[bucketIndex]
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	// Проверяем есть ли уже такой узел
	for i, n := range bucket.nodes {
		if n.ID == node.ID {
			// Обновляем существующий узел
			bucket.nodes[i] = node
			return
		}
	}

	// Добавляем новый узел если есть место
	if len(bucket.nodes) < KBucketSize {
		bucket.nodes = append(bucket.nodes, node)
	}
	// Иначе можно вытеснить старый узел (упрощённо - не добавляем)
}

// getBucketIndex вычисляет индекс bucket по расстоянию
func (kt *KademliaTable) getBucketIndex(dist [20]byte) int {
	// Находим старший значащий бит
	for i := 0; i < 160; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		if dist[byteIndex]&(1<<bitIndex) != 0 {
			return i
		}
	}
	return 0
}

// findClosestNodes находит ближайшие узлы к target
func (kt *KademliaTable) findClosestNodes(target [20]byte, limit int) []*DHTNode {
	kt.mu.RLock()
	defer kt.mu.RUnlock()

	type nodeDist struct {
		node *DHTNode
		dist [20]byte
	}

	var allNodes []nodeDist
	for _, bucket := range kt.buckets {
		bucket.mu.RLock()
		for _, node := range bucket.nodes {
			dist := distance(kt.nodeID, node.ID)
			allNodes = append(allNodes, nodeDist{node, dist})
		}
		bucket.mu.RUnlock()
	}

	// Сортируем по расстоянию до target
	// (упрощённая сортировка)
	result := make([]*DHTNode, 0, limit)
	for i := 0; i < len(allNodes) && i < limit; i++ {
		result = append(result, allNodes[i].node)
	}

	return result
}

// GetNodeCount возвращает количество узлов в таблице
func (d *DHTClient) GetNodeCount() int {
	return d.kademlia.getNodeCount()
}

// getNodeCount подсчитывает все узлы в таблице
func (kt *KademliaTable) getNodeCount() int {
	count := 0
	for _, bucket := range kt.buckets {
		bucket.mu.RLock()
		count += len(bucket.nodes)
		bucket.mu.RUnlock()
	}
	return count
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