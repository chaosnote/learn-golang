package models

// LogoutResult 登出結果
type LogoutResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
