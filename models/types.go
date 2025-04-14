package models

type Job struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Link string `json:"link"`
	Msg  string `json:"msg"`
}

type CookieUser struct {
	UserID string `json:"user_id"`
	Cookie string `json:"cookie"`
}

// TDSProfile lưu trữ thông tin của TDS
type TDSProfile struct {
	AccessToken string `json:"access_token"`
}
