package services

import (
	"cooperative-erp-lite/internal/models"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test helper untuk penjualan service
func setupPenjualanTestDB(t *testing.T) *gorm.DB {
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
		&models.Penjualan{},
		&models.ItemPenjualan{},
		&models.Pengguna{},
		&models.Anggota{},
		&models.Akun{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean up existing data
	db.Exec("TRUNCATE TABLE baris_transaksi CASCADE")
	db.Exec("TRUNCATE TABLE transaksi CASCADE")
	db.Exec("TRUNCATE TABLE item_penjualan CASCADE")
	db.Exec("TRUNCATE TABLE penjualan CASCADE")
	db.Exec("TRUNCATE TABLE produk CASCADE")
	db.Exec("TRUNCATE TABLE anggota CASCADE")
	db.Exec("TRUNCATE TABLE pengguna CASCADE")
	db.Exec("TRUNCATE TABLE akun CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	return db
}

// TestProsesPenjualan_Success tests successful sales transaction
func TestProsesPenjualan_Success(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	// Create test data
	koperasi := &models.Koperasi{
		ID:          uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	kasir := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaPengguna: "kasir01",
		Email:        "kasir@test.com",
		NamaLengkap:  "Kasir Test",
		Peran:        models.PeranKasir,
		StatusAktif:  true,
	}
	db.Create(kasir)

	// Create product with stock
	produkReq := &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       100,
	}
	produk, _ := produkService.BuatProduk(koperasi.ID, produkReq)

	// Process sale
	saleReq := &ProsesPenjualanRequest{
		Items: []ItemPenjualanRequest{
			{
				IDProduk:    produk.ID,
				Kuantitas:   2,
				HargaSatuan: 50000,
			},
		},
		JumlahBayar: 100000,
		Catatan:     "Test sale",
	}

	result, err := service.ProsesPenjualan(koperasi.ID, kasir.ID, saleReq)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.NomorPenjualan)
	assert.Equal(t, 100000.0, result.TotalBelanja)
	assert.Equal(t, 0.0, result.Kembalian)
	assert.Len(t, result.ItemPenjualan, 1)

	// Verify stock reduced
	updatedProduk, _ := produkService.DapatkanProduk(produk.ID)
	assert.Equal(t, 98, updatedProduk.Stok) // 100 - 2
}

// TestProsesPenjualan_ValidationErrors tests sales validation
func TestProsesPenjualan_ValidationErrors(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	produkReq := &BuatProdukRequest{KodeProduk: "PRD001", NamaProduk: "Test Product", Harga: 50000, Stok: 10}
	produk, _ := produkService.BuatProduk(koperasi.ID, produkReq)

	tests := []struct {
		name    string
		req     *ProsesPenjualanRequest
		wantErr bool
	}{
		{
			name: "negative payment",
			req: &ProsesPenjualanRequest{
				Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 1, HargaSatuan: 50000}},
				JumlahBayar: -1000,
			},
			wantErr: true,
		},
		{
			name: "payment too large",
			req: &ProsesPenjualanRequest{
				Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 1, HargaSatuan: 50000}},
				JumlahBayar: 1000000000, // > 999,999,999
			},
			wantErr: true,
		},
		{
			name: "insufficient payment",
			req: &ProsesPenjualanRequest{
				Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 2, HargaSatuan: 50000}},
				JumlahBayar: 50000, // Less than 100000
			},
			wantErr: true,
		},
		{
			name: "insufficient stock",
			req: &ProsesPenjualanRequest{
				Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 100, HargaSatuan: 50000}},
				JumlahBayar: 5000000,
			},
			wantErr: true,
		},
		{
			name: "zero quantity",
			req: &ProsesPenjualanRequest{
				Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 0, HargaSatuan: 50000}},
				JumlahBayar: 50000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ProsesPenjualan(koperasi.ID, kasir.ID, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestProsesPenjualan_WithChange tests sale with change calculation
func TestProsesPenjualan_WithChange(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	produkReq := &BuatProdukRequest{KodeProduk: "PRD001", NamaProduk: "Test Product", Harga: 25000, Stok: 100}
	produk, _ := produkService.BuatProduk(koperasi.ID, produkReq)

	saleReq := &ProsesPenjualanRequest{
		Items: []ItemPenjualanRequest{
			{IDProduk: produk.ID, Kuantitas: 3, HargaSatuan: 25000},
		},
		JumlahBayar: 100000, // Total: 75000, Change: 25000
	}

	result, err := service.ProsesPenjualan(koperasi.ID, kasir.ID, saleReq)

	assert.NoError(t, err)
	assert.Equal(t, 75000.0, result.TotalBelanja)
	assert.Equal(t, 100000.0, result.JumlahBayar)
	assert.Equal(t, 25000.0, result.Kembalian)
}

// TestProsesPenjualan_MultipleItems tests sale with multiple items
func TestProsesPenjualan_MultipleItems(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	// Create multiple products
	produk1, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{KodeProduk: "PRD001", NamaProduk: "Product 1", Harga: 10000, Stok: 100})
	produk2, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{KodeProduk: "PRD002", NamaProduk: "Product 2", Harga: 20000, Stok: 100})
	produk3, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{KodeProduk: "PRD003", NamaProduk: "Product 3", Harga: 30000, Stok: 100})

	saleReq := &ProsesPenjualanRequest{
		Items: []ItemPenjualanRequest{
			{IDProduk: produk1.ID, Kuantitas: 2, HargaSatuan: 10000},
			{IDProduk: produk2.ID, Kuantitas: 1, HargaSatuan: 20000},
			{IDProduk: produk3.ID, Kuantitas: 3, HargaSatuan: 30000},
		},
		JumlahBayar: 150000, // Total: 2*10k + 1*20k + 3*30k = 130k
	}

	result, err := service.ProsesPenjualan(koperasi.ID, kasir.ID, saleReq)

	assert.NoError(t, err)
	assert.Equal(t, 130000.0, result.TotalBelanja)
	assert.Len(t, result.ItemPenjualan, 3)
	assert.Equal(t, 20000.0, result.Kembalian)

	// Verify all stock reduced correctly
	p1, _ := produkService.DapatkanProduk(produk1.ID)
	p2, _ := produkService.DapatkanProduk(produk2.ID)
	p3, _ := produkService.DapatkanProduk(produk3.ID)

	assert.Equal(t, 98, p1.Stok) // 100 - 2
	assert.Equal(t, 99, p2.Stok) // 100 - 1
	assert.Equal(t, 97, p3.Stok) // 100 - 3
}

