package job

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"tds/config"
	"tds/internal/facebook"
	"tds/internal/traodoisub"
	"tds/models"
	"tds/utils"
)

func Run(cookies []models.CookieUser) {
	tds := traodoisub.NewClient(config.AccessTokenTDS)
	var wg sync.WaitGroup

	// Xử lý tài khoản Facebook đồng thời
	for _, user := range cookies {
		wg.Add(1)
		go func(user models.CookieUser) {
			defer wg.Done()
			handleUser(user, tds)
		}(user)
	}

	wg.Wait() // Chờ tất cả goroutines hoàn thành
}

func handleUser(user models.CookieUser, tds *traodoisub.TDSClient) {
	fb := facebook.NewFacebookClient(user.Cookie)
	_ = tds.Run(user.UserID) // Bắt đầu nhiệm vụ cho user
	utils.Info("🧩 Bắt đầu với: " + user.Name)
	jobCount := 0

	// Lặp mãi cho mỗi nhiệm vụ
	for {
		jobs, err := tds.GetJob("like")
		if err != nil || len(jobs) == 0 {
			utils.Warning("Không có job cho " + user.Name)
			time.Sleep(10 * time.Second)
			continue
		}
		for _, j := range jobs {
			if err := fb.Like(j.ID); err == nil {
				msg, xu, _ := tds.ConfirmJob("like", j.ID)
				utils.Success(fmt.Sprintf("%s LIKE %s | %s | %d xu", user.Name, j.ID, msg, xu))
				jobCount++
			} else {
				utils.Warning(fmt.Sprintf("%s lỗi LIKE: %s", user.Name, j.ID))
			}

			time.Sleep(time.Duration(rand.Intn(3)+2) * time.Second)

			// Nghỉ sau mỗi 10 job
			if jobCount > 0 && jobCount%config.JobsBeforeRest == 0 {
				utils.Info(fmt.Sprintf("⏸ %s nghỉ %v...", user.Name, config.RestDuration))
				time.Sleep(config.RestDuration)
			}
		}
	}
}
