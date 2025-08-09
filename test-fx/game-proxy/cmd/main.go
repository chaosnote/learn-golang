package main

import (
	"idv/chris/internal/app/config"
	"idv/chris/internal/app/logger"
	"idv/chris/internal/app/modules/vendor_a"
	"idv/chris/internal/app/server"
	"idv/chris/internal/app/services"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(config.ProvideConfig),

		// Logger Providers，分別注入 Gin 和 Service logger
		logger.GinLoggerProvider,
		logger.ServiceLoggerProvider,

		// 提供 Gin Engine
		fx.Provide(server.ProvideGinEngine),

		// Vendor Modules（可再加多個廠商Module）
		vendor_a.Module,

		// 服務管理器，會拿到 group 的 Vendors 切片
		fx.Provide(services.ProvideVendorManager),

		// 啟動 Gin Server，注入 VendorManager
		fx.Invoke(server.RegisterServer),
	)

	app.Run()
}
