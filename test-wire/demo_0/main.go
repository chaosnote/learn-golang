//
// project/
// ├── main.go
// ├── robot/
// │   ├── parts.go
// │   └── service.go
// └── wire/
//     └── wire.go
//

//
// go get github.com/google/wire
// go install github.com/google/wire/cmd/wire@latest
//

package main

import (
	"idv/chris/wire"
)

func main() {
	// 呼叫 wire 產生的函式，獲得已組裝好的機器人實例
	myRobot := wire.InitializeRobot()

	// 讓機器人執行動作
	myRobot.DoAction()
}
