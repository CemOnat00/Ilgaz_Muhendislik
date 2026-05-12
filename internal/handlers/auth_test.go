package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cemonat00/ilgaz-backend/internal/handlers"
	"github.com/cemonat00/ilgaz-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Protected Admin Endpoints
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired())
	{
		admin.GET("/mesajlar", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Mesajlar"})
		})
		admin.POST("/upload", handlers.UploadImage)
	}

	return r
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	r := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/admin/mesajlar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for request with no token, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/admin/mesajlar", nil)
	req.Header.Set("Authorization", "Bearer invalid-fake-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for request with invalid token, got %d", w.Code)
	}
}

func TestUploadImage_InvalidExtension(t *testing.T) {

	// Create a fake .exe file in memory
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "fake_hack.exe")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("fake executable content"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/admin/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Bypass auth for this specific handler test by simulating a valid token context,
	// or in this case, we just expect a 401 if auth fails, but we want to test the handler logic.
	// Since the router applies AuthRequired, we need to mock a valid token or test the handler directly.

	// Let's test the handler directly without middleware to isolate the extension test
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.UploadImage(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for .exe file upload, got %d. Body: %s", w.Code, w.Body.String())
	}
}
