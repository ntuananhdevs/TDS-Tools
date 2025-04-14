package job

import (
	"fmt"
	"tds/config"
	"time"
	"math/rand"

	"tds/models"
)

// Hàm xử lý nhiệm vụ với các tham số cấu hình
func ExecuteTasks(config config.ToolConfig, jobs map[string][]models.Job) {
	taskCount := 0

	for {
		// Thực hiện nhiệm vụ
		taskCount++
		fmt.Printf("Đang thực hiện nhiệm vụ thứ %d...\n", taskCount)

		// Tạo delay ngẫu nhiên giữa các nhiệm vụ
		delay := randomDelay(config.DelayMin, config.DelayMax)
		fmt.Printf("Delay: %d ms\n", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		// Kiểm tra nếu cần chống block
		if taskCount%config.BlockAfter == 0 {
			fmt.Println("Chống block!")
		}

		// Kiểm tra nếu cần nghỉ ngơi
		if taskCount%config.RestAfter == 0 {
			fmt.Println("Nghỉ ngơi một lát...")
			time.Sleep(2 * time.Second) // Nghỉ ngơi 2 giây
		}

		// Kiểm tra nếu cần đổi nick
		if taskCount%config.ChangeNickAfter == 0 {
			fmt.Println("Đổi Nick!")
		}

		// Kiểm tra nếu cần xóa cookie
		if taskCount%config.DeleteCookieAfter == 0 {
			fmt.Println("Xóa cookie!")
		}

		// Thực hiện các nhiệm vụ chính
		if taskCount >= 10 { // Chỉ là ví dụ, có thể dừng sau 10 nhiệm vụ
			break
		}
	}
}

// Hàm tạo delay ngẫu nhiên giữa DelayMin và DelayMax
func randomDelay(min, max int) int {
	return min + rand.Intn(max-min+1) // Đây chỉ là ví dụ đơn giản
}
