package ratelimit

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(1024, time.Second)
	
	if rl.maxBytes != 1024 {
		t.Errorf("Expected maxBytes 1024, got %d", rl.maxBytes)
	}
	
	if rl.interval != time.Second {
		t.Errorf("Expected interval 1s, got %v", rl.interval)
	}
}

func TestRateLimiterUnlimited(t *testing.T) {
	rl := NewRateLimiter(0, time.Second)
	
	// Должно разрешать без ограничений
	for i := 0; i < 1000; i++ {
		waitTime := rl.Allow(100)
		if waitTime != 0 {
			t.Errorf("Expected no wait time for unlimited, got %v", waitTime)
		}
	}
}

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter(1000, time.Second)
	
	// Первые 1000 байт должны пройти
	waitTime := rl.Allow(500)
	if waitTime != 0 {
		t.Errorf("Expected no wait time, got %v", waitTime)
	}
	
	waitTime = rl.Allow(500)
	if waitTime != 0 {
		t.Errorf("Expected no wait time, got %v", waitTime)
	}
	
	// Следующие должны ждать
	waitTime = rl.Allow(100)
	if waitTime == 0 {
		t.Error("Expected wait time when limit exceeded")
	}
}

func TestRateLimiterSetLimit(t *testing.T) {
	rl := NewRateLimiter(1000, time.Second)
	
	rl.SetLimit(2000)
	if rl.GetLimit() != 2000 {
		t.Errorf("Expected limit 2000, got %d", rl.GetLimit())
	}
}

func TestRateLimiterGetAvailable(t *testing.T) {
	rl := NewRateLimiter(1000, time.Second)
	
	available := rl.GetAvailable()
	if available != 1000 {
		t.Errorf("Expected 1000 available, got %d", available)
	}
	
	rl.Allow(300)
	available = rl.GetAvailable()
	if available != 700 {
		t.Errorf("Expected 700 available, got %d", available)
	}
}

func TestRateLimiterRefill(t *testing.T) {
	rl := NewRateLimiter(1000, 100*time.Millisecond)
	
	// Потребляем весь лимит
	rl.Allow(1000)
	
	if rl.GetAvailable() != 0 {
		t.Errorf("Expected 0 available after consuming all")
	}
	
	// Ждём пополнения
	time.Sleep(150 * time.Millisecond)
	
	available := rl.GetAvailable()
	if available < 1000 {
		t.Errorf("Expected at least 1000 after refill, got %d", available)
	}
}

func TestRateLimiterPartialAllow(t *testing.T) {
	rl := NewRateLimiter(1000, time.Second)
	
	// Потребляем 900 байт
	rl.Allow(900)
	
	// Пытаемся потребовать 200 (осталось только 100)
	waitTime := rl.Allow(200)
	
	// Должно частично разрешить (100 байт) и вернуть время ожидания
	if waitTime == 0 {
		t.Error("Expected wait time when requesting more than available")
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	rl := NewRateLimiter(1000000, time.Second)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow(100)
	}
}
