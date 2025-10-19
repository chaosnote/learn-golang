package config

import (
	"encoding/json"
	"os"

	"go.uber.org/fx"
)

func LoadConfig() (*APPConfig, error) {
	f, err := os.Open("assets/config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg APPConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	// set gin mode if provided
	if cfg.Gin.Mode != "" {
		// gin.SetMode(cfg.Gin.Mode) // caller can set
	}
	return &cfg, nil
}

func Module() fx.Option {
	return fx.Provide(LoadConfig)
}
