package sysintegration

import (
	"runtime"
	"testing"

	"golang.org/x/sys/windows/registry"
	"smirnovtorrent/config"
)

func TestIntegrationToggles(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows registry integration test on non-Windows")
	}

	testClasses := `Software\Classes\SmirnovTorrentTest`
	testRun := `Software\Classes\SmirnovTorrentTestRun`
	cleanupPaths := SetTestRegistryPaths(testClasses, testRun)
	defer func() {
		cleanupPaths()
		_ = deleteKeyRecursively(registry.CURRENT_USER, testClasses)
		_ = deleteKeyRecursively(registry.CURRENT_USER, testRun)
	}()

	// 1. Initial status in test registry should be all false
	initial, err := GetStatus()
	if err != nil {
		t.Fatalf("Failed to get initial status: %v", err)
	}
	if initial.TorrentAssociated || initial.MagnetAssociated || initial.AutostartEnabled {
		t.Fatalf("Expected all initial test statuses to be false, got: %+v", initial)
	}

	// 2. Test Torrent Association
	if err := SetTorrentAssociation(true); err != nil {
		t.Fatalf("SetTorrentAssociation(true) failed: %v", err)
	}
	if !IsTorrentAssociated() {
		t.Errorf("Expected IsTorrentAssociated to be true after enabling")
	}

	if err := SetTorrentAssociation(false); err != nil {
		t.Fatalf("SetTorrentAssociation(false) failed: %v", err)
	}
	if IsTorrentAssociated() {
		t.Errorf("Expected IsTorrentAssociated to be false after disabling")
	}

	// 3. Test Magnet Association
	if err := SetMagnetAssociation(true); err != nil {
		t.Fatalf("SetMagnetAssociation(true) failed: %v", err)
	}
	if !IsMagnetAssociated() {
		t.Errorf("Expected IsMagnetAssociated to be true after enabling")
	}

	if err := SetMagnetAssociation(false); err != nil {
		t.Fatalf("SetMagnetAssociation(false) failed: %v", err)
	}
	if IsMagnetAssociated() {
		t.Errorf("Expected IsMagnetAssociated to be false after disabling")
	}

	// 4. Test Autostart
	if err := SetAutostart(true); err != nil {
		t.Fatalf("SetAutostart(true) failed: %v", err)
	}
	if !IsAutostartEnabled() {
		t.Errorf("Expected IsAutostartEnabled to be true after enabling")
	}

	if err := SetAutostart(false); err != nil {
		t.Fatalf("SetAutostart(false) failed: %v", err)
	}
	if IsAutostartEnabled() {
		t.Errorf("Expected IsAutostartEnabled to be false after disabling")
	}

	// 5. Test ApplyIntegration
	cfg := config.AppConfig{
		AssociateTorrents: true,
		AssociateMagnets:  true,
		Autostart:         true,
	}
	if err := ApplyIntegration(cfg); err != nil {
		t.Fatalf("ApplyIntegration(all true) failed: %v", err)
	}

	st, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if !st.TorrentAssociated || !st.MagnetAssociated || !st.AutostartEnabled {
		t.Errorf("Expected all to be enabled after ApplyIntegration: %+v", st)
	}

	cfgDisabled := config.AppConfig{
		AssociateTorrents: false,
		AssociateMagnets:  false,
		Autostart:         false,
	}
	if err := ApplyIntegration(cfgDisabled); err != nil {
		t.Fatalf("ApplyIntegration(all false) failed: %v", err)
	}
	st2, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if st2.TorrentAssociated || st2.MagnetAssociated || st2.AutostartEnabled {
		t.Errorf("Expected all to be disabled after ApplyIntegration(false): %+v", st2)
	}
}
