package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	addr    = "nats://192.168.0.236:4222" // NATS 伺服器位址
	subject = "hello"                     // 定義主題 (Subject)
)

func pub_0(nc *nats.Conn, msg string) {
	// 發布消息
	var err error
	err = nc.Publish(subject, []byte(msg))
	if err != nil {
		log.Fatalf("發布消息失敗: %v", err)
	}
	fmt.Printf("已發布消息 '%s' 到主題: %s\n", msg, subject)
}

func sub_0(nc *nats.Conn) {
	var err error
	// (一般)訂閱消息
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("接收到來自主題 '%s' 的消息: %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("訂閱失敗: %v", err)
	}
	fmt.Printf("已訂閱主題: %s\n", subject)
}

func sub_1(nc *nats.Conn) {
	var err error
	// (佇列)訂閱消息
	group := "workers"
	_, err = nc.QueueSubscribe(subject, group, func(msg *nats.Msg) {
		fmt.Printf("接收到來自主題 '%s' 的消息: %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("訂閱失敗: %v", err)
	}
	fmt.Printf("已訂閱主題: %s\n", subject)
}

func main() {
	// 連接到 NATS 伺服器
	nc, err := nats.Connect(addr)
	if err != nil {
		log.Fatalf("無法連接到 NATS: %v", err)
	}
	defer nc.Close() // 確保程式結束時關閉連線
	fmt.Println("已連接到 NATS 伺服器:", addr)

	// sub_0(nc)
	sub_1(nc)

	// 延遲一秒，確保訂閱者已準備好接收消息
	time.Sleep(1 * time.Second)

	// pub_0(nc, "message 1")
	// pub_0(nc, "message 2")

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("程式結束")
}
