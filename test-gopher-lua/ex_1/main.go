package main

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	// 執行 Lua 檔案
	e := L.DoFile("./lib/script.lua")
	if e != nil {
		log.Fatal(e)
	}

	// 呼叫 Lua 函式 "add"
	e = L.CallByParam(
		lua.P{
			Fn:   L.GetGlobal("add"), // 取得 Lua 函式
			NRet: 1,                  // 預期回傳值數量
		},
		lua.LNumber(10), lua.LNumber(20), // 傳入參數 (使用 Params)
	)

	if e != nil {
		log.Fatal(e)
	}
	res := L.Get(-1)                  // 取得回傳值
	L.Pop(1)                          // 移除回傳值
	fmt.Println("add(10, 20) =", res) // 輸出結果

	// 呼叫 Lua 函式 "greet"
	e = L.CallByParam(
		lua.P{
			Fn:   L.GetGlobal("greet"),
			NRet: 1,
		},
		lua.LString("World"), // 傳入參數 (使用 Params)
	)
	if e != nil {
		log.Fatal(e)
	}
	res = L.Get(-1)
	L.Pop(1)
	fmt.Println("greet(\"World\") =", res)
}
