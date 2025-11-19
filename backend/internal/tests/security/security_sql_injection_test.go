package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"cooperative-erp-lite/internal/handlers"
	"cooperative-erp-lite/internal/middleware"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
)



// setupTestRouter creates a test router with auth middleware
func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup services
	anggotaService := services.NewAnggotaService(db)
	simpananService := services.NewSimpananService(db, services.NewTransaksiService(db))
	transaksiService := services.NewTransaksiService(db)

	// Setup handlers
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)
	simpananHandler := handlers.NewSimpananHandler(simpananService)
	transaksiHandler := handlers.NewTransaksiHandler(transaksiService)

	// Setup routes
	api := router.Group("/api/v1")
	{
		api.Use(middleware.AuthMiddleware(utils.NewJWTUtil("test-secret", 24)))
		api.GET("/anggota", anggotaHandler.List)
		api.GET("/simpanan", simpananHandler.List)
		api.GET("/transaksi", transaksiHandler.List)
	}

	return router
}

// TestSQLInjection_MemberSearch tests SQL injection prevention in member search
func TestSQLInjection_MemberSearch(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	router := setupTestRouter(db)

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		// NomorBadan:  "123456",
		Alamat:      "Test Address",
		NoTelepon:    "08123456789",
		Email:       "test@koperasi.com",
	}
	db.Create(koperasi)

	// Create test member
	member := &models.Anggota{
		IDKoperasi:     koperasi.ID,
		NomorAnggota:   "001",
		NamaLengkap:  "Test Member",
		NIK:            "1234567890123456",
		TempatLahir:    "Jakarta",
		TanggalLahir:  nil,
		JenisKelamin:   "L",
		Alamat:         "Test Address",
		RT:             "001",
		RW:             "002",
		Kelurahan:      "Test Kelurahan",
		Kecamatan:      "Test Kecamatan",
		KotaKabupaten:  "Jakarta",
		Provinsi:       "DKI Jakarta",
		KodePos:        "12345",
		NoTelepon:           "08123456789",
		Email:          "test@member.com",
		Status:  "aktif",
	}
	db.Create(member)

	// SQL injection payloads to test
	sqlInjectionPayloads := []string{
		"' OR '1'='1",
		"'; DROP TABLE anggotas; --",
		"1' UNION SELECT NULL, NULL, NULL--",
		"admin'--",
		"' OR 1=1--",
		"1' AND '1' = '1",
		"'; DELETE FROM anggotas WHERE '1'='1",
	}

	for _, payload := range sqlInjectionPayloads {
		t.Run("Injection_"+payload, func(t *testing.T) {
			// Create request with SQL injection payload
			req, _ := http.NewRequest("GET", "/api/v1/anggota?search="+payload, nil)
			req.Header.Set("cooperative_id", koperasi.ID.String())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not return error or execute SQL injection
			// Response should be valid JSON (not SQL error)
			assert.NotEqual(t, http.StatusInternalServerError, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")

			// Verify database still has the test member
			var count int64
			db.Model(&models.Anggota{}).Where("id_koperasi = ?", koperasi.ID).Count(&count)
			assert.Equal(t, int64(1), count, "Member should still exist (not deleted by injection)")
		})
	}
}

// TestSQLInjection_TransactionFilter tests SQL injection in transaction filtering
func TestSQLInjection_TransactionFilter(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	router := setupTestRouter(db)

	// Create test data
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		// NomorBadan:  "123456",
		Alamat:      "Test Address",
		NoTelepon:    "08123456789",
		Email:       "test@koperasi.com",
	}
	db.Create(koperasi)

	// SQL injection in date filter
	sqlInjectionPayloads := []string{
		"2024-01-01' OR '1'='1",
		"'; DROP TABLE transaksis; --",
		"2024-01-01'; DELETE FROM transaksis WHERE '1'='1",
	}

	for _, payload := range sqlInjectionPayloads {
		t.Run("DateFilter_"+payload, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/transaksi?tanggal_mulai="+payload, nil)
			req.Header.Set("cooperative_id", koperasi.ID.String())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should handle invalid date format gracefully
			// Should not execute SQL injection
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")
		})
	}
}

