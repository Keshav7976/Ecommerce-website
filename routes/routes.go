// File: routes/routes.go
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/keshav7976/ecommerce/controllers"
	"github.com/keshav7976/ecommerce/middleware"
)

func RegisterRoutes(r *mux.Router) {
	// Auth
	r.HandleFunc("/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	// Items (Public & Protected)
	r.HandleFunc("/items", controllers.GetItems).Methods("GET")
	r.Handle("/items", middleware.JWTAuth(http.HandlerFunc(controllers.CreateItem))).Methods("POST")
	r.Handle("/items/{id}", middleware.JWTAuth(http.HandlerFunc(controllers.UpdateItem))).Methods("PUT")
	r.Handle("/items/{id}", middleware.JWTAuth(http.HandlerFunc(controllers.DeleteItem))).Methods("DELETE")

	// Cart (Protected)
	r.Handle("/cart/add", middleware.JWTAuth(http.HandlerFunc(controllers.AddToCart))).Methods("POST")
	r.Handle("/cart", middleware.JWTAuth(http.HandlerFunc(controllers.GetCart))).Methods("GET")
	r.Handle("/cart/remove", middleware.JWTAuth(http.HandlerFunc(controllers.RemoveFromCart))).Methods("DELETE")

	// Categories
	r.HandleFunc("/categories", controllers.GetCategories).Methods("GET")
}