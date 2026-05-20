package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cemonat00/ilgaz-backend/internal/catalog"
	"github.com/cemonat00/ilgaz-backend/internal/database"
	"github.com/cemonat00/ilgaz-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// --- UPLOAD HANDLER ---

func UploadImage(c *gin.Context) {
	// Enforce 10MB file size limit
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10 MB

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya yüklenemedi veya çok büyük (Max 10MB)"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz dosya formatı. Sadece JPG, PNG ve WEBP izinlidir."})
		return
	}

	// Create unique filename
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
	savePath := filepath.Join("assets", "uploads", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi"})
		return
	}

	// Return the accessible URL
	url := fmt.Sprintf("/assets/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// --- PRODUCT HANDLERS ---

// GetAllProducts — Admin için tüm ürünleri döndürür (Pasif dahil)
func GetAllProducts(c *gin.Context) {
	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var products []models.Product
	opts := options.Find().SetSort(bson.M{"order": 1})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürünler çekilemedi"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

func GetProducts(c *gin.Context) {
	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var products []models.Product
	options := options.Find().SetSort(bson.M{"order": 1})
	// Sadece "Aktif" ürünleri döndür — "Pasif" ürünler müşteri kataloğunda gizlenir
	cursor, err := collection.Find(ctx, bson.M{"status": "Aktif"}, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürünler çekilemedi"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ürün ID formatı"})
		return
	}

	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func AddProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.ID = bson.NewObjectID()
	product.CreatedAt = time.Now()

	// Ensure Name and Category are synced for consistency
	if product.Name == "" && product.Isim != "" {
		product.Name = product.Isim
	}
	if product.Category == "" && product.Kategori != "" {
		product.Category = product.Kategori
	}

	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün eklenemedi"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ürün ID formatı"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":              product.Isim,
			"baslik":            product.Baslik,
			"isim":              product.Isim,
			"category":          product.Kategori,
			"kategori":          product.Kategori,
			"price":             product.Price,
			"stock_status":      product.StockStatus,
			"image_url":         product.ImageURL,
			"description":       product.Description,
			"status":            product.Status,
			"images":            product.Images,
			"videos":            product.Videos,
			"technical_specs":   product.TechnicalSpecs,
			"feature_boxes":     product.FeatureBoxes,
			"application_areas": product.ApplicationAreas,
			"is_industrial":     product.IsIndustrial,
			"is_home_appliance": product.IsHomeAppliance,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Güncellenecek ürün bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ürün başarıyla güncellendi"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ürün ID formatı"})
		return
	}

	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silinecek ürün bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ürün başarıyla silindi"})
}

func ReorderProducts(c *gin.Context) {
	var req struct {
		Order []string `json:"order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	collection := database.GetCollection("products")
	for i, idStr := range req.Order {
		objID, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			continue
		}
		_, _ = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"order": i}})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sıralama güncellendi"})
}

// --- MESSAGE HANDLERS ---

func CreateMessage(c *gin.Context) {
	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz mesaj verisi"})
		return
	}

	msg.ID = bson.NewObjectID()
	msg.CreatedAt = time.Now()
	msg.Status = "Unread"

	collection := database.GetCollection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mesaj gönderilemedi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Mesajınız başarıyla iletildi"})
}

func GetMessages(c *gin.Context) {
	collection := database.GetCollection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var messages []models.Message
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mesajlar çekilemedi"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var m models.Message
		cursor.Decode(&m)
		m.Tarih = m.CreatedAt // Duplicate for frontend compatibility
		messages = append(messages, m)
	}

	c.JSON(http.StatusOK, messages)
}

func MarkMessageRead(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz mesaj ID formatı"})
		return
	}

	collection := database.GetCollection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"is_read": true, "status": "Read"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncellenemedi: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mesaj bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Okundu işaretlendi"})
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz mesaj ID formatı"})
		return
	}

	collection := database.GetCollection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silinecek mesaj bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesaj silindi"})
}

func ReplyMessage(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz mesaj ID"})
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	collection := database.GetCollection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var msg models.Message
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&msg)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mesaj bulunamadı"})
		return
	}

	// E-posta gönderimi
	err = sendEmail(msg.Email, "Re: "+msg.Subject, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "E-posta gönderilemedi: " + err.Error()})
		return
	}

	// Mesaj durumunu güncelle
	_, _ = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": "Replied", "okundu_mu": true}})

	c.JSON(http.StatusOK, gin.H{"message": "Yanıt başarıyla gönderildi"})
}

func sendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT") // Örn: "587"

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		// SMTP ayarları yoksa hata vermeden logla (Geliştirme aşaması için)
		fmt.Printf("SMTP ayarları eksik. Gönderilecek mail:\nTo: %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	// Email header ve body
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	return err
}

// --- PROJECT HANDLERS ---

func GetProjects(c *gin.Context) {
	collection := database.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var projects []models.Project
	options := options.Find().SetSort(bson.M{"order": 1})
	cursor, err := collection.Find(ctx, bson.M{}, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Projeler çekilemedi"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Project
		cursor.Decode(&p)
		projects = append(projects, p)
	}

	c.JSON(http.StatusOK, projects)
}

func AddProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = bson.NewObjectID()
	project.Date = time.Now()

	collection := database.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proje eklenemedi"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz proje ID formatı"})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       project.Title,
			"customer":    project.Customer,
			"category":    project.Category,
			"subtitle":    project.SubTitle,
			"progress":    project.Progress,
			"status":      project.Status,
			"description": project.Description,
			"image_url":   project.ImageURL,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Güncellenecek proje bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proje başarıyla güncellendi"})
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz proje ID formatı"})
		return
	}

	collection := database.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silinecek proje bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proje başarıyla silindi"})
}

func ReorderProjects(c *gin.Context) {
	var req struct {
		Order []string `json:"order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	collection := database.GetCollection("projects")
	for i, idStr := range req.Order {
		objID, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			continue
		}
		_, _ = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"order": i}})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sıralama güncellendi"})
}

