package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aslam/backend/internal/database"
	"github.com/aslam/backend/internal/handlers"
	"github.com/aslam/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	// Database connection
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := database.New(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create initial admin user if no users exist
	users, err := db.ListUsers(nil)
	if err == nil && len(users) == 0 {
		adminEmail := "admin@aslam.local"
		adminPass := "Admin@123456"
		_, err := db.CreateUser(nil, adminEmail, adminPass, "admin")
		if err == nil {
			log.Printf("✓ Initial admin user created: %s / %s", adminEmail, adminPass)
		}
	}

	// Router setup
	router := gin.Default()

	// Middlewares
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Public routes
	public := router.Group("/api/auth")
	{
		public.POST("/login", handlers.Login(db))
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth endpoints
		auth := protected.Group("/auth")
		{
			auth.GET("/profile", handlers.GetProfile(db))
			auth.PUT("/profile", handlers.UpdateProfile(db))
		}

		// User management (admin only)
		users := protected.Group("/users")
		users.Use(middleware.RoleMiddleware("admin"))
		{
			users.POST("", handlers.Register(db))
			users.GET("", handlers.ListUsers(db))
			users.DELETE("/:id", handlers.DeleteUser(db))
		}

		// Product endpoints
		products := protected.Group("/products")
		{
			products.GET("", handlers.ListProducts(db))
			products.GET("/:sku", handlers.GetProduct(db))
			products.POST("", handlers.CreateProduct(db))
			products.PUT("/:sku", handlers.UpdateProduct(db))
			products.DELETE("/:sku", handlers.DeleteProduct(db))
		}

		// Shelf endpoints
		shelves := protected.Group("/shelves")
		{
			shelves.GET("", handlers.ListShelves(db))
			shelves.GET("/:id", handlers.GetShelf(db))
			shelves.POST("", handlers.CreateShelf(db))
			shelves.PUT("/:id", handlers.UpdateShelf(db))
			shelves.DELETE("/:id", handlers.DeleteShelf(db))
			shelves.POST("/:id/items", handlers.AddItemToShelf(db))
			shelves.DELETE("/:id/items/:itemId", handlers.RemoveItemFromShelf(db))
			shelves.PUT("/:id/items/:itemId", handlers.UpdateItemQuantity(db))
		}

		// Health check
		protected.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
