package traodoisub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"tds/config"
	"tds/models"
	"tds/utils"
)

const (
	BaseURL = "https://traodoisub.com/api/"
)

// TDSClient đại diện cho một client tương tác với API TraDoiSub
type TDSClient struct {
	Token      string
	HTTPClient *http.Client
	Logger     *utils.Logger
	UserID     string // ID Facebook hiện đang cấu hình
}

// NewTDSClient tạo một client mới để tương tác với TDS API
func NewTDSClient(cfg *config.Config, logger *utils.Logger) *TDSClient {
	return &TDSClient{
		Token: cfg.GetTDSToken(),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Logger: logger,
	}
}

// ConfigureAccount cấu hình tài khoản Facebook cho TDS
func (c *TDSClient) ConfigureAccount(facebookID string) error {
	c.Logger.Info(fmt.Sprintf("Đang cấu hình tài khoản Facebook ID: %s", facebookID))
	
	reqURL := BaseURL
	
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}
	
	// Thiết lập query params
	q := url.Values{}
	q.Add("fields", "run")
	q.Add("id", facebookID)
	q.Add("access_token", c.Token)
	req.URL.RawQuery = q.Encode()
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	// Log response để debug
	c.Logger.Debug(fmt.Sprintf("API Response: %s", string(body)))
	
	// Kiểm tra xem response có phải HTML không
	if len(body) > 0 && body[0] == '<' {
		return errors.New("API trả về HTML thay vì JSON, có thể token không hợp lệ hoặc URL không chính xác")
	}
	
	var configResp struct {
		Success int `json:"success"`
		Data    struct {
			ID  string `json:"id"`
			Msg string `json:"msg"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}
	
	if err := json.Unmarshal(body, &configResp); err != nil {
		return err
	}
	
	if configResp.Success != 200 {
		if configResp.Error != "" {
			return errors.New(configResp.Error)
		}
		return errors.New("không thể cấu hình tài khoản")
	}
	
	c.UserID = facebookID
	c.Logger.Info(fmt.Sprintf("Cấu hình thành công: %s", configResp.Data.Msg))
	
	return nil
}

// GetProfile lấy thông tin profile của người dùng từ TDS
func (c *TDSClient) GetProfile() (*models.TDSProfile, error) {
	c.Logger.Info("Đang lấy thông tin profile từ TDS...")
	
	reqURL := BaseURL
	
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	
	q := url.Values{}
	q.Add("fields", "profile")
	q.Add("access_token", c.Token)
	req.URL.RawQuery = q.Encode()
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Log response để debug
	c.Logger.Debug(fmt.Sprintf("API Response: %s", string(body)))
	
	// Kiểm tra xem response có phải HTML không
	if len(body) > 0 && body[0] == '<' {
		return nil, errors.New("API trả về HTML thay vì JSON, có thể token không hợp lệ hoặc URL không chính xác")
	}
	
	var profileResp struct {
		Success int               `json:"success"`
		Data    models.TDSProfile `json:"data"`
		Error   string            `json:"error,omitempty"`
	}
	
	if err := json.Unmarshal(body, &profileResp); err != nil {
		return nil, err
	}
	
	if xuStr, ok := profileResp.Data.Xu.(string); ok {
		if xu, err := strconv.Atoi(xuStr); err == nil {
			profileResp.Data.Xu = xu
		}
	}
	
	if xuDieStr, ok := profileResp.Data.XuDie.(string); ok {
		if xuDie, err := strconv.Atoi(xuDieStr); err == nil {
			profileResp.Data.XuDie = xuDie
		}
	}
	
	if profileResp.Success != 200 {
		if profileResp.Error != "" {
			return nil, errors.New(profileResp.Error)
		}
		return nil, errors.New("không thể lấy thông tin profile")
	}
	
	return &profileResp.Data, nil
}

// GetReactionJobs lấy danh sách các nhiệm vụ reaction từ TDS
func (c *TDSClient) GetReactionJobs(reactionType string) ([]models.JobInfo, error) {
	return c.getJobsByField("facebook_reaction", reactionType)
}

// GetFollowJobs lấy danh sách các nhiệm vụ follow từ TDS
func (c *TDSClient) GetFollowJobs() ([]models.JobInfo, error) {
	return c.getJobsByField("facebook_follow", "")
}

// GetShareJobs lấy danh sách các nhiệm vụ share từ TDS
func (c *TDSClient) GetShareJobs() ([]models.JobInfo, error) {
	return c.getJobsByField("facebook_share", "")
}

// GetPageJobs lấy danh sách các nhiệm vụ like page từ TDS
func (c *TDSClient) GetPageJobs() ([]models.JobInfo, error) {
	return c.getJobsByField("facebook_page", "")
}

// getJobsByField là hàm helper để lấy các nhiệm vụ theo loại
func (c *TDSClient) getJobsByField(field, jobType string) ([]models.JobInfo, error) {
	jobName := field
	if jobType != "" {
		jobName = fmt.Sprintf("%s (%s)", field, jobType)
	}
	c.Logger.Info(fmt.Sprintf("Đang lấy nhiệm vụ %s từ TDS...", jobName))
	
	if c.UserID == "" {
		return nil, errors.New("cần cấu hình tài khoản Facebook trước khi lấy nhiệm vụ")
	}
	
	reqURL := BaseURL
	
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	
	// Thiết lập query params
	q := url.Values{}
	q.Add("fields", field)
	q.Add("access_token", c.Token)
	if jobType != "" {
		q.Add("type", jobType)
	}
	req.URL.RawQuery = q.Encode()
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Log response để debug
	c.Logger.Debug(fmt.Sprintf("API Response: %s", string(body)))
	
	// Kiểm tra xem response có phải HTML không
	if len(body) > 0 && body[0] == '<' {
		return nil, errors.New("API trả về HTML thay vì JSON, có thể token không hợp lệ hoặc URL không chính xác")
	}
	
	var jobs []models.JobInfo
	if err := json.Unmarshal(body, &jobs); err != nil {
		// Kiểm tra nếu phản hồi là thông báo lỗi
		var errorResp struct {
			Success int    `json:"success"`
			Error   string `json:"error"`
		}
		
		if jsonErr := json.Unmarshal(body, &errorResp); jsonErr == nil && errorResp.Success != 200 {
			if errorResp.Error != "" {
				return nil, errors.New(errorResp.Error)
			}
			return nil, errors.New("không thể lấy danh sách nhiệm vụ hoặc danh sách rỗng")
		}
		
		return nil, err
	}
	
	// Gắn loại nhiệm vụ vào mỗi job
	for i := range jobs {
		jobs[i].Type = field
		if jobType != "" {
			jobs[i].SubType = jobType
		}
	}
	
	return jobs, nil
}

// GetAllJobs lấy tất cả các loại nhiệm vụ từ TDS
func (c *TDSClient) GetAllJobs() (*models.AvailableJobs, error) {
	c.Logger.Info("Đang lấy tất cả các nhiệm vụ từ TDS...")
	
	if c.UserID == "" {
		return nil, errors.New("cần cấu hình tài khoản Facebook trước khi lấy nhiệm vụ")
	}
	
	likeJobs, err := c.GetReactionJobs("LIKE")
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("Lỗi khi lấy nhiệm vụ LIKE: %s", err.Error()))
		likeJobs = []models.JobInfo{}
	}
	
	time.Sleep(1 * time.Second) // Tránh rate limit
	
	followJobs, err := c.GetFollowJobs()
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("Lỗi khi lấy nhiệm vụ FOLLOW: %s", err.Error()))
		followJobs = []models.JobInfo{}
	}
	
	time.Sleep(1 * time.Second) // Tránh rate limit
	
	shareJobs, err := c.GetShareJobs()
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("Lỗi khi lấy nhiệm vụ SHARE: %s", err.Error()))
		shareJobs = []models.JobInfo{}
	}
	
	time.Sleep(1 * time.Second) // Tránh rate limit
	
	pageJobs, err := c.GetPageJobs()
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("Lỗi khi lấy nhiệm vụ PAGE: %s", err.Error()))
		pageJobs = []models.JobInfo{}
	}
	
	return &models.AvailableJobs{
		LikeJobs:   likeJobs,
		FollowJobs: followJobs,
		ShareJobs:  shareJobs,
		PageJobs:   pageJobs,
	}, nil
}