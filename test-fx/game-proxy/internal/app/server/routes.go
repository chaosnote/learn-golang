package server

import (
	"net/http"

	// 模組化路由：你可新增 vendor_b 並在此掛載
	"idv/chris/internal/app/services"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 掛載全域路由與各模組自註冊路由
func RegisterRoutes(r *gin.Engine, vm *services.VendorManager) {
	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		names := []string{}
		for _, value := range vm.GetAllVendors() {
			names = append(names, value.GetName())
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "vendor": names})
	})

	// 玩家登入或註冊路由範例
	r.POST("/vendor/:name/player/login", func(c *gin.Context) {
		vendorName := c.Param("name")
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
