package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

// --- 1. Repository(DB) 層 ---
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	fmt.Println("創建 UserRepository...")
	return &UserRepository{}
}

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	fmt.Println("創建 ProductRepository...")
	return &ProductRepository{}
}

// --- 2. Service 層 ---
type UserService struct{ repo *UserRepository }

func NewUserService(repo *UserRepository) *UserService {
	fmt.Println("創建 UserService...")
	return &UserService{repo: repo}
}
func (s *UserService) GetUserByID(id string) string {
	return fmt.Sprintf("User: %s (from DB)", id)
}

type ProductService struct{ repo *ProductRepository }

func NewProductService(repo *ProductRepository) *ProductService {
	fmt.Println("創建 ProductService...")
	return &ProductService{repo: repo}
}
func (s *ProductService) GetProductByID(id string) string {
	return fmt.Sprintf("Product: %s (from inventory)", id)
}

// --- 3. Handler 層 ---
type UserHandler struct{ service *UserService }

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

type ProductHandler struct{ service *ProductService }

func NewProductHandler(service *ProductService) *ProductHandler {
	fmt.Println("創建 ProductHandler...")
	return &ProductHandler{service: service}
}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	product := h.service.GetProductByID(productID)
	fmt.Fprintf(w, "Fetched: %s\n", product)
}

// --- 4. Router Provider 的優化 ---

// 這是新的 `RouterDependencies` 結構體。
// 它可以容納任意數量的 Handler，且不會改變 NewRouter 的簽名。
// 這是一個 "in" 結構體，Dig 會自動注入它的欄位。
type RouterDependencies struct {
	dig.In
	UserHandler    *UserHandler
	ProductHandler *ProductHandler
}

// NewRouter 的簽名現在只需要一個參數：RouterDependencies。
// Dig 會負責解析這個結構體並注入所有必要的 Handler。
func NewRouter(deps RouterDependencies) *mux.Router {
	fmt.Println("創建並設定 Router...")
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", deps.UserHandler.GetUser).Methods("GET")
	r.HandleFunc("/products/{id}", deps.ProductHandler.GetProduct).Methods("GET")
	return r
}

// --- 5. 應用程式：組裝與啟動 ---

func main() {
	container := dig.New()

	// 註冊所有 Provider
	container.Provide(NewUserRepository)
	container.Provide(NewProductRepository)
	container.Provide(NewUserService)
	container.Provide(NewProductService)
	container.Provide(NewUserHandler)
	container.Provide(NewProductHandler)
	container.Provide(NewRouter)

	err := container.Invoke(func(router *mux.Router) {
		fmt.Println("伺服器啟動中，監聽 http://localhost:8080...")
		http.ListenAndServe(":8080", router)
	})

	if err != nil {
		fmt.Println("伺服器啟動失敗:", err)
	}
}
