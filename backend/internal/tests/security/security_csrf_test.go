package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"cooperative-erp-lite/internal/handlers"
	"cooperative-erp-lite/internal/middleware"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
)

// TestCSRF_TokenGeneration tests CSRF token generation
func TestCSRF_TokenGeneration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/csrf-token", middleware.GenerateCSRFTokenEndpoint)

	req, _ := http.NewRequest("GET", "/api/v1/csrf-token", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Verify token is returned
	token, exists := response["token"]
	assert.True(t, exists, "Token should be present in response")
	assert.NotEmpty(t, token, "Token should not be empty")

	// Verify token is set as cookie
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == middleware.CSRFTokenCookie {
			found = true
			assert.NotEmpty(t, cookie.Value, "Cookie value should not be empty")
			assert.True(t, cookie.HttpOnly, "Cookie should be HttpOnly")
		}
	}
	assert.True(t, found, "CSRF cookie should be set")
}

// TestCSRF_MissingToken tests that requests without CSRF token are rejected
func TestCSRF_MissingToken(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add CSRF protection middleware
	router.Use(middleware.CSRFProtection())

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.POST("/api/v1/anggota", anggotaHandler.Create)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	// Try to create member without CSRF token
	payload := map[string]interface{}{
		"namaLengkap":      "Test Member",
		"nik":              "1234567890123456",
		"tempatLahir":      "Jakarta",
		"tanggalLahir":     "1990-01-01T00:00:00Z",
		"jenisKelamin":     "L",
		"alamat":           "Test Address",
		"rt":               "001",
		"rw":               "002",
		"kelurahan":        "Test",
		"kecamatan":        "Test",
		"kotaKabupaten":    "Jakarta",
		"provinsi":         "DKI",
		"kodePos":          "12345",
		"noTelepon":        "08123456789",
		"email":            "test@member.com",
		"tanggalBergabung": "2025-01-01T00:00:00Z",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cooperative_id", koperasi.ID.String())
	// No CSRF token header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should be forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response["error"], "CSRF token missing")
}

// TestCSRF_ValidToken tests that requests with valid CSRF token succeed
func TestCSRF_ValidToken(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Generate a valid token first
	token, err := middleware.GenerateCSRFToken()
	assert.NoError(t, err)

	// Add CSRF protection middleware
	router.Use(middleware.CSRFProtection())

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.POST("/api/v1/anggota", anggotaHandler.Create)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	payload := map[string]interface{}{
		"namaLengkap":      "Test Member",
		"nik":              "1234567890123456",
		"tempatLahir":      "Jakarta",
		"tanggalLahir":     "1990-01-01T00:00:00Z",
		"jenisKelamin":     "L",
		"alamat":           "Test Address",
		"rt":               "001",
		"rw":               "002",
		"kelurahan":        "Test",
		"kecamatan":        "Test",
		"kotaKabupaten":    "Jakarta",
		"provinsi":         "DKI",
		"kodePos":          "12345",
		"noTelepon":        "08123456789",
		"email":            "test@member.com",
		"tanggalBergabung": "2025-01-01T00:00:00Z",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cooperative_id", koperasi.ID.String())
	req.Header.Set(middleware.CSRFTokenHeader, token) // Add valid CSRF token

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed (or return validation error, not CSRF error)
	assert.NotEqual(t, http.StatusForbidden, w.Code, "Should not be forbidden with valid CSRF token")
}

// TestCSRF_InvalidToken tests that requests with invalid CSRF token are rejected
func TestCSRF_InvalidToken(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(middleware.CSRFProtection())

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.POST("/api/v1/anggota", anggotaHandler.Create)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	payload := map[string]interface{}{
		"namaLengkap": "Test Member",
		"email":       "test@member.com",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cooperative_id", koperasi.ID.String())
	req.Header.Set(middleware.CSRFTokenHeader, "invalid-token-12345") // Invalid token

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should be forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response["error"], "Invalid or expired CSRF token")
}

// TestCSRF_SafeMethodsBypass tests that safe methods (GET, HEAD, OPTIONS) bypass CSRF check
func TestCSRF_SafeMethodsBypass(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(middleware.CSRFProtection())

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.GET("/api/v1/anggota", anggotaHandler.List)
	router.HEAD("/api/v1/anggota", anggotaHandler.List)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	safeMethods := []string{"GET", "HEAD", "OPTIONS"}

	for _, method := range safeMethods {
		t.Run("SafeMethod_"+method, func(t *testing.T) {
			req, _ := http.NewRequest(method, "/api/v1/anggota", nil)
			req.Header.Set("cooperative_id", koperasi.ID.String())
			// No CSRF token

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not be forbidden (safe methods don't require CSRF)
			assert.NotEqual(t, http.StatusForbidden, w.Code, method+" should not require CSRF token")
		})
	}
}

// TestCSRF_StateChangingMethods tests that state-changing methods require CSRF
func TestCSRF_StateChangingMethods(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(middleware.CSRFProtection())

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.POST("/api/v1/anggota", anggotaHandler.Create)
	router.PUT("/api/v1/anggota/:id", anggotaHandler.Update)
	router.DELETE("/api/v1/anggota/:id", anggotaHandler.Delete)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	stateChangingMethods := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/anggota"},
		{"PUT", "/api/v1/anggota/123e4567-e89b-12d3-a456-426614174000"},
		{"DELETE", "/api/v1/anggota/123e4567-e89b-12d3-a456-426614174000"},
	}

	for _, test := range stateChangingMethods {
		t.Run("RequireCSRF_"+test.method, func(t *testing.T) {
			var req *http.Request
			if test.method == "POST" || test.method == "PUT" {
				payload := map[string]interface{}{"namaLengkap": "Test"}
				jsonPayload, _ := json.Marshal(payload)
				req, _ = http.NewRequest(test.method, test.path, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(test.method, test.path, nil)
			}
			req.Header.Set("cooperative_id", koperasi.ID.String())
			// No CSRF token

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should be forbidden without CSRF token
			assert.Equal(t, http.StatusForbidden, w.Code, test.method+" should require CSRF token")

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["error"], "CSRF")
		})
	}
}

// TestCSRF_TokenReuse tests that the same token can be used multiple times
func TestCSRF_TokenReuse(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Generate a valid token
	token, err := middleware.GenerateCSRFToken()
	assert.NoError(t, err)

	router.Use(middleware.CSRFProtection())

	simpananService := services.NewSimpananService(db, services.NewTransaksiService(db))
	simpananHandler := handlers.NewSimpananHandler(simpananService)

	router.POST("/api/v1/simpanan", simpananHandler.CatatSetoran)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	// Use the same token for multiple requests
	for i := 0; i < 3; i++ {
		payload := map[string]interface{}{
			"tipeSimpanan":   "pokok",
			"jumlahSetoran":  100000,
			"tanggalTransaksi": "2025-01-01T00:00:00Z",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/simpanan", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("cooperative_id", koperasi.ID.String())
		req.Header.Set(middleware.CSRFTokenHeader, token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should not be forbidden (token can be reused)
		assert.NotEqual(t, http.StatusForbidden, w.Code, "Token should be reusable")
	}
}

// TestCSRF_FormFieldFallback tests CSRF token in form field as fallback
func TestCSRF_FormFieldFallback(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	token, err := middleware.GenerateCSRFToken()
	assert.NoError(t, err)

	router.Use(middleware.CSRFProtection())

	router.POST("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test with form data instead of JSON
	req, _ := http.NewRequest("POST", "/api/v1/test", nil)
	req.PostForm = map[string][]string{
		"csrf_token": {token},
		"name":       {"test"},
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed with form field token
	assert.NotEqual(t, http.StatusForbidden, w.Code, "Should accept CSRF token from form field")
}
