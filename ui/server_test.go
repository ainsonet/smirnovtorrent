package ui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"smirnovtorrent/config"
	"smirnovtorrent/engine"
)

func TestStaticAssets(t *testing.T) {
	// Verify index.html contains updated data-sort keys
	strippedFS, err := assetsFS.ReadFile("assets/index.html")
	if err != nil {
		t.Fatalf("Failed to read assets/index.html: %v", err)
	}
	indexContent := string(strippedFS)
	if !strings.Contains(indexContent, `data-sort="download_speed"`) {
		t.Errorf("Expected index.html to contain data-sort=\"download_speed\"")
	}
	if !strings.Contains(indexContent, `data-sort="upload_speed"`) {
		t.Errorf("Expected index.html to contain data-sort=\"upload_speed\"")
	}
	if !strings.Contains(indexContent, `data-sort="peers_connected"`) {
		t.Errorf("Expected index.html to contain data-sort=\"peers_connected\"")
	}
	if !strings.Contains(indexContent, `data-sort="seeds_connected"`) {
		t.Errorf("Expected index.html to contain data-sort=\"seeds_connected\"")
	}
	if !strings.Contains(indexContent, `data-sort="added_at"`) {
		t.Errorf("Expected index.html to contain data-sort=\"added_at\"")
	}
	if !strings.Contains(indexContent, `id="btn-verify"`) {
		t.Errorf("Expected index.html to contain id=\"btn-verify\"")
	}
	if !strings.Contains(indexContent, `id="setting-dl-limit"`) {
		t.Errorf("Expected index.html to contain id=\"setting-dl-limit\"")
	}
	if !strings.Contains(indexContent, `id="setting-ul-limit"`) {
		t.Errorf("Expected index.html to contain id=\"setting-ul-limit\"")
	}
	if !strings.Contains(indexContent, `id="setting-assoc-torrent"`) {
		t.Errorf("Expected index.html to contain id=\"setting-assoc-torrent\"")
	}
	if !strings.Contains(indexContent, `id="setting-assoc-magnet"`) {
		t.Errorf("Expected index.html to contain id=\"setting-assoc-magnet\"")
	}
	if !strings.Contains(indexContent, `id="setting-autostart"`) {
		t.Errorf("Expected index.html to contain id=\"setting-autostart\"")
	}
	if !strings.Contains(indexContent, `id="setting-minimize-tray"`) {
		t.Errorf("Expected index.html to contain id=\"setting-minimize-tray\"")
	}
	if !strings.Contains(indexContent, `id="setting-close-tray"`) {
		t.Errorf("Expected index.html to contain id=\"setting-close-tray\"")
	}
	if !strings.Contains(indexContent, `id="setting-start-minimized"`) {
		t.Errorf("Expected index.html to contain id=\"setting-start-minimized\"")
	}

	if strings.Contains(indexContent, `id="btn-add-file"`) {
		t.Errorf("Expected index.html to NOT contain id=\"btn-add-file\" in toolbar")
	}
	if strings.Contains(indexContent, `id="btn-add-magnet"`) {
		t.Errorf("Expected index.html to NOT contain id=\"btn-add-magnet\" in toolbar")
	}
	if !strings.Contains(indexContent, `id="filter-input"`) {
		t.Errorf("Expected index.html to contain id=\"filter-input\"")
	}
	if !strings.Contains(indexContent, `onclick="addTorrentFile()"`) {
		t.Errorf("Expected index.html to contain empty state addTorrentFile button")
	}
	if !strings.Contains(indexContent, `onclick="openModal('modal-magnet')"`) {
		t.Errorf("Expected index.html to contain empty state openModal('modal-magnet') button")
	}

	// Verify logo is embedded and configured
	logoBytes, err := assetsFS.ReadFile("assets/smirnovtorrent.jpg")
	if err != nil {
		t.Fatalf("Failed to read assets/smirnovtorrent.jpg: %v", err)
	}
	if len(logoBytes) == 0 {
		t.Errorf("Expected assets/smirnovtorrent.jpg to have non-zero size")
	}
	if !strings.Contains(indexContent, `src="smirnovtorrent.jpg"`) {
		t.Errorf("Expected index.html to contain brand logo src=\"smirnovtorrent.jpg\"")
	}
	if !strings.Contains(indexContent, `href="smirnovtorrent.jpg"`) {
		t.Errorf("Expected index.html to contain favicon href=\"smirnovtorrent.jpg\"")
	}
	if strings.Contains(indexContent, `points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"`) {
		t.Errorf("Expected index.html to have removed the lightning polygon")
	}

	// Verify style.css contains .status-missing and .status-checking
	cssBytes, err := assetsFS.ReadFile("assets/style.css")
	if err != nil {
		t.Fatalf("Failed to read assets/style.css: %v", err)
	}
	cssContent := string(cssBytes)
	if !strings.Contains(cssContent, ".status-missing") {
		t.Errorf("Expected style.css to contain .status-missing")
	}
	if !strings.Contains(cssContent, ".status-checking") {
		t.Errorf("Expected style.css to contain .status-checking")
	}
	if !strings.Contains(cssContent, "flex-shrink: 0") {
		t.Errorf("Expected style.css to contain flex-shrink: 0")
	}
	if !strings.Contains(cssContent, "min-width: 200px") {
		t.Errorf("Expected style.css to contain min-width: 200px")
	}
	if !strings.Contains(cssContent, "@media (max-width: 1350px)") {
		t.Errorf("Expected style.css to contain @media (max-width: 1350px)")
	}

	// Verify app.js contains clearDetails, sortField: 'added_at', and btnVerify
	jsBytes, err := assetsFS.ReadFile("assets/app.js")
	if err != nil {
		t.Fatalf("Failed to read assets/app.js: %v", err)
	}
	jsContent := string(jsBytes)
	if !strings.Contains(jsContent, "clearDetails()") {
		t.Errorf("Expected app.js to contain clearDetails()")
	}
	if !strings.Contains(jsContent, "sortField: 'added_at'") {
		t.Errorf("Expected app.js to default sortField to 'added_at'")
	}
	if !strings.Contains(jsContent, "btnVerify") {
		t.Errorf("Expected app.js to contain btnVerify")
	}
	if !strings.Contains(jsContent, "/api/torrent/verify") {
		t.Errorf("Expected app.js to call /api/torrent/verify")
	}
	if !strings.Contains(jsContent, "applyTheme(") {
		t.Errorf("Expected app.js to contain applyTheme")
	}
	if !strings.Contains(jsContent, "applyLanguage(") {
		t.Errorf("Expected app.js to contain applyLanguage")
	}

	// Verify light logo is embedded and non-zero
	lightLogoBytes, err := assetsFS.ReadFile("assets/smirnovtorrent-light.jpg")
	if err != nil {
		t.Fatalf("Failed to read assets/smirnovtorrent-light.jpg: %v", err)
	}
	if len(lightLogoBytes) == 0 {
		t.Errorf("Expected assets/smirnovtorrent-light.jpg to have non-zero size")
	}

	// Verify theme and language elements in index.html
	if !strings.Contains(indexContent, `id="setting-theme"`) {
		t.Errorf("Expected index.html to contain id=\"setting-theme\"")
	}
	if !strings.Contains(indexContent, `id="setting-language"`) {
		t.Errorf("Expected index.html to contain id=\"setting-language\"")
	}
	if !strings.Contains(indexContent, `value="light"`) {
		t.Errorf("Expected index.html to contain light theme option")
	}
	if !strings.Contains(indexContent, `value="en"`) {
		t.Errorf("Expected index.html to contain English language option")
	}
	if strings.Contains(indexContent, `value="uk"`) || strings.Contains(indexContent, "Ukrainian") || strings.Contains(indexContent, "Українська") {
		t.Errorf("Expected index.html to not contain Ukrainian language option")
	}
	if !strings.Contains(indexContent, `l === 'uk'`) {
		t.Errorf("Expected index.html script to sanitize legacy 'uk' language in localStorage")
	}
	if strings.Contains(jsContent, `"uk":`) || strings.Contains(jsContent, `'uk-UA'`) {
		t.Errorf("Expected app.js to not contain Ukrainian language dictionary")
	}

	// Verify language country emojis
	expectedFlags := []string{
		"🇷🇺 Русский",
		"🇬🇧 English",
		"🇩🇪 Deutsch",
		"🇫🇷 Français",
		"🇪🇸 Español",
		"🇮🇹 Italiano",
		"🇵🇹 Português",
		"🇨🇳 简体中文",
		"🇯🇵 日本語",
		"🇵🇱 Polski",
		"🇹🇷 Türkçe",
	}
	for _, flag := range expectedFlags {
		if !strings.Contains(indexContent, flag) {
			t.Errorf("Expected index.html to contain language option with country emoji: %s", flag)
		}
	}

	// Verify style.css contains expanded col-added and light theme rules
	if !strings.Contains(cssContent, `[data-theme="light"]`) {
		t.Errorf("Expected style.css to contain [data-theme=\"light\"]")
	}
	if !strings.Contains(cssContent, `.col-added { min-width: 195px`) {
		t.Errorf("Expected style.css to contain expanded .col-added min-width: 195px")
	}
	if !strings.Contains(cssContent, `[data-theme="light"] .sub-table tbody tr.selected`) {
		t.Errorf("Expected style.css to override selected sub-table row text color in light theme")
	}
	if !strings.Contains(cssContent, `[data-theme="light"] .sep`) {
		t.Errorf("Expected style.css to override status bar separator color in light theme")
	}

	// Verify index.html early logo swap script to prevent black logo flash on light theme
	if !strings.Contains(indexContent, `id="brand-logo"`) || !strings.Contains(indexContent, `smirnovtorrent-light.jpg`) {
		t.Errorf("Expected index.html to contain brand logo inline swap for light theme")
	}

	// Verify app.js contains in-place date updates and settings cancel reverting
	if !strings.Contains(jsContent, `tr.children[10]`) {
		t.Errorf("Expected app.js to update col-added (tr.children[10]) in renderTorrents")
	}
	if !strings.Contains(jsContent, `settingsInitialTheme`) {
		t.Errorf("Expected app.js to track settingsInitialTheme for clean cancel reverting")
	}

	// Verify TwemojiCountryFlags font is embedded and referenced in CSS and HTML
	fontBytes, err := assetsFS.ReadFile("assets/TwemojiCountryFlags.woff2")
	if err != nil {
		t.Fatalf("Failed to read assets/TwemojiCountryFlags.woff2: %v", err)
	}
	if len(fontBytes) == 0 {
		t.Errorf("Expected assets/TwemojiCountryFlags.woff2 to have non-zero size")
	}
	if !strings.Contains(cssContent, "Twemoji Country Flags") {
		t.Errorf("Expected style.css to define and use Twemoji Country Flags font")
	}
	if !strings.Contains(indexContent, "TwemojiCountryFlags.woff2") {
		t.Errorf("Expected index.html to preload TwemojiCountryFlags.woff2")
	}
}