// TestGenerateNomorPenjualan_Sequential tests sequential sales number generation
func TestGenerateNomorPenjualan_Sequential(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	db.Create(koperasi)

	today := time.Now()
	dateStr := today.Format("20060102")

	// Generate first sales number
	nomor1, err := service.GenerateNomorPenjualan(koperasi.ID, today)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("POS-%s-0001", dateStr), nomor1)

	// Create sale with first number
	penjualan1 := &models.Penjualan{
		IDKoperasi:       koperasi.ID,
		NomorPenjualan:   nomor1,
		TanggalPenjualan: today,
		TotalBelanja:     50000,
		JumlahBayar:      50000,
		IDKasir:          uuid.New(),
	}
	db.Create(penjualan1)

	// Generate second sales number
	nomor2, err := service.GenerateNomorPenjualan(koperasi.ID, today)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("POS-%s-0002", dateStr), nomor2)
}

// TestProsesPenjualan_Concurrent tests concurrent sales processing
func TestProsesPenjualan_Concurrent(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	produk, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       1000,
	})

	concurrentCount := 10
	results := make([]string, concurrentCount)
	errors := make([]error, concurrentCount)
	var wg sync.WaitGroup

	for i := 0; i < concurrentCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req := &ProsesPenjualanRequest{
				Items: []ItemPenjualanRequest{
					{IDProduk: produk.ID, Kuantitas: 1, HargaSatuan: 50000},
				},
				JumlahBayar: 50000,
			}

			result, err := service.ProsesPenjualan(koperasi.ID, kasir.ID, req)
			errors[index] = err
			if err == nil {
				results[index] = result.NomorPenjualan
			}
		}(i)
	}

	wg.Wait()

	// Verify no errors
	for i, err := range errors {
		assert.NoError(t, err, "Error at index %d", i)
	}

	// Verify all sales numbers are unique
	uniqueNumbers := make(map[string]bool)
	for _, nomor := range results {
		assert.NotEmpty(t, nomor)
		assert.False(t, uniqueNumbers[nomor], "Duplicate sales number: %s", nomor)
		uniqueNumbers[nomor] = true
	}

	assert.Equal(t, concurrentCount, len(uniqueNumbers))
}

