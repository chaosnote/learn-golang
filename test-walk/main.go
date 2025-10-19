package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

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

		fmt.Println(prefix + connector + entry.Name())

		if entry.IsDir() {
			subPrefix := prefix
			if i == len(entries)-1 {
				subPrefix += "    "
			} else {
				subPrefix += "│   "
			}
			PrintTree(filepath.Join(path, entry.Name()), subPrefix)
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
	PrintTree(pwd, "")
}
