package router

import (
	"go_learn/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/product/{id}", middleware.GetProduct).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/product", middleware.GetAllProduct).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newproduct", middleware.CreateProduct).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/product/{id}", middleware.UpdateProduct).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteproduct/{id}", middleware.DeleteProduct).Methods("DELETE", "OPTIONS")

	return router
}
