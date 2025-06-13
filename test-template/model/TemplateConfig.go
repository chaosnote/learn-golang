package model

import "path/filepath"

type TemplateConfig struct {
	Layout  string
	Page    []string
	Pattern []string
}

var DefaultTemplateDir = "./templates"

func DefaultTemplateConfig() TemplateConfig {
	return TemplateConfig{
		Layout: filepath.Join(DefaultTemplateDir, "layout.html"),
		Page:   []string{},
		Pattern: []string{
			filepath.Join(DefaultTemplateDir, "component", "*.html"),
		},
	}
}
