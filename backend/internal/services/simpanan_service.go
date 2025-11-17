package services

import (
	"cooperative-erp-lite/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SimpananService menangani logika bisnis simpanan anggota
type SimpananService struct {
	db               *gorm.DB
	transaksiService *TransaksiService
}

// NewSimpananService membuat instance baru SimpananService
func NewSimpananService(db *gorm.DB, transaksiService *TransaksiService) *SimpananService {
	return &SimpananService{
		db:               db,
		transaksiService: transaksiService,
	}
}

// CatatSetoranRequest adalah struktur request untuk catat setoran
type CatatSetoranRequest struct {
	IDAnggota        uuid.UUID           `json:"idAnggota" binding:"required"`
	TipeSimpanan     models.TipeSimpanan `json:"tipeSimpanan" binding:"required"`
	TanggalTransaksi time.Time           `json:"tanggalTransaksi" binding:"required"`
	JumlahSetoran    float64             `json:"jumlahSetoran" binding:"required,gt=0"`
	Keterangan       string              `json:"keterangan"`
}

// CatatSetoran mencatat setoran simpanan anggota
func (s *SimpananService) CatatSetoran(idKoperasi, idPengguna uuid.UUID, req *CatatSetoranRequest) (*models.SimpananResponse, error) {
	var simpanan *models.Simpanan

	// Wrap entire operation in a single transaction for data consistency
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Validasi anggota exists dan aktif
		var anggota models.Anggota
		if err := tx.Where("id = ? AND id_koperasi = ? AND status = ?", req.IDAnggota, idKoperasi, models.StatusAktif).
			First(&anggota).Error; err != nil {
			return errors.New("anggota tidak ditemukan atau tidak aktif")
		}

		// 2. Validasi jumlah setoran
		if req.JumlahSetoran <= 0 {
			return errors.New("jumlah setoran harus lebih dari 0")
		}

		// 3. Generate nomor referensi using transaction
		nomorReferensi, err := s.GenerateNomorReferensiWithTx(tx, idKoperasi, req.TanggalTransaksi)
		if err != nil {
			return err
		}

		// 4. Buat record simpanan
		simpanan = &models.Simpanan{
			IDKoperasi:       idKoperasi,
			IDAnggota:        req.IDAnggota,
			TipeSimpanan:     req.TipeSimpanan,
			TanggalTransaksi: req.TanggalTransaksi,
			JumlahSetoran:    req.JumlahSetoran,
			Keterangan:       req.Keterangan,
			NomorReferensi:   nomorReferensi,
			DibuatOleh:       idPengguna,
		}

		if err := tx.Create(simpanan).Error; err != nil {
			return errors.New("gagal mencatat setoran simpanan")
		}

		// 5. Auto-posting ke jurnal akuntansi within same transaction
		if err := s.postingSimpananWithTx(tx, idKoperasi, idPengguna, simpanan.ID); err != nil {
			return fmt.Errorf("gagal posting ke jurnal: %w", err)
		}

		return nil // Commit only if everything succeeds
	})

	if err != nil {
		return nil, err // Automatic rollback on any error
	}

	// Reload dengan relasi outside transaction
	s.db.Preload("Anggota").First(simpanan, simpanan.ID)

	response := simpanan.ToResponse()
	return &response, nil
}

// GenerateNomorReferensi menghasilkan nomor referensi setoran
// Format: SMP-YYYYMMDD-NNNN
// Uses row-level locking to prevent race conditions in concurrent requests
func (s *SimpananService) GenerateNomorReferensi(idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	var nomorReferensi string

	// Use transaction with row-level locking to prevent race conditions
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		nomorReferensi, err = s.GenerateNomorReferensiWithTx(tx, idKoperasi, tanggal)
		return err
	})

	if err != nil {
		return "", errors.New("gagal generate nomor referensi")
	}

	return nomorReferensi, nil
}

