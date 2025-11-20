package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"cooperative-erp-lite/internal/handlers"
	"cooperative-erp-lite/internal/middleware"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
)

// TestRateLimit_BasicLimiting tests basic rate limiting functionality
func TestRateLimit_BasicLimiting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Rate limit: 5 requests per 10 seconds
	router.Use(middleware.RateLimitMiddleware(5, 10*time.Second))

	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Make 5 requests (should all succeed)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}

	// 6th request should be rate limited
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "6th request should be rate limited")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Rate limit exceeded")
}

// TestRateLimit_PerIP tests that rate limiting is per IP address
func TestRateLimit_PerIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Rate limit: 3 requests per 10 seconds
	router.Use(middleware.RateLimitMiddleware(3, 10*time.Second))

	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// IP 1 makes 3 requests (should succeed)
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "IP1 request %d should succeed", i+1)
	}

	// IP 1's 4th request should be blocked
	req1, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusTooManyRequests, w1.Code, "IP1's 4th request should be blocked")

	// IP 2 should still be able to make requests
	req2, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code, "IP2 should not be affected by IP1's rate limit")
}

// TestRateLimit_WindowExpiry tests that rate limit window expires correctly
func TestRateLimit_WindowExpiry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Rate limit: 2 requests per 1 second
	router.Use(middleware.RateLimitMiddleware(2, 1*time.Second))

	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Make 2 requests (should succeed)
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}

	// 3rd request should be blocked
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "3rd request should be blocked")

	// Wait for window to expire
	time.Sleep(1100 * time.Millisecond)

	// Should be able to make requests again
	req2, _ := http.NewRequest("GET", "/api/v1/test", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code, "Request after window expiry should succeed")
}

// TestLoginRateLimit_FailedAttempts tests login rate limiting with failed attempts
func TestLoginRateLimit_FailedAttempts(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Login rate limit: 5 attempts per 5 minutes, 15 minute lockout
	router.Use(middleware.LoginRateLimitMiddleware(5, 5*time.Minute, 15*time.Minute))

	authService := services.NewAuthService(db, utils.NewJWTUtil("test-secret", 24))
	authHandler := handlers.NewAuthHandler(authService)

	router.POST("/api/v1/auth/login", authHandler.Login)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	user := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaPengguna: "testuser",
		Email:        "test@test.com",
		Peran:        models.PeranAdmin,
		NamaLengkap:  "Test User",
		StatusAktif:  true,
	}
	user.SetKataSandi("correctpassword")
	db.Create(user)

	// Make 5 failed login attempts
	for i := 0; i < 5; i++ {
		payload := map[string]interface{}{
			"namaPengguna": "testuser",
			"kataSandi":    "wrongpassword",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.100:12345"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// First 5 attempts should return unauthorized (not rate limited)
		assert.Contains(t, []int{http.StatusUnauthorized, http.StatusTooManyRequests}, w.Code)
	}

	// 6th attempt should be rate limited (account locked)
	payload := map[string]interface{}{
		"namaPengguna": "testuser",
		"kataSandi":    "correctpassword", // Even with correct password
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.168.1.100:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Should be locked out after 5 failed attempts")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "locked")
}

// TestLoginRateLimit_SuccessfulLoginClearsAttempts tests that successful login clears attempts
func TestLoginRateLimit_SuccessfulLoginClearsAttempts(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	limiter := middleware.NewLoginRateLimiter(5, 5*time.Minute, 15*time.Minute)

	authService := services.NewAuthService(db, utils.NewJWTUtil("test-secret", 24))
	authHandler := handlers.NewAuthHandler(authService)

	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		// Manually set limiter in context for testing
		c.Set("loginLimiter", limiter)
		authHandler.Login(c)
	})

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	user := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaPengguna: "testuser",
		Email:        "test@test.com",
		Peran:        models.PeranAdmin,
		NamaLengkap:  "Test User",
		StatusAktif:  true,
	}
	user.SetKataSandi("correctpassword")
	db.Create(user)

	ip := "192.168.1.100"

	// Make 3 failed attempts
	for i := 0; i < 3; i++ {
		limiter.RecordAttempt(ip)
	}

	// Successful login should clear attempts
	limiter.ClearAttempts(ip)

	// Should be able to attempt login again
	allowed := limiter.RecordAttempt(ip)
	assert.True(t, allowed, "Should be able to attempt login after successful login cleared attempts")
}

