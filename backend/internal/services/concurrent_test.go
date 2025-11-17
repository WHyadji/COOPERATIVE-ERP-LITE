package services

import (
	"cooperative-erp-lite/internal/models"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *gorm.DB {
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
		&models.Transaksi{},
		&models.BarisTransaksi{},
		&models.Simpanan{},
		&models.Penjualan{},
		&models.ItemPenjualan{},
		&models.Produk{},
		&models.Akun{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// cleanupTestData removes test data
func cleanupTestData(db *gorm.DB, koperasiID uuid.UUID) {
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Anggota{})
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Transaksi{})
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Simpanan{})
	db.Unscoped().Where("id_koperasi = ?", koperasiID).Delete(&models.Penjualan{})
	db.Unscoped().Where("id = ?", koperasiID).Delete(&models.Koperasi{})
}

// TestGenerateNomorAnggota_Concurrent tests concurrent member number generation
func TestGenerateNomorAnggota_Concurrent(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)
	defer cleanupTestData(db, koperasiID)

	service := NewAnggotaService(db)

	const numGoroutines = 100
	results := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)
	var wg sync.WaitGroup

	// Launch 100 concurrent requests
	startTime := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nomor, err := service.GenerateNomorAnggota(koperasiID)
			if err != nil {
				errors <- err
				return
			}
			results <- nomor
		}()
	}

	wg.Wait()
	close(results)
	close(errors)
	duration := time.Since(startTime)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Error generating nomor: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Fatalf("Got %d errors during concurrent generation", errorCount)
	}

	// Collect all results and check for duplicates
	generated := make(map[string]bool)
	var nomorList []string
	for nomor := range results {
		if generated[nomor] {
			t.Errorf("Duplicate number generated: %s", nomor)
		}
		generated[nomor] = true
		nomorList = append(nomorList, nomor)
	}

	// Verify all numbers are unique
	if len(generated) != numGoroutines {
		t.Errorf("Expected %d unique numbers, got %d", numGoroutines, len(generated))
	}

	// Performance check: should complete in < 10 seconds for 100 concurrent requests
	if duration > 10*time.Second {
		t.Errorf("Performance issue: took %v for %d concurrent requests (expected < 10s)", duration, numGoroutines)
	}

	t.Logf("✓ Successfully generated %d unique member numbers in %v", len(generated), duration)
	t.Logf("✓ Average time per request: %v", duration/numGoroutines)
}

// TestGenerateNomorJurnal_Concurrent tests concurrent journal number generation
func TestGenerateNomorJurnal_Concurrent(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)
	defer cleanupTestData(db, koperasiID)

	service := NewTransaksiService(db)
	tanggal := time.Now()

	const numGoroutines = 100
	results := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)
	var wg sync.WaitGroup

	// Launch 100 concurrent requests
	startTime := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nomor, err := service.GenerateNomorJurnal(koperasiID, tanggal)
			if err != nil {
				errors <- err
				return
			}
			results <- nomor
		}()
	}

	wg.Wait()
	close(results)
	close(errors)
	duration := time.Since(startTime)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Error generating nomor: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Fatalf("Got %d errors during concurrent generation", errorCount)
	}

	// Collect all results and check for duplicates
	generated := make(map[string]bool)
	for nomor := range results {
		if generated[nomor] {
			t.Errorf("Duplicate number generated: %s", nomor)
		}
		generated[nomor] = true
	}

	// Verify all numbers are unique
	if len(generated) != numGoroutines {
		t.Errorf("Expected %d unique numbers, got %d", numGoroutines, len(generated))
	}

	// Performance check
	if duration > 10*time.Second {
		t.Errorf("Performance issue: took %v for %d concurrent requests (expected < 10s)", duration, numGoroutines)
	}

	t.Logf("✓ Successfully generated %d unique journal numbers in %v", len(generated), duration)
	t.Logf("✓ Average time per request: %v", duration/numGoroutines)
}

