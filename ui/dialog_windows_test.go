//go:build windows

package ui

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestOpenURLAndOpenFileMocked(t *testing.T) {
	var openedVerb, openedTarget string
	oldShellExec := shellExecuteFn
	shellExecuteFn = func(hwnd windows.Handle, verb, file, args, cwd *uint16, showCmd int32) error {
		if verb != nil {
			openedVerb = windows.UTF16PtrToString(verb)
		}
		if file != nil {
			openedTarget = windows.UTF16PtrToString(file)
		}
		return nil
	}
	defer func() { shellExecuteFn = oldShellExec }()

	if err := OpenURL("http://127.0.0.1:55555/test"); err != nil {
		t.Fatalf("OpenURL failed: %v", err)
	}
	if openedVerb != "open" || openedTarget != "http://127.0.0.1:55555/test" {
		t.Errorf("Unexpected OpenURL call: verb=%q target=%q", openedVerb, openedTarget)
	}

	if err := OpenFile("C:\\test\\file.txt"); err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	if openedVerb != "open" || openedTarget != "C:\\test\\file.txt" {
		t.Errorf("Unexpected OpenFile call: verb=%q target=%q", openedVerb, openedTarget)
	}
}
