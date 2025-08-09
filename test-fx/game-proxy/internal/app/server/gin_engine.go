package server

import (
	"idv/chris/internal/app/server/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// 專門注入 Gin logger 的結構，標註 name:"gin_logger"
type ginLoggerIn struct {
	fx.In
	Logger *zap.Logger `name:"gin_logger"`
}

func ProvideGinEngine(in ginLoggerIn) *gin.Engine {
	r := gin.New()
	r.Use(middleware.NewLoggerMiddleware(in.Logger))
	r.Use(middleware.NewRecoveryMiddleware(in.Logger))
	return r
}
