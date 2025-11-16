package services

import (
	"cooperative-erp-lite/internal/models"
	"testing"

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
