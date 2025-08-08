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

// --- 介面與多個實作 ---
type User interface {
	GetName() string
}

type Service_CCC struct{}

func (s *Service_CCC) GetName() string { return "User from Service CCC" }

type Service_DDD struct{}

func (s *Service_DDD) GetName() string { return "User from Service DDD" }

// --- Handler 與 Server ---
type ServiceMap map[string]User

type UserHandler struct {
	services ServiceMap
}

func NewUserHandler(services ServiceMap) *UserHandler {
	return &UserHandler{services: services}
}

func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := strings.ToUpper(vars["serviceName"])
	userSvc, ok := h.services[serviceName]
	if !ok {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, userSvc.GetName())
}

func NewRouter(h *UserHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user/{serviceName}", h.HandleUser).Methods("GET")
	return router
}

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

// --- 用於 fx.In 具名注入的集合 ---
type UsersCollect struct {
	fx.In
	CCC User `name:"CCC"`
	DDD User `name:"DDD"`
}

// --- 裝飾器 ---
type userDecorator struct{ inner User }

func (d *userDecorator) GetName() string {
	return "[decorated] " + d.inner.GetName()
}

// --- main() 用 name 並 Decorate CCC ---
func main() {
	fx.New(
		// 提供具名的 User
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

		// 組裝 ServiceMap
		fx.Provide(func(users UsersCollect) ServiceMap {
			m := make(ServiceMap)
			m["CCC"] = users.CCC
			m["DDD"] = users.DDD
			return m
		}),

		// 用 fx.Decorate 包裝 name:"CCC" 的 User
		fx.Decorate(
			fx.Annotate(
				func(next User) User {
					return &userDecorator{inner: next}
				},
				fx.ParamTags(`name:"CCC"`),
				fx.ResultTags(`name:"CCC"`),
			),
		),

		// 其他組件
		fx.Provide(NewUserHandler, NewRouter, zap.NewDevelopment),
		fx.Invoke(StartServer),
	).Run()
}
