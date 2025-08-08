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

// --- main() 用 group 自動收集 ---
func main() {
	fx.New(
		// 兩個 User provider 放入同一個 group
		fx.Provide(
			fx.Annotate(
				func() User { return &Service_CCC{} },
				fx.As(new(User)),
				fx.ResultTags(`group:"users"`),
			),
			fx.Annotate(
				func() User { return &Service_DDD{} },
				fx.As(new(User)),
				fx.ResultTags(`group:"users"`),
			),
		),

		// 從 group 收集並組成 ServiceMap
		fx.Provide(
			fx.Annotate(
				func(users []User) ServiceMap {
					m := make(ServiceMap)
					for _, u := range users {
						parts := strings.Fields(u.GetName())
						key := strings.ToUpper(parts[len(parts)-1])
						m[key] = u
					}
					return m
				},
				fx.ParamTags(`group:"users"`),
			),
		),

		// 其他組件
		fx.Provide(NewUserHandler, NewRouter, zap.NewDevelopment),
		fx.Invoke(StartServer),
	).Run()
}
