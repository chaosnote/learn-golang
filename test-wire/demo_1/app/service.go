package app

import "fmt"

// MessageSender 介面：定義了「發送訊息」的行為
type MessageSender interface {
	Send(message string) error
}

// EmailSender 結構體：MessageSender 的一種具體實作
type EmailSender struct{}

func (e *EmailSender) Send(message string) error {
	fmt.Println("透過電子郵件發送:", message)
	return nil
}

// SMSSender 結構體：MessageSender 的另一種具體實作
type SMSSender struct{}

func (s *SMSSender) Send(message string) error {
	fmt.Println("透過簡訊發送:", message)
	return nil
}

// Greeter 結構體：我們的核心服務，它「依賴」於一個 MessageSender 介面
type Greeter struct {
	sender MessageSender
}

func (g *Greeter) Greet(name string) {
	message := fmt.Sprintf("你好, %s!", name)
	g.sender.Send(message)
}
