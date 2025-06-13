package main

import (
	"html/template"
	"idv/chris/model"
	"idv/chris/utils"
	"log"
	"net/http"
)

var (
	tmpl *template.Template
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	var list []model.MenuItem = []model.MenuItem{}
	list = append(list, model.MenuItem{Label: "A"})
	list = append(list, model.MenuItem{Label: "B"})
	list = append(list, model.MenuItem{Label: "C"})

	var output = map[string]any{}
	output["list"] = list
	output["title"] = "首頁"

	// 2. 執行名為 "home.html" 的模板
	// 注意：ExecuteTemplate 第二個參數是模板的名稱
	err := tmpl.ExecuteTemplate(w, "home.html", output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	var err error
	tmpl, err = utils.InitTemplate("./templates")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	log.Println("已解析的模板名稱:")
	for _, t := range tmpl.Templates() {
		log.Printf("- %s\n", t.Name())
	}

	http.HandleFunc("/", HandleHome)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
