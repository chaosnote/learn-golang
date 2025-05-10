package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 用於存放每個 UID 的條件變數和鎖的狀態
type uidState struct {
	cond   *sync.Cond
	locked bool // 標誌位，指示鎖是否被持有
	locker int  // 持有鎖的 goroutine 的 ID (用於調試，可選)
}

var uidStates = make(map[string]*uidState)
var mu sync.Mutex
var nextGoroutineID int = 1

// 模擬資料庫操作
func fetchDataFromDatabase(uid string) (string, error) {
	time.Sleep(2 * time.Second)
	return fmt.Sprintf("Data for UID: %s", uid), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "UID parameter is required", http.StatusBadRequest)
		return
	}

	goroutineID := nextGoroutineID
	nextGoroutineID++

	// 取得特定 UID 的狀態或創建一個新的
	mu.Lock()
	state, ok := uidStates[uid]
	if !ok {
		state = &uidState{
			cond: sync.NewCond(&sync.Mutex{}),
		}
		uidStates[uid] = state
	}
	mu.Unlock()

	// 取得特定 UID 的鎖
	state.cond.L.Lock()
	state.locked = true
	state.locker = goroutineID
	fmt.Printf("Goroutine %d acquired lock for UID: %s at %s\n", goroutineID, uid, time.Now().Format(time.RFC3339))

	// 使用 defer 確保在函數退出時釋放鎖和更新狀態
	defer func() {
		state.locked = false
		state.locker = 0
		state.cond.L.Unlock()
		state.cond.Signal() // 通知下一個等待的 goroutine
		fmt.Printf("Goroutine %d released lock for UID: %s at %s\n", goroutineID, uid, time.Now().Format(time.RFC3339))
	}()

	// 模擬 API 請求 (略)
	// ...

	// 檢查是否有其他相同 UID 的請求正在處理
	if state.locked && state.locker != goroutineID {
		fmt.Printf("Goroutine %d waiting for UID: %s, current locker: %d\n", goroutineID, uid, state.locker)
		state.cond.Wait() // 釋放 cond.L 並等待通知
		fmt.Printf("Goroutine %d resumed for UID: %s at %s\n", goroutineID, uid, time.Now().Format(time.RFC3339))
	}

	// 模擬資料庫請求
	data, err := fetchDataFromDatabase(uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Data for UID %s: %s\n", uid, data)
}

func main() {
	http.HandleFunc("/data", handler)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
