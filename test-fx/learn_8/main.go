package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/looplab/fsm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// LightFSMManager 管理多個紅綠燈 FSM 實例
type LightFSMManager struct {
	mu     sync.RWMutex
	fsms   map[string]*fsm.FSM
	logger *zap.Logger
}

func NewLightFSMManager(logger *zap.Logger) *LightFSMManager {
	return &LightFSMManager{
		fsms:   make(map[string]*fsm.FSM),
		logger: logger,
	}
}

// CreateOrGetFSM 取得或建立新的紅綠燈 FSM
func (m *LightFSMManager) CreateOrGetFSM(id string) *fsm.FSM {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f, ok := m.fsms[id]; ok {
		return f
	}

	f := fsm.NewFSM(
		"red",
		fsm.Events{
			{Name: "next", Src: []string{"red"}, Dst: "green"},
			{Name: "next", Src: []string{"green"}, Dst: "yellow"},
			{Name: "next", Src: []string{"yellow"}, Dst: "red"},
		},
		fsm.Callbacks{
			"next": func(ctx context.Context, e *fsm.Event) {
				m.logger.Info("FSM event triggered",
					zap.String("light_id", id),
					zap.String("from", e.Src),
					zap.String("to", e.Dst),
					zap.Any("args", e.Args),
				)
			},
		},
	)
	m.fsms[id] = f
	return f
}

// LightHandler fx.In 結構體，用於注入依賴
type LightHandler struct {
	fx.In
	Manager *LightFSMManager
}

func (h LightHandler) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fsm := h.Manager.CreateOrGetFSM(id)

	fmt.Fprintf(w, "Light %s current state: %s\n", id, fsm.Current())
}

func (h LightHandler) HandleTriggerEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	event := vars["event"]

	fsm := h.Manager.CreateOrGetFSM(id)

	err := fsm.Event(ctx, event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Event failed: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Light %s state changed to: %s\n", id, fsm.Current())
}

// 建立 HTTP Router，fx 會自動注入 LightHandler
func NewRouter(handler LightHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/lights/{id}", handler.HandleGetStatus).Methods("GET")
	r.HandleFunc("/lights/{id}/event/{event}", handler.HandleTriggerEvent).Methods("POST")
	return r
}

// 啟動 HTTP Server
func StartServer(lc fx.Lifecycle, router *mux.Router, logger *zap.Logger) {
	srv := &http.Server{Addr: ":8080", Handler: router}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server on :8080")
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server")
			return srv.Shutdown(ctx)
		},
	})
}

// # 取得紅綠燈狀態 (預設紅燈)
// curl localhost:8080/lights/A
//
// # 觸發 next 事件，狀態會依序切換
// curl -X POST localhost:8080/lights/A/event/next
//
// # 再查一次狀態
// curl localhost:8080/lights/A

func main() {
	fx.New(
		fx.Provide(
			zap.NewDevelopment,
			NewLightFSMManager,
			NewRouter,
		),
		fx.Invoke(StartServer),
	).Run()
}
