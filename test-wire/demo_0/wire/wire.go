//go:build wireinject

//go:generate wire

package wire

import (
	"idv/chris/robot"

	"github.com/google/wire"
)

// ProviderSet，將所有 robot 套件的提供者組織起來
var RobotSet = wire.NewSet(
	robot.NewRobot,
	robot.NewRobotArm,
	robot.NewRobotHead,
)

// InitializeRobot 是一個 Injector 函式，它會指示 wire 建立一個 *robot.Robot 物件
// 它會使用 RobotSet 中的提供者來滿足依賴
func InitializeRobot() *robot.Robot {
	wire.Build(RobotSet)
	return nil
}
