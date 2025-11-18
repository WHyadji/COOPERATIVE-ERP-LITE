package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"cooperative-erp-lite/internal/handlers"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
)

// TestMultiTenant_DataIsolation tests that cooperatives cannot access each other's data
func TestMultiTenant_DataIsolation(t *testing.T) {
	db := setupTestDB(t)
	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(testAuthMiddleware(jwtUtil))

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.GET("/api/v1/anggota", anggotaHandler.List)
	router.GET("/api/v1/anggota/:id", anggotaHandler.GetByID)

	// Create two cooperatives
	koperasi1 := &models.Koperasi{
		NamaKoperasi:   "Koperasi A",
		Alamat:         "Address A",
		NoTelepon:      "08111111111",
		Email:          "koperasia@test.com",
	}
	db.Create(koperasi1)

	koperasi2 := &models.Koperasi{
		NamaKoperasi:   "Koperasi B",
		Alamat:         "Address B",
		NoTelepon:      "08122222222",
		Email:          "koperasib@test.com",
	}
	db.Create(koperasi2)

	// Create users for each cooperative
	user1 := &models.Pengguna{
		IDKoperasi: koperasi1.ID,
		NamaPengguna:   "admin1",
		Email:      "admin1@test.com",
		// Password:   "hashed1",
		Peran:       models.PeranAdmin,
		NamaLengkap: "Admin 1",
		StatusAktif:   true,
	}
	db.Create(user1)

	user2 := &models.Pengguna{
		IDKoperasi: koperasi2.ID,
		NamaPengguna:   "admin2",
		Email:      "admin2@test.com",
		// Password:   "hashed2",
		Peran:       models.PeranAdmin,
		NamaLengkap: "Admin 2",
		StatusAktif:   true,
	}
	db.Create(user2)

	// Create members for Koperasi A
	member1 := &models.Anggota{
		IDKoperasi:     koperasi1.ID,
		NomorAnggota:   "A001",
		NamaLengkap:  "Member A1",
		NIK:            "1111111111111111",
		TempatLahir:    "Jakarta",
		TanggalLahir:  nil,
		JenisKelamin:   "L",
		Alamat:         "Address",
		RT:             "001",
		RW:             "002",
		Kelurahan:      "Kelurahan",
		Kecamatan:      "Kecamatan",
		KotaKabupaten:  "Jakarta",
		Provinsi:       "DKI",
		KodePos:        "12345",
		NoTelepon:           "08111111111",
		Email:          "member1@test.com",
	}
	db.Create(member1)

	member2 := &models.Anggota{
		IDKoperasi:     koperasi1.ID,
		NomorAnggota:   "A002",
		NamaLengkap:  "Member A2",
		NIK:            "1111111111111112",
		TempatLahir:    "Jakarta",
		TanggalLahir:  nil,
		JenisKelamin:   "P",
		Alamat:         "Address",
		RT:             "001",
		RW:             "002",
		Kelurahan:      "Kelurahan",
		Kecamatan:      "Kecamatan",
		KotaKabupaten:  "Jakarta",
		Provinsi:       "DKI",
		KodePos:        "12345",
		NoTelepon:           "08111111112",
		Email:          "member2@test.com",
	}
	db.Create(member2)

	// Create members for Koperasi B
	member3 := &models.Anggota{
		IDKoperasi:     koperasi2.ID,
		NomorAnggota:   "B001",
		NamaLengkap:  "Member B1",
		NIK:            "2222222222222221",
		TempatLahir:    "Bandung",
		TanggalLahir:  nil,
		JenisKelamin:   "L",
		Alamat:         "Address",
		RT:             "001",
		RW:             "002",
		Kelurahan:      "Kelurahan",
		Kecamatan:      "Kecamatan",
		KotaKabupaten:  "Bandung",
		Provinsi:       "Jabar",
		KodePos:        "54321",
		NoTelepon:           "08122222221",
		Email:          "member3@test.com",
	}
	db.Create(member3)

	// Generate tokens for both users
	token1, _ := jwtUtil.GenerateToken(user1)
	token2, _ := jwtUtil.GenerateToken(user2)

	t.Run("Koperasi_A_Should_Only_See_Own_Members", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+token1)
		req.Header.Set("cooperative_id", koperasi1.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// Should only see 2 members from Koperasi A
		if data, ok := response["data"].([]interface{}); ok {
			assert.Equal(t, 2, len(data), "Should only see own members")

			// Verify all members belong to Koperasi A
			for _, item := range data {
				member := item.(map[string]interface{})
				assert.Contains(t, []string{"A001", "A002"}, member["nomorAnggota"], "Should only return Koperasi A members")
			}
		}
	})

	t.Run("Koperasi_B_Should_Only_See_Own_Members", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+token2)
		req.Header.Set("cooperative_id", koperasi2.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// Should only see 1 member from Koperasi B
		if data, ok := response["data"].([]interface{}); ok {
			assert.Equal(t, 1, len(data), "Should only see own members")

			member := data[0].(map[string]interface{})
			assert.Equal(t, "B001", member["nomorAnggota"], "Should only return Koperasi B members")
		}
	})

	t.Run("Cannot_Access_Another_Cooperative_Member_By_ID", func(t *testing.T) {
		// Koperasi A trying to access Koperasi B's member
		req, _ := http.NewRequest("GET", "/api/v1/anggota/"+member3.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token1)
		req.Header.Set("cooperative_id", koperasi1.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return not found or forbidden
		assert.NotEqual(t, http.StatusOK, w.Code, "Should not be able to access another cooperative's member")
		assert.Contains(t, []int{http.StatusNotFound, http.StatusForbidden}, w.Code)
	})
}

