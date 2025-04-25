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

	// 定義請求主題和訊息內容
	requestSubject := "MY.REQUEST.MANUAL"
	requestPayload := []byte("這是一個使用 PublishRequest 的請求")

	// 產生一個唯一的臨時回覆主題
	replySubject := nats.NewInbox()

	// 建立一個訂閱來接收回覆
	sub, err := nc.Subscribe(replySubject, func(msg *nats.Msg) {
		log.Printf("收到回覆 (PublishRequest): %s", string(msg.Data))
	})
	if err != nil {
		log.Fatalf("訂閱回覆主題失敗: %v", err)
	}
	defer sub.Unsubscribe()

	// 發布請求，並指定回覆主題
	err = nc.PublishRequest(requestSubject, replySubject, requestPayload)
	if err != nil {
		log.Fatalf("發布請求失敗 (PublishRequest): %v", err)
	}
	log.Printf("已發布請求 (PublishRequest) 到 %s，等待回覆到 %s", requestSubject, replySubject)

	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		log.Printf("在 %s 內沒有收到回覆 (PublishRequest)", timeout)
	}
}
