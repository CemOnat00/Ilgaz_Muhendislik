package routes

import (
	"github.com/gin-gonic/gin"
	"ilgaz-backend/internal/handlers"
)

// SetupRoutes configures the API endpoints
func SetupRoutes(router *gin.Engine) {
	
	// API Group
	api := router.Group("/api")
	{
		// Public Endpoints
		api.GET("/urunler", handlers.GetUrunler)
		api.GET("/urunler/:id", handlers.GetUrun)
		api.POST("/mesajlar", handlers.PostMesaj)

		// Admin Endpoints
		admin := api.Group("/admin")
		{
			admin.POST("/login", handlers.AdminLogin)
			admin.GET("/mesajlar", handlers.GetAdminMesajlar)
		}
	}
}
