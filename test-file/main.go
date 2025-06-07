package main

import (
	"fmt"
	"os"
	"time"
)

const filename = "test_file.txt"

func writeToFile(content string) error {
	// os.O_TRUNC 會清空檔案內容
	// os.O_CREATE 會在檔案不存在時建立
	// os.O_WRONLY 只寫入
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("無法打開/建立檔案: %w", err)
	}
	defer file.Close() // 確保檔案最終會被關閉

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("寫入檔案失敗: %w", err)
	}

	// 强制将文件内容同步到磁盘，以尽量减少操作系统缓存的影响
	// 在某些文件系统和操作系统上，这对于确保立即看到更改至关重要
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("同步檔案到磁碟失敗: %w", err)
	}

	fmt.Printf("已寫入內容到檔案 '%s':\n---\n%s---\n", filename, content)
	return nil
}

func readFromFile() ([]byte, error) {
	// os.ReadFile 是 ioutil.ReadFile 的替代品
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("讀取檔案失敗: %w", err)
	}
	return data, nil
}

func main() {
	// 1. 寫入初始內容
	initialContent := "這是檔案的初始內容。\n第一行。\n第二行。"
	err := writeToFile(initialContent)
	if err != nil {
		fmt.Printf("錯誤: %s\n", err)
		return
	}

	// 2. 第一次讀取並列印
	fmt.Println("\n第一次讀取檔案內容:")
	data, err := readFromFile()
	if err != nil {
		fmt.Printf("錯誤: %s\n", err)
		return
	}
	fmt.Printf("---\n%s---\n", string(data))

	fmt.Println("\n=======================================================")
	fmt.Printf("請在接下來的 10 秒內，手動修改 '%s' 檔案的內容。\n", filename)
	fmt.Println("例如，將內容改為：")
	fmt.Println("新的內容來了！")
	fmt.Println("這是第三行。")
	fmt.Println("=======================================================")
	time.Sleep(10 * time.Second) // 給你時間手動修改檔案

	// 3. 第二次讀取並列印
	fmt.Println("\n第二次讀取檔案內容 (手動修改後):")
	dataAfterChange, err := readFromFile()
	if err != nil {
		fmt.Printf("錯誤: %s\n", err)
		return
	}
	fmt.Printf("---\n%s---\n", string(dataAfterChange))

	// 清理：刪除測試檔案
	defer func() {
		if err := os.Remove(filename); err != nil {
			fmt.Printf("清理檔案 '%s' 失敗: %s\n", filename, err)
		} else {
			fmt.Printf("\n已刪除測試檔案 '%s'.\n", filename)
		}
	}()

	fmt.Println("\n程式執行結束。")
}
