package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cemonat00/ilgaz-backend/internal/database"
	"github.com/cemonat00/ilgaz-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// maskName "Ahmet Yılmaz" → "A*** Y***" dönüşümü yapar.
// Her kelimenin ilk harfi korunur, geri kalanı *** ile maskelenir.
func maskName(fullName string) string {
	words := strings.Fields(strings.TrimSpace(fullName))
	if len(words) == 0 {
		return "***"
	}
	masked := make([]string, len(words))
	for i, w := range words {
		if utf8.RuneCountInString(w) <= 1 {
			masked[i] = w + "***"
		} else {
			// İlk runeyi al, geri kalanı maskele
			runes := []rune(w)
			masked[i] = string(runes[0]) + "***"
		}
	}
	return strings.Join(masked, " ")
}

// toAdminReview — Review modelini admin response struct'ına dönüştürür.
func toAdminReview(r models.Review) models.AdminReview {
	return models.AdminReview{
		ID:          r.ID.Hex(),
		ProductID:   r.ProductID.Hex(),
		ProductName: r.ProductName,
		Name:        r.Name,
		Phone:       r.Phone,
		Email:       r.Email,
		Rating:      r.Rating,
		Content:     r.Content,
		IsApproved:  r.IsApproved,
		AdminReply:  r.AdminReply,
		HelpfulYes:  r.HelpfulYes,
		HelpfulNo:   r.HelpfulNo,
		CreatedAt:   r.CreatedAt,
	}
}

// toPublicReview — Review modelini maskeli public response struct'ına dönüştürür.
func toPublicReview(r models.Review, clientIP string) models.PublicReview {
	hasVoted := false
	for _, ip := range r.VoterIPs {
		if ip == clientIP {
			hasVoted = true
			break
		}
	}
	return models.PublicReview{
		ID:         r.ID.Hex(),
		MaskedName: maskName(r.Name),
		Rating:     r.Rating,
		Content:    r.Content,
		AdminReply: r.AdminReply,
		HelpfulYes: r.HelpfulYes,
		HelpfulNo:  r.HelpfulNo,
		CreatedAt:  r.CreatedAt,
		HasVoted:   hasVoted,
	}
}

// --- PUBLIC HANDLERS ---

