package interfaces

import (
	"context"
	"idv/chris/internal/app/models"
)

// VendorGameService 定義所有第三方遊戲廠商服務需實作的方法
type VendorGameService interface {
	GetName() string
	RegisterOrLogin(ctx context.Context, playerID string) (*models.Player, error)
	TransferOut(ctx context.Context, playerID string, amount float64) (*models.Wallet, error)
	TransferIn(ctx context.Context, playerID string, amount float64) (*models.Wallet, error)
	GetGameEntry(ctx context.Context) (*models.GameEntry, error)
	Logout(ctx context.Context, playerID string) (*models.LogoutResult, error)
}
