package main

import (
	"encoding/json"
	"fmt"
	"os"

	"tds/config"
	"tds/internal/traodoisub"
	"tds/models"
	"tds/utils"
)

func main() {
	// Khởi tạo logger
	logger := utils.NewLogger("debug") // Đổi thành "debug" để xem chi tiết API response
	logger.Info("Khởi động công cụ TraDoiSub...")

	// Tải cấu hình từ file .env
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Lỗi khi tải cấu hình: %s", err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("Đã tải cấu hình. Token: %s...", cfg.GetTDSToken()[:10]))

	// Đọc cookie Facebook từ file
	logger.Info("Đang đọc cookie Facebook từ file...")
	fbCookies, err := loadFacebookCookies("../cookie.json")
	if err != nil {
		logger.Error(fmt.Sprintf("Lỗi khi đọc cookie: %s", err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("Đã đọc %d tài khoản Facebook từ file cookie", len(fbCookies)))

	// Trích xuất UID từ tài khoản Facebook đầu tiên
	if len(fbCookies) > 0 {
		uid := fbCookies[0].UID // Lấy UID trực tiếp từ cấu trúc cookie
		if uid != "" {
			logger.Info(fmt.Sprintf("UID của tài khoản Facebook: %s", uid))
		} else {
			logger.Error("Không tìm thấy UID trong cookie của tài khoản Facebook.")
			return
		}
	} else {
		logger.Error("Không có tài khoản Facebook nào trong file cookie")
		return
	}

	// Khởi tạo TDS client
	tdsClient := traodoisub.NewTDSClient(cfg, logger)

	// Lấy thông tin profile TDS
	profile, err := tdsClient.GetProfile()
	if err != nil {
		logger.Error(fmt.Sprintf("Lỗi khi lấy thông tin profile: %s", err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("Xin chào %s! Số xu hiện tại: %d", profile.UserName, profile.Xu))

	// Cấu hình tài khoản Facebook đầu tiên
	if len(fbCookies) > 0 {
		err = tdsClient.ConfigureAccount(fbCookies[0].UID) // Truyền UID để cấu hình tài khoản Facebook
		if err != nil {
			logger.Error(fmt.Sprintf("Lỗi khi cấu hình tài khoản Facebook: %s", err.Error()))
			return
		}
	} else {
		logger.Error("Không có tài khoản Facebook nào trong file cookie")
		return
	}

	// Lấy danh sách nhiệm vụ
	logger.Info("Đang lấy danh sách nhiệm vụ từ TraDoiSub...")
	jobs, err := tdsClient.GetAllJobs()
	if err != nil {
		logger.Error(fmt.Sprintf("Lỗi khi lấy danh sách nhiệm vụ: %s", err.Error()))
		return
	}

	// Hiển thị danh sách nhiệm vụ
	displayJobs("LIKE", jobs.LikeJobs, logger)
	displayJobs("FOLLOW", jobs.FollowJobs, logger)
	displayJobs("SHARE", jobs.ShareJobs, logger)
	displayJobs("PAGE", jobs.PageJobs, logger)

	// Hiển thị tổng số nhiệm vụ
	totalJobs := len(jobs.LikeJobs) + len(jobs.FollowJobs) + len(jobs.ShareJobs) + len(jobs.PageJobs)
	logger.Info(fmt.Sprintf("Tổng cộng: %d nhiệm vụ", totalJobs))
}

// loadFacebookCookies đọc cookie Facebook từ file JSON
func loadFacebookCookies(filePath string) ([]models.FacebookCookie, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cookies []models.FacebookCookie
	err = json.Unmarshal(data, &cookies)
	if err != nil {
		return nil, err
	}

	return cookies, nil
}

// displayJobs hiển thị danh sách nhiệm vụ theo loại
func displayJobs(jobType string, jobs []models.JobInfo, logger *utils.Logger) {
	if len(jobs) == 0 {
		logger.Info(fmt.Sprintf("Không có nhiệm vụ %s", jobType))
		return
	}

	logger.Info(fmt.Sprintf("=== %s (%d nhiệm vụ) ===", jobType, len(jobs)))
	for i, job := range jobs {
		// Kiểm tra trạng thái nhiệm vụ
		if job.JobStatus != "ACTIVE" {
			logger.Warn(fmt.Sprintf("Nhiệm vụ %s có trạng thái không hợp lệ: %s", jobType, job.JobStatus))
			continue
		}

		// Hiển thị nhiệm vụ nếu hợp lệ
		logger.Info(fmt.Sprintf("%d. ID: %s, Link: %s, Xu: %d", i+1, job.ID, job.LinkPost, job.Reward))
	}
}
