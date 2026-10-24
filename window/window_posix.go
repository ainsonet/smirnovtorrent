//go:build !windows

package window

func (dw *DesktopWindow) initPlatform(title string, width, height int) {}
func (dw *DesktopWindow) cleanupPlatform()                             {}
func (dw *DesktopWindow) Show()                                        {}
func (dw *DesktopWindow) Hide()                                        {}
func (dw *DesktopWindow) setThemePlatform(theme string)                {}
func setAppUserModelIDPlatform()                                       {}
