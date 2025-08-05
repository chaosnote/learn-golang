package main

import (
	"fmt"
	"idv/chris/wire" // 引入 wire 套件
)

func main() {
	fmt.Println("--- 透過電子郵件發送 ---")
	// 呼叫 wire 生成的函式，獲得已注入 EmailSender 的 Greeter 服務
	emailGreeter := wire.InitializeGreeter()
	emailGreeter.Greet("Alice")

	fmt.Println("\n--- 透過簡訊發送 ---")
	// 呼叫 wire 生成的另一個函式，獲得已注入 SMSSender 的 Greeter 服務
	smsGreeter := wire.InitializeSMSGreeter()
	smsGreeter.Greet("Bob")
}
