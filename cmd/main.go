// @title Slack Clone API
// @version 1.0
// @description This is a backend for a Slack clone with chat, audio/video calling.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @contact.name API Support
// @contact.email support@slackclone.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"
	"os"

	_ "slack-clone-go-next/docs" // Import generated docs
	"slack-clone-go-next/internal/database"
	"slack-clone-go-next/router"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting server...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}

	// Initialize database (skip if SKIP_DB=true for testing)
	if os.Getenv("SKIP_DB") != "true" {
		if err := database.InitDB(); err != nil {
			log.Printf("Database connection failed: %v", err)
			log.Println("To skip database and test server only, set SKIP_DB=true in environment")
			log.Fatal("Exiting due to database connection failure")
		}
		defer database.CloseDB()
		fmt.Println("Database connected successfully")
	} else {
		fmt.Println("Skipping database connection (SKIP_DB=true)")
	}

	router := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	fmt.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html\n", port)
	router.Run(":" + port)
}
