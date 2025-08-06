package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

// --- 1. Handler 層 ---
type UserService struct{}

func NewUserService() *UserService { return &UserService{} }

type UserHandler struct{ service *UserService }

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fetched: User from public API")
}

func (h *UserHandler) GetUserPrivate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fetched: User from internal API")
}

// --- 2. 兩個 Router Provider ---
type PublicRouterDeps struct {
	dig.In
	UserHandler *UserHandler
}

func NewPublicRouter(deps PublicRouterDeps) *mux.Router {
	fmt.Println("創建 public Router...")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/users/{id}", deps.UserHandler.GetUser).Methods("GET")
	return r
}

type InternalRouterDeps struct {
	dig.In
	UserHandler *UserHandler
}

func NewInternalRouter(deps InternalRouterDeps) *mux.Router {
	fmt.Println("創建 internal Router...")
	r := mux.NewRouter()
	r.HandleFunc("/internal/api/v1/users/{id}", deps.UserHandler.GetUserPrivate).Methods("GET")
	return r
}

// --- 3. ServerLauncher 的參數 struct（嵌入 dig.In）---
type ServerLauncherParams struct {
	dig.In
	PublicRouter   *mux.Router `name:"public"`
	InternalRouter *mux.Router `name:"internal"`
}

// 實際的 Launcher 結構體（不含 dig.In）
type ServerLauncher struct {
	PublicRouter   *mux.Router
	InternalRouter *mux.Router
}

// 建構 ServerLauncher（傳入值型別，回傳指標 OK）
func NewServerLauncher(p ServerLauncherParams) *ServerLauncher {
	fmt.Println("創建 ServerLauncher...")
	return &ServerLauncher{
		PublicRouter:   p.PublicRouter,
		InternalRouter: p.InternalRouter,
	}
}

func (s *ServerLauncher) Start() {
	fmt.Println("伺服器啟動中...")

	go func() {
		fmt.Println("Public 服務監聽在 :8080...")
		http.ListenAndServe(":8080", s.PublicRouter)
	}()

	go func() {
		fmt.Println("Internal 服務監聽在 :8081...")
		http.ListenAndServe(":8081", s.InternalRouter)
	}()

	select {}
}

// --- 4. main 組裝與啟動 ---
func main() {
	container := dig.New()

	// 註冊依賴
	container.Provide(NewUserService)
	container.Provide(NewUserHandler)

	container.Provide(NewPublicRouter, dig.Name("public"))
	container.Provide(NewInternalRouter, dig.Name("internal"))

	// 提供 ServerLauncher（參數 struct 分開）
	container.Provide(NewServerLauncher)

	// 啟動服務
	err := container.Invoke(func(launcher *ServerLauncher) {
		launcher.Start()
	})

	if err != nil {
		fmt.Println("伺服器啟動失敗:", err)
	}
}
