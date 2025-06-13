package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath" // 用於處理路徑
)

type PageData struct {
	Title     string
	Countries []string
}

func LoadTemplates() (*template.Template, error) {
	tmpl := template.New("") // 這邊會產生 空白 名稱 - 不影響

	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			log.Println(path)
			_, err := tmpl.ParseFiles(path)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func main() {
	// 1. 解析所有模板文件
	// 使用 ParseGlob 載入所有 .html 文件，包括子目錄中的文件
	// "**/*.html" 表示匹配 templates/ 目錄及其所有子目錄下的 .html 文件
	// tmpl, err := template.ParseGlob(filepath.Join("templates", "**", "*.html"))
	// if err != nil {
	// 	log.Fatalf("Error parsing templates: %v", err)
	// }
	// tmpl, err = tmpl.ParseFiles(filepath.Join("templates", "layout.html"))
	// if err != nil {
	// 	log.Fatalf("Error parsing templates: %v", err)
	// }

	tmpl, err := LoadTemplates()
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	log.Println("已解析的模板名稱:")
	for _, t := range tmpl.Templates() {
		log.Printf("- %s\n", t.Name())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:     "首頁",
			Countries: []string{"美國", "加拿大", "墨西哥", "日本"},
		}

		// 2. 執行名為 "home.html" 的模板
		// 注意：ExecuteTemplate 第二個參數是模板的名稱，它是在 ParseGlob 時根據檔案名自動定義的
		// 如果 home.html 繼承了 layout.html，那麼最終渲染的是 home.html 內容
		err = tmpl.ExecuteTemplate(w, "home.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
