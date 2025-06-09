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

package main

import (
	"fmt"
	"log"
	"os"

	"slack-clone-go-next/router"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	fmt.Println("Hello, World!")
	router := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on port", port)
	router.Run(":" + port)
}
