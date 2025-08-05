//go:build wireinject

//go:generate wire

package wire

import (
	"idv/chris/app"

	"github.com/google/wire"
)

// Injector 函式，它定義了最終回傳的物件型別 (*app.Greeter)
// 這告訴 wire 我們的目標是組裝一個完整的 Greeter 服務
func InitializeGreeter() *app.Greeter {
	// wire.Build 會自動根據 InitializeGreeter 的回傳型別 (*app.Greeter)
	// 找到 NewGreeter 這個 Provider，然後再根據 NewGreeter 的參數 (MessageSender 介面)
	// 找到 NewEmailSender 這個 Provider 來滿足依賴。
	wire.Build(app.NewGreeter, app.NewEmailSender)
	return nil // 這只是佔位符
}

// 如果我們想換成用簡訊發送，只需要修改這裡，換成 NewSMSSender 即可
func InitializeSMSGreeter() *app.Greeter {
	wire.Build(app.NewGreeter, app.NewSMSSender)
	return nil
}