// GenerateNomorReferensiWithTx menghasilkan nomor referensi menggunakan existing transaction
// Format: SMP-YYYYMMDD-NNNN
func (s *SimpananService) GenerateNomorReferensiWithTx(tx *gorm.DB, idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	tanggalStr := tanggal.Format("20060102")
	tanggalDate := tanggal.Format("2006-01-02")

	// Lock and get the last deposit number for this date
	var lastSimpanan models.Simpanan
	err := tx.Where("id_koperasi = ? AND DATE(tanggal_transaksi) = ?", idKoperasi, tanggalDate).
		Order("nomor_referensi DESC").
		Limit(1).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&lastSimpanan).Error

	nomorUrut := 1

	// If there's a previous deposit, parse and increment
	if err == nil && lastSimpanan.NomorReferensi != "" {
		// Extract number from SMP-20250116-0001
		var parsedTanggal string
		var parsedUrut int
		_, scanErr := fmt.Sscanf(lastSimpanan.NomorReferensi, "SMP-%s-%04d", &parsedTanggal, &parsedUrut)
		if scanErr == nil && parsedTanggal == tanggalStr {
			nomorUrut = parsedUrut + 1
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	nomorReferensi := fmt.Sprintf("SMP-%s-%04d", tanggalStr, nomorUrut)
	return nomorReferensi, nil
}

// DapatkanSemuaTransaksiSimpanan mengambil daftar transaksi simpanan
func (s *SimpananService) DapatkanSemuaTransaksiSimpanan(idKoperasi uuid.UUID, tipeSimpanan string, idAnggota *uuid.UUID, tanggalMulai, tanggalAkhir string, page, pageSize int) ([]models.SimpananResponse, int64, error) {
	var simpananList []models.Simpanan
	var total int64

	query := s.db.Model(&models.Simpanan{}).Where("id_koperasi = ?", idKoperasi)

	// Apply filters
	if tipeSimpanan != "" {
		query = query.Where("tipe_simpanan = ?", tipeSimpanan)
	}
	if idAnggota != nil {
		query = query.Where("id_anggota = ?", *idAnggota)
	}
	if tanggalMulai != "" {
		query = query.Where("tanggal_transaksi >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("tanggal_transaksi <= ?", tanggalAkhir)
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("tanggal_transaksi DESC").
		Preload("Anggota").
		Find(&simpananList).Error

	if err != nil {
		return nil, 0, errors.New("gagal mengambil daftar transaksi simpanan")
	}

	// Convert to response
	responses := make([]models.SimpananResponse, len(simpananList))
	for i, simpanan := range simpananList {
		responses[i] = simpanan.ToResponse()
	}

	return responses, total, nil
}

// DapatkanSaldoAnggota mengambil saldo simpanan per anggota
func (s *SimpananService) DapatkanSaldoAnggota(idAnggota uuid.UUID) (*models.SaldoSimpananAnggota, error) {
	// Validasi anggota exists
	var anggota models.Anggota
	err := s.db.Where("id = ?", idAnggota).First(&anggota).Error
	if err != nil {
		return nil, errors.New("anggota tidak ditemukan")
	}

	// Hitung saldo per tipe simpanan
	type SaldoByTipe struct {
		TipeSimpanan models.TipeSimpanan
		Total        float64
	}

	var saldoList []SaldoByTipe
	err = s.db.Model(&models.Simpanan{}).
		Select("tipe_simpanan, COALESCE(SUM(jumlah_setoran), 0) as total").
		Where("id_anggota = ?", idAnggota).
		Group("tipe_simpanan").
		Scan(&saldoList).Error

	if err != nil {
		return nil, errors.New("gagal menghitung saldo simpanan")
	}

	// Build response
	saldo := &models.SaldoSimpananAnggota{
		IDAnggota:    idAnggota,
		NomorAnggota: anggota.NomorAnggota,
		NamaAnggota:  anggota.NamaLengkap,
	}

	for _, item := range saldoList {
		switch item.TipeSimpanan {
		case models.SimpananPokok:
			saldo.SimpananPokok = item.Total
		case models.SimpananWajib:
			saldo.SimpananWajib = item.Total
		case models.SimpananSukarela:
			saldo.SimpananSukarela = item.Total
		}
	}

	saldo.TotalSimpanan = saldo.SimpananPokok + saldo.SimpananWajib + saldo.SimpananSukarela

	return saldo, nil
}

// DapatkanRingkasanSimpanan mengambil ringkasan total simpanan koperasi
func (s *SimpananService) DapatkanRingkasanSimpanan(idKoperasi uuid.UUID) (*models.RingkasanSimpanan, error) {
	type SaldoByTipe struct {
		TipeSimpanan models.TipeSimpanan
		Total        float64
	}

	var saldoList []SaldoByTipe
	err := s.db.Model(&models.Simpanan{}).
		Select("tipe_simpanan, COALESCE(SUM(jumlah_setoran), 0) as total").
		Where("id_koperasi = ?", idKoperasi).
		Group("tipe_simpanan").
		Scan(&saldoList).Error

	if err != nil {
		return nil, errors.New("gagal menghitung ringkasan simpanan")
	}

	// Hitung jumlah anggota yang memiliki simpanan
	var jumlahAnggota int64
	s.db.Model(&models.Simpanan{}).
		Where("id_koperasi = ?", idKoperasi).
		Distinct("id_anggota").
		Count(&jumlahAnggota)

	// Build response
	ringkasan := &models.RingkasanSimpanan{
		JumlahAnggota: jumlahAnggota,
	}

	for _, item := range saldoList {
		switch item.TipeSimpanan {
		case models.SimpananPokok:
			ringkasan.TotalSimpananPokok = item.Total
		case models.SimpananWajib:
			ringkasan.TotalSimpananWajib = item.Total
		case models.SimpananSukarela:
			ringkasan.TotalSimpananSukarela = item.Total
		}
	}

	ringkasan.TotalSemuaSimpanan = ringkasan.TotalSimpananPokok + ringkasan.TotalSimpananWajib + ringkasan.TotalSimpananSukarela

	return ringkasan, nil
}

// DapatkanLaporanSaldoAnggota mengambil laporan saldo semua anggota
func (s *SimpananService) DapatkanLaporanSaldoAnggota(idKoperasi uuid.UUID) ([]models.SaldoSimpananAnggota, error) {
	// Dapatkan semua anggota aktif
	var anggotaList []models.Anggota
	err := s.db.Where("id_koperasi = ? AND status = ?", idKoperasi, models.StatusAktif).
		Order("nomor_anggota ASC").
		Find(&anggotaList).Error

	if err != nil {
		return nil, errors.New("gagal mengambil daftar anggota")
	}

	// Hitung saldo untuk setiap anggota
	var laporan []models.SaldoSimpananAnggota

	for _, anggota := range anggotaList {
		saldo, err := s.DapatkanSaldoAnggota(anggota.ID)
		if err != nil {
			continue // Skip jika error
		}
		laporan = append(laporan, *saldo)
	}

	return laporan, nil
}

// HitungTotalSimpananByTipe menghitung total simpanan berdasarkan tipe
func (s *SimpananService) HitungTotalSimpananByTipe(idKoperasi uuid.UUID, tipeSimpanan models.TipeSimpanan) (float64, error) {
	var total float64
	err := s.db.Model(&models.Simpanan{}).
		Select("COALESCE(SUM(jumlah_setoran), 0)").
		Where("id_koperasi = ? AND tipe_simpanan = ?", idKoperasi, tipeSimpanan).
		Scan(&total).Error

	if err != nil {
		return 0, errors.New("gagal menghitung total simpanan")
	}

	return total, nil
}

// postingSimpananWithTx creates journal entry for simpanan within an existing transaction
// This ensures atomicity - if posting fails, simpanan record is also rolled back
func (s *SimpananService) postingSimpananWithTx(tx *gorm.DB, idKoperasi, idPengguna, idSimpanan uuid.UUID) error {
	// Get simpanan data
	var simpanan models.Simpanan
	if err := tx.First(&simpanan, idSimpanan).Error; err != nil {
		return errors.New("simpanan tidak ditemukan")
	}

	// Determine account based on simpanan type
	var kodeAkunModal string
	switch simpanan.TipeSimpanan {
	case models.SimpananPokok:
		kodeAkunModal = "3101" // Simpanan Pokok
	case models.SimpananWajib:
		kodeAkunModal = "3102" // Simpanan Wajib
	case models.SimpananSukarela:
		kodeAkunModal = "3103" // Simpanan Sukarela
	default:
		return fmt.Errorf("tipe simpanan tidak valid: %s", simpanan.TipeSimpanan)
	}

	// Get required accounts
	var akunKas, akunModal models.Akun
	if err := tx.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1101").First(&akunKas).Error; err != nil {
		return errors.New("akun kas tidak ditemukan")
	}
	if err := tx.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, kodeAkunModal).First(&akunModal).Error; err != nil {
		return fmt.Errorf("akun simpanan tidak ditemukan: %s", kodeAkunModal)
	}

	// Generate journal number within the same transaction
	tanggalStr := simpanan.TanggalTransaksi.Format("20060102")
	tanggalDate := simpanan.TanggalTransaksi.Format("2006-01-02")

	var lastTransaksi models.Transaksi
	err := tx.Where("id_koperasi = ? AND DATE(tanggal_transaksi) = ?", idKoperasi, tanggalDate).
		Order("nomor_jurnal DESC").
		Limit(1).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&lastTransaksi).Error

	nomorUrut := 1
	if err == nil && lastTransaksi.NomorJurnal != "" {
		var parsedTanggal string
		var parsedUrut int
		_, scanErr := fmt.Sscanf(lastTransaksi.NomorJurnal, "JRN-%s-%04d", &parsedTanggal, &parsedUrut)
		if scanErr == nil && parsedTanggal == tanggalStr {
			nomorUrut = parsedUrut + 1
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	nomorJurnal := fmt.Sprintf("JRN-%s-%04d", tanggalStr, nomorUrut)

	// Create journal entry
	transaksi := &models.Transaksi{
		IDKoperasi:       idKoperasi,
		NomorJurnal:      nomorJurnal,
		TanggalTransaksi: simpanan.TanggalTransaksi,
		Deskripsi:        fmt.Sprintf("Setoran %s - %s", simpanan.TipeSimpanan, simpanan.NomorReferensi),
		NomorReferensi:   simpanan.NomorReferensi,
		TipeTransaksi:    "simpanan",
		TotalDebit:       simpanan.JumlahSetoran,
		TotalKredit:      simpanan.JumlahSetoran,
		StatusBalanced:   true,
		DibuatOleh:       idPengguna,
	}

	if err := tx.Create(transaksi).Error; err != nil {
		return fmt.Errorf("gagal membuat jurnal: %w", err)
	}

	// Create journal lines
	barisTransaksi := []models.BarisTransaksi{
		{
			IDTransaksi:  transaksi.ID,
			IDAkun:       akunKas.ID,
			JumlahDebit:  simpanan.JumlahSetoran,
			JumlahKredit: 0,
			Keterangan:   "Penerimaan setoran simpanan",
		},
		{
			IDTransaksi:  transaksi.ID,
			IDAkun:       akunModal.ID,
			JumlahDebit:  0,
			JumlahKredit: simpanan.JumlahSetoran,
			Keterangan:   "Simpanan anggota",
		},
	}

	for _, baris := range barisTransaksi {
		if err := tx.Create(&baris).Error; err != nil {
			return fmt.Errorf("gagal membuat baris jurnal: %w", err)
		}
	}

	// Update simpanan with transaction ID
	simpanan.IDTransaksi = &transaksi.ID
	if err := tx.Save(&simpanan).Error; err != nil {
		return fmt.Errorf("gagal update simpanan dengan ID transaksi: %w", err)
	}

	return nil
}