// --- AUTH HANDLERS ---

func AdminLogin(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Eksik bilgi"})
		return
	}

	collection := database.GetCollection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var admin models.Admin
	err := collection.FindOne(ctx, bson.M{"username": loginReq.Username}).Decode(&admin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz kullanıcı adı veya şifre"})
		return
	}

	// Secure bcrypt password check
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz kullanıcı adı veya şifre"})
		return
	}

	// Generate real JWT token
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-key-for-dev"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": admin.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulamadı"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:   tokenString,
		Message: "Giriş başarılı",
	})
}

func ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz email"})
		return
	}

	collection := database.GetCollection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var admin models.Admin
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&admin)
	if err != nil {
		// Don't reveal if email exists for security, but for this project we'll return success
		c.JSON(http.StatusOK, gin.H{"message": "Sıfırlama linki gönderildi"})
		return
	}

	// Generate mock reset token
	resetToken := "token-" + bson.NewObjectID().Hex()
	_, _ = collection.UpdateOne(ctx, bson.M{"_id": admin.ID}, bson.M{"$set": bson.M{"reset_token": resetToken}})

	// In a real app, send email. Here we return the URL for testing.
	resetURL := "/admin-reset-password.html?token=" + resetToken

	c.JSON(http.StatusOK, gin.H{
		"message":   "Sıfırlama linki gönderildi",
		"reset_url": resetURL, // Dev mode convenience
	})
}

func ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	collection := database.GetCollection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var admin models.Admin
	err := collection.FindOne(ctx, bson.M{"reset_token": req.Token}).Decode(&admin)
	if err != nil || req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veya süresi dolmuş token"})
		return
	}

	// Update password and clear token
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	_, err = collection.UpdateOne(ctx, bson.M{"_id": admin.ID}, bson.M{
		"$set":   bson.M{"password": string(hashedPassword)},
		"$unset": bson.M{"reset_token": ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Şifre başarıyla güncellendi"})
}

// --- SETTINGS HANDLERS ---

// GetSettings — Site ayarlarını döndürür (Public). Ayar dokümanı yoksa varsayılanları verir.
func GetSettings(c *gin.Context) {
	collection := database.GetCollection("settings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var settings models.SiteSettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		c.JSON(http.StatusOK, models.DefaultSettings())
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings — Site ayarlarını günceller (Admin). Tek ayar dokümanını upsert eder.
func UpdateSettings(c *gin.Context) {
	var settings models.SiteSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ayar verisi"})
		return
	}

	if settings.MapZoom < 1 || settings.MapZoom > 21 {
		settings.MapZoom = 16
	}

	collection := database.GetCollection("settings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"company_name":  settings.CompanyName,
			"email":         settings.Email,
			"phone":         settings.Phone,
			"support_email": settings.SupportEmail,
			"address":       settings.Address,
			"map_location":  settings.MapLocation,
			"map_zoom":      settings.MapZoom,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, bson.M{}, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ayarlar kaydedilemedi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ayarlar başarıyla kaydedildi"})
}

// --- CATALOG & LEAD HANDLERS ---

// slugify Türkçe metni dosya adı için ascii slug'a çevirir.
func slugify(s string) string {
	repl := strings.NewReplacer(
		"ç", "c", "Ç", "c", "ğ", "g", "Ğ", "g", "ı", "i", "İ", "i",
		"ö", "o", "Ö", "o", "ş", "s", "Ş", "s", "ü", "u", "Ü", "u", " ", "-",
	)
	s = strings.ToLower(repl.Replace(s))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// DownloadCatalog — Lead bilgisini kaydeder, IP rate-limit uygular ve PDF kataloğu döndürür.
func DownloadCatalog(c *gin.Context) {
	var req struct {
		Name     string `json:"ad_soyad"`
		Phone    string `json:"telefon"`
		Email    string `json:"email"`
		Category string `json:"kategori"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz form verisi"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.Category = strings.TrimSpace(req.Category)

	if req.Name == "" || req.Phone == "" || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ad-soyad, telefon ve e-posta zorunludur"})
		return
	}
	if !strings.Contains(req.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçerli bir e-posta adresi giriniz"})
		return
	}

	ip := c.ClientIP()
	leadsCol := database.GetCollection("leads")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// IP rate-limit: aynı IP'den son 20 saniyede kayıt → engelle (bot / hızlı tıklama)
	if recent, _ := leadsCol.CountDocuments(ctx, bson.M{
		"ip_address": ip,
		"created_at": bson.M{"$gt": time.Now().Add(-20 * time.Second)},
	}); recent > 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Çok sık deneme yaptınız, lütfen birkaç saniye bekleyin."})
		return
	}
	// Saatlik üst sınır: aynı IP'den son 1 saatte 20+ indirme → engelle
	if hourly, _ := leadsCol.CountDocuments(ctx, bson.M{
		"ip_address": ip,
		"created_at": bson.M{"$gt": time.Now().Add(-1 * time.Hour)},
	}); hourly >= 20 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "İndirme limitine ulaşıldı, lütfen daha sonra tekrar deneyin."})
		return
	}

	category := req.Category
	if category == "" {
		category = "Genel"
	}

	// Aktif ürünleri çek ve gerekiyorsa kategoriye göre filtrele
	prodCol := database.GetCollection("products")
	cursor, err := prodCol.Find(ctx, bson.M{"status": "Aktif"}, options.Find().SetSort(bson.M{"order": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürünler okunamadı"})
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	for cursor.Next(ctx) {
		var p models.Product
		if cursor.Decode(&p) == nil {
			if category == "Genel" || p.Kategori == category {
				products = append(products, p)
			}
		}
	}

	pdfBytes, err := catalog.Generate(products, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Katalog oluşturulamadı"})
		return
	}

	// Lead kaydı — her başarılı indirme bir potansiyel müşteri / indirme olayıdır
	lead := models.Lead{
		ID:        bson.NewObjectID(),
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Category:  category,
		IPAddress: ip,
		CreatedAt: time.Now(),
	}
	_, _ = leadsCol.InsertOne(ctx, lead)

	filename := "ilgaz-katalog.pdf"
	if category != "Genel" {
		filename = "ilgaz-katalog-" + slugify(category) + ".pdf"
	}
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetLeads — Tüm potansiyel müşteri/indirme kayıtlarını döndürür (Admin).
func GetLeads(c *gin.Context) {
	collection := database.GetCollection("leads")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıtlar çekilemedi"})
		return
	}
	defer cursor.Close(ctx)

	leads := []models.Lead{}
	for cursor.Next(ctx) {
		var l models.Lead
		if cursor.Decode(&l) == nil {
			leads = append(leads, l)
		}
	}
	c.JSON(http.StatusOK, leads)
}

// DeleteLead — Bir potansiyel müşteri kaydını siler (Admin).
func DeleteLead(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID formatı"})
		return
	}

	collection := database.GetCollection("leads")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi: " + err.Error()})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kayıt bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kayıt silindi"})
}
