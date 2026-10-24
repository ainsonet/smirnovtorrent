//go:build windows

package window

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"sync"
	"unsafe"

	"github.com/TheTitanrain/w32"
	"golang.org/x/sys/windows"

	"smirnovtorrent/ui"
)

var (
	modUser32                    = windows.NewLazySystemDLL("user32.dll")
	procCreatePopupMenu          = modUser32.NewProc("CreatePopupMenu")
	procAppendMenuW              = modUser32.NewProc("AppendMenuW")
	procTrackPopupMenu           = modUser32.NewProc("TrackPopupMenu")
	procDestroyMenu              = modUser32.NewProc("DestroyMenu")
	procGetCursorPos             = modUser32.NewProc("GetCursorPos")
	procSetForegroundWindow      = modUser32.NewProc("SetForegroundWindow")
	procPostMessageW             = modUser32.NewProc("PostMessageW")
	procLoadIconW                = modUser32.NewProc("LoadIconW")
	procSetWindowPos             = modUser32.NewProc("SetWindowPos")
	procSendMessageW             = modUser32.NewProc("SendMessageW")
	procSetClassLongPtrW         = modUser32.NewProc("SetClassLongPtrW")
	procSetClassLongW            = modUser32.NewProc("SetClassLongW")
	procCreateIcon               = modUser32.NewProc("CreateIcon")
	procDestroyIcon              = modUser32.NewProc("DestroyIcon")
	procCreateIconFromResourceEx = modUser32.NewProc("CreateIconFromResourceEx")
	procGetSystemMetrics         = modUser32.NewProc("GetSystemMetrics")

	modShell32                                  = windows.NewLazySystemDLL("shell32.dll")
	procShell_NotifyIconW                       = modShell32.NewProc("Shell_NotifyIconW")
	procSetCurrentProcessExplicitAppUserModelID = modShell32.NewProc("SetCurrentProcessExplicitAppUserModelID")

	modDwmapi                 = windows.NewLazySystemDLL("dwmapi.dll")
	procDwmSetWindowAttribute = modDwmapi.NewProc("DwmSetWindowAttribute")
)

const (
	WM_CLOSE         = 0x0010
	WM_DESTROY       = 0x0002
	WM_COMMAND       = 0x0111
	WM_SYSCOMMAND    = 0x0112
	SC_MINIMIZE      = 0xF020
	WM_LBUTTONUP     = 0x0202
	WM_LBUTTONDBLCLK = 0x0203
	WM_RBUTTONUP     = 0x0205

	WM_TRAY_CALLBACK = 0x8001 // WM_APP + 1

	NIM_ADD    = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002

	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004

	MF_STRING    = 0x00000000
	MF_SEPARATOR = 0x00000800

	TPM_RIGHTBUTTON = 0x0002
	TPM_RETURNCMD   = 0x0100

	ID_TRAY_SHOW = 1001
	ID_TRAY_EXIT = 1002

	DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1 = 19
	DWMWA_USE_IMMERSIVE_DARK_MODE             = 20
	DWMWA_BORDER_COLOR                        = 34
	DWMWA_CAPTION_COLOR                       = 35
	DWMWA_TEXT_COLOR                          = 36

	SWP_NOSIZE       = 0x0001
	SWP_NOMOVE       = 0x0002
	SWP_NOZORDER     = 0x0004
	SWP_FRAMECHANGED = 0x0020

	WM_SETICON = 0x0080
	WM_GETICON = 0x007F

	ICON_SMALL = 0
	ICON_BIG   = 1

	WS_EX_DLGMODALFRAME = 0x00000001
	WS_EX_APPWINDOW     = 0x00040000

	GCLP_HICONSM = ^uintptr(33) // -34
	GCLP_HICON   = ^uintptr(13) // -14

	SM_CXICON   = 11
	SM_CYICON   = 12
	SM_CXSMICON = 49
	SM_CYSMICON = 50
)

type NOTIFYICONDATAW struct {
	cbSize            uint32
	hWnd              uintptr
	uID               uint32
	uFlags            uint32
	uCallbackMessage  uint32
	hIcon             uintptr
	szTip             [128]uint16
	dwState           uint32
	dwStateMask       uint32
	szInfo            [256]uint16
	uTimeoutOrVersion uint32
	szInfoTitle       [64]uint16
	dwInfoFlags       uint32
	guidItem          windows.GUID
	hBalloonIcon      uintptr
}

