package server

import (
	"idv/chris/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 將所有路由註冊到 Gin 引擎
func RegisterRoutes(r *gin.Engine, vm *services.VendorManager) {
	// 健康檢查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 玩家登入或註冊路由範例
	r.POST("/vendor/:vendorName/player/login", func(c *gin.Context) {
		vendorName := c.Param("vendorName")
		playerID := c.PostForm("player_id")

		svc, err := vm.GetVendorService(vendorName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		player, err := svc.RegisterOrLogin(c.Request.Context(), playerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, player)
	})
}
