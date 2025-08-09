package models

import "time"

// WalletTransaction 公用交易模型
type WalletTransaction struct {
	TransactionID string    `json:"transaction_id"`
	PlayerID      string    `json:"player_id"`
	Amount        float64   `json:"amount"`
	Direction     string    `json:"direction"` // "in" or "out"
	Timestamp     time.Time `json:"timestamp"`
}
