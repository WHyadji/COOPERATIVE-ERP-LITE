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
		&models.Produk{},
		&models.Penjualan{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate models: %v", err)
		return nil
	}

	return db
}

// cleanupTestDB cleans up test database
func cleanupTestDB(t *testing.T, db *gorm.DB) {
	if db == nil {
		return
	}

	// Clean up all test data
	db.Exec("DELETE FROM penjualan")
	db.Exec("DELETE FROM produk")
	db.Exec("DELETE FROM transaksi")
	db.Exec("DELETE FROM simpanan")
	db.Exec("DELETE FROM akun")
	db.Exec("DELETE FROM anggota")
	db.Exec("DELETE FROM pengguna")
	db.Exec("DELETE FROM koperasi")
}
