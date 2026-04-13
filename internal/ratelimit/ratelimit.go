package ratelimit

import (
	"sync"
	"time"
)

// RateLimiter ограничивает скорость передачи данных
type RateLimiter struct {
	mu         sync.Mutex
	bytesLeft  int64
	maxBytes   int64
	lastRefill time.Time
	interval   time.Duration
}

// NewRateLimiter создаёт новый ограничитель скорости
// maxBytes - максимальное количество байт за интервал
// interval - интервал времени (обычно 1 секунда)
func NewRateLimiter(maxBytes int64, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		bytesLeft:  maxBytes,
		maxBytes:   maxBytes,
		lastRefill: time.Now(),
		interval:   interval,
	}
}

// Allow запрашивает разрешение на передачу n байт
// Возвращает время через которое можно попробовать снова (0 если можно сразу)
func (r *RateLimiter) Allow(n int64) time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Если лимит 0 или отрицательный - без ограничений
	if r.maxBytes <= 0 {
		return 0
	}

	// Пополняем запас байтов
	r.refill()

	// Если достаточно байтов - разрешаем
	if r.bytesLeft >= n {
		r.bytesLeft -= n
		return 0
	}

	// Если нет - вычисляем время ожидания
	if r.bytesLeft > 0 {
		// Частично разрешаем
		r.bytesLeft = 0
	}

	// Время до следующего пополнения
	waitTime := r.interval - time.Since(r.lastRefill)
	if waitTime < 0 {
		waitTime = 0
	}
	return waitTime
}

// refill пополняет запас байтов
func (r *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	// Сколько интервалов прошло
	intervals := int64(elapsed / r.interval)
	if intervals == 0 {
		return
	}

	// Пополняем запас
	r.bytesLeft += intervals * r.maxBytes
	if r.bytesLeft > r.maxBytes {
		r.bytesLeft = r.maxBytes
	}

	// Обновляем время последнего пополнения
	r.lastRefill = now
}

// SetLimit устанавливает новый лимит
func (r *RateLimiter) SetLimit(maxBytes int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.maxBytes = maxBytes
	if r.bytesLeft > maxBytes {
		r.bytesLeft = maxBytes
	}
}

// GetLimit возвращает текущий лимит
func (r *RateLimiter) GetLimit() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.maxBytes
}

// GetAvailable возвращает доступное количество байт
func (r *RateLimiter) GetAvailable() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.refill()
	return r.bytesLeft
}

// Wait ожидает разрешения на передачу n байт
func (r *RateLimiter) Wait(n int64) {
	for {
		waitTime := r.Allow(n)
		if waitTime == 0 {
			return
		}
		time.Sleep(waitTime)
	}
}
