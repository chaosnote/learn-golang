package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AppConfig 是全域配置結構
type AppConfig struct {
	Server  ServerConfig            `yaml:"server"`
	Vendors map[string]VendorConfig `yaml:"vendors"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

// 之後可調整為 json 對應 map[string]interface{}
type VendorConfig struct {
	BaseURL    string `yaml:"base_url"`
	APIKey     string `yaml:"api_key"`
	Secret     string `yaml:"secret"`
	TimeoutSec int    `yaml:"timeout_sec"`
	Currency   string `yaml:"currency"`
}

// ProvideConfig 從 configs/config.yaml 讀取設定
func ProvideConfig() (*AppConfig, error) {
	f, err := os.Open("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("open config.yaml: %w", err)
	}
	defer f.Close()

	var cfg AppConfig
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config.yaml: %w", err)
	}
	return &cfg, nil
}
