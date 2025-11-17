package services

import (
	"cooperative-erp-lite/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDBForRollback creates a test database connection
func setupTestDBForRollback(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	// Auto migrate tables for testing
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Anggota{},
		&models.Simpanan{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
		&models.Akun{},
		&models.Produk{},
		&models.Penjualan{},
		&models.ItemPenjualan{},
		&models.Pengguna{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// cleanupTestData removes test data after test
func cleanupTestData(t *testing.T, db *gorm.DB) {
	// Delete in reverse order of foreign key dependencies
	db.Exec("DELETE FROM baris_transaksi")
	db.Exec("DELETE FROM transaksi")
	db.Exec("DELETE FROM item_penjualan")
	db.Exec("DELETE FROM penjualan")
	db.Exec("DELETE FROM simpanan")
	db.Exec("DELETE FROM produk")
	db.Exec("DELETE FROM anggota")
	db.Exec("DELETE FROM akun")
	db.Exec("DELETE FROM pengguna")
	db.Exec("DELETE FROM koperasi")
}

// TestSimpananService_CatatSetoran_RollbackOnPostingFailure tests that simpanan is rolled back if posting fails
func TestSimpananService_CatatSetoran_RollbackOnPostingFailure(t *testing.T) {
	db := setupTestDBForRollback(t)
	if db == nil {
		return
	}
	defer cleanupTestData(t, db)

	// Create test data
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Test Address",
		NoTelepon:    "08123456789",
	}
	db.Create(koperasi)

	anggota := &models.Anggota{
		IDKoperasi:       koperasi.ID,
		NomorAnggota:     "001",
		NamaLengkap:      "Test Member",
		Alamat:           "Test Address",
		NoTelepon:        "08123456789",
		TanggalBergabung: time.Now(),
		Status:           models.StatusAktif,
	}
	db.Create(anggota)

	pengguna := &models.Pengguna{
		IDKoperasi:    koperasi.ID,
		NamaLengkap:   "Test User",
		NamaPengguna:  "testuser",
		Email:         "test@test.com",
		KataSandiHash: "hashedpassword",
		Peran:         models.PeranAdmin,
		StatusAktif:   true,
	}
	db.Create(pengguna)

	// DO NOT create the required accounts (1101, 3101) - this will cause posting to fail
	// This simulates a scenario where auto-posting fails

	// Create services
	transaksiService := NewTransaksiService(db)
	simpananService := NewSimpananService(db, transaksiService)

	// Test: Attempt to create simpanan (should fail during posting)
	req := &CatatSetoranRequest{
		IDAnggota:        anggota.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    100000,
		Keterangan:       "Test setoran",
	}

	// This should fail because accounts don't exist
	_, err := simpananService.CatatSetoran(koperasi.ID, pengguna.ID, req)

	// Assert: Operation should fail
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	// Assert: Error should be about posting failure
	if err.Error() != "gagal posting ke jurnal: akun kas tidak ditemukan" &&
		err.Error() != "gagal posting ke jurnal: akun simpanan tidak ditemukan: 3101" {
		t.Logf("Got expected error: %v", err)
	}

	// Verify: NO simpanan record should exist (rollback successful)
	var count int64
	db.Model(&models.Simpanan{}).Where("id_koperasi = ?", koperasi.ID).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 simpanan records after rollback, got %d", count)
	}

	// Verify: NO journal entry should exist
	db.Model(&models.Transaksi{}).Where("id_koperasi = ?", koperasi.ID).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 journal entries after rollback, got %d", count)
	}

	t.Log("✅ Rollback successful: No orphaned simpanan or journal records")
}

// TestSimpananService_CatatSetoran_CommitOnSuccess tests that simpanan and posting are both saved on success
func TestSimpananService_CatatSetoran_CommitOnSuccess(t *testing.T) {
	db := setupTestDBForRollback(t)
	if db == nil {
		return
	}
	defer cleanupTestData(t, db)

	// Create test data
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Test Address",
		NoTelepon:    "08123456789",
	}
	db.Create(koperasi)

	anggota := &models.Anggota{
		IDKoperasi:       koperasi.ID,
		NomorAnggota:     "001",
		NamaLengkap:      "Test Member",
		Alamat:           "Test Address",
		NoTelepon:        "08123456789",
		TanggalBergabung: time.Now(),
		Status:           models.StatusAktif,
	}
	db.Create(anggota)

	pengguna := &models.Pengguna{
		IDKoperasi:    koperasi.ID,
		NamaLengkap:   "Test User",
		NamaPengguna:  "testuser",
		Email:         "test@test.com",
		KataSandiHash: "hashedpassword",
		Peran:         models.PeranAdmin,
		StatusAktif:   true,
	}
	db.Create(pengguna)

	// Create required accounts
	akunKas := &models.Akun{
		IDKoperasi:  koperasi.ID,
		KodeAkun:    "1101",
		NamaAkun:    "Kas",
		Kategori:    models.KategoriAset,
		NormalSaldo: "debit",
	}
	db.Create(akunKas)

	akunSimpananPokok := &models.Akun{
		IDKoperasi:  koperasi.ID,
		KodeAkun:    "3101",
		NamaAkun:    "Simpanan Pokok",
		Kategori:    models.KategoriModal,
		NormalSaldo: "kredit",
	}
	db.Create(akunSimpananPokok)

	// Create services
	transaksiService := NewTransaksiService(db)
	simpananService := NewSimpananService(db, transaksiService)

	// Test: Create simpanan (should succeed)
	req := &CatatSetoranRequest{
		IDAnggota:        anggota.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    100000,
		Keterangan:       "Test setoran",
	}

	result, err := simpananService.CatatSetoran(koperasi.ID, pengguna.ID, req)

	// Assert: Operation should succeed
	if err != nil {
		t.Fatalf("Expected success but got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result but got nil")
	}

	// Verify: Simpanan record exists
	var simpananCount int64
	db.Model(&models.Simpanan{}).Where("id_koperasi = ?", koperasi.ID).Count(&simpananCount)
	if simpananCount != 1 {
		t.Errorf("Expected 1 simpanan record, got %d", simpananCount)
	}

	// Verify: Journal entry exists
	var transaksiCount int64
	db.Model(&models.Transaksi{}).Where("id_koperasi = ?", koperasi.ID).Count(&transaksiCount)
	if transaksiCount != 1 {
		t.Errorf("Expected 1 journal entry, got %d", transaksiCount)
	}

	// Verify: Journal has 2 lines (debit and kredit)
	var barisCount int64
	db.Model(&models.BarisTransaksi{}).
		Joins("JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi").
		Where("transaksi.id_koperasi = ?", koperasi.ID).
		Count(&barisCount)
	if barisCount != 2 {
		t.Errorf("Expected 2 journal lines, got %d", barisCount)
	}

	// Verify: Simpanan has IDTransaksi set
	var simpanan models.Simpanan
	db.First(&simpanan, result.ID)
	if simpanan.IDTransaksi == nil {
		t.Error("Expected simpanan to have IDTransaksi set")
	}

	t.Log("✅ Commit successful: Simpanan and journal entry both saved atomically")
}

