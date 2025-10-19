package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouterParams struct {
	fx.In

	Engine *gin.Engine
}

func RegisterRoutes(p RouterParams) {
	r := p.Engine.Group("/api/v1")
	// 假設 ArticleController 已以 fx.Provide 註冊
	// 這裡使用 gin handler 範例 (注入 controller via closure below)
	// 實際上 controller 會提供一個 constructor 註冊到 fx，並在此使用 fx.Invoke 注入
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