type POINT struct {
	X int32
	Y int32
}

var (
	windowMu        sync.Mutex
	globalWindow    *DesktopWindow
	originalWndProc uintptr
	forceQuit       bool
	emptyIconMu     sync.Mutex
	emptyIcon       uintptr
)

func getEmptyIcon() uintptr {
	emptyIconMu.Lock()
	defer emptyIconMu.Unlock()
	if emptyIcon != 0 {
		return emptyIcon
	}
	if procCreateIcon.Find() != nil {
		return 0
	}
	// Standard 16x16 1-bpp icon: 16 bits (2 bytes) per scanline * 16 rows = 32 bytes.
	// AND mask: all 1s (0xFF) -> preserve background (100% transparent).
	// XOR mask: all 0s (0x00) -> leave pixels unchanged.
	andMask := make([]byte, 32)
	for i := range andMask {
		andMask[i] = 0xFF
	}
	xorMask := make([]byte, 32)
	hIcon, _, _ := procCreateIcon.Call(
		0,
		16, 16,
		1, 1,
		uintptr(unsafe.Pointer(&andMask[0])),
		uintptr(unsafe.Pointer(&xorMask[0])),
	)
	emptyIcon = hIcon
	return emptyIcon
}

func customWndProc(hwnd, msg, wp, lp uintptr) uintptr {
	switch msg {
	case WM_SYSCOMMAND:
		if (wp & 0xFFF0) == SC_MINIMIZE {
			windowMu.Lock()
			gw := globalWindow
			windowMu.Unlock()
			if gw != nil && gw.cfgManager != nil && gw.cfgManager.Get().MinimizeToTray {
				gw.Hide()
				return 0
			}
		}
	case WM_CLOSE:
		windowMu.Lock()
		gw := globalWindow
		fq := forceQuit
		windowMu.Unlock()
		if !fq && gw != nil && gw.cfgManager != nil && gw.cfgManager.Get().CloseToTray {
			gw.Hide()
			return 0
		}
	case WM_TRAY_CALLBACK:
		switch lp {
		case WM_LBUTTONUP, WM_LBUTTONDBLCLK:
			windowMu.Lock()
			gw := globalWindow
			windowMu.Unlock()
			if gw != nil {
				gw.Show()
			}
			return 0
		case WM_RBUTTONUP:
			windowMu.Lock()
			gw := globalWindow
			windowMu.Unlock()
			if gw != nil {
				gw.showTrayMenu()
			}
			return 0
		}
	case WM_GETICON:
		windowMu.Lock()
		gw := globalWindow
		windowMu.Unlock()
		if gw != nil {
			if (wp == uintptr(ICON_SMALL) || wp == 2) && gw.currentIconSmall != 0 {
				return gw.currentIconSmall
			}
			if wp == uintptr(ICON_BIG) && gw.currentIconBig != 0 {
				return gw.currentIconBig
			}
			if gw.currentIconSmall != 0 {
				return gw.currentIconSmall
			}
			if gw.currentIconBig != 0 {
				return gw.currentIconBig
			}
		}
		// Return transparent icon handle to prevent Windows/DWM from falling back to IDI_APPLICATION
		if hEmpty := getEmptyIcon(); hEmpty != 0 {
			return hEmpty
		}
		return 0
	case WM_SETICON:
		windowMu.Lock()
		gw := globalWindow
		windowMu.Unlock()
		if gw != nil && (lp == gw.currentIconSmall || lp == gw.currentIconBig) && lp != 0 {
			break
		}
		if gw != nil && (gw.currentIconSmall != 0 || gw.currentIconBig != 0) {
			// Suppress external attempts to overwrite custom theme icon
			return 0
		}
		hEmpty := getEmptyIcon()
		if hEmpty != 0 && lp != hEmpty {
			// Suppress attempts to set default or non-transparent icon in caption
			return 0
		}
	}

	if originalWndProc != 0 {
		return w32.CallWindowProc(originalWndProc, w32.HWND(hwnd), uint32(msg), wp, lp)
	}
	return w32.DefWindowProc(w32.HWND(hwnd), uint32(msg), wp, lp)
}

