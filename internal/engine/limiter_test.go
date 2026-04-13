package engine

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(1024, 512)
	
	if rl == nil {
		t.Fatal("Expected non-nil rate limiter")
	}
	
	if rl.downloadLimiter == nil {
		t.Error("Expected download limiter to be created")
	}
	
	if rl.uploadLimiter == nil {
		t.Error("Expected upload limiter to be created")
	}
}

func TestRateLimiterUnlimited(t *testing.T) {
	rl := NewRateLimiter(0, 0)
	
	// Должно работать без ограничений
	done := make(chan bool)
	go func() {
		rl.WaitDownload(1000000)
		rl.WaitUpload(1000000)
		done <- true
	}()
	
	select {
	case <-done:
		// Успешно
	case <-time.After(2 * time.Second):
		t.Error("Rate limiter should not block when unlimited")
	}
}

func TestRateLimiterSetLimits(t *testing.T) {
	rl := NewRateLimiter(1000, 500)
	
	rl.SetMaxDownloadRate(2000)
	if rl.GetMaxDownloadRate() != 2000 {
		t.Errorf("Expected download rate 2000, got %d", rl.GetMaxDownloadRate())
	}
	
	rl.SetMaxUploadRate(1000)
	if rl.GetMaxUploadRate() != 1000 {
		t.Errorf("Expected upload rate 1000, got %d", rl.GetMaxUploadRate())
	}
}

func TestRateLimiterUnlimit(t *testing.T) {
	rl := NewRateLimiter(1000, 500)
	
	rl.Unlimit()
	
	if rl.GetMaxDownloadRate() != 0 {
		t.Error("Expected download rate 0 after unlimit")
	}
	
	if rl.GetMaxUploadRate() != 0 {
		t.Error("Expected upload rate 0 after unlimit")
	}
}

func TestRateLimiterWaitDownload(t *testing.T) {
	rl := NewRateLimiter(1000, 0) // 1000 bytes/sec download
	
	// Первые 1000 байт должны пройти быстро
	start := time.Now()
	rl.WaitDownload(500)
	rl.WaitDownload(500)
	elapsed := time.Since(start)
	
	if elapsed > 100*time.Millisecond {
		t.Logf("Warning: download took longer than expected: %v", elapsed)
	}
	
// Следующие байты должны ждать
	start = time.Now()
	rl.WaitDownload(100)
	elapsed = time.Since(start)
	
	if elapsed < 500*time.Millisecond {
		t.Logf("Expected to wait for rate limit, waited: %v", elapsed)
	}
}

func TestRateLimiterGetRates(t *testing.T) {
	rl := NewRateLimiter(1000, 500)
	
	dlRate := rl.GetDownloadRate()
	if dlRate != 1000 {
		t.Errorf("Expected download available 1000, got %d", dlRate)
	}
	
	ulRate := rl.GetUploadRate()
	if ulRate != 500 {
		t.Errorf("Expected upload available 500, got %d", ulRate)
	}
}

func TestFormatRate(t *testing.T) {
	tests := []struct {
		rate     int64
		expected string
	}{
		{100, "100 B/s"},
		{1024, "1.0 KB/s"},
		{1536, "1.5 KB/s"},
		{1048576, "1.0 MB/s"},
	}
	
	for _, test := range tests {
		result := FormatRate(test.rate)
		if result != test.expected {
			t.Errorf("FormatRate(%d) = %s, expected %s", test.rate, result, test.expected)
		}
	}
}

func BenchmarkRateLimiterDownload(b *testing.B) {
	rl := NewRateLimiter(1000000, 0) // 1 MB/s
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.WaitDownload(100)
	}
}

func BenchmarkRateLimiterUpload(b *testing.B) {
	rl := NewRateLimiter(0, 1000000) // 1 MB/s
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.WaitUpload(100)
	}
}
