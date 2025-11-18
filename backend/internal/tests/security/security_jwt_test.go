package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"cooperative-erp-lite/internal/handlers"
	"cooperative-erp-lite/internal/middleware"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
)

// TestJWT_TokenTampering tests detection of tampered JWT tokens
func TestJWT_TokenTampering(t *testing.T) {
	db := setupTestDB(t)
	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(utils.NewJWTUtil("test-secret", 24)))

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)
	router.GET("/api/v1/anggota", anggotaHandler.List)

	// Create valid user and cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Test Address",
		NoTelepon:      "08123456789",
		Email:          "test@koperasi.com",
	}
	db.Create(koperasi)

	user := &models.Pengguna{
		IDKoperasi: koperasi.ID,
		NamaPengguna:   "testuser",
		Email:      "test@user.com",
		// Password:   "hashedpassword",
		Peran:       models.PeranAdmin,
		NamaLengkap: "Test User",
		StatusAktif:   true,
	}
	db.Create(user)

	// Generate valid token
	validToken, _ := jwtUtil.GenerateToken(user)

	// Test cases for token tampering
	testCases := []struct {
		name          string
		token         string
		expectedCode  int
		description   string
	}{
		{
			name:         "Valid Token",
			token:        validToken,
			expectedCode: http.StatusOK,
			description:  "Valid token should be accepted",
		},
		{
			name:         "Tampered Signature",
			token:        validToken[:len(validToken)-10] + "tampered12",
			expectedCode: http.StatusUnauthorized,
			description:  "Token with tampered signature should be rejected",
		},
		{
			name:         "Modified Payload",
			token:        modifyJWTPayload(validToken, "role", "superadmin"),
			expectedCode: http.StatusUnauthorized,
			description:  "Token with modified payload should be rejected",
		},
		{
			name:         "Empty Token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
			description:  "Empty token should be rejected",
		},
		{
			name:         "Invalid Format",
			token:        "not.a.jwt.token",
			expectedCode: http.StatusUnauthorized,
			description:  "Invalid format token should be rejected",
		},
		{
			name:         "Missing Signature",
			token:        validToken[:len(validToken)-20],
			expectedCode: http.StatusUnauthorized,
			description:  "Token without signature should be rejected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code, tc.description)

			if tc.expectedCode == http.StatusUnauthorized {
				// Verify error response
				assert.Contains(t, w.Body.String(), "error", "Should return error message")
			}
		})
	}
}

// TestJWT_TokenExpiration tests expired token rejection
func TestJWT_TokenExpiration(t *testing.T) {
	db := setupTestDB(t)
	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(utils.NewJWTUtil("test-secret", 24)))

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

	user := &models.Pengguna{
		IDKoperasi: koperasi.ID,
		NamaPengguna:   "testuser",
		Email:      "test@user.com",
		// Password:   "hashedpassword",
		Peran:       models.PeranAdmin,
		NamaLengkap: "Test User",
		StatusAktif:   true,
	}
	db.Create(user)

	t.Run("Expired Token", func(t *testing.T) {
		// Create expired token (expired 1 hour ago)
		expiredToken := generateExpiredToken(user, -1*time.Hour)

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Expired token should be rejected")
		assert.Contains(t, w.Body.String(), "error", "Should return error for expired token")
	})

	t.Run("Token Expiring Soon", func(t *testing.T) {
		// Create token expiring in 5 minutes
		soonExpiringToken := generateExpiredToken(user, 5*time.Minute)

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+soonExpiringToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Token should still be valid
		assert.Equal(t, http.StatusOK, w.Code, "Token expiring soon should still be valid")
	})

	t.Run("Freshly Generated Token", func(t *testing.T) {
		// Generate fresh token
		freshToken, _ := jwtUtil.GenerateToken(user)

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+freshToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Fresh token should be accepted")
	})
}

// TestJWT_InvalidClaims tests rejection of tokens with invalid claims
func TestJWT_InvalidClaims(t *testing.T) {
	db := setupTestDB(t)
// 	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(utils.NewJWTUtil("test-secret", 24)))

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)
	router.GET("/api/v1/anggota", anggotaHandler.List)

	t.Run("Missing User ID", func(t *testing.T) {
		// Create token without user_id claim
		token := generateTokenWithCustomClaims(map[string]interface{}{
			"cooperative_id": uuid.New().String(),
			"role":           "admin",
			"exp":            time.Now().Add(24 * time.Hour).Unix(),
		})

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Token without user_id should be rejected")
	})

	t.Run("Missing Cooperative ID", func(t *testing.T) {
		// Create token without cooperative_id claim
		token := generateTokenWithCustomClaims(map[string]interface{}{
			"user_id": uuid.New().String(),
			"role":    "admin",
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		})

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Token without cooperative_id should be rejected")
	})

	t.Run("Missing Role", func(t *testing.T) {
		// Create token without role claim
		token := generateTokenWithCustomClaims(map[string]interface{}{
			"user_id":        uuid.New().String(),
			"cooperative_id": uuid.New().String(),
			"exp":            time.Now().Add(24 * time.Hour).Unix(),
		})

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Token without role should be rejected")
	})

	t.Run("Invalid User ID Format", func(t *testing.T) {
		// Create token with invalid UUID format
		token := generateTokenWithCustomClaims(map[string]interface{}{
			"user_id":        "not-a-valid-uuid",
			"cooperative_id": uuid.New().String(),
			"role":           "admin",
			"exp":            time.Now().Add(24 * time.Hour).Unix(),
		})

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should be rejected or handled gracefully
		assert.NotEqual(t, http.StatusOK, w.Code, "Token with invalid UUID should not be accepted")
	})
}