// TestGenerateNomorReferensi_Concurrent tests concurrent deposit reference number generation
func TestGenerateNomorReferensi_Concurrent(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)
	defer cleanupTestData(db, koperasiID)

	service := NewSimpananService(db, nil)
	tanggal := time.Now()

	const numGoroutines = 100
	results := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)
	var wg sync.WaitGroup

	// Launch 100 concurrent requests
	startTime := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nomor, err := service.GenerateNomorReferensi(koperasiID, tanggal)
			if err != nil {
				errors <- err
				return
			}
			results <- nomor
		}()
	}

	wg.Wait()
	close(results)
	close(errors)
	duration := time.Since(startTime)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Error generating nomor: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Fatalf("Got %d errors during concurrent generation", errorCount)
	}

	// Collect all results and check for duplicates
	generated := make(map[string]bool)
	for nomor := range results {
		if generated[nomor] {
			t.Errorf("Duplicate number generated: %s", nomor)
		}
		generated[nomor] = true
	}

	// Verify all numbers are unique
	if len(generated) != numGoroutines {
		t.Errorf("Expected %d unique numbers, got %d", numGoroutines, len(generated))
	}

	// Performance check
	if duration > 10*time.Second {
		t.Errorf("Performance issue: took %v for %d concurrent requests (expected < 10s)", duration, numGoroutines)
	}

	t.Logf("✓ Successfully generated %d unique deposit reference numbers in %v", len(generated), duration)
	t.Logf("✓ Average time per request: %v", duration/numGoroutines)
}

// TestGenerateNomorPenjualan_Concurrent tests concurrent sales number generation
func TestGenerateNomorPenjualan_Concurrent(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)
	defer cleanupTestData(db, koperasiID)

	service := NewPenjualanService(db, nil, nil)
	tanggal := time.Now()

	const numGoroutines = 100
	results := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)
	var wg sync.WaitGroup

	// Launch 100 concurrent requests
	startTime := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nomor, err := service.GenerateNomorPenjualan(koperasiID, tanggal)
			if err != nil {
				errors <- err
				return
			}
			results <- nomor
		}()
	}

	wg.Wait()
	close(results)
	close(errors)
	duration := time.Since(startTime)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Error generating nomor: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Fatalf("Got %d errors during concurrent generation", errorCount)
	}

	// Collect all results and check for duplicates
	generated := make(map[string]bool)
	for nomor := range results {
		if generated[nomor] {
			t.Errorf("Duplicate number generated: %s", nomor)
		}
		generated[nomor] = true
	}

	// Verify all numbers are unique
	if len(generated) != numGoroutines {
		t.Errorf("Expected %d unique numbers, got %d", numGoroutines, len(generated))
	}

	// Performance check
	if duration > 10*time.Second {
		t.Errorf("Performance issue: took %v for %d concurrent requests (expected < 10s)", duration, numGoroutines)
	}

	t.Logf("✓ Successfully generated %d unique sales numbers in %v", len(generated), duration)
	t.Logf("✓ Average time per request: %v", duration/numGoroutines)
}

// BenchmarkGenerateNomorAnggota benchmarks member number generation
func BenchmarkGenerateNomorAnggota(b *testing.B) {
	db := setupTestDB(&testing.T{})
	if db == nil {
		b.Skip("Database not available")
		return
	}

	koperasiID := uuid.New()
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)
	defer cleanupTestData(db, koperasiID)

	service := NewAnggotaService(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateNomorAnggota(koperasiID)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

// TestAllNumberGeneration_Integration tests all number generators work correctly together
func TestAllNumberGeneration_Integration(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)
	defer cleanupTestData(db, koperasiID)

	// Test that all services can generate numbers concurrently without conflicts
	var wg sync.WaitGroup
	errors := make(chan error, 4)

	// Test member numbers
	wg.Add(1)
	go func() {
		defer wg.Done()
		service := NewAnggotaService(db)
		for i := 0; i < 10; i++ {
			_, err := service.GenerateNomorAnggota(koperasiID)
			if err != nil {
				errors <- fmt.Errorf("anggota error: %w", err)
				return
			}
		}
	}()

	// Test journal numbers
	wg.Add(1)
	go func() {
		defer wg.Done()
		service := NewTransaksiService(db)
		for i := 0; i < 10; i++ {
			_, err := service.GenerateNomorJurnal(koperasiID, time.Now())
			if err != nil {
				errors <- fmt.Errorf("jurnal error: %w", err)
				return
			}
		}
	}()

	// Test deposit numbers
	wg.Add(1)
	go func() {
		defer wg.Done()
		service := NewSimpananService(db, nil)
		for i := 0; i < 10; i++ {
			_, err := service.GenerateNomorReferensi(koperasiID, time.Now())
			if err != nil {
				errors <- fmt.Errorf("simpanan error: %w", err)
				return
			}
		}
	}()

	// Test sales numbers
	wg.Add(1)
	go func() {
		defer wg.Done()
		service := NewPenjualanService(db, nil, nil)
		for i := 0; i < 10; i++ {
			_, err := service.GenerateNomorPenjualan(koperasiID, time.Now())
			if err != nil {
				errors <- fmt.Errorf("penjualan error: %w", err)
				return
			}
		}
	}()

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Integration test error: %v", err)
	}

	t.Logf("✓ All number generators working correctly together")
}
