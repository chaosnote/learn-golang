package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/looplab/fsm"
)

type FSMFactory func(id string) *fsm.FSM

func NewFSMFactory() FSMFactory {
	return func(id string) *fsm.FSM {
		return fsm.NewFSM(
			"red",
			fsm.Events{
				{Name: "next", Src: []string{"red"}, Dst: "green"},
				{Name: "next", Src: []string{"green"}, Dst: "yellow"},
				{Name: "next", Src: []string{"yellow"}, Dst: "red"},
			},
			// 直接放函式，簽名為 func(ctx context.Context, e *fsm.Event)
			fsm.Callbacks{
				"enter_state": func(ctx context.Context, e *fsm.Event) {
					fmt.Printf("[紅綠燈 %s] 進入狀態: %s\n", id, e.Dst)
				},
			},
		)
	}
}

type TrafficController struct {
	fx.In
	MakeFSM FSMFactory
}

func (tc TrafficController) InitLights() {
	ctx := context.Background()

	lightA := tc.MakeFSM("A")
	lightB := tc.MakeFSM("B")

	// v2 的 Event 需帶 context
	lightA.Event(ctx, "next")
	lightA.Event(ctx, "next")
	lightB.Event(ctx, "next")
	lightA.Event(ctx, "next")
	lightB.Event(ctx, "next")
}

var Module = fx.Module("traffic",
	fx.Provide(NewFSMFactory),
	fx.Invoke(func(tc TrafficController) {
		tc.InitLights()
	}),
)

func main() {
	fx.New(Module).Run()
}
