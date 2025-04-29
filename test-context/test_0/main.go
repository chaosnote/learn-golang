package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	fmt.Printf("Worker %d started\n", id)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d received cancel signal\n", id)
			return
		default:
			fmt.Printf("Worker %d is working...\n", id)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// 創建一個可以取消的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 啟動幾個 worker goroutine
	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	// 模擬主程序運行一段時間
	time.Sleep(5 * time.Second)

	// 調用 cancel 函數，通知所有 worker 停止工作
	fmt.Println("Sending cancel signal...")
	cancel()

	// 等待一段時間，確保所有 worker 都已停止
	time.Sleep(1 * time.Second)
	fmt.Println("All workers stopped.")
}
