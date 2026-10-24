package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"smirnovtorrent/config"
	"smirnovtorrent/engine"
	"smirnovtorrent/ui"
	"smirnovtorrent/window"
)

func tryForwardToRunningInstance(portFile string) bool {
	data, err := os.ReadFile(portFile)
	if err != nil {
		return false
	}
	portStr := strings.TrimSpace(string(data))
	if portStr == "" {
		return false
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%s/api/stats", portStr))
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			_ = resp.Body.Close()
		}
		return false
	}
	_ = resp.Body.Close()

	// Forward any arguments (magnet links or torrent files)
	for _, arg := range os.Args[1:] {
		trimmed := strings.Trim(strings.TrimSpace(arg), `"'`)
		if strings.HasPrefix(strings.ToLower(trimmed), "magnet:") {
			body, _ := json.Marshal(map[string]string{"magnet_uri": trimmed})
			if resp, err := client.Post(fmt.Sprintf("http://127.0.0.1:%s/api/add-magnet", portStr), "application/json", bytes.NewReader(body)); err == nil && resp != nil {
				_ = resp.Body.Close()
			}
		} else if strings.HasSuffix(strings.ToLower(trimmed), ".torrent") {
			absPath, err := filepath.Abs(trimmed)
			if err == nil {
				trimmed = absPath
			}
			body, _ := json.Marshal(map[string]string{"file_path": trimmed})
			if resp, err := client.Post(fmt.Sprintf("http://127.0.0.1:%s/api/add-torrent-path", portStr), "application/json", bytes.NewReader(body)); err == nil && resp != nil {
				_ = resp.Body.Close()
			}
		}
	}

	hasMinimizedArg := false
	for _, arg := range os.Args[1:] {
		trimmed := strings.Trim(strings.TrimSpace(arg), `"'`)
		if trimmed == "--minimized" || trimmed == "--autostart" {
			hasMinimizedArg = true
			break
		}
	}

	// Bring running window to front unless started minimized
	if !hasMinimizedArg {
		if resp, err := client.Post(fmt.Sprintf("http://127.0.0.1:%s/api/window/show", portStr), "application/json", nil); err == nil && resp != nil {
			_ = resp.Body.Close()
		}
	}
	return true
}

func main() {
	window.SetAppUserModelID()

	configDir := config.GetConfigDir()
	portFile := filepath.Join(configDir, "app.port")

	// Check if SmirnovTorrent is already running; if so, forward arguments and focus existing window
	if tryForwardToRunningInstance(portFile) {
		log.Println("[SmirnovTorrent] Forwarded request to running instance. Exiting.")
		return
	}

	log.Println("[SmirnovTorrent] Starting application...")

	// 1. Load Configuration
	cfgManager := config.NewConfigManager()
	log.Printf("[SmirnovTorrent] Download directory: %s", cfgManager.Get().DownloadDir)

	// 2. Start Torrent Engine
	eng, err := engine.NewEngine(cfgManager)
	if err != nil {
		log.Fatalf("[SmirnovTorrent] Failed to initialize engine: %v", err)
	}
	defer eng.Close()

	// 3. Start Embedded Local UI Server
	server, err := ui.NewServer(eng, cfgManager)
	if err != nil {
		log.Fatalf("[SmirnovTorrent] Failed to initialize UI server: %v", err)
	}
	defer server.Close()

	go func() {
		if err := server.Start(); err != nil {
			log.Printf("[SmirnovTorrent] UI Server stopped: %v", err)
		}
	}()

	// Wait briefly for server to bind and save port file
	time.Sleep(100 * time.Millisecond)
	appURL := server.GetURL()
	log.Printf("[SmirnovTorrent] UI ready at %s", appURL)

	_ = os.WriteFile(portFile, []byte(fmt.Sprintf("%d", server.GetPort())), 0644)
	defer func() {
		_ = os.Remove(portFile)
	}()

	// 4. Handle graceful shutdown signals (Ctrl+C if run in terminal)
	var win *window.DesktopWindow
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("[SmirnovTorrent] Shutting down...")
		if win != nil {
			win.Destroy()
		}
		_ = os.Remove(portFile)
		_ = server.Close()
		eng.Close()
		os.Exit(0)
	}()

	// Handle initial command line arguments (e.g. launched via file association or magnet link)
	startMinimized := false
	for _, arg := range os.Args[1:] {
		trimmed := strings.Trim(strings.TrimSpace(arg), `"'`)
		if trimmed == "--minimized" || trimmed == "--autostart" {
			startMinimized = true
		} else if strings.HasPrefix(strings.ToLower(trimmed), "magnet:") {
			_, _ = eng.AddMagnet(trimmed)
		} else if strings.HasSuffix(strings.ToLower(trimmed), ".torrent") {
			absPath, err := filepath.Abs(trimmed)
			if err == nil {
				trimmed = absPath
			}
			if _, err := os.Stat(trimmed); err == nil {
				_, _ = eng.AddTorrentFile(trimmed)
			}
		}
	}

	// 5. Create and run native Desktop Window
	win = window.NewDesktopWindow("SmirnovTorrent", appURL, 1480, 850, false, cfgManager)
	if win == nil {
		// Fallback: If WebView2 fails to initialize, notify user with URL and open default browser
		fmt.Printf("SmirnovTorrent is running in browser mode: %s\n", appURL)
		_ = ui.OpenURL(appURL)
		select {}
	}

	defer win.Destroy()

	server.SetOnWindowShow(func() {
		win.Show()
	})
	server.SetOnThemeChange(func(theme string) {
		win.SetTheme(theme)
	})

	if startMinimized || cfgManager.Get().StartMinimized {
		win.Hide()
		log.Println("[SmirnovTorrent] Started minimized in system tray.")
	} else {
		log.Println("[SmirnovTorrent] Desktop window opened.")
	}

	win.Run()

	log.Println("[SmirnovTorrent] Exited normally.")
}
