package services

import (
	"cooperative-erp-lite/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupLaporanTestDB creates a test database for laporan service
func setupLaporanTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	// Auto-migrate all required models
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Akun{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
		&models.Anggota{},
		&models.Simpanan{},
		&models.Penjualan{},
		&models.ItemPenjualan{},
		&models.Produk{},
		&models.Pengguna{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean up existing data
	db.Exec("TRUNCATE TABLE baris_transaksi CASCADE")
	db.Exec("TRUNCATE TABLE transaksi CASCADE")
	db.Exec("TRUNCATE TABLE item_penjualan CASCADE")
	db.Exec("TRUNCATE TABLE penjualan CASCADE")
	db.Exec("TRUNCATE TABLE simpanan CASCADE")
	db.Exec("TRUNCATE TABLE anggota CASCADE")
	db.Exec("TRUNCATE TABLE produk CASCADE")
	db.Exec("TRUNCATE TABLE akun CASCADE")
	db.Exec("TRUNCATE TABLE pengguna CASCADE")
	db.Exec("TRUNCATE TABLE koperasi CASCADE")

	return db
}

// createTestKoperasiForLaporan creates a test cooperative with COA
func createTestKoperasiForLaporan(_ *testing.T, db *gorm.DB) (*models.Koperasi, *AkunService) {
	koperasi := &models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@koperasi.com",
		NoTelepon:    "081234567890",
	}
	db.Create(koperasi)

	// Create basic COA
	accounts := []models.Akun{
		{ID: uuid.New(), IDKoperasi: koperasi.ID, KodeAkun: "1101", NamaAkun: "Kas", TipeAkun: models.AkunAktiva, NormalSaldo: "DEBIT"},
		{ID: uuid.New(), IDKoperasi: koperasi.ID, KodeAkun: "1102", NamaAkun: "Bank", TipeAkun: models.AkunAktiva, NormalSaldo: "DEBIT"},
		{ID: uuid.New(), IDKoperasi: koperasi.ID, KodeAkun: "2101", NamaAkun: "Utang", TipeAkun: models.AkunKewajiban, NormalSaldo: "KREDIT"},
		{ID: uuid.New(), IDKoperasi: koperasi.ID, KodeAkun: "3101", NamaAkun: "Modal", TipeAkun: models.AkunModal, NormalSaldo: "KREDIT"},
		{ID: uuid.New(), IDKoperasi: koperasi.ID, KodeAkun: "4101", NamaAkun: "Pendapatan Penjualan", TipeAkun: models.AkunPendapatan, NormalSaldo: "KREDIT"},
		{ID: uuid.New(), IDKoperasi: koperasi.ID, KodeAkun: "5101", NamaAkun: "Beban Gaji", TipeAkun: models.AkunBeban, NormalSaldo: "DEBIT"},
	}

	for _, akun := range accounts {
		db.Create(&akun)
	}

	akunService := NewAkunService(db)
	return koperasi, akunService
}

// TestGenerateLaporanPosisiKeuangan tests balance sheet generation
func TestGenerateLaporanPosisiKeuangan(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)
	laporanService := NewLaporanService(db, akunService, nil, nil)

	// Create test transactions
	transaksi := &models.Transaksi{
		ID:               uuid.New(),
		IDKoperasi:       koperasi.ID,
		NomorJurnal:      "JU001",
		TanggalTransaksi: time.Now(),
		Deskripsi:        "Test Transaction",
	}
	db.Create(transaksi)

	// Get account IDs
	var kasAkun, modalAkun models.Akun
	db.Where("kode_akun = ? AND id_koperasi = ?", "1101", koperasi.ID).First(&kasAkun)
	db.Where("kode_akun = ? AND id_koperasi = ?", "3101", koperasi.ID).First(&modalAkun)

	// Create transaction lines: Debit Kas 1000000, Kredit Modal 1000000
	barisTransaksi := []models.BarisTransaksi{
		{IDTransaksi: transaksi.ID, IDAkun: kasAkun.ID, JumlahDebit: 1000000, JumlahKredit: 0},
		{IDTransaksi: transaksi.ID, IDAkun: modalAkun.ID, JumlahDebit: 0, JumlahKredit: 1000000},
	}
	for _, baris := range barisTransaksi {
		db.Create(&baris)
	}

	t.Run("successful balance sheet generation", func(t *testing.T) {
		laporan, err := laporanService.GenerateLaporanPosisiKeuangan(koperasi.ID, "")

		assert.NoError(t, err)
		assert.NotNil(t, laporan)
		assert.Equal(t, float64(1000000), laporan.TotalAset)
		assert.Equal(t, float64(1000000), laporan.TotalModal)
		assert.True(t, len(laporan.Aset) > 0)
		assert.True(t, len(laporan.Modal) > 0)
	})

	t.Run("with specific date", func(t *testing.T) {
		tanggalPer := time.Now().Format("2006-01-02")
		laporan, err := laporanService.GenerateLaporanPosisiKeuangan(koperasi.ID, tanggalPer)

		assert.NoError(t, err)
		assert.NotNil(t, laporan)
	})

	t.Run("invalid date format", func(t *testing.T) {
		laporan, err := laporanService.GenerateLaporanPosisiKeuangan(koperasi.ID, "invalid-date")

		assert.Error(t, err)
		assert.Nil(t, laporan)
		assert.Contains(t, err.Error(), "format tanggal tidak valid")
	})
}

