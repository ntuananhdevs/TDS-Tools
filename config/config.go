// config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config chứa các thiết lập cho ứng dụng
type Config struct {
	ACCESS_TOKEN       string        // Token xác thực cho TDS API
	RequestDelay   time.Duration // Độ trễ giữa các request để tránh bị chặn
	MaxConcurrency int           // Số lượng tài khoản Facebook có thể chạy đồng thời
	LogLevel       string        // Mức độ log (debug, info, warn, error)
}

// LoadConfig tải cấu hình từ file .env
func LoadConfig() (*Config, error) {
	// Tải file .env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: Không tìm thấy file .env, sử dụng biến môi trường")
	}

	// Lấy token từ biến môi trường
	tdsToken := os.Getenv("ACCESS_TOKEN")
	if tdsToken == "" {
		return nil, &ConfigError{"TDS_TOKEN không được cung cấp trong file .env"}
	}

	// Lấy các cấu hình khác với giá trị mặc định
	delayStr := os.Getenv("REQUEST_DELAY")
	delay := 5 * time.Second // Mặc định 5 giây
	if delayStr != "" {
		if delayInt, err := strconv.Atoi(delayStr); err == nil {
			delay = time.Duration(delayInt) * time.Second
		}
	}

	maxConcurrencyStr := os.Getenv("MAX_CONCURRENCY")
	maxConcurrency := 3 // Mặc định 3 tài khoản đồng thời
	if maxConcurrencyStr != "" {
		if concInt, err := strconv.Atoi(maxConcurrencyStr); err == nil {
			maxConcurrency = concInt
		}
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // Mặc định log level
	}

	return &Config{
		ACCESS_TOKEN:       tdsToken,
		RequestDelay:   delay,
		MaxConcurrency: maxConcurrency,
		LogLevel:       logLevel,
	}, nil
}

// ConfigError là lỗi tùy chỉnh cho vấn đề cấu hình
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}

// GetTDSToken trả về token TDS
func (c *Config) GetTDSToken() string {
	return c.ACCESS_TOKEN
}

// GetRequestDelay trả về độ trễ giữa các request
func (c *Config) GetRequestDelay() time.Duration {
	return c.RequestDelay
}

// GetMaxConcurrency trả về số lượng tài khoản tối đa chạy đồng thời
func (c *Config) GetMaxConcurrency() int {
	return c.MaxConcurrency
}

// GetLogLevel trả về mức độ log
func (c *Config) GetLogLevel() string {
	return c.LogLevel
}