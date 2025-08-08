package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// --- 1. Handler 層 ---
type UserService struct{}

func NewUserService() *UserService { return &UserService{} }

type UserHandler struct{ service *UserService }

func NewUserHandler(service *UserService) *UserHandler { return &UserHandler{service: service} }
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fetched: User from Fx-powered public API")
}
func (h *UserHandler) GetUserPrivate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fetched: User from Fx-powered internal API")
}

// --- 2. 兩個 Router Provider ---
type PublicRouterDeps struct {
	fx.In
	UserHandler *UserHandler
}

func NewPublicRouter(deps PublicRouterDeps) *mux.Router {
	fmt.Println("創建 public Router...")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/users/{id}", deps.UserHandler.GetUser).Methods("GET")
	return r
}

type InternalRouterDeps struct {
	fx.In
	UserHandler *UserHandler
}

func NewInternalRouter(deps InternalRouterDeps) *mux.Router {
	fmt.Println("創建 internal Router...")
	r := mux.NewRouter()
	r.HandleFunc("/internal/api/v1/users/{id}", deps.UserHandler.GetUserPrivate).Methods("GET")
	return r
}

// --- 3. 啟動器：這是 Fx 的核心，它處理生命週期 ---

// 這是啟動服務的參數結構體，使用 fx.In 來注入命名依賴
type ServerParams struct {
	// 當你需要「一次接多個，還要有條件（name/group/optional）」的時候，就用 fx.In。
	// 其他情況直接用函式參數就好。
	fx.In
	Lifecycle      fx.Lifecycle
	PublicRouter   *mux.Router `name:"public"`
	InternalRouter *mux.Router `name:"internal"`
}

// StartServers 函數將被 fx.Invoke 呼叫，並接收 ServerParams 結構體
func StartServers(p ServerParams) {
	fmt.Println("Fx 啟動器被呼叫...")

	// 啟動 public server
	publicServer := &http.Server{Addr: ":8080", Handler: p.PublicRouter}
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Public 服務監聽在 :8080...")
			go publicServer.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("正在優雅地關閉 Public 服務...")
			return publicServer.Shutdown(ctx)
		},
	})

	// 啟動 internal server
	internalServer := &http.Server{Addr: ":8081", Handler: p.InternalRouter}
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Internal 服務監聽在 :8081...")
			go internalServer.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("正在優雅地關閉 Internal 服務...")
			return internalServer.Shutdown(ctx)
		},
	})
}

// fx.New 函式的運作原理
// fx.New 函式的作用是建立一個 Fx 應用程式容器（Container）。它會接收一個或多個 fx.Option，這些選項會告訴 Fx 容器要如何組裝你的應用程式。
//
// 你可以把 fx.New 想像成一個總指揮部，而 fx.Option 就像是指令，告訴總指揮部：
//
// fx.Provide: 「請提供這些服務或物件。」
//
// fx.Invoke: 「當所有東西都準備好了，請調用這個函式來啟動應用程式。」
//
// --- 4. main 組裝與啟動 ---
func main() {
	// fx 內建 DIg 、主要負責整個生命週期控管
	fx.New(
		fx.Provide(
			NewUserService,
			NewUserHandler,
			// 在 fx 裡有兩種特殊 tag key：
			//
			// name:"..."
			//
			// 把一個值標記為具名（named value）。
			// 容器裡同一個型別 + 不同 name 可以共存。
			// 消費端必須精確指定 name:"..." 才能拿到。
			//
			// group:"..."
			//
			// 把多個相同型別的值歸入同一群組（value group）。
			// 容器會把同一 group 的所有值自動收集成 slice []T。
			// 順序無保證（unordered）。
			// 消費端必須宣告 []T \group:"..."`` 才會收到整組。
			//
			fx.Annotate(NewPublicRouter, fx.ResultTags(`name:"public"`)),
			fx.Annotate(NewInternalRouter, fx.ResultTags(`name:"internal"`)),
		),
		// 這裡的 fx.Invoke 只呼叫 StartServers 函數
		// Fx 會自動為 StartServers 注入 ServerParams 結構體中的所有依賴
		fx.Invoke(StartServers),
	).Run()
}
