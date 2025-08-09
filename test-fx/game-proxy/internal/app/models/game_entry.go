package models

// GameEntry 第三方遊戲入口資訊
type GameEntry struct {
	URL     string `json:"url"`
	Expired int64  `json:"expired_unix"` // Unix timestamp
	Token   string `json:"token"`
}
