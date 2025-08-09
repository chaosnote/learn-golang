package vendor_a

import (
	"net/http"

	"idv/chris/internal/app/services"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 由 vendor_a 模組負責註冊自己路由
func RegisterRoutes(r *gin.Engine, vm *services.VendorManager) {
	group := r.Group("/vendor/vendor_a")

	group.POST("/player/login", func(c *gin.Context) {
		playerID := c.PostForm("player_id")

		svc, err := vm.GetVendorService("vendor_a")
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
