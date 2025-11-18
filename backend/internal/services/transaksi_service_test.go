package services

import (
	"cooperative-erp-lite/internal/models"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDBForTransaksi creates a test database connection
func setupTestDBForTransaksi(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	// Auto migrate tables for testing
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
		&models.Akun{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TestValidasiTransaksi_FloatingPointPrecision tests that floating-point precision issues don't cause false rejections
func TestValidasiTransaksi_FloatingPointPrecision(t *testing.T) {
	db := setupTestDBForTransaksi(t)
	if db == nil {
		return
	}

	service := NewTransaksiService(db)

	tests := []struct {
		name           string
		barisTransaksi []BuatBarisTransaksiRequest
		shouldPass     bool
		description    string
	}{
		{
			name: "Classic floating-point issue: 0.1 + 0.1 + 0.1 should equal 0.3",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 0.1, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0.1, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0.1, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 0.3},
			},
			shouldPass:  true,
			description: "Should pass with epsilon tolerance (0.1+0.1+0.1 = 0.30000000000000004 ≈ 0.3)",
		},
		{
			name: "Multiple small amounts that could accumulate floating-point errors",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 10.5, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 15.3, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 8.7, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 34.5},
			},
			shouldPass:  true,
			description: "Multiple decimal additions should balance within epsilon",
		},
		{
			name: "Real-world scenario: Rp 33,800 split across multiple items",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 10000.0, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 15500.0, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 8300.0, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 33800.0},
			},
			shouldPass:  true,
			description: "Real-world transaction amounts should balance exactly",
		},
		{
			name: "Edge case: Very small amounts within epsilon tolerance",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 0.01, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 0.01},
			},
			shouldPass:  true,
			description: "Minimum transaction amount (1 sen) should be valid",
		},
		{
			name: "Actually unbalanced: Off by more than epsilon (0.50)",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 100.00, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 99.50},
			},
			shouldPass:  false,
			description: "Truly unbalanced transaction should fail (off by 0.50 > epsilon 0.01)",
		},
		{
			name: "Unbalanced: Off by 0.02 (greater than epsilon)",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 100.00, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 99.98},
			},
			shouldPass:  false,
			description: "Unbalanced by 0.02 should fail (> epsilon 0.01)",
		},
		{
			name: "Edge case: Balanced within epsilon (off by 0.005)",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 100.005, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 100.00},
			},
			shouldPass:  true,
			description: "Off by 0.005 should pass (< epsilon 0.01)",
		},
		{
			name: "Zero transaction should fail",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 0},
			},
			shouldPass:  false,
			description: "Transaction with zero amounts should be rejected",
		},
		{
			name: "Complex multi-line transaction with floating-point arithmetic",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 33.33, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 33.33, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 33.34, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 100.00},
			},
			shouldPass:  true,
			description: "Splitting amounts like 100/3 should balance within epsilon",
		},
		{
			name: "Large amounts with decimal precision",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 9999999.99, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 9999999.99},
			},
			shouldPass:  true,
			description: "Large amounts should not cause precision issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidasiTransaksi(tt.barisTransaksi)

			if tt.shouldPass && err != nil {
				t.Errorf("Expected validation to PASS but got error: %v\n"+
					"Description: %s\n"+
					"This indicates floating-point precision issue not properly handled",
					err, tt.description)
			}

			if !tt.shouldPass && err == nil {
				t.Errorf("Expected validation to FAIL but it passed\n"+
					"Description: %s\n"+
					"This indicates validation is too lenient",
					tt.description)
			}
		})
	}
}

// TestValidasiTransaksi_EdgeCases tests edge cases beyond floating-point precision
func TestValidasiTransaksi_EdgeCases(t *testing.T) {
	db := setupTestDBForTransaksi(t)
	if db == nil {
		return
	}

	service := NewTransaksiService(db)

	tests := []struct {
		name           string
		barisTransaksi []BuatBarisTransaksiRequest
		shouldPass     bool
		expectedError  string
	}{
		{
			name: "Less than 2 lines should fail",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 100, JumlahKredit: 0},
			},
			shouldPass:    false,
			expectedError: "minimal 2 baris",
		},
		{
			name: "Both debit and credit on same line should fail",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 100, JumlahKredit: 50},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 100},
			},
			shouldPass:    false,
			expectedError: "tidak boleh memiliki debit dan kredit sekaligus",
		},
		{
			name: "Line with zero debit and zero credit should fail",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 100, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 0},
			},
			shouldPass:    false,
			expectedError: "harus memiliki nilai debit atau kredit",
		},
		{
			name: "Only debit lines (no credit) should fail",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 100, JumlahKredit: 0},
				{IDAkun: uuid.New(), JumlahDebit: 50, JumlahKredit: 0},
			},
			shouldPass:    false,
			expectedError: "minimal satu baris debit dan satu baris kredit",
		},
		{
			name: "Only credit lines (no debit) should fail",
			barisTransaksi: []BuatBarisTransaksiRequest{
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 100},
				{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 50},
			},
			shouldPass:    false,
			expectedError: "minimal satu baris debit dan satu baris kredit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidasiTransaksi(tt.barisTransaksi)

			if tt.shouldPass && err != nil {
				t.Errorf("Expected validation to pass but got error: %v", err)
			}

			if !tt.shouldPass && err == nil {
				t.Errorf("Expected validation to fail but it passed")
			}
		})
	}
}

