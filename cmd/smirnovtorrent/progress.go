package main

import (
	"fmt"
	"os"
	"time"
)

// ProgressBar простой прогресс-бар для терминала
type ProgressBar struct {
	width      int
	startTime  time.Time
	lastUpdate time.Time
}

// NewProgressBar создаёт новый прогресс-бар
func NewProgressBar(width int) *ProgressBar {
	return &ProgressBar{
		width:      width,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

// Show отображает прогресс
func (pb *ProgressBar) Show(progress float64, current, total int, activePeers int, downloadSpeed float64) {
	// Очищаем строку
	fmt.Print("\r")

	// Заполняем прогресс-бар
	barWidth := pb.width - 10 // Оставляем место для процентов
	filled := int(progress / 100.0 * float64(barWidth))

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"

	// Скорость в байтах/секунду
	speedStr := formatBytesFloat(downloadSpeed) + "/s"

	// Время оставшееся
	elapsed := time.Since(pb.startTime).Seconds()
	var eta string
	if progress > 0 {
		remaining := (100.0 - progress) / progress * elapsed
		eta = formatDuration(time.Duration(remaining) * time.Second)
	} else {
		eta = "--:--"
	}

	// Выводим прогресс
	fmt.Printf("%s %.1f%% %d/%d | Peers: %d | %s | ETA: %s",
		bar, progress, current, total, activePeers, speedStr, eta)

	pb.lastUpdate = time.Now()
}

// Finish завершает прогресс-бар
func (pb *ProgressBar) Finish() {
	fmt.Println()
	fmt.Println()
	fmt.Println("Download completed!")
	fmt.Printf("Total time: %s\n", formatDuration(time.Since(pb.startTime)))
}

// formatBytes форматирует размер в байтах
func formatBytes(bytes float64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%.1f B", bytes)
	}
	div, exp := float64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", bytes/div, suffixes[exp])
}

// formatDuration форматирует длительность
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// ClearLine очищает текущую строку
func ClearLine() {
	fmt.Print("\r\x1b[K")
}

// GetCurrentTime возвращает текущее время для логов
func GetCurrentTime() string {
	return time.Now().Format("15:04:05")
}

// IsTerminal проверяет работает ли в терминале
func IsTerminal() bool {
	// Простая проверка - если stdout это терминал
	stat, _ := os.Stdout.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}