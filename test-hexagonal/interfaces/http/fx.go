package http

import (
	"idv/chris/interfaces/http/controllers"

	"go.uber.org/fx"
)

// 匯出 HTTP 模組 (包含 controllers)
func Module() fx.Option {
	return fx.Options(
		controllers.Module(),
	)
}