// TestBeforeSaveHook_FloatingPointPrecision tests the BeforeSave hook handles floating-point correctly
func TestBeforeSaveHook_FloatingPointPrecision(t *testing.T) {
	db := setupTestDBForTransaksi(t)
	if db == nil {
		return
	}

	tests := []struct {
		name            string
		totalDebit      float64
		totalKredit     float64
		expectedBalance bool
		description     string
	}{
		{
			name:            "Exact match should be balanced",
			totalDebit:      100.0,
			totalKredit:     100.0,
			expectedBalance: true,
			description:     "Perfect equality",
		},
		{
			name:            "Floating-point precision issue: 0.3 vs 0.30000000000000004",
			totalDebit:      0.1 + 0.1 + 0.1, // = 0.30000000000000004
			totalKredit:     0.3,
			expectedBalance: true,
			description:     "Should be balanced with epsilon tolerance",
		},
		{
			name:            "Off by more than epsilon should not be balanced",
			totalDebit:      100.0,
			totalKredit:     99.5,
			expectedBalance: false,
			description:     "Off by 0.5 > epsilon 0.01",
		},
		{
			name:            "Zero amounts should not be balanced",
			totalDebit:      0.0,
			totalKredit:     0.0,
			expectedBalance: false,
			description:     "Empty transaction",
		},
		{
			name:            "Within epsilon tolerance (0.005 difference)",
			totalDebit:      100.005,
			totalKredit:     100.0,
			expectedBalance: true,
			description:     "0.005 difference < epsilon 0.01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaksi := &models.Transaksi{
				BarisTransaksi: []models.BarisTransaksi{
					{JumlahDebit: tt.totalDebit, JumlahKredit: 0},
					{JumlahDebit: 0, JumlahKredit: tt.totalKredit},
				},
			}

			// Call BeforeSave hook
			err := transaksi.BeforeSave(db)
			if err != nil {
				t.Fatalf("BeforeSave returned error: %v", err)
			}

			// Check StatusBalanced
			if transaksi.StatusBalanced != tt.expectedBalance {
				t.Errorf("Expected StatusBalanced=%v but got %v\n"+
					"Description: %s\n"+
					"TotalDebit: %.17f, TotalKredit: %.17f",
					tt.expectedBalance, transaksi.StatusBalanced,
					tt.description, tt.totalDebit, tt.totalKredit)
			}
		})
	}
}

