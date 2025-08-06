package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

// --- 1. Repository：與資料庫互動 ---
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	fmt.Println("創建 UserRepository...")
	return &UserRepository{}
}

// --- 2. Service：業務邏輯 ---
type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	fmt.Println("創建 UserService...")
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id string) string {
	// 這裡應該呼叫 repo 獲取資料
	return fmt.Sprintf("User: %s (from DB)", id)
}

// --- 3. Handler：處理 HTTP 請求 ---
type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	fmt.Println("創建 UserHandler...")
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user := h.service.GetUserByID(userID)
	fmt.Fprintf(w, "Fetched: %s\n", user)
}

// --- 4. Router：路由設定 ---
// 這個 Provider 負責建立並設定 Router
func NewRouter(handler *UserHandler) *mux.Router {
	fmt.Println("創建並設定 Router...")
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	return r
}

// --- 5. 應用程式：啟動伺服器 ---
func main() {
	// 建立 Dig 容器
	container := dig.New()

	// 註冊所有 Provider
	container.Provide(NewUserRepository)
	container.Provide(NewUserService)
	container.Provide(NewUserHandler)
	container.Provide(NewRouter)

	// 使用 Invoke 取得 Router，並啟動 Web 服務
	err := container.Invoke(func(router *mux.Router) {
		fmt.Println("伺服器啟動中，監聽 http://localhost:8080...")
		http.ListenAndServe(":8080", router)
	})

	if err != nil {
		fmt.Println("伺服器啟動失敗:", err)
	}
}
