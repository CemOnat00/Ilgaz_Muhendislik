package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/cemonat00/ilgaz-backend/internal/database"
	"github.com/cemonat00/ilgaz-backend/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 0. Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Uyarı: .env dosyası bulunamadı, sistem değişkenleri kullanılacak")
	}

	// 1. Initialize Database
	database.ConnectDB()
	database.SeedAdmin()
	database.SeedProducts()
	database.SeedSettings()

	// 2. Setup Gin Router
	r := gin.Default()

	// HTML sayfaları için tarayıcı önbelleğini kapat — her zaman güncel sürüm sunulur
	r.Use(func(c *gin.Context) {
		p := c.Request.URL.Path
		if p == "/" || p == "/admin" || strings.HasSuffix(p, ".html") {
			c.Header("Cache-Control", "no-cache, must-revalidate")
		}
		c.Next()
	})

	// 3. Static Files (Frontend)
	r.Static("/css", "./public/css")
	r.Static("/js", "./public/js")
	r.Static("/assets", "./assets")
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/index.html", "./public/index.html")
	r.StaticFile("/hakkimizda.html", "./public/hakkimizda.html")
	r.StaticFile("/hizmetler.html", "./public/hizmetler.html")
	r.StaticFile("/urunler.html", "./public/urunler.html")
	r.StaticFile("/urun_detay.html", "./public/urun_detay.html")
	r.StaticFile("/iletisim.html", "./public/iletisim.html")
	
	// Admin Pages
	r.StaticFile("/admin", "./public/admin-login.html")
	r.StaticFile("/admin-login.html", "./public/admin-login.html")
	r.StaticFile("/admin-forgot-password.html", "./public/admin-forgot-password.html")
	r.StaticFile("/admin-reset-password.html", "./public/admin-reset-password.html")
	r.StaticFile("/admin-dashboard.html", "./public/admin-dashboard.html")
	r.StaticFile("/admin-products.html", "./public/admin-products.html")
	r.StaticFile("/admin-messages.html", "./public/admin-messages.html")
	r.StaticFile("/admin-reviews.html", "./public/admin-reviews.html")
	r.StaticFile("/admin-settings.html", "./public/admin-settings.html")
	r.StaticFile("/admin-leads.html", "./public/admin-leads.html")

	// 4. Setup API Routes
	routes.SetupRoutes(r)

	// 5. Start Server
	port := ":8080"
	log.Printf("Sunucu %s portunda başlatılıyor...", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
