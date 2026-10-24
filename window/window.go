package window

import (
	"github.com/jchv/go-webview2"
	"smirnovtorrent/config"
)

type DesktopWindow struct {
	view             webview2.WebView
	hwnd             uintptr
	cfgManager       *config.ConfigManager
	trayAdded        bool
	iconDarkSmall    uintptr
	iconDarkBig      uintptr
	iconLightSmall   uintptr
	iconLightBig     uintptr
	currentIconSmall uintptr
	currentIconBig   uintptr
	currentTheme     string
}

// SetAppUserModelID registers the application identifier with Windows shell.
func SetAppUserModelID() {
	setAppUserModelIDPlatform()
}

func NewDesktopWindow(title string, url string, width, height int, debug bool, cfgManager *config.ConfigManager) *DesktopWindow {
	w := webview2.New(debug)
	if w == nil {
		return nil
	}

	w.SetTitle(title)
	w.SetSize(width, height, webview2.HintNone)
	w.Navigate(url)

	dw := &DesktopWindow{
		view:       w,
		hwnd:       uintptr(w.Window()),
		cfgManager: cfgManager,
	}

	dw.initPlatform(title, width, height)
	return dw
}

func (dw *DesktopWindow) SetTheme(theme string) {
	dw.setThemePlatform(theme)
}

func (dw *DesktopWindow) Run() {
	if dw.view != nil {
		dw.view.Run()
	}
}

func (dw *DesktopWindow) Destroy() {
	dw.cleanupPlatform()
	if dw.view != nil {
		dw.view.Destroy()
		dw.view = nil
	}
}
