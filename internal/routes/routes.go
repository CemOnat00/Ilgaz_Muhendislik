package routes

import (
	"github.com/cemonat00/ilgaz-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// API Group
	api := r.Group("/api")
	{
		// Public Endpoints
		api.GET("/urunler", handlers.GetProducts)
		api.GET("/urunler/:id", handlers.GetProductByID)
		api.POST("/mesajlar", handlers.CreateMessage)
		api.GET("/projeler", handlers.GetProjects)

		// Admin Auth
		api.POST("/admin/login", handlers.AdminLogin)
		api.POST("/admin/forgot-password", handlers.ForgotPassword)
		api.POST("/admin/reset-password", handlers.ResetPassword)

		// Protected Admin Endpoints (Middleware can be added later)
		admin := api.Group("/admin")
		{
			admin.POST("/upload", handlers.UploadImage)
			admin.GET("/mesajlar", handlers.GetMessages)
			admin.PUT("/mesajlar/:id/read", handlers.MarkMessageRead)
			admin.DELETE("/mesajlar/:id", handlers.DeleteMessage)
			admin.POST("/urunler", handlers.AddProduct)
			admin.PUT("/urunler/:id", handlers.UpdateProduct)
			admin.DELETE("/urunler/:id", handlers.DeleteProduct)
			admin.POST("/projeler", handlers.AddProject)
			admin.PUT("/projeler/:id", handlers.UpdateProject)
			admin.DELETE("/projeler/:id", handlers.DeleteProject)
		}
	}
}
