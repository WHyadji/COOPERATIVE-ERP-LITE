package integration

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestSimpananIntegration_TransactionCycle tests complete simpanan transaction flow
func TestSimpananIntegration_TransactionCycle(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	// Setup koperasi
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi Simpanan",
		NoTelepon:    "08111111115",
		Email:        "simpanan@test.com",
	}
	if err := db.Create(koperasi).Error; err != nil {
		t.Fatalf("Failed to create koperasi: %v", err)
	}

	// Setup Chart of Accounts (required for journal posting)
	setupChartOfAccounts(db, koperasiID)

	// Create member
	anggota := &models.Anggota{
		IDKoperasi:    koperasiID,
		NomorAnggota:  "SIM001",
		NamaLengkap:  "Test Member",
		NIK:           "1234567890123456",
		TempatLahir:   "Jakarta",
		TanggalLahir:  nil,
		JenisKelamin:  "L",
		Alamat:        "Test",
		RT:            "001",
		RW:            "002",
		Kelurahan:     "Test",
		Kecamatan:     "Test",
		KotaKabupaten: "Jakarta",
		Provinsi:      "DKI",
		KodePos:       "12345",
		NoTelepon:          "08123456789",
		Email:         "member@test.com",
		Status: models.StatusAktif,
	}
	if err := db.Create(anggota).Error; err != nil {
		t.Fatalf("Failed to create member: %v", err)
	}

	// Create services
	transaksiService := services.NewTransaksiService(db)
	simpananService := services.NewSimpananService(db, transaksiService)

	// Step 1: Record Simpanan Pokok (Principal Share)
	pokokReq := &services.CatatSetoranRequest{
		IDAnggota:         anggota.ID,
		TipeSimpanan:     models.SimpananPokok,
		JumlahSetoran:    100000,
		TanggalTransaksi:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Keterangan:        "Simpanan Pokok Awal",
	}

	pokokSimpanan, err := simpananService.CatatSetoran(koperasiID, uuid.Nil, pokokReq)
	if err != nil {
		t.Fatalf("Failed to record simpanan pokok: %v", err)
	}

	if pokokSimpanan.JumlahSetoran != 100000 {
		t.Errorf("Expected amount 100000, got %.0f", pokokSimpanan.JumlahSetoran)
	}

	t.Logf("✓ Simpanan Pokok recorded: Rp %.0f", pokokSimpanan.JumlahSetoran)

	// Step 2: Record Simpanan Wajib (Mandatory Share) - Month 1
	wajibReq1 := &services.CatatSetoranRequest{
		IDAnggota:         anggota.ID,
		TipeSimpanan:     models.SimpananWajib,
		JumlahSetoran:    50000,
		TanggalTransaksi:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Keterangan:        "Simpanan Wajib Januari",
	}

	wajibSimpanan1, err := simpananService.CatatSetoran(koperasiID, uuid.Nil, wajibReq1)
	if err != nil {
		t.Fatalf("Failed to record simpanan wajib: %v", err)
	}

	t.Logf("✓ Simpanan Wajib Januari recorded: Rp %.0f", wajibSimpanan1.JumlahSetoran)

	// Step 3: Record Simpanan Wajib - Month 2
	wajibReq2 := &services.CatatSetoranRequest{
		IDAnggota:         anggota.ID,
		TipeSimpanan:     models.SimpananWajib,
		JumlahSetoran:    50000,
		TanggalTransaksi:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Keterangan:        "Simpanan Wajib Februari",
	}

	_, err = simpananService.CatatSetoran(koperasiID, uuid.Nil, wajibReq2)
	if err != nil {
		t.Fatalf("Failed to record simpanan wajib month 2: %v", err)
	}

	t.Logf("✓ Simpanan Wajib Februari recorded")

	// Step 4: Record Simpanan Sukarela (Voluntary Share)
	sukarelaReq := &services.CatatSetoranRequest{
		IDAnggota:         anggota.ID,
		TipeSimpanan:     models.SimpananSukarela,
		JumlahSetoran:    75000,
		TanggalTransaksi:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Keterangan:        "Simpanan Sukarela",
	}

	_, err = simpananService.CatatSetoran(koperasiID, uuid.Nil, sukarelaReq)
	if err != nil {
		t.Fatalf("Failed to record simpanan sukarela: %v", err)
	}

	t.Logf("✓ Simpanan Sukarela recorded")

	// Step 5: Get member balance summary
	saldo, err := simpananService.DapatkanSaldoAnggota(anggota.ID)
	if err != nil {
		t.Fatalf("Failed to get member balance: %v", err)
	}

	expectedTotal := float64(100000 + 50000 + 50000 + 75000) // 275000
	if saldo.TotalSimpanan != expectedTotal {
		t.Errorf("Expected total %.0f, got %.0f", expectedTotal, saldo.TotalSimpanan)
	}
	if saldo.SimpananPokok != 100000 {
		t.Errorf("Expected pokok 100000, got %.0f", saldo.SimpananPokok)
	}
	if saldo.SimpananWajib != 100000 {
		t.Errorf("Expected wajib 100000, got %.0f", saldo.SimpananWajib)
	}
	if saldo.SimpananSukarela != 75000 {
		t.Errorf("Expected sukarela 75000, got %.0f", saldo.SimpananSukarela)
	}

	t.Logf("✓ Member balance calculated correctly:")
	t.Logf("  - Pokok: Rp %.0f", saldo.SimpananPokok)
	t.Logf("  - Wajib: Rp %.0f", saldo.SimpananWajib)
	t.Logf("  - Sukarela: Rp %.0f", saldo.SimpananSukarela)
	t.Logf("  - Total: Rp %.0f", saldo.TotalSimpanan)

	// Step 6: Get transaction history
	riwayat, _, err := simpananService.DapatkanSemuaTransaksiSimpanan(koperasiID, "", &anggota.ID, "", "", 1, 10)
	if err != nil {
		t.Fatalf("Failed to get transaction history: %v", err)
	}

	if len(riwayat) != 4 {
		t.Errorf("Expected 4 transactions, got %d", len(riwayat))
	}

	t.Logf("✓ Transaction history retrieved: %d transactions", len(riwayat))

	t.Log("✅ Complete simpanan transaction cycle integration test passed")
}

