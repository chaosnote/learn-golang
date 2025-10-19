package http_server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// 匯出 HTTP Server 模組（含 Engine 與啟動流程）
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewGinEngine),
		fx.Invoke(RegisterServer),
	)
}

func NewGinEngine() *gin.Engine {
	r := gin.Default()
	return r
}

func RegisterServer(lc fx.Lifecycle, logger *zap.Logger, engine *gin.Engine) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("HTTP server starting", zap.String("addr", server.Addr))
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("HTTP server stopping")
			return server.Shutdown(ctx)
		},
	})
}
