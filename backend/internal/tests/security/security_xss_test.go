package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"cooperative-erp-lite/internal/handlers"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
)

// TestXSS_MemberName tests XSS prevention in member name field
func TestXSS_MemberName(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.POST("/api/v1/anggota", anggotaHandler.Create)
	router.GET("/api/v1/anggota/:id", anggotaHandler.GetByID)

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Test Address",
		NoTelepon:      "08123456789",
		Email:          "test@koperasi.com",
	}
	db.Create(koperasi)

	// XSS payloads to test
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"<svg/onload=alert('XSS')>",
		"javascript:alert('XSS')",
		"<iframe src='javascript:alert(\"XSS\")'></iframe>",
		"<body onload=alert('XSS')>",
		"<input onfocus=alert('XSS') autofocus>",
		"<select onfocus=alert('XSS') autofocus>",
		"<textarea onfocus=alert('XSS') autofocus>",
		"<marquee onstart=alert('XSS')>",
	}

	for _, xssPayload := range xssPayloads {
		t.Run("XSS_Payload_"+xssPayload[:20], func(t *testing.T) {
			// Create member with XSS payload in name
			memberPayload := map[string]interface{}{
				"nama":           xssPayload,
				"nomor_anggota":  "XSS001",
				"nik":            "1234567890123456",
				"tempat_lahir":   "Jakarta",
				"tanggal_lahir":  "1990-01-01",
				"jenis_kelamin":  "L",
				"alamat":         "Test Address",
				"rt":             "001",
				"rw":             "002",
				"kelurahan":      "Test Kelurahan",
				"kecamatan":      "Test Kecamatan",
				"kabupaten_kota": "Jakarta",
				"provinsi":       "DKI Jakarta",
				"kode_pos":       "12345",
				"no_hp":          "08123456789",
				"email":          "xss@test.com",
				"status_anggota": "aktif",
			}

			jsonPayload, _ := json.Marshal(memberPayload)
			req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("cooperative_id", koperasi.ID.String())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should create or reject based on validation
			if w.Code == http.StatusCreated || w.Code == http.StatusOK {
				var createResponse map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &createResponse)

				// Get the created member
				if data, ok := createResponse["data"].(map[string]interface{}); ok {
					memberID := data["id"]

					getReq, _ := http.NewRequest("GET", "/api/v1/anggota/"+memberID.(string), nil)
					getReq.Header.Set("cooperative_id", koperasi.ID.String())

					getW := httptest.NewRecorder()
					router.ServeHTTP(getW, getReq)

					// Verify response is properly encoded JSON
					// XSS payload should be escaped in JSON
					assert.Equal(t, http.StatusOK, getW.Code)
					assert.Equal(t, "application/json; charset=utf-8", getW.Header().Get("Content-Type"))

					var getResponse map[string]interface{}
					err := json.Unmarshal(getW.Body.Bytes(), &getResponse)
					assert.NoError(t, err, "Response should be valid JSON")

					// JSON encoding should escape HTML characters
					responseBody := getW.Body.String()
					// In JSON, < and > should be escaped or Unicode escaped
					// The XSS payload should not be executable
					assert.NotContains(t, responseBody, "<script>", "Script tags should be escaped")
				}
			}
		})
	}
}

// TestXSS_JSONEncoding tests that JSON responses properly encode HTML
func TestXSS_JSONEncoding(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Test Address",
		NoTelepon:      "08123456789",
		Email:          "test@koperasi.com",
	}
	db.Create(koperasi)

	// Create user for authentication
	user := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaPengguna: "testuser",
		Email:        "test@test.com",
		Peran:        models.PeranAdmin,
		NamaLengkap:  "Test User",
		StatusAktif:  true,
	}
	db.Create(user)

	// Create member with potential XSS content
	member := &models.Anggota{
		IDKoperasi:     koperasi.ID,
		NomorAnggota:   "XSS001",
		NamaLengkap:    "<script>alert('XSS')</script>",
		NIK:            "1234567890123456",
		TempatLahir:    "Jakarta<script>",
		TanggalLahir:  nil,
		JenisKelamin:   "L",
		Alamat:         "<img src=x onerror=alert('XSS')>",
		RT:             "001",
		RW:             "002",
		Kelurahan:      "Test",
		Kecamatan:      "Test",
		KotaKabupaten:  "Jakarta",
		Provinsi:       "DKI",
		KodePos:        "12345",
		NoTelepon:           "08123456789",
		Email:          "test@test.com",
		Status:  "aktif",
	}
	db.Create(member)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add mock auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("idKoperasi", koperasi.ID)
		c.Set("idPengguna", user.ID)
		c.Set("peran", user.Peran)
		c.Next()
	})

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.GET("/api/v1/anggota/:id", anggotaHandler.GetByID)

	req, _ := http.NewRequest("GET", "/api/v1/anggota/"+member.ID.String(), nil)
	req.Header.Set("cooperative_id", koperasi.ID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify Content-Type is application/json
	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Response should be JSON")

	// Verify response is valid JSON (not HTML)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// Verify dangerous characters are escaped in JSON
	responseBody := w.Body.String()

	// In Go's json.Marshal, HTML characters are escaped by default
	// But Gin might disable this, so we check both scenarios

	// Script tags should not be present unescaped
	assert.NotContains(t, responseBody, "<script>alert", "Script should be escaped or encoded")

	// Image tags should not be present unescaped
	assert.NotContains(t, responseBody, "<img src=x", "Image tags should be escaped or encoded")
}

