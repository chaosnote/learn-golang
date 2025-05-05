package main

/*
二種系統對靜態檔的路徑設定不相同
win   放置於專案底下
linux 放置於 lib 底下
*/

/*
#cgo windows LDFLAGS: -L. -lcustom_lib.dll
#cgo linux LDFLAGS: lib/custom_lib.so

#include <stdlib.h>
#include "custom_lib.h"
*/
import "C"

func main() {
	C.Message()
}