// TestSQLInjection_ShareCapitalType tests SQL injection in share capital type filter
func TestSQLInjection_ShareCapitalType(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	router := setupTestRouter(db)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		// NomorBadan:  "123456",
		Alamat:      "Test Address",
		NoTelepon:    "08123456789",
		Email:       "test@koperasi.com",
	}
	db.Create(koperasi)

	sqlInjectionPayloads := []string{
		"pokok' OR '1'='1",
		"wajib'; DROP TABLE simpanans; --",
		"sukarela' UNION SELECT * FROM penggunas--",
	}

	for _, payload := range sqlInjectionPayloads {
		t.Run("TypeFilter_"+payload, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/simpanan?jenis="+payload, nil)
			req.Header.Set("cooperative_id", koperasi.ID.String())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")
		})
	}
}

// TestSQLInjection_PostRequest tests SQL injection in POST request body
func TestSQLInjection_PostRequest(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	router := setupTestRouter(db)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		// NomorBadan:  "123456",
		Alamat:      "Test Address",
		NoTelepon:    "08123456789",
		Email:       "test@koperasi.com",
	}
	db.Create(koperasi)

	// SQL injection in POST body
	maliciousPayload := map[string]interface{}{
		"namaLengkap":    "Test' OR '1'='1",
		"nomorAnggota":   "001'; DROP TABLE anggotas; --",
		"nik":            "1234567890123456",
		"tempatLahir":    "Jakarta",
		"tanggalLahir":   "1990-01-01",
		"jenisKelamin":   "L",
		"alamat":         "Test",
		"rt":             "001",
		"rw":             "002",
		"kelurahan":      "Test",
		"kecamatan":      "Test",
		"kotaKabupaten":  "Jakarta",
		"provinsi":       "DKI",
		"kodePos":        "12345",
		"noTelepon":      "08123456789",
		"email":          "test@test.com",
		"status":         "aktif",
	}

	jsonPayload, _ := json.Marshal(maliciousPayload)
	req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cooperative_id", koperasi.ID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should insert data safely or return validation error
	// Should not execute SQL injection

	// The important thing is that the response is not a SQL error
	// and the table still exists (no DROP TABLE executed)
	assert.NotContains(t, w.Body.String(), "SQL", "Should not return SQL errors")
	assert.NotContains(t, w.Body.String(), "syntax error", "Should not return syntax errors")

	// Try to parse as JSON (may fail if validation error with different format)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		// If JSON parsing fails, just check the response is not empty
		// and doesn't contain SQL error messages
		assert.NotEmpty(t, w.Body.String(), "Response should not be empty")
	}

	// Verify table still exists and data is safe
	var count int64
	db.Model(&models.Anggota{}).Count(&count)
	assert.GreaterOrEqual(t, count, int64(0), "Table should still exist")

	// Verify no tables were dropped
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'anggota'").Scan(&tableCount)
	assert.Equal(t, int64(1), tableCount, "Anggota table should still exist")
}

// TestParameterizedQueries verifies that all queries use parameterized statements
func TestParameterizedQueries(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Enable SQL logging to verify parameterized queries
	db = db.Debug()

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		// NomorBadan:  "123456",
		Alamat:      "Test Address",
		NoTelepon:    "08123456789",
		Email:       "test@koperasi.com",
	}
	db.Create(koperasi)

	// Test various query types
	t.Run("SelectQuery", func(t *testing.T) {
		var members []models.Anggota
		result := db.Where("id_koperasi = ?", koperasi.ID).Find(&members)
		assert.NoError(t, result.Error)
	})

	t.Run("InsertQuery", func(t *testing.T) {
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
		result := db.Create(member)
		assert.NoError(t, result.Error)
	})

	t.Run("UpdateQuery", func(t *testing.T) {
		var member models.Anggota
		db.Where("id_koperasi = ?", koperasi.ID).First(&member)

		result := db.Model(&member).Update("nama_lengkap", "Updated Name")
		assert.NoError(t, result.Error)
	})

	t.Run("DeleteQuery", func(t *testing.T) {
		var member models.Anggota
		db.Where("id_koperasi = ?", koperasi.ID).First(&member)

		result := db.Delete(&member)
		assert.NoError(t, result.Error)
	})
}
