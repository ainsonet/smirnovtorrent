package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	Progress      float64 `json:"progress"`
	Downloaded    int64   `json:"downloaded"`
	TotalSize     int64   `json:"totalSize"`
	ActivePeers   int     `json:"activePeers"`
	DownloadSpeed float64 `json:"downloadSpeed"`
	UploadSpeed   float64 `json:"uploadSpeed"`
	Status        string  `json:"status"`
	TorrentName   string  `json:"torrentName"`
	Path          string  `json:"path"`
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
	// Открываем браузер после запуска сервера
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(fmt.Sprintf("http://localhost:%d", w.port))
	}()

	// Регистрируем обработчики
	http.HandleFunc("/", w.handleIndex)
	http.HandleFunc("/api/status", w.handleAPIStatus)
	http.HandleFunc("/api/add", w.handleAPIAdd)
	http.HandleFunc("/api/start", w.handleAPIStart)
	http.HandleFunc("/api/stop", w.handleAPIStop)
	http.HandleFunc("/api/remove", w.handleAPIRemove)
	http.HandleFunc("/api/pause", w.handleAPIPause)
	http.HandleFunc("/api/resume", w.handleAPIResume)
	http.HandleFunc("/api/select-file", w.handleAPISelectFile)
	
	addr := fmt.Sprintf(":%d", w.port)
	log.Printf("Web UI starting on http://localhost%s", addr)
	log.Printf("Opening browser automatically...")
	
	return http.ListenAndServe(addr, nil)
}

// openBrowser открывает браузер
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
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
		Path string  `json:"path"`
		Size int64   `json:"size"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(rw, "Path is required", http.StatusBadRequest)
		return
	}

	// Реалистичные демо-данные
	totalSize := req.Size
	if totalSize == 0 {
		totalSize = 1024 * 1024 * 100 // 100MB по умолчанию
	}

	w.mu.Lock()
	w.status.TorrentName = filepath.Base(req.Path)
	w.status.Path = req.Path
	w.status.Status = "downloading"
	w.status.Progress = 0
	w.status.Downloaded = 0
	w.status.TotalSize = totalSize
	w.status.DownloadSpeed = 0
	w.status.UploadSpeed = 0
	w.status.ActivePeers = 0
	w.mu.Unlock()

	addLog("Added torrent: " + filepath.Base(req.Path))

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Torrent added successfully",
	})
}

// handleAPISelectFile обрабатывает выбор файла
func (w *WebUI) handleAPISelectFile(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Для демо возвращаем путь к файлу
	// В реальной версии здесь будет диалог выбора файла
	path := "demo.torrent"
	size := int64(524288000) // 500MB

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"path": path,
		"size": size,
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

	// Реалистичная симуляция загрузки
	if w.status.Status == "downloading" && w.status.TotalSize > 0 {
		// Увеличиваем прогресс реалистично (не больше 99%)
		increment := 0.3 // 0.3% каждые 2 секунды
		if w.status.Progress < 99.0 {
			w.status.Progress += increment
		}
		
		// Расчитываем downloaded на основе прогресса
		w.status.Downloaded = int64(float64(w.status.TotalSize) * w.status.Progress / 100.0)
		
		// Скорости
		w.status.DownloadSpeed = 2.5 * 1024 * 1024 // 2.5 MB/s
		w.status.UploadSpeed = 0.5 * 1024 * 1024   // 0.5 MB/s
		w.status.ActivePeers = 15 + rand.Intn(10)
	}

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
