package integration

import (
	"bytes"
	"cooperative-erp-lite/internal/config"
	"cooperative-erp-lite/internal/handlers"
	"cooperative-erp-lite/internal/middleware"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// setupPortalTestRouter membuat router untuk testing portal anggota
func setupPortalTestRouter(t *testing.T) (*gin.Engine, *utils.JWTUtil, uuid.UUID, uuid.UUID) {
	gin.SetMode(gin.TestMode)

	db := config.GetDB()
	require.NotNil(t, db)

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)

	// Initialize services
	portalService := services.NewPortalAnggotaService(db, jwtUtil)
	portalHandler := handlers.NewPortalAnggotaHandler(portalService)

	router := gin.New()
	router.Use(gin.Recovery())

	// Portal routes
	portal := router.Group("/api/v1/portal")
	{
		// Public login
		portal.POST("/login", portalHandler.Login)

		// Protected routes
		protected := portal.Group("")
		protected.Use(middleware.AuthAnggotaMiddleware(jwtUtil))
		{
			protected.GET("/profile", portalHandler.GetProfile)
			protected.GET("/saldo", portalHandler.GetSaldo)
			protected.GET("/riwayat", portalHandler.GetRiwayat)
			protected.PUT("/ubah-pin", portalHandler.UbahPIN)
		}
	}

	// Create test koperasi
	koperasi := models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Koperasi Test Portal",
	}
	err := db.Create(&koperasi).Error
	require.NoError(t, err)

	// Create test anggota with PIN
	pin := "123456"
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	require.NoError(t, err)

	tanggalBergabung, _ := time.Parse("2006-01-02", "2024-01-01")

	anggota := models.Anggota{
		ID:               uuid.New(),
		IDKoperasi:       koperasi.ID,
		NomorAnggota:     "A001",
		NamaLengkap:      "Test Member Portal",
		TanggalBergabung: tanggalBergabung,
		Status:           models.StatusAktif,
		PINPortal:        string(hashedPIN),
	}
	err = db.Create(&anggota).Error
	require.NoError(t, err)

	// Cleanup function
	t.Cleanup(func() {
		db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", koperasi.ID)
		db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
		db.Exec("DELETE FROM koperasi WHERE id = ?", koperasi.ID)
	})

	return router, jwtUtil, koperasi.ID, anggota.ID
}

// TestPortalAnggotaLogin tests member portal login
func TestPortalAnggotaLogin(t *testing.T) {
	router, _, koperasiID, _ := setupPortalTestRouter(t)

	tests := []struct {
		name           string
		koperasiID     string
		nomorAnggota   string
		pin            string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "Successful login",
			koperasiID:     koperasiID.String(),
			nomorAnggota:   "A001",
			pin:            "123456",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["token"])
				anggota := data["anggota"].(map[string]interface{})
				assert.Equal(t, "A001", anggota["nomorAnggota"])
				assert.Equal(t, "Test Member Portal", anggota["namaLengkap"])
			},
		},
		{
			name:           "Invalid PIN",
			koperasiID:     koperasiID.String(),
			nomorAnggota:   "A001",
			pin:            "wrong-pin",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "Invalid nomor anggota",
			koperasiID:     koperasiID.String(),
			nomorAnggota:   "INVALID",
			pin:            "123456",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "Missing koperasi ID",
			koperasiID:     "",
			nomorAnggota:   "A001",
			pin:            "123456",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]string{
				"nomorAnggota": tt.nomorAnggota,
				"pin":          tt.pin,
			}
			bodyBytes, _ := json.Marshal(body)

			url := "/api/v1/portal/login"
			if tt.koperasiID != "" {
				url += "?idKoperasi=" + tt.koperasiID
			}

			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			tt.checkResponse(t, response)
		})
	}
}

// TestPortalAnggotaProfile tests getting member profile
func TestPortalAnggotaProfile(t *testing.T) {
	router, jwtUtil, koperasiID, anggotaID := setupPortalTestRouter(t)

	// Create test anggota for token
	db := config.GetDB()
	var anggota models.Anggota
	err := db.Where("id = ?", anggotaID).First(&anggota).Error
	require.NoError(t, err)

	// Generate token
	token, err := jwtUtil.GenerateTokenAnggota(&anggota)
	require.NoError(t, err)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "Get profile successfully",
			token:          token,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "A001", data["nomorAnggota"])
				assert.Equal(t, "Test Member Portal", data["namaLengkap"])
				assert.Equal(t, koperasiID.String(), data["idKoperasi"])
			},
		},
		{
			name:           "Missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "Invalid token",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/profile", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			tt.checkResponse(t, response)
		})
	}
}