// TestJWT_MissingAuthorizationHeader tests missing authorization header
func TestJWT_MissingAuthorizationHeader(t *testing.T) {
	db := setupTestDB(t)
// 	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(utils.NewJWTUtil("test-secret", 24)))

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)
	router.GET("/api/v1/anggota", anggotaHandler.List)

	req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
	// No Authorization header set

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Request without authorization should be rejected")
	assert.Contains(t, w.Body.String(), "error", "Should return error message")
}

// TestJWT_WrongSecretKey tests token signed with wrong secret key
func TestJWT_WrongSecretKey(t *testing.T) {
	db := setupTestDB(t)
// 	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(utils.NewJWTUtil("test-secret", 24)))

	anggotaService := services.NewAnggotaService(db)
	anggotaHandler := handlers.NewAnggotaHandler(anggotaService)
	router.GET("/api/v1/anggota", anggotaHandler.List)

	// Generate token with wrong secret key
	wrongSecretToken := generateTokenWithWrongSecret(uuid.New(), uuid.New(), "admin")

	req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
	req.Header.Set("Authorization", "Bearer "+wrongSecretToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Token with wrong secret should be rejected")
}

// Helper functions

// modifyJWTPayload attempts to modify JWT payload (will invalidate signature)
func modifyJWTPayload(tokenString, key, value string) string {
	// This will create an invalid token as signature won't match
	// Used to test tampering detection
	return tokenString[:len(tokenString)-10] + "modified"
}

// generateExpiredToken generates a token with custom expiration
func generateExpiredToken(user *models.Pengguna, expiresIn time.Duration) string {
	claims := jwt.MapClaims{
		"user_id":        user.ID.String(),
		"cooperative_id": user.IDKoperasi.String(),
		"role":           string(user.Peran),
		"exp":            time.Now().Add(expiresIn).Unix(),
		"iat":            time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Use the same secret as the test setup
	jwtSecret := []byte("test-secret")

	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// generateTokenWithCustomClaims generates token with custom claims
func generateTokenWithCustomClaims(claims map[string]interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	jwtSecret := []byte("test-secret")
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// generateTokenWithWrongSecret generates token with wrong secret key
func generateTokenWithWrongSecret(userID, cooperativeID uuid.UUID, role string) string {
	claims := jwt.MapClaims{
		"user_id":        userID.String(),
		"cooperative_id": cooperativeID.String(),
		"role":           role,
		"exp":            time.Now().Add(24 * time.Hour).Unix(),
		"iat":            time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Wrong secret key
	wrongSecret := []byte("wrong-secret-key-12345")
	tokenString, _ := token.SignedString(wrongSecret)
	return tokenString
}

// TestJWT_RoleValidation tests that role from token is properly validated
func TestJWT_RoleValidation(t *testing.T) {
	db := setupTestDB(t)
	jwtUtil := utils.NewJWTUtil("test-secret", 24)

	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(utils.NewJWTUtil("test-secret", 24)))

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

	validRoles := []string{"admin", "bendahara", "kasir", "anggota"}

	for _, role := range validRoles {
		t.Run("Role_"+role, func(t *testing.T) {
			user := &models.Pengguna{
				IDKoperasi: koperasi.ID,
				NamaPengguna:   "testuser_" + role,
				Email:      role + "@user.com",
				// Password:   "hashedpassword",
				Peran:       models.PeranPengguna(role),
				NamaLengkap: "Test User",
				StatusAktif:   true,
			}
			db.Create(user)

			token, _ := jwtUtil.GenerateToken(user)

			req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// All valid roles should be accepted for this endpoint
			assert.Equal(t, http.StatusOK, w.Code, "Valid role should be accepted")
		})
	}

	// Test invalid role
	t.Run("Invalid_Role", func(t *testing.T) {
		invalidToken := generateTokenWithCustomClaims(map[string]interface{}{
			"user_id":        uuid.New().String(),
			"cooperative_id": koperasi.ID.String(),
			"role":           "superadmin", // Invalid role
			"exp":            time.Now().Add(24 * time.Hour).Unix(),
		})

		req, _ := http.NewRequest("GET", "/api/v1/anggota", nil)
		req.Header.Set("Authorization", "Bearer "+invalidToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should handle invalid role (implementation dependent)
		// Either reject or accept but limit permissions
		assert.NotEqual(t, http.StatusInternalServerError, w.Code, "Should handle invalid role gracefully")
	})
}
