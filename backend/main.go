package main

// Package Documentation
// This is the main package for the ISX Portfolio backend.
// It initializes the database, sets up the Gin router, and starts the server.
// It also handles Google OAuth and JWT authentication.
// It uses the Gin framework for routing and middleware.
// It uses the GORM ORM for database operations.
// It uses the GORM Migrator for database migrations.
// It uses the GORM Logger for logging database operations.
// It uses the GORM Config for database configuration.
// It uses the GORM Dialect for database dialect.
// It uses the GORM Debugger for debugging database operations.
// It uses the GORM Validator for validating database operations.
// It uses the GORM Validator for validating database operations.

// Importing the necessary packages
import (
	"isxportfolio-backend/config"
	"isxportfolio-backend/handlers"
	"isxportfolio-backend/jobs"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// Main function
// This is the main function for the ISX Portfolio backend.
// It initializes the database, sets up the Gin router, and starts the server.
// It also handles Google OAuth and JWT authentication.
// It uses the Gin framework for routing and middleware.
// It uses the GORM ORM for database operations.
// It uses the GORM Migrator for database migrations.
// It uses the GORM Logger for logging database operations.
// It uses the GORM Config for database configuration.
// It uses the GORM Dialect for database dialect.
// It uses the GORM Debugger for debugging database operations.
// It uses the GORM Validator for validating database operations.

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

	// Update CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize and start the market news job
	newsJob := jobs.NewMarketNewsJob()
	newsJob.Start()
	defer newsJob.Stop()

	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// Auth routes
	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/callback", handlers.GoogleCallback)
	r.GET("/auth/user", handlers.GetCurrentUser)

	// Initialize JWT before starting the server
	config.InitJWT()

	// Setup routes
	setupRoutes(r)

	// Start server
	r.Run(":8000")
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Your existing routes...

		// Market news routes
		market := api.Group("/market")
		{
			newsHandler := handlers.NewMarketNewsHandler()
			market.GET("/news", newsHandler.GetMarketNews)
			market.POST("/news/refresh", newsHandler.RefreshMarketNews)
		}
	}
}