// TestGenerateLaporanLabaRugi tests income statement generation
func TestGenerateLaporanLabaRugi(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)
	laporanService := NewLaporanService(db, akunService, nil, nil)

	// Create test transactions
	transaksi := &models.Transaksi{
		ID:               uuid.New(),
		IDKoperasi:       koperasi.ID,
		NomorJurnal:      "JU002",
		TanggalTransaksi: time.Now(),
		Deskripsi:        "Sales Transaction",
	}
	db.Create(transaksi)

	// Get account IDs
	var kasAkun, pendapatanAkun, bebanAkun models.Akun
	db.Where("kode_akun = ? AND id_koperasi = ?", "1101", koperasi.ID).First(&kasAkun)
	db.Where("kode_akun = ? AND id_koperasi = ?", "4101", koperasi.ID).First(&pendapatanAkun)
	db.Where("kode_akun = ? AND id_koperasi = ?", "5101", koperasi.ID).First(&bebanAkun)

	// Create transaction: Pendapatan 5000000, Beban 2000000, Profit 3000000
	barisTransaksi := []models.BarisTransaksi{
		{IDTransaksi: transaksi.ID, IDAkun: kasAkun.ID, JumlahDebit: 5000000, JumlahKredit: 0},
		{IDTransaksi: transaksi.ID, IDAkun: pendapatanAkun.ID, JumlahDebit: 0, JumlahKredit: 5000000},
	}
	for _, baris := range barisTransaksi {
		db.Create(&baris)
	}

	// Create expense transaction
	transaksi2 := &models.Transaksi{
		ID:               uuid.New(),
		IDKoperasi:       koperasi.ID,
		NomorJurnal:      "JU003",
		TanggalTransaksi: time.Now(),
		Deskripsi:        "Expense Transaction",
	}
	db.Create(transaksi2)

	barisTransaksi2 := []models.BarisTransaksi{
		{IDTransaksi: transaksi2.ID, IDAkun: bebanAkun.ID, JumlahDebit: 2000000, JumlahKredit: 0},
		{IDTransaksi: transaksi2.ID, IDAkun: kasAkun.ID, JumlahDebit: 0, JumlahKredit: 2000000},
	}
	for _, baris := range barisTransaksi2 {
		db.Create(&baris)
	}

	t.Run("successful income statement generation", func(t *testing.T) {
		tanggalMulai := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		tanggalAkhir := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

		laporan, err := laporanService.GenerateLaporanLabaRugi(koperasi.ID, tanggalMulai, tanggalAkhir)

		assert.NoError(t, err)
		assert.NotNil(t, laporan)
		assert.Equal(t, float64(5000000), laporan.TotalPendapatan)
		assert.Equal(t, float64(2000000), laporan.TotalBeban)
		assert.Equal(t, float64(3000000), laporan.LabaRugiBersih)
	})

	t.Run("invalid start date format", func(t *testing.T) {
		laporan, err := laporanService.GenerateLaporanLabaRugi(koperasi.ID, "invalid", "2025-01-01")

		assert.Error(t, err)
		assert.Nil(t, laporan)
		assert.Contains(t, err.Error(), "format tanggal mulai tidak valid")
	})

	t.Run("invalid end date format", func(t *testing.T) {
		laporan, err := laporanService.GenerateLaporanLabaRugi(koperasi.ID, "2025-01-01", "invalid")

		assert.Error(t, err)
		assert.Nil(t, laporan)
		assert.Contains(t, err.Error(), "format tanggal akhir tidak valid")
	})
}

