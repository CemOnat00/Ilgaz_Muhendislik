package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ilgaz-backend/internal/models"
)

// --- IN-MEMORY MOCK DATA ---

var mockUrunler = []models.Urun{
	{
		ID:         "1",
		Kategori:   "Hvac",
		Baslik:     "Yüksek Verimli Pompa Sistemi",
		Aciklama:   "Endüstriyel tesisler için enerji tasarruflu akıllı sirkülasyon pompası.",
		Resim:      "https://images.unsplash.com/photo-1581092160562-40aa08e78837?auto=format&fit=crop&q=80",
		Ozellikler: []string{"A Sınıfı Enerji", "Sessiz Çalışma"},
		Durum:      "Aktif",
	},
	{
		ID:         "2",
		Kategori:   "Elektrik",
		Baslik:     "Akıllı Ana Dağıtım Panosu",
		Aciklama:   "Tesisin tüm elektrik yükünü güvenle yöneten, uzaktan izlenebilir akıllı pano sistemi.",
		Resim:      "https://images.unsplash.com/photo-1581092335397-9583eb92d232?auto=format&fit=crop&q=80",
		Ozellikler: []string{"Uzaktan İzleme", "Aşırı Yük Koruması"},
		Durum:      "Aktif",
	},
}

var mockMesajlar = []models.Mesaj{
	{
		ID:       "1",
		AdSoyad:  "Ahmet Yılmaz",
		Email:    "ahmet@yilmaz-const.com",
		Konu:     "Project Inquiry - Site 2A",
		Mesaj:    "I would like to request a quote for our new commercial complex in Bahçelievler.",
		Tarih:    time.Now().Add(-2 * time.Hour),
		OkunduMu: false,
	},
}

// --- HANDLERS ---

// GetUrunler returns all products
func GetUrunler(c *gin.Context) {
	c.JSON(http.StatusOK, mockUrunler)
}

// GetUrun returns a specific product by ID
func GetUrun(c *gin.Context) {
	id := c.Param("id")
	for _, urun := range mockUrunler {
		if urun.ID == id {
			c.JSON(http.StatusOK, urun)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// AddUrun adds a new product
func AddUrun(c *gin.Context) {
	var yeniUrun models.Urun
	if err := c.ShouldBindJSON(&yeniUrun); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	yeniUrun.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	mockUrunler = append([]models.Urun{yeniUrun}, mockUrunler...)
	c.JSON(http.StatusCreated, yeniUrun)
}

// UpdateUrun updates an existing product
func UpdateUrun(c *gin.Context) {
	id := c.Param("id")
	var guncelUrun models.Urun
	if err := c.ShouldBindJSON(&guncelUrun); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, urun := range mockUrunler {
		if urun.ID == id {
			guncelUrun.ID = id
			mockUrunler[i] = guncelUrun
			c.JSON(http.StatusOK, guncelUrun)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// DeleteUrun deletes a product by ID
func DeleteUrun(c *gin.Context) {
	id := c.Param("id")
	for i, urun := range mockUrunler {
		if urun.ID == id {
			// Remove from slice
			mockUrunler = append(mockUrunler[:i], mockUrunler[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// PostMesaj handles incoming contact form submissions
func PostMesaj(c *gin.Context) {
	var yeniMesaj models.Mesaj
	
	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&yeniMesaj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign ID and timestamp
	yeniMesaj.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	yeniMesaj.Tarih = time.Now()
	yeniMesaj.OkunduMu = false

	// Add to in-memory store
	mockMesajlar = append([]models.Mesaj{yeniMesaj}, mockMesajlar...) // Add to beginning

	c.JSON(http.StatusCreated, gin.H{
		"message": "Mesaj başarıyla alındı",
		"data":    yeniMesaj,
	})
}

// GetAdminMesajlar returns all messages for the admin panel
func GetAdminMesajlar(c *gin.Context) {
	c.JSON(http.StatusOK, mockMesajlar)
}

// AdminLogin handles a simple mock authentication
func AdminLogin(c *gin.Context) {
	var loginReq models.LoginRequest
	
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
		return
	}

	// Mock check
	if loginReq.KullaniciAdi == "admin" && loginReq.Sifre == "1234" {
		c.JSON(http.StatusOK, models.LoginResponse{
			Token:   "mock-jwt-token-12345",
			Message: "Giriş başarılı",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Hatalı kullanıcı adı veya şifre"})
}
