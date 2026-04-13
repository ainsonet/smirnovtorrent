package logger

import (
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	logger := New(&Config{
		Level: INFO,
	})
	defer logger.Close()
	
	// DEBUG не должен выводиться при уровне INFO
	// (визуально проверяем что нет вывода)
	logger.Debug("This should not appear")
	
	// INFO должен выводиться
	logger.Info("Info message")
	
	// WARN должен выводиться
	logger.Warn("Warning message")
	
	// ERROR должен выводиться
	logger.Error("Error message")
}

func TestLoggerWithPrefix(t *testing.T) {
	logger := New(&Config{
		Level:  DEBUG,
		Prefix: "[TEST]",
	})
	defer logger.Close()
	
	logger.Info("Message with prefix")
}

func TestLoggerSetLevel(t *testing.T) {
	logger := New(&Config{
		Level: INFO,
	})
	defer logger.Close()
	
	// Меняем уровень на ERROR
	logger.SetLevel(ERROR)
	
	// INFO не должен выводиться
	logger.Info("This should not appear")
	
	// ERROR должен выводиться
	logger.Error("This should appear")
}

func TestGlobalLogger(t *testing.T) {
	// Проверяем что глобальный логгер существует
	logger := GetGlobal()
	if logger == nil {
		t.Fatal("Expected global logger to exist")
	}
	
	// Используем helper функции
	Info("Global info message")
	Warn("Global warning message")
	Error("Global error message")
}

func BenchmarkLogger(b *testing.B) {
	logger := New(&Config{
		Level: INFO,
	})
	defer logger.Close()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Test message %d", i)
	}
}
