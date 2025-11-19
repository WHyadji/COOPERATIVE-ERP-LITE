// ============================================================================
// Test Data Seeding Script for E2E Testing
// Creates test cooperative, member, and transactions for Playwright tests
// ============================================================================

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cooperative-erp-lite/internal/models"
)

// Database connection configuration
const (
	DBHost     = "localhost"
	DBUser     = "postgres"
	DBPassword = "postgres"
	DBName     = "koperasi_erp"
	DBPort     = "5432"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("E2E Test Data Seeding Script")
	fmt.Println("===========================================")

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DBHost, DBUser, DBPassword, DBName, DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✓ Connected to database")

	// Start transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Failed to start transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("Panic occurred, rolling back: %v", r)
		}
	}()

	// 1. Create or find test cooperative
	fmt.Println("\n1. Creating test cooperative...")
	koperasi := createTestKoperasi(tx)
	fmt.Printf("   ✓ Koperasi ID: %s\n", koperasi.ID)

	// 2. Create test member (A001 with PIN 123456)
	fmt.Println("\n2. Creating test member...")
	anggota := createTestMember(tx, koperasi.ID)
	fmt.Printf("   ✓ Member: %s (%s)\n", anggota.NamaLengkap, anggota.NomorAnggota)
	fmt.Printf("   ✓ PIN: 123456 (hashed)\n")

	// 3. Create initial balance
	fmt.Println("\n3. Creating initial balance...")
	saldo := createInitialBalance(tx, anggota.ID, koperasi.ID)
	fmt.Printf("   ✓ Simpanan Pokok: Rp %,.0f\n", saldo.SimpananPokok)
	fmt.Printf("   ✓ Simpanan Wajib: Rp %,.0f\n", saldo.SimpananWajib)
	fmt.Printf("   ✓ Simpanan Sukarela: Rp %,.0f\n", saldo.SimpananSukarela)

	// 4. Create sample transactions
	fmt.Println("\n4. Creating sample transactions...")
	transactions := createSampleTransactions(tx, anggota.ID, koperasi.ID)
	fmt.Printf("   ✓ Created %d transactions\n", len(transactions))

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("\n===========================================")
	fmt.Println("✓ Test data seeding completed successfully!")
	fmt.Println("===========================================")
	fmt.Println("\nTest Credentials:")
	fmt.Println("  Nomor Anggota: A001")
	fmt.Println("  PIN: 123456")
	fmt.Println("\nYou can now run E2E tests with:")
	fmt.Println("  cd frontend && npx playwright test")
	fmt.Println("===========================================")
}

