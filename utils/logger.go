package utils

import (
	"fmt"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func Info(msg string) {
	fmt.Printf("%s[%s] %s%s\n", Blue, timestamp(), msg, Reset)
}

func Success(msg string) {
	fmt.Printf("%s[%s] ✅ %s%s\n", Green, timestamp(), msg, Reset)
}

func Warning(msg string) {
	fmt.Printf("%s[%s] ⚠️ %s%s\n", Yellow, timestamp(), msg, Reset)
}

func Error(msg string) {
	fmt.Printf("%s[%s] ❌ %s%s\n", Red, timestamp(), msg, Reset)
}
