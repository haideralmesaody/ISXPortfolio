package main

import (
	"isxportfolio-backend/config"
	"isxportfolio-backend/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDB()
	config.TestDatabaseWrite()

	// Debug: Print environment variables
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	log.Printf("Client ID: %s", clientID)
	log.Printf("Client Secret length: %d", len(clientSecret))

	// Initialize Google OAuth with environment variables
	config.InitGoogleOAuth(clientID, clientSecret)

	// Create Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// Auth routes
	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/callback", handlers.GoogleCallback)
	r.GET("/auth/user", handlers.GetCurrentUser)

	// Start server
	r.Run(":8000")
}
