package controllers

import (
	"net/http"

	"idv/chris/app/article"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// ControllerParam 用 Fx 注入依賴
type ControllerParam struct {
	fx.In
	Engine *gin.Engine
	UC     *article.ArticleUseCase
}

// NewArticleController 負責註冊路由
func NewArticleController(p ControllerParam) {
	group := p.Engine.Group("/api/v1/article")
	group.GET("/hello", func(c *gin.Context) {
		// msg := p.UC.Hello(context.Background())
		c.JSON(http.StatusOK, gin.H{"message": "hello"})
	})
}

// 匯出 Fx 模組
func Module() fx.Option {
	return fx.Options(
		fx.Invoke(NewArticleController),
	)
}
