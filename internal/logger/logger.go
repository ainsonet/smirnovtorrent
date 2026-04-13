package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Level уровень логирования
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// String возвращает строковое представление уровня
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger структурированный логгер
type Logger struct {
	level    Level
	prefix   string
	mu       sync.Mutex
	file     *os.File
	fileMu   sync.Mutex
}

// Config конфигурация логгера
type Config struct {
	Level  Level
	Prefix string
	File   string // путь к файлу логов (опционально)
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Level:  INFO,
		Prefix: "",
		File:   "",
	}
}

// New создаёт новый логгер
func New(config *Config) *Logger {
	logger := &Logger{
		level:  config.Level,
		prefix: config.Prefix,
	}
	
	// Открываем файл для логов если указан
	if config.File != "" {
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Warning: failed to open log file: %v", err)
		} else {
			logger.file = file
		}
	}
	
	return logger
}

// log записывает сообщение в лог
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	
	// Форматируем сообщение
	logEntry := fmt.Sprintf("[%s] [%s] %s %s", 
		timestamp, 
		level.String(),
		l.prefix,
		message)
	
	// Пишем в stdout
	fmt.Println(logEntry)
	
	// Пишем в файл если открыт
	if l.file != nil {
		l.fileMu.Lock()
		l.file.WriteString(logEntry + "\n")
		l.fileMu.Unlock()
	}
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info логирует информационное сообщение
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn логирует предупреждение
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error логирует ошибку
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// SetLevel устанавливает уровень логирования
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Close закрывает логгер
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Global логгер по умолчанию
var globalLogger = New(DefaultConfig())

// SetGlobal устанавливает глобальный логгер
func SetGlobal(logger *Logger) {
	globalLogger = logger
}

// GetGlobal возвращает глобальный логгер
func GetGlobal() *Logger {
	return globalLogger
}

// Helper функции для удобства
func Debug(format string, args ...interface{}) {
	globalLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	globalLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	globalLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	globalLogger.Error(format, args...)
}
