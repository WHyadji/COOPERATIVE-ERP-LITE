package services

import (
	"context"
	"cooperative-erp-lite/internal/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// queryCounter implements a custom logger to count database queries
type queryCounter struct {
	count int
}

func (qc *queryCounter) LogMode(level logger.LogLevel) logger.Interface {
	return qc
}

func (qc *queryCounter) Info(ctx context.Context, msg string, data ...interface{}) {}

func (qc *queryCounter) Warn(ctx context.Context, msg string, data ...interface{}) {}

func (qc *queryCounter) Error(ctx context.Context, msg string, data ...interface{}) {}

func (qc *queryCounter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	// Only count SELECT queries (ignore BEGIN, COMMIT, etc.)
	if len(sql) > 6 && sql[:6] == "SELECT" {
		qc.count++
	}
}

// setupTestDBWithQueryCounter creates a test database with query counting
func setupTestDBWithQueryCounter(t *testing.T) (*gorm.DB, *queryCounter) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	counter := &queryCounter{}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: counter,
	})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil, nil
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Anggota{},
		&models.Simpanan{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db, counter
}

// TestN1QueryCount verifies that DapatkanLaporanSaldoAnggota uses ≤3 queries
// regardless of the number of members (eliminating N+1 problem)
func TestN1QueryCount(t *testing.T) {
	db, counter := setupTestDBWithQueryCounter(t)
	if db == nil {
		return
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Query Count Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create 100 test members with deposits
	members := make([]models.Anggota, 100)
	for i := 0; i < 100; i++ {
		members[i] = models.Anggota{
			IDKoperasi:   koperasi.ID,
			NomorAnggota: fmt.Sprintf("A%04d", i+1),
			NamaLengkap:  fmt.Sprintf("Anggota Test %d", i+1),
			Status:       models.StatusAktif,
		}
	}
	db.CreateInBatches(members, 100)

	// Create deposits
	deposits := make([]models.Simpanan, 0, 300)
	for _, member := range members {
		deposits = append(deposits,
			models.Simpanan{
				IDKoperasi:       koperasi.ID,
				IDAnggota:        member.ID,
				TipeSimpanan:     models.SimpananPokok,
				TanggalTransaksi: time.Now(),
				JumlahSetoran:    100000,
				NomorReferensi:   fmt.Sprintf("SMP-%s-0001", time.Now().Format("20060102")),
			},
			models.Simpanan{
				IDKoperasi:       koperasi.ID,
				IDAnggota:        member.ID,
				TipeSimpanan:     models.SimpananWajib,
				TanggalTransaksi: time.Now(),
				JumlahSetoran:    50000,
				NomorReferensi:   fmt.Sprintf("SMP-%s-0002", time.Now().Format("20060102")),
			},
			models.Simpanan{
				IDKoperasi:       koperasi.ID,
				IDAnggota:        member.ID,
				TipeSimpanan:     models.SimpananSukarela,
				TanggalTransaksi: time.Now(),
				JumlahSetoran:    25000,
				NomorReferensi:   fmt.Sprintf("SMP-%s-0003", time.Now().Format("20060102")),
			},
		)
	}
	db.CreateInBatches(deposits, 100)

	// Reset counter before testing
	counter.count = 0

	// Execute the optimized function
	service := NewSimpananService(db, nil)
	laporan, err := service.DapatkanLaporanSaldoAnggota(koperasi.ID)

	// Verify no errors
	assert.NoError(t, err)
	assert.NotNil(t, laporan)
	assert.Equal(t, 100, len(laporan))

	// Verify query count is ≤3 (1 main query + 1 for members without deposits + potential overhead)
	// This is the key test: with N+1 pattern, this would be 101+ queries
	assert.LessOrEqual(t, counter.count, 3,
		"Expected ≤3 queries but got %d. N+1 problem may not be fully resolved.", counter.count)

	t.Logf("Query count: %d (target: ≤3)", counter.count)

	// Cleanup
	db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
	db.Delete(koperasi)
}

// TestDapatkanLaporanSaldoAnggota_Correctness verifies data accuracy
func TestDapatkanLaporanSaldoAnggota_Correctness(t *testing.T) {
	db, _ := setupTestDBWithQueryCounter(t)
	if db == nil {
		return
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Correctness Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test members
	member1 := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0001",
		NamaLengkap:  "Test Member 1",
		Status:       models.StatusAktif,
	}
	member2 := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0002",
		NamaLengkap:  "Test Member 2",
		Status:       models.StatusAktif,
	}
	member3 := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0003",
		NamaLengkap:  "Test Member 3 - No Deposits",
		Status:       models.StatusAktif,
	}
	db.Create(member1)
	db.Create(member2)
	db.Create(member3)

	// Create deposits for member 1
	db.Create(&models.Simpanan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        member1.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    100000,
		NomorReferensi:   "SMP-TEST-001",
	})
	db.Create(&models.Simpanan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        member1.ID,
		TipeSimpanan:     models.SimpananWajib,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    50000,
		NomorReferensi:   "SMP-TEST-002",
	})

	// Create deposits for member 2
	db.Create(&models.Simpanan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        member2.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    200000,
		NomorReferensi:   "SMP-TEST-003",
	})
	db.Create(&models.Simpanan{
		IDKoperasi:       koperasi.ID,
		IDAnggota:        member2.ID,
		TipeSimpanan:     models.SimpananSukarela,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    75000,
		NomorReferensi:   "SMP-TEST-004",
	})

	// Execute function
	service := NewSimpananService(db, nil)
	laporan, err := service.DapatkanLaporanSaldoAnggota(koperasi.ID)

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, laporan)
	assert.Equal(t, 3, len(laporan), "Should include all 3 members")

	// Find members in report
	var report1, report2, report3 *models.SaldoSimpananAnggota
	for i := range laporan {
		switch laporan[i].NomorAnggota {
		case "A0001":
			report1 = &laporan[i]
		case "A0002":
			report2 = &laporan[i]
		case "A0003":
			report3 = &laporan[i]
		}
	}

	// Verify member 1
	assert.NotNil(t, report1)
	assert.Equal(t, float64(100000), report1.SimpananPokok)
	assert.Equal(t, float64(50000), report1.SimpananWajib)
	assert.Equal(t, float64(0), report1.SimpananSukarela)
	assert.Equal(t, float64(150000), report1.TotalSimpanan)

	// Verify member 2
	assert.NotNil(t, report2)
	assert.Equal(t, float64(200000), report2.SimpananPokok)
	assert.Equal(t, float64(0), report2.SimpananWajib)
	assert.Equal(t, float64(75000), report2.SimpananSukarela)
	assert.Equal(t, float64(275000), report2.TotalSimpanan)

	// Verify member 3 (no deposits)
	assert.NotNil(t, report3)
	assert.Equal(t, float64(0), report3.SimpananPokok)
	assert.Equal(t, float64(0), report3.SimpananWajib)
	assert.Equal(t, float64(0), report3.SimpananSukarela)
	assert.Equal(t, float64(0), report3.TotalSimpanan)

	// Cleanup
	db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
	db.Delete(koperasi)
}

