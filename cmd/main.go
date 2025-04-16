package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"tds/config"
	"tds/internal/facebook"
	"tds/internal/traodoisub"
	"tds/models"
	"tds/utils"
)

func main() {
	// Khởi tạo random seed
	rand.Seed(time.Now().UnixNano())

	// Khởi tạo logger
	logger := utils.NewLogger("info")
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
	var fbUID string
	if len(fbCookies) > 0 {
		fbUID = fbCookies[0].UID
		if fbUID != "" {
			logger.Info(fmt.Sprintf("UID của tài khoản Facebook: %s", fbUID))
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
	if fbUID != "" {
		err = tdsClient.ConfigureAccount(fbUID)
		if err != nil {
			logger.Error(fmt.Sprintf("Lỗi khi cấu hình tài khoản Facebook: %s", err.Error()))
			return
		}
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

	// Xử lý nhiệm vụ
	processTasks(cfg, logger, fbCookies[0], tdsClient, jobs)

	// Hiển thị tổng số nhiệm vụ
	totalJobs := len(jobs.LikeJobs) + len(jobs.FollowJobs) + len(jobs.ShareJobs) + len(jobs.PageJobs)
	logger.Info(fmt.Sprintf("Tổng cộng: %d nhiệm vụ", totalJobs))
}

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

func displayJobs(jobType string, jobs []models.JobInfo, logger *utils.Logger) {
	if len(jobs) == 0 {
		logger.Info(fmt.Sprintf("Không có nhiệm vụ %s", jobType))
		return
	}

	logger.Info(fmt.Sprintf("=== %s (%d nhiệm vụ) ===", jobType, len(jobs)))
	for i, job := range jobs {
		var codeStr string
		if job.Code != "" {
			codeStr = fmt.Sprintf(", Code: %s", job.Code)
		}
		var linkStr string
		if job.Link != "" {
			linkStr = fmt.Sprintf(", Link: %s", job.Link)
		}
		logger.Info(fmt.Sprintf("%d. ID: %s%s%s, Xu: %d", i+1, job.ID, codeStr, linkStr, job.Reward))
	}
}

func processTasks(cfg *config.Config, logger *utils.Logger, fbCookie models.FacebookCookie, tdsClient *traodoisub.TDSClient, jobs *models.AvailableJobs) {
	// Xử lý nhiệm vụ FOLLOW
	if len(jobs.FollowJobs) > 0 {
		logger.Info("Tiến hành thực hiện và nhận xu cho nhiệm vụ FOLLOW...")
		for _, job := range jobs.FollowJobs {
			taskSuccess := facebook.PerformFacebookTask("follow", job.ID, fbCookie.Cookie, logger, cfg)
			if taskSuccess {
				logger.Info(fmt.Sprintf("Đã thực hiện (HTTP Request) theo dõi người dùng ID: %s", job.ID))
				time.Sleep(cfg.GetRequestDelay())
				// Nhận xu sau khi theo dõi
				logger.Info(fmt.Sprintf("Tiến hành nhận xu cho nhiệm vụ FOLLOW ID: %s...", job.ID))
				claimResp, err := tdsClient.ClaimCoin("facebook_follow", "facebook_api")
				if err != nil {
					logger.Error(fmt.Sprintf("Lỗi khi nhận xu FOLLOW ID: %s: %s", job.ID, err.Error()))
				} else {
					logger.Info(fmt.Sprintf("Nhận xu FOLLOW ID: %s thành công. %s", job.ID, claimResp.Data.Msg))
				}
			} else {
				logger.Warn(fmt.Sprintf("Thực hiện (HTTP Request) theo dõi người dùng ID: %s thất bại.", job.ID))
			}
			time.Sleep(cfg.GetRequestDelay())
		}
	}

	// Xử lý nhiệm vụ SHARE
	if len(jobs.ShareJobs) > 0 {
		logger.Info("Tiến hành thực hiện và nhận xu cho nhiệm vụ SHARE...")
		for _, job := range jobs.ShareJobs {
			taskSuccess := facebook.PerformFacebookTask("share", job.Link, fbCookie.Cookie, logger, cfg)
			if taskSuccess {
				logger.Info(fmt.Sprintf("Đã thực hiện (HTTP Request) chia sẻ bài viết ID: %s", job.ID))
				time.Sleep(cfg.GetRequestDelay())
				claimResp, err := tdsClient.ClaimCoin("facebook_share", job.Code)
				if err != nil {
					logger.Error(fmt.Sprintf("Lỗi khi nhận xu SHARE (ID: %s, Code: %s): %s", job.ID, job.Code, err.Error()))
				} else {
					logger.Info(fmt.Sprintf("Nhận xu SHARE (ID: %s, Code: %s) thành công. %s", job.ID, job.Code, claimResp.Data.Msg))
				}
			} else {
				logger.Warn(fmt.Sprintf("Thực hiện (HTTP Request) chia sẻ bài viết ID: %s thất bại.", job.ID))
			}
			time.Sleep(cfg.GetRequestDelay())
		}
	}

	// Xử lý nhiệm vụ PAGE
	if len(jobs.PageJobs) > 0 {
		logger.Info("Tiến hành thực hiện và nhận xu cho nhiệm vụ PAGE...")
		for _, job := range jobs.PageJobs {
			taskSuccess := facebook.PerformFacebookTask("like_page", job.ID, fbCookie.Cookie, logger, cfg)
			if taskSuccess {
				logger.Info(fmt.Sprintf("Đã thực hiện (HTTP Request) thích trang ID: %s", job.ID))
				time.Sleep(cfg.GetRequestDelay())
				claimResp, err := tdsClient.ClaimCoin("facebook_page", "facebook_api")
				if err != nil {
					logger.Error(fmt.Sprintf("Lỗi khi nhận xu PAGE ID: %s: %s", job.ID, err.Error()))
				} else {
					logger.Info(fmt.Sprintf("Nhận xu PAGE ID: %s thành công. %s", job.ID, claimResp.Data.Msg))
				}
			} else {
				logger.Warn(fmt.Sprintf("Thực hiện (HTTP Request) thích trang ID: %s thất bại.", job.ID))
			}
			time.Sleep(cfg.GetRequestDelay())
		}
	}

	// Xử lý nhiệm vụ LIKE
	if len(jobs.LikeJobs) > 0 {
		logger.Info("Tiến hành thực hiện và nhận xu cho nhiệm vụ LIKE...")
		for _, job := range jobs.LikeJobs {
			taskSuccess := facebook.PerformFacebookTask("like", job.Link, fbCookie.Cookie, logger, cfg)
			if taskSuccess {
				logger.Info(fmt.Sprintf("Đã thực hiện (HTTP Request) thích bài viết ID: %s", job.ID))
				time.Sleep(cfg.GetRequestDelay())
				claimResp, err := tdsClient.ClaimCoin("facebook_reaction", job.Code)
				if err != nil {
					logger.Error(fmt.Sprintf("Lỗi khi nhận xu LIKE (ID: %s, Code: %s): %s", job.ID, job.Code, err.Error()))
				} else {
					logger.Info(fmt.Sprintf("Nhận xu LIKE (ID: %s, Code: %s) thành công. %s", job.ID, job.Code, claimResp.Data.Msg))
				}
			} else {
				logger.Warn(fmt.Sprintf("Thực hiện (HTTP Request) thích bài viết ID: %s thất bại.", job.ID))
			}
			time.Sleep(cfg.GetRequestDelay())
		}
	}
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}