package services

import (
	"cooperative-erp-lite/internal/models"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupLaporanBenchmarkDB creates a test database for laporan benchmarks
func setupLaporanBenchmarkDB(b *testing.B) (*gorm.DB, uuid.UUID) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database: %v", err)
		return nil, uuid.Nil
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Akun{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
	)
	if err != nil {
		b.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Benchmark Laporan Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	return db, koperasi.ID
}

// seedAccounts creates test chart of accounts
func seedAccounts(db *gorm.DB, idKoperasi uuid.UUID, count int) []models.Akun {
	// Clear existing data
	db.Exec("DELETE FROM baris_transaksi WHERE id_akun IN (SELECT id FROM akun WHERE id_koperasi = ?)", idKoperasi)
	db.Exec("DELETE FROM akun WHERE id_koperasi = ?", idKoperasi)

	accounts := make([]models.Akun, 0, count)

	// Create different types of accounts
	accountTypes := []struct {
		tipe   models.TipeAkun
		prefix string
		normal string
	}{
		{models.AkunAset, "1", "debit"},
		{models.AkunKewajiban, "2", "kredit"},
		{models.AkunModal, "3", "kredit"},
		{models.AkunPendapatan, "4", "kredit"},
		{models.AkunBeban, "5", "debit"},
	}

	for i := 0; i < count; i++ {
		typeIdx := i % len(accountTypes)
		accType := accountTypes[typeIdx]

		account := models.Akun{
			IDKoperasi:  idKoperasi,
			KodeAkun:    fmt.Sprintf("%s%03d", accType.prefix, i),
			NamaAkun:    fmt.Sprintf("%s Account %d", accType.tipe, i),
			TipeAkun:    accType.tipe,
			NormalSaldo: accType.normal,
		}
		accounts = append(accounts, account)
	}

	db.CreateInBatches(accounts, 100)
	return accounts
}

// seedTransactions creates test transactions with line items
func seedTransactions(db *gorm.DB, idKoperasi uuid.UUID, accounts []models.Akun, transactionCount int) {
	// Clear existing transactions
	db.Exec("DELETE FROM baris_transaksi WHERE id_transaksi IN (SELECT id FROM transaksi WHERE id_koperasi = ?)", idKoperasi)
	db.Exec("DELETE FROM transaksi WHERE id_koperasi = ?", idKoperasi)

	transactions := make([]models.Transaksi, transactionCount)
	for i := 0; i < transactionCount; i++ {
		transactions[i] = models.Transaksi{
			IDKoperasi:       idKoperasi,
			NomorTransaksi:   fmt.Sprintf("TRX-%04d", i+1),
			TanggalTransaksi: time.Now().AddDate(0, 0, -(i % 30)), // Spread over 30 days
			Keterangan:       fmt.Sprintf("Test Transaction %d", i+1),
			TotalDebit:       100000,
			TotalKredit:      100000,
		}
	}
	db.CreateInBatches(transactions, 100)

	// Create line items for each transaction
	lineItems := make([]models.BarisTransaksi, 0, transactionCount*2)
	for i, txn := range transactions {
		// Debit entry
		debitAccount := accounts[i%len(accounts)]
		lineItems = append(lineItems, models.BarisTransaksi{
			IDTransaksi: txn.ID,
			IDAkun:      debitAccount.ID,
			JumlahDebit: 100000,
			JumlahKredit: 0,
		})

		// Credit entry
		creditAccount := accounts[(i+1)%len(accounts)]
		lineItems = append(lineItems, models.BarisTransaksi{
			IDTransaksi: txn.ID,
			IDAkun:      creditAccount.ID,
			JumlahDebit: 0,
			JumlahKredit: 100000,
		})
	}
	db.CreateInBatches(lineItems, 100)
}

// BenchmarkGenerateLaporanPosisiKeuangan_100Accounts benchmarks balance sheet with 100 accounts
func BenchmarkGenerateLaporanPosisiKeuangan_100Accounts(b *testing.B) {
	db, idKoperasi := setupLaporanBenchmarkDB(b)
	if db == nil {
		return
	}

	accounts := seedAccounts(db, idKoperasi, 100)
	seedTransactions(db, idKoperasi, accounts, 500)

	akunService := NewAkunService(db)
	service := NewLaporanService(db, akunService, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateLaporanPosisiKeuangan(idKoperasi, "")
		if err != nil {
			b.Fatalf("Error generating report: %v", err)
		}
	}
}

// BenchmarkGenerateLaporanPosisiKeuangan_500Accounts benchmarks balance sheet with 500 accounts
// Target: < 200ms per operation
func BenchmarkGenerateLaporanPosisiKeuangan_500Accounts(b *testing.B) {
	db, idKoperasi := setupLaporanBenchmarkDB(b)
	if db == nil {
		return
	}

	accounts := seedAccounts(db, idKoperasi, 500)
	seedTransactions(db, idKoperasi, accounts, 2000)

	akunService := NewAkunService(db)
	service := NewLaporanService(db, akunService, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateLaporanPosisiKeuangan(idKoperasi, "")
		if err != nil {
			b.Fatalf("Error generating report: %v", err)
		}
	}
}

// BenchmarkGenerateLaporanLabaRugi_100Accounts benchmarks income statement
func BenchmarkGenerateLaporanLabaRugi_100Accounts(b *testing.B) {
	db, idKoperasi := setupLaporanBenchmarkDB(b)
	if db == nil {
		return
	}

	accounts := seedAccounts(db, idKoperasi, 100)
	seedTransactions(db, idKoperasi, accounts, 500)

	akunService := NewAkunService(db)
	service := NewLaporanService(db, akunService, nil, nil)

	tanggalMulai := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	tanggalAkhir := time.Now().Format("2006-01-02")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateLaporanLabaRugi(idKoperasi, tanggalMulai, tanggalAkhir)
		if err != nil {
			b.Fatalf("Error generating report: %v", err)
		}
	}
}

// BenchmarkGenerateLaporanLabaRugi_500Accounts benchmarks income statement with 500 accounts
// Target: < 200ms per operation
func BenchmarkGenerateLaporanLabaRugi_500Accounts(b *testing.B) {
	db, idKoperasi := setupLaporanBenchmarkDB(b)
	if db == nil {
		return
	}

	accounts := seedAccounts(db, idKoperasi, 500)
	seedTransactions(db, idKoperasi, accounts, 2000)

	akunService := NewAkunService(db)
	service := NewLaporanService(db, akunService, nil, nil)

	tanggalMulai := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	tanggalAkhir := time.Now().Format("2006-01-02")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateLaporanLabaRugi(idKoperasi, tanggalMulai, tanggalAkhir)
		if err != nil {
			b.Fatalf("Error generating report: %v", err)
		}
	}
}
