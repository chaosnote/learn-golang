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

// --- 2. 兩個 Router Provider (它們都提供 *mux.Router) ---
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

// --- 3. 啟動器 (Launcher) ---
type ServerLauncher struct {
	dig.In
	PublicRouter   *mux.Router `name:"public"`
	InternalRouter *mux.Router `name:"internal"`
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

// 這裡就是修正的地方：參數 deps 是一個值，而不是指針
func NewServerLauncher(deps ServerLauncher) *ServerLauncher {
	fmt.Println("創建 ServerLauncher...")
	return &ServerLauncher{
		PublicRouter:   deps.PublicRouter,
		InternalRouter: deps.InternalRouter,
	}
}

// --- 4. 應用程式：啟動與組裝 ---
func main() {
	container := dig.New()

	container.Provide(NewUserService)
	container.Provide(NewUserHandler)

	container.Provide(NewPublicRouter, dig.Name("public"))
	container.Provide(NewInternalRouter, dig.Name("internal"))

	// 註冊啟動器，它會回傳一個 *ServerLauncher
	container.Provide(NewServerLauncher)

	err := container.Invoke(func(launcher *ServerLauncher) {
		launcher.Start()
	})

	if err != nil {
		fmt.Println("伺服器啟動失敗:", err)
	}
}
