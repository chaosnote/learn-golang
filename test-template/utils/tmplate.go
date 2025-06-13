package utils

import (
	"html/template"
	"io/fs"
	"path/filepath"
)

// 解析模板文件
//
// 使用 ParseGlob 載入所有 .html 文件，包括子目錄中的文件
// "**/*.html" 表示匹配 templates/**/*.html 文件
// tmpl, err := template.ParseGlob(filepath.Join("templates", "**", "*.html"))
// if err != nil {
// 	log.Fatalf("Error parsing templates: %v", err)
// }
//
// tmpl, err = tmpl.ParseFiles(filepath.Join("templates", "layout.html"))
// if err != nil {
// 	log.Fatalf("Error parsing templates: %v", err)
// }

func InitTemplate(dir_path string) (*template.Template, error) {
	tmpl := template.New("")

	err := filepath.WalkDir(dir_path, func(file_path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(file_path) != ".html" {
			return nil
		}
		_, err = tmpl.ParseFiles(file_path)
		return err
	})

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
