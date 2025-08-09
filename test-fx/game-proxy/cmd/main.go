package main

import (
	"idv/chris/internal/app/config"
	"idv/chris/internal/app/logger"
	"idv/chris/internal/app/modules/vendor_a"
	"idv/chris/internal/app/server"
	"idv/chris/internal/app/services"

	"go.uber.org/fx"
)

// main 組合所有 fx providers/modules，啟動應用
func main() {
	app := fx.New(
		// 1. Config
		fx.Provide(config.ProvideConfig),

		// 2. Loggers (命名)
		logger.GinLoggerProvider,
		logger.ServiceLoggerProvider,

		// 3. Vendor modules (每個 module 必須把其 VendorGameService 放到 group:"vendors")
		vendor_a.Module,

		// 4. VendorManager (從 group:"vendors" 收集所有實作)
		fx.Provide(services.ProvideVendorManager),

		// 5. Gin engine (使用命名 gin logger)
		fx.Provide(server.ProvideGinEngine),

		// 6. 啟動 HTTP Server (RegisterServer 會使用 fx.Lifecycle)
		fx.Invoke(server.RegisterServer),
	)

	app.Run()
}
