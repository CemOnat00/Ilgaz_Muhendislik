package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cemonat00/ilgaz-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func getAuthToken(t *testing.T) string {
	loginData := models.LoginRequest{
		Username: "testadmin",
		Password: "testpass123",
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp models.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Token
}

func TestProductCRUD(t *testing.T) {
	token := getAuthToken(t)
	var productID string

	// 1. Create Product
	newProduct := models.Product{
		Baslik:   "Test Product",
		Isim:     "Test Name",
		Kategori: "Vana Grubu",
		Price:    100.0,
		Status:   "Aktif",
	}
	body, _ := json.Marshal(newProduct)
	req, _ := http.NewRequest("POST", "/api/admin/urunler", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdProduct models.Product
	json.Unmarshal(w.Body.Bytes(), &createdProduct)
	productID = createdProduct.ID.Hex()
	assert.NotEmpty(t, productID)

	// 2. Get Products (Public)
	req, _ = http.NewRequest("GET", "/api/urunler", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Update Product
	updatedProduct := createdProduct
	updatedProduct.Isim = "Updated Name"
	body, _ = json.Marshal(updatedProduct)
	req, _ = http.NewRequest("PUT", "/api/admin/urunler/"+productID, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Delete Product
	req, _ = http.NewRequest("DELETE", "/api/admin/urunler/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Negative Path: Get Non-existent
	req, _ = http.NewRequest("GET", "/api/urunler/"+productID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMessageCRUD(t *testing.T) {
	// 1. Create Message (Public)
	newMsg := models.Message{
		Name:    "Tester",
		Email:   "test@test.com",
		Subject: "Test Subject",
		Content: "Test Message",
	}
	body, _ := json.Marshal(newMsg)
	req, _ := http.NewRequest("POST", "/api/mesajlar", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Get Messages (Admin)
	token := getAuthToken(t)
	req, _ = http.NewRequest("GET", "/api/admin/mesajlar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	var messages []models.Message
	json.Unmarshal(w.Body.Bytes(), &messages)
	assert.GreaterOrEqual(t, len(messages), 1)
	
	msgID := messages[0].ID.Hex()

	// 3. Mark Read
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/admin/mesajlar/%s/read", msgID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Delete Message
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/admin/mesajlar/%s", msgID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
