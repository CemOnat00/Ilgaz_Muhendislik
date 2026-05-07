package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cemonat00/ilgaz-backend/internal/models"
)

// --- IN-MEMORY MOCK DATA ---

var mockUrunler = []models.Urun{
	{
		ID:          "1",
		Kategori:    "Vana Grubu",
		Baslik:      "Küresel Vana - Tip 101",
		Isim:        "Endüstriyel Küresel Vana",
		Aciklama:    "Yüksek basınçlı buhar hatları için özel sızdırmazlık teknolojisi.",
		Resim:       "https://images.unsplash.com/photo-1581092160562-40aa08e78837?auto=format&fit=crop&q=80",
		Fiyat:       1250.00,
		StoktaVarMi: true,
		Ozellikler:  []string{"Endüstriyel Tip", "Paslanmaz Çelik"},
		Durum:       "Aktif",
	},
	{
		ID:          "2",
		Kategori:    "Pompa Sistemleri",
		Baslik:      "Sirkülasyon Pompası XP",
		Isim:        "XP Serisi Sirkülasyon Pompası",
		Aciklama:    "Isıtma ve soğutma sistemlerinde yüksek verimli enerji tasarrufu.",
		Resim:       "https://images.unsplash.com/photo-1581092335397-9583eb92d232?auto=format&fit=crop&q=80",
		Fiyat:       3400.00,
		StoktaVarMi: true,
		Ozellikler:  []string{"Yeni Ürün", "A Sınıfı Enerji"},
		Durum:       "Aktif",
	},
	{
		ID:          "3",
		Kategori:    "Isı Değiştiriciler",
		Baslik:      "Plakalı Eşanjör - HE-50",
		Isim:        "Yüksek Verimli Plakalı Eşanjör",
		Aciklama:    "Kompakt tasarım ile maksimum ısı transfer verimliliği.",
		Resim:       "https://images.unsplash.com/photo-1581092430616-94e671651e14?auto=format&fit=crop&q=80",
		Fiyat:       5600.00,
		StoktaVarMi: false,
		Ozellikler:  []string{"Endüstriyel Tip"},
		Durum:       "Aktif",
	},
	{
		ID:          "4",
		Kategori:    "Otomasyon",
		Baslik:      "PLC Kontrol Paneli",
		Isim:        "Akıllı Proses Kontrol Ünitesi",
		Aciklama:    "Fabrika otomasyonu için programlanabilir mantıksal denetleyici.",
		Resim:       "https://images.unsplash.com/photo-1581092580497-e0d23cbdf1dc?auto=format&fit=crop&q=80",
		Fiyat:       8900.00,
		StoktaVarMi: true,
		Ozellikler:  []string{"Yeni Ürün", "Endüstriyel Tip"},
		Durum:       "Aktif",
	},
	{
		ID:          "5",
		Kategori:    "Yedek Parça",
		Baslik:      "Conta Takımı - Viton",
		Isim:        "Yüksek Isı Dayanımlı Conta Takımı",
		Aciklama:    "Kimyasal dayanımı yüksek, sızdırmazlık garantili conta seti.",
		Resim:       "https://images.unsplash.com/photo-1581092334651-ddf26d9a1930?auto=format&fit=crop&q=80",
		Fiyat:       450.00,
		StoktaVarMi: true,
		Ozellikler:  []string{"Endüstriyel Tip"},
		Durum:       "Aktif",
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

var mockProjeler = []models.Proje{
	{
		ID:       "1",
		Baslik:   "Modern Konut Kompleksi",
		Musteri:  "Yılmaz İnşaat",
		Kategori: "Elektrik & Mekanik",
		Durum:    "Devam Ediyor",
		Tarih:    time.Now(),
		Aciklama: "500 dairelik konut projesinin tüm elektrik altyapı işleri.",
	},
}

// --- HANDLERS ---

// GetProjeler returns all projects
func GetProjeler(c *gin.Context) {
	c.JSON(http.StatusOK, mockProjeler)
}

// AddProje adds a new project
func AddProje(c *gin.Context) {
	var yeniProje models.Proje
	if err := c.ShouldBindJSON(&yeniProje); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	yeniProje.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	yeniProje.Tarih = time.Now()
	mockProjeler = append([]models.Proje{yeniProje}, mockProjeler...)
	c.JSON(http.StatusCreated, yeniProje)
}

// UpdateProje updates an existing project
func UpdateProje(c *gin.Context) {
	id := c.Param("id")
	var guncelProje models.Proje
	if err := c.ShouldBindJSON(&guncelProje); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, p := range mockProjeler {
		if p.ID == id {
			guncelProje.ID = id
			mockProjeler[i] = guncelProje
			c.JSON(http.StatusOK, guncelProje)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
}

// DeleteProje deletes a project
func DeleteProje(c *gin.Context) {
	id := c.Param("id")
	for i, p := range mockProjeler {
		if p.ID == id {
			mockProjeler = append(mockProjeler[:i], mockProjeler[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
}

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
