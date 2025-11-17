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

// setupKoperasiTestDB creates a test database for koperasi service
func setupKoperasiTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.Koperasi{}, &models.Anggota{}, &models.Pengguna{}, &models.Produk{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean up existing data
	db.Exec("TRUNCATE TABLE koperasi CASCADE")
	db.Exec("TRUNCATE TABLE anggota CASCADE")
	db.Exec("TRUNCATE TABLE pengguna CASCADE")
	db.Exec("TRUNCATE TABLE produk CASCADE")

	return db
}

// TestBuatKoperasi tests cooperative creation
func TestBuatKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	t.Run("successful creation", func(t *testing.T) {
		req := &BuatKoperasiRequest{
			NamaKoperasi:   "Koperasi Test",
			Alamat:         "Jl. Test No. 123",
			NoTelepon:      "081234567890",
			Email:          "test@koperasi.com",
			TahunBukuMulai: 2025,
		}

		result, err := service.BuatKoperasi(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Koperasi Test", result.NamaKoperasi)
		assert.Equal(t, "Jl. Test No. 123", result.Alamat)
		assert.Equal(t, "081234567890", result.NoTelepon)
		assert.Equal(t, "test@koperasi.com", result.Email)
		assert.Equal(t, 2025, result.TahunBukuMulai)
		assert.NotEqual(t, uuid.Nil, result.ID)
	})

	t.Run("minimal fields", func(t *testing.T) {
		req := &BuatKoperasiRequest{
			NamaKoperasi: "Koperasi Minimal",
		}

		result, err := service.BuatKoperasi(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Koperasi Minimal", result.NamaKoperasi)
	})
}

// TestDapatkanKoperasi tests getting cooperative by ID
func TestDapatkanKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	// Create test cooperative
	koperasi := &models.Koperasi{
		ID:             uuid.New(),
		NamaKoperasi:   "Test Koperasi",
		Alamat:         "Jl. Test",
		NoTelepon:      "081234567890",
		Email:          "test@koperasi.com",
		TahunBukuMulai: 2025,
	}
	db.Create(koperasi)

	t.Run("existing cooperative", func(t *testing.T) {
		result, err := service.DapatkanKoperasi(koperasi.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, koperasi.ID, result.ID)
		assert.Equal(t, "Test Koperasi", result.NamaKoperasi)
	})

	t.Run("non-existing cooperative", func(t *testing.T) {
		result, err := service.DapatkanKoperasi(uuid.New())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "koperasi tidak ditemukan")
	})
}

// TestPerbaruiKoperasi tests cooperative updates
func TestPerbaruiKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	// Create test cooperative
	koperasi := &models.Koperasi{
		ID:             uuid.New(),
		NamaKoperasi:   "Original Name",
		Alamat:         "Original Address",
		NoTelepon:      "081234567890",
		Email:          "original@koperasi.com",
		TahunBukuMulai: 2024,
	}
	db.Create(koperasi)

	t.Run("update all fields", func(t *testing.T) {
		req := &PerbaruiKoperasiRequest{
			NamaKoperasi:   "Updated Name",
			Alamat:         "Updated Address",
			NoTelepon:      "082345678901",
			Email:          "updated@koperasi.com",
			LogoURL:        "https://example.com/logo.png",
			TahunBukuMulai: 2025,
		}

		result, err := service.PerbaruiKoperasi(koperasi.ID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Name", result.NamaKoperasi)
		assert.Equal(t, "Updated Address", result.Alamat)
		assert.Equal(t, "082345678901", result.NoTelepon)
		assert.Equal(t, "updated@koperasi.com", result.Email)
		assert.Equal(t, "https://example.com/logo.png", result.LogoURL)
		assert.Equal(t, 2025, result.TahunBukuMulai)
	})

	t.Run("partial update", func(t *testing.T) {
		req := &PerbaruiKoperasiRequest{
			NamaKoperasi: "Partial Update",
		}

		result, err := service.PerbaruiKoperasi(koperasi.ID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Partial Update", result.NamaKoperasi)
		// Other fields should remain unchanged
		assert.Equal(t, "updated@koperasi.com", result.Email)
	})

	t.Run("non-existing cooperative", func(t *testing.T) {
		req := &PerbaruiKoperasiRequest{
			NamaKoperasi: "Should Fail",
		}

		result, err := service.PerbaruiKoperasi(uuid.New(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "koperasi tidak ditemukan")
	})
}

// TestDapatkanSemuaKoperasi tests getting all cooperatives
func TestDapatkanSemuaKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	// Create multiple test cooperatives
	cooperatives := []models.Koperasi{
		{
			ID:           uuid.New(),
			NamaKoperasi: "Koperasi A",
			NoTelepon:    "081234567890",
			Email:        "a@koperasi.com",
		},
		{
			ID:           uuid.New(),
			NamaKoperasi: "Koperasi B",
			NoTelepon:    "081234567891",
			Email:        "b@koperasi.com",
		},
		{
			ID:           uuid.New(),
			NamaKoperasi: "Koperasi C",
			NoTelepon:    "081234567892",
			Email:        "c@koperasi.com",
		},
	}

	for _, koperasi := range cooperatives {
		db.Create(&koperasi)
	}

	t.Run("get all cooperatives", func(t *testing.T) {
		results, err := service.DapatkanSemuaKoperasi()

		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.GreaterOrEqual(t, len(results), 3)
	})
}

// TestGetSemuaKoperasi tests the English wrapper with pagination
func TestGetSemuaKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	// Create test cooperatives
	for i := 0; i < 5; i++ {
		koperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Koperasi Test",
			NoTelepon:    "081234567890",
			Email:        "test@koperasi.com",
		}
		db.Create(koperasi)
	}

	t.Run("first page", func(t *testing.T) {
		results, total, err := service.GetSemuaKoperasi(1, 2)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, results, 2)
	})

	t.Run("second page", func(t *testing.T) {
		results, total, err := service.GetSemuaKoperasi(2, 2)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, results, 2)
	})

	t.Run("last page", func(t *testing.T) {
		results, total, err := service.GetSemuaKoperasi(3, 2)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, results, 1)
	})

	t.Run("empty page", func(t *testing.T) {
		results, total, err := service.GetSemuaKoperasi(10, 2)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, results, 0)
	})
}

