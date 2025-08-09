package server

import (
	"context"
	"fmt"
	"net/http"

	"idv/chris/internal/app/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// RegisterServer 使用 fx.Lifecycle 啟動與關閉 HTTP Server
func RegisterServer(lc fx.Lifecycle, engine *gin.Engine, vm *services.VendorManager) {
	// 註冊路由（可視需求拆分或擴充）
	RegisterRoutes(engine, vm)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
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
