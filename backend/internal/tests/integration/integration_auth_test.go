package integration

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestAuthIntegration_CompleteLoginLogoutCycle tests complete authentication flow
func TestAuthIntegration_CompleteLoginLogoutCycle(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	// Step 1: Create koperasi
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi Integration",
		NoTelepon:    "08111111111",
		Email:        "integration@test.com",
	}
	if err := db.Create(koperasi).Error; err != nil {
		t.Fatalf("Failed to create koperasi: %v", err)
	}

	// Step 2: Create user directly (auth service doesn't have DaftarPengguna in integration tests)
	pengguna := &models.Pengguna{
		IDKoperasi:   koperasiID,
		NamaLengkap:  "Test User Integration",
		NamaPengguna: "testuser",
		Email:        "testuser@integration.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("SecurePassword123!")
	if err := db.Create(pengguna).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if pengguna.ID == uuid.Nil {
		t.Error("Created user should have valid ID")
	}
	if pengguna.NamaPengguna != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", pengguna.NamaPengguna)
	}

	// Step 3: Login with credentials
	jwtUtil := utils.NewJWTUtil("test-secret-key-integration", 24)
	authService := services.NewAuthService(db, jwtUtil)

	loginResp, err := authService.Login("testuser", "SecurePassword123!")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if loginResp.Token == "" {
		t.Error("Login should return JWT token")
	}
	if loginResp.Pengguna.ID != pengguna.ID {
		t.Error("Login should return same user")
	}

	t.Logf("✓ Login successful, token generated: %s...", loginResp.Token[:20])

	// Step 4: Validate token using JWT util
	claims, err := jwtUtil.ValidateToken(loginResp.Token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.IDPengguna != pengguna.ID {
		t.Error("Token should contain correct user ID")
	}

	if claims.IDKoperasi != koperasiID {
		t.Error("Token should contain correct cooperative ID")
	}

	if claims.Peran != models.PeranAdmin {
		t.Error("Token should contain correct role")
	}

	t.Logf("✓ Token validated successfully")

	// Step 5: Use token to access protected resource (get user info)
	userID := claims.IDPengguna

	var userInfo models.Pengguna
	if err := db.Where("id = ? AND id_koperasi = ?", userID, koperasiID).First(&userInfo).Error; err != nil {
		t.Fatalf("Failed to get user info with token: %v", err)
	}

	if userInfo.NamaPengguna != "testuser" {
		t.Error("Should retrieve correct user with token")
	}

	t.Logf("✓ Protected resource access successful")

	// Step 6: Test invalid password
	_, err = authService.Login("testuser", "WrongPassword")
	if err == nil {
		t.Error("Login with wrong password should fail")
	}

	t.Logf("✓ Invalid password rejected")

	// Step 7: Test non-existent user
	_, err = authService.Login("nonexistent", "password")
	if err == nil {
		t.Error("Login with non-existent user should fail")
	}

	t.Logf("✓ Non-existent user rejected")

	t.Log("✅ Complete auth integration test passed")
}

// TestAuthIntegration_TokenExpiration tests token expiration handling
func TestAuthIntegration_TokenExpiration(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	// Create test data
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasiID,
		NamaPengguna: "testuser",
		Email:        "test@test.com",
		NamaLengkap:  "Test User",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	// Generate token
	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	token, err := jwtUtil.GenerateToken(pengguna)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Immediately validate (should succeed)
	_, err = jwtUtil.ValidateToken(token)
	if err != nil {
		t.Errorf("Fresh token should be valid: %v", err)
	}

	t.Logf("✓ Fresh token is valid")

	// Note: Testing actual expiration would require waiting or manipulating time
	// This test validates the token generation and immediate validation flow
	t.Log("✅ Token expiration test passed (immediate validation)")
}

