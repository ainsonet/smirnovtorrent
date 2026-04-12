package engine

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter ограничивает скорость передачи
type RateLimiter struct {
	maxDownloadRate int64 // байт в секунду
	maxUploadRate   int64 // байт в секунду
	currentDL       int64
	currentUL       int64
	mu              sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
	ticker          *time.Ticker
}

// NewRateLimiter создаёт новый лимитер скорости
func NewRateLimiter(maxDownloadRate, maxUploadRate int64) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &RateLimiter{
		maxDownloadRate: maxDownloadRate,
		maxUploadRate:   maxUploadRate,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start запускает цикл лимитирования
func (rl *RateLimiter) Start() {
	rl.ticker = time.NewTicker(1 * time.Second)
	
	go func() {
		for {
			select {
			case <-rl.ctx.Done():
				rl.ticker.Stop()
				return
			case <-rl.ticker.C:
				rl.mu.Lock()
				// Сбрасываем счётчики
				rl.currentDL = 0
				rl.currentUL = 0
				rl.mu.Unlock()
			}
		}
	}()
}

// Stop останавливает лимитер
func (rl *RateLimiter) Stop() {
	rl.cancel()
}

// RecordDownload записывает загруженные байты
func (rl *RateLimiter) RecordDownload(bytes int64) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.maxDownloadRate > 0 {
		rl.currentDL += bytes
		
		// Проверяем лимит
		if rl.currentDL > rl.maxDownloadRate {
			// Ждём до следующего тика
			rl.waitForReset()
		}
	}

	return nil
}

// RecordUpload записывает отправленные байты
func (rl *RateLimiter) RecordUpload(bytes int64) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.maxUploadRate > 0 {
		rl.currentUL += bytes
		
		// Проверяем лимит
		if rl.currentUL > rl.maxUploadRate {
			// Ждём до следующего тика
			rl.waitForReset()
		}
	}

	return nil
}

// waitForReset ждёт сброса счётчиков
func (rl *RateLimiter) waitForReset() {
	// Простая реализация - ждём 1 секунду
	// В production можно использовать более умный подход
	time.Sleep(1 * time.Second)
}

// SetMaxDownloadRate устанавливает максимальную скорость загрузки
func (rl *RateLimiter) SetMaxDownloadRate(rate int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.maxDownloadRate = rate
}

// SetMaxUploadRate устанавливает максимальную скорость отдачи
func (rl *RateLimiter) SetMaxUploadRate(rate int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.maxUploadRate = rate
}

// GetMaxDownloadRate возвращает максимальную скорость загрузки
func (rl *RateLimiter) GetMaxDownloadRate() int64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.maxDownloadRate
}

// GetMaxUploadRate возвращает максимальную скорость отдачи
func (rl *RateLimiter) GetMaxUploadRate() int64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.maxUploadRate
}

// GetDownloadRate возвращает текущую скорость загрузки
func (rl *RateLimiter) GetDownloadRate() int64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentDL
}

// GetUploadRate возвращает текущую скорость отдачи
func (rl *RateLimiter) GetUploadRate() int64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentUL
}

// Unlimiter снимает все ограничения
func (rl *RateLimiter) Unlimit() {
	rl.SetMaxDownloadRate(0)
	rl.SetMaxUploadRate(0)
}

// FormatRate форматирует скорость для отображения
func FormatRate(bytesPerSecond int64) string {
	const unit = 1024
	if bytesPerSecond < unit {
		return formatBytesFloat(float64(bytesPerSecond)) + "/s"
	}
	div, exp := int64(unit), 0
	for n := bytesPerSecond / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB"}
	return formatBytesFloat(float64(bytesPerSecond)/float64(div)) + suffixes[exp] + "/s"
}

// formatBytesFloat форматирует байты
func formatBytesFloat(bytes float64) string {
	const unit = 1024
	if bytes < unit {
		return formatFloat(bytes, 0) + " B"
	}
	div, exp := float64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB", "TB"}
	return formatFloat(bytes/div, 1) + " " + suffixes[exp]
}

// formatFloat форматирует float
func formatFloat(f float64, precision int) string {
	if precision == 0 {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprintf("%.1f", f)
}
