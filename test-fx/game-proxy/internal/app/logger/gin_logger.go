package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ProvideGinLogger 建立適合 Gin middleware 使用的 zap.Logger，命名為 "gin_logger"
func ProvideGinLogger() (logger *zap.Logger, _ error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.NameKey = "logger"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build()
}

var GinLoggerProvider = fx.Provide(
	fx.Annotate(
		ProvideGinLogger,
		fx.ResultTags(`name:"gin_logger"`),
	),
)
