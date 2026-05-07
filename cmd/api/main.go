package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/cemonat00/ilgaz-backend/internal/routes"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// CORS Middleware (Basic implementation for local dev)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup API Routes
	routes.SetupRoutes(router)

	// Serve Static Files from /public directory
	// Serve individual files at root
	router.StaticFile("/", "./public/index.html")
	router.StaticFile("/index.html", "./public/index.html")
	router.StaticFile("/hakkimizda.html", "./public/hakkimizda.html")
	router.StaticFile("/hizmetler.html", "./public/hizmetler.html")
	router.StaticFile("/urunler.html", "./public/urunler.html")
	router.StaticFile("/urun_detay.html", "./public/urun_detay.html")
	router.StaticFile("/iletisim.html", "./public/iletisim.html")
	
	// Admin pages
	router.StaticFile("/admin-login.html", "./public/admin-login.html")
	router.StaticFile("/admin-dashboard.html", "./public/admin-dashboard.html")
	router.StaticFile("/admin-products.html", "./public/admin-products.html")
	router.StaticFile("/admin-messages.html", "./public/admin-messages.html")
	router.StaticFile("/admin-projects.html", "./public/admin-projects.html")
	router.StaticFile("/admin-settings.html", "./public/admin-settings.html")

	// Serve directories
	router.Static("/css", "./public/css")
	router.Static("/js", "./public/js")
	router.Static("/assets", "./public/assets")

	// Start server on port 8080
	log.Println("Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
