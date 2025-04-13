package config

import (
	"os"
	"time"
)

var (
	AccessTokenTDS      = getAccessToken()  // Đọc token từ file configtds.txt
	DelayMin            = 3 * time.Second
	DelayMax            = 6 * time.Second
	JobsBeforeRest      = 10
	RestDuration        = 15 * time.Second
	JobsBeforeSwitchAcc = 30
)

// Hàm để lấy token TDS từ file configtds.txt
func getAccessToken() string {
	// Đọc token từ file assets/configtds.txt
	data, err := os.ReadFile("../assets/configtds.txt")
	if err != nil {
		panic("Không tìm thấy token TDS. Đặt file assets/configtds.txt hoặc kiểm tra quyền đọc file")
	}
	return string(data)
}
