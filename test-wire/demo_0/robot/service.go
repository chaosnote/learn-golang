package robot

import "fmt"

// Robot 結構體，它依賴 RobotArm 和 RobotHead
type Robot struct {
	Arm  RobotArm
	Head RobotHead
}

// NewRobot 是提供者，它將依賴注入到 Robot 結構體中
func NewRobot(arm RobotArm, head RobotHead) *Robot {
	return &Robot{
		Arm:  arm,
		Head: head,
	}
}

// DoAction 讓機器人執行動作
func (r *Robot) DoAction() {
	fmt.Println("機器人開始行動！")
	fmt.Println(r.Arm.Move())
	fmt.Println(r.Head.Speak())
}
