package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"tds/config"
	"tds/internal/job"
	"tds/internal/traodoisub"
	"tds/models"
	"tds/utils"
)

func main() {
	// Nhập các giá trị cấu hình
	delay_min, err := getIntInput("Nhập Delay Min (ms): ")
	if err != nil {
		log.Fatal(err)
	}
	delay_max, err := getIntInput("Nhập Delay Max (ms): ")
	if err != nil {
		log.Fatal(err)
	}
	block_after, err := getIntInput("Sau Bao Nhiêu Nhiệm Vụ Thì Chống Block: ")
	if err != nil {
		log.Fatal(err)
	}
	rest_after, err := getIntInput("Sau Bao Nhiêu Nhiệm Vụ Thì Nghỉ Ngơi: ")
	if err != nil {
		log.Fatal(err)
	}
	change_nick_after, err := getIntInput("Sau Bao Nhiêu Nhiệm Vụ Thì Đổi Nick: ")
	if err != nil {
		log.Fatal(err)
	}
	delete_cookie_after, err := getIntInput("Lỗi Bao Nhiêu Nhiệm Vụ Thì Xóa Cookie: ")
	if err != nil {
		log.Fatal(err)
	}

	// Tạo cấu hình với các giá trị vừa nhập
	config := config.ToolConfig{
		DelayMin:           delay_min,
		DelayMax:           delay_max,
		BlockAfter:         block_after,
		RestAfter:          rest_after,
		ChangeNickAfter:    change_nick_after,
		DeleteCookieAfter:  delete_cookie_after,
	}

	// In thông tin cấu hình
	fmt.Println("\nCấu hình đã được thiết lập:")
	fmt.Printf("Delay Min: %d ms\n", config.DelayMin)
	fmt.Printf("Delay Max: %d ms\n", config.DelayMax)
	fmt.Printf("Block after %d tasks\n", config.BlockAfter)
	fmt.Printf("Rest after %d tasks\n", config.RestAfter)
	fmt.Printf("Change Nick after %d tasks\n", config.ChangeNickAfter)
	fmt.Printf("Delete Cookie after %d tasks\n", config.DeleteCookieAfter)

	// Lấy và in ra cookies
	cookies := loadCookies("../cookie.json")
	fmt.Println("\nThông tin Cookie:")
	for _, cookie := range cookies {
		fmt.Printf("Cookie Name: %s, Cookie Value: %s\n", cookie.UserID, cookie.Cookie)
	}

	// Khởi tạo TDSClient
	tds := traodoisub.NewClient(config.AccessTokenTDS)

	// Lấy tất cả nhiệm vụ từ Traodoisub
	jobs, err := tds.GetAllJobs()
	if err != nil {
		utils.Error("Không thể lấy nhiệm vụ: " + err.Error())
		log.Fatal(err)
	}

	// Gọi hàm để thực hiện các nhiệm vụ
	job.ExecuteTasks(config, jobs)
}

// Hàm nhận đầu vào từ người dùng và chuyển đổi thành số nguyên
func getIntInput(prompt string) (int, error) {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	// Chuyển đổi sang kiểu int
	num, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("input không hợp lệ: %v", err)
	}
	return num, nil
}

// Hàm đọc cookies từ file
func loadCookies(path string) []models.CookieUser {
	file, err := os.ReadFile(path)
	if err != nil {
		utils.Error("Không thể đọc cookie file: " + err.Error())
		os.Exit(1)
	}

	var cookies []models.CookieUser
	err = json.Unmarshal(file, &cookies)
	if err != nil {
		utils.Error("Cookie file không hợp lệ: " + err.Error())
		os.Exit(1)
	}

	return cookies
}