// createTestKoperasi creates or finds the test cooperative
func createTestKoperasi(db *gorm.DB) *models.Koperasi {
	var koperasi models.Koperasi

	// Try to find existing test cooperative
	err := db.Where("nomor_badan_hukum = ?", "TEST-E2E-001").First(&koperasi).Error
	if err == nil {
		fmt.Println("   → Using existing test cooperative")
		return &koperasi
	}

	// Create new test cooperative
	koperasi = models.Koperasi{
		ID:               uuid.New(),
		NamaKoperasi:     "Koperasi Test E2E",
		NomorBadanHukum:  "TEST-E2E-001",
		TanggalBerdiri:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Alamat:           "Jl. Test E2E No. 123",
		Kelurahan:        "Test Kelurahan",
		Kecamatan:        "Test Kecamatan",
		KotaKabupaten:    "Test City",
		Provinsi:         "Test Province",
		KodePos:          "12345",
		NomorTelepon:     "08123456789",
		Email:            "test@e2e.com",
		Website:          "https://test-e2e.com",
		JumlahAnggota:    1,
		TotalAset:        5000000,
		SimpananPokok:    1000000,
		SimpananWajib:    500000,
		SimpananSukarela: 200000,
		Aktif:            true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(&koperasi).Error; err != nil {
		log.Fatalf("Failed to create test cooperative: %v", err)
	}

	return &koperasi
}

// createTestMember creates the test member (A001)
func createTestMember(db *gorm.DB, koperasiID uuid.UUID) *models.Anggota {
	var anggota models.Anggota

	// Try to find existing test member
	err := db.Where("nomor_anggota = ? AND id_koperasi = ?", "A001", koperasiID).First(&anggota).Error
	if err == nil {
		fmt.Println("   → Using existing test member")
		return &anggota
	}

	// Hash the PIN (123456)
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash PIN: %v", err)
	}

	// Create new test member
	tanggalLahir := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	tanggalBergabung := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	anggota = models.Anggota{
		ID:               uuid.New(),
		IDKoperasi:       koperasiID,
		NomorAnggota:     "A001",
		NamaLengkap:      "Test Member Portal",
		NIK:              "1234567890123456",
		JenisKelamin:     "L",
		TempatLahir:      "Jakarta",
		TanggalLahir:     tanggalLahir,
		Alamat:           "Jl. Test Member No. 1",
		RT:               "001",
		RW:               "002",
		Kelurahan:        "Test Kelurahan",
		Kecamatan:        "Test Kecamatan",
		KotaKabupaten:    "Jakarta",
		Provinsi:         "DKI Jakarta",
		KodePos:          "12345",
		NomorTelepon:     "081234567890",
		Email:            "test.member@email.com",
		Pekerjaan:        "Karyawan Swasta",
		TanggalBergabung: tanggalBergabung,
		Status:           models.StatusAktif,
		PIN:              string(hashedPIN),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(&anggota).Error; err != nil {
		log.Fatalf("Failed to create test member: %v", err)
	}

	return &anggota
}

// createInitialBalance creates the initial balance for test member
func createInitialBalance(db *gorm.DB, anggotaID, koperasiID uuid.UUID) *models.SaldoSimpananAnggota {
	var saldo models.SaldoSimpananAnggota

	// Try to find existing balance
	err := db.Where("id_anggota = ? AND id_koperasi = ?", anggotaID, koperasiID).First(&saldo).Error
	if err == nil {
		fmt.Println("   → Using existing balance")
		return &saldo
	}

	// Create new balance
	saldo = models.SaldoSimpananAnggota{
		ID:               uuid.New(),
		IDKoperasi:       koperasiID,
		IDAnggota:        anggotaID,
		SimpananPokok:    1000000,  // Rp 1,000,000
		SimpananWajib:    2500000,  // Rp 2,500,000
		SimpananSukarela: 500000,   // Rp 500,000
		TotalSimpanan:    4000000,  // Rp 4,000,000
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(&saldo).Error; err != nil {
		log.Fatalf("Failed to create initial balance: %v", err)
	}

	return &saldo
}

// createSampleTransactions creates sample transactions for testing
func createSampleTransactions(db *gorm.DB, anggotaID, koperasiID uuid.UUID) []models.Simpanan {
	// Check if transactions already exist
	var count int64
	db.Model(&models.Simpanan{}).Where("id_anggota = ? AND id_koperasi = ?", anggotaID, koperasiID).Count(&count)
	if count > 0 {
		fmt.Println("   → Using existing transactions")
		var existing []models.Simpanan
		db.Where("id_anggota = ? AND id_koperasi = ?", anggotaID, koperasiID).Find(&existing)
		return existing
	}

	transactions := []models.Simpanan{
		// Simpanan Pokok (initial deposit)
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SP-2024-001",
			TipeSimpanan:     models.SimpananPokok,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			Jumlah:           1000000,
			Keterangan:       "Setoran Simpanan Pokok",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		// Simpanan Wajib (monthly deposits)
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SW-2024-001",
			TipeSimpanan:     models.SimpananWajib,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Jumlah:           500000,
			Keterangan:       "Setoran Simpanan Wajib Januari 2024",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SW-2024-002",
			TipeSimpanan:     models.SimpananWajib,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC),
			Jumlah:           500000,
			Keterangan:       "Setoran Simpanan Wajib Februari 2024",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SW-2024-003",
			TipeSimpanan:     models.SimpananWajib,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC),
			Jumlah:           500000,
			Keterangan:       "Setoran Simpanan Wajib Maret 2024",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SW-2024-004",
			TipeSimpanan:     models.SimpananWajib,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 4, 15, 10, 0, 0, 0, time.UTC),
			Jumlah:           500000,
			Keterangan:       "Setoran Simpanan Wajib April 2024",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SW-2024-005",
			TipeSimpanan:     models.SimpananWajib,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC),
			Jumlah:           500000,
			Keterangan:       "Setoran Simpanan Wajib Mei 2024",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		// Simpanan Sukarela (voluntary deposits)
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SS-2024-001",
			TipeSimpanan:     models.SimpananSukarela,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC),
			Jumlah:           200000,
			Keterangan:       "Setoran Simpanan Sukarela",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:               uuid.New(),
			IDKoperasi:       koperasiID,
			IDAnggota:        anggotaID,
			NomorReferensi:   "SS-2024-002",
			TipeSimpanan:     models.SimpananSukarela,
			TipeTransaksi:    models.TipeSetoran,
			TanggalTransaksi: time.Date(2024, 3, 20, 10, 0, 0, 0, time.UTC),
			Jumlah:           300000,
			Keterangan:       "Setoran Simpanan Sukarela",
			MetodePembayaran: "tunai",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	// Insert all transactions
	for _, transaction := range transactions {
		if err := db.Create(&transaction).Error; err != nil {
			log.Fatalf("Failed to create transaction: %v", err)
		}
	}

	return transactions
}