func (dw *DesktopWindow) setThemePlatform(theme string) {
	if dw == nil || dw.hwnd == 0 {
		return
	}
	isDark := theme != "light"
	setWindowTheme(dw.hwnd, isDark)
	dw.applyWindowIcon(theme)
}

func (dw *DesktopWindow) applyConfigTheme() {
	if dw == nil || dw.hwnd == 0 {
		return
	}
	theme := "dark"
	if dw.cfgManager != nil && dw.cfgManager.Get().Theme == "light" {
		theme = "light"
	}
	setWindowTheme(dw.hwnd, theme != "light")
	dw.applyWindowIcon(theme)
}

func setWindowTheme(hwnd uintptr, isDark bool) {
	if hwnd == 0 {
		return
	}
	if procDwmSetWindowAttribute.Find() != nil {
		return
	}

	darkMode := int32(0)
	if isDark {
		darkMode = 1
	}

	_, _, _ = procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)
	_, _, _ = procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)

	var captionColor, textColor, borderColor uint32
	if isDark {
		captionColor = 0x001d1817 // #17181d
		textColor = 0x00f7f3f1    // #f1f3f7
		borderColor = 0x001d1817  // #17181d
	} else {
		captionColor = 0x00f9f5f1 // #f1f5f9
		textColor = 0x002a170f    // #0f172a
		borderColor = 0x00f0e8e2  // #e2e8f0 in 0x00BBGGRR format
	}

	_, _, _ = procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(DWMWA_CAPTION_COLOR),
		uintptr(unsafe.Pointer(&captionColor)),
		unsafe.Sizeof(captionColor),
	)
	_, _, _ = procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(DWMWA_TEXT_COLOR),
		uintptr(unsafe.Pointer(&textColor)),
		unsafe.Sizeof(textColor),
	)
	_, _, _ = procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(DWMWA_BORDER_COLOR),
		uintptr(unsafe.Pointer(&borderColor)),
		unsafe.Sizeof(borderColor),
	)

	if procSetWindowPos.Find() == nil {
		_, _, _ = procSetWindowPos.Call(
			hwnd,
			0,
			0, 0, 0, 0,
			SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED,
		)
	}
}

// setWindowDarkMode configures the native Windows window frame/titlebar to a dark theme
// matching the SmirnovTorrent interface color (#17181d).
func setWindowDarkMode(hwnd uintptr) {
	setWindowTheme(hwnd, true)
}

// centerWindow positions and sizes the window centered on the monitor where it resides,
// clamping to the available monitor work area so it never overflows off-screen.
func centerWindow(hwnd uintptr, width, height int) {
	if hwnd == 0 {
		return
	}
	hMonitor := w32.MonitorFromWindow(w32.HWND(hwnd), w32.MONITOR_DEFAULTTONEAREST)
	var mi w32.MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	if hMonitor != 0 && w32.GetMonitorInfo(hMonitor, &mi) {
		workW := int(mi.RcWork.Right - mi.RcWork.Left)
		workH := int(mi.RcWork.Bottom - mi.RcWork.Top)
		if width > workW {
			width = workW
		}
		if height > workH {
			height = workH
		}
		x := int(mi.RcWork.Left) + (workW-width)/2
		y := int(mi.RcWork.Top) + (workH-height)/2
		w32.SetWindowPos(w32.HWND(hwnd), 0, x, y, width, height, w32.SWP_NOZORDER|w32.SWP_FRAMECHANGED)
	}
}

