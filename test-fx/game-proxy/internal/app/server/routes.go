package server

import (
	"net/http"

	"idv/chris/internal/app/modules/vendor_a" // 模組化路由：你可新增 vendor_b 並在此掛載
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

	// 各模組自行管理自己的路由，這裡只負責呼叫它們
	vendor_a.RegisterRoutes(r, vm)
}
