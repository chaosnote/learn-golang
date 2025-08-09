package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AppConfig 是應用全域設定結構
type AppConfig struct {
	Server  ServerConfig            `yaml:"server"`
	Vendors map[string]VendorConfig `yaml:"vendors"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type VendorConfig struct {
	Name       string `yaml:"name"`
	BaseURL    string `yaml:"base_url"`
	APIKey     string `yaml:"api_key"`
	Secret     string `yaml:"secret"`
	TimeoutSec int    `yaml:"timeout_sec"`
	Currency   string `yaml:"currency"`
}

// ProvideConfig 讀取配置檔並回傳 AppConfig
func ProvideConfig() (*AppConfig, error) {
	file, err := os.Open("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config.yaml: %w", err)
	}
	defer file.Close()

	var cfg AppConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config.yaml: %w", err)
	}

	return &cfg, nil
}
