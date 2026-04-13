package main

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// DownloadInfo информация о загрузке
type DownloadInfo struct {
	Name          string
	Progress      float32
	Downloaded    int64
	Uploaded      int64
	DownloadSpeed float64
	UploadSpeed   float64
	Peers         int
	Status        string
}

// MainWindow главное окно
type MainWindow struct {
	downloads   []DownloadInfo
	downloadsMu sync.RWMutex
	list        *widget.List
}

// NewMainWindow создаёт главное окно
func NewMainWindow() *MainWindow {
	mw := &MainWindow{
		downloads: make([]DownloadInfo, 0),
	}
	mw.list = widget.NewList(
		func() int { return len(mw.downloads) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Name"),
				widget.NewProgressBar(),
				widget.NewLabel("0 KB/s"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(mw.downloads) {
				return
			}
			
			dl := mw.downloads[id]
			box := obj.(*fyne.Container)
			
			// Заголовок
			title := box.Objects[0].(*widget.Label)
			title.SetText(dl.Name)
			
			// Прогресс
			progress := box.Objects[1].(*widget.ProgressBar)
			progress.SetValue(dl.Progress)
			
			// Статус
			status := box.Objects[2].(*widget.Label)
			status.SetText(fmt.Sprintf("↓ %s/s ↑ %s/s | Peers: %d | %s",
				formatSpeed(dl.DownloadSpeed),
				formatSpeed(dl.UploadSpeed),
				dl.Peers,
				dl.Status))
		},
	)
	
	return mw
}

// CreateUI создаёт интерфейс
func (mw *MainWindow) CreateUI() fyne.CanvasObject {
	// Поле ввода
	input := widget.NewEntry()
	input.SetPlaceHolder("Paste magnet link or .torrent path...")
	
	// Кнопка добавить
	addBtn := widget.NewButtonWithIcon("Add Torrent", theme.ContentAddIcon(), func() {
		if input.Text == "" {
			return
		}
		
		mw.downloadsMu.Lock()
		mw.downloads = append(mw.downloads, DownloadInfo{
			Name:          input.Text,
			Progress:      0,
			Downloaded:    0,
			Uploaded:      0,
			DownloadSpeed: 0,
			UploadSpeed:   0,
			Peers:         0,
			Status:        "downloading",
		})
		mw.downloadsMu.Unlock()
		
		mw.list.Refresh()
		input.SetText("")
	})
	
	// Статистика
	downloadSpeed := widget.NewLabel("0 KB/s")
	uploadSpeed := widget.NewLabel("0 KB/s")
	peers := widget.NewLabel("0")
	torrents := widget.NewLabel("0")
	
	stats := container.NewHBox(
		container.NewVBox(
			widget.NewLabel("Download"),
			downloadSpeed,
		),
		layout.NewSpacer(),
		container.NewVBox(
			widget.NewLabel("Upload"),
			uploadSpeed,
		),
		layout.NewSpacer(),
		container.NewVBox(
			widget.NewLabel("Peers"),
			peers,
		),
		layout.NewSpacer(),
		container.NewVBox(
			widget.NewLabel("Torrents"),
			torrents,
		),
	)
	
	// Основной контейнер
	mainContainer := container.NewVBox(
		// Заголовок
		container.NewHBox(
			widget.NewLabelWithStyle("🧲 SmirnovTorrent", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
		),
		
		// Добавление
		container.NewVBox(
			widget.NewLabel("Add Torrent"),
			container.NewHBox(input, addBtn),
		),
		
		// Статистика
		container.NewVBox(
			widget.NewSeparator(),
			stats,
			widget.NewSeparator(),
		),
		
		// Список загрузок
		container.NewVBox(
			widget.NewLabel("Active Torrents"),
			mw.list,
		),
	)
	
	return mainContainer
}

// UpdateStats обновляет статистику
func (mw *MainWindow) UpdateStats() {
	mw.downloadsMu.RLock()
	defer mw.downloadsMu.RUnlock()
	
	var totalDL, totalUL float64
	var totalPeers int
	
	for _, dl := range mw.downloads {
		if dl.Status == "downloading" {
			totalDL += dl.DownloadSpeed
			totalUL += dl.UploadSpeed
			totalPeers += dl.Peers
		}
	}
	
	// Обновляем данные для демо
	for i := range mw.downloads {
		if mw.downloads[i].Status == "downloading" {
			mw.downloads[i].Progress += 0.5
			if mw.downloads[i].Progress > 100 {
				mw.downloads[i].Progress = 100
			}
			mw.downloads[i].DownloadSpeed = float64(512*1024 + i*1024*1024)
			mw.downloads[i].UploadSpeed = float64(256*1024 + i*512*1024)
			mw.downloads[i].Peers = 5 + i*3
		}
	}
	
	mw.list.Refresh()
}

// formatSpeed форматирует скорость
func formatSpeed(bytes float64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%.0f B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", bytes/1024)
	}
	return fmt.Sprintf("%.1f MB", bytes/(1024*1024))
}

func main() {
	a := app.New()
	w := a.NewWindow("SmirnovTorrent v1.0.0")
	
	// Создаём главное окно
	mainWin := NewMainWindow()
	
	w.SetContent(mainWin.CreateUI())
	w.Resize(fyne.NewSize(900, 700))

	// Запускаем обновление каждые 2 секунды
	go func() {
		for {
			time.Sleep(2 * time.Second)
			mainWin.UpdateStats()
		}
	}()
	
	w.ShowAndRun()
}
	}()

	w.ShowAndRun()
}

// customTheme кастомная тема
type customTheme struct{}

func (c *customTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.RGBA{15, 23, 42, 255}
	case theme.ColorNamePrimary:
		return color.RGBA{99, 102, 241, 255}
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}

func (c *customTheme) Font(style fyne.TextStyle) fyne.Font {
	return theme.DefaultTheme().Font(style)
}

func (c *customTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (c *customTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
