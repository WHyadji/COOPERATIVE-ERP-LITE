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

// Test helper untuk produk service
func setupProdukTestDB(t *testing.T) *gorm.DB {
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
		&models.Produk{},
		&models.ItemPenjualan{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean up existing data
	db.Exec("TRUNCATE TABLE item_penjualan CASCADE")
	db.Exec("TRUNCATE TABLE produk CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	return db
}

// TestBuatProduk_Success tests successful product creation
func TestBuatProduk_Success(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)

	// Create test cooperative
	koperasi := &models.Koperasi{
		ID:          uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	req := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Kategori:   "Elektronik",
		Deskripsi:  "Test product description",
		Harga:      50000,
		HargaBeli:  40000,
		Stok:       100,
		StokMinimum: 10,
		Satuan:     "pcs",
		Barcode:    "1234567890",
	}

	result, err := service.BuatProduk(koperasi.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, "PRD001", result.KodeProduk)
	assert.Equal(t, "Test Product", result.NamaProduk)
	assert.Equal(t, 50000.0, result.Harga)
	assert.Equal(t, 100, result.Stok)
	assert.True(t, result.StatusAktif)
}

// TestBuatProduk_ValidationErrors tests product creation validation
func TestBuatProduk_ValidationErrors(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	tests := []struct {
		name    string
		req     *BuatProdukRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "empty kode produk",
			req: &BuatProdukRequest{
				KodeProduk: "",
				NamaProduk: "Test Product",
				Harga:      50000,
			},
			wantErr: true,
		},
		{
			name: "too short nama produk",
			req: &BuatProdukRequest{
				KodeProduk: "PRD001",
				NamaProduk: "AB",
				Harga:      50000,
			},
			wantErr: true,
		},
		{
			name: "negative harga",
			req: &BuatProdukRequest{
				KodeProduk: "PRD001",
				NamaProduk: "Test Product",
				Harga:      -1000,
			},
			wantErr: true,
		},
		{
			name: "harga too large",
			req: &BuatProdukRequest{
				KodeProduk: "PRD001",
				NamaProduk: "Test Product",
				Harga:      1000000000, // > 999,999,999
			},
			wantErr: true,
		},
		{
			name: "valid minimum data",
			req: &BuatProdukRequest{
				KodeProduk: "PRD002",
				NamaProduk: "Valid Product",
				Harga:      1000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.BuatProduk(koperasi.ID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBuatProduk_DuplicateKode tests duplicate product code validation
func TestBuatProduk_DuplicateKode(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create first product
	req1 := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Product 1",
		Harga:      50000,
	}
	_, err := service.BuatProduk(koperasi.ID, req1)
	assert.NoError(t, err)

	// Try to create product with same code
	req2 := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Product 2",
		Harga:      60000,
	}
	_, err = service.BuatProduk(koperasi.ID, req2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sudah digunakan")
}

// TestBuatProduk_MultiTenant tests multi-tenant product isolation
func TestBuatProduk_MultiTenant(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)

	// Create two cooperatives
	koperasi1 := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Koperasi 1", Email: "kop1@test.com", NoTelepon: "081234567890"}
	koperasi2 := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Koperasi 2", Email: "kop2@test.com", NoTelepon: "081234567891"}
	db.Create(koperasi1)
	db.Create(koperasi2)

	// Create product with same code in different cooperatives (should work)
	req1 := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Product Koperasi 1",
		Harga:      50000,
	}
	result1, err := service.BuatProduk(koperasi1.ID, req1)
	assert.NoError(t, err)
	assert.NotNil(t, result1)

	req2 := &BuatProdukRequest{
		KodeProduk: "PRD001",  // Same code but different cooperative
		NamaProduk: "Product Koperasi 2",
		Harga:      60000,
	}
	result2, err := service.BuatProduk(koperasi2.ID, req2)
	assert.NoError(t, err)
	assert.NotNil(t, result2)

	// Verify products are isolated
	assert.NotEqual(t, result1.ID, result2.ID)
}

// TestDapatkanSemuaProduk tests listing products with filters
func TestDapatkanSemuaProduk(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create test products
	products := []BuatProdukRequest{
		{KodeProduk: "EL001", NamaProduk: "Laptop", Kategori: "Elektronik", Harga: 5000000},
		{KodeProduk: "EL002", NamaProduk: "Mouse", Kategori: "Elektronik", Harga: 50000},
		{KodeProduk: "FD001", NamaProduk: "Mie Instant", Kategori: "Makanan", Harga: 3000},
	}

	for _, p := range products {
		service.BuatProduk(koperasi.ID, &p)
	}

	t.Run("get all products", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaProduk(koperasi.ID, "", "", nil, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 3)
	})

	t.Run("filter by category", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaProduk(koperasi.ID, "Elektronik", "", nil, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)
	})

	t.Run("search by name", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaProduk(koperasi.ID, "", "Laptop", nil, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, results, 1)
		assert.Equal(t, "Laptop", results[0].NamaProduk)
	})

	t.Run("pagination", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaProduk(koperasi.ID, "", "", nil, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 2)
	})
}