// TestPortalAnggotaSaldo tests getting member balance
func TestPortalAnggotaSaldo(t *testing.T) {
	router, jwtUtil, koperasiID, anggotaID := setupPortalTestRouter(t)

	db := config.GetDB()

	tanggalTransaksi, _ := time.Parse("2006-01-02", "2024-01-01")

	// Create test simpanan
	simpanan := []models.Simpanan{
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			TipeSimpanan:     models.SimpananPokok,
			TanggalTransaksi: tanggalTransaksi,
			JumlahSetoran:    1000000,
			Keterangan:       "Simpanan Pokok",
			NomorReferensi:   "SPK-001",
		},
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			TipeSimpanan:     models.SimpananWajib,
			TanggalTransaksi: tanggalTransaksi,
			JumlahSetoran:    500000,
			Keterangan:       "Simpanan Wajib Januari",
			NomorReferensi:   "SWJ-001",
		},
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			TipeSimpanan:     models.SimpananSukarela,
			TanggalTransaksi: tanggalTransaksi,
			JumlahSetoran:    200000,
			Keterangan:       "Simpanan Sukarela",
			NomorReferensi:   "SSK-001",
		},
	}

	for _, s := range simpanan {
		err := db.Create(&s).Error
		require.NoError(t, err)
	}

	// Get anggota for token
	var anggota models.Anggota
	err := db.Where("id = ?", anggotaID).First(&anggota).Error
	require.NoError(t, err)

	token, err := jwtUtil.GenerateTokenAnggota(&anggota)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/portal/saldo", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1000000), data["simpananPokok"])
	assert.Equal(t, float64(500000), data["simpananWajib"])
	assert.Equal(t, float64(200000), data["simpananSukarela"])
	assert.Equal(t, float64(1700000), data["totalSimpanan"])
	assert.Equal(t, "A001", data["nomorAnggota"])
	assert.Equal(t, "Test Member Portal", data["namaAnggota"])
}

// TestPortalAnggotaRiwayat tests getting transaction history
func TestPortalAnggotaRiwayat(t *testing.T) {
	router, jwtUtil, koperasiID, anggotaID := setupPortalTestRouter(t)

	db := config.GetDB()

	tanggalTransaksi, _ := time.Parse("2006-01-02", "2024-01-01")

	// Create multiple transactions
	for i := 1; i <= 15; i++ {
		simpanan := models.Simpanan{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			TipeSimpanan:     models.SimpananWajib,
			TanggalTransaksi: tanggalTransaksi,
			JumlahSetoran:    float64(i * 100000),
			Keterangan:       "Test Transaction",
			NomorReferensi:   "TEST-" + string(rune(i)),
		}
		err := db.Create(&simpanan).Error
		require.NoError(t, err)
	}

	// Get anggota for token
	var anggota models.Anggota
	err := db.Where("id = ?", anggotaID).First(&anggota).Error
	require.NoError(t, err)

	token, err := jwtUtil.GenerateTokenAnggota(&anggota)
	require.NoError(t, err)

	tests := []struct {
		name          string
		page          string
		pageSize      string
		expectedCount int
	}{
		{
			name:          "First page default size",
			page:          "1",
			pageSize:      "10",
			expectedCount: 10,
		},
		{
			name:          "Second page",
			page:          "2",
			pageSize:      "10",
			expectedCount: 5,
		},
		{
			name:          "Large page size",
			page:          "1",
			pageSize:      "20",
			expectedCount: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/portal/riwayat?page=" + tt.page + "&pageSize=" + tt.pageSize
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.True(t, response["success"].(bool))
			data := response["data"].([]interface{})
			assert.Equal(t, tt.expectedCount, len(data))

			// Check pagination metadata
			pagination := response["pagination"].(map[string]interface{})
			assert.NotNil(t, pagination)
			assert.Equal(t, float64(15), pagination["totalItems"])
		})
	}
}

// TestPortalAnggotaUbahPIN tests changing member PIN
func TestPortalAnggotaUbahPIN(t *testing.T) {
	router, jwtUtil, _, anggotaID := setupPortalTestRouter(t)

	db := config.GetDB()
	var anggota models.Anggota
	err := db.Where("id = ?", anggotaID).First(&anggota).Error
	require.NoError(t, err)

	token, err := jwtUtil.GenerateTokenAnggota(&anggota)
	require.NoError(t, err)

	tests := []struct {
		name           string
		pinLama        string
		pinBaru        string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "Successful PIN change",
			pinLama:        "123456",
			pinBaru:        "654321",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
			},
		},
		{
			name:           "Wrong old PIN",
			pinLama:        "wrong",
			pinBaru:        "654321",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
			},
		},
		{
			name:           "Invalid new PIN length",
			pinLama:        "654321", // After previous test, PIN is now 654321
			pinBaru:        "123",     // Too short
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]string{
				"pinLama": tt.pinLama,
				"pinBaru": tt.pinBaru,
			}
			bodyBytes, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/portal/ubah-pin", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			tt.checkResponse(t, response)
		})
	}
}
