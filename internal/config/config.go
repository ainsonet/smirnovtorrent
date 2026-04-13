package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config конфигурация приложения
type Config struct {
	// Peer settings
	PeerPort int `json:"peer_port"`
	MaxPeers int `json:"max_peers"`
	
	// Download settings
	OutputDir string `json:"output_dir"`
	DownloadRateLimit int64 `json:"download_rate_limit"` // bytes/sec, 0 = unlimited
	UploadRateLimit int64 `json:"upload_rate_limit"`     // bytes/sec, 0 = unlimited
	
	// Features
	EnableDHT bool `json:"enable_dht"`
	EnablePEX bool `json:"enable_pex"`
	EnableEncryption bool `json:"enable_encryption"`
	EnableResume bool `json:"enable_resume"`
	
	// Web UI
	WebUIPort int `json:"webui_port"`
	WebUIHost string `json:"webui_host"`
	EnableWebUI bool `json:"enable_webui"`
	
	// Advanced
	NumWorkers int `json:"num_workers"`
	SeedRatio float64 `json:"seed_ratio"` // 0 = seed forever
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		PeerPort: 6881,
		MaxPeers: 50,
		OutputDir: "",
		DownloadRateLimit: 0,
		UploadRateLimit: 0,
		EnableDHT: true,
		EnablePEX: true,
		EnableEncryption: true,
		EnableResume: true,
		WebUIPort: 8080,
		WebUIHost: "localhost",
		EnableWebUI: false,
		NumWorkers: 4,
		SeedRatio: 0,
	}
}

// Load загружает конфигурацию из файла
func Load(path string) (*Config, error) {
	config := DefaultConfig()
	
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Файл не существует, используем дефолтную конфигурацию
			return config, nil
		}
		return nil, err
	}
	
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	
	return config, nil
}

// Save сохраняет конфигурацию в файл
func (c *Config) Save(path string) error {
	// Создаём директорию если не существует
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// GetConfigPath возвращает путь к конфигурационному файлу
func GetConfigPath() (string, error) {
	// Проверяем текущую директорию
	localPath := "smirnovtorrent.json"
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}
	
	// Проверяем домашнюю директорию
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(homeDir, ".smirnovtorrent")
	configPath := filepath.Join(configDir, "config.json")
	
	return configPath, nil
}
