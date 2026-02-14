package main

import (
	"log"
	"os"
	"sound-horee-backend/config"
	"sound-horee-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main is the entry point of the Sound Horee Backend Service.
func main() {
	// 1. Bootstrapping environment variables.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	// 2. Initialize database connection pool.
	config.ConnectDatabase()

	// 3. Initialize HTTP Router.
	r := gin.Default()

	// 4. Register API routes.
	routes.SetupRoutes(r)

	// 5. Start Server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Service starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Fatal: Failed to start server: %v", err)
	}
}
