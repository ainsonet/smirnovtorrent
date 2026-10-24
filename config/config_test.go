package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestConfigManagerDefaultsAndSave(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cm := NewConfigManager()
	cfg := cm.Get()

	if cfg.MaxEstablished != 80 {
		t.Errorf("Expected default MaxEstablished 80, got %d", cfg.MaxEstablished)
	}
	if !cfg.EnableDHT {
		t.Errorf("Expected default EnableDHT to be true")
	}

	// Update with new values
	cfg.MaxEstablished = 150
	cfg.ListenPort = 55555
	if err := cm.Update(cfg); err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Reload from disk
	cm2 := NewConfigManager()
	loadedCfg := cm2.Get()
	if loadedCfg.MaxEstablished != 150 {
		t.Errorf("Expected reloaded MaxEstablished to be 150, got %d", loadedCfg.MaxEstablished)
	}
	if loadedCfg.ListenPort != 55555 {
		t.Errorf("Expected reloaded ListenPort to be 55555, got %d", loadedCfg.ListenPort)
	}

	// Update with empty download dir should revert to default
	loadedCfg.DownloadDir = ""
	if err := cm2.Update(loadedCfg); err != nil {
		t.Fatalf("Failed to update config with empty download dir: %v", err)
	}
	if cm2.Get().DownloadDir == "" {
		t.Errorf("Expected DownloadDir to fallback to default, got empty")
	}
}

func TestConfigManagerConcurrentUpdates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_concurrent_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cm := NewConfigManager()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		val := i
		go func() {
			defer wg.Done()
			cfg := cm.Get()
			cfg.MaxEstablished = 50 + val
			_ = cm.Update(cfg)
		}()
		go func() {
			defer wg.Done()
			_ = cm.Get()
		}()
	}
	wg.Wait()

	// Settings file should be valid JSON
	settingsPath := filepath.Join(tempDir, "SmirnovTorrent", "settings.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Fatalf("Expected settings.json to exist: %v", err)
	}
}

func TestConfigClamping(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_clamp_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cm := NewConfigManager()
	cfg := cm.Get()
	cfg.ListenPort = 99999 // Out of range
	cfg.DownloadLimitBps = -500
	cfg.UploadLimitBps = -1000

	if err := cm.Update(cfg); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated := cm.Get()
	if updated.ListenPort != 0 {
		t.Errorf("Expected ListenPort clamped to 0, got %d", updated.ListenPort)
	}
	if updated.MaxEstablished != 80 {
		t.Errorf("Expected MaxEstablished clamped to 80, got %d", updated.MaxEstablished)
	}
	if updated.DownloadLimitBps != 0 {
		t.Errorf("Expected DownloadLimitBps clamped to 0, got %d", updated.DownloadLimitBps)
	}
	if updated.UploadLimitBps != 0 {
		t.Errorf("Expected UploadLimitBps clamped to 0, got %d", updated.UploadLimitBps)
	}
	if !updated.MinimizeToTray || !updated.CloseToTray {
		t.Errorf("Expected default tray options to be true")
	}

	// Test Load() with invalid on-disk values
	settingsPath := filepath.Join(tempDir, "SmirnovTorrent", "settings.json")
	_ = os.WriteFile(settingsPath, []byte(`{"max_established": 0, "listen_port": -10, "download_limit_bps": -200, "upload_limit_bps": -300}`), 0644)
	cmLoaded := NewConfigManager()
	loaded := cmLoaded.Get()
	if loaded.MaxEstablished != 80 {
		t.Errorf("Expected Load to clamp MaxEstablished to 80, got %d", loaded.MaxEstablished)
	}
	if loaded.ListenPort != 0 {
		t.Errorf("Expected Load to clamp ListenPort to 0, got %d", loaded.ListenPort)
	}
	if loaded.DownloadLimitBps != 0 {
		t.Errorf("Expected Load to clamp DownloadLimitBps to 0, got %d", loaded.DownloadLimitBps)
	}
	if loaded.UploadLimitBps != 0 {
		t.Errorf("Expected Load to clamp UploadLimitBps to 0, got %d", loaded.UploadLimitBps)
	}
}

func TestGetDefaultDownloadDir(t *testing.T) {
	dir := GetDefaultDownloadDir()
	if dir == "" || dir == "." {
		t.Errorf("Expected valid non-empty download dir, got %q", dir)
	}
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		t.Errorf("Expected default download dir to exist and be a directory: %v", err)
	}
}

func TestThemeAndLanguageConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_theme_lang_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cm := NewConfigManager()
	cfg := cm.Get()
	if cfg.Theme != "dark" {
		t.Errorf("Expected default theme to be 'dark', got %q", cfg.Theme)
	}
	if cfg.Language != "ru" {
		t.Errorf("Expected default language to be 'ru', got %q", cfg.Language)
	}

	// Update to light theme and English language
	cfg.Theme = "light"
	cfg.Language = "en"
	if err := cm.Update(cfg); err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Reload from disk
	cm2 := NewConfigManager()
	reloaded := cm2.Get()
	if reloaded.Theme != "light" {
		t.Errorf("Expected reloaded theme to be 'light', got %q", reloaded.Theme)
	}
	if reloaded.Language != "en" {
		t.Errorf("Expected reloaded language to be 'en', got %q", reloaded.Language)
	}

	// Clamping invalid theme and empty language
	reloaded.Theme = "neon_blue"
	reloaded.Language = ""
	if err := cm2.Update(reloaded); err != nil {
		t.Fatalf("Failed to update with invalid values: %v", err)
	}
	clamped := cm2.Get()
	if clamped.Theme != "dark" {
		t.Errorf("Expected invalid theme clamped to 'dark', got %q", clamped.Theme)
	}
	if clamped.Language != "ru" {
		t.Errorf("Expected empty language clamped to 'ru', got %q", clamped.Language)
	}

	// Ukrainian language must be rejected and clamped to default 'ru'
	if IsValidLanguage("uk") {
		t.Errorf("Expected 'uk' to not be a valid language")
	}
	reloaded.Language = "uk"
	if err := cm2.Update(reloaded); err != nil {
		t.Fatalf("Failed to update with 'uk' language: %v", err)
	}
	clampedUk := cm2.Get()
	if clampedUk.Language != "ru" {
		t.Errorf("Expected 'uk' language clamped to 'ru', got %q", clampedUk.Language)
	}

	// Legacy settings.json on disk containing 'uk' must be clamped to 'ru' upon Load() and rewritten on disk
	legacySettingsPath := filepath.Join(tempDir, "SmirnovTorrent", "settings.json")
	if err := os.WriteFile(legacySettingsPath, []byte("{\n  \"theme\": \"dark\",\n  \"language\": \"uk\"\n}"), 0644); err != nil {
		t.Fatalf("Failed to write legacy settings file: %v", err)
	}
	cmLegacy := NewConfigManager()
	if cmLegacy.Get().Language != "ru" {
		t.Errorf("Expected legacy 'uk' language clamped to 'ru' on Load(), got %q", cmLegacy.Get().Language)
	}
	savedDiskData, err := os.ReadFile(legacySettingsPath)
	if err != nil {
		t.Fatalf("Failed to read back settings file: %v", err)
	}
	if strings.Contains(string(savedDiskData), `"uk"`) {
		t.Errorf("Expected settings.json on disk to no longer contain 'uk', got: %s", string(savedDiskData))
	}
	if !strings.Contains(string(savedDiskData), `"ru"`) {
		t.Errorf("Expected settings.json on disk to contain 'ru', got: %s", string(savedDiskData))
	}
}