// removeWindowIcon removes the default generic application icon (IDI_APPLICATION)
// from the native Windows title bar / caption ("верхняя шторка").
func (dw *DesktopWindow) removeWindowIcon() {
	if dw.hwnd == 0 {
		return
	}

	// 1. Set WS_EX_DLGMODALFRAME (0x0001) which instructs Windows DWM not to display
	// an icon in the title bar (caption), and preserve/add WS_EX_APPWINDOW (0x00040000)
	// to ensure the window remains visible in the Windows taskbar.
	exStyle := w32.GetWindowLong(w32.HWND(dw.hwnd), w32.GWL_EXSTYLE)
	w32.SetWindowLong(w32.HWND(dw.hwnd), w32.GWL_EXSTYLE, uint32(exStyle)|WS_EX_DLGMODALFRAME|WS_EX_APPWINDOW)

	// 2. Set window icon to a 100% transparent icon so Windows DWM will not fall back
	// to the generic application icon (IDI_APPLICATION).
	hEmpty := getEmptyIcon()
	if hEmpty != 0 && procSendMessageW.Find() == nil {
		_, _, _ = procSendMessageW.Call(dw.hwnd, WM_SETICON, uintptr(ICON_SMALL), hEmpty)
		_, _, _ = procSendMessageW.Call(dw.hwnd, WM_SETICON, uintptr(ICON_BIG), hEmpty)
	}

	// 3. Update window class icons to transparent icon
	if hEmpty != 0 {
		if procSetClassLongPtrW.Find() == nil {
			_, _, _ = procSetClassLongPtrW.Call(dw.hwnd, GCLP_HICONSM, hEmpty)
			_, _, _ = procSetClassLongPtrW.Call(dw.hwnd, GCLP_HICON, hEmpty)
		} else if procSetClassLongW.Find() == nil {
			_, _, _ = procSetClassLongW.Call(dw.hwnd, GCLP_HICONSM, hEmpty)
			_, _, _ = procSetClassLongW.Call(dw.hwnd, GCLP_HICON, hEmpty)
		}
	}

	// 4. Force non-client frame recalculation to apply changes immediately
	if procSetWindowPos.Find() == nil {
		_, _, _ = procSetWindowPos.Call(
			dw.hwnd,
			0,
			0, 0, 0, 0,
			SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED,
		)
	}
}

func init() {
	setAppUserModelID()
}

func setAppUserModelIDPlatform() {
	setAppUserModelID()
}

func setAppUserModelID() {
	if procSetCurrentProcessExplicitAppUserModelID.Find() == nil {
		appID, _ := windows.UTF16PtrFromString("SmirnovTorrent.Client")
		_, _, _ = procSetCurrentProcessExplicitAppUserModelID.Call(uintptr(unsafe.Pointer(appID)))
	}
}

func resizeAreaAverage(src image.Image, dstW, dstH int) *image.RGBA {
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	for dy := 0; dy < dstH; dy++ {
		sy0 := bounds.Min.Y + (dy * srcH / dstH)
		sy1 := bounds.Min.Y + ((dy + 1) * srcH / dstH)
		if sy1 <= sy0 {
			sy1 = sy0 + 1
		}
		for dx := 0; dx < dstW; dx++ {
			sx0 := bounds.Min.X + (dx * srcW / dstW)
			sx1 := bounds.Min.X + ((dx + 1) * srcW / dstW)
			if sx1 <= sx0 {
				sx1 = sx0 + 1
			}

			var rSum, gSum, bSum uint64
			var count uint64
			for y := sy0; y < sy1; y++ {
				for x := sx0; x < sx1; x++ {
					r, g, b, _ := src.At(x, y).RGBA()
					rSum += uint64(r >> 8)
					gSum += uint64(g >> 8)
					bSum += uint64(b >> 8)
					count++
				}
			}
			if count > 0 {
				dst.SetRGBA(dx, dy, color.RGBA{
					R: uint8(rSum / count),
					G: uint8(gSum / count),
					B: uint8(bSum / count),
					A: 255,
				})
			}
		}
	}
	return dst
}

func applyRoundedCorners(img *image.RGBA, radius float64) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if radius <= 0 || w <= 0 || h <= 0 {
		return
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var cx, cy float64
			inCorner := false

			if float64(x) < radius && float64(y) < radius {
				cx, cy = radius, radius
				inCorner = true
			} else if float64(x) >= float64(w)-radius && float64(y) < radius {
				cx, cy = float64(w)-1-radius, radius
				inCorner = true
			} else if float64(x) < radius && float64(y) >= float64(h)-radius {
				cx, cy = radius, float64(h)-1-radius
				inCorner = true
			} else if float64(x) >= float64(w)-radius && float64(y) >= float64(h)-radius {
				cx, cy = float64(w)-1-radius, float64(h)-1-radius
				inCorner = true
			}

			if inCorner {
				dx := float64(x) - cx
				dy := float64(y) - cy
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > radius+0.5 {
					img.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
				} else if dist > radius-0.5 {
					factor := radius + 0.5 - dist
					c := img.RGBAAt(x, y)
					newA := uint8(float64(c.A) * factor)
					img.SetRGBA(x, y, color.RGBA{R: c.R, G: c.G, B: c.B, A: newA})
				}
			}
		}
	}
}

