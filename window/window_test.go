package window

import (
	"os"
	"testing"
	"unsafe"

	"github.com/TheTitanrain/w32"
	"golang.org/x/sys/windows"

	"smirnovtorrent/config"
)

func TestDesktopWindowConfigBinding(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "win_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.MinimizeToTray = true
	cfg.CloseToTray = true
	_ = cfgMgr.Update(cfg)

	dw := &DesktopWindow{
		cfgManager: cfgMgr,
	}

	// Test safe calls on mock/headless window
	dw.initPlatform("Test", 1480, 850)
	dw.Show()
	dw.Hide()
	dw.cleanupPlatform()
	dw.Destroy()

	if !dw.cfgManager.Get().MinimizeToTray {
		t.Errorf("Expected MinimizeToTray to remain true")
	}
}

func TestCustomWndProcInterceptions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "win_wndproc_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	cfg := cfgMgr.Get()
	cfg.MinimizeToTray = true
	cfg.CloseToTray = true
	_ = cfgMgr.Update(cfg)

	dw := &DesktopWindow{
		cfgManager: cfgMgr,
	}

	windowMu.Lock()
	globalWindow = dw
	forceQuit = false
	windowMu.Unlock()

	defer func() {
		windowMu.Lock()
		globalWindow = nil
		forceQuit = false
		windowMu.Unlock()
	}()

	// 1. SC_MINIMIZE should be intercepted and return 0 when MinimizeToTray is true
	ret := customWndProc(0, WM_SYSCOMMAND, SC_MINIMIZE, 0)
	if ret != 0 {
		t.Errorf("Expected customWndProc to return 0 for SC_MINIMIZE when MinimizeToTray is true, got %v", ret)
	}

	// 2. WM_CLOSE should be intercepted and return 0 when CloseToTray is true
	retClose := customWndProc(0, WM_CLOSE, 0, 0)
	if retClose != 0 {
		t.Errorf("Expected customWndProc to return 0 for WM_CLOSE when CloseToTray is true, got %v", retClose)
	}

	// 3. WM_TRAY_CALLBACK with WM_LBUTTONUP should return 0
	retTrayClick := customWndProc(0, WM_TRAY_CALLBACK, 0, WM_LBUTTONUP)
	if retTrayClick != 0 {
		t.Errorf("Expected customWndProc to return 0 for WM_TRAY_CALLBACK WM_LBUTTONUP, got %v", retTrayClick)
	}

	// 4. WM_TRAY_CALLBACK with WM_LBUTTONDBLCLK should return 0
	retTrayDbl := customWndProc(0, WM_TRAY_CALLBACK, 0, WM_LBUTTONDBLCLK)
	if retTrayDbl != 0 {
		t.Errorf("Expected customWndProc to return 0 for WM_TRAY_CALLBACK WM_LBUTTONDBLCLK, got %v", retTrayDbl)
	}

	// 5. WM_GETICON should return transparent icon handle to suppress generic icon fallback
	retGetIcon := customWndProc(0, WM_GETICON, 0, 0)
	hEmpty := getEmptyIcon()
	if retGetIcon != hEmpty || retGetIcon == 0 {
		t.Errorf("Expected customWndProc to return transparent icon handle %v for WM_GETICON, got %v", hEmpty, retGetIcon)
	}

	// 6. WM_SETICON with non-matching lParam should return 0 to block setting default/generic icons
	retSetIcon := customWndProc(0, WM_SETICON, 0, 0x1234)
	if retSetIcon != 0 {
		t.Errorf("Expected customWndProc to return 0 for WM_SETICON with non-matching lParam, got %v", retSetIcon)
	}
}

func TestSetWindowDarkMode(t *testing.T) {
	// Must safely handle null HWND
	setWindowDarkMode(0)

	// Must safely handle non-zero test HWND without panic
	setWindowDarkMode(0x12345678)

	// Test on real HWND
	hwnd := w32.CreateWindowEx(
		0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr("Test"),
		w32.WS_OVERLAPPEDWINDOW,
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 300,
		0, 0, 0, nil,
	)
	if hwnd != 0 {
		defer w32.DestroyWindow(hwnd)

		darkMode := int32(1)
		r1, _, _ := procDwmSetWindowAttribute.Call(
			uintptr(hwnd),
			uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE),
			uintptr(unsafe.Pointer(&darkMode)),
			unsafe.Sizeof(darkMode),
		)
		t.Logf("DWMWA_USE_IMMERSIVE_DARK_MODE ret: 0x%x", r1)

		captionColor := uint32(0x001d1817)
		r2, _, _ := procDwmSetWindowAttribute.Call(
			uintptr(hwnd),
			uintptr(DWMWA_CAPTION_COLOR),
			uintptr(unsafe.Pointer(&captionColor)),
			unsafe.Sizeof(captionColor),
		)
		t.Logf("DWMWA_CAPTION_COLOR ret: 0x%x", r2)

		textColor := uint32(0x00f7f3f1)
		r3, _, _ := procDwmSetWindowAttribute.Call(
			uintptr(hwnd),
			uintptr(DWMWA_TEXT_COLOR),
			uintptr(unsafe.Pointer(&textColor)),
			unsafe.Sizeof(textColor),
		)
		t.Logf("DWMWA_TEXT_COLOR ret: 0x%x", r3)

		setWindowDarkMode(uintptr(hwnd))
	}
}

