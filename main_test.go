package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestIconFiles(t *testing.T) {
	icoPath := "smirnovtorrent.ico"
	data, err := os.ReadFile(icoPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", icoPath, err)
	}
	if len(data) < 6 {
		t.Fatalf("%s is too small to be a valid ICO file", icoPath)
	}

	// ICO header: reserved (0), type (1 for ICO), count
	reserved := binary.LittleEndian.Uint16(data[0:2])
	icoType := binary.LittleEndian.Uint16(data[2:4])
	count := binary.LittleEndian.Uint16(data[4:6])

	if reserved != 0 || icoType != 1 {
		t.Errorf("Invalid ICO header in %s: reserved=%d, type=%d", icoPath, reserved, icoType)
	}
	if count < 4 {
		t.Errorf("Expected at least 4 icon sizes in %s, got %d", icoPath, count)
	}

	// Verify syso files exist for Windows builds
	sysoAmd64 := "rsrc_windows_amd64.syso"
	amd64Data, err := os.ReadFile(sysoAmd64)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", sysoAmd64, err)
	}
	if len(amd64Data) == 0 {
		t.Errorf("Expected %s to have non-zero size", sysoAmd64)
	}

	syso386 := "rsrc_windows_386.syso"
	if _, err := os.Stat(syso386); err != nil {
		t.Errorf("Expected %s to exist: %v", syso386, err)
	}

	sysoArm64 := "rsrc_windows_arm64.syso"
	if _, err := os.Stat(sysoArm64); err != nil {
		t.Errorf("Expected %s to exist: %v", sysoArm64, err)
	}
}

func TestExecutableIconEmbedded(t *testing.T) {
	exePath := "SmirnovTorrent.exe"
	data, err := os.ReadFile(exePath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("SmirnovTorrent.exe not found; skipping binary resource test")
		}
		t.Fatalf("Failed to read %s: %v", exePath, err)
	}

	// Verify MZ header
	if len(data) < 64 || !bytes.Equal(data[:2], []byte("MZ")) {
		t.Fatalf("%s is not a valid PE binary", exePath)
	}

	peOffset := binary.LittleEndian.Uint32(data[0x3c:0x40])
	if int(peOffset)+24 > len(data) {
		t.Fatalf("Invalid PE header offset: %d", peOffset)
	}

	if !bytes.Equal(data[peOffset:peOffset+4], []byte("PE\x00\x00")) {
		t.Fatalf("Missing PE signature at offset %d", peOffset)
	}

	numSections := binary.LittleEndian.Uint16(data[peOffset+6 : peOffset+8])
	optHdrSize := binary.LittleEndian.Uint16(data[peOffset+20 : peOffset+22])
	sectOffset := int(peOffset) + 24 + int(optHdrSize)

	hasRsrc := false
	for i := 0; i < int(numSections); i++ {
		secStart := sectOffset + i*40
		if secStart+40 > len(data) {
			break
		}
		name := string(bytes.TrimRight(data[secStart:secStart+8], "\x00"))
		if name == ".rsrc" {
			hasRsrc = true
			break
		}
	}

	if !hasRsrc {
		t.Errorf("Expected %s to contain .rsrc section with embedded icon resources", exePath)
	}
}

func TestTryForwardToRunningInstanceNoPortFile(t *testing.T) {
	nonExistentPortFile := filepath.Join(t.TempDir(), "nonexistent.port")
	if tryForwardToRunningInstance(nonExistentPortFile) {
		t.Errorf("Expected false when port file does not exist")
	}
}
