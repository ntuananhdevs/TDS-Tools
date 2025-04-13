package models

type Job struct {
	ID   string `json:"id"`
	Link string `json:"link"`
	Msg  string `json:"msg,omitempty"`
	Type string `json:"type,omitempty"`
}

type CookieUser struct {
	Cookie string `json:"cookie"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

type TDSProfile struct {
	User  string `json:"user"`
	Xu    int    `json:"xu"`
	Xudie int    `json:"xudie"`
}