// TestSimpananIntegration_BalanceReporting tests balance reporting across members
func TestSimpananIntegration_BalanceReporting(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)

	// Setup Chart of Accounts (required for journal posting)
	setupChartOfAccounts(db, koperasiID)

	// Create multiple members
	members := []struct {
		nomor string
		nama  string
		pokok int64
		wajib int64
	}{
		{"M001", "Member 1", 100000, 50000},
		{"M002", "Member 2", 100000, 75000},
		{"M003", "Member 3", 150000, 60000},
	}

	transaksiService := services.NewTransaksiService(db)
	simpananService := services.NewSimpananService(db, transaksiService)

	for _, m := range members {
		// Create member
		anggota := &models.Anggota{
			IDKoperasi:    koperasiID,
			NomorAnggota:  m.nomor,
			NamaLengkap:  m.nama,
			NIK:           "1234567890123456",
			TempatLahir:   "Jakarta",
			TanggalLahir:  nil,
			JenisKelamin:  "L",
			Alamat:        "Test",
			RT:            "001",
			RW:            "002",
			Kelurahan:     "Test",
			Kecamatan:     "Test",
			KotaKabupaten: "Jakarta",
			Provinsi:      "DKI",
			KodePos:       "12345",
			NoTelepon:          "08123456789",
			Email:         m.nama + "@test.com",
			Status: models.StatusAktif,
		}
		db.Create(anggota)

		// Record simpanan pokok
		pokokReq := &services.CatatSetoranRequest{
			IDAnggota:         anggota.ID,
			TipeSimpanan:     models.SimpananPokok,
			JumlahSetoran:    float64(m.pokok),
			TanggalTransaksi:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Keterangan:        "Simpanan Pokok",
		}
		simpananService.CatatSetoran(koperasiID, uuid.Nil, pokokReq)

		// Record simpanan wajib
		wajibReq := &services.CatatSetoranRequest{
			IDAnggota:         anggota.ID,
			TipeSimpanan:     models.SimpananWajib,
			JumlahSetoran:    float64(m.wajib),
			TanggalTransaksi:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Keterangan:        "Simpanan Wajib",
		}
		simpananService.CatatSetoran(koperasiID, uuid.Nil, wajibReq)
	}

	// Get balance report for all members
	laporanSaldo, err := simpananService.DapatkanLaporanSaldoAnggota(koperasiID)
	if err != nil {
		t.Fatalf("Failed to get balance report: %v", err)
	}

	if len(laporanSaldo) != 3 {
		t.Errorf("Expected 3 members in report, got %d", len(laporanSaldo))
	}

	var totalPokok, totalWajib, grandTotal float64
	for _, saldo := range laporanSaldo {
		totalPokok += saldo.SimpananPokok
		totalWajib += saldo.SimpananWajib
		grandTotal += saldo.TotalSimpanan

		t.Logf("  %s: Pokok=Rp%.0f, Wajib=Rp%.0f, Total=Rp%.0f",
			saldo.NamaAnggota, saldo.SimpananPokok, saldo.SimpananWajib, saldo.TotalSimpanan)
	}

	expectedTotal := float64((100000+100000+150000) + (50000+75000+60000))
	if grandTotal != expectedTotal {
		t.Errorf("Expected grand total %.0f, got %.0f", expectedTotal, grandTotal)
	}

	t.Logf("✓ Balance report totals:")
	t.Logf("  - Total Pokok: Rp %.0f", totalPokok)
	t.Logf("  - Total Wajib: Rp %.0f", totalWajib)
	t.Logf("  - Grand Total: Rp %.0f", grandTotal)

	t.Log("✅ Balance reporting integration test passed")
}

