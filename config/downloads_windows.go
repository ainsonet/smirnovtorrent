//go:build windows

package config

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func getDefaultDownloadDirPlatform() string {
	// 1. Try Windows Registry (handles localized/redirected Downloads folder like "Загрузки")
	if regPath := getWindowsDownloadsFromRegistry(); regPath != "" {
		_ = os.MkdirAll(regPath, 0755)
		if fi, err := os.Stat(regPath); err == nil && fi.IsDir() {
			return regPath
		}
	}

	// 2. Try standard USERPROFILE paths
	userProfile := os.Getenv("USERPROFILE")
	if userProfile != "" {
		downloads := filepath.Join(userProfile, "Downloads")
		_ = os.MkdirAll(downloads, 0755)
		if fi, err := os.Stat(downloads); err == nil && fi.IsDir() {
			return downloads
		}

		// Russian localized directory check
		zagruzki := filepath.Join(userProfile, "Загрузки")
		if fi, err := os.Stat(zagruzki); err == nil && fi.IsDir() {
			return zagruzki
		}
	}

	// 3. UserHomeDir fallback
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		downloads := filepath.Join(home, "Downloads")
		_ = os.MkdirAll(downloads, 0755)
		if fi, err := os.Stat(downloads); err == nil && fi.IsDir() {
			return downloads
		}
		return home
	}

	return "."
}

func getWindowsDownloadsFromRegistry() string {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()

	// Known folder GUID for FOLDERID_Downloads: {374DE290-123F-4565-9164-39C4925E467B}
	val, _, err := k.GetStringValue("{374DE290-123F-4565-9164-39C4925E467B}")
	if err == nil && val != "" {
		expanded, err := registry.ExpandString(val)
		if err == nil && expanded != "" {
			return expanded
		}
		return val
	}

	// Alternate known folder GUID: {7d83ee9b-2244-4e70-b1f5-4e3982e44d62}
	val, _, err = k.GetStringValue("{7d83ee9b-2244-4e70-b1f5-4e3982e44d62}")
	if err == nil && val != "" {
		expanded, err := registry.ExpandString(val)
		if err == nil && expanded != "" {
			return expanded
		}
		return val
	}

	return ""
}