func createIconFromRGBA(img *image.RGBA) uintptr {
	if img == nil || procCreateIconFromResourceEx.Find() != nil {
		return 0
	}
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	if w <= 0 || h <= 0 {
		return 0
	}

	andRowBytes := ((w + 31) / 32) * 4
	andMaskSize := andRowBytes * h
	xorSize := w * h * 4
	totalSize := 40 + xorSize + andMaskSize

	buf := make([]byte, totalSize)

	// BITMAPINFOHEADER (40 bytes)
	binary.LittleEndian.PutUint32(buf[0:4], 40)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(w))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(h*2))
	binary.LittleEndian.PutUint16(buf[12:14], 1)
	binary.LittleEndian.PutUint16(buf[14:16], 32)
	binary.LittleEndian.PutUint32(buf[16:20], 0) // BI_RGB
	binary.LittleEndian.PutUint32(buf[20:24], uint32(xorSize))

	offset := 40
	for y := h - 1; y >= 0; y-- {
		for x := 0; x < w; x++ {
			c := img.RGBAAt(x, y)
			a := uint32(c.A)
			buf[offset] = byte((uint32(c.B)*a + 127) / 255)
			buf[offset+1] = byte((uint32(c.G)*a + 127) / 255)
			buf[offset+2] = byte((uint32(c.R)*a + 127) / 255)
			buf[offset+3] = c.A
			offset += 4

			if c.A < 128 {
				rowOffset := 40 + xorSize + (h-1-y)*andRowBytes
				byteIdx := rowOffset + x/8
				bitIdx := uint(7 - (x % 8))
				buf[byteIdx] |= (1 << bitIdx)
			}
		}
	}

	hIcon, _, _ := procCreateIconFromResourceEx.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		1,          // TRUE (Icon)
		0x00030000, // dwVersion
		uintptr(w),
		uintptr(h),
		0, // LR_DEFAULTCOLOR
	)
	return hIcon
}

func createIconFromJPEG(jpegBytes []byte, dstW, dstH int) uintptr {
	if len(jpegBytes) == 0 {
		return 0
	}
	src, err := jpeg.Decode(bytes.NewReader(jpegBytes))
	if err != nil {
		return 0
	}
	rgba := resizeAreaAverage(src, dstW, dstH)
	applyRoundedCorners(rgba, float64(dstW)*0.22)
	return createIconFromRGBA(rgba)
}

func (dw *DesktopWindow) loadIcons() {
	smW := 16
	smH := 16
	bgW := 32
	bgH := 32

	if procGetSystemMetrics.Find() == nil {
		rSmW, _, _ := procGetSystemMetrics.Call(SM_CXSMICON)
		rSmH, _, _ := procGetSystemMetrics.Call(SM_CYSMICON)
		rBgW, _, _ := procGetSystemMetrics.Call(SM_CXICON)
		rBgH, _, _ := procGetSystemMetrics.Call(SM_CYICON)
		if rSmW > 0 {
			smW = int(rSmW)
		}
		if rSmH > 0 {
			smH = int(rSmH)
		}
		if rBgW > 0 {
			bgW = int(rBgW)
		}
		if rBgH > 0 {
			bgH = int(rBgH)
		}
	}
	if bgW < 32 {
		bgW = 32
		bgH = 32
	}
	if smW < 16 {
		smW = 16
		smH = 16
	}

	darkBytes := ui.GetLogoDark()
	if len(darkBytes) == 0 {
		darkBytes, _ = os.ReadFile("smirnovtorrent.jpg")
	}

	lightBytes := ui.GetLogoLight()
	if len(lightBytes) == 0 {
		lightBytes, _ = os.ReadFile("smirnovtorrent-light.jpg")
	}

	if len(darkBytes) > 0 {
		dw.iconDarkSmall = createIconFromJPEG(darkBytes, smW, smH)
		dw.iconDarkBig = createIconFromJPEG(darkBytes, bgW, bgH)
	}

	if len(lightBytes) > 0 {
		dw.iconLightSmall = createIconFromJPEG(lightBytes, smW, smH)
		dw.iconLightBig = createIconFromJPEG(lightBytes, bgW, bgH)
	}
}