func TestCenterWindow(t *testing.T) {
	// Must safely handle null HWND
	centerWindow(0, 1480, 850)

	// Test on real HWND
	hwnd := w32.CreateWindowEx(
		0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr("TestCenter"),
		w32.WS_OVERLAPPEDWINDOW,
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 300,
		0, 0, 0, nil,
	)
	if hwnd != 0 {
		defer w32.DestroyWindow(hwnd)
		centerWindow(uintptr(hwnd), 1480, 850)

		r := w32.GetWindowRect(hwnd)
		if r != nil {
			width := r.Right - r.Left
			height := r.Bottom - r.Top
			if width <= 0 || height <= 0 {
				t.Errorf("Expected positive window dimensions, got %dx%d", width, height)
			}
		}
	}
}

func TestDlgModalFrameStyle(t *testing.T) {
	hwnd := w32.CreateWindowEx(
		0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr("TestStyles"),
		w32.WS_OVERLAPPEDWINDOW,
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 300,
		0, 0, 0, nil,
	)
	if hwnd == 0 {
		t.Skip("Cannot create test window")
	}
	defer w32.DestroyWindow(hwnd)

	styleBefore := w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	exBefore := w32.GetWindowLong(hwnd, w32.GWL_EXSTYLE)

	// Apply WS_EX_DLGMODALFRAME
	w32.SetWindowLong(hwnd, w32.GWL_EXSTYLE, uint32(exBefore)|0x0001)

	styleAfter := w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	exAfter := w32.GetWindowLong(hwnd, w32.GWL_EXSTYLE)

	t.Logf("Style before: 0x%x, after: 0x%x", styleBefore, styleAfter)
	t.Logf("ExStyle before: 0x%x, after: 0x%x", exBefore, exAfter)

	if uint32(styleAfter)&uint32(w32.WS_MINIMIZEBOX) == 0 {
		t.Errorf("WS_MINIMIZEBOX was lost")
	}
	if uint32(styleAfter)&uint32(w32.WS_MAXIMIZEBOX) == 0 {
		t.Errorf("WS_MAXIMIZEBOX was lost")
	}
}

func TestRemoveWindowIcon(t *testing.T) {
	// Must safely handle zero hwnd
	dwNil := &DesktopWindow{}
	dwNil.removeWindowIcon()

	// Test on real HWND
	hwnd := w32.CreateWindowEx(
		0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr("TestRemoveIcon"),
		w32.WS_OVERLAPPEDWINDOW,
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 300,
		0, 0, 0, nil,
	)
	if hwnd == 0 {
		t.Skip("Cannot create test window")
	}
	defer w32.DestroyWindow(hwnd)

	dw := &DesktopWindow{
		hwnd: uintptr(hwnd),
	}
	dw.removeWindowIcon()

	exStyle := w32.GetWindowLong(hwnd, w32.GWL_EXSTYLE)
	if uint32(exStyle)&WS_EX_DLGMODALFRAME == 0 {
		t.Errorf("Expected WS_EX_DLGMODALFRAME to be set in extended style")
	}
	if uint32(exStyle)&WS_EX_APPWINDOW == 0 {
		t.Errorf("Expected WS_EX_APPWINDOW to be set in extended style")
	}

	hEmpty := getEmptyIcon()
	if hEmpty == 0 {
		t.Errorf("Expected getEmptyIcon() to produce a valid non-zero icon handle")
	}
}

