// internal/job/executor.go

package job

import (
	"fmt"
	"sync"
	"time"

	"tds/config"
	"tds/internal/facebook"
	"tds/internal/traodoisub"
	"tds/models"
	"tds/utils"
)

// JobExecutor struct để quản lý việc thực thi nhiệm vụ
type JobExecutor struct {
	config    *config.Config
	logger    *utils.Logger
	tdsClient *traodoisub.TDSClient
	fbCookies []models.FacebookCookie // Danh sách cookie của các tài khoản Facebook
}

// NewJobExecutor tạo một JobExecutor mới
func NewJobExecutor(cfg *config.Config, logger *utils.Logger, tdsClient *traodoisub.TDSClient, fbCookies []models.FacebookCookie) *JobExecutor {
	return &JobExecutor{
		config:    cfg,
		logger:    logger,
		tdsClient: tdsClient,
		fbCookies: fbCookies,
	}
}

// ExecuteJobs sẽ xử lý việc lấy và thực hiện nhiệm vụ cho tất cả các tài khoản
func (e *JobExecutor) ExecuteJobs(jobs *models.AvailableJobs) {
	var wg sync.WaitGroup

	for _, fbCookie := range e.fbCookies {
		wg.Add(1) // Tăng bộ đếm cho mỗi Goroutine

		// Chạy một Goroutine cho mỗi tài khoản Facebook
		go func(cookie models.FacebookCookie) {
			defer wg.Done() // Giảm bộ đếm khi Goroutine hoàn thành

			e.logger.Info(fmt.Sprintf("Bắt đầu xử lý nhiệm vụ cho tài khoản Facebook UID: %s", cookie.UID))

			// Thực hiện nhiệm vụ FOLLOW
			if len(jobs.FollowJobs) > 0 {
				e.processFollowJobs(cookie, jobs.FollowJobs)
			}

			// Thực hiện nhiệm vụ SHARE
			if len(jobs.ShareJobs) > 0 {
				e.processShareJobs(cookie, jobs.ShareJobs)
			}

			// Thực hiện nhiệm vụ PAGE
			if len(jobs.PageJobs) > 0 {
				e.processPageJobs(cookie, jobs.PageJobs)
			}

			// Thực hiện nhiệm vụ LIKE
			if len(jobs.LikeJobs) > 0 {
				e.processLikeJobs(cookie, jobs.LikeJobs)
			}

			e.logger.Info(fmt.Sprintf("Hoàn thành xử lý nhiệm vụ cho tài khoản Facebook UID: %s", cookie.UID))

		}(fbCookie) // Truyền bản sao của cookie vào Goroutine
	}

	wg.Wait() // Chờ tất cả các Goroutine hoàn thành
	e.logger.Info("Tất cả các tài khoản đã hoàn thành xử lý nhiệm vụ.")
}

func (e *JobExecutor) processFollowJobs(cookie models.FacebookCookie, followJobs []models.JobInfo) {
	e.logger.Info(fmt.Sprintf("Tiến hành thực hiện và nhận xu cho %d nhiệm vụ FOLLOW (UID: %s)...", len(followJobs), cookie.UID))
	for _, job := range followJobs {
		taskSuccess := facebook.PerformFacebookTask("follow", job.ID, cookie.Cookie, e.logger, e.config)
		if taskSuccess {
			e.logger.Info(fmt.Sprintf("[UID: %s] Đã thực hiện (HTTP Request) theo dõi người dùng ID: %s", cookie.UID, job.ID))
			time.Sleep(e.config.GetRequestDelay())
			// Nhận xu sau khi theo dõi
			e.logger.Info(fmt.Sprintf("[UID: %s] Tiến hành nhận xu cho nhiệm vụ FOLLOW ID: %s...", cookie.UID, job.ID))
			claimResp, err := e.tdsClient.ClaimCoin("facebook_follow", "facebook_api")
			if err != nil {
				e.logger.Error(fmt.Sprintf("[UID: %s] Lỗi khi nhận xu FOLLOW ID: %s: %s", cookie.UID, job.ID, err.Error()))
			} else {
				e.logger.Info(fmt.Sprintf("[UID: %s] Nhận xu FOLLOW ID: %s thành công. %s", cookie.UID, job.ID, claimResp.Data.Msg))
			}
		} else {
			e.logger.Warn(fmt.Sprintf("[UID: %s] Thực hiện (HTTP Request) theo dõi người dùng ID: %s thất bại.", cookie.UID, job.ID))
		}
		time.Sleep(e.config.GetRequestDelay())
	}
}