func (dw *DesktopWindow) applyWindowIcon(theme string) {
	if dw == nil || dw.hwnd == 0 {
		return
	}

	isLight := theme == "light"
	hSmall := dw.iconDarkSmall
	hBig := dw.iconDarkBig
	if isLight {
		if dw.iconLightSmall != 0 {
			hSmall = dw.iconLightSmall
		}
		if dw.iconLightBig != 0 {
			hBig = dw.iconLightBig
		}
	} else {
		if dw.iconDarkSmall != 0 {
			hSmall = dw.iconDarkSmall
		}
		if dw.iconDarkBig != 0 {
			hBig = dw.iconDarkBig
		}
	}

	windowMu.Lock()
	dw.currentIconSmall = hSmall
	dw.currentIconBig = hBig
	dw.currentTheme = theme
	windowMu.Unlock()

	// 1. Remove WS_EX_DLGMODALFRAME so caption shows icon, and ensure WS_EX_APPWINDOW is set for taskbar
	exStyle := w32.GetWindowLong(w32.HWND(dw.hwnd), w32.GWL_EXSTYLE)
	w32.SetWindowLong(w32.HWND(dw.hwnd), w32.GWL_EXSTYLE, (uint32(exStyle)&^WS_EX_DLGMODALFRAME)|WS_EX_APPWINDOW)

	// 2. Set window icons (small for caption / titlebar, big for taskbar)
	if procSendMessageW.Find() == nil {
		if hSmall != 0 {
			_, _, _ = procSendMessageW.Call(dw.hwnd, WM_SETICON, uintptr(ICON_SMALL), hSmall)
		}
		if hBig != 0 {
			_, _, _ = procSendMessageW.Call(dw.hwnd, WM_SETICON, uintptr(ICON_BIG), hBig)
		}
	}

	// 3. Set window class icons
	if procSetClassLongPtrW.Find() == nil {
		if hSmall != 0 {
			_, _, _ = procSetClassLongPtrW.Call(dw.hwnd, GCLP_HICONSM, hSmall)
		}
		if hBig != 0 {
			_, _, _ = procSetClassLongPtrW.Call(dw.hwnd, GCLP_HICON, hBig)
		}
	} else if procSetClassLongW.Find() == nil {
		if hSmall != 0 {
			_, _, _ = procSetClassLongW.Call(dw.hwnd, GCLP_HICONSM, hSmall)
		}
		if hBig != 0 {
			_, _, _ = procSetClassLongW.Call(dw.hwnd, GCLP_HICON, hBig)
		}
	}

	// 4. Update tray icon if present
	if dw.trayAdded && hSmall != 0 {
		var nid NOTIFYICONDATAW
		nid.cbSize = uint32(unsafe.Sizeof(nid))
		nid.hWnd = dw.hwnd
		nid.uID = 1
		nid.uFlags = NIF_ICON
		nid.hIcon = hSmall
		_, _, _ = procShell_NotifyIconW.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&nid)))
	}

	// 5. Force frame recalculation so title bar redraws immediately with new icon
	if procSetWindowPos.Find() == nil {
		_, _, _ = procSetWindowPos.Call(
			dw.hwnd,
			0,
			0, 0, 0, 0,
			SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED,
		)
	}
}

func (dw *DesktopWindow) initPlatform(title string, width, height int) {
	if dw.hwnd == 0 {
		return
	}

	setAppUserModelID()

	windowMu.Lock()
	globalWindow = dw
	forceQuit = false
	windowMu.Unlock()

	// Subclass window to intercept minimize and close
	cb := windows.NewCallback(customWndProc)
	originalWndProc = w32.SetWindowLongPtr(w32.HWND(dw.hwnd), w32.GWLP_WNDPROC, cb)

	dw.loadIcons()
	dw.applyConfigTheme()
	centerWindow(dw.hwnd, width, height)

	// Add tray icon
	dw.addTrayIcon(title)
}

func (dw *DesktopWindow) addTrayIcon(title string) {
	hIcon := dw.currentIconSmall
	if hIcon == 0 {
		var hinstance windows.Handle
		_ = windows.GetModuleHandleEx(0, nil, &hinstance)
		hIcon, _, _ = procLoadIconW.Call(uintptr(hinstance), 32512)
		if hIcon == 0 {
			hIcon, _, _ = procLoadIconW.Call(0, 32512) // IDI_APPLICATION
		}
	}

	var nid NOTIFYICONDATAW
	nid.cbSize = uint32(unsafe.Sizeof(nid))
	nid.hWnd = dw.hwnd
	nid.uID = 1
	nid.uFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	nid.uCallbackMessage = WM_TRAY_CALLBACK
	nid.hIcon = hIcon

	tip, _ := windows.UTF16FromString(title)
	copy(nid.szTip[:], tip)

	r, _, _ := procShell_NotifyIconW.Call(NIM_ADD, uintptr(unsafe.Pointer(&nid)))
	if r != 0 {
		dw.trayAdded = true
	}
}

