package config

import (
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.PeerPort != 6881 {
		t.Errorf("Expected peer port 6881, got %d", config.PeerPort)
	}
	
	if config.MaxPeers != 50 {
		t.Errorf("Expected max peers 50, got %d", config.MaxPeers)
	}
	
	if config.EnableDHT != true {
		t.Error("Expected DHT enabled by default")
	}
	
	if config.EnablePEX != true {
		t.Error("Expected PEX enabled by default")
	}
	
	if config.EnableEncryption != true {
		t.Error("Expected encryption enabled by default")
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	config, err := Load("/nonexistent/path/config.json")
	
	if err != nil {
		t.Fatalf("Load should not fail for non-existent file: %v", err)
	}
	
	if config == nil {
		t.Fatal("Expected default config for non-existent file")
	}
}

func TestLoadAndSave(t *testing.T) {
	// Создаём временный файл
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.json")
	
	// Создаём и сохраняем конфигурацию
	config := DefaultConfig()
	config.PeerPort = 9999
	config.MaxPeers = 100
	config.EnableDHT = false
	
	if err := config.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Загружаем конфигурацию
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Проверяем значения
	if loaded.PeerPort != 9999 {
		t.Errorf("Expected peer port 9999, got %d", loaded.PeerPort)
	}
	
	if loaded.MaxPeers != 100 {
		t.Errorf("Expected max peers 100, got %d", loaded.MaxPeers)
	}
	
	if loaded.EnableDHT != false {
		t.Error("Expected DHT disabled")
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}
	
	if path == "" {
		t.Error("Expected non-empty config path")
	}
}

func TestConfigSaveInvalidPath(t *testing.T) {
	config := DefaultConfig()
	
	// Пытаемся сохранить в недопустимый путь (пустой файл)
	err := config.Save("")
	
	// Ожидаем ошибку
	if err == nil {
		t.Error("Expected error when saving to empty path")
	}
}