// TestSimpananIntegration_FilterByType tests filtering transactions by type
func TestSimpananIntegration_FilterByType(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)

	// Setup Chart of Accounts (required for journal posting)
	setupChartOfAccounts(db, koperasiID)

	anggota := &models.Anggota{
		IDKoperasi:    koperasiID,
		NomorAnggota:  "F001",
		NamaLengkap:  "Filter Test Member",
		NIK:           "1234567890123456",
		TempatLahir:   "Jakarta",
		TanggalLahir:  nil,
		JenisKelamin:  "L",
		Alamat:        "Test",
		RT:            "001",
		RW:            "002",
		Kelurahan:     "Test",
		Kecamatan:     "Test",
		KotaKabupaten: "Jakarta",
		Provinsi:      "DKI",
		KodePos:       "12345",
		NoTelepon:          "08123456789",
		Email:         "filter@test.com",
		Status: models.StatusAktif,
	}
	db.Create(anggota)

	transaksiService := services.NewTransaksiService(db)
	simpananService := services.NewSimpananService(db, transaksiService)

	// Create transactions of different types
	types := []models.TipeSimpanan{
		models.SimpananPokok,
		models.SimpananWajib,
		models.SimpananWajib,
		models.SimpananSukarela,
		models.SimpananSukarela,
		models.SimpananSukarela,
	}

	for i, jenis := range types {
		req := &services.CatatSetoranRequest{
			IDAnggota:         anggota.ID,
			TipeSimpanan:     jenis,
			JumlahSetoran:    float64(10000 * (i + 1)),
			TanggalTransaksi:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Keterangan:        "Test " + string(jenis),
		}
		simpananService.CatatSetoran(koperasiID, uuid.Nil, req)
	}

	// Filter by Simpanan Pokok
	pokokOnly, total, err := simpananService.DapatkanSemuaTransaksiSimpanan(koperasiID, string(models.SimpananPokok), &anggota.ID, "", "", 1, 10)
	if err != nil {
		t.Fatalf("Failed to filter by pokok: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected 1 pokok transaction, got %d", total)
	}

	t.Logf("✓ Filter pokok: %d transactions", len(pokokOnly))

	// Filter by Simpanan Wajib
	wajibOnly, total, err := simpananService.DapatkanSemuaTransaksiSimpanan(koperasiID, string(models.SimpananWajib), &anggota.ID, "", "", 1, 10)
	if err != nil {
		t.Fatalf("Failed to filter by wajib: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected 2 wajib transactions, got %d", total)
	}

	t.Logf("✓ Filter wajib: %d transactions", len(wajibOnly))

	// Filter by Simpanan Sukarela
	sukarelaOnly, total, err := simpananService.DapatkanSemuaTransaksiSimpanan(koperasiID, string(models.SimpananSukarela), &anggota.ID, "", "", 1, 10)
	if err != nil {
		t.Fatalf("Failed to filter by sukarela: %v", err)
	}

	if total != 3 {
		t.Errorf("Expected 3 sukarela transactions, got %d", total)
	}

	t.Logf("✓ Filter sukarela: %d transactions", len(sukarelaOnly))

	t.Log("✅ Filter by type integration test passed")
}

