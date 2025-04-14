package utils

import (
	"fmt"
	"time"
)

// LogLevel đại diện cho mức độ log
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

// Logger đại diện cho một logger với màu sắc
type Logger struct {
	MinLevel LogLevel
}

// NewLogger tạo một logger mới với mức log tối thiểu
func NewLogger(levelStr string) *Logger {
	var level LogLevel
	
	switch levelStr {
	case "debug":
		level = DEBUG
	case "info":
		level = INFO
	case "warn":
		level = WARN
	case "error":
		level = ERROR
	default:
		level = INFO
	}
	
	return &Logger{
		MinLevel: level,
	}
}

// Debug ghi log ở mức debug
func (l *Logger) Debug(message string) {
	if l.MinLevel <= DEBUG {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s%s [DEBUG] %s%s\n", Cyan, timestamp, message, Reset)
	}
}

// Info ghi log ở mức info
func (l *Logger) Info(message string) {
	if l.MinLevel <= INFO {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s%s [INFO] %s%s\n", Green, timestamp, message, Reset)
	}
}

// Warn ghi log ở mức warn
func (l *Logger) Warn(message string) {
	if l.MinLevel <= WARN {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s%s [WARN] %s%s\n", Yellow, timestamp, message, Reset)
	}
}

// Error ghi log ở mức error
func (l *Logger) Error(message string) {
	if l.MinLevel <= ERROR {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s%s [ERROR] %s%s\n", Red, timestamp, message, Reset)
	}
}
