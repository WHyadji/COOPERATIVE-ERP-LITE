package services

import (
	"cooperative-erp-lite/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupPenggunaTestDB creates a test database for pengguna service
func setupPenggunaTestDB(t *testing.T) *gorm.DB {
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

// TestBuatPengguna tests user creation
func TestBuatPengguna(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

	// Create test cooperative
	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	t.Run("successful user creation", func(t *testing.T) {
		req := &BuatPenggunaRequest{
			NamaLengkap:  "Test User",
			NamaPengguna: "testuser",
			Email:        "test@example.com",
			KataSandi:    "password123",
			Peran:        models.PeranAdmin,
		}

		result, err := service.BuatPengguna(koperasi.ID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test User", result.NamaLengkap)
		assert.Equal(t, "testuser", result.NamaPengguna)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, models.PeranAdmin, result.Peran)
		assert.True(t, result.StatusAktif)

		// Password is stored hashed (private field, can't verify directly)
	})

	t.Run("duplicate username", func(t *testing.T) {
		// First user
		req1 := &BuatPenggunaRequest{
			NamaLengkap:  "User One",
			NamaPengguna: "duplicate",
			Email:        "user1@example.com",
			KataSandi:    "password123",
			Peran:        models.PeranAdmin,
		}
		service.BuatPengguna(koperasi.ID, req1)

		// Try to create duplicate
		req2 := &BuatPenggunaRequest{
			NamaLengkap:  "User Two",
			NamaPengguna: "duplicate",
			Email:        "user2@example.com",
			KataSandi:    "password123",
			Peran:        models.PeranKasir,
		}
		result, err := service.BuatPengguna(koperasi.ID, req2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "nama pengguna sudah digunakan")
	})

	t.Run("validation errors", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *BuatPenggunaRequest
			wantErr string
		}{
			{
				name: "short full name",
				req: &BuatPenggunaRequest{
					NamaLengkap:  "AB",
					NamaPengguna: "testuser2",
					Email:        "test2@example.com",
					KataSandi:    "password123",
					Peran:        models.PeranAdmin,
				},
				wantErr: "nama lengkap",
			},
			{
				name: "short username",
				req: &BuatPenggunaRequest{
					NamaLengkap:  "Test User",
					NamaPengguna: "ab",
					Email:        "test3@example.com",
					KataSandi:    "password123",
					Peran:        models.PeranAdmin,
				},
				wantErr: "nama pengguna",
			},
			{
				name: "invalid email",
				req: &BuatPenggunaRequest{
					NamaLengkap:  "Test User",
					NamaPengguna: "testuser3",
					Email:        "invalid-email",
					KataSandi:    "password123",
					Peran:        models.PeranAdmin,
				},
				wantErr: "email",
			},
			{
				name: "short password",
				req: &BuatPenggunaRequest{
					NamaLengkap:  "Test User",
					NamaPengguna: "testuser4",
					Email:        "test4@example.com",
					KataSandi:    "12345",
					Peran:        models.PeranAdmin,
				},
				wantErr: "kata sandi",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := service.BuatPengguna(koperasi.ID, tt.req)
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})
}

// TestDapatkanSemuaPengguna tests user listing with filters
func TestDapatkanSemuaPengguna(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

	// Create test cooperative
	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	// Create test users with different roles and status
	users := []struct {
		namaLengkap  string
		namaPengguna string
		email        string
		peran        models.PeranPengguna
		statusAktif  bool
	}{
		{"Admin User", "admin1", "admin@test.com", models.PeranAdmin, true},
		{"Treasurer User", "bendahara1", "bendahara@test.com", models.PeranBendahara, true},
		{"Cashier User", "kasir1", "kasir@test.com", models.PeranKasir, true},
		{"Inactive User", "inactive1", "inactive@test.com", models.PeranKasir, false},
	}

	for _, u := range users {
		pengguna := &models.Pengguna{
			IDKoperasi:   koperasi.ID,
			NamaLengkap:  u.namaLengkap,
			NamaPengguna: u.namaPengguna,
			Email:        u.email,
			Peran:        u.peran,
			StatusAktif:  u.statusAktif,
		}
		pengguna.SetKataSandi("password123")
		db.Create(pengguna)
	}

	t.Run("get all users", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaPengguna(koperasi.ID, "", nil, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, results, 4)
	})

	t.Run("filter by role", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaPengguna(koperasi.ID, string(models.PeranKasir), nil, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)
	})

	t.Run("filter by active status", func(t *testing.T) {
		statusAktif := true
		results, total, err := service.DapatkanSemuaPengguna(koperasi.ID, "", &statusAktif, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 3)
	})

	t.Run("pagination", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaPengguna(koperasi.ID, "", nil, 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, results, 2)
	})
}