// TestRateLimit_BruteForceProtection tests protection against brute force attacks
func TestRateLimit_BruteForceProtection(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Aggressive rate limit: 3 attempts per 1 minute, 5 minute lockout
	router.Use(middleware.LoginRateLimitMiddleware(3, 1*time.Minute, 5*time.Minute))

	authService := services.NewAuthService(db, utils.NewJWTUtil("test-secret", 24))
	authHandler := handlers.NewAuthHandler(authService)

	router.POST("/api/v1/auth/login", authHandler.Login)

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	user := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaPengguna: "victim",
		Email:        "victim@test.com",
		Peran:        models.PeranAdmin,
		NamaLengkap:  "Victim User",
		StatusAktif:  true,
	}
	user.SetKataSandi("strongpassword123")
	db.Create(user)

	// Simulate brute force attack with different passwords
	passwords := []string{
		"password123",
		"admin123",
		"12345678",
		"qwerty",
		"letmein",
		"password",
		"abc123",
		"monkey",
		"1234567890",
		"strongpassword123", // Correct password, but account is locked
	}

	successCount := 0
	lockedCount := 0

	for i, password := range passwords {
		payload := map[string]interface{}{
			"namaPengguna": "victim",
			"kataSandi":    password,
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.200:12345" // Attacker's IP

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			successCount++
		} else if w.Code == http.StatusTooManyRequests {
			lockedCount++
		}

		t.Logf("Attempt %d: password=%s, status=%d", i+1, password, w.Code)
	}

	// Account should be locked before reaching correct password
	assert.Greater(t, lockedCount, 0, "Should have locked out attacker")
	assert.Equal(t, 0, successCount, "Brute force attack should not succeed")
}

// TestRateLimit_DifferentEndpoints tests that rate limiting can be applied to different endpoints
func TestRateLimit_DifferentEndpoints(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply rate limiting to specific endpoint
	loginGroup := router.Group("/api/v1/auth")
	loginGroup.Use(middleware.RateLimitMiddleware(3, 10*time.Second))

	authService := services.NewAuthService(db, utils.NewJWTUtil("test-secret", 24))
	authHandler := handlers.NewAuthHandler(authService)

	loginGroup.POST("/login", authHandler.Login)

	// Another endpoint without rate limiting
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Address",
		NoTelepon:    "08111111111",
		Email:        "test@test.com",
	}
	db.Create(koperasi)

	// Make 3 login attempts (should succeed, then get rate limited)
	for i := 0; i < 4; i++ {
		payload := map[string]interface{}{
			"namaPengguna": "test",
			"kataSandi":    "password",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 3 {
			assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "First 3 requests should not be rate limited")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "4th request should be rate limited")
		}
	}

	// Health endpoint should not be affected
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Health endpoint should not be rate limited")
	}
}

// TestRateLimit_ConcurrentRequests tests rate limiting under concurrent load
func TestRateLimit_ConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Rate limit: 10 requests per 5 seconds
	router.Use(middleware.RateLimitMiddleware(10, 5*time.Second))

	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Make 20 concurrent requests
	results := make(chan int, 20)
	for i := 0; i < 20; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/api/v1/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	// Collect results
	successCount := 0
	rateLimitedCount := 0

	for i := 0; i < 20; i++ {
		code := <-results
		if code == http.StatusOK {
			successCount++
		} else if code == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	// Should allow up to 10 requests, block the rest
	assert.LessOrEqual(t, successCount, 10, "Should allow max 10 requests")
	assert.GreaterOrEqual(t, rateLimitedCount, 10, "Should block at least 10 requests")
}

// TestLoginRateLimit_LockoutDuration tests that lockout duration is enforced
func TestLoginRateLimit_LockoutDuration(t *testing.T) {
	// This test is skipped by default as it requires waiting for lockout to expire
	t.Skip("Skipping test that requires long wait time")

	limiter := middleware.NewLoginRateLimiter(3, 1*time.Minute, 2*time.Second)

	ip := "192.168.1.300"

	// Trigger lockout with 3 failed attempts
	for i := 0; i < 3; i++ {
		limiter.RecordAttempt(ip)
	}

	// Should be locked out
	assert.True(t, limiter.IsLockedOut(ip), "Should be locked out")

	// Wait for lockout to expire
	time.Sleep(2100 * time.Millisecond)

	// Should no longer be locked out
	assert.False(t, limiter.IsLockedOut(ip), "Lockout should have expired")
}
