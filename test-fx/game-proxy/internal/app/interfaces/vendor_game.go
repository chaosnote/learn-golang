package interfaces

import (
	"context"
	"idv/chris/internal/app/models"
)

// VendorGameService 定義每個廠商需要實作的行為（你要求的完整方法都在）
type VendorGameService interface {
	// 1. 註冊或登入玩家（若廠商要做分兩階段可自行在實作處理）
	RegisterOrLogin(ctx context.Context, playerID string) (*models.Player, error)

	// 2. 我方點數轉出（扣我方點數並呼叫第三方）
	TransferOut(ctx context.Context, playerID string, amount float64) (*models.WalletTransaction, error)

	// 3. 第三方點數轉入（從對方扣點並增加我方）
	TransferIn(ctx context.Context, playerID string, amount float64) (*models.WalletTransaction, error)

	// 4. 取得遊戲入口（URL / token）
	GetGameEntry(ctx context.Context, playerID string) (*models.GameEntry, error)

	// 5. 玩家登出第三方（被動通知）
	Logout(ctx context.Context, playerID string) (*models.LogoutResult, error)

	// 回傳廠商識別名稱（例如 "vendor_a"），讓 VendorManager 可以用來索引
	GetName() string
}
