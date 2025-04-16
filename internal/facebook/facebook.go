package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tds/config"
	"tds/utils"
)

// PageLikeResponse struct để chứa response khi thích trang
type PageLikeResponse struct {
	Data struct {
		PageLike struct {
			Page struct {
				ProfilePlusForDelegatePage struct {
					ProfileAction struct {
						Typename string `json:"__typename"`
						IsProfileAction string `json:"__isProfileAction"`
						IconImage struct {
							Height int    `json:"height"`
							Scale  int    `json:"scale"`
							URI    string `json:"uri"`
							Width  int    `json:"width"`
						} `json:"icon_image"`
						ID    string `json:"id"`
						Title struct {
							Text string `json:"text"`
						} `json:"title"`
					} `json:"profile_action"`
					FollowingStatus struct {
						Typename      string `json:"__typename"`
						IsProfileAction string `json:"__isProfileAction"`
						IconImage       struct {
							Height int    `json:"height"`
							Scale  int    `json:"scale"`
							URI    string `json:"uri"`
							Width  int    `json:"width"`
						} `json:"icon_image"`
						ID        string `json:"id"`
						Title     struct {
							Text string `json:"text"`
						} `json:"title"`
						IsActive  bool `json:"is_active"`
					} `json:"following_status"`
					ID string `json:"id"`
				} `json:"profile_plus_for_delegate_page"`
				IsViewerFan    bool   `json:"is_viewer_fan"`
				SubscribeStatus string `json:"subscribe_status"`
				ID              string `json:"id"`
			} `json:"page"`
		} `json:"page_like"`
	} `json:"data"`
	Extensions struct {
		IsFinal bool `json:"is_final"`
	} `json:"extensions"`
}

func PerformFacebookTask(taskType string, taskTarget string, cookieString string, logger *utils.Logger, cfg *config.Config) bool {
	logger.Info(fmt.Sprintf("[HTTP Request] Thực hiện '%s' với target '%s'", taskType, taskTarget))

	client := &http.Client{}
	var req *http.Request
	var err error
	var body string // Biến để chứa body cho các request POST

	headers := map[string]string{
		"Cookie":    cookieString,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36", // Quan trọng
		// Thêm các headers chung khác nếu cần
	}

	switch taskType {
	case "like":
		// **NGHIÊN CỨU:** Tìm URL và dữ liệu POST cần thiết để thích một bài viết có ID là taskTarget
		// Đây thường là một request POST
		body = fmt.Sprintf("fb_dtsg=YOUR_FB_DTSG&jazoest=YOUR_JAZOEST&post_id=%s&__user=%s&__a=1&__dyn=...", taskTarget, strings.SplitN(cookieString, "c_user=", 2)[1][:strings.Index(strings.SplitN(cookieString, "c_user=", 2)[1], ";")]) // Cần điền thông tin thật
		req, err = http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(body)) // Thay ENDPOINT thật
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	case "follow":
		followURL := fmt.Sprintf("https://m.facebook.com/%s", taskTarget) // Ví dụ URL
		// **NGHIÊN CỨU:** Tìm URL và dữ liệu POST cần thiết để theo dõi người dùng
		req, err = http.NewRequest("GET", followURL, nil) // Có thể là POST
		// Có thể cần thêm action vào URL hoặc body
	case "share":
		// **NGHIÊN CỨU:** Việc share phức tạp hơn. Cần xác định URL và các tham số để share một bài viết có link là taskTarget
		req, err = http.NewRequest("POST", "FACEBOOK_SHARE_ENDPOINT", nil) // Rất có thể là POST với dữ liệu
		headers["Content-Type"] = "application/x-www-form-urlencoded"
		// Cần xây dựng body chứa thông tin link share
	case "like_page":
		body = fmt.Sprintf("fb_dtsg=YOUR_FB_DTSG&jazoest=YOUR_JAZOEST&page_id=%s&__user=%s&__a=1&__dyn=...", taskTarget, strings.SplitN(cookieString, "c_user=", 2)[1][:strings.Index(strings.SplitN(cookieString, "c_user=", 2)[1], ";")]) // Cần điền thông tin thật
		req, err = http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(body)) // Thay ENDPOINT thật
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	default:
		logger.Warn(fmt.Sprintf("Loại nhiệm vụ Facebook không được hỗ trợ: %s", taskType))
		return false
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Lỗi tạo request Facebook (%s): %v", taskType, err))
		return false
	}

	// Thêm headers vào request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Lỗi gửi request Facebook (%s) đến %s: %v", taskType, req.URL.String(), err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Info(fmt.Sprintf("[HTTP Response] %s - Status: %d", taskType, resp.StatusCode))
		if taskType == "like_page" {
			var pageLikeResponse PageLikeResponse
			err = json.NewDecoder(resp.Body).Decode(&pageLikeResponse)
			if err == nil && pageLikeResponse.Data.PageLike.Page.IsViewerFan {
				return true
			} else {
				logger.Warn(fmt.Sprintf("Thích trang không thành công hoặc lỗi giải mã response: %v", err))
				return false
			}
		}
		time.Sleep(time.Second * 2) // Mô phỏng thời gian xử lý cho các loại khác
		return true
	} else {
		logger.Warn(fmt.Sprintf("[HTTP Response] %s - Status: %d", taskType, resp.StatusCode))
		return false
	}
}