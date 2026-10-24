//go:build windows

package sysintegration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	modShell32         = windows.NewLazySystemDLL("shell32.dll")
	procSHChangeNotify = modShell32.NewProc("SHChangeNotify")
)

const (
	shcneAssocChanged = 0x08000000
	shcnfIdList       = 0x0000
)

var (
	pathMu        sync.RWMutex
	classesSubkey = `Software\Classes`
	runSubkey     = `Software\Microsoft\Windows\CurrentVersion\Run`
)

// SetTestRegistryPaths overrides registry subkeys for isolated testing.
func SetTestRegistryPaths(classes, run string) func() {
	pathMu.Lock()
	oldClasses := classesSubkey
	oldRun := runSubkey
	classesSubkey = classes
	runSubkey = run
	pathMu.Unlock()

	return func() {
		pathMu.Lock()
		classesSubkey = oldClasses
		runSubkey = oldRun
		pathMu.Unlock()
	}
}

func getClassesSubkey() string {
	pathMu.RLock()
	defer pathMu.RUnlock()
	return classesSubkey
}

func getRunSubkey() string {
	pathMu.RLock()
	defer pathMu.RUnlock()
	return runSubkey
}

func notifyAssocChanged() {
	_, _, _ = procSHChangeNotify.Call(shcneAssocChanged, shcnfIdList, 0, 0)
}

func getExePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Clean(exe), nil
}

// GetStatus returns the current integration status from the Windows registry.
func GetStatus() (IntegrationStatus, error) {
	return IntegrationStatus{
		TorrentAssociated: IsTorrentAssociated(),
		MagnetAssociated:  IsMagnetAssociated(),
		AutostartEnabled:  IsAutostartEnabled(),
	}, nil
}