// TestHapusKoperasi tests cooperative deletion
func TestHapusKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	t.Run("successful deletion", func(t *testing.T) {
		koperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "To Delete",
			NoTelepon:    "081234567890",
			Email:        "delete@koperasi.com",
		}
		db.Create(koperasi)

		err := service.HapusKoperasi(koperasi.ID)

		assert.NoError(t, err)

		// Verify soft delete
		var count int64
		db.Model(&models.Koperasi{}).Where("id = ?", koperasi.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("non-existing cooperative", func(t *testing.T) {
		err := service.HapusKoperasi(uuid.New())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "koperasi tidak ditemukan")
	})
}

// TestDapatkanStatistikKoperasi tests cooperative statistics
func TestDapatkanStatistikKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	// Create test cooperative
	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		NoTelepon:    "081234567890",
		Email:        "test@koperasi.com",
	}
	db.Create(koperasi)

	// Create test data
	// Create members
	for i := 0; i < 3; i++ {
		anggota := &models.Anggota{
			IDKoperasi:   koperasi.ID,
			NamaLengkap:  "Member",
			NomorAnggota: "A001",
			JenisKelamin: "L",
			Status:       models.StatusAktif,
		}
		db.Create(anggota)
	}

	// Create users
	for i := 0; i < 2; i++ {
		pengguna := &models.Pengguna{
			IDKoperasi:   koperasi.ID,
			NamaLengkap:  "User",
			NamaPengguna: "user",
			Email:        "user@example.com",
			Peran:        models.PeranAdmin,
			StatusAktif:  true,
		}
		db.Create(pengguna)
	}

	// Create products
	for i := 0; i < 5; i++ {
		produk := &models.Produk{
			IDKoperasi:  koperasi.ID,
			NamaProduk:  "Product",
			KodeProduk:  "P001",
			Harga:       10000,
			HargaBeli:   8000,
			Stok:        100,
			StatusAktif: true,
		}
		db.Create(produk)
	}

	t.Run("get statistics", func(t *testing.T) {
		stats, err := service.DapatkanStatistikKoperasi(koperasi.ID)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(3), stats["jumlahAnggota"])
		assert.Equal(t, int64(2), stats["jumlahPengguna"])
		assert.Equal(t, int64(5), stats["jumlahProduk"])
	})

	t.Run("non-existing cooperative", func(t *testing.T) {
		stats, err := service.DapatkanStatistikKoperasi(uuid.New())

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "koperasi tidak ditemukan")
	})
}

// TestGetStatistikKoperasi tests the English wrapper
func TestGetStatistikKoperasi(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		NoTelepon:    "081234567890",
		Email:        "test@koperasi.com",
	}
	db.Create(koperasi)

	t.Run("get statistics via English wrapper", func(t *testing.T) {
		stats, err := service.GetStatistikKoperasi(koperasi.ID)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Contains(t, stats, "jumlahAnggota")
		assert.Contains(t, stats, "jumlahPengguna")
		assert.Contains(t, stats, "jumlahProduk")
	})
}

// TestGetKoperasiByID tests the English wrapper for getting cooperative
func TestGetKoperasiByID(t *testing.T) {
	db := setupKoperasiTestDB(t)
	if db == nil {
		return
	}

	service := NewKoperasiService(db)

	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		NoTelepon:    "081234567890",
		Email:        "test@koperasi.com",
	}
	db.Create(koperasi)

	t.Run("get by ID via English wrapper", func(t *testing.T) {
		result, err := service.GetKoperasiByID(koperasi.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, koperasi.ID, result.ID)
		assert.Equal(t, "Test Koperasi", result.NamaKoperasi)
	})

	t.Run("non-existing cooperative", func(t *testing.T) {
		result, err := service.GetKoperasiByID(uuid.New())

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// BenchmarkBuatKoperasi benchmarks cooperative creation
func BenchmarkBuatKoperasi(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{})
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	service := NewKoperasiService(db)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &BuatKoperasiRequest{
			NamaKoperasi:   "Benchmark Koperasi",
			Alamat:         "Jl. Benchmark",
			NoTelepon:      "081234567890",
			Email:          uuid.New().String() + "@benchmark.com",
			TahunBukuMulai: 2025,
		}
		service.BuatKoperasi(req)
	}
}

// BenchmarkDapatkanStatistikKoperasi benchmarks statistics generation
func BenchmarkDapatkanStatistikKoperasi(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{}, &models.Anggota{}, &models.Pengguna{}, &models.Produk{})
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	service := NewKoperasiService(db)

	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Benchmark Koperasi",
		NoTelepon:    "081234567890",
		Email:        "bench@koperasi.com",
	}
	db.Create(koperasi)

	// Create test data
	for i := 0; i < 100; i++ {
		anggota := &models.Anggota{
			IDKoperasi:   koperasi.ID,
			NamaLengkap:  "Member",
			NomorAnggota: "A001",
			JenisKelamin: "L",
			Status:       models.StatusAktif,
		}
		db.Create(anggota)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		service.DapatkanStatistikKoperasi(koperasi.ID)
	}
}
