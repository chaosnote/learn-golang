package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ProvideServiceLogger 建立適合 Service 層業務使用的 zap.Logger，命名為 "service_logger"
func ProvideServiceLogger() (logger *zap.Logger, _ error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.Encoding = "console"
	cfg.EncoderConfig.TimeKey = "T"
	cfg.EncoderConfig.LevelKey = "L"
	cfg.EncoderConfig.NameKey = "N"
	cfg.EncoderConfig.MessageKey = "M"
	cfg.EncoderConfig.CallerKey = "C"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build()
}

var ServiceLoggerProvider = fx.Provide(
	fx.Annotate(
		ProvideServiceLogger,
		fx.ResultTags(`name:"service_logger"`),
	),
)