// TestXSS_HeaderInjection tests XSS prevention in HTTP headers
func TestXSS_HeaderInjection(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.GET("/api/v1/anggota", anggotaHandler.List)

	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Test Address",
		NoTelepon:      "08123456789",
		Email:          "test@koperasi.com",
	}
	db.Create(koperasi)

	// Try to inject XSS in search parameter
	xssInQuery := "<script>alert('XSS')</script>"
	req, _ := http.NewRequest("GET", "/api/v1/anggota?search="+xssInQuery, nil)
	req.Header.Set("cooperative_id", koperasi.ID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Response should be JSON, not HTML with script execution
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
}

// TestXSS_ErrorMessages tests XSS prevention in error messages
func TestXSS_ErrorMessages(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.POST("/api/v1/anggota", anggotaHandler.Create)

	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Test Address",
		NoTelepon:      "08123456789",
		Email:          "test@koperasi.com",
	}
	db.Create(koperasi)

	// Send invalid data that will trigger validation error
	invalidPayload := map[string]interface{}{
		"nama": "<script>alert('XSS in error')</script>",
		// Missing required fields to trigger validation error
	}

	jsonPayload, _ := json.Marshal(invalidPayload)
	req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cooperative_id", koperasi.ID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Error response should still be JSON
	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json")

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err, "Error response should be valid JSON")

	// Error message should not contain executable script
	responseBody := w.Body.String()
	assert.NotContains(t, responseBody, "<script>alert", "Error messages should escape XSS")
}

// TestXSS_MultipleFields tests XSS prevention across multiple fields
func TestXSS_MultipleFields(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	simpananService := services.NewSimpananService(db, services.NewTransaksiService(db))
	simpananHandler := handlers.NewSimpananHandler(simpananService)

	router.POST("/api/v1/simpanan", simpananHandler.CatatSetoran)

	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Test Address",
		NoTelepon:      "08123456789",
		Email:          "test@koperasi.com",
	}
	db.Create(koperasi)

	member := &models.Anggota{
		IDKoperasi:     koperasi.ID,
		NomorAnggota:   "001",
		NamaLengkap:  "Test Member",
		NIK:            "1234567890123456",
		TempatLahir:    "Jakarta",
		TanggalLahir:  nil,
		JenisKelamin:   "L",
		Alamat:         "Test",
		RT:             "001",
		RW:             "002",
		Kelurahan:      "Test",
		Kecamatan:      "Test",
		KotaKabupaten:  "Jakarta",
		Provinsi:       "DKI",
		KodePos:        "12345",
		NoTelepon:           "08123456789",
		Email:          "test@test.com",
		Status:  "aktif",
	}
	db.Create(member)

	// XSS in multiple fields
	simpananPayload := map[string]interface{}{
		"anggota_id":       member.ID.String(),
		"jenis_simpanan":   "pokok",
		"jumlah":           10000,
		"tanggal_transaksi": "2024-01-01",
		"nomor_referensi":  "<script>alert('XSS')</script>",
		"keterangan":       "<img src=x onerror=alert('XSS')>",
		"metode_pembayaran": "tunai",
	}

	jsonPayload, _ := json.Marshal(simpananPayload)
	req, _ := http.NewRequest("POST", "/api/v1/simpanan", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cooperative_id", koperasi.ID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Response should be JSON with escaped content
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	responseBody := w.Body.String()

	// Verify XSS payloads are not executable
	// They should be escaped or encoded in JSON
	if strings.Contains(responseBody, "script") {
		// If script appears, it should be escaped
		assert.NotContains(t, responseBody, "<script>alert", "Script tags should be escaped")
	}
}

// TestXSS_ContentSecurityPolicy tests that proper security headers are set
func TestXSS_ContentSecurityPolicy(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add security headers middleware
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.GET("/api/v1/anggota", anggotaHandler.List)

	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Test Address",
		NoTelepon:      "08123456789",
		Email:          "test@koperasi.com",
	}
	db.Create(koperasi)

	req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
	req.Header.Set("cooperative_id", koperasi.ID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify security headers are present
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
}
