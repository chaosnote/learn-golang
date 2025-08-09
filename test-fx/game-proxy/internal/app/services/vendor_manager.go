package services

import (
	"errors"
	"idv/chris/internal/app/interfaces"
	"sync"

	"go.uber.org/fx"
)

// VendorManager 保留 map[string]VendorGameService 的原始設計
type VendorManager struct {
	vendors map[string]interfaces.VendorGameService
	mu      sync.RWMutex
}

// ProvideVendorManager 使用 fx.In 接收 group:"vendors" 的切片
func ProvideVendorManager(params struct {
	fx.In
	Vendors []interfaces.VendorGameService `group:"vendors"`
}) *VendorManager {
	vm := &VendorManager{
		vendors: make(map[string]interfaces.VendorGameService),
	}
	for _, v := range params.Vendors {
		vm.vendors[v.GetName()] = v
	}
	return vm
}

// GetVendorService 依名稱取出服務
func (vm *VendorManager) GetVendorService(name string) (interfaces.VendorGameService, error) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	if svc, ok := vm.vendors[name]; ok {
		return svc, nil
	}
	return nil, errors.New("vendor service not found")
}

// GetAllVendors 返回所有註冊的廠商（供動態掛載路由之用）
func (vm *VendorManager) GetAllVendors() []interfaces.VendorGameService {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	res := make([]interfaces.VendorGameService, 0, len(vm.vendors))
	for _, v := range vm.vendors {
		res = append(res, v)
	}
	return res
}