// TestPerbaruiProduk tests product updates
func TestPerbaruiProduk(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create product
	createReq := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Original Name",
		Harga:      50000,
	}
	produk, _ := service.BuatProduk(koperasi.ID, createReq)

	t.Run("update price", func(t *testing.T) {
		updateReq := &PerbaruiProdukRequest{
			Harga: 60000,
		}
		result, err := service.PerbaruiProduk(koperasi.ID, produk.ID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, 60000.0, result.Harga)
	})

	t.Run("update invalid price", func(t *testing.T) {
		updateReq := &PerbaruiProdukRequest{
			Harga: 1000000000, // Too large
		}
		_, err := service.PerbaruiProduk(koperasi.ID, produk.ID, updateReq)
		assert.Error(t, err)
	})

	t.Run("deactivate product", func(t *testing.T) {
		statusFalse := false
		updateReq := &PerbaruiProdukRequest{
			StatusAktif: &statusFalse,
		}
		result, err := service.PerbaruiProduk(koperasi.ID, produk.ID, updateReq)
		assert.NoError(t, err)
		assert.False(t, result.StatusAktif)
	})
}

// TestKurangiStok tests stock reduction
func TestKurangiStok(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create product with stock
	createReq := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       100,
	}
	produk, _ := service.BuatProduk(koperasi.ID, createReq)

	t.Run("reduce stock success", func(t *testing.T) {
		err := service.KurangiStok(produk.ID, 30)
		assert.NoError(t, err)

		// Verify stock reduced
		updated, _ := service.DapatkanProduk(produk.ID)
		assert.Equal(t, 70, updated.Stok)
	})

	t.Run("insufficient stock", func(t *testing.T) {
		err := service.KurangiStok(produk.ID, 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak mencukupi")
	})

	t.Run("non-existent product", func(t *testing.T) {
		err := service.KurangiStok(uuid.New(), 10)
		assert.Error(t, err)
	})
}

// TestTambahStok tests stock addition
func TestTambahStok(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	createReq := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       100,
	}
	produk, _ := service.BuatProduk(koperasi.ID, createReq)

	t.Run("add stock success", func(t *testing.T) {
		err := service.TambahStok(produk.ID, 50)
		assert.NoError(t, err)

		updated, _ := service.DapatkanProduk(produk.ID)
		assert.Equal(t, 150, updated.Stok)
	})
}

// TestCekStokTersedia tests stock availability check
func TestCekStokTersedia(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	createReq := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       100,
	}
	produk, _ := service.BuatProduk(koperasi.ID, createReq)

	tests := []struct {
		name     string
		jumlah   int
		expected bool
	}{
		{"stock available", 50, true},
		{"stock exactly equal", 100, true},
		{"stock insufficient", 101, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tersedia, err := service.CekStokTersedia(produk.ID, tt.jumlah)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, tersedia)
		})
	}
}

// TestHapusProduk tests product deletion
func TestHapusProduk(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	t.Run("delete product without sales", func(t *testing.T) {
		createReq := &BuatProdukRequest{
			KodeProduk: "PRD001",
			NamaProduk: "Test Product",
			Harga:      50000,
		}
		produk, _ := service.BuatProduk(koperasi.ID, createReq)

		err := service.HapusProduk(koperasi.ID, produk.ID)
		assert.NoError(t, err)

		// Verify deleted (soft delete)
		_, err = service.DapatkanProduk(produk.ID)
		assert.Error(t, err)
	})

	t.Run("delete product with sales (should fail)", func(t *testing.T) {
		createReq := &BuatProdukRequest{
			KodeProduk: "PRD002",
			NamaProduk: "Test Product 2",
			Harga:      50000,
		}
		produk, _ := service.BuatProduk(koperasi.ID, createReq)

		// Create item penjualan (simulating sales)
		penjualan := &models.Penjualan{
			IDKoperasi:       koperasi.ID,
			NomorPenjualan:   "POS-001",
			TotalBelanja:     50000,
			JumlahBayar:      50000,
			IDKasir:          uuid.New(),
		}
		db.Create(penjualan)

		itemPenjualan := &models.ItemPenjualan{
			IDPenjualan: penjualan.ID,
			IDProduk:    produk.ID,
			NamaProduk:  produk.NamaProduk,
			Kuantitas:   1,
			HargaSatuan: 50000,
		}
		db.Create(itemPenjualan)

		err := service.HapusProduk(koperasi.ID, produk.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sudah pernah dijual")
	})
}

// TestDapatkanProdukStokRendah tests low stock products
func TestDapatkanProdukStokRendah(t *testing.T) {
	db := setupProdukTestDB(t)
	if db == nil {
		return
	}

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	// Create products with varying stock levels
	products := []BuatProdukRequest{
		{KodeProduk: "PRD001", NamaProduk: "Low Stock 1", Harga: 10000, Stok: 5, StokMinimum: 10},
		{KodeProduk: "PRD002", NamaProduk: "Low Stock 2", Harga: 10000, Stok: 3, StokMinimum: 5},
		{KodeProduk: "PRD003", NamaProduk: "Good Stock", Harga: 10000, Stok: 100, StokMinimum: 10},
	}

	for _, p := range products {
		service.BuatProduk(koperasi.ID, &p)
	}

	results, err := service.DapatkanProdukStokRendah(koperasi.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // Only 2 products with low stock
}

// BenchmarkBuatProduk benchmarks product creation
func BenchmarkBuatProduk(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{}, &models.Produk{})
	db.Exec("TRUNCATE TABLE produk CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &BuatProdukRequest{
			KodeProduk: fmt.Sprintf("PRD%05d", i),
			NamaProduk: fmt.Sprintf("Product %d", i),
			Harga:      50000,
		}
		_, _ = service.BuatProduk(koperasi.ID, req)
	}
}

// BenchmarkKurangiStok benchmarks stock reduction
func BenchmarkKurangiStok(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{}, &models.Produk{})
	db.Exec("TRUNCATE TABLE produk CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	service := NewProdukService(db)
	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	req := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       1000000, // Large stock for benchmarking
	}
	produk, _ := service.BuatProduk(koperasi.ID, req)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = service.KurangiStok(produk.ID, 1)
	}
}
