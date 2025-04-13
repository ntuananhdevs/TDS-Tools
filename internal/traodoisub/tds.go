package traodoisub

import (
	"encoding/json"  // Import json package
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"tds/models"
)

type TDSClient struct {
	Token string
}

func NewClient(token string) *TDSClient {
	return &TDSClient{Token: token}
}

// Lấy tất cả nhiệm vụ từ Traodoisub (gọi từ API)
func (c *TDSClient) GetAllJobs() (map[string][]models.Job, error) {
	jobTypes := []string{"like", "likegiare", "likesieure", "page", "reaction", "comment", "share", "follow", "group"}

	jobs := make(map[string][]models.Job)

	// Lấy tất cả nhiệm vụ cho mỗi loại job
	for _, jobType := range jobTypes {
		jobList, err := c.GetJob(jobType)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy nhiệm vụ %s: %v", jobType, err)
		}
		jobs[jobType] = jobList
	}

	return jobs, nil
}

// Lấy job theo loại (like, comment, share,...)
func (c *TDSClient) GetJob(jobType string) ([]models.Job, error) {
	// Đảm bảo không có ký tự không hợp lệ trong token
	encodedToken := strings.TrimSpace(c.Token)  // Loại bỏ các ký tự thừa (như \r\n)

	// Tạo URL hợp lệ
	baseURL := "https://traodoisub.com/api/"
	apiURL := fmt.Sprintf("%s?fields=%s&access_token=%s", baseURL, jobType, url.QueryEscape(encodedToken))

	// Gửi yêu cầu HTTP để lấy dữ liệu nhiệm vụ
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Giải mã JSON từ response
	var jobs []models.Job
	_ = json.NewDecoder(resp.Body).Decode(&jobs)

	// Kiểm tra nếu không có nhiệm vụ
	if len(jobs) == 0 {
		return nil, fmt.Errorf("Không có nhiệm vụ cho %s", jobType)
	}

	// In chi tiết nhiệm vụ để kiểm tra
	for _, job := range jobs {
		fmt.Printf("Nhiệm vụ ID: %s - Link: %s - Msg: %s\n", job.ID, job.Link, job.Msg)
	}

	return jobs, nil
}

// Xác nhận job đã hoàn thành (sau khi làm nhiệm vụ)
func (c *TDSClient) ConfirmJob(jobType, id string) (string, int, error) {
	resp, err := http.Get(fmt.Sprintf("https://traodoisub.com/api/coin/?type=%s&id=%s&access_token=%s", jobType, id, c.Token))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Msg string `json:"msg"`
			Xu  int    `json:"xu"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	return result.Data.Msg, result.Data.Xu, nil
}

// Hàm bắt đầu chạy nhiệm vụ cho user, sử dụng token và lấy nhiệm vụ
func (c *TDSClient) Run(userID string) error {
	// Gọi API để bắt đầu nhiệm vụ cho user
	_, err := http.Get(fmt.Sprintf("https://traodoisub.com/api/?fields=run&id=%s&access_token=%s", userID, c.Token))
	return err
}
