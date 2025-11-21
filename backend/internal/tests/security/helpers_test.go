package security

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cooperative-erp-lite/internal/models"
)

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return nil
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Pengguna{},
		&models.Anggota{},
		&models.Akun{},
		&models.Simpanan{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
		&models.Produk{},
		&models.Penjualan{},
		&models.ItemPenjualan{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate models: %v", err)
		return nil
	}

	return db
}

// cleanupTestDB cleans up test database
// Uses Unscoped() to delete soft-deleted records and respects foreign key constraints
func cleanupTestDB(t *testing.T, db *gorm.DB) {
	if db == nil {
		return
	}

	// Clean up in correct order to respect foreign key constraints
	// Unscoped() ensures we delete soft-deleted records too (prevents duplicate key errors)

	// 1. Delete child records first (those with FKs to other tables)
	db.Exec("DELETE FROM item_penjualan")
	db.Unscoped().Delete(&models.Penjualan{})
	db.Exec("DELETE FROM baris_transaksi")
	db.Unscoped().Delete(&models.Simpanan{})

	// 2. Delete transaction and product records
	db.Unscoped().Delete(&models.Transaksi{})
	db.Unscoped().Delete(&models.Produk{})

	// 3. Delete accounting records
	db.Unscoped().Delete(&models.Akun{})

	// 4. Delete member and user records
	db.Unscoped().Delete(&models.Anggota{})
	db.Unscoped().Delete(&models.Pengguna{})

	// 5. Finally delete the cooperative itself
	db.Unscoped().Delete(&models.Koperasi{})
}
