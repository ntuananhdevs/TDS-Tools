package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ACCESS_TOKEN   string
	RequestDelay   time.Duration
	MaxConcurrency int
	LogLevel       string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: Không tìm thấy file .env, sử dụng biến môi trường")
	}

	tdsToken := os.Getenv("ACCESS_TOKEN")
	if tdsToken == "" {
		return nil, &ConfigError{"TDS_TOKEN không được cung cấp trong file .env"}
	}

	delayStr := os.Getenv("REQUEST_DELAY")
	delay := 5 * time.Second
	if delayStr != "" {
		if delayInt, err := strconv.Atoi(delayStr); err == nil {
			delay = time.Duration(delayInt) * time.Second
		}
	}

	maxConcurrencyStr := os.Getenv("MAX_CONCURRENCY")
	maxConcurrency := 3
	if maxConcurrencyStr != "" {
		if concInt, err := strconv.Atoi(maxConcurrencyStr); err == nil {
			maxConcurrency = concInt
		}
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		ACCESS_TOKEN:   tdsToken,
		RequestDelay:   delay,
		MaxConcurrency: maxConcurrency,
		LogLevel:       logLevel,
	}, nil
}

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}

func (c *Config) GetTDSToken() string {
	return c.ACCESS_TOKEN
}

func (c *Config) GetRequestDelay() time.Duration {
	return c.RequestDelay
}

func (c *Config) GetMaxConcurrency() int {
	return c.MaxConcurrency
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}