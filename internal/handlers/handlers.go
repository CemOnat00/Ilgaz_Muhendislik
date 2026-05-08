package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/cemonat00/ilgaz-backend/internal/database"
	"github.com/cemonat00/ilgaz-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"path/filepath"
	"fmt"
)

// --- UPLOAD HANDLER ---

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya yüklenemedi"})
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

func GetProducts(c *gin.Context) {
	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var products []models.Product
	cursor, err := collection.Find(ctx, bson.M{})
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
	objID, _ := bson.ObjectIDFromHex(id)

	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
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
	objID, _ := bson.ObjectIDFromHex(id)

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
			"name":         product.Name,
			"isim":         product.Isim,
			"category":     product.Category,
			"kategori":     product.Kategori,
			"price":        product.Price,
			"stock_status": product.StockStatus,
			"image_url":    product.ImageURL,
			"description":  product.Description,
			"status":       product.Status,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ürün güncellendi"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	objID, _ := bson.ObjectIDFromHex(id)

	collection := database.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ürün silindi"})
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
	objID, _ := bson.ObjectIDFromHex(id)

	collection := database.GetCollection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"is_read": true, "status": "Read"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Okundu işaretlendi"})
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	objID, _ := bson.ObjectIDFromHex(id)

	collection := database.GetCollection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesaj silindi"})
}

// --- PROJECT HANDLERS ---

func GetProjects(c *gin.Context) {
	collection := database.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var projects []models.Project
	cursor, err := collection.Find(ctx, bson.M{})
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
	objID, _ := bson.ObjectIDFromHex(id)

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

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme başarısız"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proje güncellendi"})
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	objID, _ := bson.ObjectIDFromHex(id)

	collection := database.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silme işlemi başarısız"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proje silindi"})
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
	
	// Simple password check (In production, use bcrypt)
	if err != nil || admin.Password != loginReq.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz kullanıcı adı veya şifre"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:   "mock-jwt-token-" + admin.Username,
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
	_, err = collection.UpdateOne(ctx, bson.M{"_id": admin.ID}, bson.M{
		"$set":   bson.M{"password": req.NewPassword},
		"$unset": bson.M{"reset_token": ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Şifre başarıyla güncellendi"})
}
