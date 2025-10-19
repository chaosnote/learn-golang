package utils

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewZapLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger.Named("app"), nil
}

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewZapLogger),
	)
}
