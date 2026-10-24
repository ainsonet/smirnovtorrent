//go:build windows

package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func openInNativeExplorer(cleaned string) error {
	safePath := strings.ReplaceAll(cleaned, `"`, "")
	fi, err := os.Stat(safePath)
	cmd := exec.Command("explorer.exe")
	targetArg := safePath
	if strings.HasSuffix(targetArg, `\`) {
		targetArg += `\`
	}
	if err == nil && fi.IsDir() {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf(`explorer.exe "%s"`, targetArg),
		}
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf(`explorer.exe /select,"%s"`, targetArg),
		}
	}
	return cmd.Start()
}

var shellExecuteFn = windows.ShellExecute

func openNativeFile(cleaned string) error {
	safePath := strings.ReplaceAll(cleaned, `"`, "")
	verb, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	file, err := windows.UTF16PtrFromString(safePath)
	if err != nil {
		return err
	}
	return shellExecuteFn(0, verb, file, nil, nil, windows.SW_SHOWNORMAL)
}

func openNativeURL(urlStr string) error {
	safeURL := strings.ReplaceAll(urlStr, `"`, "")
	verb, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	target, err := windows.UTF16PtrFromString(safeURL)
	if err != nil {
		return err
	}
	return shellExecuteFn(0, verb, target, nil, nil, windows.SW_SHOWNORMAL)
}
