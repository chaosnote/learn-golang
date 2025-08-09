package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ProvideGinLogger 建立給 Gin middleware 用的 logger，會被命名為 "gin_logger"
func ProvideGinLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg.Build()
}

// Fx option：把 ProvideGinLogger 包起來並 result tag 為 name:"gin_logger"
var GinLoggerProvider = fx.Provide(
	fx.Annotate(
		ProvideGinLogger,
		fx.ResultTags(`name:"gin_logger"`),
	),
)