// TestConcurrentTransactionCreation_FirstOfDay tests the CRITICAL "first transaction of day" scenario
// This is where the original SELECT FOR UPDATE approach FAILED because no rows existed to lock.
// PostgreSQL advisory locks solve this by allowing us to lock even when no rows exist.
//
// BUG SCENARIO WITHOUT ADVISORY LOCK:
// Thread A: SELECT FOR UPDATE WHERE date='2025-01-20' → ErrRecordNotFound (NO LOCK!)
// Thread B: SELECT FOR UPDATE WHERE date='2025-01-20' → ErrRecordNotFound (NO LOCK!)
// Thread A: Generate JRN-20250120-0001
// Thread B: Generate JRN-20250120-0001 (DUPLICATE!)
// Thread A: INSERT → SUCCESS
// Thread B: INSERT → DUPLICATE KEY ERROR ❌
//
// FIX WITH ADVISORY LOCK:
// Thread A: pg_advisory_xact_lock(hash) → ACQUIRED
// Thread B: pg_advisory_xact_lock(hash) → WAITING...
// Thread A: Generate JRN-20250120-0001, INSERT, COMMIT → LOCK RELEASED
// Thread B: Lock acquired, Generate JRN-20250120-0002, INSERT, COMMIT ✅
func TestConcurrentTransactionCreation_FirstOfDay(t *testing.T) {
	db := setupTestDBForTransaksi(t)
	if db == nil {
		return
	}

	// Clean up to simulate "first transaction of the day"
	db.Exec("DELETE FROM baris_transaksi")
	db.Exec("DELETE FROM transaksi")
	db.Exec("DELETE FROM akun")
	db.Exec("DELETE FROM koperasi")

	// Create test cooperative
	koperasi := models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@example.com",
		NoTelepon:    "08123456789",
	}
	if err := db.Create(&koperasi).Error; err != nil {
		t.Fatalf("Failed to create test cooperative: %v", err)
	}

	// Create test accounts
	akunKas := models.Akun{
		ID:          uuid.New(),
		IDKoperasi:  koperasi.ID,
		KodeAkun:    "1101",
		NamaAkun:    "Kas",
		TipeAkun:    models.AkunAset,
		NormalSaldo: "debit",
	}
	akunModal := models.Akun{
		ID:          uuid.New(),
		IDKoperasi:  koperasi.ID,
		KodeAkun:    "3101",
		NamaAkun:    "Simpanan Pokok",
		TipeAkun:    models.AkunModal,
		NormalSaldo: "kredit",
	}
	db.Create(&akunKas)
	db.Create(&akunModal)

	service := NewTransaksiService(db)
	idPengguna := uuid.New()

	// CRITICAL: 10 concurrent "first transactions" - this would FAIL without advisory lock
	concurrentCount := 10
	errChan := make(chan error, concurrentCount)
	doneChan := make(chan bool, concurrentCount)

	for i := 0; i < concurrentCount; i++ {
		go func(index int) {
			req := &BuatTransaksiRequest{
				TanggalTransaksi: mustParseTime("2025-01-20"),
				Deskripsi:        fmt.Sprintf("First transaction #%d", index),
				TipeTransaksi:    "test",
				BarisTransaksi: []BuatBarisTransaksiRequest{
					{IDAkun: akunKas.ID, JumlahDebit: 100000.0},
					{IDAkun: akunModal.ID, JumlahKredit: 100000.0},
				},
			}

			_, err := service.BuatTransaksi(koperasi.ID, idPengguna, req)
			if err != nil {
				errChan <- err
			} else {
				doneChan <- true
			}
		}(i)
	}

	// Collect results
	var errors []error
	successCount := 0
	for i := 0; i < concurrentCount; i++ {
		select {
		case err := <-errChan:
			errors = append(errors, err)
		case <-doneChan:
			successCount++
		}
	}

	// CRITICAL: All must succeed (no duplicate key errors)
	if len(errors) > 0 {
		t.Errorf("FAILED 'first-of-day' test: %d errors:", len(errors))
		for i, err := range errors {
			t.Errorf("  %d: %v", i+1, err)
		}
		t.Fatal("Advisory lock FAILED to prevent race condition")
	}

	// Verify sequential journal numbers
	var transaksiList []models.Transaksi
	db.Where("id_koperasi = ? AND DATE(tanggal_transaksi) = ?", koperasi.ID, "2025-01-20").
		Order("nomor_jurnal ASC").Find(&transaksiList)

	if len(transaksiList) != concurrentCount {
		t.Fatalf("Expected %d transactions, got %d", concurrentCount, len(transaksiList))
	}

	for i, transaksi := range transaksiList {
		expected := fmt.Sprintf("JRN-20250120-%04d", i+1)
		if transaksi.NomorJurnal != expected {
			t.Errorf("Expected %s, got %s", expected, transaksi.NomorJurnal)
		}
	}

	t.Logf("✅ SUCCESS: All %d concurrent first-of-day transactions succeeded", concurrentCount)

	// Clean up
	db.Exec("DELETE FROM baris_transaksi")
	db.Exec("DELETE FROM transaksi")
	db.Exec("DELETE FROM akun WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM koperasi WHERE id = ?", koperasi.ID)
}

