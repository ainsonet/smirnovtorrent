package tracker

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	Interval   int         `json:"interval"`
	Peers      []PeerInfo  `json:"peers"`
	Failure    string      `json:"failure reason"`
	Warnings   []string    `json:"warnings"`
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
	query.Set("info_hash", params.InfoHash)
	query.Set("peer_id", params.PeerID)
	query.Set("port", strconv.Itoa(int(params.Port)))
	query.Set("downloaded", strconv.FormatInt(params.Downloaded, 10))
	query.Set("left", strconv.FormatInt(params.Left, 10))
	query.Set("uploaded", strconv.FormatInt(params.Uploaded, 10))

	if params.Event != "" {
		query.Set("event", params.Event)
	}

	baseURL.RawQuery = query.Encode()

	resp, err := http.Get(baseURL.String())
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

	var announceResp AnnounceResponse
	if err := json.Unmarshal(body, &announceResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if announceResp.Failure != "" {
		return nil, fmt.Errorf("tracker error: %s", announceResp.Failure)
	}

	return &announceResp, nil
}

// GetPeers это удобная обёртка для Announce с дефолтными параметрами
func (t *Tracker) GetPeers(infoHash string, peerID string, port uint16) ([]PeerInfo, error) {
	params := AnnounceParams{
		InfoHash:   infoHash,
		PeerID:     peerID,
		Port:       port,
		Downloaded: 0,
		Left:       0,
		Uploaded:   0,
		Event:      "started",
	}

	resp, err := t.Announce(params)
	if err != nil {
		return nil, err
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