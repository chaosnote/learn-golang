package app

// NewEmailSender 是一個 Provider，負責建立 EmailSender 實例
// 它回傳一個 MessageSender 介面型別，這讓 Greeter 可以接收
func NewEmailSender() MessageSender {
	return &EmailSender{}
}

// NewSMSSender 是一個 Provider，負責建立 SMSSender 實例
// 同樣回傳 MessageSender 介面型別
func NewSMSSender() MessageSender {
	return &SMSSender{}
}

// NewGreeter 是一個 Provider，它接受 MessageSender 介面作為參數
// 注意：它依賴於介面，這就是DI的核心！
func NewGreeter(sender MessageSender) *Greeter {
	return &Greeter{sender: sender}
}
