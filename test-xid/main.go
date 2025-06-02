package main

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

func run_0() {
	// 生成一個新的 XID
	id := xid.New()
	fmt.Println("生成的 XID:", id.String())

	// --- 提取 XID 的資訊 ---

	// 提取時間戳：XID 的第一個部分就是時間戳
	// 它會回傳一個 time.Time 物件
	timestamp := id.Time()
	fmt.Println("XID 的時間戳:", timestamp.Format(time.RFC3339)) // 格式化輸出時間

	// 提取機器識別碼 (Machine ID)
	machineID := id.Machine()
	fmt.Printf("XID 的機器識別碼 (byte array): %v\n", machineID)
	// 如果你想看到十六進制表示，可以這樣做：
	fmt.Printf("XID 的機器識別碼 (hex): %x\n", machineID)

	// 提取進程 ID (Process ID)
	processID := id.Pid()
	fmt.Println("XID 的進程 ID:", processID)

	// 提取計數器 (Counter)
	counter := id.Counter()
	fmt.Println("XID 的計數器:", counter)

	// --- 從字串解析 XID ---

	// 假設你有一個 XID 的字串形式，你想把它解析回來
	xidStr := id.String() // 這是上面範例中生成的其中一個 XID
	parsedID, err := xid.FromString(xidStr)
	if err != nil {
		fmt.Println("解析 XID 字串失敗:", err)
		return
	}
	fmt.Println("\n從字串解析的 XID:", parsedID.String())
	fmt.Println("解析後的時間戳:", parsedID.Time().Format(time.RFC3339))
}

func run_1() {
	id := xid.New()
	counter := id.Counter()
	fmt.Printf("%v : %v\n", counter, id.String())
}

func main() {
	if false {
		for i := 0; i < 10; i++ {
			run_0()
			println("======")
		}
	}
	for i := 0; i < 10; i++ {
		run_1()
		println("======")
	}
}
