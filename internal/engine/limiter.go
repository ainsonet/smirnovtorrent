package engine

import (
	"fmt"
	"sync"
	"time"

	"smirnovtorrent/internal/ratelimit"
)

// RateLimiter ограничивает скорость передачи
type RateLimiter struct {
	downloadLimiter *ratelimit.RateLimiter
	uploadLimiter   *ratelimit.RateLimiter
	mu              sync.RWMutex
}

// NewRateLimiter создаёт новый лимитер скорости
// maxDownloadRate - максимальная скорость загрузки в байт/сек (0 = без ограничений)
// maxUploadRate - максимальная скорость отдачи в байт/сек (0 = без ограничений)
func NewRateLimiter(maxDownloadRate, maxUploadRate int64) *RateLimiter {
	return &RateLimiter{
		downloadLimiter: ratelimit.NewRateLimiter(maxDownloadRate, time.Second),
		uploadLimiter:   ratelimit.NewRateLimiter(maxUploadRate, time.Second),
	}
}

// Start запускает цикл лимитирования
func (rl *RateLimiter) Start() {
	// В новой реализации лимитер работает автоматически через Allow/Wait
}

// Stop останавливает лимитер
func (rl *RateLimiter) Stop() {
	// В новой реализации нет необходимости в остановке
}

// WaitDownload ожидает разрешения на загрузку n байт
func (rl *RateLimiter) WaitDownload(n int64) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	rl.downloadLimiter.Wait(n)
}

// WaitUpload ожидает разрешения на отдачу n байт
func (rl *RateLimiter) WaitUpload(n int64) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	rl.uploadLimiter.Wait(n)
}

// SetMaxDownloadRate устанавливает максимальную скорость загрузки
func (rl *RateLimiter) SetMaxDownloadRate(rate int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.downloadLimiter.SetLimit(rate)
}

// SetMaxUploadRate устанавливает максимальную скорость отдачи
func (rl *RateLimiter) SetMaxUploadRate(rate int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.uploadLimiter.SetLimit(rate)
}

// GetMaxDownloadRate возвращает максимальную скорость загрузки
func (rl *RateLimiter) GetMaxDownloadRate() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.downloadLimiter.GetLimit()
}

// GetMaxUploadRate возвращает максимальную скорость отдачи
func (rl *RateLimiter) GetMaxUploadRate() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.uploadLimiter.GetLimit()
}

// GetDownloadRate возвращает доступное количество байт для загрузки
func (rl *RateLimiter) GetDownloadRate() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.downloadLimiter.GetAvailable()
}

// GetUploadRate возвращает доступное количество байт для отдачи
func (rl *RateLimiter) GetUploadRate() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.uploadLimiter.GetAvailable()
}

// Unlimit снимает все ограничения
func (rl *RateLimiter) Unlimit() {
	rl.SetMaxDownloadRate(0)
	rl.SetMaxUploadRate(0)
}

// FormatRate форматирует скорость для отображения
func FormatRate(bytesPerSecond int64) string {
	const unit = 1024
	if bytesPerSecond < unit {
		return fmt.Sprintf("%d B/s", bytesPerSecond)
	}
	div, exp := int64(unit), 0
	for n := bytesPerSecond / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB"}
	return fmt.Sprintf("%.1f %s/s", float64(bytesPerSecond)/float64(div), suffixes[exp])
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