// TestDapatkanPengguna tests getting single user
func TestDapatkanPengguna(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

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
	db.Create(pengguna)

	t.Run("existing user", func(t *testing.T) {
		result, err := service.DapatkanPengguna(pengguna.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, pengguna.ID, result.ID)
		assert.Equal(t, "Test User", result.NamaLengkap)
	})

	t.Run("non-existing user", func(t *testing.T) {
		result, err := service.DapatkanPengguna(uuid.New())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "pengguna tidak ditemukan")
	})
}

// TestPerbaruiPengguna tests user updates
func TestPerbaruiPengguna(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Original Name",
		NamaPengguna: "originaluser",
		Email:        "original@example.com",
		Peran:        models.PeranKasir,
		StatusAktif:  true,
	}
	db.Create(pengguna)

	t.Run("successful update", func(t *testing.T) {
		statusAktif := false
		req := &PerbaruiPenggunaRequest{
			NamaLengkap: "Updated Name",
			Email:       "updated@example.com",
			Peran:       models.PeranBendahara,
			StatusAktif: &statusAktif,
		}

		result, err := service.PerbaruiPengguna(koperasi.ID, pengguna.ID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Name", result.NamaLengkap)
		assert.Equal(t, "updated@example.com", result.Email)
		assert.Equal(t, models.PeranBendahara, result.Peran)
		assert.False(t, result.StatusAktif)
	})

	t.Run("multi-tenant validation", func(t *testing.T) {
		// Try to update user from different cooperative
		otherKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Other Koperasi",
			Email:        "other@koperasi.com",
			NoTelepon:    "081234567891",
		}
		db.Create(otherKoperasi)

		req := &PerbaruiPenggunaRequest{
			NamaLengkap: "Hacked Name",
		}

		result, err := service.PerbaruiPengguna(otherKoperasi.ID, pengguna.ID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan atau tidak memiliki akses")
	})

	t.Run("validation errors", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *PerbaruiPenggunaRequest
			wantErr string
		}{
			{
				name: "invalid email",
				req: &PerbaruiPenggunaRequest{
					Email: "invalid-email",
				},
				wantErr: "email",
			},
			{
				name: "short name",
				req: &PerbaruiPenggunaRequest{
					NamaLengkap: "AB",
				},
				wantErr: "nama lengkap",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := service.PerbaruiPengguna(koperasi.ID, pengguna.ID, tt.req)
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})
}