func TestServerSettingsHandler(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_settings_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.ListenPort = 0
	_ = cfgMgr.Update(cfg)

	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// GET /api/settings
	reqGet := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	wGet := httptest.NewRecorder()
	srv.handleSettings(wGet, reqGet)
	if wGet.Code != http.StatusOK {
		t.Fatalf("GET /api/settings returned status %d", wGet.Code)
	}

	var curCfg config.AppConfig
	if err := json.Unmarshal(wGet.Body.Bytes(), &curCfg); err != nil {
		t.Fatalf("Failed to unmarshal GET settings: %v", err)
	}

	if curCfg.Theme != "dark" {
		t.Errorf("Expected default Theme to be 'dark', got %q", curCfg.Theme)
	}
	if curCfg.Language != "ru" {
		t.Errorf("Expected default Language to be 'ru', got %q", curCfg.Language)
	}

	themeChanged := false
	var updatedTheme string
	srv.SetOnThemeChange(func(th string) {
		themeChanged = true
		updatedTheme = th
	})

	// POST /api/settings with updated max_established, theme, language
	curCfg.MaxEstablished = 120
	curCfg.Theme = "light"
	curCfg.Language = "en"
	body, _ := json.Marshal(curCfg)
	reqPost := httptest.NewRequest(http.MethodPost, "/api/settings", bytes.NewReader(body))
	wPost := httptest.NewRecorder()
	srv.handleSettings(wPost, reqPost)
	if wPost.Code != http.StatusOK {
		t.Fatalf("POST /api/settings returned status %d", wPost.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(wPost.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to unmarshal POST response: %v", err)
	}
	if res["success"] != true {
		t.Errorf("Expected success: true, got %v", res)
	}

	if cfgMgr.Get().MaxEstablished != 120 {
		t.Errorf("Expected MaxEstablished to be 120, got %d", cfgMgr.Get().MaxEstablished)
	}
	if cfgMgr.Get().Theme != "light" {
		t.Errorf("Expected Theme to be 'light', got %q", cfgMgr.Get().Theme)
	}
	if cfgMgr.Get().Language != "en" {
		t.Errorf("Expected Language to be 'en', got %q", cfgMgr.Get().Language)
	}
	if !themeChanged || updatedTheme != "light" {
		t.Errorf("Expected SetOnThemeChange callback to fire with 'light', got called=%v val=%q", themeChanged, updatedTheme)
	}

	// POST /api/settings with uk language must be rejected/clamped to ru
	curCfg.Language = "uk"
	bodyUk, _ := json.Marshal(curCfg)
	reqPostUk := httptest.NewRequest(http.MethodPost, "/api/settings", bytes.NewReader(bodyUk))
	wPostUk := httptest.NewRecorder()
	srv.handleSettings(wPostUk, reqPostUk)
	if wPostUk.Code != http.StatusOK {
		t.Fatalf("POST /api/settings with 'uk' returned status %d", wPostUk.Code)
	}
	if cfgMgr.Get().Language != "ru" {
		t.Errorf("Expected Language to be clamped to 'ru', got %q", cfgMgr.Get().Language)
	}

	// GET /api/settings should also return ru
	reqGet2 := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	wGet2 := httptest.NewRecorder()
	srv.handleSettings(wGet2, reqGet2)
	var reloadedCfg config.AppConfig
	if err := json.Unmarshal(wGet2.Body.Bytes(), &reloadedCfg); err != nil {
		t.Fatalf("Failed to unmarshal GET settings: %v", err)
	}
	if reloadedCfg.Language != "ru" {
		t.Errorf("Expected GET settings language to be 'ru', got %q", reloadedCfg.Language)
	}
}

func TestServerGetTorrentsAndStats(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_api_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.ListenPort = 0
	_ = cfgMgr.Update(cfg)

	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// GET /api/torrents
	req := httptest.NewRequest(http.MethodGet, "/api/torrents", nil)
	w := httptest.NewRecorder()
	srv.handleGetTorrents(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/torrents returned status %d", w.Code)
	}

	// GET /api/stats
	reqStats := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	wStats := httptest.NewRecorder()
	srv.handleGetStats(wStats, reqStats)
	if wStats.Code != http.StatusOK {
		t.Fatalf("GET /api/stats returned status %d", wStats.Code)
	}
}

func TestServerOpenFolderNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_open_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.ListenPort = 0
	_ = cfgMgr.Update(cfg)

	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Non-existent hash should return JSON success: false
	req := httptest.NewRequest(http.MethodPost, "/api/open-folder?hash=invalidhash", nil)
	w := httptest.NewRecorder()
	srv.handleOpenFolder(w, req)

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if res["success"] != false {
		t.Errorf("Expected success: false for non-existent hash, got %v", res)
	}
}

func TestExplorerDirectoryDetection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "explorer_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Directory check
	fiDir, err := os.Stat(tempDir)
	if err != nil || !fiDir.IsDir() {
		t.Fatalf("Expected tempDir to be directory")
	}

	// File check
	filePath := filepath.Join(tempDir, "file.txt")
	_ = os.WriteFile(filePath, []byte("data"), 0644)
	fiFile, err := os.Stat(filePath)
	if err != nil || fiFile.IsDir() {
		t.Fatalf("Expected filePath to not be directory")
	}
}

