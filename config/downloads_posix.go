//go:build !windows

package config

import (
	"os"
	"path/filepath"
)

func getDefaultDownloadDirPlatform() string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		downloads := filepath.Join(home, "Downloads")
		_ = os.MkdirAll(downloads, 0755)
		if fi, err := os.Stat(downloads); err == nil && fi.IsDir() {
			return downloads
		}
		return home
	}
	return "."
}