func (dw *DesktopWindow) removeTrayIcon() {
	if !dw.trayAdded || dw.hwnd == 0 {
		return
	}
	var nid NOTIFYICONDATAW
	nid.cbSize = uint32(unsafe.Sizeof(nid))
	nid.hWnd = dw.hwnd
	nid.uID = 1
	_, _, _ = procShell_NotifyIconW.Call(NIM_DELETE, uintptr(unsafe.Pointer(&nid)))
	dw.trayAdded = false
}

func (dw *DesktopWindow) showTrayMenu() {
	if dw == nil || dw.hwnd == 0 {
		return
	}
	var pt POINT
	_, _, _ = procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	hMenu, _, _ := procCreatePopupMenu.Call()
	if hMenu == 0 {
		return
	}
	defer procDestroyMenu.Call(hMenu)

	openLabel := "Open SmirnovTorrent"
	exitLabel := "Exit"
	if dw.cfgManager != nil && dw.cfgManager.Get().Language == "ru" {
		openLabel = "Открыть SmirnovTorrent"
		exitLabel = "Выход"
	}

	openText, _ := windows.UTF16PtrFromString(openLabel)
	exitText, _ := windows.UTF16PtrFromString(exitLabel)

	_, _, _ = procAppendMenuW.Call(hMenu, MF_STRING, ID_TRAY_SHOW, uintptr(unsafe.Pointer(openText)))
	_, _, _ = procAppendMenuW.Call(hMenu, MF_SEPARATOR, 0, 0)
	_, _, _ = procAppendMenuW.Call(hMenu, MF_STRING, ID_TRAY_EXIT, uintptr(unsafe.Pointer(exitText)))

	_, _, _ = procSetForegroundWindow.Call(dw.hwnd)
	cmd, _, _ := procTrackPopupMenu.Call(hMenu, TPM_RIGHTBUTTON|TPM_RETURNCMD, uintptr(pt.X), uintptr(pt.Y), 0, dw.hwnd, 0)
	_, _, _ = procPostMessageW.Call(dw.hwnd, 0, 0, 0) // WM_NULL required to dismiss menu properly

	if cmd == ID_TRAY_SHOW {
		dw.Show()
	} else if cmd == ID_TRAY_EXIT {
		windowMu.Lock()
		forceQuit = true
		windowMu.Unlock()
		dw.Destroy()
	}
}

func (dw *DesktopWindow) Show() {
	if dw.hwnd != 0 {
		dw.applyConfigTheme()
		w32.ShowWindow(w32.HWND(dw.hwnd), w32.SW_RESTORE)
		w32.SetForegroundWindow(w32.HWND(dw.hwnd))
	}
}

func (dw *DesktopWindow) Hide() {
	if dw.hwnd != 0 {
		w32.ShowWindow(w32.HWND(dw.hwnd), w32.SW_HIDE)
	}
}

func (dw *DesktopWindow) cleanupPlatform() {
	dw.removeTrayIcon()
	windowMu.Lock()
	if globalWindow == dw {
		globalWindow = nil
	}
	windowMu.Unlock()
	dw.hwnd = 0

	if dw.iconDarkSmall != 0 {
		_, _, _ = procDestroyIcon.Call(dw.iconDarkSmall)
		dw.iconDarkSmall = 0
	}
	if dw.iconDarkBig != 0 {
		_, _, _ = procDestroyIcon.Call(dw.iconDarkBig)
		dw.iconDarkBig = 0
	}
	if dw.iconLightSmall != 0 {
		_, _, _ = procDestroyIcon.Call(dw.iconLightSmall)
		dw.iconLightSmall = 0
	}
	if dw.iconLightBig != 0 {
		_, _, _ = procDestroyIcon.Call(dw.iconLightBig)
		dw.iconLightBig = 0
	}
	dw.currentIconSmall = 0
	dw.currentIconBig = 0
}
