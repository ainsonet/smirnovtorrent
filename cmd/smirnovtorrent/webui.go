package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"smirnovtorrent/internal/engine"
)

// WebUI представляет веб-интерфейс
type WebUI struct {
	engine     *engine.DownloadEngine
	port       int
	mu         sync.RWMutex
	status     DownloadStatus
	tmpl       *template.Template
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
	// Создаём шаблоны
	w.tmpl = template.New("index")
	
	// Регистрируем обработчики
	http.HandleFunc("/", w.handleIndex)
	http.HandleFunc("/api/status", w.handleAPIStatus)
	http.HandleFunc("/api/start", w.handleAPIStart)
	http.HandleFunc("/api/stop", w.handleAPIStop)
	
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

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SmirnovTorrent Web UI</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            padding: 30px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }
        h1 { color: #333; margin-bottom: 10px; }
        .subtitle { color: #666; margin-bottom: 30px; }
        .status-card {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .progress-bar {
            width: 100%;
            height: 30px;
            background: #e9ecef;
            border-radius: 15px;
            overflow: hidden;
            margin: 15px 0;
        }
        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #667eea, #764ba2);
            transition: width 0.3s ease;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        .stat-item {
            background: white;
            padding: 15px;
            border-radius: 8px;
            text-align: center;
        }
        .stat-value { font-size: 24px; font-weight: bold; color: #667eea; }
        .stat-label { font-size: 12px; color: #666; margin-top: 5px; }
        .controls {
            display: flex;
            gap: 10px;
            margin-top: 20px;
        }
        button {
            flex: 1;
            padding: 12px 24px;
            border: none;
            border-radius: 6px;
            font-size: 16px;
            cursor: pointer;
            transition: all 0.2s;
        }
        .btn-start { background: #667eea; color: white; }
        .btn-start:hover { background: #5568d3; }
        .btn-stop { background: #dc3545; color: white; }
        .btn-stop:hover { background: #c82333; }
        .btn:disabled { opacity: 0.5; cursor: not-allowed; }
        .log {
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 15px;
            border-radius: 8px;
            margin-top: 20px;
            max-height: 200px;
            overflow-y: auto;
            font-family: 'Courier New', monospace;
            font-size: 13px;
        }
        .status-indicator {
            display: inline-block;
            width: 10px;
            height: 10px;
            border-radius: 50%;
            margin-right: 8px;
        }
        .status-running { background: #28a745; }
        .status-stopped { background: #dc3545; }
        .status-paused { background: #ffc107; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌊 SmirnovTorrent</h1>
        <p class="subtitle">Web Interface</p>
        
        <div class="status-card">
            <h2><span class="status-indicator status-running" id="statusIndicator"></span><span id="torrentName">No active torrent</span></h2>
            
            <div class="progress-bar">
                <div class="progress-fill" id="progressFill" style="width: 0%">0%</div>
            </div>
            
            <div class="stats">
                <div class="stat-item">
                    <div class="stat-value" id="downloaded">0 MB</div>
                    <div class="stat-label">Downloaded</div>
                </div>
                <div class="stat-item">
                    <div class="stat-value" id="totalSize">0 MB</div>
                    <div class="stat-label">Total Size</div>
                </div>
                <div class="stat-item">
                    <div class="stat-value" id="peers">0</div>
                    <div class="stat-label">Active Peers</div>
                </div>
                <div class="stat-item">
                    <div class="stat-value" id="downloadSpeed">0 KB/s</div>
                    <div class="stat-label">Download Speed</div>
                </div>
            </div>
        </div>
        
        <div class="controls">
            <button class="btn-start" onclick="startDownload()" id="startBtn">Start Download</button>
            <button class="btn-stop" onclick="stopDownload()" id="stopBtn" disabled>Stop</button>
        </div>
        
        <div class="log" id="log">
            <div>Web UI initialized...</div>
        </div>
    </div>
    
    <script>
        let updateInterval;
        
        function formatBytes(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }
        
        async function updateStatus() {
            try {
                const response = await fetch('/api/status');
                const data = await response.json();
                
                document.getElementById('progressFill').style.width = data.progress + '%';
                document.getElementById('progressFill').textContent = data.progress.toFixed(1) + '%';
                document.getElementById('downloaded').textContent = formatBytes(data.downloaded);
                document.getElementById('totalSize').textContent = formatBytes(data.totalSize);
                document.getElementById('peers').textContent = data.activePeers;
                document.getElementById('downloadSpeed').textContent = formatBytes(data.downloadSpeed) + '/s';
                document.getElementById('torrentName').textContent = data.torrentName || 'No active torrent';
                
                addLog('Status updated: ' + data.progress.toFixed(1) + '%');
            } catch (error) {
                addLog('Error updating status: ' + error.message);
            }
        }
        
        async function startDownload() {
            try {
                const response = await fetch('/api/start', { method: 'POST' });
                const data = await response.json();
                addLog('Start: ' + data.message);
                
                document.getElementById('startBtn').disabled = true;
                document.getElementById('stopBtn').disabled = false;
                
                if (!updateInterval) {
                    updateInterval = setInterval(updateStatus, 1000);
                }
            } catch (error) {
                addLog('Error starting: ' + error.message);
            }
        }
        
        async function stopDownload() {
            try {
                const response = await fetch('/api/stop', { method: 'POST' });
                const data = await response.json();
                addLog('Stop: ' + data.message);
                
                document.getElementById('startBtn').disabled = false;
                document.getElementById('stopBtn').disabled = true;
                
                if (updateInterval) {
                    clearInterval(updateInterval);
                    updateInterval = null;
                }
            } catch (error) {
                addLog('Error stopping: ' + error.message);
            }
        }
        
        function addLog(message) {
            const log = document.getElementById('log');
            const time = new Date().toLocaleTimeString();
            log.innerHTML += '<div>[' + time + '] ' + message + '</div>';
            log.scrollTop = log.scrollHeight;
        }
        
        // Auto-update status every 2 seconds
        updateInterval = setInterval(updateStatus, 2000);
        updateStatus();
    </script>
</body>
</html>`

	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprint(rw, html)
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