// TestSimpananIntegration_DateRangeFilter tests filtering by date range
func TestSimpananIntegration_DateRangeFilter(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)

	// Setup Chart of Accounts (required for journal posting)
	setupChartOfAccounts(db, koperasiID)

	anggota := &models.Anggota{
		IDKoperasi:    koperasiID,
		NomorAnggota:  "D001",
		NamaLengkap:  "Date Filter Member",
		NIK:           "1234567890123456",
		TempatLahir:   "Jakarta",
		TanggalLahir:  nil,
		JenisKelamin:  "L",
		Alamat:        "Test",
		RT:            "001",
		RW:            "002",
		Kelurahan:     "Test",
		Kecamatan:     "Test",
		KotaKabupaten: "Jakarta",
		Provinsi:      "DKI",
		KodePos:       "12345",
		NoTelepon:          "08123456789",
		Email:         "date@test.com",
		Status: models.StatusAktif,
	}
	db.Create(anggota)

	transaksiService := services.NewTransaksiService(db)
	simpananService := services.NewSimpananService(db, transaksiService)

	// Create transactions across different months
	dates := []time.Time{
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 3, 5, 0, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		req := &services.CatatSetoranRequest{
			IDAnggota:         anggota.ID,
			TipeSimpanan:     models.SimpananWajib,
			JumlahSetoran:    50000,
			TanggalTransaksi:  date,
			Keterangan:        "Transaction on " + date.Format("2006-01-02"),
		}
		simpananService.CatatSetoran(koperasiID, uuid.Nil, req)
	}

	// Filter January only
	januaryTransactions, total, err := simpananService.DapatkanSemuaTransaksiSimpanan(koperasiID, "", &anggota.ID, "2025-01-01", "2025-01-31", 1, 10)
	if err != nil {
		t.Fatalf("Failed to filter January: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected 2 January transactions, got %d", total)
	}

	t.Logf("✓ January filter: %d transactions", len(januaryTransactions))

	// Filter February only
	februaryTransactions, total, err := simpananService.DapatkanSemuaTransaksiSimpanan(koperasiID, "", &anggota.ID, "2025-02-01", "2025-02-28", 1, 10)
	if err != nil {
		t.Fatalf("Failed to filter February: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected 2 February transactions, got %d", total)
	}

	t.Logf("✓ February filter: %d transactions", len(februaryTransactions))

	// Filter Q1 (January to March)
	q1Transactions, total, err := simpananService.DapatkanSemuaTransaksiSimpanan(koperasiID, "", &anggota.ID, "2025-01-01", "2025-03-31", 1, 10)
	if err != nil {
		t.Fatalf("Failed to filter Q1: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected 5 Q1 transactions, got %d", total)
	}

	t.Logf("✓ Q1 filter: %d transactions", len(q1Transactions))

	t.Log("✅ Date range filter integration test passed")
}
