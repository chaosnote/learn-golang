package vendor_a

import (
	"idv/chris/internal/app/config"
	"idv/chris/internal/app/interfaces"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type serviceLoggerIn struct {
	fx.In
	Logger *zap.Logger `name:"service_logger"`
}

var Module = fx.Module(
	"vendor_a",
	fx.Provide(
		func(cfg *config.AppConfig) config.VendorConfig {
			return cfg.Vendors["vendor_a"]
		},
		fx.Annotate(
			func(vendorCfg config.VendorConfig, in serviceLoggerIn) interfaces.VendorGameService {
				client := NewAPIClient(vendorCfg)
				return NewVendorAService(vendorCfg, client, in.Logger)
			},
			fx.ResultTags(`group:"vendors"`), // 這裡關鍵加 group 標籤
		),
	),
)
