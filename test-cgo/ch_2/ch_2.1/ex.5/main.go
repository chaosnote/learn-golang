package main

//#include "hello.h"
import "C"

func main() {
	C.SayHello(C.CString("Hello, World\n"))
	C.SayHello(C.CString("Make EZ Way\n"))
}
