package models

// TDSProfile chứa thông tin profile người dùng TDS
type TDSProfile struct {
	ID       string      `json:"id"`
	UserName string      `json:"user_name"`
	Xu       interface{} `json:"xu"`
	XuDie    interface{} `json:"xudie"`
}

// FacebookCookie chứa thông tin cookie của một tài khoản Facebook
type FacebookCookie struct {
	UID    string `json:"uid"` // Trường UID đã thêm vào
	Cookie string `json:"cookie"` // Trường Cookie
}

// JobInfo chứa thông tin về một nhiệm vụ từ TDS
type JobInfo struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	LinkPost  string `json:"linkpost"` // Remove this line
	Link      string `json:"link"`     // Add this line with json tag "link"
	Type      string `json:"type"`    // Loại nhiệm vụ: facebook_reaction, facebook_follow...
	SubType   string `json:"sub_type"` // Loại phụ: LIKE, LOVE, HAHA...
	JobStatus string `json:"job_status"`
	Reward    int    `json:"reward"`
	Code      string `json:"code"`
}

// AvailableJobs chứa danh sách các nhiệm vụ theo loại
type AvailableJobs struct {
	LikeJobs   []JobInfo `json:"like_jobs"`
	FollowJobs []JobInfo `json:"follow_jobs"`
	ShareJobs  []JobInfo `json:"share_jobs"`
	PageJobs   []JobInfo `json:"page_jobs"`
}
// ClaimCoinResponse chứa thông tin phản hồi khi nhận xu
type ClaimCoinResponse struct {
	Success int `json:"success"`
	Data    struct {
		Xu        int    `json:"xu"`
		JobSuccess int    `json:"job_success"`
		XuThem    int    `json:"xu_them"`
		Msg       string `json:"msg"`
	} `json:"data"`
	Msg string `json:"msg,omitempty"` // Thêm trường Msg ở ngoài data
}