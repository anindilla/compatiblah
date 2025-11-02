package main

import (
	"log"
	"os"
	"compatiblah/backend/db"
	"compatiblah/backend/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	dbPath := "compatiblah.db"
	if err := db.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Check for Gemini API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	// Setup Gin router
	r := gin.Default()

	// Manual CORS middleware - handle ALL requests including OPTIONS preflight
	r.Use(func(c *gin.Context) {
		// Set CORS headers for all origins (production-friendly)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept, Authorization, X-Requested-With, X-CSRF-Token")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("Access-Control-Max-Age", "43200")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// Handle OPTIONS preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Also use gin-cors as backup (redundant but ensures compatibility)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"X-CSRF-Token",
	}
	config.AllowCredentials = false
	config.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	config.MaxAge = 12 * 60 * 60 // 12 hours
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api")
	{
		api.POST("/assess", handlers.AssessCompatibility)
		api.POST("/assess/category", handlers.AssessCategory)
		api.GET("/assessment/:id", handlers.GetAssessment)
		api.GET("/assessments", handlers.GetAllAssessments)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Root endpoint - helpful message
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Compatiblah API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health": "/health",
				"assess": "POST /api/assess",
				"get_assessment": "GET /api/assessment/:id",
				"get_all": "GET /api/assessments",
			},
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

