//go:build !windows

package ui

import (
	"os/exec"
	"runtime"
)

func openInNativeExplorer(cleaned string) error {
	if runtime.GOOS == "darwin" {
		return exec.Command("open", cleaned).Start()
	}
	return exec.Command("xdg-open", cleaned).Start()
}

func openNativeFile(cleaned string) error {
	if runtime.GOOS == "darwin" {
		return exec.Command("open", cleaned).Start()
	}
	return exec.Command("xdg-open", cleaned).Start()
}

func openNativeURL(urlStr string) error {
	if runtime.GOOS == "darwin" {
		return exec.Command("open", urlStr).Start()
	}
	return exec.Command("xdg-open", urlStr).Start()
}
