package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupAuthTestDB creates a test database for auth service
func setupAuthTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.Koperasi{}, &models.Pengguna{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean up existing data
	db.Exec("TRUNCATE TABLE pengguna CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	return db
}

// TestLogin_Success tests successful login
func TestLogin_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	if db == nil {
		return
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	// Create test cooperative and user
	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Email:        "test@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	// Test login
	result, err := service.Login("testuser", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "Test User", result.Pengguna.NamaLengkap)
	assert.Equal(t, "testuser", result.Pengguna.NamaPengguna)
}

// TestLogin_InvalidCredentials tests login with wrong credentials
func TestLogin_InvalidCredentials(t *testing.T) {
	db := setupAuthTestDB(t)
	if db == nil {
		return
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Email:        "test@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{"wrong password", "testuser", "wrongpassword", true},
		{"wrong username", "wronguser", "password123", true},
		{"empty password", "testuser", "", true},
		{"correct credentials", "testuser", "password123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(tt.username, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLogin_InactiveUser tests login with inactive user
func TestLogin_InactiveUser(t *testing.T) {
	db := setupAuthTestDB(t)
	if db == nil {
		return
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Inactive User",
		NamaPengguna: "inactiveuser",
		Email:        "inactive@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  false, // Inactive
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	_, err := service.Login("inactiveuser", "password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nama pengguna atau kata sandi salah")
}

// TestValidasiToken tests token validation
func TestValidasiToken(t *testing.T) {
	db := setupAuthTestDB(t)
	if db == nil {
		return
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Email:        "test@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	// Login to get token
	loginResult, _ := service.Login("testuser", "password123")

	t.Run("valid token", func(t *testing.T) {
		claims, err := service.ValidasiToken(loginResult.Token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, pengguna.ID, claims.IDPengguna)
		assert.Equal(t, koperasi.ID, claims.IDKoperasi)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := service.ValidasiToken("invalid.token.string")
		assert.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		// Create a token with very short expiry
		shortJWT := utils.NewJWTUtil("test-secret-key", 0)
		expiredToken, _ := shortJWT.GenerateToken(pengguna)
		time.Sleep(10 * time.Millisecond)

		_, err := service.ValidasiToken(expiredToken)
		assert.Error(t, err)
	})
}

// TestRefreshToken tests token refresh
func TestRefreshToken(t *testing.T) {
	db := setupAuthTestDB(t)
	if db == nil {
		return
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Email:        "test@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	// Login to get initial token
	loginResult, _ := service.Login("testuser", "password123")
	oldToken := loginResult.Token

	// Refresh token
	newToken, err := service.RefreshToken(oldToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEqual(t, oldToken, newToken) // Should be different
}

// TestUbahKataSandi tests password change
func TestUbahKataSandi(t *testing.T) {
	db := setupAuthTestDB(t)
	if db == nil {
		return
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Email:        "test@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("oldpassword")
	db.Create(pengguna)

	t.Run("successful password change", func(t *testing.T) {
		err := service.UbahKataSandi(pengguna.ID.String(), "oldpassword", "newpassword123")
		assert.NoError(t, err)

		// Verify can login with new password
		_, err = service.Login("testuser", "newpassword123")
		assert.NoError(t, err)

		// Verify cannot login with old password
		_, err = service.Login("testuser", "oldpassword")
		assert.Error(t, err)
	})

	t.Run("wrong old password", func(t *testing.T) {
		err := service.UbahKataSandi(pengguna.ID.String(), "wrongoldpassword", "newpassword")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kata sandi lama tidak sesuai")
	})

	t.Run("new password too short", func(t *testing.T) {
		err := service.UbahKataSandi(pengguna.ID.String(), "newpassword123", "12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimal 6 karakter")
	})
}

// TestDapatkanProfilPengguna tests getting user profile
func TestDapatkanProfilPengguna(t *testing.T) {
	db := setupAuthTestDB(t)
	if db == nil {
		return
	}

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Email:        "test@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	db.Create(pengguna)

	t.Run("existing user", func(t *testing.T) {
		profile, err := service.DapatkanProfilPengguna(pengguna.ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, "Test User", profile.NamaLengkap)
	})

	t.Run("non-existing user", func(t *testing.T) {
		profile, err := service.DapatkanProfilPengguna(uuid.New().String())
		assert.Error(t, err)
		assert.Nil(t, profile)
	})
}

// BenchmarkLogin benchmarks login operation
func BenchmarkLogin(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{}, &models.Pengguna{})
	db.Exec("TRUNCATE TABLE pengguna CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	jwtUtil := utils.NewJWTUtil("test-secret-key", 24)
	service := NewAuthService(db, jwtUtil)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Email:        "test@example.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = service.Login("testuser", "password123")
	}
}
