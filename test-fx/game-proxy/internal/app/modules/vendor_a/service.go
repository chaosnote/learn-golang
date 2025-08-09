package vendor_a

import (
	"context"
	"idv/chris/internal/app/config"
	"idv/chris/internal/app/models"
	"time"

	"go.uber.org/zap"
)

// VendorAService 實作 interfaces.VendorGameService
type VendorAService struct {
	cfg    config.VendorConfig
	client *APIClient
	logger *zap.Logger
}

func NewVendorAService(cfg config.VendorConfig, client *APIClient, logger *zap.Logger) *VendorAService {
	return &VendorAService{
		cfg:    cfg,
		client: client,
		logger: logger,
	}
}

// GetName 回傳供 VendorManager 索引的名稱（須與 config 裡的 key 一致）
func (s *VendorAService) GetName() string {
	return s.cfg.Name
}

func (s *VendorAService) RegisterOrLogin(ctx context.Context, playerID string) (*models.Player, error) {
	s.logger.Info("VendorA RegisterOrLogin", zap.String("playerID", playerID))
	// 實作應該呼叫 s.client 並 map 回 models.Player；這裡示範回傳樣板
	return &models.Player{
		ID:       playerID,
		Username: "vendor_a_user",
		Token:    "example-token",
	}, nil
}

func (s *VendorAService) TransferOut(ctx context.Context, playerID string, amount float64) (*models.WalletTransaction, error) {
	// 範例：建立交易紀錄
	return &models.WalletTransaction{
		TransactionID: "tx-out-123",
		PlayerID:      playerID,
		Amount:        amount,
		Direction:     "out",
		Timestamp:     time.Now(),
	}, nil
}

func (s *VendorAService) TransferIn(ctx context.Context, playerID string, amount float64) (*models.WalletTransaction, error) {
	return &models.WalletTransaction{
		TransactionID: "tx-in-123",
		PlayerID:      playerID,
		Amount:        amount,
		Direction:     "in",
		Timestamp:     time.Now(),
	}, nil
}

func (s *VendorAService) GetGameEntry(ctx context.Context, playerID string) (*models.GameEntry, error) {
	// 範例：產生帶 token 的遊戲入口
	return &models.GameEntry{
		URL:     s.cfg.BaseURL + "/game-entry?player=" + playerID,
		Expired: time.Now().Add(5 * time.Minute).Unix(),
		Token:   "entry-token-abc",
	}, nil
}

func (s *VendorAService) Logout(ctx context.Context, playerID string) (*models.LogoutResult, error) {
	// 被動登出示範
	return &models.LogoutResult{
		Success: true,
		Message: "logout ok",
	}, nil
}

// APIClient skeleton（實務上可加 httpclient、簽名、timeout 等）
type APIClient struct {
	cfg config.VendorConfig
}

func NewAPIClient(cfg config.VendorConfig) *APIClient {
	return &APIClient{cfg: cfg}
}