// IsTorrentAssociated checks if HKCU contains .torrent registration pointing to this app.
func IsTorrentAssociated() bool {
	exe, err := getExePath()
	if err != nil {
		return false
	}

	classes := getClassesSubkey()
	dotTorrentPath := classes + `\.torrent`
	progIdPath := classes + `\SmirnovTorrent.torrent`

	k, err := registry.OpenKey(registry.CURRENT_USER, dotTorrentPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetStringValue("")
	if err != nil || val != "SmirnovTorrent.torrent" {
		return false
	}

	cmdKey, err := registry.OpenKey(registry.CURRENT_USER, progIdPath+`\shell\open\command`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer cmdKey.Close()

	cmdVal, _, err := cmdKey.GetStringValue("")
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(cmdVal), strings.ToLower(exe))
}

// SetTorrentAssociation enables or disables .torrent file registration in HKCU.
func SetTorrentAssociation(enable bool) error {
	exe, err := getExePath()
	if err != nil {
		return err
	}

	classes := getClassesSubkey()
	dotTorrentPath := classes + `\.torrent`
	progIdPath := classes + `\SmirnovTorrent.torrent`

	if !enable {
		// Only remove our ProgID from .torrent; do not obliterate key if another app was associated
		kDot, err := registry.OpenKey(registry.CURRENT_USER, dotTorrentPath, registry.ALL_ACCESS)
		if err == nil {
			val, _, _ := kDot.GetStringValue("")
			if val == "SmirnovTorrent.torrent" {
				_ = kDot.DeleteValue("")
			}
			kOpenWith, errOw := registry.OpenKey(kDot, "OpenWithProgids", registry.ALL_ACCESS)
			if errOw == nil {
				_ = kOpenWith.DeleteValue("SmirnovTorrent.torrent")
				kOpenWith.Close()
			}
			kDot.Close()

			// If .torrent has no more subkeys and values, clean it up
			kCheck, errCheck := registry.OpenKey(registry.CURRENT_USER, dotTorrentPath, registry.READ)
			if errCheck == nil {
				subkeys, _ := kCheck.ReadSubKeyNames(-1)
				values, _ := kCheck.ReadValueNames(-1)
				kCheck.Close()
				if len(subkeys) == 0 && len(values) == 0 {
					_ = registry.DeleteKey(registry.CURRENT_USER, dotTorrentPath)
				}
			}
		}

		_ = deleteKeyRecursively(registry.CURRENT_USER, progIdPath)
		notifyAssocChanged()
		return nil
	}

	// 1. HKCU\Software\Classes\.torrent
	kDot, _, err := registry.CreateKey(registry.CURRENT_USER, dotTorrentPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create .torrent registry key: %w", err)
	}
	defer kDot.Close()
	_ = kDot.SetStringValue("", "SmirnovTorrent.torrent")

	kOpenWith, _, err := registry.CreateKey(kDot, "OpenWithProgids", registry.ALL_ACCESS)
	if err == nil {
		_ = kOpenWith.SetStringValue("SmirnovTorrent.torrent", "")
		kOpenWith.Close()
	}

	// 2. HKCU\Software\Classes\SmirnovTorrent.torrent
	kProg, _, err := registry.CreateKey(registry.CURRENT_USER, progIdPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create SmirnovTorrent.torrent key: %w", err)
	}
	defer kProg.Close()
	_ = kProg.SetStringValue("", "BitTorrent File")

	kIcon, _, err := registry.CreateKey(kProg, "DefaultIcon", registry.ALL_ACCESS)
	if err == nil {
		_ = kIcon.SetStringValue("", fmt.Sprintf(`"%s",0`, exe))
		kIcon.Close()
	}

	kCmd, _, err := registry.CreateKey(kProg, `shell\open\command`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create shell command key: %w", err)
	}
	defer kCmd.Close()
	_ = kCmd.SetStringValue("", fmt.Sprintf(`"%s" "%%1"`, exe))

	notifyAssocChanged()
	return nil
}

// IsMagnetAssociated checks if HKCU contains magnet: protocol registration pointing to this app.
func IsMagnetAssociated() bool {
	exe, err := getExePath()
	if err != nil {
		return false
	}

	classes := getClassesSubkey()
	magnetCmdPath := classes + `\magnet\shell\open\command`

	k, err := registry.OpenKey(registry.CURRENT_USER, magnetCmdPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	cmdVal, _, err := k.GetStringValue("")
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(cmdVal), strings.ToLower(exe))
}

// SetMagnetAssociation enables or disables magnet: protocol registration in HKCU.
func SetMagnetAssociation(enable bool) error {
	exe, err := getExePath()
	if err != nil {
		return err
	}

	classes := getClassesSubkey()
	magnetPath := classes + `\magnet`

	if !enable {
		// Only delete magnet key if SmirnovTorrent is currently the handler
		if IsMagnetAssociated() {
			_ = deleteKeyRecursively(registry.CURRENT_USER, magnetPath)
			notifyAssocChanged()
		}
		return nil
	}

	kMag, _, err := registry.CreateKey(registry.CURRENT_USER, magnetPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create magnet key: %w", err)
	}
	defer kMag.Close()

	_ = kMag.SetStringValue("", "URL:BitTorrent Magnet URL")
	_ = kMag.SetStringValue("URL Protocol", "")

	kIcon, _, err := registry.CreateKey(kMag, "DefaultIcon", registry.ALL_ACCESS)
	if err == nil {
		_ = kIcon.SetStringValue("", fmt.Sprintf(`"%s",0`, exe))
		kIcon.Close()
	}

	kCmd, _, err := registry.CreateKey(kMag, `shell\open\command`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create magnet command key: %w", err)
	}
	defer kCmd.Close()
	_ = kCmd.SetStringValue("", fmt.Sprintf(`"%s" "%%1"`, exe))

	notifyAssocChanged()
	return nil
}

// IsAutostartEnabled checks if SmirnovTorrent is set to launch on Windows startup in HKCU.
func IsAutostartEnabled() bool {
	exe, err := getExePath()
	if err != nil {
		return false
	}

	runKey := getRunSubkey()
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetStringValue("SmirnovTorrent")
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(val), strings.ToLower(exe))
}

// SetAutostart toggles startup launch entry in HKCU\Software\Microsoft\Windows\CurrentVersion\Run.
func SetAutostart(enable bool) error {
	exe, err := getExePath()
	if err != nil {
		return err
	}

	runKey := getRunSubkey()
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open Run registry key: %w", err)
	}
	defer k.Close()

	if enable {
		val := fmt.Sprintf(`"%s" --minimized`, exe)
		return k.SetStringValue("SmirnovTorrent", val)
	}

	err = k.DeleteValue("SmirnovTorrent")
	if err == registry.ErrNotExist {
		return nil
	}
	return err
}

func deleteKeyRecursively(root registry.Key, path string) error {
	if path == "" {
		return nil
	}
	k, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}

	subkeys, err := k.ReadSubKeyNames(-1)
	k.Close()
	if err == nil {
		for _, sub := range subkeys {
			_ = deleteKeyRecursively(root, path+`\`+sub)
		}
	}
	return registry.DeleteKey(root, path)
}