func TestThemeIconLoadingAndSwitching(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "win_icon_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	t.Setenv("APPDATA", tempDir)

	cfgMgr := config.NewConfigManager()
	dw := &DesktopWindow{
		cfgManager: cfgMgr,
	}

	dw.loadIcons()
	defer dw.cleanupPlatform()

	if dw.iconDarkSmall == 0 {
		t.Errorf("Expected iconDarkSmall to be non-zero")
	}
	if dw.iconDarkBig == 0 {
		t.Errorf("Expected iconDarkBig to be non-zero")
	}
	if dw.iconLightSmall == 0 {
		t.Errorf("Expected iconLightSmall to be non-zero")
	}
	if dw.iconLightBig == 0 {
		t.Errorf("Expected iconLightBig to be non-zero")
	}
	if dw.iconDarkSmall == dw.iconLightSmall {
		t.Errorf("Expected dark and light small icons to have distinct handles")
	}

	hwnd := w32.CreateWindowEx(
		0,
		windows.StringToUTF16Ptr("STATIC"),
		windows.StringToUTF16Ptr("TestIconSwitch"),
		w32.WS_OVERLAPPEDWINDOW,
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 300,
		0, 0, 0, nil,
	)
	if hwnd == 0 {
		t.Skip("Cannot create test window")
	}
	defer w32.DestroyWindow(hwnd)

	dw.hwnd = uintptr(hwnd)

	// Apply dark theme
	dw.applyWindowIcon("dark")
	if dw.currentIconSmall != dw.iconDarkSmall {
		t.Errorf("Expected currentIconSmall to be iconDarkSmall, got %v", dw.currentIconSmall)
	}
	if dw.currentIconBig != dw.iconDarkBig {
		t.Errorf("Expected currentIconBig to be iconDarkBig, got %v", dw.currentIconBig)
	}

	exStyle := w32.GetWindowLong(hwnd, w32.GWL_EXSTYLE)
	if uint32(exStyle)&WS_EX_DLGMODALFRAME != 0 {
		t.Errorf("Expected WS_EX_DLGMODALFRAME to NOT be set when window icon is enabled")
	}
	if uint32(exStyle)&WS_EX_APPWINDOW == 0 {
		t.Errorf("Expected WS_EX_APPWINDOW to be set for taskbar display")
	}

	// Apply light theme
	dw.applyWindowIcon("light")
	if dw.currentIconSmall != dw.iconLightSmall {
		t.Errorf("Expected currentIconSmall to be iconLightSmall, got %v", dw.currentIconSmall)
	}
	if dw.currentIconBig != dw.iconLightBig {
		t.Errorf("Expected currentIconBig to be iconLightBig, got %v", dw.currentIconBig)
	}

	// Test customWndProc returns loaded theme icons when globalWindow is set
	windowMu.Lock()
	globalWindow = dw
	windowMu.Unlock()
	defer func() {
		windowMu.Lock()
		globalWindow = nil
		windowMu.Unlock()
	}()

	retSm := customWndProc(0, WM_GETICON, uintptr(ICON_SMALL), 0)
	if retSm != dw.iconLightSmall {
		t.Errorf("Expected customWndProc WM_GETICON ICON_SMALL to return light small icon %v, got %v", dw.iconLightSmall, retSm)
	}
	retSm2 := customWndProc(0, WM_GETICON, 2, 0)
	if retSm2 != dw.iconLightSmall {
		t.Errorf("Expected customWndProc WM_GETICON ICON_SMALL2 (2) to return light small icon %v, got %v", dw.iconLightSmall, retSm2)
	}
	retBg := customWndProc(0, WM_GETICON, uintptr(ICON_BIG), 0)
	if retBg != dw.iconLightBig {
		t.Errorf("Expected customWndProc WM_GETICON ICON_BIG to return light big icon %v, got %v", dw.iconLightBig, retBg)
	}

	// Verify custom theme icons are protected from external overrides
	if retOverride := customWndProc(0, WM_SETICON, 0, 0x9999); retOverride != 0 {
		t.Errorf("Expected customWndProc WM_SETICON with foreign handle to return 0, got %v", retOverride)
	}

	// Switch back to dark theme
	dw.applyWindowIcon("dark")
	retSmDark := customWndProc(0, WM_GETICON, uintptr(ICON_SMALL), 0)
	if retSmDark != dw.iconDarkSmall {
		t.Errorf("Expected customWndProc WM_GETICON ICON_SMALL to return dark small icon %v, got %v", dw.iconDarkSmall, retSmDark)
	}
	retSm2Dark := customWndProc(0, WM_GETICON, 2, 0)
	if retSm2Dark != dw.iconDarkSmall {
		t.Errorf("Expected customWndProc WM_GETICON ICON_SMALL2 (2) to return dark small icon %v, got %v", dw.iconDarkSmall, retSm2Dark)
	}
	retBgDark := customWndProc(0, WM_GETICON, uintptr(ICON_BIG), 0)
	if retBgDark != dw.iconDarkBig {
		t.Errorf("Expected customWndProc WM_GETICON ICON_BIG to return dark big icon %v, got %v", dw.iconDarkBig, retBgDark)
	}
}
