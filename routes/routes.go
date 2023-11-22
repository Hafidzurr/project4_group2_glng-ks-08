package routes

import (
	"fmt"
	"net/http"

	"github.com/Hafidzurr/project4_group2_glng-ks-08/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	// Handler baru untuk root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API Project 4 Kelompok 2")
	})
	// User routes
	router.HandleFunc("/users/register", controllers.RegisterUser(db)).Methods("POST")
	router.HandleFunc("/users/login", controllers.LoginUser(db)).Methods("POST")
	router.HandleFunc("/users/topup", controllers.TopUpUser(db)).Methods("PATCH")

	// Category routes
	router.HandleFunc("/categories", controllers.CreateCategory(db)).Methods("POST")
	router.HandleFunc("/categories", controllers.GetCategories(db)).Methods("GET")
	router.HandleFunc("/categories/{categoryId}", controllers.UpdateCategory(db)).Methods("PATCH")
	router.HandleFunc("/categories/{categoryId}", controllers.DeleteCategory(db)).Methods("DELETE")

	// Product routes
	router.HandleFunc("/products", controllers.CreateProduct(db)).Methods("POST")
	router.HandleFunc("/products", controllers.GetProducts(db)).Methods("GET")
	router.HandleFunc("/products/{productId}", controllers.UpdateProduct(db)).Methods("PUT")
	router.HandleFunc("/products/{productId}", controllers.DeleteProduct(db)).Methods("DELETE")

	// TransactionHistory routes
	router.HandleFunc("/transactions", controllers.CreateTransaction(db)).Methods("POST")
	router.HandleFunc("/transactions/my-transactions", controllers.GetMyTransactions(db)).Methods("GET")
	router.HandleFunc("/transactions/user-transactions", controllers.GetUserTransactions(db)).Methods("GET")
}
