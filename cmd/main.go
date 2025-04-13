package main

import (
	"encoding/json"
	"fmt"
	"os"
	"tds/config"
	"tds/internal/traodoisub"
	"tds/internal/facebook"
	"tds/models"
	"tds/utils"
	"strings"
	// "time"
	// "math/rand"
)

func main() {
	// Lấy token TDS từ cấu hình (đọc từ file)
	utils.Info("🚀 Đang khởi động...")

	// Đọc cookie từ file (lấy các tài khoản Facebook)
	cookies := loadCookies("../assets/Cookie_FB.txt")
	utils.Info("🚀 Bắt đầu chạy...")

	// Khởi tạo TDSClient
	tds := traodoisub.NewClient(config.AccessTokenTDS)

	// Lấy tất cả nhiệm vụ từ Traodoisub
	jobs, err := tds.GetAllJobs()
	if err != nil {
		utils.Error("Không thể lấy nhiệm vụ: " + err.Error())
		os.Exit(1)
	}

	// Hiển thị tất cả nhiệm vụ cho người dùng chọn
	fmt.Println("Danh sách các nhiệm vụ có sẵn:")
	jobTypes := []string{"like", "likegiare", "likesieure", "reaction", "comment", "share", "follow", "group", "page"}
	var availableJobs []string

	for _, jobType := range jobTypes {
		if jobList, found := jobs[jobType]; found && len(jobList) > 0 {
			availableJobs = append(availableJobs, jobType)
			fmt.Printf("- %s\n", jobType)
		}
	}

	if len(availableJobs) == 0 {
		fmt.Println("Không có nhiệm vụ nào có sẵn để chạy.")
		return
	}

	// Cho người dùng chọn nhiệm vụ cần thực hiện
	fmt.Print("\nChọn nhiệm vụ cần chạy (cách nhau bằng dấu phẩy): ")
	var userInput string
	fmt.Scanln(&userInput)

	// Tách các nhiệm vụ người dùng muốn chọn
	selectedJobs := strings.Split(userInput, ",")
	for i, job := range selectedJobs {
		selectedJobs[i] = strings.TrimSpace(job) // Loại bỏ các ký tự thừa (space)
	}

	// Chạy các nhiệm vụ đã chọn
	for _, jobType := range selectedJobs {
		if contains(availableJobs, jobType) {
			fmt.Printf("\nĐang thực hiện nhiệm vụ: %s\n", jobType)
			runJobForAllCookies(jobType, cookies, tds)
		} else {
			fmt.Printf("Nhiệm vụ '%s' không có sẵn.\n", jobType)
		}
	}
}

// Hàm kiểm tra nhiệm vụ có trong danh sách nhiệm vụ có sẵn không
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Hàm chạy các nhiệm vụ cho tất cả tài khoản Facebook
func runJobForAllCookies(jobType string, cookies []models.CookieUser, tds *traodoisub.TDSClient) {
	for _, user := range cookies {
		// Lấy dữ liệu Facebook cho user
		fb := facebook.NewFacebookClient(user.Cookie)
		utils.Info("🧩 Bắt đầu với: " + user.Name)

		// Lấy tất cả nhiệm vụ của loại jobType
		jobs, err := tds.GetJob(jobType)
		if err != nil {
			utils.Warning(fmt.Sprintf("Không có job cho %s: %s", user.Name, jobType))
			continue // Tiếp tục tìm nhiệm vụ khác
		}

		// Chạy các nhiệm vụ cho tài khoản
		for _, j := range jobs {
			if err := fb.Like(j.ID); err == nil {
				msg, xu, _ := tds.ConfirmJob(jobType, j.ID)
				utils.Success(fmt.Sprintf("%s LIKE %s | %s | %d xu", user.Name, j.ID, msg, xu))
			} else {
				utils.Warning(fmt.Sprintf("%s lỗi LIKE: %s", user.Name, j.ID))
			}
		}
	}
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

