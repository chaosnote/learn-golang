package main

/*
#include <stdlib.h>
#include <stdio.h>

void printString(const char* s) {
    printf("%s", s);
}
*/
import "C"
import "unsafe"

func printString(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs)) // 需手動釋放

	C.printString(cs)
}

func main() {
	s := "hello"
	printString(s)
}