func TestServerMissingHashHandlers(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_missing_hash_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	// 1. handlePause without hash
	reqPause := httptest.NewRequest(http.MethodPost, "/api/torrent/pause", nil)
	wPause := httptest.NewRecorder()
	srv.handlePause(wPause, reqPause)
	var res map[string]interface{}
	_ = json.Unmarshal(wPause.Body.Bytes(), &res)
	if res["success"] != false {
		t.Errorf("Expected success: false for empty hash in pause")
	}

	// 2. handleResume without hash
	reqResume := httptest.NewRequest(http.MethodPost, "/api/torrent/resume", nil)
	wResume := httptest.NewRecorder()
	srv.handleResume(wResume, reqResume)
	res = nil
	_ = json.Unmarshal(wResume.Body.Bytes(), &res)
	if res["success"] != false {
		t.Errorf("Expected success: false for empty hash in resume")
	}

	// 3. handleDelete without hash
	reqDelete := httptest.NewRequest(http.MethodPost, "/api/torrent/delete", nil)
	wDelete := httptest.NewRecorder()
	srv.handleDelete(wDelete, reqDelete)
	res = nil
	_ = json.Unmarshal(wDelete.Body.Bytes(), &res)
	if res["success"] != false {
		t.Errorf("Expected success: false for empty hash in delete")
	}

	// 4. handleSequential without hash
	reqSeq := httptest.NewRequest(http.MethodPost, "/api/torrent/sequential", nil)
	wSeq := httptest.NewRecorder()
	srv.handleSequential(wSeq, reqSeq)
	res = nil
	_ = json.Unmarshal(wSeq.Body.Bytes(), &res)
	if res["success"] != false {
		t.Errorf("Expected success: false for empty hash in sequential")
	}

	// 5. handleVerify without hash
	reqVerify := httptest.NewRequest(http.MethodPost, "/api/torrent/verify", nil)
	wVerify := httptest.NewRecorder()
	srv.handleVerify(wVerify, reqVerify)
	res = nil
	_ = json.Unmarshal(wVerify.Body.Bytes(), &res)
	if res["success"] != false {
		t.Errorf("Expected success: false for empty hash in verify")
	}
}

