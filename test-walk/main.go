package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	file_path = "./dist/output.txt"
)

var (
	allow_map = map[string]bool{
		".go":   true,
		".json": true,
		".html": true,
	}

	pass = flag.Bool("pass", false, "所以檔案全寫入檔案")
)

func write_file(content []byte) {
	// 1. 使用 os.OpenFile 開啟檔案，設定 O_APPEND, O_CREATE, O_WRONLY
	// O_APPEND: 追加模式
	// O_CREATE: 如果檔案不存在則建立
	// O_WRONLY: 唯寫模式
	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("無法開啟檔案: %v", err))
	}
	defer file.Close() // 確保在函式結束時關閉檔案

	// 2. 寫入數據
	_, err = file.Write([]byte(string(content) + "\n"))
	if err != nil {
		panic(fmt.Sprintf("寫入檔案失敗: %v", err))
	}
}

func PrintTree(path string, prefix string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("無法讀取目錄 %s: %v\n", path, err)
		return
	}

	for i, entry := range entries {
		connector := "├── "
		if i == len(entries)-1 {
			connector = "└── "
		}

		output_path := prefix + connector + entry.Name()
		write_file([]byte(output_path))

		if entry.IsDir() {
			subPrefix := prefix
			if i == len(entries)-1 {
				subPrefix += "    "
			} else {
				subPrefix += "│   "
			}
			PrintTree(filepath.Join(path, entry.Name()), subPrefix)
		} else if allow_map[filepath.Ext(entry.Name())] || *pass {
			content, err := os.ReadFile(filepath.Join(path, entry.Name()))
			if err != nil {
				panic(fmt.Sprintf("寫入檔案失敗: %v", err))
			}
			write_file(content)
		}
	}
}

func main() {
	pwd, e := os.Getwd()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	dir_path := flag.String("dir_path", pwd, "預設路徑")
	flag.Parse()

	pwd = *dir_path
	fmt.Println(pwd)

	os.WriteFile(file_path, []byte(""), 0777)

	PrintTree(pwd, "")
}
