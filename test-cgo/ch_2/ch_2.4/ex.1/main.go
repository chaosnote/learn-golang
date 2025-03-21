package main

/*
#include "stdio.h"
#include "stdlib.h"

static int add(int a, int b) {
    return a+b;
}
*/
import "C"

func main() {
	res := C.add(1, 1)
	println(int64(res))
}
