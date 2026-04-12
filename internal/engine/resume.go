package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ResumeData данные для продолжения загрузки
type ResumeData struct {
	InfoHash      string     `json:"info_hash"`
	TorrentName   string     `json:"name,omitempty"`
	TotalSize     int64      `json:"total_size,omitempty"`
	PieceLength   int32      `json:"piece_length,omitempty"`
	CompletedPieces []int    `json:"completed_pieces"`
	Downloaded    int64      `json:"downloaded"`
	Uploaded      int64      `json:"uploaded"`
	StartTime     int64      `json:"start_time"`
	LastSaveTime  int64      `json:"last_save_time"`
	Peers         []PeerResumeInfo `json:"peers,omitempty"`
}

// PeerResumeInfo информация о пире для сохранения
type PeerResumeInfo struct {
	IP        string `json:"ip"`
	Port      uint16 `json:"port"`
	Encrypted bool   `json:"encrypted"`
}

// ResumeManager управляет сохранением/загрузкой прогресса
type ResumeManager struct {
	data        ResumeData
	filePath    string
	mu          sync.RWMutex
	autoSave    bool
	saveTicker  *time.Ticker
}

// NewResumeManager создаёт менеджер продолжения
func NewResumeManager(infoHash string, outputDir string) *ResumeManager {
	// Создаём путь к файлу сохранения
	filename := fmt.Sprintf("%s.resume", infoHash)
	dir := filepath.Join(outputDir, ".smirnovtorrent")
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		// Если не удалось создать директорию, используем текущую
		dir = "."
	}

	filePath := filepath.Join(dir, filename)

	return &ResumeManager{
		data: ResumeData{
			InfoHash: infoHash,
		},
		filePath: filePath,
		autoSave: true,
	}
}

// Load загружает данные из файла
func (rm *ResumeManager) Load() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	data, err := os.ReadFile(rm.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Файл не существует - это нормально для новой загрузки
			return nil
		}
		return fmt.Errorf("failed to read resume file: %w", err)
	}

	if err := json.Unmarshal(data, &rm.data); err != nil {
		return fmt.Errorf("failed to parse resume data: %w", err)
	}

	return nil
}

// Save сохраняет данные в файл
func (rm *ResumeManager) Save() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	data, err := json.MarshalIndent(rm.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal resume data: %w", err)
	}

	if err := os.WriteFile(rm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write resume file: %w", err)
	}

	return nil
}

// StartAutoSave запускает автоматическое сохранение
func (rm *ResumeManager) StartAutoSave(interval time.Duration) {
	rm.saveTicker = time.NewTicker(interval)
	
	go func() {
		for range rm.saveTicker.C {
			rm.Save()
		}
	}()
}

// StopAutoSave останавливает автоматическое сохранение
func (rm *ResumeManager) StopAutoSave() {
	if rm.saveTicker != nil {
		rm.saveTicker.Stop()
	}
}

// MarkPieceComplete помечает кусок как завершённый
func (rm *ResumeManager) MarkPieceComplete(pieceIndex int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Проверяем что кусок ещё не завершён
	for _, idx := range rm.data.CompletedPieces {
		if idx == pieceIndex {
			return // Уже завершён
		}
	}

	rm.data.CompletedPieces = append(rm.data.CompletedPieces, pieceIndex)
	rm.data.Downloaded += int64(pieceIndex) // Упрощённо
	
	if rm.autoSave {
		go rm.Save()
	}
}

// IsPieceComplete проверяет завершён ли кусок
func (rm *ResumeManager) IsPieceComplete(pieceIndex int) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	for _, idx := range rm.data.CompletedPieces {
		if idx == pieceIndex {
			return true
		}
	}
	return false
}

// GetCompletedPieces возвращает список завершённых кусков
func (rm *ResumeManager) GetCompletedPieces() []int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	result := make([]int, len(rm.data.CompletedPieces))
	copy(result, rm.data.CompletedPieces)
	return result
}

// SetDownloaded устанавливает размер загруженных данных
func (rm *ResumeManager) SetDownloaded(bytes int64) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.data.Downloaded = bytes
}

// SetUploaded устанавливает размер отправленных данных
func (rm *ResumeManager) SetUploaded(bytes int64) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.data.Uploaded = bytes
}

// GetDownloaded возвращает размер загруженных данных
func (rm *ResumeManager) GetDownloaded() int64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.data.Downloaded
}

// GetUploaded возвращает размер отправленных данных
func (rm *ResumeManager) GetUploaded() int64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.data.Uploaded
}

// Delete удаляет файл сохранения
func (rm *ResumeManager) Delete() error {
	rm.StopAutoSave()
	
	if err := os.Remove(rm.filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete resume file: %w", err)
	}
	
	return nil
}

// SetAutoSave включает/выключает авто-сохранение
func (rm *ResumeManager) SetAutoSave(enabled bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.autoSave = enabled
}

// GetProgress возвращает прогресс загрузки
func (rm *ResumeManager) GetProgress(totalPieces int) float64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if totalPieces == 0 {
		return 0
	}

	return float64(len(rm.data.CompletedPieces)) / float64(totalPieces) * 100
}

// AddPeer добавляет пирa в список
func (rm *ResumeManager) AddPeer(ip string, port uint16, encrypted bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	// Проверяем есть ли уже такой пир
	for _, p := range rm.data.Peers {
		if p.IP == ip && p.Port == port {
			return
		}
	}
	
	rm.data.Peers = append(rm.data.Peers, PeerResumeInfo{
		IP:        ip,
		Port:      port,
		Encrypted: encrypted,
	})
}

// GetPeers возвращает список сохранённых пиров
func (rm *ResumeManager) GetPeers() []PeerResumeInfo {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	peers := make([]PeerResumeInfo, len(rm.data.Peers))
	copy(peers, rm.data.Peers)
	return peers
}

// SetTorrentInfo устанавливает информацию о торренте
func (rm *ResumeManager) SetTorrentInfo(name string, totalSize int64, pieceLength int32) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	rm.data.TorrentName = name
	rm.data.TotalSize = totalSize
	rm.data.PieceLength = pieceLength
}

// Stop останавливает авто-сохранение и сохраняет финальное состояние
func (rm *ResumeManager) Stop() error {
	rm.StopAutoSave()
	return rm.Save()
}

// UpdateCompletePieces обновляет список завершённых кусков
func (rm *ResumeManager) UpdateCompletePieces(pieces []int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	rm.data.CompletedPieces = make([]int, len(pieces))
	copy(rm.data.CompletedPieces, pieces)
	rm.data.LastSaveTime = time.Now().Unix()
}

// UpdateDownloaded обновляет размер загруженных данных
func (rm *ResumeManager) UpdateDownloaded(downloaded int64) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.data.Downloaded = downloaded
	rm.data.LastSaveTime = time.Now().Unix()
}

// Cleanup удаляет все файлы сохранения в директории
func CleanupResumeFiles(outputDir string) error {
	dir := filepath.Join(outputDir, ".smirnovtorrent")
	
	files, err := filepath.Glob(filepath.Join(dir, "*.resume"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}