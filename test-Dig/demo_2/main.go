package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

// --- 1. Client 與 Repository 層 ---
// 負責與第三方 API 互動
type ThirdPartyClient struct{}

func NewThirdPartyClient() *ThirdPartyClient {
	fmt.Println("創建 ThirdPartyClient...")
	return &ThirdPartyClient{}
}

func (c *ThirdPartyClient) FetchDataFromAPI(id string) (string, error) {
	// 模擬從第三方 API 取得資料
	return fmt.Sprintf("Data from API for ID: %s", id), nil
}

// 負責與資料庫互動
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	fmt.Println("創建 UserRepository...")
	return &UserRepository{}
}

func (r *UserRepository) ValidateUserInDB(id string) bool {
	// 模擬在資料庫中驗證用戶
	fmt.Printf("在資料庫中驗證用戶 ID: %s...\n", id)
	// 這裡假設所有用戶 ID 100 以上的都存在
	return id >= "100"
}

// --- 2. Service 層 ---
// User Service 現在同時依賴 ThirdPartyClient 和 UserRepository
type UserService struct {
	client *ThirdPartyClient
	repo   *UserRepository
}

// NewUserService 的參數列表現在包含兩個依賴
func NewUserService(client *ThirdPartyClient, repo *UserRepository) *UserService {
	fmt.Println("創建 UserService...")
	return &UserService{client: client, repo: repo}
}

func (s *UserService) GetAndValidateUser(id string) (string, error) {
	// 業務邏輯：先從第三方 API 取得資料
	apiData, err := s.client.FetchDataFromAPI(id)
	if err != nil {
		return "", err
	}

	// 接著用資料庫進行驗證
	isValid := s.repo.ValidateUserInDB(id)
	if !isValid {
		return "", fmt.Errorf("用戶 ID: %s 驗證失敗", id)
	}

	return fmt.Sprintf("API資料: %s, 驗證結果: OK", apiData), nil
}

// --- 3. Handler 層 ---
// Handler 只依賴 Service，不知道 Service 的內部如何實現
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

	user, err := h.service.GetAndValidateUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Fetched: %s\n", user)
}

// --- 4. Router Provider（使用結構體封裝依賴）---
type RouterDependencies struct {
	dig.In
	UserHandler *UserHandler
}

func NewRouter(deps RouterDependencies) *mux.Router {
	fmt.Println("創建並設定 Router...")
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", deps.UserHandler.GetUser).Methods("GET")
	return r
}

// --- 5. 應用程式：組裝與啟動 ---
func main() {
	container := dig.New()

	// 註冊所有 Provider
	container.Provide(NewThirdPartyClient) // 新增的 Provider
	container.Provide(NewUserRepository)
	container.Provide(NewUserService) // Dig 會自動處理這個 Provider 的兩個依賴
	container.Provide(NewUserHandler)
	container.Provide(NewRouter)

	err := container.Invoke(func(router *mux.Router) {
		fmt.Println("伺服器啟動中，監聽 http://localhost:8080...")
		http.ListenAndServe(":8080", router)
	})

	if err != nil {
		fmt.Println("伺服器啟動失敗:", err)
	}
}
