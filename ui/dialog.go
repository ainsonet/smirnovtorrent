package ui

import (
	"path/filepath"

	"github.com/sqweek/dialog"
)

// PickTorrentFile opens a native Windows file dialog to choose a .torrent file
func PickTorrentFile() (string, error) {
	filename, err := dialog.File().
		Filter("Torrent files (*.torrent)", "torrent").
		Filter("All files (*.*)", "*").
		Title("Select .torrent file").
		Load()
	if err != nil {
		return "", err
	}
	return filename, nil
}

// PickFolder opens a native folder selection dialog
func PickFolder(defaultPath string) (string, error) {
	d := dialog.Directory().Title("Select Download Directory")
	if defaultPath != "" {
		d = d.SetStartDir(defaultPath)
	}
	return d.Browse()
}

// OpenInExplorer opens Windows Explorer at the given file or directory
func OpenInExplorer(path string) error {
	cleaned := filepath.Clean(path)
	return openInNativeExplorer(cleaned)
}

// OpenFile opens a file using the default system handler
func OpenFile(path string) error {
	cleaned := filepath.Clean(path)
	return openNativeFile(cleaned)
}

// OpenURL opens an HTTP/HTTPS URL in the default web browser without filepath mangling
func OpenURL(urlStr string) error {
	return openNativeURL(urlStr)
}
