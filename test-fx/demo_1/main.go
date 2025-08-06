package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap" // 匯入 zap
)

// --- 1. Service 和 Handler 層 ---
type UserService struct{}

func NewUserService() *UserService { return &UserService{} }

// UserHandler 現在依賴 *zap.Logger
type UserHandler struct {
	logger  *zap.Logger
	service *UserService
}

// NewUserHandler 的參數現在需要一個 *zap.Logger
func NewUserHandler(logger *zap.Logger, service *UserService) *UserHandler {
	return &UserHandler{logger: logger, service: service}
}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 在這裡使用 logger 記錄資訊
	h.logger.Info("處理 GetUser 請求", zap.String("path", r.URL.Path))
	fmt.Fprintf(w, "Fetched: User from Fx-powered public API")
}
func (h *UserHandler) GetUserPrivate(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("處理 GetUserPrivate 請求", zap.String("path", r.URL.Path))
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

// --- 3. Zap 的 Provider ---
// 這是我們提供 *zap.Logger 實例的 Provider
func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment() // 使用開發環境的日誌配置
	if err != nil {
		return nil, fmt.Errorf("無法創建 zap logger: %w", err)
	}
	return logger, nil
}

// --- 4. 啟動器 ---
type ServerParams struct {
	fx.In
	Lifecycle      fx.Lifecycle
	PublicRouter   *mux.Router `name:"public"`
	InternalRouter *mux.Router `name:"internal"`
}

func StartServers(p ServerParams) {
	fmt.Println("Fx 啟動器被呼叫...")

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

// --- 5. main 組裝與啟動 ---
func main() {
	fx.New(
		fx.Provide(
			NewLogger, // 註冊 Zap Logger
			NewUserService,
			NewUserHandler, // Fx 會自動注入 NewLogger 提供的 *zap.Logger
			fx.Annotate(NewPublicRouter, fx.ResultTags(`name:"public"`)),
			fx.Annotate(NewInternalRouter, fx.ResultTags(`name:"internal"`)),
		),
		fx.Invoke(StartServers),
	).Run()
}
