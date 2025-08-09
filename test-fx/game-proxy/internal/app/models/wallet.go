package models

// Wallet 範例錢包結構
type Wallet struct {
	PlayerID string  `json:"player_id"`
	Balance  float64 `json:"balance"`
}