// TestGenerateLaporanArusKas tests cash flow statement generation
func TestGenerateLaporanArusKas(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)
	laporanService := NewLaporanService(db, akunService, nil, nil)

	t.Run("successful cash flow generation", func(t *testing.T) {
		tanggalMulai := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		tanggalAkhir := time.Now().Format("2006-01-02")

		laporan, err := laporanService.GenerateLaporanArusKas(koperasi.ID, tanggalMulai, tanggalAkhir)

		assert.NoError(t, err)
		assert.NotNil(t, laporan)
		assert.NotNil(t, laporan.ArusKasOperasional)
		assert.NotNil(t, laporan.ArusKasInvestasi)
		assert.NotNil(t, laporan.ArusKasPendanaan)
	})

	t.Run("invalid date format", func(t *testing.T) {
		laporan, err := laporanService.GenerateLaporanArusKas(koperasi.ID, "invalid", "2025-01-01")

		assert.Error(t, err)
		assert.Nil(t, laporan)
		assert.Contains(t, err.Error(), "format tanggal")
	})

	t.Run("kas account not found", func(t *testing.T) {
		// Create new cooperative without kas account
		newKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Koperasi Tanpa Kas",
			Email:        "nokas@koperasi.com",
			NoTelepon:    "081234567891",
		}
		db.Create(newKoperasi)

		tanggalMulai := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		tanggalAkhir := time.Now().Format("2006-01-02")

		laporan, err := laporanService.GenerateLaporanArusKas(newKoperasi.ID, tanggalMulai, tanggalAkhir)

		assert.Error(t, err)
		assert.Nil(t, laporan)
		assert.Contains(t, err.Error(), "akun kas tidak ditemukan")
	})
}

// TestGenerateLaporanTransaksiHarian tests daily transaction report
func TestGenerateLaporanTransaksiHarian(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)

	// Create mock services
	simpananService := &SimpananService{db: db}
	penjualanService := &PenjualanService{db: db}

	laporanService := NewLaporanService(db, akunService, simpananService, penjualanService)

	// Create test data
	tanggalHariIni := time.Now().Format("2006-01-02")

	// Create penjualan
	pengguna := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Kasir",
		NamaPengguna: "kasir1",
		Email:        "kasir@example.com",
		Peran:        models.PeranKasir,
		StatusAktif:  true,
	}
	db.Create(pengguna)

	penjualan := &models.Penjualan{
		IDKoperasi:       koperasi.ID,
		IDKasir:          pengguna.ID,
		NomorPenjualan:   "POS001",
		TanggalPenjualan: time.Now(),
		TotalBelanja:     100000,
		JumlahBayar:      100000,
		Kembalian:        0,
	}
	db.Create(penjualan)

	t.Run("successful daily report", func(t *testing.T) {
		laporan, err := laporanService.GenerateLaporanTransaksiHarian(koperasi.ID, tanggalHariIni)

		assert.NoError(t, err)
		assert.NotNil(t, laporan)
		assert.Equal(t, int64(1), laporan.JumlahPenjualan)
	})

	t.Run("invalid date format", func(t *testing.T) {
		laporan, err := laporanService.GenerateLaporanTransaksiHarian(koperasi.ID, "invalid-date")

		assert.Error(t, err)
		assert.Nil(t, laporan)
		assert.Contains(t, err.Error(), "format tanggal tidak valid")
	})

	t.Run("kas account not found", func(t *testing.T) {
		newKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Koperasi No Kas",
			Email:        "nokas2@koperasi.com",
			NoTelepon:    "081234567892",
		}
		db.Create(newKoperasi)

		laporan, err := laporanService.GenerateLaporanTransaksiHarian(newKoperasi.ID, tanggalHariIni)

		assert.Error(t, err)
		assert.Nil(t, laporan)
		assert.Contains(t, err.Error(), "akun kas tidak ditemukan")
	})
}

// TestGetDashboardStats tests dashboard statistics
func TestGetDashboardStats(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)
	simpananService := &SimpananService{db: db}
	penjualanService := &PenjualanService{db: db}
	laporanService := NewLaporanService(db, akunService, simpananService, penjualanService)

	// Create test anggota
	anggota := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NamaLengkap:  "Test Member",
		NomorAnggota: "A001",
		JenisKelamin: "L",
		Status:       models.StatusAktif,
	}
	db.Create(anggota)

	t.Run("get dashboard stats", func(t *testing.T) {
		stats, err := laporanService.GetDashboardStats(koperasi.ID)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Contains(t, stats, "totalAnggota")
		assert.Equal(t, int64(1), stats["totalAnggota"])
	})
}

