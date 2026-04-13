package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"smirnovtorrent/internal/engine"
)

// WebUI представляет веб-интерфейс
type WebUI struct {
	engine   *engine.DownloadEngine
	port     int
	mu       sync.RWMutex
	status   DownloadStatus
}

// DownloadStatus статус загрузки для Web UI
type DownloadStatus struct {
	Progress     float64 `json:"progress"`
	Downloaded   int64   `json:"downloaded"`
	TotalSize    int64   `json:"totalSize"`
	ActivePeers  int     `json:"activePeers"`
	DownloadSpeed float64 `json:"downloadSpeed"`
	UploadSpeed  float64 `json:"uploadSpeed"`
	Status       string  `json:"status"`
	TorrentName  string  `json:"torrentName"`
}

// NewWebUI создаёт новый веб-интерфейс
func NewWebUI(eng *engine.DownloadEngine, port int) *WebUI {
	return &WebUI{
		engine: eng,
		port:   port,
		status: DownloadStatus{
			Status: "initializing",
		},
	}
}

// Start запускает веб-сервер
func (w *WebUI) Start() error {
	// Регистрируем обработчики
	http.HandleFunc("/", w.handleIndex)
	http.HandleFunc("/api/status", w.handleAPIStatus)
	http.HandleFunc("/api/add", w.handleAPIAdd)
	http.HandleFunc("/api/start", w.handleAPIStart)
	http.HandleFunc("/api/stop", w.handleAPIStop)
	http.HandleFunc("/api/remove", w.handleAPIRemove)
	http.HandleFunc("/api/pause", w.handleAPIPause)
	http.HandleFunc("/api/resume", w.handleAPIResume)
	
	addr := fmt.Sprintf(":%d", w.port)
	log.Printf("Web UI starting on http://localhost%s", addr)
	
	return http.ListenAndServe(addr, nil)
}

// handleIndex обрабатывает главную страницу
func (w *WebUI) handleIndex(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	}

	// Читаем встроенный HTML файл
	htmlContent, err := webUIAssets.ReadFile("webui.html")
	if err != nil {
		// Если файл не найден, показываем простой HTML
		rw.Header().Set("Content-Type", "text/html")
		fmt.Fprint(rw, "<html><body><h1>SmirnovTorrent Web UI</h1><p>Error loading interface</p></body></html>")
		return
	}

	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprint(rw, string(htmlContent))
}

// handleAPIStatus обрабатывает API запрос статуса
func (w *WebUI) handleAPIStatus(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.RLock()
	status := w.status
	w.mu.RUnlock()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(status)
}

// handleAPIStart обрабатывает API запрос запуска
func (w *WebUI) handleAPIStart(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.Lock()
	w.status.Status = "running"
	w.mu.Unlock()

	addLog("Download started via Web UI")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Download started",
	})
}

// handleAPIStop обрабатывает API запрос остановки
func (w *WebUI) handleAPIStop(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.Lock()
	w.status.Status = "stopped"
	w.mu.Unlock()

	addLog("Download stopped via Web UI")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Download stopped",
	})
}

// handleAPIAdd обрабатывает добавление торрента
func (w *WebUI) handleAPIAdd(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(rw, "Path is required", http.StatusBadRequest)
		return
	}

	w.mu.Lock()
	w.status.TorrentName = req.Path
	w.status.Status = "downloading"
	w.status.Progress = 0
	w.status.Downloaded = 0
	w.status.TotalSize = 100 * 1024 * 1024 // 100MB demo
	w.mu.Unlock()

	addLog("Added torrent: " + req.Path)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Torrent added successfully",
		"id":      "demo-" + fmt.Sprint(time.Now().Unix()),
	})
}

// handleAPIRemove обрабатывает удаление торрента
func (w *WebUI) handleAPIRemove(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.Lock()
	w.status.Status = "stopped"
	w.status.TorrentName = ""
	w.mu.Unlock()

	addLog("Torrent removed")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Torrent removed",
	})
}

// handleAPIPause обрабатывает паузу
func (w *WebUI) handleAPIPause(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.Lock()
	w.status.Status = "paused"
	w.mu.Unlock()

	addLog("Download paused")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Download paused",
	})
}

// handleAPIResume обрабатывает продолжение
func (w *WebUI) handleAPIResume(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.Lock()
	w.status.Status = "downloading"
	w.mu.Unlock()

	addLog("Download resumed")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Download resumed",
	})
}

// UpdateStatus обновляет статус загрузки
func (w *WebUI) UpdateStatus(status engine.DownloadStatus, torrentName string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.status.Progress = status.Progress
	w.status.Downloaded = status.Downloaded
	w.status.TotalSize = status.TotalSize
	w.status.ActivePeers = status.ActivePeers
	w.status.DownloadSpeed = status.DownloadSpeed
	w.status.TorrentName = torrentName
}

var logFile *os.File
var logMu sync.Mutex

func addLog(message string) {
	logMu.Lock()
	defer logMu.Unlock()

	log.Println(message)

	if logFile != nil {
		fmt.Fprintf(logFile, "%s %s\n", os.Getenv("TIME"), message)
	}
}

// InitLog инициализирует лог файл
func InitLog() error {
	logDir := filepath.Join(os.TempDir(), "smirnovtorrent")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logDir, "webui.log")
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	return nil
}
