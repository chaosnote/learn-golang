package models

// LogoutResult 玩家登出第三方遊戲結果
type LogoutResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
