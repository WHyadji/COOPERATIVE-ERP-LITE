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

// setupAnggotaTestDB creates a test database connection for anggota tests
func setupAnggotaTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Anggota{},
		&models.Simpanan{},
		&models.Penjualan{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TestHapusAnggota_WithSimpananTransactions verifies that members with simpanan transactions cannot be deleted
func TestHapusAnggota_WithSimpananTransactions(t *testing.T) {
	db := setupAnggotaTestDB(t)
	if db == nil {
		return
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Delete Member Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test member
	member := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0001",
		NamaLengkap:  "Test Member with Simpanan",
		Status:       models.StatusAktif,
	}
	db.Create(member)

	// Create simpanan transaction for this member
	simpanan := &models.Simpanan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        member.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    100000,
		NomorReferensi:   "SMP-TEST-001",
	}
	db.Create(simpanan)

	// Try to delete the member - should fail
	service := NewAnggotaService(db)
	err := service.HapusAnggota(koperasi.ID, member.ID)

	// Verify deletion is blocked with correct error message
	assert.Error(t, err)
	assert.Equal(t, "tidak dapat menghapus anggota yang memiliki transaksi simpanan", err.Error())

	// Verify member still exists (not soft deleted)
	var memberCheck models.Anggota
	result := db.Unscoped().Where("id = ?", member.ID).First(&memberCheck)
	assert.NoError(t, result.Error)
	assert.False(t, memberCheck.TanggalDihapus.Valid, "Member should not be soft deleted")

	// Cleanup
	db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
	db.Delete(koperasi)
}

// TestHapusAnggota_WithPenjualanTransactions verifies that members with penjualan transactions cannot be deleted
func TestHapusAnggota_WithPenjualanTransactions(t *testing.T) {
	db := setupAnggotaTestDB(t)
	if db == nil {
		return
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Delete Member with Sales Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test member
	member := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0002",
		NamaLengkap:  "Test Member with Penjualan",
		Status:       models.StatusAktif,
	}
	db.Create(member)

	// Create penjualan transaction for this member
	penjualan := &models.Penjualan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        &member.ID,
		TanggalPenjualan: time.Now(),
		NomorPenjualan:   "POS-TEST-001",
		TotalBelanja:     50000,
		MetodePembayaran: models.PembayaranTunai,
		JumlahBayar:      50000,
		IDKasir:          member.ID, // Use member ID as kasir for test
	}
	db.Create(penjualan)

	// Try to delete the member - should fail
	service := NewAnggotaService(db)
	err := service.HapusAnggota(koperasi.ID, member.ID)

	// Verify deletion is blocked with correct error message
	assert.Error(t, err)
	assert.Equal(t, "tidak dapat menghapus anggota yang memiliki transaksi penjualan", err.Error())

	// Verify member still exists (not soft deleted)
	var memberCheck models.Anggota
	result := db.Unscoped().Where("id = ?", member.ID).First(&memberCheck)
	assert.NoError(t, result.Error)
	assert.False(t, memberCheck.TanggalDihapus.Valid, "Member should not be soft deleted")

	// Cleanup
	db.Exec("DELETE FROM penjualan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
	db.Delete(koperasi)
}

// TestHapusAnggota_WithBothTransactionTypes verifies that members with both simpanan and penjualan cannot be deleted
func TestHapusAnggota_WithBothTransactionTypes(t *testing.T) {
	db := setupAnggotaTestDB(t)
	if db == nil {
		return
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Delete Member with Both Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test member
	member := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0003",
		NamaLengkap:  "Test Member with Both Transactions",
		Status:       models.StatusAktif,
	}
	db.Create(member)

	// Create both simpanan and penjualan transactions
	simpanan := &models.Simpanan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        member.ID,
		TipeSimpanan:     models.SimpananWajib,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    50000,
		NomorReferensi:   "SMP-TEST-002",
	}
	db.Create(simpanan)

	penjualan := &models.Penjualan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        &member.ID,
		TanggalPenjualan: time.Now(),
		NomorPenjualan:   "POS-TEST-002",
		TotalBelanja:     75000,
		MetodePembayaran: models.PembayaranTunai,
		JumlahBayar:      75000,
		IDKasir:          member.ID, // Use member ID as kasir for test
	}
	db.Create(penjualan)

	// Try to delete the member - should fail (simpanan check comes first)
	service := NewAnggotaService(db)
	err := service.HapusAnggota(koperasi.ID, member.ID)

	// Verify deletion is blocked
	assert.Error(t, err)
	// Should fail on simpanan check first
	assert.Equal(t, "tidak dapat menghapus anggota yang memiliki transaksi simpanan", err.Error())

	// Cleanup
	db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM penjualan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
	db.Delete(koperasi)
}

