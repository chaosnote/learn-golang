package main

//int sum(int a, int b) { return a+b; }
import "C"

// 學習指命
// go tool cgo main.go

func main() {
	println(C.sum(1, 1))
}
