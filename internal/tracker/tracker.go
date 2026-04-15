package tracker

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"smirnovtorrent/pkg/bencode"
)

// PeerInfo информация о пире
type PeerInfo struct {
	IP       string
	Port     uint16
	PeerID   string
	Downloaded int64
	Left     int64
	Uploaded int64
}

// AnnounceParams параметры для announce запроса
type AnnounceParams struct {
	InfoHash   string
	PeerID     string
	Port       uint16
	Downloaded int64
	Left       int64
	Uploaded   int64
	Event      string // started, completed, stopped
}

// AnnounceResponse ответ от трекера
type AnnounceResponse struct {
	Interval   int         `bencode:"interval,omitempty"`
	Peers      []PeerInfo  `bencode:"peers,omitempty"`
	PeersCompact string    `bencode:"peers,omitempty"` // компактный формат
	Failure    string      `bencode:"failure reason,omitempty"`
	Warnings   []string    `bencode:"warning message,omitempty"`
}

// Tracker клиент для работы с трекером
type Tracker struct {
	announceURL string
}

// NewTracker создаёт новый трекер
func NewTracker(announceURL string) *Tracker {
	return &Tracker{
		announceURL: announceURL,
	}
}

// Announce отправляет запрос к трекеру и получает список пиров
func (t *Tracker) Announce(params AnnounceParams) (*AnnounceResponse, error) {
	baseURL, err := url.Parse(t.announceURL)
	if err != nil {
		return nil, fmt.Errorf("invalid tracker URL: %w", err)
	}

	query := baseURL.Query()
	
	// Info hash нужно закодировать как raw bytes в URL
	// Проверяем формат - может быть hex строка (40 символов) или raw bytes (20 байт)
	var infoHashBytes []byte
	if len(params.InfoHash) == 40 {
		// Hex формат (например: "378eb779eb59bb66b666f25fc1ecc70fb151aa60")
		var err error
		infoHashBytes, err = hex.DecodeString(params.InfoHash)
		if err != nil {
			return nil, fmt.Errorf("invalid info hash format: %w", err)
		}
	} else if len(params.InfoHash) == 20 {
		// Уже сырые байты
		infoHashBytes = []byte(params.InfoHash)
	} else {
		return nil, fmt.Errorf("invalid info hash length: expected 20 bytes or 40 hex chars, got %d", len(params.InfoHash))
	}
	
	// URL encode raw bytes
	query.Set("info_hash", string(infoHashBytes))
	query.Set("peer_id", params.PeerID)
	query.Set("port", strconv.Itoa(int(params.Port)))
	query.Set("downloaded", strconv.FormatInt(params.Downloaded, 10))
	query.Set("left", strconv.FormatInt(params.Left, 10))
	query.Set("uploaded", strconv.FormatInt(params.Uploaded, 10))

	if params.Event != "" {
		query.Set("event", params.Event)
	}

	baseURL.RawQuery = query.Encode()

	// Увеличиваем таймаут для трекеров
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("tracker request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tracker returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Пытаемся распарсить как bencode (стандарт BitTorrent)
	var announceResp AnnounceResponse
	if val, err := bencode.Unmarshal(body); err == nil {
		// Успешно распарсили как bencode
		if dict, ok := val.(bencode.Dict); ok {
			// Извлекаем поля вручную
			if interval, ok := dict["interval"].(bencode.Int); ok {
				announceResp.Interval = int(interval)
			}
			if failure, ok := dict["failure reason"].(bencode.String); ok {
				announceResp.Failure = string(failure)
			}
			if warning, ok := dict["warning message"].(bencode.String); ok {
				announceResp.Warnings = []string{string(warning)}
			}
			
			// Пытаемся получить пиры
			if peers, ok := dict["peers"]; ok {
				if peersStr, ok := peers.(bencode.String); ok {
					announceResp.PeersCompact = string(peersStr)
				} else if peersList, ok := peers.(bencode.List); ok {
					// Парсим список пиров
					for _, peerVal := range peersList {
						if peerDict, ok := peerVal.(bencode.Dict); ok {
							peer := PeerInfo{}
							if ip, ok := peerDict["ip"].(bencode.String); ok {
								peer.IP = string(ip)
							}
							if port, ok := peerDict["port"].(bencode.Int); ok {
								peer.Port = uint16(port)
							}
							if peer.IP != "" {
								announceResp.Peers = append(announceResp.Peers, peer)
							}
						}
					}
				}
			}
		}
	} else {
		// Если bencode не работает, пробуем как JSON
		if jsonErr := json.Unmarshal(body, &announceResp); jsonErr != nil {
			return nil, fmt.Errorf("failed to parse response (bencode: %v, json: %v)", err, jsonErr)
		}
	}

	if announceResp.Failure != "" {
		return nil, fmt.Errorf("tracker error: %s", announceResp.Failure)
	}

	// Если пиры в компактном формате, парсим их
	if announceResp.PeersCompact != "" {
		peers := parseCompactPeers(announceResp.PeersCompact)
		announceResp.Peers = append(announceResp.Peers, peers...)
	}

	return &announceResp, nil
}

// GetPeers это удобная обёртка для Announce с дефолтными параметрами
func (t *Tracker) GetPeers(infoHash string, peerID string, port uint16) ([]PeerInfo, error) {
	// Info hash должен быть в URL-safe формате для трекера
	// Трекеры ожидают raw bytes, закодированные в URL
	params := AnnounceParams{
		InfoHash:   infoHash, // уже в hex формате
		PeerID:     peerID,
		Port:       port,
		Downloaded: 0,
		Left:       0,
		Uploaded:   0,
		Event:      "started",
	}

	resp, err := t.Announce(params)
	if err != nil {
		return nil, fmt.Errorf("tracker error: %w", err)
	}
	
	return resp.Peers, nil
}

// ParsePeerURL парсит URL трекера из announce
func ParsePeerURL(announce string) (*Tracker, error) {
	if announce == "" {
		return nil, fmt.Errorf("announce URL is empty")
	}

	return NewTracker(announce), nil
}

// Helper functions

// parseCompactPeers парсит компактный формат пиров (6 байт на пира: 4 IP + 2 порт)
func parseCompactPeers(compact string) []PeerInfo {
	peers := []PeerInfo{}
	
	// Пытаемся как hex строку
	data, err := hex.DecodeString(compact)
	if err != nil {
		// Если не hex, пробуем как сырые байты
		data = []byte(compact)
	}
	
	// Каждый пир 6 байт
	for i := 0; i+6 <= len(data); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", 
			data[i], data[i+1], data[i+2], data[i+3])
		port := uint16(data[i+4])<<8 | uint16(data[i+5])
		
		peers = append(peers, PeerInfo{
			IP:   ip,
			Port: port,
		})
	}
	
	return peers
}

// EncodeInfoHash конвертирует info hash в формат для трекера (URL-safe base64)
func EncodeInfoHash(hash []byte) string {
	return hex.EncodeToString(hash)
}

// ParsePeerInfo извлекает информацию о пире из ответа трекера
func ParsePeerInfo(peerData map[string]interface{}) (*PeerInfo, error) {
	peer := &PeerInfo{}

	if ip, ok := peerData["ip"].(string); ok {
		peer.IP = ip
	}

	if port, ok := peerData["port"].(float64); ok {
		peer.Port = uint16(port)
	}

	if peerID, ok := peerData["peer id"].(string); ok {
		peer.PeerID = peerID
	}

	if downloaded, ok := peerData["downloaded"].(float64); ok {
		peer.Downloaded = int64(downloaded)
	}

	if left, ok := peerData["left"].(float64); ok {
		peer.Left = int64(left)
	}

	if uploaded, ok := peerData["uploaded"].(float64); ok {
		peer.Uploaded = int64(uploaded)
	}

	return peer, nil
}