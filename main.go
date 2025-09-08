// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/keshav7976/ecommerce/config"
	"github.com/keshav7976/ecommerce/routes"
	"github.com/keshav7976/ecommerce/utils"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the JWT key from environment variable
	utils.InitJWT()

	// Initialize the database connection
	config.ConnectDB()

	router := mux.NewRouter()

	// Register all routes from the routes package
	routes.RegisterRoutes(router)

	// Get allowed origins from environment variable
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Server listening on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
