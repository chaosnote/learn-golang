package vendor_a

import (
	"context"
	"idv/chris/internal/app/config"
	"idv/chris/internal/app/models"

	"go.uber.org/zap"
)

// VendorAService 是 VendorA 的服務實作
type VendorAService struct {
	cfg           config.VendorConfig
	client        *APIClient
	serviceLogger *zap.Logger
}

func NewVendorAService(cfg config.VendorConfig, client *APIClient, logger *zap.Logger) *VendorAService {
	return &VendorAService{
		cfg:           cfg,
		client:        client,
		serviceLogger: logger,
	}
}

// RegisterOrLogin 範例實作（可改為真實呼叫第三方 API）
func (s *VendorAService) RegisterOrLogin(ctx context.Context, playerID string) (*models.Player, error) {
	s.serviceLogger.Info("VendorA RegisterOrLogin called", zap.String("playerID", playerID))
	return &models.Player{
		ID:       playerID,
		Username: "vendor_a_user",
		Token:    "example-token",
	}, nil
}

func (s *VendorAService) TransferOut(ctx context.Context, playerID string, amount float64) (*models.Wallet, error) {
	return nil, nil
}

func (s *VendorAService) TransferIn(ctx context.Context, playerID string, amount float64) (*models.Wallet, error) {
	return nil, nil
}

func (s *VendorAService) GetGameEntry(ctx context.Context) (*models.GameEntry, error) {
	return nil, nil
}

func (s *VendorAService) Logout(ctx context.Context, playerID string) (*models.LogoutResult, error) {
	return nil, nil
}

// APIClient 是示範的第三方 API 客戶端結構，實作可自行擴充
type APIClient struct {
	cfg config.VendorConfig
}

func NewAPIClient(cfg config.VendorConfig) *APIClient {
	return &APIClient{cfg: cfg}
}
