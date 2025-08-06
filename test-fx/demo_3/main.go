package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// --- 1. Service 和 Handler 層 ---
type UserService struct{}

func NewUserService() *UserService { return &UserService{} }

type UserHandler struct {
	logger  *zap.Logger
	service *UserService
}

func NewUserHandler(logger *zap.Logger, service *UserService) *UserHandler {
	return &UserHandler{logger: logger, service: service}
}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
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

// --- 3. Zap 的 Provider (依環境區分) ---
type LoggerConfig struct {
	LogFile string
}

func NewDevLogger() (*zap.Logger, error) {
	fmt.Println("創建開發環境的 Logger (輸出至控制台)...")
	logger, err := zap.NewDevelopment()
	return logger, err
}

func NewProdLogger(cfg LoggerConfig) (*zap.Logger, error) {
	fmt.Println("創建正式環境的 Logger (輸出至檔案)...")
	logDir := filepath.Dir(cfg.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("無法創建日誌目錄: %w", err)
	}

	zapConfig := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{cfg.LogFile},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("無法建立正式環境 Logger: %w", err)
	}
	return logger, nil
}

// 正確使用 fx.In 結構體方式注入命名 Logger
type LoggerDeps struct {
	fx.In
	DevLogger  *zap.Logger `name:"dev"`
	ProdLogger *zap.Logger `name:"prod"`
}

func NewLoggerForEnv(deps LoggerDeps) *zap.Logger {
	env := os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}
	if env == "dev" {
		return deps.DevLogger
	}
	return deps.ProdLogger
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

// --- 5. Fx 的模組化組件 ---
var LoggerModule = fx.Module(
	"logger",
	fx.Provide(
		fx.Annotate(NewDevLogger, fx.ResultTags(`name:"dev"`)),
		fx.Annotate(NewProdLogger, fx.ResultTags(`name:"prod"`)),
		NewLoggerForEnv,
	),
	fx.Invoke(func(lifecycle fx.Lifecycle, logger *zap.Logger) {
		lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				fmt.Println("正在同步日誌...")
				return logger.Sync()
			},
		})
	}),
)

var UserModule = fx.Module(
	"user",
	fx.Provide(
		NewUserService,
		NewUserHandler,
	),
)

var RouterModule = fx.Module(
	"router",
	fx.Provide(
		fx.Annotate(NewPublicRouter, fx.ResultTags(`name:"public"`)),
		fx.Annotate(NewInternalRouter, fx.ResultTags(`name:"internal"`)),
	),
)

// --- 6. main 組裝與啟動 ---
func main() {
	app := fx.New(
		fx.Supply(
			LoggerConfig{LogFile: "./tmp/app.log"},
		),
		LoggerModule,
		UserModule,
		RouterModule,
		fx.Invoke(StartServers),
	)

	app.Run()
}
