package infra

import (
	"go.uber.org/fx"

	"idv/chris/infra/persistence"
)

func Module() fx.Option {
	return fx.Options(
		persistence.Module(),
	)
}
