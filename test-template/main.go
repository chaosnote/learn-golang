package main

import (
	"idv/chris/model"
	"idv/chris/utils"
	"log"
	"net/http"
	"path/filepath"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	config := model.DefaultTemplateConfig()
	config.Page = append(config.Page, filepath.Join(model.DefaultTemplateDir, "page", "home.html"))
	tmpl, err := utils.RenderTemplate(config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var list []model.MenuItem = []model.MenuItem{}
	list = append(list, model.MenuItem{Label: "A"})
	list = append(list, model.MenuItem{Label: "B"})
	list = append(list, model.MenuItem{Label: "C"})

	var output = map[string]any{}
	output["list"] = list
	output["title"] = "Home"

	// 2. 執行名為 "home.html" 的模板
	// 注意：ExecuteTemplate 第二個參數是模板的名稱
	err = tmpl.ExecuteTemplate(w, "home.html", output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleCountry(w http.ResponseWriter, r *http.Request) {
	config := model.DefaultTemplateConfig()
	config.Page = append(config.Page, filepath.Join(model.DefaultTemplateDir, "page", "countries.html"))
	tmpl, err := utils.RenderTemplate(config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var list []model.CountryItem = []model.CountryItem{}
	list = append(list, model.CountryItem{Code: "00", Label: "A"})
	list = append(list, model.CountryItem{Code: "01", Label: "B"})

	var output = map[string]any{}
	output["list"] = list
	output["title"] = "Country"

	err = tmpl.ExecuteTemplate(w, "countries.html", output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/country", HandleCountry)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
