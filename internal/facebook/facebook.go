package facebook

import (
	"fmt"
	"net/http"
)

// Hàm Like một bài viết
func Like(jobID string) error {
	url := fmt.Sprintf("https://facebook.com/like?jobID=%s", jobID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error liking post: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã phản hồi
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to like post, status code: %d", resp.StatusCode)
	}

	return nil
}