func TestServerAddMagnetValidation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_magnet_val_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	// Empty magnet
	bodyEmpty, _ := json.Marshal(AddMagnetRequest{MagnetURI: "   "})
	req := httptest.NewRequest(http.MethodPost, "/api/add-magnet", bytes.NewReader(bodyEmpty))
	w := httptest.NewRecorder()
	srv.handleAddMagnet(w, req)

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	if res["success"] != false {
		t.Errorf("Expected success: false for empty magnet URI")
	}
}

func TestServerSystemIntegrationHandler(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_sys_integ_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	// GET /api/system/integration
	reqGet := httptest.NewRequest(http.MethodGet, "/api/system/integration", nil)
	wGet := httptest.NewRecorder()
	srv.handleSystemIntegration(wGet, reqGet)
	if wGet.Code != http.StatusOK {
		t.Fatalf("GET /api/system/integration returned status %d", wGet.Code)
	}

	// POST /api/system/integration with invalid JSON
	reqBad := httptest.NewRequest(http.MethodPost, "/api/system/integration", bytes.NewReader([]byte("{invalid")))
	wBad := httptest.NewRecorder()
	srv.handleSystemIntegration(wBad, reqBad)
	var resBad map[string]interface{}
	_ = json.Unmarshal(wBad.Body.Bytes(), &resBad)
	if resBad["success"] != false {
		t.Errorf("Expected success: false for invalid json in system integration")
	}
}

