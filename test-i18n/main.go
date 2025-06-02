package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os" // 引入 os 套件
	"path/filepath"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// setupI18n 用於初始化 go-i18n Bundle，並載入語系檔案
func setupI18n(langDir string) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)           // 設定預設語系，當找不到翻譯時使用
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal) // 註冊 JSON 解析器

	// 遍歷語系目錄，載入所有 .json 檔案
	// 使用 os.ReadDir 替代 io/ioutil.ReadDir
	files, err := os.ReadDir(langDir)
	if err != nil {
		return nil, fmt.Errorf("無法讀取語系目錄 %s: %w", langDir, err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(langDir, file.Name())

			// 使用 os.ReadFile 替代 io/ioutil.ReadFile
			// LoadMessageFile 函式現在可以直接接收一個 io.Reader，但 go-i18n 的 v2 版本
			// 的 LoadMessageFile 仍然是接收 string 路徑，所以我們直接傳遞路徑即可。
			// 如果你未來看到 LoadMessageFile 變成需要 io.Reader，則需先讀取檔案內容
			// bytes, err := os.ReadFile(filePath)
			// if err != nil { /* 處理錯誤 */ }
			// bundle.ParseMessageFileBytes(bytes, filePath) // 可能會是類似這樣

			// 由於 go-i18n/v2/i18n 的 LoadMessageFile 仍支援直接讀取路徑
			// 且內部會處理檔案讀取，所以這裡不需要我們手動 os.ReadFile。
			// 它的警告是針對直接使用 io/ioutil.ReadFile 或 io/ioutil.WriteFile 等函數。
			// 我們僅需確保如果我們**直接讀取檔案內容**，就使用 os.ReadFile。
			// 對於 LoadMessageFile 這種接收路徑的函式，它的內部實現會處理檔案讀取，
			// go-i18n 庫本身會處理棄用函式的替換。

			// 因此，最直接的修改是針對 ReadDir。
			// 如果未來 go-i18n 更新其 LoadMessageFile 簽名，才需要改為 os.ReadFile。
			// 目前，這裡的改動主要是針對 os.ReadDir。

			_, err := bundle.LoadMessageFile(filePath)
			if err != nil {
				log.Printf("警告: 無法載入語系檔案 %s: %v", filePath, err)
			} else {
				log.Printf("成功載入語系檔案: %s", filePath)
			}
		}
	}
	return bundle, nil
}

func main() {
	// 語系檔案所在的目錄
	langDir := "i18n"

	// 初始化 i18n Bundle
	bundle, err := setupI18n(langDir)
	if err != nil {
		log.Fatalf("初始化 i18n 失敗: %v", err)
	}

	// 示範多語系翻譯
	languages := []language.Tag{
		language.English,
		language.TraditionalChinese,
		language.MustParse("zh-Hant"), // 也可以這樣解析
	}

	for _, lang := range languages {
		fmt.Printf("\n--- 嘗試語系: %s ---\n", lang.String())

		// 建立 Localizer
		localizer := i18n.NewLocalizer(bundle, lang.String())

		// 翻譯 "hello" 訊息
		helloMsg, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "hello",
			TemplateData: map[string]interface{}{
				"Name": "Go 程式",
			},
		})
		if err != nil {
			log.Printf("翻譯 'hello' 失敗 (%s): %v", lang.String(), err)
		} else {
			fmt.Printf("翻譯 'hello': %s\n", helloMsg)
		}

		// 翻譯 "welcome_message"
		welcomeMsg, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "welcome_message",
		})
		if err != nil {
			log.Printf("翻譯 'welcome_message' 失敗 (%s): %v", lang.String(), err)
		} else {
			fmt.Printf("翻譯 'welcome_message': %s\n", welcomeMsg)
		}

		// 翻譯 "item_count" (單數)
		itemCountOne, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID:   "item_count",
			PluralCount: 1, // 這裡設定單數
			TemplateData: map[string]interface{}{
				"Count": 1,
			},
		})
		if err != nil {
			log.Printf("翻譯 'item_count' (單數) 失敗 (%s): %v", lang.String(), err)
		} else {
			fmt.Printf("翻譯 'item_count' (單數): %s\n", itemCountOne)
		}

		// 翻譯 "item_count" (複數)
		itemCountOther, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID:   "item_count",
			PluralCount: 5, // 這裡設定複數
			TemplateData: map[string]interface{}{
				"Count": 5,
			},
		})
		if err != nil {
			log.Printf("翻譯 'item_count' (複數) 失敗 (%s): %v", lang.String(), err)
		} else {
			fmt.Printf("翻譯 'item_count' (複數): %s\n", itemCountOther)
		}

		// 翻譯 "current_time"
		currentTimeMsg, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "current_time",
			TemplateData: map[string]interface{}{
				"Time": time.Now().Format("2006-01-02 15:04:05"),
			},
		})
		if err != nil {
			log.Printf("翻譯 'current_time' 失敗 (%s): %v", lang.String(), err)
		} else {
			fmt.Printf("翻譯 'current_time': %s\n", currentTimeMsg)
		}
	}

	// 也可以根據環境變數來獲取語系
	fmt.Printf("\n--- 根據環境變數 (LANG) 嘗試翻譯 ---\n")
	os.Setenv("LANG", "zh-Hant") // 模擬設定環境變數
	localizer := i18n.NewLocalizer(bundle, os.Getenv("LANG"))
	helloMsgEnv, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "hello",
		TemplateData: map[string]interface{}{
			"Name": "使用者",
		},
	})
	if err != nil {
		log.Printf("翻譯 'hello' (環境變數) 失敗: %v", err)
	} else {
		fmt.Printf("環境變數語系翻譯 'hello': %s\n", helloMsgEnv)
	}

	// 清除環境變數設定
	os.Unsetenv("LANG")
}