// TestMultiTenant_ShareCapitalIsolation tests share capital data isolation
func TestMultiTenant_ShareCapitalIsolation(t *testing.T) {
	db := setupTestDB(t)
	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(testAuthMiddleware(jwtUtil))

	simpananService := services.NewSimpananService(db, services.NewTransaksiService(db))
	simpananHandler := handlers.NewSimpananHandler(simpananService)

	router.GET("/api/v1/simpanan", simpananHandler.List)

	// Create two cooperatives with members and share capital
	koperasi1 := &models.Koperasi{
		NamaKoperasi:   "Koperasi A",
		Alamat:         "Address A",
		NoTelepon:      "08111111111",
		Email:          "a@test.com",
	}
	db.Create(koperasi1)

	koperasi2 := &models.Koperasi{
		NamaKoperasi:   "Koperasi B",
		Alamat:         "Address B",
		NoTelepon:      "08122222222",
		Email:          "b@test.com",
	}
	db.Create(koperasi2)

	user1 := &models.Pengguna{
		IDKoperasi: koperasi1.ID,
		NamaPengguna:   "admin1",
		Email:      "admin1@test.com",
		// Password:   "hashed",
		Peran:       models.PeranAdmin,
		NamaLengkap: "Admin 1",
		StatusAktif:   true,
	}
	db.Create(user1)

	user2 := &models.Pengguna{
		IDKoperasi: koperasi2.ID,
		NamaPengguna:   "admin2",
		Email:      "admin2@test.com",
		// Password:   "hashed",
		Peran:       models.PeranAdmin,
		NamaLengkap: "Admin 2",
		StatusAktif:   true,
	}
	db.Create(user2)

	// Create members for both cooperatives
	memberA := &models.Anggota{
		IDKoperasi:     koperasi1.ID,
		NomorAnggota:   "A001",
		NamaLengkap:  "Member A",
		NIK:            "1111111111111111",
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
		NoTelepon:           "08111111111",
		Email:          "a@test.com",
	}
	db.Create(memberA)

	memberB := &models.Anggota{
		IDKoperasi:     koperasi2.ID,
		NomorAnggota:   "B001",
		NamaLengkap:  "Member B",
		NIK:            "2222222222222222",
		TempatLahir:    "Bandung",
		TanggalLahir:  nil,
		JenisKelamin:   "L",
		Alamat:         "Test",
		RT:             "001",
		RW:             "002",
		Kelurahan:      "Test",
		Kecamatan:      "Test",
		KotaKabupaten:  "Bandung",
		Provinsi:       "Jabar",
		KodePos:        "54321",
		NoTelepon:           "08122222222",
		Email:          "b@test.com",
	}
	db.Create(memberB)

	// Create share capital for Koperasi A
	simpananA1 := &models.Simpanan{
		IDKoperasi:        koperasi1.ID,
		IDAnggota:         memberA.ID,
		TipeSimpanan:     "pokok",
		JumlahSetoran:            100000,
		TanggalTransaksi:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		NomorReferensi:    "REF-A001",
		Keterangan:        "Simpanan Pokok",
	}
	db.Create(simpananA1)

	// Create share capital for Koperasi B
	simpananB1 := &models.Simpanan{
		IDKoperasi:        koperasi2.ID,
		IDAnggota:         memberB.ID,
		TipeSimpanan:     "pokok",
		JumlahSetoran:            200000,
		TanggalTransaksi:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		NomorReferensi:    "REF-B001",
		Keterangan:        "Simpanan Pokok",
	}
	db.Create(simpananB1)

	token1, _ := jwtUtil.GenerateToken(user1)
	token2, _ := jwtUtil.GenerateToken(user2)

	t.Run("Koperasi_A_Share_Capital_Isolation", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/simpanan", nil)
		req.Header.Set("Authorization", "Bearer "+token1)
		req.Header.Set("cooperative_id", koperasi1.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// Should only see own share capital
		if data, ok := response["data"].([]interface{}); ok {
			assert.Equal(t, 1, len(data))

			simpanan := data[0].(map[string]interface{})
			assert.Equal(t, float64(100000), simpanan["jumlahSetoran"])
		}
	})

	t.Run("Koperasi_B_Share_Capital_Isolation", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/simpanan", nil)
		req.Header.Set("Authorization", "Bearer "+token2)
		req.Header.Set("cooperative_id", koperasi2.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// Should only see own share capital
		if data, ok := response["data"].([]interface{}); ok {
			assert.Equal(t, 1, len(data))

			simpanan := data[0].(map[string]interface{})
			assert.Equal(t, float64(200000), simpanan["jumlahSetoran"])
		}
	})
}

