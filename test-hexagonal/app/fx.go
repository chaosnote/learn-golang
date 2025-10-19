package app

import (
	"idv/chris/app/article"
	"idv/chris/app/node"

	"go.uber.org/fx"
)

// Module 匯出整個 Application 層的 fx.Option。
// 未來可在這裡註冊所有 UseCase。
func Module() fx.Option {
	return fx.Options(
		article.Module(),
		node.Module(),
	)
}
