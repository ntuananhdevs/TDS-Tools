package models

// TDSProfile chứa thông tin profile người dùng TDS
type TDSProfile struct {
	ID       string      `json:"id"`
	UserName string      `json:"user_name"`
	Xu       interface{} `json:"xu"`
	XuDie    interface{}         `json:"xudie"`
}

// FacebookCookie chứa thông tin cookie của một tài khoản Facebook
type FacebookCookie struct {
	UID     string `json:"uid"`     // Trường UID đã thêm vào
	Cookie string `json:"cookie"`  // Trường Cookie
}
// JobInfo chứa thông tin về một nhiệm vụ từ TDS
type JobInfo struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	LinkPost  string `json:"linkpost"`
	Type      string `json:"type"`      // Loại nhiệm vụ: facebook_reaction, facebook_follow...
	SubType   string `json:"sub_type"`  // Loại phụ: LIKE, LOVE, HAHA...
	JobStatus string `json:"job_status"`
	Reward    int    `json:"reward"`
}

// AvailableJobs chứa danh sách các nhiệm vụ theo loại
type AvailableJobs struct {
	LikeJobs   []JobInfo `json:"like_jobs"`
	FollowJobs []JobInfo `json:"follow_jobs"`
	ShareJobs  []JobInfo `json:"share_jobs"`
	PageJobs   []JobInfo `json:"page_jobs"`
}
