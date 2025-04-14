package utils

import (
	"fmt"
	"time"
)

// In log với timestamp
func Info(msg string) {
	fmt.Printf("[INFO] %s - %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func Error(msg string) {
	fmt.Printf("[ERROR] %s - %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}
