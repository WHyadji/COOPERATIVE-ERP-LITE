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