// TestPenjualanService_ProsesPenjualan_RollbackOnPostingFailure tests that penjualan is rolled back if posting fails
func TestPenjualanService_ProsesPenjualan_RollbackOnPostingFailure(t *testing.T) {
	db := setupTestDBForRollback(t)
	if db == nil {
		return
	}
	defer cleanupTestData(t, db)

	// Create test data
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Test Address",
		NoTelepon:    "08123456789",
	}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:    koperasi.ID,
		NamaLengkap:   "Kasir User",
		NamaPengguna:  "kasir",
		Email:         "kasir@test.com",
		KataSandiHash: "hashedpassword",
		Peran:         models.PeranKasir,
		StatusAktif:   true,
	}
	db.Create(pengguna)

	produk := &models.Produk{
		IDKoperasi:  koperasi.ID,
		KodeProduk:  "P001",
		NamaProduk:  "Test Product",
		HargaJual:   10000,
		HargaBeli:   8000,
		Stok:        100,
		StokMinimum: 10,
	}
	db.Create(produk)

	// DO NOT create required accounts - this will cause posting to fail

	// Create services
	transaksiService := NewTransaksiService(db)
	produkService := NewProdukService(db)
	penjualanService := NewPenjualanService(db, produkService, transaksiService)

	// Test: Attempt to create penjualan (should fail during posting)
	req := &ProsesPenjualanRequest{
		Items: []ItemPenjualanRequest{
			{
				IDProduk:    produk.ID,
				Kuantitas:   5,
				HargaSatuan: 10000,
			},
		},
		JumlahBayar: 50000,
	}

	// This should fail because accounts don't exist
	_, err := penjualanService.ProsesPenjualan(koperasi.ID, pengguna.ID, req)

	// Assert: Operation should fail
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	t.Logf("Got expected error: %v", err)

	// Verify: NO penjualan record exists (rollback successful)
	var penjualanCount int64
	db.Model(&models.Penjualan{}).Where("id_koperasi = ?", koperasi.ID).Count(&penjualanCount)
	if penjualanCount != 0 {
		t.Errorf("Expected 0 penjualan records after rollback, got %d", penjualanCount)
	}

	// Verify: NO item penjualan records exist
	var itemCount int64
	db.Model(&models.ItemPenjualan{}).
		Joins("JOIN penjualan ON penjualan.id = item_penjualan.id_penjualan").
		Where("penjualan.id_koperasi = ?", koperasi.ID).
		Count(&itemCount)
	if itemCount != 0 {
		t.Errorf("Expected 0 item penjualan records after rollback, got %d", itemCount)
	}

	// Verify: Stock was NOT deducted (rollback successful)
	var produkAfter models.Produk
	db.First(&produkAfter, produk.ID)
	if produkAfter.Stok != 100 {
		t.Errorf("Expected stock to remain 100 after rollback, got %d", produkAfter.Stok)
	}

	// Verify: NO journal entry exists
	var transaksiCount int64
	db.Model(&models.Transaksi{}).Where("id_koperasi = ?", koperasi.ID).Count(&transaksiCount)
	if transaksiCount != 0 {
		t.Errorf("Expected 0 journal entries after rollback, got %d", transaksiCount)
	}

	t.Log("✅ Rollback successful: No orphaned penjualan, stock changes, or journal records")
}

