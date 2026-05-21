package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cemonat00/ilgaz-backend/internal/database"
	"github.com/cemonat00/ilgaz-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"strings"
)

func TestQAWebPortalRequirements(t *testing.T) {
	token := getAuthToken(t)

	// Clean up database collections before starting QA tests
	ctx := context.Background()
	_ = database.GetCollection("products").Drop(ctx)
	_ = database.GetCollection("messages").Drop(ctx)
	_ = database.GetCollection("reviews").Drop(ctx)
	_ = database.GetCollection("leads").Drop(ctx)
	_ = database.GetCollection("settings").Drop(ctx)

	// Seed fresh test settings and products
	database.SeedSettings()
	database.SeedProducts()
	database.SeedAdmin()

	// -------------------------------------------------------------
	// REQUIREMENT 1: Ürün Kataloğu Görüntüleme ve Detay İnceleme
	// -------------------------------------------------------------
	t.Run("Req1_ProductCatalogDisplayAndSorting", func(t *testing.T) {
		// Verify public GET /api/urunler returns active products only and sorted by order
		req, _ := http.NewRequest("GET", "/api/urunler", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var products []models.Product
		err := json.Unmarshal(w.Body.Bytes(), &products)
		assert.NoError(t, err)

		// Assert status is active and order is sorted ascending
		assert.GreaterOrEqual(t, len(products), 1)
		for i, p := range products {
			assert.Equal(t, "Aktif", p.Status)
			if i > 0 {
				assert.True(t, products[i].Order >= products[i-1].Order, "Ürünler order sırasına göre sıralanmalı")
			}
		}

		// Verify single product detail GET /api/urunler/:id
		pID := products[0].ID.Hex()
		reqDetail, _ := http.NewRequest("GET", "/api/urunler/"+pID, nil)
		wDetail := httptest.NewRecorder()
		router.ServeHTTP(wDetail, reqDetail)
		assert.Equal(t, http.StatusOK, wDetail.Code)

		var singleProduct models.Product
		json.Unmarshal(wDetail.Body.Bytes(), &singleProduct)
		assert.NotEmpty(t, singleProduct.TechnicalSpecs)
		assert.NotEmpty(t, singleProduct.FeatureBoxes)
		assert.NotEmpty(t, singleProduct.ApplicationAreas)
	})

	// -------------------------------------------------------------
	// REQUIREMENT 2: İletişim ve Talep Gönderme
	// -------------------------------------------------------------
	t.Run("Req2_ContactSubmissionAndMap", func(t *testing.T) {
		msgData := map[string]string{
			"ad_soyad": "Deneme Kullanici",
			"email":    "deneme@ilgaz.com",
			"konu":     "Teknik Destek",
			"mesaj":    "Harita entegrasyonu ve OpenStreetMap test edilmektedir.",
		}
		body, _ := json.Marshal(msgData)
		req, _ := http.NewRequest("POST", "/api/mesajlar", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify it was saved in database as Unread
		messagesCol := database.GetCollection("messages")
		var savedMsg models.Message
		err := messagesCol.FindOne(ctx, bson.M{"email": "deneme@ilgaz.com"}).Decode(&savedMsg)
		assert.NoError(t, err)
		assert.Equal(t, "Unread", savedMsg.Status)
	})

	// -------------------------------------------------------------
	// REQUIREMENT 3: Üyeliksiz Yorum ve Değerlendirme Sistemi
	// -------------------------------------------------------------
	t.Run("Req3_AnonymousReviewsAndKVKK", func(t *testing.T) {
		// Get a test product ID
		productsCol := database.GetCollection("products")
		var prod models.Product
		productsCol.FindOne(ctx, bson.M{"status": "Aktif"}).Decode(&prod)
		prodID := prod.ID.Hex()

		// Case A: Missing KVKK Checkbox
		invalidReview := map[string]interface{}{
			"product_id": prodID,
			"ad_soyad":   "Veli Han",
			"telefon":    "05551234567",
			"email":      "veli@han.com",
			"rating":     5,
			"content":    "Harika bir ürün kalitesi, çok beğendik.",
			"kvkk_onay":  false,
		}
		bodyInvalid, _ := json.Marshal(invalidReview)
		reqInv, _ := http.NewRequest("POST", "/api/yorumlar", bytes.NewBuffer(bodyInvalid))
		wInv := httptest.NewRecorder()
		router.ServeHTTP(wInv, reqInv)
		assert.Equal(t, http.StatusBadRequest, wInv.Code)

		// Case B: Valid Review (KVKK Approved) -> Should be saved as is_approved: false
		validReview := map[string]interface{}{
			"product_id": prodID,
			"ad_soyad":   "Veli Han",
			"telefon":    "05551234567",
			"email":      "veli@han.com",
			"rating":     5,
			"content":    "Harika bir ürün kalitesi, çok beğendik.",
			"kvkk_onay":  true,
		}
		bodyValid, _ := json.Marshal(validReview)
		reqVal, _ := http.NewRequest("POST", "/api/yorumlar", bytes.NewBuffer(bodyValid))
		reqVal.Header.Set("X-Forwarded-For", "192.168.1.100") // Simulate IP
		reqVal.RemoteAddr = "192.168.1.100:1234"
		wVal := httptest.NewRecorder()
		router.ServeHTTP(wVal, reqVal)
		assert.Equal(t, http.StatusCreated, wVal.Code)

		// Verify database state: is_approved: false
		reviewsCol := database.GetCollection("reviews")
		var savedReview models.Review
		err := reviewsCol.FindOne(ctx, bson.M{"email": "veli@han.com"}).Decode(&savedReview)
		assert.NoError(t, err)
		assert.False(t, savedReview.IsApproved)

		// Public product reviews endpoint GET /api/yorumlar/:productId should not list it
		reqPublic, _ := http.NewRequest("GET", "/api/yorumlar/"+prodID, nil)
		wPublic := httptest.NewRecorder()
		router.ServeHTTP(wPublic, reqPublic)
		assert.Equal(t, http.StatusOK, wPublic.Code)
		var publicResp map[string]interface{}
		json.Unmarshal(wPublic.Body.Bytes(), &publicResp)
		assert.Equal(t, float64(0), publicResp["total"])

		// Approve it via Admin API
		reviewID := savedReview.ID.Hex()
		reqApprove, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/yorumlar/%s/onayla", reviewID), nil)
		reqApprove.Header.Set("Authorization", "Bearer "+token)
		wApprove := httptest.NewRecorder()
		router.ServeHTTP(wApprove, reqApprove)
		assert.Equal(t, http.StatusOK, wApprove.Code)

		// Add Admin Reply
		replyData := map[string]string{"yanit": "Ilgaz Mühendislik olarak geri bildiriminiz için teşekkür ederiz."}
		bodyReply, _ := json.Marshal(replyData)
		reqReply, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/yorumlar/%s/yanit", reviewID), bytes.NewBuffer(bodyReply))
		reqReply.Header.Set("Authorization", "Bearer "+token)
		wReply := httptest.NewRecorder()
		router.ServeHTTP(wReply, reqReply)
		assert.Equal(t, http.StatusOK, wReply.Code)

		// Verify public view shows masked name ("V*** H***"), admin reply, and average rating
		wPublic2 := httptest.NewRecorder()
		router.ServeHTTP(wPublic2, reqPublic)
		assert.Equal(t, http.StatusOK, wPublic2.Code)
		
		var publicResp2 map[string]interface{}
		json.Unmarshal(wPublic2.Body.Bytes(), &publicResp2)
		assert.Equal(t, float64(1), publicResp2["total"])
		assert.Equal(t, float64(5), publicResp2["avg_rating"])

		reviewsList := publicResp2["reviews"].([]interface{})
		r0 := reviewsList[0].(map[string]interface{})
		assert.Equal(t, "V*** H***", r0["ad_soyad"], "Name must be masked")
		assert.Equal(t, "Ilgaz Mühendislik olarak geri bildiriminiz için teşekkür ederiz.", r0["admin_reply"])

		// Vote Helpful
		reqVote, _ := http.NewRequest("POST", fmt.Sprintf("/api/yorumlar/%s/oy", reviewID), bytes.NewBuffer([]byte(`{"oy": "evet"}`)))
		reqVote.Header.Set("X-Forwarded-For", "192.168.1.101")
		reqVote.RemoteAddr = "192.168.1.101:1234"
		wVote := httptest.NewRecorder()
		router.ServeHTTP(wVote, reqVote)
		assert.Equal(t, http.StatusOK, wVote.Code)

		// Vote again from same IP -> Should return 409
		reqVote2, _ := http.NewRequest("POST", fmt.Sprintf("/api/yorumlar/%s/oy", reviewID), bytes.NewBuffer([]byte(`{"oy": "evet"}`)))
		reqVote2.Header.Set("X-Forwarded-For", "192.168.1.101")
		reqVote2.RemoteAddr = "192.168.1.101:1234"
		wVote2 := httptest.NewRecorder()
		router.ServeHTTP(wVote2, reqVote2)
		assert.Equal(t, http.StatusConflict, wVote2.Code)
	})

	// -------------------------------------------------------------
	// REQUIREMENT 4: Ürün Kataloğu Dinamik PDF İndirme Sistemi
	// -------------------------------------------------------------
	t.Run("Req4_DynamicPDFCatalogDownload", func(t *testing.T) {
		// Case A: Missing KVKK Onay
		catalogReqFail := map[string]interface{}{
			"ad_soyad":  "Katalog Test Fail",
			"telefon":   "05559998877",
			"email":     "katalog_fail@test.com",
			"kategori":  "Vana Grubu",
			"kvkk_onay": false,
		}
		bodyFail, _ := json.Marshal(catalogReqFail)
		reqFail, _ := http.NewRequest("POST", "/api/katalog/indir", bytes.NewBuffer(bodyFail))
		wFail := httptest.NewRecorder()
		router.ServeHTTP(wFail, reqFail)
		assert.Equal(t, http.StatusBadRequest, wFail.Code)

		// Case B: Valid KVKK Onay
		catalogReq := map[string]interface{}{
			"ad_soyad":  "Katalog Test",
			"telefon":   "05559998877",
			"email":     "katalog@test.com",
			"kategori":  "Vana Grubu",
			"kvkk_onay": true,
		}
		body, _ := json.Marshal(catalogReq)
		req, _ := http.NewRequest("POST", "/api/katalog/indir", bytes.NewBuffer(body))
		req.Header.Set("X-Forwarded-For", "192.168.1.102")
		req.RemoteAddr = "192.168.1.102:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment; filename=\"ilgaz-katalog-vana-grubu.pdf\"")
		assert.NotEmpty(t, w.Body.Bytes())
	})

	// -------------------------------------------------------------
	// REQUIREMENT 5: Potansiyel Müşteri (Lead) Havuzu ve CRM Entegrasyonu
	// -------------------------------------------------------------
	t.Run("Req5_LeadCRMTrackingAndDeduplication", func(t *testing.T) {
		leadsCol := database.GetCollection("leads")

		// Verify lead created from Review Form ("review_form")
		var reviewLead models.Lead
		err := leadsCol.FindOne(ctx, bson.M{"email": "veli@han.com"}).Decode(&reviewLead)
		assert.NoError(t, err)
		assert.Equal(t, "review_form", reviewLead.Source)

		// Verify lead created from Catalog Download ("catalog_download")
		var catalogLead models.Lead
		err = leadsCol.FindOne(ctx, bson.M{"email": "katalog@test.com"}).Decode(&catalogLead)
		
		// NOTE: THIS WILL FAIL if Source is not set correctly by handlers.go
		assert.NoError(t, err, "Lead should exist for katalog@test.com")
		assert.Equal(t, "catalog_download", catalogLead.Source, "Catalog lead source must be catalog_download")

		// Test email deduplication on Catalog Download
		// Try to download again with same email
		catalogReq2 := map[string]interface{}{
			"ad_soyad":  "Katalog Test Mükerrer",
			"telefon":   "05559998877",
			"email":     "katalog@test.com",
			"kategori":  "Vana Grubu",
			"kvkk_onay": true,
		}
		body2, _ := json.Marshal(catalogReq2)
		req2, _ := http.NewRequest("POST", "/api/katalog/indir", bytes.NewBuffer(body2))
		req2.Header.Set("X-Forwarded-For", "192.168.1.103")
		req2.RemoteAddr = "192.168.1.103:1234"
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		count, _ := leadsCol.CountDocuments(ctx, bson.M{"email": "katalog@test.com"})
		// NOTE: THIS WILL FAIL if catalog download lead deduplication is missing
		assert.Equal(t, int64(1), count, "Email deduplication must prevent duplicate lead entries")
	})

	// -------------------------------------------------------------
	// REQUIREMENT 6: Yönetici (Admin) Hesap Yönetimi ve Güvenlik
	// -------------------------------------------------------------
	t.Run("Req6_SecurityAndRateLimiting", func(t *testing.T) {
		// A. Forgot Password & Reset Token
		forgotReq := models.ForgotPasswordRequest{Email: "admin@ilgazmuhendislik.com"}
		bodyForgot, _ := json.Marshal(forgotReq)
		reqForgot, _ := http.NewRequest("POST", "/api/admin/forgot-password", bytes.NewBuffer(bodyForgot))
		wForgot := httptest.NewRecorder()
		router.ServeHTTP(wForgot, reqForgot)
		assert.Equal(t, http.StatusOK, wForgot.Code)

		var forgotResp map[string]string
		json.Unmarshal(wForgot.Body.Bytes(), &forgotResp)
		resetURL := forgotResp["reset_url"]
		assert.Contains(t, resetURL, "token=")

		// Parse token
		parts := strings.Split(resetURL, "token=")
		if len(parts) < 2 {
			t.Fatalf("Reset URL does not contain token: %q", resetURL)
		}
		tokenVal := parts[1]

		// Reset Password
		resetReq := models.ResetPasswordRequest{
			Token:       tokenVal,
			NewPassword: "newsecurepass",
		}
		bodyReset, _ := json.Marshal(resetReq)
		reqReset, _ := http.NewRequest("POST", "/api/admin/reset-password", bytes.NewBuffer(bodyReset))
		wReset := httptest.NewRecorder()
		router.ServeHTTP(wReset, reqReset)
		assert.Equal(t, http.StatusOK, wReset.Code)

		// Login with new password
		loginReq := models.LoginRequest{
			Username: "admin",
			Password: "newsecurepass",
		}
		bodyLogin, _ := json.Marshal(loginReq)
		reqLogin, _ := http.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(bodyLogin))
		wLogin := httptest.NewRecorder()
		router.ServeHTTP(wLogin, reqLogin)
		assert.Equal(t, http.StatusOK, wLogin.Code)

		// Restore seed admin password for future tests
		database.SeedAdmin()

		// B. IP Rate Limiting (20s limit)
		// Try to submit review again immediately from same IP (192.168.1.100)
		productsCol := database.GetCollection("products")
		var prod models.Product
		_ = productsCol.FindOne(ctx, bson.M{}).Decode(&prod)

		dupReview := map[string]interface{}{
			"product_id": prod.ID.Hex(),
			"ad_soyad":   "Rate Limit Test",
			"telefon":    "05550000000",
			"email":      "limit@test.com",
			"rating":     4,
			"content":    "Bu yorum rate limit testi için gönderilmektedir.",
			"kvkk_onay":  true,
		}
		bodyDup, _ := json.Marshal(dupReview)
		reqDup, _ := http.NewRequest("POST", "/api/yorumlar", bytes.NewBuffer(bodyDup))
		reqDup.Header.Set("X-Forwarded-For", "192.168.1.100") // Same IP as review in Req3
		reqDup.RemoteAddr = "192.168.1.100:1234"
		wDup := httptest.NewRecorder()
		router.ServeHTTP(wDup, reqDup)
		assert.Equal(t, http.StatusTooManyRequests, wDup.Code, "Should trigger rate limit (20s)")

		// C. Vote IP Rate Limiting (20s limit)
		// Try to vote on a review again within 20s from same voting IP
		var review models.Review
		database.GetCollection("reviews").FindOne(ctx, bson.M{"is_approved": true}).Decode(&review)
		reviewID := review.ID.Hex()

		reqVoteDup, _ := http.NewRequest("POST", fmt.Sprintf("/api/yorumlar/%s/oy", reviewID), bytes.NewBuffer([]byte(`{"oy": "evet"}`)))
		reqVoteDup.Header.Set("X-Forwarded-For", "192.168.1.101") // Same voting IP as Req3
		reqVoteDup.RemoteAddr = "192.168.1.101:1234"
		wVoteDup := httptest.NewRecorder()
		router.ServeHTTP(wVoteDup, reqVoteDup)
		
		// NOTE: THIS WILL FAIL if vote rate limiting is not implemented in review_handlers.go
		assert.Equal(t, http.StatusTooManyRequests, wVoteDup.Code, "Vote operations must be rate limited (20s)")
	})

	// -------------------------------------------------------------
	// REQUIREMENT 7: Yönetici Paneli - Ürün Yönetimi
	// -------------------------------------------------------------
	t.Run("Req7_AdminProductCRUDAndReorder", func(t *testing.T) {
		// Add product
		newProd := models.Product{
			Isim:        "Yeni Admin Ürün",
			Baslik:      "Yeni Admin Ürün - Başlık",
			Kategori:    "Vana Grubu",
			Price:       5000,
			StockStatus: true,
			Status:      "Aktif",
			Order:       10,
		}
		body, _ := json.Marshal(newProd)
		req, _ := http.NewRequest("POST", "/api/admin/urunler", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdProd models.Product
		json.Unmarshal(w.Body.Bytes(), &createdProd)
		assert.NotEmpty(t, createdProd.ID)
		createdID := createdProd.ID.Hex()

		// Update product
		createdProd.Isim = "Güncel Admin Ürün"
		bodyUp, _ := json.Marshal(createdProd)
		reqUp, _ := http.NewRequest("PUT", "/api/admin/urunler/"+createdID, bytes.NewBuffer(bodyUp))
		reqUp.Header.Set("Authorization", "Bearer "+token)
		wUp := httptest.NewRecorder()
		router.ServeHTTP(wUp, reqUp)
		assert.Equal(t, http.StatusOK, wUp.Code)

		// Reorder products
		reorderReq := map[string][]string{
			"order": {createdID},
		}
		bodyReorder, _ := json.Marshal(reorderReq)
		reqReorder, _ := http.NewRequest("PUT", "/api/admin/urunler/reorder", bytes.NewBuffer(bodyReorder))
		reqReorder.Header.Set("Authorization", "Bearer "+token)
		wReorder := httptest.NewRecorder()
		router.ServeHTTP(wReorder, reqReorder)
		assert.Equal(t, http.StatusOK, wReorder.Code)

		// Delete product
		reqDel, _ := http.NewRequest("DELETE", "/api/admin/urunler/"+createdID, nil)
		reqDel.Header.Set("Authorization", "Bearer "+token)
		wDel := httptest.NewRecorder()
		router.ServeHTTP(wDel, reqDel)
		assert.Equal(t, http.StatusOK, wDel.Code)
	})

	// -------------------------------------------------------------
	// REQUIREMENT 8: Yönetici Paneli - İletişim Ayarları & Harita Konumlandırma
	// -------------------------------------------------------------
	t.Run("Req8_AdminSettingsAndMapConfig", func(t *testing.T) {
		// Update settings
		newSettings := models.SiteSettings{
			CompanyName:  "Ilgaz Mühendislik QA",
			Email:        "qa@ilgaz.com",
			Phone:        "+90 212 999 88 77",
			SupportEmail: "qasupport@ilgaz.com",
			Address:      "QA Caddesi No:5, Samsun",
			MapLocation:  "41.331,36.280",
			MapZoom:      15,
		}
		body, _ := json.Marshal(newSettings)
		req, _ := http.NewRequest("PUT", "/api/admin/ayarlar", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify settings on public get
		reqGet, _ := http.NewRequest("GET", "/api/ayarlar", nil)
		wGet := httptest.NewRecorder()
		router.ServeHTTP(wGet, reqGet)
		assert.Equal(t, http.StatusOK, wGet.Code)

		var savedSettings models.SiteSettings
		json.Unmarshal(wGet.Body.Bytes(), &savedSettings)
		assert.Equal(t, "Ilgaz Mühendislik QA", savedSettings.CompanyName)
		assert.Equal(t, "41.331,36.280", savedSettings.MapLocation)
		assert.Equal(t, 15, savedSettings.MapZoom)

		// Restore default settings
		database.GetCollection("settings").Drop(ctx)
		database.SeedSettings()
	})

	// -------------------------------------------------------------
	// REQUIREMENT 9: Gelen Kutusu ve Mesaj Yönetimi
	// -------------------------------------------------------------
	t.Run("Req9_InboxAndSMTPReply", func(t *testing.T) {
		// Get messages list
		req, _ := http.NewRequest("GET", "/api/admin/mesajlar", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var messages []models.Message
		json.Unmarshal(w.Body.Bytes(), &messages)
		assert.GreaterOrEqual(t, len(messages), 1)

		msgID := messages[0].ID.Hex()

		// Reply to message (SMTP Mock)
		replyReq := map[string]string{"message": "Talep alındı, size dönüş sağlayacağız."}
		body, _ := json.Marshal(replyReq)
		reqReply, _ := http.NewRequest("POST", fmt.Sprintf("/api/admin/mesajlar/%s/reply", msgID), bytes.NewBuffer(body))
		reqReply.Header.Set("Authorization", "Bearer "+token)
		wReply := httptest.NewRecorder()
		router.ServeHTTP(wReply, reqReply)
		assert.Equal(t, http.StatusOK, wReply.Code)
	})

	// -------------------------------------------------------------
	// REQUIREMENT 10: Yönetici Paneli - Yorum Yönetimi
	// -------------------------------------------------------------
	t.Run("Req10_AdminReviewManagement", func(t *testing.T) {
		// Get reviews list
		req, _ := http.NewRequest("GET", "/api/admin/yorumlar", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var adminReviews []models.AdminReview
		json.Unmarshal(w.Body.Bytes(), &adminReviews)
		assert.GreaterOrEqual(t, len(adminReviews), 1)
	})

	// -------------------------------------------------------------
	// REQUIREMENT 11: Dashboard ve İstatistik Raporlama
	// -------------------------------------------------------------
	t.Run("Req11_DashboardTrends", func(t *testing.T) {
		// Dashboard collects stats from multiple endpoints: /api/urunler, /api/admin/yorumlar?durum=beklemede, /api/admin/mesajlar, /api/admin/leads
		reqLeads, _ := http.NewRequest("GET", "/api/admin/leads", nil)
		reqLeads.Header.Set("Authorization", "Bearer "+token)
		wLeads := httptest.NewRecorder()
		router.ServeHTTP(wLeads, reqLeads)
		assert.Equal(t, http.StatusOK, wLeads.Code)

		var leads []models.Lead
		json.Unmarshal(wLeads.Body.Bytes(), &leads)
		assert.GreaterOrEqual(t, len(leads), 1)
	})
}