// TestGenerateNeracaSaldo tests trial balance generation
func TestGenerateNeracaSaldo(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)
	laporanService := NewLaporanService(db, akunService, nil, nil)

	// Create balanced transaction
	transaksi := &models.Transaksi{
		ID:               uuid.New(),
		IDKoperasi:       koperasi.ID,
		NomorJurnal:      "JU004",
		TanggalTransaksi: time.Now(),
		Deskripsi:        "Balanced Transaction",
	}
	db.Create(transaksi)

	var kasAkun, modalAkun models.Akun
	db.Where("kode_akun = ? AND id_koperasi = ?", "1101", koperasi.ID).First(&kasAkun)
	db.Where("kode_akun = ? AND id_koperasi = ?", "3101", koperasi.ID).First(&modalAkun)

	// Create balanced entry
	barisTransaksi := []models.BarisTransaksi{
		{IDTransaksi: transaksi.ID, IDAkun: kasAkun.ID, JumlahDebit: 1000000, JumlahKredit: 0},
		{IDTransaksi: transaksi.ID, IDAkun: modalAkun.ID, JumlahDebit: 0, JumlahKredit: 1000000},
	}
	for _, baris := range barisTransaksi {
		db.Create(&baris)
	}

	t.Run("successful trial balance", func(t *testing.T) {
		tanggalPer := time.Now().Format("2006-01-02")
		result, err := laporanService.GenerateNeracaSaldo(koperasi.ID, tanggalPer)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result, "items")
		assert.Contains(t, result, "totalDebit")
		assert.Contains(t, result, "totalKredit")
		assert.Contains(t, result, "isBalanced")
		assert.Equal(t, result["totalDebit"], result["totalKredit"])
		assert.True(t, result["isBalanced"].(bool))
	})
}

// TestGenerateBukuBesar tests general ledger generation
func TestGenerateBukuBesar(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)
	laporanService := NewLaporanService(db, akunService, nil, nil)

	var kasAkun models.Akun
	db.Where("kode_akun = ? AND id_koperasi = ?", "1101", koperasi.ID).First(&kasAkun)

	t.Run("successful general ledger", func(t *testing.T) {
		tanggalMulai := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		tanggalAkhir := time.Now().Format("2006-01-02")

		result, err := laporanService.GenerateBukuBesar(koperasi.ID, kasAkun.ID, tanggalMulai, tanggalAkhir)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result, "akun")
		assert.Contains(t, result, "periode")
		assert.Contains(t, result, "transaksi")
		assert.Contains(t, result, "saldoAkhir")
	})

	t.Run("multi-tenant validation", func(t *testing.T) {
		otherKoperasi := &models.Koperasi{
			ID:           uuid.New(),
			NamaKoperasi: "Other Koperasi",
			Email:        "other@koperasi.com",
			NoTelepon:    "081234567893",
		}
		db.Create(otherKoperasi)

		tanggalMulai := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		tanggalAkhir := time.Now().Format("2006-01-02")

		result, err := laporanService.GenerateBukuBesar(otherKoperasi.ID, kasAkun.ID, tanggalMulai, tanggalAkhir)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan atau tidak memiliki akses")
	})

	t.Run("non-existing account", func(t *testing.T) {
		tanggalMulai := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		tanggalAkhir := time.Now().Format("2006-01-02")

		result, err := laporanService.GenerateBukuBesar(koperasi.ID, uuid.New(), tanggalMulai, tanggalAkhir)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestGenerateLaporanPerubahanModal tests statement of changes in equity (placeholder)
func TestGenerateLaporanPerubahanModal(t *testing.T) {
	db := setupLaporanTestDB(t)
	if db == nil {
		return
	}

	koperasi, akunService := createTestKoperasiForLaporan(t, db)
	laporanService := NewLaporanService(db, akunService, nil, nil)

	t.Run("placeholder implementation", func(t *testing.T) {
		tanggalMulai := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		tanggalAkhir := time.Now().Format("2006-01-02")

		result, err := laporanService.GenerateLaporanPerubahanModal(koperasi.ID, tanggalMulai, tanggalAkhir)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result, "message")
	})
}
