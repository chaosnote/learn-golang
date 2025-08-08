package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// --- 1. 定義介面與多個實作 ---

// User 定義了使用者服務的介面。
type User interface {
	GetName() string
}

// Service_CCC 是 User 介面的一個實作。
type Service_CCC struct{}

func (s *Service_CCC) GetName() string {
	return "User from Service CCC"
}

// Service_DDD 是 User 介面的一個實作。
type Service_DDD struct{}

func (s *Service_DDD) GetName() string {
	return "User from Service DDD"
}

// --- 2. Handler 和 Server ---

// ServiceMap 是一個將服務名稱映射到 User 介面的 map。
type ServiceMap map[string]User

// UserHandler 持有 ServiceMap。
type UserHandler struct {
	services ServiceMap
}

// NewUserHandler 接收一個包含所有服務的 map。
func NewUserHandler(services ServiceMap) *UserHandler {
	return &UserHandler{services: services}
}

// HandleUser 處理 HTTP 請求。
func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := strings.ToUpper(vars["serviceName"])

	// 透過 map 動態查找服務，避免 switch-case 語句。
	userSvc, ok := h.services[serviceName]
	if !ok {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, userSvc.GetName())
}

// NewRouter 建立 mux.Router 並掛上 handler。
func NewRouter(h *UserHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user/{serviceName}", h.HandleUser).Methods("GET")
	return router
}

// StartServer 啟動 HTTP Server。
func StartServer(lifecycle fx.Lifecycle, router *mux.Router) {
	server := &http.Server{Addr: ":8080", Handler: router}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Server is running on :8080")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

// --- 3. Fx 提供者和主程式 ---

// UsersCollect 封裝具名注入的 User 實例。
type UsersCollect struct {
	fx.In
	CCC User `name:"CCC"`
	DDD User `name:"DDD"`
}

// curl http://localhost:8080/user/ccc
// curl http://localhost:8080/user/ddd
// curl http://localhost:8080/user/aaa
// ∟ 失敗測試

func main() {
	fx.New(
		// 提供所有具名的 User 實作
		// 簡單來說，匿名函式是「做事的」，而 fx.Annotate 則是「描述這件事該怎麼做的」。
		// 如果沒有 fx.Annotate，fx.Provide 只能直接提供匿名函式，此時 fx 容器會：
		// 看到有兩個匿名函式都返回 User 介面。
		// 無法區分這兩個實作，因為它們沒有名字。
		// 當有其他組件需要注入 User 時，fx 會因為不確定要注入哪一個而報錯。
		// 因此，fx.Annotate 就像是一個標籤工具，它讓你可以更靈活地控制如何提供和注入服務，特別是當你有多個相同介面的實作時。
		fx.Provide(
			fx.Annotate(
				func() User { return &Service_CCC{} },
				fx.As(new(User)),
				fx.ResultTags(`name:"CCC"`),
			),
			fx.Annotate(
				func() User { return &Service_DDD{} },
				fx.As(new(User)),
				fx.ResultTags(`name:"DDD"`),
			),
		),

		// 用 fx.Provide 組裝 ServiceMap（不能用 fx.Decorate）
		fx.Provide(func(users UsersCollect) ServiceMap {
			m := make(ServiceMap)
			m["CCC"] = users.CCC
			m["DDD"] = users.DDD
			return m
		}),

		// 其他組件
		fx.Provide(NewUserHandler, NewRouter, zap.NewDevelopment),
		fx.Invoke(StartServer),
	).Run()
}
