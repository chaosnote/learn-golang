package services

import (
	"errors"
	"idv/chris/internal/app/interfaces"
	"sync"

	"go.uber.org/fx"
)

// VendorManager 負責管理多個第三方廠商服務
type VendorManager struct {
	vendors map[string]interfaces.VendorGameService
	mu      sync.RWMutex
}

// ProvideVendorManager Fx 提供 VendorManager，接收多個 VendorGameService 實例 (fx.In + group)
func ProvideVendorManager(params struct {
	fx.In
	Vendors []interfaces.VendorGameService `group:"vendors"` // 這裡用 group
}) *VendorManager {
	vm := &VendorManager{
		vendors: make(map[string]interfaces.VendorGameService),
	}
	for _, v := range params.Vendors {
		vm.vendors[v.GetName()] = v
	}
	return vm
}

// GetVendorService 根據廠商名稱取得對應服務
func (vm *VendorManager) GetVendorService(name string) (interfaces.VendorGameService, error) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	if svc, ok := vm.vendors[name]; ok {
		return svc, nil
	}
	return nil, errors.New("vendor service not found")
}
