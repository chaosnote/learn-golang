package main

/*
#include <stdio.h>

struct A {
    int i;
    float f;
    int type; // type 是 Go 语言的关键字
};

*/
import "C"
import "fmt"

func main() {
	var a C.struct_A
	fmt.Println(a.i)
	fmt.Println(a.f)
	fmt.Println(a._type)
}

//export helloString
func helloString(s string) {}

//export helloSlice
func helloSlice(s []byte) {}
