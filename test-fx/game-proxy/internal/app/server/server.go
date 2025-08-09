package server

import (
	"context"
	"fmt"
	"net/http"

	"idv/chris/internal/app/config"
	"idv/chris/internal/app/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// RegisterServer 使用 fx.Lifecycle 啟動/關閉 HTTP Server
func RegisterServer(lc fx.Lifecycle, engine *gin.Engine, cfg *config.AppConfig, vm *services.VendorManager) {
	// 在啟動前註冊路由（確保 VendorManager 已提供）
	RegisterRoutes(engine, vm)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