// TestConcurrentTransactionCreation_RaceConditionFix tests that the race condition
// in GenerateNomorJurnal is fixed when creating multiple transactions concurrently
func TestConcurrentTransactionCreation_RaceConditionFix(t *testing.T) {
	db := setupTestDBForTransaksi(t)
	if db == nil {
		return
	}

	// Clean up any existing data
	db.Exec("DELETE FROM baris_transaksi")
	db.Exec("DELETE FROM transaksi")
	db.Exec("DELETE FROM akun")
	db.Exec("DELETE FROM koperasi")

	// Create test cooperative
	koperasi := models.Koperasi{
		ID:           uuid.New(),
		NamaKoperasi: "Test Koperasi",
		Email:        "test@example.com",
		NoTelepon:    "08123456789",
	}
	if err := db.Create(&koperasi).Error; err != nil {
		t.Fatalf("Failed to create test cooperative: %v", err)
	}

	// Create test accounts (Kas and Modal)
	akunKas := models.Akun{
		ID:          uuid.New(),
		IDKoperasi:  koperasi.ID,
		KodeAkun:    "1101",
		NamaAkun:    "Kas",
		TipeAkun:    models.AkunAset,
		NormalSaldo: "debit",
	}
	akunModal := models.Akun{
		ID:          uuid.New(),
		IDKoperasi:  koperasi.ID,
		KodeAkun:    "3101",
		NamaAkun:    "Simpanan Pokok",
		TipeAkun:    models.AkunModal,
		NormalSaldo: "kredit",
	}
	if err := db.Create(&akunKas).Error; err != nil {
		t.Fatalf("Failed to create kas account: %v", err)
	}
	if err := db.Create(&akunModal).Error; err != nil {
		t.Fatalf("Failed to create modal account: %v", err)
	}

	service := NewTransaksiService(db)
	idPengguna := uuid.New()

	// Number of concurrent transactions to create
	concurrentCount := 10

	// Use channels to coordinate goroutines
	errChan := make(chan error, concurrentCount)
	doneChan := make(chan bool, concurrentCount)

	// Create concurrent transactions on the same date
	for i := 0; i < concurrentCount; i++ {
		go func(index int) {
			req := &BuatTransaksiRequest{
				TanggalTransaksi: mustParseTime("2025-01-16"),
				Deskripsi:        "Concurrent test transaction",
				NomorReferensi:   "",
				TipeTransaksi:    "test",
				BarisTransaksi: []BuatBarisTransaksiRequest{
					{
						IDAkun:      akunKas.ID,
						JumlahDebit: 100000.0,
						Keterangan:  "Test debit",
					},
					{
						IDAkun:       akunModal.ID,
						JumlahKredit: 100000.0,
						Keterangan:   "Test kredit",
					},
				},
			}

			_, err := service.BuatTransaksi(koperasi.ID, idPengguna, req)
			if err != nil {
				errChan <- err
			} else {
				doneChan <- true
			}
		}(i)
	}

	// Collect results
	var errors []error
	successCount := 0

	for i := 0; i < concurrentCount; i++ {
		select {
		case err := <-errChan:
			errors = append(errors, err)
		case <-doneChan:
			successCount++
		}
	}

	// Verify that all transactions were created successfully
	if len(errors) > 0 {
		t.Errorf("Expected all %d concurrent transactions to succeed, but %d failed with errors:",
			concurrentCount, len(errors))
		for i, err := range errors {
			t.Errorf("  Error %d: %v", i+1, err)
		}
	}

	if successCount != concurrentCount {
		t.Errorf("Expected %d successful transactions, got %d", concurrentCount, successCount)
	}

	// Verify that journal numbers are sequential and unique
	var transaksiList []models.Transaksi
	db.Where("id_koperasi = ? AND DATE(tanggal_transaksi) = ?", koperasi.ID, "2025-01-16").
		Order("nomor_jurnal ASC").
		Find(&transaksiList)

	if len(transaksiList) != concurrentCount {
		t.Errorf("Expected %d transactions in database, found %d", concurrentCount, len(transaksiList))
	}

	// Check that journal numbers are sequential: JRN-20250116-0001, JRN-20250116-0002, etc.
	expectedNomor := make(map[string]bool)
	for i := 1; i <= concurrentCount; i++ {
		expectedNomor[mustFormatJournalNumber(i)] = false
	}

	for _, transaksi := range transaksiList {
		if _, exists := expectedNomor[transaksi.NomorJurnal]; exists {
			expectedNomor[transaksi.NomorJurnal] = true
		} else {
			t.Errorf("Unexpected journal number: %s", transaksi.NomorJurnal)
		}
	}

	// Verify all expected journal numbers were found
	for nomor, found := range expectedNomor {
		if !found {
			t.Errorf("Missing expected journal number: %s", nomor)
		}
	}

	// Clean up
	db.Exec("DELETE FROM baris_transaksi")
	db.Exec("DELETE FROM transaksi")
	db.Exec("DELETE FROM akun WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM koperasi WHERE id = ?", koperasi.ID)
}

// Helper function to parse time for tests
func mustParseTime(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return t
}

// Helper function to format expected journal number
func mustFormatJournalNumber(seq int) string {
	return fmt.Sprintf("JRN-20250116-%04d", seq)
}

// BenchmarkValidasiTransaksi benchmarks the validation performance
func BenchmarkValidasiTransaksi(b *testing.B) {
	db, _ := gorm.Open(postgres.Open(""), &gorm.Config{})
	service := NewTransaksiService(db)

	barisTransaksi := []BuatBarisTransaksiRequest{
		{IDAkun: uuid.New(), JumlahDebit: 0.1, JumlahKredit: 0},
		{IDAkun: uuid.New(), JumlahDebit: 0.1, JumlahKredit: 0},
		{IDAkun: uuid.New(), JumlahDebit: 0.1, JumlahKredit: 0},
		{IDAkun: uuid.New(), JumlahDebit: 0, JumlahKredit: 0.3},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.ValidasiTransaksi(barisTransaksi)
	}
}
