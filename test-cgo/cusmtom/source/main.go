package main

/*
#include <stdio.h>
#include <stdlib.h>

// 假設有一個 C 程式需要呼叫 Go 的 Message 函式
// 可以使用以下方式呼叫：
// extern void Message();
// Message();
*/
import "C"
import (
	"fmt"
)

// Message 函式會印出 "Hello from Go DLL"
//
//export Message
func Message() {
	fmt.Println("Hello from Go DLL")
}

func main() {
	// 如果需要測試，可以在這裡呼叫 Message()
	// 但通常情況下，應該由 C 程式呼叫
	// C.Message()
}
