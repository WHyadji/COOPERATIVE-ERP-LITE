package services

import (
	"cooperative-erp-lite/internal/models"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test helper untuk akun service
func setupAkunTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Akun{},
		&models.BarisTransaksi{},
		&models.Transaksi{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean up existing data
	db.Exec("TRUNCATE TABLE baris_transaksi CASCADE")
	db.Exec("TRUNCATE TABLE transaksi CASCADE")
	db.Exec("TRUNCATE TABLE akun CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	return db
}

// TestBuatAkun_Success tests successful account creation
func TestBuatAkun_Success(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)

	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	req := &BuatAkunRequest{
		KodeAkun: "1101",
		NamaAkun: "Kas",
		TipeAkun: models.AkunAset,
		Deskripsi: "Kas di tangan",
	}

	result, err := service.BuatAkun(koperasi.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1101", result.KodeAkun)
	assert.Equal(t, "Kas", result.NamaAkun)
	assert.Equal(t, models.AkunAset, result.TipeAkun)
	assert.Equal(t, "debit", result.NormalSaldo) // Auto-set berdasarkan tipe
	assert.True(t, result.StatusAktif)
}

// TestBuatAkun_ValidationErrors tests account creation validation
func TestBuatAkun_ValidationErrors(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	tests := []struct {
		name    string
		req     *BuatAkunRequest
		wantErr bool
	}{
		{
			name: "invalid kode akun (contains letters)",
			req: &BuatAkunRequest{
				KodeAkun: "ABC1",
				NamaAkun: "Test Account",
				TipeAkun: models.AkunAset,
			},
			wantErr: true,
		},
		{
			name: "too short nama akun",
			req: &BuatAkunRequest{
				KodeAkun: "1101",
				NamaAkun: "AB",
				TipeAkun: models.AkunAset,
			},
			wantErr: true,
		},
		{
			name: "valid account",
			req: &BuatAkunRequest{
				KodeAkun: "1101",
				NamaAkun: "Valid Account",
				TipeAkun: models.AkunAset,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.BuatAkun(koperasi.ID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBuatAkun_DuplicateCode tests duplicate account code validation
func TestBuatAkun_DuplicateCode(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create first account
	req1 := &BuatAkunRequest{
		KodeAkun: "1101",
		NamaAkun: "Kas",
		TipeAkun: models.AkunAset,
	}
	_, err := service.BuatAkun(koperasi.ID, req1)
	assert.NoError(t, err)

	// Try to create account with same code
	req2 := &BuatAkunRequest{
		KodeAkun: "1101",
		NamaAkun: "Bank",
		TipeAkun: models.AkunAset,
	}
	_, err = service.BuatAkun(koperasi.ID, req2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sudah digunakan")
}

// TestBuatAkun_WithParent tests account creation with parent
func TestBuatAkun_WithParent(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create parent account
	parentReq := &BuatAkunRequest{
		KodeAkun: "1100",
		NamaAkun: "Aset Lancar",
		TipeAkun: models.AkunAset,
	}
	parent, err := service.BuatAkun(koperasi.ID, parentReq)
	assert.NoError(t, err)

	// Create child account
	childReq := &BuatAkunRequest{
		KodeAkun: "1101",
		NamaAkun: "Kas",
		TipeAkun: models.AkunAset,
		IDInduk:  &parent.ID,
	}
	child, err := service.BuatAkun(koperasi.ID, childReq)
	assert.NoError(t, err)
	assert.NotNil(t, child.IDInduk)
	assert.Equal(t, parent.ID, *child.IDInduk)
}

// TestBuatAkun_InvalidParent tests account creation with non-existent parent
func TestBuatAkun_InvalidParent(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	invalidParentID := uuid.New()
	req := &BuatAkunRequest{
		KodeAkun: "1101",
		NamaAkun: "Kas",
		TipeAkun: models.AkunAset,
		IDInduk:  &invalidParentID,
	}

	_, err := service.BuatAkun(koperasi.ID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "induk tidak ditemukan")
}

// TestBuatAkun_AutoSetNormalSaldo tests automatic normal saldo setting
func TestBuatAkun_AutoSetNormalSaldo(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	tests := []struct {
		name            string
		tipeAkun        models.TipeAkun
		expectedNormal  string
	}{
		{"aset should be debit", models.AkunAset, "debit"},
		{"beban should be debit", models.AkunBeban, "debit"},
		{"kewajiban should be kredit", models.AkunKewajiban, "kredit"},
		{"modal should be kredit", models.AkunModal, "kredit"},
		{"pendapatan should be kredit", models.AkunPendapatan, "kredit"},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &BuatAkunRequest{
				KodeAkun: fmt.Sprintf("%d", 1000+i),
				NamaAkun: fmt.Sprintf("Test Account %d", i),
				TipeAkun: tt.tipeAkun,
			}

			result, err := service.BuatAkun(koperasi.ID, req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedNormal, result.NormalSaldo)
		})
	}
}

// TestDapatkanSemuaAkun tests listing accounts with filters
func TestDapatkanSemuaAkun(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create test accounts
	accounts := []BuatAkunRequest{
		{KodeAkun: "1101", NamaAkun: "Kas", TipeAkun: models.AkunAset},
		{KodeAkun: "1102", NamaAkun: "Bank", TipeAkun: models.AkunAset},
		{KodeAkun: "2101", NamaAkun: "Hutang", TipeAkun: models.AkunKewajiban},
		{KodeAkun: "4101", NamaAkun: "Pendapatan", TipeAkun: models.AkunPendapatan},
	}

	for _, acc := range accounts {
		service.BuatAkun(koperasi.ID, &acc)
	}

	t.Run("get all accounts", func(t *testing.T) {
		results, err := service.DapatkanSemuaAkun(koperasi.ID, "", nil)
		assert.NoError(t, err)
		assert.Len(t, results, 4)
	})

	t.Run("filter by tipe aset", func(t *testing.T) {
		results, err := service.DapatkanSemuaAkun(koperasi.ID, string(models.AkunAset), nil)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("filter by status aktif", func(t *testing.T) {
		aktif := true
		results, err := service.DapatkanSemuaAkun(koperasi.ID, "", &aktif)
		assert.NoError(t, err)
		assert.Len(t, results, 4) // All created accounts are active
	})
}

// TestPerbaruiAkun tests account updates
func TestPerbaruiAkun(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	createReq := &BuatAkunRequest{
		KodeAkun: "1101",
		NamaAkun: "Original Name",
		TipeAkun: models.AkunAset,
	}
	akun, _ := service.BuatAkun(koperasi.ID, createReq)

	t.Run("update name", func(t *testing.T) {
		updateReq := &PerbaruiAkunRequest{
			NamaAkun: "Updated Name",
		}
		result, err := service.PerbaruiAkun(koperasi.ID, akun.ID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", result.NamaAkun)
	})

	t.Run("deactivate account", func(t *testing.T) {
		statusFalse := false
		updateReq := &PerbaruiAkunRequest{
			StatusAktif: &statusFalse,
		}
		result, err := service.PerbaruiAkun(koperasi.ID, akun.ID, updateReq)
		assert.NoError(t, err)
		assert.False(t, result.StatusAktif)
	})
}

// TestHapusAkun tests account deletion with validation
func TestHapusAkun(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	t.Run("delete account without transactions", func(t *testing.T) {
		createReq := &BuatAkunRequest{
			KodeAkun: "1101",
			NamaAkun: "Kas",
			TipeAkun: models.AkunAset,
		}
		akun, _ := service.BuatAkun(koperasi.ID, createReq)

		err := service.HapusAkun(koperasi.ID, akun.ID)
		assert.NoError(t, err)

		// Verify soft delete
		_, err = service.DapatkanAkun(akun.ID)
		assert.Error(t, err)
	})

	t.Run("delete account with sub-accounts (should fail)", func(t *testing.T) {
		// Create parent
		parentReq := &BuatAkunRequest{
			KodeAkun: "1100",
			NamaAkun: "Aset Lancar",
			TipeAkun: models.AkunAset,
		}
		parent, _ := service.BuatAkun(koperasi.ID, parentReq)

		// Create child
		childReq := &BuatAkunRequest{
			KodeAkun: "1101",
			NamaAkun: "Kas",
			TipeAkun: models.AkunAset,
			IDInduk:  &parent.ID,
		}
		service.BuatAkun(koperasi.ID, childReq)

		err := service.HapusAkun(koperasi.ID, parent.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sub-akun")
	})
}

// TestInisialisasiCOADefault tests default COA initialization
func TestInisialisasiCOADefault(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	err := service.InisialisasiCOADefault(koperasi.ID)
	assert.NoError(t, err)

	// Verify COA created
	akuns, err := service.DapatkanSemuaAkun(koperasi.ID, "", nil)
	assert.NoError(t, err)
	assert.Greater(t, len(akuns), 20) // Should have many default accounts

	// Verify specific accounts exist
	t.Run("verify ASET account exists", func(t *testing.T) {
		akun, err := service.DapatkanAkunByKode(koperasi.ID, "1000")
		assert.NoError(t, err)
		assert.Equal(t, "ASET", akun.NamaAkun)
		assert.Equal(t, models.AkunAset, akun.TipeAkun)
		assert.Equal(t, "debit", akun.NormalSaldo)
	})

	t.Run("verify Kas account exists", func(t *testing.T) {
		akun, err := service.DapatkanAkunByKode(koperasi.ID, "1101")
		assert.NoError(t, err)
		assert.Equal(t, "Kas", akun.NamaAkun)
	})

	t.Run("verify Simpanan Pokok account exists", func(t *testing.T) {
		akun, err := service.DapatkanAkunByKode(koperasi.ID, "3101")
		assert.NoError(t, err)
		assert.Equal(t, "Simpanan Pokok", akun.NamaAkun)
		assert.Equal(t, models.AkunModal, akun.TipeAkun)
		assert.Equal(t, "kredit", akun.NormalSaldo)
	})

	t.Run("verify Penjualan account exists", func(t *testing.T) {
		akun, err := service.DapatkanAkunByKode(koperasi.ID, "4101")
		assert.NoError(t, err)
		assert.Equal(t, "Penjualan", akun.NamaAkun)
		assert.Equal(t, models.AkunPendapatan, akun.TipeAkun)
	})

	t.Run("verify all account types exist", func(t *testing.T) {
		aset, _ := service.DapatkanSemuaAkun(koperasi.ID, string(models.AkunAset), nil)
		kewajiban, _ := service.DapatkanSemuaAkun(koperasi.ID, string(models.AkunKewajiban), nil)
		modal, _ := service.DapatkanSemuaAkun(koperasi.ID, string(models.AkunModal), nil)
		pendapatan, _ := service.DapatkanSemuaAkun(koperasi.ID, string(models.AkunPendapatan), nil)
		beban, _ := service.DapatkanSemuaAkun(koperasi.ID, string(models.AkunBeban), nil)

		assert.Greater(t, len(aset), 0, "Should have ASET accounts")
		assert.Greater(t, len(kewajiban), 0, "Should have KEWAJIBAN accounts")
		assert.Greater(t, len(modal), 0, "Should have MODAL accounts")
		assert.Greater(t, len(pendapatan), 0, "Should have PENDAPATAN accounts")
		assert.Greater(t, len(beban), 0, "Should have BEBAN accounts")
	})
}

// TestInisialisasiCOADefault_Idempotent tests COA initialization is transaction-safe
func TestInisialisasiCOADefault_Idempotent(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// First initialization should succeed
	err := service.InisialisasiCOADefault(koperasi.ID)
	assert.NoError(t, err)

	count1, _ := service.DapatkanSemuaAkun(koperasi.ID, "", nil)

	// Second initialization should fail (duplicate codes)
	err = service.InisialisasiCOADefault(koperasi.ID)
	assert.Error(t, err) // Should error due to duplicates

	count2, _ := service.DapatkanSemuaAkun(koperasi.ID, "", nil)

	// Count should remain the same (transaction rolled back)
	assert.Equal(t, len(count1), len(count2))
}

// TestDapatkanHierarkiAkun tests hierarchical COA structure
func TestDapatkanHierarkiAkun(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Initialize default COA
	err := service.InisialisasiCOADefault(koperasi.ID)
	assert.NoError(t, err)

	hierarchy, err := service.DapatkanHierarkiAkun(koperasi.ID)
	assert.NoError(t, err)
	assert.Greater(t, len(hierarchy), 0)
}

// TestMultiTenant_AccountIsolation tests multi-tenant account isolation
func TestMultiTenant_AccountIsolation(t *testing.T) {
	db := setupAkunTestDB(t)
	if db == nil {
		return
	}

	service := NewAkunService(db)

	// Create two cooperatives
	koperasi1 := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Koperasi 1", Email: "kop1@test.com", NoTelepon: "081234567890"}
	koperasi2 := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Koperasi 2", Email: "kop2@test.com", NoTelepon: "081234567891"}
	db.Create(koperasi1)
	db.Create(koperasi2)

	// Create accounts in both cooperatives with same code
	req1 := &BuatAkunRequest{KodeAkun: "1101", NamaAkun: "Kas Koperasi 1", TipeAkun: models.AkunAset}
	req2 := &BuatAkunRequest{KodeAkun: "1101", NamaAkun: "Kas Koperasi 2", TipeAkun: models.AkunAset}

	akun1, err := service.BuatAkun(koperasi1.ID, req1)
	assert.NoError(t, err)

	akun2, err := service.BuatAkun(koperasi2.ID, req2)
	assert.NoError(t, err)

	// Verify accounts are different
	assert.NotEqual(t, akun1.ID, akun2.ID)
	assert.Equal(t, "Kas Koperasi 1", akun1.NamaAkun)
	assert.Equal(t, "Kas Koperasi 2", akun2.NamaAkun)

	// Verify cross-tenant access denied
	t.Run("cross-tenant update should fail", func(t *testing.T) {
		updateReq := &PerbaruiAkunRequest{NamaAkun: "Hacker Name"}
		_, err := service.PerbaruiAkun(koperasi2.ID, akun1.ID, updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak memiliki akses")
	})

	t.Run("cross-tenant delete should fail", func(t *testing.T) {
		err := service.HapusAkun(koperasi2.ID, akun1.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak memiliki akses")
	})
}

// BenchmarkBuatAkun benchmarks account creation
func BenchmarkBuatAkun(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{}, &models.Akun{})
	db.Exec("TRUNCATE TABLE akun CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	service := NewAkunService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &BuatAkunRequest{
			KodeAkun: fmt.Sprintf("%04d", i),
			NamaAkun: fmt.Sprintf("Account %d", i),
			TipeAkun: models.AkunAset,
		}
		_, _ = service.BuatAkun(koperasi.ID, req)
	}
}

// BenchmarkInisialisasiCOADefault benchmarks default COA initialization
func BenchmarkInisialisasiCOADefault(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{}, &models.Akun{})

	service := NewAkunService(db)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		db.Exec("TRUNCATE TABLE akun CASCADE")
		db.Exec("TRUNCATE TABLE koperasi CASCADE")
		koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
		db.Create(koperasi)
		b.StartTimer()

		_ = service.InisialisasiCOADefault(koperasi.ID)
	}
}
