package config

import (
	"cooperative-erp-lite/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database instance
var DB *gorm.DB

// InitDatabase menginisialisasi koneksi database dan migrasi
func InitDatabase(config *Config) error {
	var err error

	// Konfigurasi GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Buka koneksi database
	DB, err = gorm.Open(postgres.Open(config.GetDSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("gagal menghubungkan ke database: %w", err)
	}

	log.Println("Koneksi database berhasil")

	// Auto migrate semua model
	err = AutoMigrate()
	if err != nil {
		return fmt.Errorf("gagal melakukan migrasi database: %w", err)
	}

	log.Println("Migrasi database berhasil")

	return nil
}

// AutoMigrate melakukan auto migration untuk semua model
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Koperasi{},
		&models.Pengguna{},
		&models.Anggota{},
		&models.Simpanan{},
		&models.Akun{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
		&models.Produk{},
		&models.Penjualan{},
		&models.ItemPenjualan{},
	)
}

// GetDB mengembalikan instance database
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase menutup koneksi database
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