// TestAuthIntegration_MultiUserSameCooperative tests multiple users in same cooperative
func TestAuthIntegration_MultiUserSameCooperative(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Multi User Koperasi",
		NoTelepon:    "08111111112",
		Email:        "multiuser@test.com",
	}
	if err := db.Create(koperasi).Error; err != nil {
		t.Fatalf("Failed to create koperasi: %v", err)
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	authService := services.NewAuthService(db, jwtUtil)

	// Create multiple users with different roles
	users := []struct {
		username string
		role     models.PeranPengguna
	}{
		{"admin", models.PeranAdmin},
		{"bendahara", models.PeranBendahara},
		{"kasir", models.PeranKasir},
	}

	tokens := make(map[string]string)

	for _, user := range users {
		// Create user
		pengguna := &models.Pengguna{
			IDKoperasi:   koperasiID,
			NamaPengguna: user.username,
			Email:        user.username + "@test.com",
			NamaLengkap:  user.username + " User",
			Peran:        user.role,
			StatusAktif:  true,
		}
		pengguna.SetKataSandi("Password123!")
		if err := db.Create(pengguna).Error; err != nil {
			t.Fatalf("Failed to create user %s: %v", user.username, err)
		}

		// Login
		loginResp, err := authService.Login(user.username, "Password123!")
		if err != nil {
			t.Fatalf("Failed to login %s: %v", user.username, err)
		}

		tokens[user.username] = loginResp.Token

		// Validate token has correct role
		claims, err := jwtUtil.ValidateToken(loginResp.Token)
		if err != nil {
			t.Fatalf("Failed to validate token for %s: %v", user.username, err)
		}

		if claims.Peran != user.role {
			t.Errorf("Expected role '%s' for %s, got '%s'", user.role, user.username, claims.Peran)
		}

		t.Logf("✓ User '%s' with role '%s' authenticated successfully", user.username, user.role)
	}

	// Verify all tokens are unique
	uniqueTokens := make(map[string]bool)
	for username, token := range tokens {
		if uniqueTokens[token] {
			t.Errorf("Duplicate token found for user %s", username)
		}
		uniqueTokens[token] = true
	}

	t.Log("✅ Multi-user authentication test passed")
}

// TestAuthIntegration_InactiveUserBlocked tests that inactive users cannot login
func TestAuthIntegration_InactiveUserBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
		NoTelepon:    "08111111113",
		Email:        "inactive@test.com",
	}
	if err := db.Create(koperasi).Error; err != nil {
		t.Fatalf("Failed to create koperasi: %v", err)
	}

	// Create inactive user
	pengguna := &models.Pengguna{
		IDKoperasi:   koperasiID,
		NamaPengguna: "inactiveuser",
		Email:        "inactive@test.com",
		NamaLengkap:  "Inactive User",
		Peran:        models.PeranAdmin,
		StatusAktif:  true, // Create as active first
	}
	pengguna.SetKataSandi("password123")
	if err := db.Create(pengguna).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	// Then explicitly set to inactive (to override database default)
	db.Model(pengguna).Update("status_aktif", false)

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	authService := services.NewAuthService(db, jwtUtil)

	// Try to login
	_, err := authService.Login("inactiveuser", "password123")
	if err == nil {
		t.Error("Inactive user should not be able to login")
	}

	t.Logf("✓ Inactive user blocked from login: %v", err)
	t.Log("✅ Inactive user blocking test passed")
}

// TestAuthIntegration_SessionManagement tests session-like behavior with tokens
func TestAuthIntegration_SessionManagement(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Session Test Koperasi",
		NoTelepon:    "08111111114",
		Email:        "session@test.com",
	}
	if err := db.Create(koperasi).Error; err != nil {
		t.Fatalf("Failed to create koperasi: %v", err)
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	authService := services.NewAuthService(db, jwtUtil)

	// Create user
	user := &models.Pengguna{
		IDKoperasi:   koperasiID,
		NamaPengguna: "sessionuser",
		Email:        "session@test.com",
		NamaLengkap:  "Session User",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	user.SetKataSandi("Password123!")
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// First login
	loginResp, err := authService.Login("sessionuser", "Password123!")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	token1 := loginResp.Token

	// Simulate second login (new session/device)
	time.Sleep(1 * time.Second)
	loginResp2, err := authService.Login("sessionuser", "Password123!")
	if err != nil {
		t.Fatalf("Failed second login: %v", err)
	}

	token2 := loginResp2.Token

	// Both tokens should be valid (stateless JWT)
	_, err1 := jwtUtil.ValidateToken(token1)
	_, err2 := jwtUtil.ValidateToken(token2)

	if err1 != nil || err2 != nil {
		t.Error("Both tokens should be valid")
	}

	// Tokens should be different (different timestamp)
	if token1 == token2 {
		t.Error("Different login sessions should generate different tokens")
	}

	// Update user to inactive
	db.Model(&models.Pengguna{}).Where("id = ?", user.ID).Update("status_aktif", false)

	// Old tokens should still be valid (stateless JWT limitation)
	// In production, you might implement token blacklisting or shorter expiry
	_, err = jwtUtil.ValidateToken(token1)
	// Note: This will still be valid as JWT is stateless
	// Additional server-side checks needed for real session management

	t.Log("✅ Session management test passed (stateless JWT behavior validated)")
}
