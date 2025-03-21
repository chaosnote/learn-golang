package main

// C 語言移至外部定義

/*
void SayHello(const char* s);
*/
import "C"

func main() {
	C.SayHello(C.CString("Hello, World\n"))
}
