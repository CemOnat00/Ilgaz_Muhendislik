package routes

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/cemonat00/ilgaz-backend/internal/handlers"
)

// AuthMiddleware checks for a valid mock token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		
		// In a real app, you would validate a JWT here
		if token == "" || !strings.HasPrefix(token, "Bearer mock-jwt-token") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bu işlem için yetkiniz yok"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetupRoutes configures the API endpoints
func SetupRoutes(router *gin.Engine) {
	
	// API Group
	api := router.Group("/api")
	{
		// Public Endpoints
		api.GET("/urunler", handlers.GetUrunler)
		api.GET("/urunler/:id", handlers.GetUrun)
		api.GET("/projeler", handlers.GetProjeler)
		api.POST("/mesajlar", handlers.PostMesaj)

		// Admin Login (Public)
		api.POST("/admin/login", handlers.AdminLogin)

		// Protected Admin Endpoints
		admin := api.Group("/admin")
		admin.Use(AuthMiddleware())
		{
			admin.GET("/mesajlar", handlers.GetAdminMesajlar)
			
			// Product CRUD
			admin.POST("/urunler", handlers.AddUrun)
			admin.PUT("/urunler/:id", handlers.UpdateUrun)
			admin.DELETE("/urunler/:id", handlers.DeleteUrun)

			// Project CRUD
			admin.POST("/projeler", handlers.AddProje)
			admin.PUT("/projeler/:id", handlers.UpdateProje)
			admin.DELETE("/projeler/:id", handlers.DeleteProje)
		}
	}
}