func (e *JobExecutor) processShareJobs(cookie models.FacebookCookie, shareJobs []models.JobInfo) {
	e.logger.Info(fmt.Sprintf("Tiến hành thực hiện và nhận xu cho %d nhiệm vụ SHARE (UID: %s)...", len(shareJobs), cookie.UID))
	for _, job := range shareJobs {
		taskSuccess := facebook.PerformFacebookTask("share", job.LinkPost, cookie.Cookie, e.logger, e.config)
		if taskSuccess {
			e.logger.Info(fmt.Sprintf("[UID: %s] Đã thực hiện (HTTP Request) chia sẻ bài viết ID: %s", cookie.UID, job.ID))
			time.Sleep(e.config.GetRequestDelay())
			claimResp, err := e.tdsClient.ClaimCoin("facebook_share", job.Code)
			if err != nil {
				e.logger.Error(fmt.Sprintf("[UID: %s] Lỗi khi nhận xu SHARE (ID: %s, Code: %s): %s", cookie.UID, job.ID, job.Code, err.Error()))
			} else {
				e.logger.Info(fmt.Sprintf("[UID: %s] Nhận xu SHARE (ID: %s, Code: %s) thành công. %s", cookie.UID, job.ID, job.Code, claimResp.Data.Msg))
			}
		} else {
			e.logger.Warn(fmt.Sprintf("[UID: %s] Thực hiện (HTTP Request) chia sẻ bài viết ID: %s thất bại.", cookie.UID, job.ID))
		}
		time.Sleep(e.config.GetRequestDelay())
	}
}

func (e *JobExecutor) processPageJobs(cookie models.FacebookCookie, pageJobs []models.JobInfo) {
	e.logger.Info(fmt.Sprintf("Tiến hành thực hiện và nhận xu cho %d nhiệm vụ LIKE PAGE (UID: %s)...", len(pageJobs), cookie.UID))
	for _, job := range pageJobs {
		taskSuccess := facebook.PerformFacebookTask("like_page", job.ID, cookie.Cookie, e.logger, e.config)
		if taskSuccess {
			e.logger.Info(fmt.Sprintf("[UID: %s] Đã thực hiện (HTTP Request) thích trang ID: %s", cookie.UID, job.ID))
			time.Sleep(e.config.GetRequestDelay())
			claimResp, err := e.tdsClient.ClaimCoin("facebook_page", "facebook_api")
			if err != nil {
				e.logger.Error(fmt.Sprintf("[UID: %s] Lỗi khi nhận xu LIKE PAGE ID: %s: %s", cookie.UID, job.ID, err.Error()))
			} else {
				e.logger.Info(fmt.Sprintf("[UID: %s] Nhận xu LIKE PAGE ID: %s thành công. %s", cookie.UID, job.ID, claimResp.Data.Msg))
			}
		} else {
			e.logger.Warn(fmt.Sprintf("[UID: %s] Thực hiện (HTTP Request) thích trang ID: %s thất bại.", cookie.UID, job.ID))
		}
		time.Sleep(e.config.GetRequestDelay())
	}
}

func (e *JobExecutor) processLikeJobs(cookie models.FacebookCookie, likeJobs []models.JobInfo) {
	e.logger.Info(fmt.Sprintf("Tiến hành thực hiện và nhận xu cho %d nhiệm vụ LIKE (UID: %s)...", len(likeJobs), cookie.UID))
	for _, job := range likeJobs {
		taskSuccess := facebook.PerformFacebookTask("like", job.LinkPost, cookie.Cookie, e.logger, e.config)
		if taskSuccess {
			e.logger.Info(fmt.Sprintf("[UID: %s] Đã thực hiện (HTTP Request) thích bài viết ID: %s", cookie.UID, job.ID))
			time.Sleep(e.config.GetRequestDelay())
			claimResp, err := e.tdsClient.ClaimCoin("facebook_reaction", job.Code)
			if err != nil {
				e.logger.Error(fmt.Sprintf("[UID: %s] Lỗi khi nhận xu LIKE (ID: %s, Code: %s): %s", cookie.UID, job.ID, job.Code, err.Error()))
			} else {
				e.logger.Info(fmt.Sprintf("[UID: %s] Nhận xu LIKE (ID: %s, Code: %s) thành công. %s", cookie.UID, job.ID, job.Code, claimResp.Data.Msg))
			}
		} else {
			e.logger.Warn(fmt.Sprintf("[UID: %s] Thực hiện (HTTP Request) thích bài viết ID: %s thất bại.", cookie.UID, job.ID))
		}
		time.Sleep(e.config.GetRequestDelay())
	}
}