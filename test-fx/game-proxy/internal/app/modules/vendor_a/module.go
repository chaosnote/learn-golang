package vendor_a

import (
	"idv/chris/internal/app/config"
	"idv/chris/internal/app/interfaces"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// vendorADeps 用於接收命名的 service_logger
type vendorADeps struct {
	fx.In
	Logger *zap.Logger `name:"service_logger"`
}

// Module 以 fx.Module 包裝 VendorA 的 providers
var Module = fx.Module(
	"vendor_a",
	fx.Provide(
		// 提供 vendor_a 的設定（從全域 AppConfig 取）
		func(cfg *config.AppConfig) config.VendorConfig {
			return cfg.Vendors["vendor_a"]
		},
		// 提供一個實作 interfaces.VendorGameService 並把它加入 group:"vendors"
		fx.Annotate(
			func(vendorCfg config.VendorConfig, deps vendorADeps) interfaces.VendorGameService {
				client := NewAPIClient(vendorCfg)
				return NewVendorAService(vendorCfg, client, deps.Logger)
			},
			// 把回傳值打包到 group:"vendors"，讓 ProvideVendorManager 能拿到切片
			fx.ResultTags(`group:"vendors"`),
		),
	),
)
