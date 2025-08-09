package server

import (
	"idv/chris/internal/app/server/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ginLoggerIn 用來命名注入 gin logger
type ginLoggerIn struct {
	fx.In
	Logger *zap.Logger `name:"gin_logger"`
}

// ProvideGinEngine 只在這裡定義一次，回傳 *gin.Engine
func ProvideGinEngine(in ginLoggerIn) *gin.Engine {
	r := gin.New()
	// 註冊 middleware（使用命名注入的 gin logger）
	r.Use(middleware.NewLoggerMiddleware(in.Logger))
	r.Use(middleware.NewRecoveryMiddleware(in.Logger))
	return r
}