func TestServerSettingsRateLimits(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_settings_limits_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	cur := cfgMgr.Get()
	cur.DownloadLimitBps = 1048576 // 1MB/s
	cur.UploadLimitBps = 524288    // 512KB/s
	cur.MinimizeToTray = true
	cur.CloseToTray = true

	body, _ := json.Marshal(cur)
	req := httptest.NewRequest(http.MethodPost, "/api/settings", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.handleSettings(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/settings returned status %d", w.Code)
	}

	dl, ul := eng.GetRateLimits()
	if dl != 1048576 {
		t.Errorf("Expected engine download limit 1048576, got %d", dl)
	}
	if ul != 524288 {
		t.Errorf("Expected engine upload limit 524288, got %d", ul)
	}
}

func TestServerWindowShow(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_win_show_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	called := false
	srv.SetOnWindowShow(func() {
		called = true
	})

	req := httptest.NewRequest(http.MethodPost, "/api/window/show", nil)
	w := httptest.NewRecorder()
	srv.handleWindowShow(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleWindowShow returned status %d", w.Code)
	}
	if !called {
		t.Errorf("Expected OnWindowShow callback to be called")
	}
}

func TestServerFileAndRedownloadHandlers(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ui_file_handlers_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	// Missing hash parameter checks
	req := httptest.NewRequest(http.MethodGet, "/api/torrent/redownload", nil)
	w := httptest.NewRecorder()
	srv.handleRedownload(w, req)
	if !strings.Contains(w.Body.String(), "missing hash") {
		t.Errorf("Expected missing hash error, got: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/torrent/open-file", nil)
	w = httptest.NewRecorder()
	srv.handleOpenFile(w, req)
	if !strings.Contains(w.Body.String(), "missing hash or file_index") {
		t.Errorf("Expected missing parameter error, got: %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/torrent/show-file-in-folder", nil)
	w = httptest.NewRecorder()
	srv.handleShowFileInFolder(w, req)
	if !strings.Contains(w.Body.String(), "missing hash or file_index") {
		t.Errorf("Expected missing parameter error, got: %s", w.Body.String())
	}
}

func TestGetLogoAssetBytes(t *testing.T) {
	dark := GetLogoDark()
	if len(dark) == 0 {
		t.Fatalf("Expected GetLogoDark to return non-empty byte slice")
	}
	light := GetLogoLight()
	if len(light) == 0 {
		t.Fatalf("Expected GetLogoLight to return non-empty byte slice")
	}
}

func TestWindowThemeEndpoint(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "st_theme_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	eng, err := engine.NewEngine(cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer eng.Close()

	srv, err := NewServer(eng, cfgMgr)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	var receivedTheme string
	srv.SetOnThemeChange(func(th string) {
		receivedTheme = th
	})

	// 1. POST /api/window/theme with "light"
	bodyLight := []byte(`{"theme":"light"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/window/theme", bytes.NewReader(bodyLight))
	w := httptest.NewRecorder()
	srv.handleWindowTheme(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/window/theme returned status %d", w.Code)
	}
	if receivedTheme != "light" {
		t.Errorf("Expected receivedTheme to be 'light', got %q", receivedTheme)
	}

	// 2. POST /api/window/theme with "dark"
	bodyDark := []byte(`{"theme":"dark"}`)
	reqDark := httptest.NewRequest(http.MethodPost, "/api/window/theme", bytes.NewReader(bodyDark))
	wDark := httptest.NewRecorder()
	srv.handleWindowTheme(wDark, reqDark)
	if wDark.Code != http.StatusOK {
		t.Fatalf("POST /api/window/theme dark returned status %d", wDark.Code)
	}
	if receivedTheme != "dark" {
		t.Errorf("Expected receivedTheme to be 'dark', got %q", receivedTheme)
	}

	// 3. GET /api/window/theme should return 405 Method Not Allowed
	reqGet := httptest.NewRequest(http.MethodGet, "/api/window/theme", nil)
	wGet := httptest.NewRecorder()
	srv.handleWindowTheme(wGet, reqGet)
	if wGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected GET /api/window/theme to return 405, got %d", wGet.Code)
	}

	// 4. Invalid JSON should return error
	reqBad := httptest.NewRequest(http.MethodPost, "/api/window/theme", bytes.NewReader([]byte("{invalid")))
	wBad := httptest.NewRecorder()
	srv.handleWindowTheme(wBad, reqBad)
	if !strings.Contains(wBad.Body.String(), "invalid json") {
		t.Errorf("Expected invalid json error, got: %s", wBad.Body.String())
	}
}