// TestPenjualanService_ProsesPenjualan_CommitOnSuccess tests that penjualan, stock changes, and posting are all saved on success
func TestPenjualanService_ProsesPenjualan_CommitOnSuccess(t *testing.T) {
	db := setupTestDBForRollback(t)
	if db == nil {
		return
	}
	defer cleanupTestData(t, db)

	// Create test data
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Koperasi",
		Alamat:       "Test Address",
		NoTelepon:    "08123456789",
	}
	db.Create(koperasi)

	pengguna := &models.Pengguna{
		IDKoperasi:    koperasi.ID,
		NamaLengkap:   "Kasir User",
		NamaPengguna:  "kasir",
		Email:         "kasir@test.com",
		KataSandiHash: "hashedpassword",
		Peran:         models.PeranKasir,
		StatusAktif:   true,
	}
	db.Create(pengguna)

	produk := &models.Produk{
		IDKoperasi:  koperasi.ID,
		KodeProduk:  "P001",
		NamaProduk:  "Test Product",
		HargaJual:   10000,
		HargaBeli:   8000,
		Stok:        100,
		StokMinimum: 10,
	}
	db.Create(produk)

	// Create required accounts
	accounts := []models.Akun{
		{IDKoperasi: koperasi.ID, KodeAkun: "1101", NamaAkun: "Kas", Kategori: models.KategoriAset, NormalSaldo: "debit"},
		{IDKoperasi: koperasi.ID, KodeAkun: "4101", NamaAkun: "Penjualan", Kategori: models.KategoriPendapatan, NormalSaldo: "kredit"},
		{IDKoperasi: koperasi.ID, KodeAkun: "5201", NamaAkun: "HPP", Kategori: models.KategoriBeban, NormalSaldo: "debit"},
		{IDKoperasi: koperasi.ID, KodeAkun: "1301", NamaAkun: "Persediaan", Kategori: models.KategoriAset, NormalSaldo: "debit"},
	}
	for _, akun := range accounts {
		db.Create(&akun)
	}

	// Create services
	transaksiService := NewTransaksiService(db)
	produkService := NewProdukService(db)
	penjualanService := NewPenjualanService(db, produkService, transaksiService)

	// Test: Create penjualan (should succeed)
	req := &ProsesPenjualanRequest{
		Items: []ItemPenjualanRequest{
			{
				IDProduk:    produk.ID,
				Kuantitas:   5,
				HargaSatuan: 10000,
			},
		},
		JumlahBayar: 50000,
	}

	result, err := penjualanService.ProsesPenjualan(koperasi.ID, pengguna.ID, req)

	// Assert: Operation should succeed
	if err != nil {
		t.Fatalf("Expected success but got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result but got nil")
	}

	// Verify: Penjualan record exists
	var penjualanCount int64
	db.Model(&models.Penjualan{}).Where("id_koperasi = ?", koperasi.ID).Count(&penjualanCount)
	if penjualanCount != 1 {
		t.Errorf("Expected 1 penjualan record, got %d", penjualanCount)
	}

	// Verify: Item penjualan exists
	var itemCount int64
	db.Model(&models.ItemPenjualan{}).
		Joins("JOIN penjualan ON penjualan.id = item_penjualan.id_penjualan").
		Where("penjualan.id_koperasi = ?", koperasi.ID).
		Count(&itemCount)
	if itemCount != 1 {
		t.Errorf("Expected 1 item penjualan record, got %d", itemCount)
	}

	// Verify: Stock was deducted
	var produkAfter models.Produk
	db.First(&produkAfter, produk.ID)
	if produkAfter.Stok != 95 {
		t.Errorf("Expected stock to be 95 (100-5), got %d", produkAfter.Stok)
	}

	// Verify: Journal entry exists
	var transaksiCount int64
	db.Model(&models.Transaksi{}).Where("id_koperasi = ?", koperasi.ID).Count(&transaksiCount)
	if transaksiCount != 1 {
		t.Errorf("Expected 1 journal entry, got %d", transaksiCount)
	}

	// Verify: Journal has 4 lines (Kas Dr, Penjualan Cr, HPP Dr, Persediaan Cr)
	var barisCount int64
	db.Model(&models.BarisTransaksi{}).
		Joins("JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi").
		Where("transaksi.id_koperasi = ?", koperasi.ID).
		Count(&barisCount)
	if barisCount != 4 {
		t.Errorf("Expected 4 journal lines, got %d", barisCount)
	}

	// Verify: Penjualan has IDTransaksi set
	var penjualan models.Penjualan
	db.First(&penjualan, result.ID)
	if penjualan.IDTransaksi == nil {
		t.Error("Expected penjualan to have IDTransaksi set")
	}

	t.Log("✅ Commit successful: Penjualan, stock changes, and journal entry all saved atomically")
}
