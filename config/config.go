package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

// Cấu hình Tool
type ToolConfig struct {
	DelayMin          int    // Delay tối thiểu giữa các nhiệm vụ
	DelayMax          int    // Delay tối đa giữa các nhiệm vụ
	BlockAfter        int    // Sau bao nhiêu nhiệm vụ thì chống block
	RestAfter         int    // Sau bao nhiêu nhiệm vụ thì nghỉ ngơi
	ChangeNickAfter   int    // Sau bao nhiêu nhiệm vụ thì đổi nick
	DeleteCookieAfter int    // Sau bao nhiêu nhiệm vụ thì xóa cookie
	AccessTokenTDS    string // Token TDS
}

// Hàm nạp cấu hình từ .env hoặc file cấu hình
func LoadConfig() ToolConfig {
	// Load tệp .env để lấy thông tin môi trường nếu có
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Trả về cấu hình với các tham số từ .env hoặc mặc định
	return ToolConfig{
		AccessTokenTDS: os.Getenv("ACCESS_TOKEN"), // Tải TDS token từ biến môi trường
	}
}