// TestHapusAnggota_WithoutTransactions verifies that members without transactions can be deleted successfully
func TestHapusAnggota_WithoutTransactions(t *testing.T) {
	db := setupAnggotaTestDB(t)
	if db == nil {
		return
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Delete Clean Member Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test member without any transactions
	member := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0004",
		NamaLengkap:  "Test Member Without Transactions",
		Status:       models.StatusAktif,
	}
	db.Create(member)

	// Delete the member - should succeed
	service := NewAnggotaService(db)
	err := service.HapusAnggota(koperasi.ID, member.ID)

	// Verify deletion succeeded
	assert.NoError(t, err)

	// Verify member is soft deleted
	var memberCheck models.Anggota
	result := db.Where("id = ?", member.ID).First(&memberCheck)
	assert.Error(t, result.Error, "Should not find member in normal query (soft deleted)")

	// Verify it exists in unscoped query (soft delete)
	result = db.Unscoped().Where("id = ?", member.ID).First(&memberCheck)
	assert.NoError(t, result.Error)
	assert.True(t, memberCheck.TanggalDihapus.Valid, "Member should be soft deleted")

	// Cleanup
	db.Exec("DELETE FROM anggota WHERE id = ?", member.ID)
	db.Delete(koperasi)
}

// TestHapusAnggota_NonExistentMember verifies proper error handling for non-existent members
func TestHapusAnggota_NonExistentMember(t *testing.T) {
	db := setupAnggotaTestDB(t)
	if db == nil {
		return
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Non-Existent Member Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Try to delete a non-existent member
	service := NewAnggotaService(db)
	fakeID := uuid.New()
	err := service.HapusAnggota(koperasi.ID, fakeID)

	// Verify proper error is returned
	assert.Error(t, err)
	assert.Equal(t, "anggota tidak ditemukan atau tidak memiliki akses", err.Error())

	// Cleanup
	db.Delete(koperasi)
}

// TestHapusAnggota_MultiTenantIsolation verifies that members from other cooperatives cannot be deleted
func TestHapusAnggota_MultiTenantIsolation(t *testing.T) {
	db := setupAnggotaTestDB(t)
	if db == nil {
		return
	}

	// Create two cooperatives
	koperasi1 := &models.Koperasi{
		NamaKoperasi: "Koperasi 1",
		Alamat:       "Address 1",
	}
	koperasi2 := &models.Koperasi{
		NamaKoperasi: "Koperasi 2",
		Alamat:       "Address 2",
	}
	db.Create(koperasi1)
	db.Create(koperasi2)

	// Create member for koperasi 1
	member := &models.Anggota{
		IDKoperasi:   koperasi1.ID,
		NomorAnggota: "A0005",
		NamaLengkap:  "Member of Koperasi 1",
		Status:       models.StatusAktif,
	}
	db.Create(member)

	// Try to delete member from koperasi 2 context - should fail
	service := NewAnggotaService(db)
	err := service.HapusAnggota(koperasi2.ID, member.ID)

	// Verify deletion is blocked due to multi-tenant isolation
	assert.Error(t, err)
	assert.Equal(t, "anggota tidak ditemukan atau tidak memiliki akses", err.Error())

	// Verify member still exists
	var memberCheck models.Anggota
	result := db.Where("id = ?", member.ID).First(&memberCheck)
	assert.NoError(t, result.Error)
	assert.False(t, memberCheck.TanggalDihapus.Valid, "Member should not be deleted")

	// Cleanup
	db.Exec("DELETE FROM anggota WHERE id_koperasi IN (?, ?)", koperasi1.ID, koperasi2.ID)
	db.Delete(koperasi1)
	db.Delete(koperasi2)
}