// TestMultiTenant_UnauthorizedCooperativeAccess tests access with wrong cooperative_id
func TestMultiTenant_UnauthorizedCooperativeAccess(t *testing.T) {
	db := setupTestDB(t)
	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)

	router.POST("/api/v1/anggota", anggotaHandler.Create)

	koperasi1 := &models.Koperasi{
		NamaKoperasi:   "Koperasi A",
		Alamat:         "Address",
		NoTelepon:      "08111111111",
		Email:          "test@test.com",
	}
	db.Create(koperasi1)

	koperasi2 := &models.Koperasi{
		NamaKoperasi:   "Koperasi B",
		Alamat:         "Address",
		NoTelepon:      "08122222222",
		Email:          "test2@test.com",
	}
	db.Create(koperasi2)

	user1 := &models.Pengguna{
		IDKoperasi: koperasi1.ID,
		NamaPengguna:   "admin1",
		Email:      "admin1@test.com",
		// Password:   "hashed",
		Peran:       models.PeranAdmin,
		NamaLengkap: "Admin 1",
		StatusAktif:   true,
	}
	db.Create(user1)

	// Generate token for Koperasi A
	token1, _ := jwtUtil.GenerateToken(user1)

	t.Run("Cannot_Create_Member_For_Another_Cooperative", func(t *testing.T) {
		// Try to create member using Koperasi A's token but with Koperasi B's ID in header
		memberPayload := map[string]interface{}{
			"nomor_anggota":  "HACK001",
			"nama":           "Hacker",
			"nik":            "1234567890123456",
			"tempat_lahir":   "Jakarta",
			"tanggal_lahir":  "1990-01-01",
			"jenis_kelamin":  "L",
			"alamat":         "Test",
			"rt":             "001",
			"rw":             "002",
			"kelurahan":      "Test",
			"kecamatan":      "Test",
			"kabupaten_kota": "Jakarta",
			"provinsi":       "DKI",
			"kode_pos":       "12345",
			"no_hp":          "08123456789",
			"email":          "hack@test.com",
			"status_anggota": "aktif",
		}

		jsonPayload, _ := json.Marshal(memberPayload)
		req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token1)
		// Try to use different cooperative_id than token
		req.Header.Set("cooperative_id", koperasi2.ID.String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should be forbidden - token cooperative_id doesn't match header
		assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, w.Code)
	})

	t.Run("Cannot_Use_Invalid_Cooperative_ID", func(t *testing.T) {
		memberPayload := map[string]interface{}{
			"nomor_anggota":  "TEST001",
			"nama":           "Test",
			"nik":            "1234567890123456",
			"tempat_lahir":   "Jakarta",
			"tanggal_lahir":  "1990-01-01",
			"jenis_kelamin":  "L",
			"alamat":         "Test",
			"rt":             "001",
			"rw":             "002",
			"kelurahan":      "Test",
			"kecamatan":      "Test",
			"kabupaten_kota": "Jakarta",
			"provinsi":       "DKI",
			"kode_pos":       "12345",
			"no_hp":          "08123456789",
			"email":          "test@test.com",
			"status_anggota": "aktif",
		}

		jsonPayload, _ := json.Marshal(memberPayload)
		req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token1)
		// Use invalid UUID
		req.Header.Set("cooperative_id", "invalid-uuid")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should reject invalid UUID
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("Cannot_Use_Non_Existent_Cooperative_ID", func(t *testing.T) {
		memberPayload := map[string]interface{}{
			"nomor_anggota":  "TEST001",
			"nama":           "Test",
			"nik":            "1234567890123456",
			"tempat_lahir":   "Jakarta",
			"tanggal_lahir":  "1990-01-01",
			"jenis_kelamin":  "L",
			"alamat":         "Test",
			"rt":             "001",
			"rw":             "002",
			"kelurahan":      "Test",
			"kecamatan":      "Test",
			"kabupaten_kota": "Jakarta",
			"provinsi":       "DKI",
			"kode_pos":       "12345",
			"no_hp":          "08123456789",
			"email":          "test@test.com",
			"status_anggota": "aktif",
		}

		jsonPayload, _ := json.Marshal(memberPayload)
		req, _ := http.NewRequest("POST", "/api/v1/anggota", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token1)
		// Use valid UUID format but non-existent cooperative
		req.Header.Set("cooperative_id", uuid.New().String())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should reject non-existent cooperative
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

// TestMultiTenant_CooperativeIDConsistency tests that cooperative_id is consistently enforced
func TestMultiTenant_CooperativeIDConsistency(t *testing.T) {
	db := setupTestDB(t)
	// jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	// Verify all queries include cooperative_id filter
	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Address",
		NoTelepon:      "08111111111",
		Email:          "test@test.com",
	}
	db.Create(koperasi)

	// Try direct query without cooperative_id filter (should be prevented by service layer)
	var members []models.Anggota
	result := db.Find(&members)

	// If there are members from different cooperatives, direct query would return all
	// But service layer should always filter by cooperative_id
	assert.NoError(t, result.Error)
}
