// hello.go
package main

/*
	使用 cgo 與 定義檔
*/

import "C"

import "fmt"

//export SayHello
func SayHello(s *C.char) {
	fmt.Print(C.GoString(s))
}
