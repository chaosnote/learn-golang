package robot

// RobotArm 是樂高手臂的介面
type RobotArm interface {
	Move() string
}

// robotArm 實現了 RobotArm 介面
type robotArm struct{}

func (a *robotArm) Move() string {
	return "手臂在動..."
}

// NewRobotArm 是提供者，用於建立一個樂高手臂
func NewRobotArm() RobotArm {
	return &robotArm{}
}

// RobotHead 是樂高頭的介面
type RobotHead interface {
	Speak() string
}

// robotHead 實現了 RobotHead 介面
type robotHead struct{}

func (h *robotHead) Speak() string {
	return "頭部在說話..."
}

// NewRobotHead 是提供者，用於建立一個樂高頭
func NewRobotHead() RobotHead {
	return &robotHead{}
}
