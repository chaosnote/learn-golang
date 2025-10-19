package main

import (
	"idv/chris/app"
	"idv/chris/cmd/http_server"
	"idv/chris/config"
	"idv/chris/infra"
	"idv/chris/utils"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Module(),
		utils.Module(),
		infra.Module(),
		app.Module(),
		http_server.Module(), // ✅ 啟動 Gin server
	)
	app.Run()
}
