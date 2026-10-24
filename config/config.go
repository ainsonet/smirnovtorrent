package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AppConfig stores application settings
type AppConfig struct {
	DownloadDir       string `json:"download_dir"`
	ListenPort        int    `json:"listen_port"`
	EnableDHT         bool   `json:"enable_dht"`
	EnablePEX         bool   `json:"enable_pex"`
	EnableUPnP        bool   `json:"enable_upnp"`
	MaxEstablished    int    `json:"max_established"`
	DownloadLimitBps  int64  `json:"download_limit_bps"`
	UploadLimitBps    int64  `json:"upload_limit_bps"`
	AssociateTorrents bool   `json:"associate_torrents"`
	AssociateMagnets  bool   `json:"associate_magnets"`
	Autostart         bool   `json:"autostart"`
	MinimizeToTray    bool   `json:"minimize_to_tray"`
	CloseToTray       bool   `json:"close_to_tray"`
	StartMinimized    bool   `json:"start_minimized"`
	Theme             string `json:"theme"`
	Language          string `json:"language"`
}

type ConfigManager struct {
	mu         sync.RWMutex
	fileMu     sync.Mutex
	cfg        AppConfig
	configPath string
}

func GetDefaultDownloadDir() string {
	return getDefaultDownloadDirPlatform()
}

func GetConfigDir() string {
	appData := os.Getenv("APPDATA")
	if appData != "" {
		dir := filepath.Join(appData, "SmirnovTorrent")
		_ = os.MkdirAll(dir, 0755)
		return dir
	}
	if ucd, err := os.UserConfigDir(); err == nil && ucd != "" {
		dir := filepath.Join(ucd, "SmirnovTorrent")
		_ = os.MkdirAll(dir, 0755)
		return dir
	}
	dir := filepath.Join(".", ".smirnovtorrent")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func NewConfigManager() *ConfigManager {
	configDir := GetConfigDir()
	configPath := filepath.Join(configDir, "settings.json")

	defaultCfg := AppConfig{
		DownloadDir:       GetDefaultDownloadDir(),
		ListenPort:        0, // 0 means random available port
		EnableDHT:         true,
		EnablePEX:         true,
		EnableUPnP:        true,
		MaxEstablished:    80,
		DownloadLimitBps:  0, // 0 = unlimited
		UploadLimitBps:    0,
		AssociateTorrents: false,
		AssociateMagnets:  false,
		Autostart:         false,
		MinimizeToTray:    true,
		CloseToTray:       true,
		StartMinimized:    false,
		Theme:             "dark",
		Language:          "ru",
	}

	cm := &ConfigManager{
		cfg:        defaultCfg,
		configPath: configPath,
	}

	cm.Load()
	return cm
}

// IsValidLanguage checks whether a given language code is supported.
func IsValidLanguage(lang string) bool {
	switch lang {
	case "ru", "en", "de", "fr", "es", "it", "pt", "zh", "ja", "pl", "tr":
		return true
	default:
		return false
	}
}

func (cm *ConfigManager) Load() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return
	}

	var loaded AppConfig
	if err := json.Unmarshal(data, &loaded); err == nil {
		if loaded.DownloadDir == "" {
			loaded.DownloadDir = GetDefaultDownloadDir()
		}
		if loaded.MaxEstablished <= 0 {
			loaded.MaxEstablished = 80
		}
		if loaded.ListenPort < 0 || loaded.ListenPort > 65535 {
			loaded.ListenPort = 0
		}
		if loaded.DownloadLimitBps < 0 {
			loaded.DownloadLimitBps = 0
		}
		if loaded.UploadLimitBps < 0 {
			loaded.UploadLimitBps = 0
		}
		if loaded.Theme != "light" && loaded.Theme != "dark" {
			loaded.Theme = "dark"
		}
		clampedLang := false
		if !IsValidLanguage(loaded.Language) {
			loaded.Language = "ru"
			clampedLang = true
		}
		cm.cfg = loaded
		if clampedLang {
			if marshaled, err := json.MarshalIndent(cm.cfg, "", "  "); err == nil {
				_ = atomicWriteFile(cm.configPath, marshaled, 0644)
			}
		}
	}
}

func (cm *ConfigManager) Save() error {
	cm.fileMu.Lock()
	defer cm.fileMu.Unlock()

	cm.mu.RLock()
	data, err := json.MarshalIndent(cm.cfg, "", "  ")
	cm.mu.RUnlock()
	if err != nil {
		return err
	}

	return atomicWriteFile(cm.configPath, data, 0644)
}

func (cm *ConfigManager) Get() AppConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cfg
}

func (cm *ConfigManager) Update(newCfg AppConfig) error {
	cm.fileMu.Lock()
	defer cm.fileMu.Unlock()

	cm.mu.Lock()
	if newCfg.DownloadDir == "" {
		newCfg.DownloadDir = GetDefaultDownloadDir()
	}
	if newCfg.ListenPort < 0 || newCfg.ListenPort > 65535 {
		newCfg.ListenPort = 0
	}
	if newCfg.MaxEstablished <= 0 {
		newCfg.MaxEstablished = 80
	}
	if newCfg.DownloadLimitBps < 0 {
		newCfg.DownloadLimitBps = 0
	}
	if newCfg.UploadLimitBps < 0 {
		newCfg.UploadLimitBps = 0
	}
	if newCfg.Theme != "light" && newCfg.Theme != "dark" {
		newCfg.Theme = "dark"
	}
	if !IsValidLanguage(newCfg.Language) {
		newCfg.Language = "ru"
	}
	cm.cfg = newCfg
	data, err := json.MarshalIndent(cm.cfg, "", "  ")
	cm.mu.Unlock()
	if err != nil {
		return err
	}

	return atomicWriteFile(cm.configPath, data, 0644)
}

func atomicWriteFile(filename string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	_ = os.MkdirAll(dir, 0755)

	tmpFile, err := os.CreateTemp(dir, "tmp-*.json")
	if err != nil {
		return os.WriteFile(filename, data, perm)
	}
	tmpName := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	_ = os.Chmod(tmpName, perm)

	var renameErr error
	for i := 0; i < 5; i++ {
		renameErr = os.Rename(tmpName, filename)
		if renameErr == nil {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	if writeErr := os.WriteFile(filename, data, perm); writeErr == nil {
		return nil
	}
	return renameErr
}
