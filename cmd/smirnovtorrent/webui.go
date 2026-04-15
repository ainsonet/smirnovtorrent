package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"smirnovtorrent/internal/engine"
	"smirnovtorrent/internal/parser"
)

// ActiveTorrent активный торрент
type ActiveTorrent struct {
	Torrent     *parser.Torrent
	Engine      *engine.AnacrolixEngine
	Status      DownloadStatus
	DownloadDir string
	lastUpdate  time.Time
	lastBytes   int64
}

// WebUI представляет веб-интерфейс
type WebUI struct {
	mu     sync.RWMutex
	port   int
	torrents map[string]*ActiveTorrent // hash -> torrent
}

// DownloadStatus статус загрузки для Web UI
type DownloadStatus struct {
	Progress      float64 `json:"progress"`
	Downloaded    int64   `json:"downloaded"`
	Uploaded      int64   `json:"uploaded"`
	TotalSize     int64   `json:"totalSize"`
	ActivePeers   int     `json:"activePeers"`
	DownloadSpeed float64 `json:"downloadSpeed"`
	UploadSpeed   float64 `json:"uploadSpeed"`
	Status        string  `json:"status"`
	TorrentName   string  `json:"torrentName"`
	Path          string  `json:"path"`
}

// NewWebUI создаёт новый веб-интерфейс
func NewWebUI(port int) *WebUI {
	return &WebUI{
		port:     port,
		torrents: make(map[string]*ActiveTorrent),
	}
}

// Start запускает веб-сервер
func (w *WebUI) Start() error {
	// Открываем браузер после запуска сервера
	go func() {
		time.Sleep(1 * time.Second)
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
	http.HandleFunc("/api/open-folder", w.handleAPIOpenFolder)
	http.HandleFunc("/logo.png", w.handleLogo)
	
	addr := fmt.Sprintf(":%d", w.port)
	log.Printf("Web UI starting on http://localhost%s", addr)
	log.Printf("Opening browser automatically in 1 second...")
	
	return http.ListenAndServe(addr, nil)
}

// openBrowser открывает браузер
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("cmd", "/c", "start", "", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

// handleLogo отдаёт логотип
func (w *WebUI) handleLogo(rw http.ResponseWriter, r *http.Request) {
	logoPath := "C:\\Users\\user\\Documents\\Visual Studio Code\\SmirnovTorrent\\logo.png"
	
	data, err := os.ReadFile(logoPath)
	if err != nil {
		// Если файл не найден, отдаём пустой ответ
		http.NotFound(rw, r)
		return
	}

	rw.Header().Set("Content-Type", "image/png")
	rw.Write(data)
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
	defer w.mu.RUnlock()

	// Возвращаем статус первого торрента или "no torrent"
	if len(w.torrents) == 0 {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"status":         "stopped",
			"progress":       0,
			"downloaded":     0,
			"uploaded":       0,
			"totalSize":      0,
			"activePeers":    0,
			"downloadSpeed":  0,
			"uploadSpeed":    0,
			"torrentName":    "",
		})
		return
	}

	// Берём первый торрент
	for _, t := range w.torrents {
		// Обновляем статус из callback
		if t.Status.Progress >= 100 {
			t.Status.Status = "completed"
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"status":         t.Status.Status,
			"progress":       t.Status.Progress,
			"downloaded":     t.Status.Downloaded,
			"uploaded":       t.Status.Uploaded,
			"totalSize":      t.Status.TotalSize,
			"activePeers":    t.Status.ActivePeers,
			"downloadSpeed":  t.Status.DownloadSpeed,
			"uploadSpeed":    t.Status.UploadSpeed,
			"torrentName":    t.Status.TorrentName,
			"path":           t.Status.Path,
		})
		return
	}
}

// handleAPIStart обрабатывает API запрос запуска
func (w *WebUI) handleAPIStart(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	defer w.mu.Unlock()

	for hash, t := range w.torrents {
		if t.Engine != nil {
			t.Engine.Stop()
		}
		delete(w.torrents, hash)
	}

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

	// Проверяем тип контента
	if r.Header.Get("Content-Type") == "application/json" {
		// Старый способ - JSON с путем
		w.handleAddTorrentJSON(rw, r)
	} else {
		// Новый способ - multipart/form-data с файлом
		w.handleAddTorrentFile(rw, r)
	}
}

