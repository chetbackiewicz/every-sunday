package main

import (
	"log"
	"os"

	"github.com/chetbackiewicz/every-sunday-backend/internal/database"
	"github.com/chetbackiewicz/every-sunday-backend/internal/handlers"
	"github.com/chetbackiewicz/every-sunday-backend/internal/middleware"
	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_NAME", "every_sunday_dev"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database connection
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repos := repositories.NewRepositories(db)

	// Initialize handlers
	handlers := handlers.NewHandlers(repos)

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	setupRoutes(r, handlers)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine, h *handlers.Handlers) {
	// CORS middleware
	r.Use(middleware.CORSMiddleware())

	// API v1 routes
	api := r.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/refresh", h.Auth.RefreshToken)
		auth.POST("/logout", h.Auth.Logout)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// User info
		protected.GET("/auth/me", h.Auth.Me)

		// Budget routes
		budgets := protected.Group("/budgets")
		{
			budgets.GET("", h.Budget.List)
			budgets.GET("/:month", h.Budget.Get)
			budgets.POST("", h.Budget.Create)
			budgets.PUT("/:month", h.Budget.Update)
			budgets.DELETE("/:month", h.Budget.Delete)
		}

		// Category routes
		categories := protected.Group("/budgets/:month/categories")
		{
			categories.GET("", h.Category.List)
			categories.POST("", h.Category.Create)
			categories.PUT("/:id", h.Category.Update)
			categories.DELETE("/:id", h.Category.Delete)
		}

		// CSV file upload
		protected.POST("/budgets/:month/files", h.File.Upload)

		// Transaction routes
		transactions := protected.Group("/budgets/:month/transactions")
		{
			transactions.GET("", h.Transaction.List)
			transactions.PUT("/:id", h.Transaction.Update)
			transactions.DELETE("/:id", h.Transaction.Delete)
		}

		// Cell reference routes
		references := protected.Group("/budgets/:month/references")
		{
			references.POST("", h.CellReference.Create)
			references.DELETE("/:id", h.CellReference.Delete)
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
