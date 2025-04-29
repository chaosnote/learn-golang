package main

import (
	"context"
	"fmt"
	"time"
)

func operationWithTimeout(ctx context.Context) error {
	fmt.Println("Starting operation...")
	// 模擬一個需要較長時間的操作
	time.Sleep(3 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("Operation cancelled due to timeout")
		return ctx.Err() // 返回 context 的錯誤信息
	default:
		fmt.Println("Operation completed successfully")
		return nil
	}
}

func main() {
	// 創建一個帶有 2 秒超時的 context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 確保在函數退出時取消 context，釋放資源

	err := operationWithTimeout(ctx) // lock
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