// TestHapusPengguna tests user deletion
func TestHapusPengguna(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	t.Run("successful deletion", func(t *testing.T) {
		pengguna := &models.Pengguna{
			IDKoperasi:   koperasi.ID,
			NamaLengkap:  "To Delete",
			NamaPengguna: "todelete",
			Email:        "delete@example.com",
			Peran:        models.PeranKasir,
			StatusAktif:  true,
		}
		db.Create(pengguna)

		err := service.HapusPengguna(koperasi.ID, pengguna.ID)

		assert.NoError(t, err)

		// Verify soft delete
		var count int64
		db.Model(&models.Pengguna{}).Where("id = ?", pengguna.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("multi-tenant validation", func(t *testing.T) {
		pengguna := &models.Pengguna{
			IDKoperasi:   koperasi.ID,
			NamaLengkap:  "User",
			NamaPengguna: "user1",
			Email:        "user1@example.com",
			Peran:        models.PeranKasir,
			StatusAktif:  true,
		}
		db.Create(pengguna)

		// Try to delete from different cooperative
		otherKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Other Koperasi",
			Email:        "other@koperasi.com",
			NoTelepon:    "081234567891",
		}
		db.Create(otherKoperasi)

		err := service.HapusPengguna(otherKoperasi.ID, pengguna.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan atau tidak memiliki akses")
	})
}

// TestUbahKataSandiPengguna tests password change by admin
func TestUbahKataSandiPengguna(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

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
		Peran:        models.PeranKasir,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("oldpassword")
	db.Create(pengguna)

	t.Run("successful password change", func(t *testing.T) {
		err := service.UbahKataSandiPengguna(koperasi.ID, pengguna.ID, "newpassword123")

		assert.NoError(t, err)

		// Verify password was changed
		var updated models.Pengguna
		db.Where("id = ?", pengguna.ID).First(&updated)
		assert.True(t, updated.CekKataSandi("newpassword123"))
		assert.False(t, updated.CekKataSandi("oldpassword"))
	})

	t.Run("password too short", func(t *testing.T) {
		err := service.UbahKataSandiPengguna(koperasi.ID, pengguna.ID, "12345")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kata sandi minimal 6 karakter")
	})

	t.Run("multi-tenant validation", func(t *testing.T) {
		otherKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Other Koperasi",
			Email:        "other@koperasi.com",
			NoTelepon:    "081234567891",
		}
		db.Create(otherKoperasi)

		err := service.UbahKataSandiPengguna(otherKoperasi.ID, pengguna.ID, "newpassword123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan atau tidak memiliki akses")
	})
}

// TestResetKataSandi tests password reset to default
func TestResetKataSandi(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

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
		Peran:        models.PeranKasir,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("oldpassword")
	db.Create(pengguna)

	t.Run("successful reset", func(t *testing.T) {
		defaultPassword, err := service.ResetKataSandi(koperasi.ID, pengguna.ID)

		assert.NoError(t, err)
		assert.Equal(t, "testuser123", defaultPassword)

		// Verify password was reset
		var updated models.Pengguna
		db.Where("id = ?", pengguna.ID).First(&updated)
		assert.True(t, updated.CekKataSandi("testuser123"))
	})

	t.Run("multi-tenant validation", func(t *testing.T) {
		otherKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Other Koperasi",
			Email:        "other@koperasi.com",
			NoTelepon:    "081234567891",
		}
		db.Create(otherKoperasi)

		_, err := service.ResetKataSandi(otherKoperasi.ID, pengguna.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan atau tidak memiliki akses")
	})
}

// TestDapatkanPenggunaByUsername tests getting user by username
func TestDapatkanPenggunaByUsername(t *testing.T) {
	db := setupPenggunaTestDB(t)
	if db == nil {
		return
	}

	service := NewPenggunaService(db)

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
	db.Create(pengguna)

	t.Run("existing username", func(t *testing.T) {
		result, err := service.DapatkanPenggunaByUsername(koperasi.ID, "testuser")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "testuser", result.NamaPengguna)
		assert.Equal(t, "Test User", result.NamaLengkap)
	})

	t.Run("non-existing username", func(t *testing.T) {
		result, err := service.DapatkanPenggunaByUsername(koperasi.ID, "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "pengguna tidak ditemukan")
	})

	t.Run("multi-tenant isolation", func(t *testing.T) {
		otherKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Other Koperasi",
			Email:        "other@koperasi.com",
			NoTelepon:    "081234567891",
		}
		db.Create(otherKoperasi)

		result, err := service.DapatkanPenggunaByUsername(otherKoperasi.ID, "testuser")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// BenchmarkBuatPengguna benchmarks user creation
func BenchmarkBuatPengguna(b *testing.B) {
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

	service := NewPenggunaService(db)

	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Benchmark Koperasi",
		Email:        "bench@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &BuatPenggunaRequest{
			NamaLengkap:  "Benchmark User",
			NamaPengguna: uuid.New().String()[:8],
			Email:        uuid.New().String() + "@example.com",
			KataSandi:    "password123",
			Peran:        models.PeranKasir,
		}
		service.BuatPengguna(koperasi.ID, req)
	}
}
