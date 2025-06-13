package utils

import (
	"html/template"
	"idv/chris/model"
)

func RenderTemplate(config model.TemplateConfig) (tmpl *template.Template, err error) {
	tmpl, err = template.ParseFiles(append([]string{config.Layout}, config.Page...)...)
	if err != nil {
		return
	}
	for _, value := range config.Pattern {
		tmpl, err = tmpl.ParseGlob(value)
		if err != nil {
			return
		}
	}
	return
}
