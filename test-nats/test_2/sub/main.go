package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	addr = "nats://192.168.0.236:4222" // NATS 伺服器位址
)

func main() {
	// 連線到 NATS 伺服器
	nc, err := nats.Connect(addr)
	if err != nil {
		log.Fatalf("無法連線到 NATS: %v", err)
	}
	defer nc.Close()

	// 訂閱請求主題
	_, err = nc.Subscribe("MY.REQUEST.MANUAL", func(msg *nats.Msg) {
		log.Printf("收到請求 (PublishRequest): %s，回覆到: %s", string(msg.Data), msg.Reply)

		// 模擬處理請求
		time.Sleep(1 * time.Second)

		// 確保回覆主題不為空，然後發送回覆
		if msg.Reply != "" {
			err := nc.Publish(msg.Reply, []byte("這是 PublishRequest 的回覆"))
			if err != nil {
				log.Printf("發送回覆失敗 (PublishRequest): %v", err)
			}
		} else {
			log.Println("請求沒有指定回覆主題")
		}
	})
	if err != nil {
		log.Fatalf("訂閱失敗: %v", err)
	}
	log.Println("已訂閱 MY.REQUEST.MANUAL 主題，等待 PublishRequest...")

	time.Sleep(10 * time.Second)
}
