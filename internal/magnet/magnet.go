package magnet

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// MagnetLink структура magnet ссылки
type MagnetLink struct {
	InfoHash    string
	DisplayName string
	Trackers    []string
	PEX         bool // Peer Exchange
	DHT         bool // Distributed Hash Table
}

// Parse разбирает magnet ссылку
func Parse(magnet string) (*MagnetLink, error) {
	if !strings.HasPrefix(magnet, "magnet:?") {
		return nil, fmt.Errorf("invalid magnet link: must start with 'magnet:?'")
	}

	u, err := url.Parse(magnet)
	if err != nil {
		return nil, fmt.Errorf("failed to parse magnet link: %w", err)
	}

	link := &MagnetLink{}

	// Парсим параметры
	for k, v := range u.Query() {
		switch k {
		case "xt":
			// Exact Topic - info hash
			for _, xt := range v {
				if strings.HasPrefix(xt, "urn:btih:") {
					hash := strings.TrimPrefix(xt, "urn:btih:")
					link.InfoHash = strings.ToLower(hash)
				}
			}

		case "tr":
			// Tracker URLs
			link.Trackers = append(link.Trackers, v...)

		case "dn":
			// Display name
			if len(v) > 0 {
				link.DisplayName = v[0]
			}

		case "xt.1", "xt.2":
			// Дополнительные xt параметры (для multi-hash)
			// Пока игнорируем

		case "pk":
			// Public key (для DHT)
			// Пока игнорируем

		case "x.pe":
			// Peer exchange
			link.PEX = true

		case "dht":
			// Distributed Hash Table
			if len(v) > 0 && v[0] == "on" {
				link.DHT = true
			}
		}
	}

	// Проверяем обязательные параметры
	if link.InfoHash == "" {
		return nil, fmt.Errorf("missing info hash (xt parameter)")
	}

	// Валидируем info hash
	if err := link.validateInfoHash(); err != nil {
		return nil, err
	}

	return link, nil
}

// validateInfoHash проверяет что info hash валиден
func (m *MagnetLink) validateInfoHash() error {
	// Info hash может быть:
	// 1. 40 символов hex (SHA-1)
	// 2. 32 символа base32 (для совместимости)
	// 3. URL-safe base64 (реже)

	if len(m.InfoHash) == 40 {
		// Проверяем что это hex
		hexPattern := regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
		if !hexPattern.MatchString(m.InfoHash) {
			return fmt.Errorf("invalid info hash format: expected 40 hex characters")
		}
		m.InfoHash = strings.ToLower(m.InfoHash)
		return nil
	}

	if len(m.InfoHash) == 32 {
		// Base32 hash - конвертируем в hex
		// Пока оставляем как есть, конвертация будет в parser
		return nil
	}

	return fmt.Errorf("invalid info hash length: %d (expected 40 hex or 32 base32)", len(m.InfoHash))
}

// String возвращает magnet ссылку как строку
func (m *MagnetLink) String() string {
	var sb strings.Builder
	sb.WriteString("magnet:?")
	sb.WriteString("xt=urn:btih:")
	sb.WriteString(m.InfoHash)

	if m.DisplayName != "" {
		sb.WriteString("&dn=")
		sb.WriteString(url.QueryEscape(m.DisplayName))
	}

	for _, tracker := range m.Trackers {
		sb.WriteString("&tr=")
		sb.WriteString(url.QueryEscape(tracker))
	}

	if m.PEX {
		sb.WriteString("&x.pe=1")
	}

	if m.DHT {
		sb.WriteString("&dht=on")
	}

	return sb.String()
}

// ToTrackerURL конвертирует в URL трекера
func (m *MagnetLink) ToTrackerURL(port uint16) string {
	var sb strings.Builder
	if len(m.Trackers) > 0 {
		sb.WriteString(m.Trackers[0])
	} else {
		// Дефолтный трекер
		sb.WriteString("http://tracker.example.com/announce")
	}

	sb.WriteString("?")
	sb.WriteString("info_hash=")
	sb.WriteString(m.InfoHash)
	sb.WriteString("&")
	sb.WriteString("peer_id=-")
	sb.WriteString("-SMRV0100-")
	sb.WriteString(fmt.Sprintf("&port=%d", port))
	sb.WriteString("&")
	sb.WriteString("uploaded=0&")
	sb.WriteString("downloaded=0&")
	sb.WriteString("left=0")

	return sb.String()
}

// Helper functions

// ExtractInfoHash извлекает info hash из magnet ссылки
func ExtractInfoHash(magnet string) (string, error) {
	link, err := Parse(magnet)
	if err != nil {
		return "", err
	}
	return link.InfoHash, nil
}

// ExtractDisplayName извлекает имя из magnet ссылки
func ExtractDisplayName(magnet string) (string, error) {
	link, err := Parse(magnet)
	if err != nil {
		return "", err
	}
	return link.DisplayName, nil
}

// IsMagnetLink проверяет что строка это magnet ссылка
func IsMagnetLink(s string) bool {
	return strings.HasPrefix(s, "magnet:?")
}

// BuildMagnetLink создаёт magnet ссылку из параметров
func BuildMagnetLink(infoHash string, displayName string, trackers []string) string {
	link := &MagnetLink{
		InfoHash:    infoHash,
		DisplayName: displayName,
		Trackers:    trackers,
		DHT:         true,
	}
	return link.String()
}