// TestDapatkanSemuaPenjualan tests listing sales with filters
func TestDapatkanSemuaPenjualan(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	produk, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       1000,
	})

	// Create multiple sales
	for i := 0; i < 5; i++ {
		req := &ProsesPenjualanRequest{
			Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 1, HargaSatuan: 50000}},
			JumlahBayar: 50000,
		}
		service.ProsesPenjualan(koperasi.ID, kasir.ID, req)
	}

	t.Run("get all sales", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaPenjualan(koperasi.ID, "", "", nil, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, results, 5)
	})

	t.Run("pagination", func(t *testing.T) {
		results, total, err := service.DapatkanSemuaPenjualan(koperasi.ID, "", "", nil, 1, 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, results, 3)
	})

	t.Run("filter by cashier", func(t *testing.T) {
		_, total, err := service.DapatkanSemuaPenjualan(koperasi.ID, "", "", &kasir.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
	})
}

// TestHitungTotalPenjualan tests sales summary calculation
func TestHitungTotalPenjualan(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	produk, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       1000,
	})

	// Create 3 sales
	for i := 0; i < 3; i++ {
		req := &ProsesPenjualanRequest{
			Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 2, HargaSatuan: 50000}},
			JumlahBayar: 100000,
		}
		service.ProsesPenjualan(koperasi.ID, kasir.ID, req)
	}

	today := time.Now().Format("2006-01-02")
	summary, err := service.HitungTotalPenjualan(koperasi.ID, today, today)

	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 300000.0, summary["totalPenjualan"])
	assert.Equal(t, int64(3), summary["jumlahTransaksi"])
	assert.Equal(t, 100000.0, summary["rataRata"])
}

// TestDapatkanPenjualanHariIni tests today's sales summary
func TestDapatkanPenjualanHariIni(t *testing.T) {
	db := setupPenjualanTestDB(t)
	if db == nil {
		return
	}

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	produk, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       1000,
	})

	req := &ProsesPenjualanRequest{
		Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 1, HargaSatuan: 50000}},
		JumlahBayar: 50000,
	}
	service.ProsesPenjualan(koperasi.ID, kasir.ID, req)

	summary, err := service.DapatkanPenjualanHariIni(koperasi.ID)

	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 50000.0, summary["totalPenjualan"])
	assert.Equal(t, int64(1), summary["jumlahTransaksi"])
}

// BenchmarkProsesPenjualan benchmarks sales processing
func BenchmarkProsesPenjualan(b *testing.B) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database")
		return
	}

	db.AutoMigrate(&models.Koperasi{}, &models.Produk{}, &models.Penjualan{}, &models.ItemPenjualan{}, &models.Pengguna{})
	db.Exec("TRUNCATE TABLE item_penjualan CASCADE")
	db.Exec("TRUNCATE TABLE penjualan CASCADE")
	db.Exec("TRUNCATE TABLE produk CASCADE")
	db.Exec("TRUNCATE TABLE pengguna CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	produkService := NewProdukService(db)
	transaksiService := NewTransaksiService(db)
	service := NewPenjualanService(db, produkService, transaksiService)

	koperasi := &models.Koperasi{ID: uuid.New(), NamaKoperasi: "Test", Email: "test@test.com", NoTelepon: "081234567890"}
	kasir := &models.Pengguna{IDKoperasi: koperasi.ID, NamaPengguna: "kasir", Email: "kasir@test.com", NamaLengkap: "Kasir", Peran: models.PeranKasir, StatusAktif: true}
	db.Create(koperasi)
	db.Create(kasir)

	produk, _ := produkService.BuatProduk(koperasi.ID, &BuatProdukRequest{
		KodeProduk: "PRD001",
		NamaProduk: "Test Product",
		Harga:      50000,
		Stok:       100000,
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := &ProsesPenjualanRequest{
			Items:       []ItemPenjualanRequest{{IDProduk: produk.ID, Kuantitas: 1, HargaSatuan: 50000}},
			JumlahBayar: 50000,
		}
		_, _ = service.ProsesPenjualan(koperasi.ID, kasir.ID, req)
	}
}