// TestDapatkanLaporanSaldoAnggota_MultiTenant verifies multi-tenant isolation
func TestDapatkanLaporanSaldoAnggota_MultiTenant(t *testing.T) {
	db, _ := setupTestDBWithQueryCounter(t)
	if db == nil {
		return
	}

	// Create two cooperatives
	koperasi1 := &models.Koperasi{
		NamaKoperasi: "Koperasi 1",
		Alamat:       "Address 1",
	}
	koperasi2 := &models.Koperasi{
		NamaKoperasi: "Koperasi 2",
		Alamat:       "Address 2",
	}
	db.Create(koperasi1)
	db.Create(koperasi2)

	// Create members for each cooperative
	member1 := &models.Anggota{
		IDKoperasi:   koperasi1.ID,
		NomorAnggota: "A0001",
		NamaLengkap:  "Member Koperasi 1",
		Status:       models.StatusAktif,
	}
	member2 := &models.Anggota{
		IDKoperasi:   koperasi2.ID,
		NomorAnggota: "A0001",
		NamaLengkap:  "Member Koperasi 2",
		Status:       models.StatusAktif,
	}
	db.Create(member1)
	db.Create(member2)

	// Create deposits
	db.Create(&models.Simpanan{
		IDKoperasi:       koperasi1.ID,
		IDAnggota:        member1.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    100000,
		NomorReferensi:   "SMP-K1-001",
	})
	db.Create(&models.Simpanan{
		IDKoperasi:       koperasi2.ID,
		IDAnggota:        member2.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    200000,
		NomorReferensi:   "SMP-K2-001",
	})

	// Get report for koperasi 1
	service := NewSimpananService(db, nil)
	laporan1, err := service.DapatkanLaporanSaldoAnggota(koperasi1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(laporan1))
	assert.Equal(t, "A0001", laporan1[0].NomorAnggota)
	assert.Equal(t, float64(100000), laporan1[0].SimpananPokok)

	// Get report for koperasi 2
	laporan2, err := service.DapatkanLaporanSaldoAnggota(koperasi2.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(laporan2))
	assert.Equal(t, "A0001", laporan2[0].NomorAnggota)
	assert.Equal(t, float64(200000), laporan2[0].SimpananPokok)

	// Verify isolation: each cooperative only sees their own data
	assert.NotEqual(t, laporan1[0].IDAnggota, laporan2[0].IDAnggota)

	// Cleanup
	db.Exec("DELETE FROM simpanan WHERE id_koperasi IN (?, ?)", koperasi1.ID, koperasi2.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi IN (?, ?)", koperasi1.ID, koperasi2.ID)
	db.Delete(koperasi1)
	db.Delete(koperasi2)
}

// TestCatatSetoran_SimpananPokokUniqueness verifies that Simpanan Pokok can only be paid once
func TestCatatSetoran_SimpananPokokUniqueness(t *testing.T) {
	db, _ := setupTestDBWithQueryCounter(t)
	if db == nil {
		return
	}

	// Auto migrate additional tables needed for this test
	err := db.AutoMigrate(
		&models.Pengguna{},
		&models.Akun{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate additional tables: %v", err)
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Simpanan Pokok Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test member
	member := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0001",
		NamaLengkap:  "Test Member",
		Status:       models.StatusAktif,
	}
	db.Create(member)

	// Create test user for transaction
	user := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		Email:        "test@example.com",
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Peran:        models.PeranAdmin,
	}
	user.SetKataSandi("password123")
	db.Create(user)

	// Create TransaksiService (required for SimpananService)
	transaksiService := NewTransaksiService(db)
	service := NewSimpananService(db, transaksiService)

	// First deposit - should succeed
	req1 := &CatatSetoranRequest{
		IDAnggota:        member.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    100000,
		Keterangan:       "Simpanan pokok pertama",
	}

	result1, err := service.CatatSetoran(koperasi.ID, user.ID, req1)
	assert.NoError(t, err)
	assert.NotNil(t, result1)
	assert.Equal(t, models.SimpananPokok, result1.TipeSimpanan)
	assert.Equal(t, float64(100000), result1.JumlahSetoran)

	// Second deposit of Simpanan Pokok - should fail with specific error
	req2 := &CatatSetoranRequest{
		IDAnggota:        member.ID,
		TipeSimpanan:     models.SimpananPokok,
		TanggalTransaksi: time.Now(),
		JumlahSetoran:    50000,
		Keterangan:       "Simpanan pokok kedua (tidak boleh)",
	}

	result2, err := service.CatatSetoran(koperasi.ID, user.ID, req2)
	assert.Error(t, err)
	assert.Nil(t, result2)
	assert.Equal(t, "anggota sudah membayar simpanan pokok", err.Error())

	// Verify only one Simpanan Pokok record exists
	var count int64
	db.Model(&models.Simpanan{}).
		Where("id_anggota = ? AND tipe_simpanan = ?", member.ID, models.SimpananPokok).
		Count(&count)
	assert.Equal(t, int64(1), count)

	// Cleanup
	db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM baris_transaksi WHERE id_transaksi IN (SELECT id FROM transaksi WHERE id_koperasi = ?)", koperasi.ID)
	db.Exec("DELETE FROM transaksi WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM pengguna WHERE id_koperasi = ?", koperasi.ID)
	db.Delete(koperasi)
}

// TestCatatSetoran_OtherSimpananTypes verifies that Simpanan Wajib and Sukarela can be paid multiple times
func TestCatatSetoran_OtherSimpananTypes(t *testing.T) {
	db, _ := setupTestDBWithQueryCounter(t)
	if db == nil {
		return
	}

	// Auto migrate additional tables
	err := db.AutoMigrate(
		&models.Pengguna{},
		&models.Akun{},
		&models.Transaksi{},
		&models.BarisTransaksi{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate additional tables: %v", err)
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Test Multiple Deposits Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test member
	member := &models.Anggota{
		IDKoperasi:   koperasi.ID,
		NomorAnggota: "A0001",
		NamaLengkap:  "Test Member",
		Status:       models.StatusAktif,
	}
	db.Create(member)

	// Create test user
	user := &models.Pengguna{
		IDKoperasi:   koperasi.ID,
		Email:        "test@example.com",
		NamaLengkap:  "Test User",
		NamaPengguna: "testuser",
		Peran:        models.PeranAdmin,
	}
	user.SetKataSandi("password123")
	db.Create(user)

	// Create services
	transaksiService := NewTransaksiService(db)
	service := NewSimpananService(db, transaksiService)

	// Test multiple Simpanan Wajib deposits - should all succeed
	for i := 1; i <= 3; i++ {
		req := &CatatSetoranRequest{
			IDAnggota:        member.ID,
			TipeSimpanan:     models.SimpananWajib,
			TanggalTransaksi: time.Now().AddDate(0, i-1, 0), // Different months
			JumlahSetoran:    50000,
			Keterangan:       fmt.Sprintf("Simpanan wajib bulan %d", i),
		}

		result, err := service.CatatSetoran(koperasi.ID, user.ID, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.SimpananWajib, result.TipeSimpanan)
	}

	// Test multiple Simpanan Sukarela deposits - should all succeed
	for i := 1; i <= 2; i++ {
		req := &CatatSetoranRequest{
			IDAnggota:        member.ID,
			TipeSimpanan:     models.SimpananSukarela,
			TanggalTransaksi: time.Now(),
			JumlahSetoran:    25000,
			Keterangan:       fmt.Sprintf("Simpanan sukarela %d", i),
		}

		result, err := service.CatatSetoran(koperasi.ID, user.ID, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.SimpananSukarela, result.TipeSimpanan)
	}

	// Verify counts
	var countWajib, countSukarela int64
	db.Model(&models.Simpanan{}).
		Where("id_anggota = ? AND tipe_simpanan = ?", member.ID, models.SimpananWajib).
		Count(&countWajib)
	db.Model(&models.Simpanan{}).
		Where("id_anggota = ? AND tipe_simpanan = ?", member.ID, models.SimpananSukarela).
		Count(&countSukarela)

	assert.Equal(t, int64(3), countWajib, "Should allow multiple Simpanan Wajib deposits")
	assert.Equal(t, int64(2), countSukarela, "Should allow multiple Simpanan Sukarela deposits")

	// Cleanup
	db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM baris_transaksi WHERE id_transaksi IN (SELECT id FROM transaksi WHERE id_koperasi = ?)", koperasi.ID)
	db.Exec("DELETE FROM transaksi WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", koperasi.ID)
	db.Exec("DELETE FROM pengguna WHERE id_koperasi = ?", koperasi.ID)
	db.Delete(koperasi)
}
