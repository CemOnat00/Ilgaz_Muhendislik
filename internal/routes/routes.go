package routes

import (
	"github.com/cemonat00/ilgaz-backend/internal/handlers"
	"github.com/cemonat00/ilgaz-backend/internal/middleware"
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
		api.GET("/ayarlar", handlers.GetSettings)
		api.POST("/katalog/indir", handlers.DownloadCatalog)

		// Yorum (Review) Public Endpoints
		api.POST("/yorumlar", handlers.CreateReview)
		api.GET("/yorumlar/:productId", handlers.GetProductReviews)
		api.POST("/yorumlar/:id/oy", handlers.VoteHelpful)

		// Admin Auth
		api.POST("/admin/login", handlers.AdminLogin)
		api.POST("/admin/forgot-password", handlers.ForgotPassword)
		api.POST("/admin/reset-password", handlers.ResetPassword)

		// Protected Admin Endpoints (Middleware applied)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired())
		{
			admin.POST("/upload", handlers.UploadImage)
			admin.GET("/mesajlar", handlers.GetMessages)
			admin.PUT("/mesajlar/:id/read", handlers.MarkMessageRead)
			admin.DELETE("/mesajlar/:id", handlers.DeleteMessage)
			admin.POST("/mesajlar/:id/reply", handlers.ReplyMessage)
			admin.GET("/urunler", handlers.GetAllProducts)
			admin.POST("/urunler", handlers.AddProduct)
			admin.PUT("/urunler/reorder", handlers.ReorderProducts)
			admin.PUT("/urunler/:id", handlers.UpdateProduct)
			admin.DELETE("/urunler/:id", handlers.DeleteProduct)
			admin.PUT("/ayarlar", handlers.UpdateSettings)
			admin.GET("/leads", handlers.GetLeads)
			admin.DELETE("/leads/:id", handlers.DeleteLead)

			// Yorum (Review) Admin Endpoints
			admin.GET("/yorumlar", handlers.GetAllReviews)
			admin.PUT("/yorumlar/:id/onayla", handlers.ApproveReview)
			admin.PUT("/yorumlar/:id/yanit", handlers.ReplyReview)
			admin.DELETE("/yorumlar/:id", handlers.DeleteReview)
		}
	}
}

