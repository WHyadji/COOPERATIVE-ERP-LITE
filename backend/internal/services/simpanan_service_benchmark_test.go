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

// setupBenchmarkDB creates a test database with realistic data
func setupBenchmarkDB(b *testing.B) (*gorm.DB, uuid.UUID, uuid.UUID) {
	dsn := "host=localhost user=postgres password=postgres dbname=koperasi_erp_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	// Use silent logger for benchmarks to avoid I/O overhead
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database: %v", err)
		return nil, uuid.Nil, uuid.Nil
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&models.Koperasi{},
		&models.Anggota{},
		&models.Simpanan{},
	)
	if err != nil {
		b.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create test cooperative
	koperasi := &models.Koperasi{
		NamaKoperasi: "Benchmark Test Koperasi",
		Alamat:       "Test Address",
	}
	db.Create(koperasi)

	// Create test user
	pengguna := &models.Pengguna{
		IDKoperasi: koperasi.ID,
		NamaLengkap: "Benchmark User",
		Email:      "benchmark@test.com",
	}
	db.Create(pengguna)

	return db, koperasi.ID, pengguna.ID
}

// seedMembers creates N test members with deposits
func seedMembers(db *gorm.DB, idKoperasi uuid.UUID, count int) {
	// Clear existing data
	db.Exec("DELETE FROM simpanan WHERE id_koperasi = ?", idKoperasi)
	db.Exec("DELETE FROM anggota WHERE id_koperasi = ?", idKoperasi)

	members := make([]models.Anggota, count)
	for i := 0; i < count; i++ {
		members[i] = models.Anggota{
			IDKoperasi:   idKoperasi,
			NomorAnggota: fmt.Sprintf("A%04d", i+1),
			NamaLengkap:  fmt.Sprintf("Anggota Benchmark %d", i+1),
			Status:       models.StatusAktif,
		}
	}

	// Batch insert members
	db.CreateInBatches(members, 100)

	// Create deposits for each member
	deposits := make([]models.Simpanan, 0, count*3)
	for _, member := range members {
		// Simpanan Pokok
		deposits = append(deposits, models.Simpanan{
			IDKoperasi:       idKoperasi,
			IDAnggota:        member.ID,
			TipeSimpanan:     models.SimpananPokok,
			TanggalTransaksi: time.Now().AddDate(0, -1, 0),
			JumlahSetoran:    100000,
			NomorReferensi:   fmt.Sprintf("SMP-%s-0001", time.Now().Format("20060102")),
		})

		// Simpanan Wajib
		deposits = append(deposits, models.Simpanan{
			IDKoperasi:       idKoperasi,
			IDAnggota:        member.ID,
			TipeSimpanan:     models.SimpananWajib,
			TanggalTransaksi: time.Now().AddDate(0, 0, -15),
			JumlahSetoran:    50000,
			NomorReferensi:   fmt.Sprintf("SMP-%s-0002", time.Now().Format("20060102")),
		})

		// Simpanan Sukarela (only for some members)
		if member.ID.String()[0]%2 == 0 {
			deposits = append(deposits, models.Simpanan{
				IDKoperasi:       idKoperasi,
				IDAnggota:        member.ID,
				TipeSimpanan:     models.SimpananSukarela,
				TanggalTransaksi: time.Now().AddDate(0, 0, -7),
				JumlahSetoran:    25000,
				NomorReferensi:   fmt.Sprintf("SMP-%s-0003", time.Now().Format("20060102")),
			})
		}
	}

	// Batch insert deposits
	db.CreateInBatches(deposits, 100)
}

// BenchmarkGenerateLaporanSaldoAnggota_100Members benchmarks report generation with 100 members
func BenchmarkGenerateLaporanSaldoAnggota_100Members(b *testing.B) {
	db, idKoperasi, _ := setupBenchmarkDB(b)
	if db == nil {
		return
	}

	seedMembers(db, idKoperasi, 100)
	service := NewSimpananService(db, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.DapatkanLaporanSaldoAnggota(idKoperasi)
		if err != nil {
			b.Fatalf("Error generating report: %v", err)
		}
	}
}

// BenchmarkGenerateLaporanSaldoAnggota_1000Members benchmarks report generation with 1000 members
// Target: < 100ms per operation
func BenchmarkGenerateLaporanSaldoAnggota_1000Members(b *testing.B) {
	db, idKoperasi, _ := setupBenchmarkDB(b)
	if db == nil {
		return
	}

	seedMembers(db, idKoperasi, 1000)
	service := NewSimpananService(db, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.DapatkanLaporanSaldoAnggota(idKoperasi)
		if err != nil {
			b.Fatalf("Error generating report: %v", err)
		}
	}
}

// BenchmarkGenerateLaporanSaldoAnggota_5000Members benchmarks report generation with 5000 members
// This tests scalability for larger cooperatives
func BenchmarkGenerateLaporanSaldoAnggota_5000Members(b *testing.B) {
	db, idKoperasi, _ := setupBenchmarkDB(b)
	if db == nil {
		return
	}

	seedMembers(db, idKoperasi, 5000)
	service := NewSimpananService(db, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.DapatkanLaporanSaldoAnggota(idKoperasi)
		if err != nil {
			b.Fatalf("Error generating report: %v", err)
		}
	}
}
