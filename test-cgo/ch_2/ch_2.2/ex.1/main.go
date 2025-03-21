package main

// 導入 C header 建議獨立行
// C header 函式定義方式並不受限於 Golang ( 例: 大小寫 )

/*
#include <stdio.h>

void printint(int v) {
    printf("printint: %d\n", v);
}
*/
import "C"

func main() {
	v := 42
	C.printint(C.int(v))
}