// handleAddTorrentJSON обрабатывает JSON запрос
func (w *WebUI) handleAddTorrentJSON(rw http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
		Size int64  `json:"size"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(rw, "Path is required", http.StatusBadRequest)
		return
	}

	w.processTorrent(rw, req.Path)
}

// handleAddTorrentFile обрабатывает загрузку файла
func (w *WebUI) handleAddTorrentFile(rw http.ResponseWriter, r *http.Request) {
	// Ограничиваем размер 10MB
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("torrent")
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to get file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаём временный файл
	tmpFile, err := os.CreateTemp("", "*.torrent")
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to create temp file: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Копируем содержимое
	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(rw, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
		return
	}
	tmpFile.Close()

	w.processTorrent(rw, tmpFile.Name())
}

// processTorrent обрабатывает торрент файл
func (w *WebUI) processTorrent(rw http.ResponseWriter, torrentPath string) {
	// Парсим .torrent файл чтобы получить имя
	torrent, err := parser.ParseFile(torrentPath)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to parse torrent file: %v", err), http.StatusBadRequest)
		return
	}

	// Проверяем валидность
	if err := torrent.IsValid(); err != nil {
		http.Error(rw, fmt.Sprintf("Invalid torrent: %v", err), http.StatusBadRequest)
		return
	}

	// Определяем директорию для загрузки
	outputDir := filepath.Dir(torrentPath)
	outputDir = filepath.Join(outputDir, torrent.Info.Name)

	// Создаём движок загрузки на базе anacrolix
	eng, err := engine.NewAnacrolixEngine(torrentPath, outputDir)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to create engine: %v", err), http.StatusBadRequest)
		return
	}

	// Загружаем торрент
	if err := eng.LoadTorrent(torrentPath); err != nil {
		http.Error(rw, fmt.Sprintf("Failed to load torrent: %v", err), http.StatusBadRequest)
		return
	}

	w.mu.Lock()

	// Создаём запись торрента
	torrentHash := torrent.Info.InfoHash
	torrentEntry := &ActiveTorrent{
		Torrent:     torrent,
		Engine:      eng,
		DownloadDir: outputDir,
		Status: DownloadStatus{
			TorrentName: torrent.Info.Name,
			Path:        torrentPath,
			Status:      "downloading",
			TotalSize:   torrent.TotalSize(),
		},
		lastUpdate: time.Now(),
		lastBytes:  0,
	}
	w.torrents[torrentHash] = torrentEntry

	w.mu.Unlock()

	// Устанавливаем callback для обновления прогресса
	eng.SetProgressCallback(func(progress float64, current, total, peers int, speed float64) {
		w.mu.Lock()
		if t, ok := w.torrents[torrentHash]; ok {
			t.Status.Progress = progress
			t.Status.ActivePeers = peers
			t.Status.DownloadSpeed = speed
			t.Status.Downloaded = int64(progress / 100.0 * float64(torrent.TotalSize()))
			t.Status.Status = "downloading"
		}
		w.mu.Unlock()
	})

	// Запускаем загрузку в фоновом режиме
	go func() {
		log.Printf("Starting download: %s", torrent.Info.Name)
		if err := eng.Start(); err != nil {
			log.Printf("Download error: %v", err)
			w.mu.Lock()
			if t, ok := w.torrents[torrentHash]; ok {
				t.Status.Status = "error"
			}
			w.mu.Unlock()
		} else {
			log.Printf("Download completed: %s", torrent.Info.Name)
			w.mu.Lock()
			if t, ok := w.torrents[torrentHash]; ok {
				t.Status.Status = "completed"
				t.Status.Progress = 100
				t.Status.Downloaded = torrent.TotalSize()
			}
			w.mu.Unlock()
		}
	}()

	addLog("Added torrent: " + torrent.Info.Name)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Torrent added successfully",
		"name":    torrent.Info.Name,
		"hash":    torrentHash,
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
	defer w.mu.Unlock()

	// Останавливаем все торренты
	for hash, t := range w.torrents {
		if t.Engine != nil {
			t.Engine.Stop()
		}
		delete(w.torrents, hash)
	}

	addLog("All torrents removed")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "All torrents removed",
	})
}

// handleAPIPause обрабатывает паузу
func (w *WebUI) handleAPIPause(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, t := range w.torrents {
		t.Status.Status = "paused"
		// В реальной версии нужно остановить engine
		if t.Engine != nil {
			t.Engine.Stop()
		}
	}

	addLog("All downloads paused")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Downloads paused",
	})
}

// handleAPIResume обрабатывает продолжение
func (w *WebUI) handleAPIResume(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, t := range w.torrents {
		t.Status.Status = "downloading"
		// В реальной версии нужно перезапустить engine
		go func(torrent *ActiveTorrent) {
			if torrent.Engine != nil {
				torrent.Engine.Start()
			}
		}(t)
	}

	addLog("All downloads resumed")

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Downloads resumed",
	})
}

// handleAPIOpenFolder открывает папку с загруженными файлами
func (w *WebUI) handleAPIOpenFolder(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.mu.RLock()
	var folderPath string
	for _, t := range w.torrents {
		if t.DownloadDir != "" {
			folderPath = t.DownloadDir
			break
		}
	}
	w.mu.RUnlock()

	if folderPath == "" {
		http.Error(rw, "No active torrent or folder not found", http.StatusNotFound)
		return
	}

	// Открываем проводник с папкой (асинхронно)
	cmd := exec.Command("explorer.exe", folderPath)
	if err := cmd.Start(); err != nil {
		// Игнорируем ошибку - explorer может запуститься асинхронно
		log.Printf("Explorer command started (may show error but folder should open): %v", err)
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Folder opened",
		"path":    folderPath,
	})
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
