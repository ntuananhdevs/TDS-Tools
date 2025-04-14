package traodoisub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tds/models"
)

type TDSClient struct {
	Token string
}

func NewClient(token string) *TDSClient {
	return &TDSClient{
		Token: token,
	}
}

// Lấy tất cả nhiệm vụ từ Traodoisub
func (c *TDSClient) GetAllJobs() (map[string][]models.Job, error) {
	jobTypes := []string{"like", "share", "follow"} // Các loại nhiệm vụ
	jobs := make(map[string][]models.Job)

	for _, jobType := range jobTypes {
		jobList, err := c.GetJob(jobType) // Sử dụng phương thức GetJob của TDSClient
		if err != nil {
			return nil, fmt.Errorf("error getting job %s: %v", jobType, err)
		}
		jobs[jobType] = jobList
	}

	return jobs, nil
}

// Lấy job theo loại (like, share, follow, etc.)
func (c *TDSClient) GetJob(jobType string) ([]models.Job, error) {
	// Tạo URL yêu cầu
	url := fmt.Sprintf("https://traodoisub.com/api/?fields=%s&access_token=%s", jobType, c.Token)

	// Gửi yêu cầu HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	// Giải mã JSON trả về thành các đối tượng Job
	var jobs []models.Job
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&jobs); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return jobs, nil
}
