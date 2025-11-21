package integration

import (
	"cooperative-erp-lite/internal/models"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB creates a test database connection for integration tests
func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
		return nil
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Pengguna{},
		&models.Anggota{},
		&models.Simpanan{},
		&models.Akun{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
		&models.Produk{},
		&models.Penjualan{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate models: %v", err)
	}

	return db
}

// cleanupTestData cleans up test data for a specific cooperative
// Uses Unscoped() to delete soft-deleted records as well
func cleanupTestData(db *gorm.DB, koperasiID uuid.UUID) {
	// Clean up in correct order to respect foreign key constraints
	// Unscoped() ensures we delete soft-deleted records too (prevents duplicate key errors)

	// 1. Delete child records first (those with FKs to other tables)
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Penjualan{})
	db.Unscoped().Exec("DELETE FROM baris_transaksi WHERE id_transaksi IN (SELECT id FROM transaksi WHERE id_koperasi = ?)", koperasiID)
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Simpanan{})

	// 2. Delete transaction and product records
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Transaksi{})
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Produk{})

	// 3. Delete accounting records
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Akun{})

	// 4. Delete member and user records
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Anggota{})
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Pengguna{})

	// 5. Finally delete the cooperative itself
	db.Unscoped().Where("id = ?", koperasiID).Delete(&models.Koperasi{})
}

// setupChartOfAccounts creates standard Chart of Accounts for simpanan tests
func setupChartOfAccounts(db *gorm.DB, koperasiID uuid.UUID) {
	// First delete any existing accounts for this koperasi (including soft-deleted ones)
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Akun{})

	accounts := []models.Akun{
		{
			IDKoperasi:  koperasiID,
			KodeAkun:    "1101",
			NamaAkun:    "Kas",
			TipeAkun:    models.AkunAktiva,
			NormalSaldo: "DEBIT",
			Deskripsi:   "Kas dan setara kas",
		},
		{
			IDKoperasi:  koperasiID,
			KodeAkun:    "3101",
			NamaAkun:    "Modal Simpanan Pokok",
			TipeAkun:    models.AkunModal,
			NormalSaldo: "KREDIT",
			Deskripsi:   "Modal dari simpanan pokok anggota",
		},
		{
			IDKoperasi:  koperasiID,
			KodeAkun:    "3102",
			NamaAkun:    "Modal Simpanan Wajib",
			TipeAkun:    models.AkunModal,
			NormalSaldo: "KREDIT",
			Deskripsi:   "Modal dari simpanan wajib anggota",
		},
		{
			IDKoperasi:  koperasiID,
			KodeAkun:    "3103",
			NamaAkun:    "Modal Simpanan Sukarela",
			TipeAkun:    models.AkunModal,
			NormalSaldo: "KREDIT",
			Deskripsi:   "Modal dari simpanan sukarela anggota",
		},
	}

	for _, akun := range accounts {
		db.Create(&akun)
	}
}
