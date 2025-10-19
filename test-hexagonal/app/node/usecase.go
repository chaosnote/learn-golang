package node

import (
	"context"

	"go.uber.org/fx"
)

type NodeUseCase struct{}

func NewNodeUseCase() *NodeUseCase {
	return &NodeUseCase{}
}

func (u *NodeUseCase) Ping(ctx context.Context) string {
	return "node ok"
}

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewNodeUseCase),
	)
}