// CreateReview — POST /api/yorumlar
// Yorum oluşturur; IP rate limit + lead kaydı uygular.
func CreateReview(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id"`
		Name      string `json:"ad_soyad"`
		Phone     string `json:"telefon"`
		Email     string `json:"email"`
		Rating    int    `json:"rating"`
		Content   string `json:"content"`
		KVKKOk    bool   `json:"kvkk_onay"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek verisi"})
		return
	}

	// --- Alan Doğrulama ---
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.Content = strings.TrimSpace(req.Content)

	if req.Name == "" || req.Phone == "" || req.Email == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tüm alanları doldurunuz"})
		return
	}
	if !strings.Contains(req.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçerli bir e-posta adresi giriniz"})
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Puan 1 ile 5 arasında olmalıdır"})
		return
	}
	if len([]rune(req.Content)) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Değerlendirme metni en az 10 karakter olmalıdır"})
		return
	}
	if !req.KVKKOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KVKK Aydınlatma Metnini onaylamanız zorunludur"})
		return
	}

	// Ürün ID doğrulama
	productObjID, err := bson.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ürün ID"})
		return
	}

	ip := c.ClientIP()
	reviewsCol := database.GetCollection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// --- IP Rate Limiting: 20 saniye içinde aynı IP'den yorum ---
	if recent, _ := reviewsCol.CountDocuments(ctx, bson.M{
		"ip_address": ip,
		"created_at": bson.M{"$gt": time.Now().Add(-20 * time.Second)},
	}); recent > 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Çok sık deneme yaptınız, lütfen birkaç saniye bekleyin."})
		return
	}

	// --- IP Rate Limiting: saatte 20 yorum ---
	if hourly, _ := reviewsCol.CountDocuments(ctx, bson.M{
		"ip_address": ip,
		"created_at": bson.M{"$gt": time.Now().Add(-1 * time.Hour)},
	}); hourly >= 20 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Saatlik yorum limitine ulaşıldı, lütfen daha sonra tekrar deneyin."})
		return
	}

	// Ürün adını getir (denormalize için)
	productsCol := database.GetCollection("products")
	var product models.Product
	productName := ""
	if err := productsCol.FindOne(ctx, bson.M{"_id": productObjID}).Decode(&product); err == nil {
		productName = product.Isim
		if productName == "" {
			productName = product.Name
		}
	}

	now := time.Now()
	review := models.Review{
		ID:          bson.NewObjectID(),
		ProductID:   productObjID,
		ProductName: productName,
		Name:        req.Name,
		Phone:       req.Phone,
		Email:       req.Email,
		Rating:      req.Rating,
		Content:     req.Content,
		IsApproved:  false, // Varsayılan: beklemede
		AdminReply:  "",
		HelpfulYes:  0,
		HelpfulNo:   0,
		VoterIPs:    []string{},
		IPAddress:   ip,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if _, err := reviewsCol.InsertOne(ctx, review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Değerlendirme kaydedilemedi"})
		return
	}

	// --- Lead Kaydı (aynı email varsa güncelle, yoksa ekle) ---
	leadsCol := database.GetCollection("leads")
	lead := models.Lead{
		ID:        bson.NewObjectID(),
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Category:  productName,
		Source:    "review_form",
		IPAddress: ip,
		CreatedAt: now,
	}
	// Aynı email'den zaten lead varsa duplicate oluşturma
	existingCount, _ := leadsCol.CountDocuments(ctx, bson.M{"email": req.Email})
	if existingCount == 0 {
		_, _ = leadsCol.InsertOne(ctx, lead)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Değerlendirmeniz alındı. Onaylandıktan sonra yayına girecektir."})
}

// GetProductReviews — GET /api/yorumlar/:productId
// Ürüne ait onaylı yorumları (faydalı oy sırası ile) döndürür.
func GetProductReviews(c *gin.Context) {
	productIDStr := c.Param("productId")
	productObjID, err := bson.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ürün ID"})
		return
	}

	collection := database.GetCollection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Sadece onaylı yorumlar; önce helpful_yes azalan, sonra tarih azalan sıra
	filter := bson.M{"product_id": productObjID, "is_approved": true}
	opts := options.Find().SetSort(bson.D{
		{Key: "helpful_yes", Value: -1},
		{Key: "helpful_no", Value: 1},
		{Key: "created_at", Value: -1},
	})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorumlar çekilemedi"})
		return
	}
	defer cursor.Close(ctx)

	clientIP := c.ClientIP()
	var reviews []models.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorumlar işlenemedi"})
		return
	}

	// Ortalama hesapla
	totalRating := 0
	for _, r := range reviews {
		totalRating += r.Rating
	}
	avgRating := 0.0
	if len(reviews) > 0 {
		avgRating = float64(totalRating) / float64(len(reviews))
		// Bir ondalık basamağa yuvarla
		avgRating = float64(int(avgRating*10+0.5)) / 10
	}

	// Public forma dönüştür
	publicReviews := make([]models.PublicReview, len(reviews))
	for i, r := range reviews {
		publicReviews[i] = toPublicReview(r, clientIP)
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews":    publicReviews,
		"total":      len(publicReviews),
		"avg_rating": avgRating,
	})
}

// VoteHelpful — POST /api/yorumlar/:id/oy
// Faydalı / faydasız oy işlemi. IP tekrar oyunu engeller.
func VoteHelpful(c *gin.Context) {
	reviewID := c.Param("id")
	objID, err := bson.ObjectIDFromHex(reviewID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz yorum ID"})
		return
	}

	var req struct {
		Vote string `json:"oy"` // "evet" veya "hayir"
	}
	if err := c.ShouldBindJSON(&req); err != nil || (req.Vote != "evet" && req.Vote != "hayir") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz oy değeri (evet/hayir)"})
		return
	}

	ip := c.ClientIP()
	collection := database.GetCollection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// IP daha önce oy verdi mi?
	existing, err := collection.FindOne(ctx, bson.M{
		"_id":      objID,
		"voter_ips": ip,
	}).Raw()
	if err == nil && existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Bu yoruma zaten oy verdiniz"})
		return
	}

	// Sadece onaylı yorumlara oy verilebilir
	var review models.Review
	if err := collection.FindOne(ctx, bson.M{"_id": objID, "is_approved": true}).Decode(&review); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Yorum bulunamadı"})
		return
	}

	// Oyu kaydet
	updateField := "helpful_yes"
	if req.Vote == "hayir" {
		updateField = "helpful_no"
	}

	update := bson.M{
		"$inc":  bson.M{updateField: 1},
		"$push": bson.M{"voter_ips": ip},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Oy kaydedilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Oyunuz kaydedildi"})
}

// --- ADMIN HANDLERS ---

// GetAllReviews — GET /api/admin/yorumlar (JWT korumalı)
// Tüm yorumları döndürür (hassas alanlarla birlikte, durum filtresi opsiyonel).
func GetAllReviews(c *gin.Context) {
	collection := database.GetCollection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	// Opsiyonel durum filtresi: ?durum=beklemede veya ?durum=onaylanmis
	if status := c.Query("durum"); status == "beklemede" {
		filter["is_approved"] = false
	} else if status == "onaylanmis" {
		filter["is_approved"] = true
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorumlar çekilemedi"})
		return
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorumlar işlenemedi"})
		return
	}

	result := make([]models.AdminReview, len(reviews))
	for i, r := range reviews {
		result[i] = toAdminReview(r)
	}

	c.JSON(http.StatusOK, result)
}

// ApproveReview — PUT /api/admin/yorumlar/:id/onayla (JWT korumalı)
// Yorumu onaylar ve sitede yayına alır.
func ApproveReview(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz yorum ID"})
		return
	}

	collection := database.GetCollection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"is_approved": true,
			"updated_at":  time.Now(),
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorum onaylanamadı"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Yorum bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Yorum onaylandı ve yayına alındı"})
}

// ReplyReview — PUT /api/admin/yorumlar/:id/yanit (JWT korumalı)
// Admin yanıtını kaydeder.
func ReplyReview(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz yorum ID"})
		return
	}

	var req struct {
		Reply string `json:"yanit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}
	req.Reply = strings.TrimSpace(req.Reply)
	if req.Reply == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yanıt metni boş olamaz"})
		return
	}

	collection := database.GetCollection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"admin_reply": req.Reply,
			"updated_at":  time.Now(),
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yanıt kaydedilemedi"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Yorum bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Yanıt başarıyla kaydedildi"})
}

// DeleteReview — DELETE /api/admin/yorumlar/:id (JWT korumalı)
// Yorumu kalıcı olarak siler.
func DeleteReview(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz yorum ID"})
		return
	}

	collection := database.GetCollection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Yorum silinemedi"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silinecek yorum bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Yorum silindi"})
}
