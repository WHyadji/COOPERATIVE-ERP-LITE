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
	db                *gorm.DB
	transaksiService  *TransaksiService
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
	IDAnggota        uuid.UUID          `json:"idAnggota" binding:"required"`
	TipeSimpanan     models.TipeSimpanan `json:"tipeSimpanan" binding:"required"`
	TanggalTransaksi time.Time          `json:"tanggalTransaksi" binding:"required"`
	JumlahSetoran    float64            `json:"jumlahSetoran" binding:"required,gt=0"`
	Keterangan       string             `json:"keterangan"`
}

// CatatSetoran mencatat setoran simpanan anggota
func (s *SimpananService) CatatSetoran(idKoperasi, idPengguna uuid.UUID, req *CatatSetoranRequest) (*models.SimpananResponse, error) {
	// Validasi anggota exists dan aktif
	var anggota models.Anggota
	err := s.db.Where("id = ? AND id_koperasi = ? AND status = ?", req.IDAnggota, idKoperasi, models.StatusAktif).
		First(&anggota).Error
	if err != nil {
		return nil, errors.New("anggota tidak ditemukan atau tidak aktif")
	}

	// Validasi jumlah setoran
	if req.JumlahSetoran <= 0 {
		return nil, errors.New("jumlah setoran harus lebih dari 0")
	}

	// Generate nomor referensi
	nomorReferensi, err := s.GenerateNomorReferensi(idKoperasi, req.TanggalTransaksi)
	if err != nil {
		return nil, err
	}

	// Buat record simpanan
	simpanan := &models.Simpanan{
		IDKoperasi:       idKoperasi,
		IDAnggota:        req.IDAnggota,
		TipeSimpanan:     req.TipeSimpanan,
		TanggalTransaksi: req.TanggalTransaksi,
		JumlahSetoran:    req.JumlahSetoran,
		Keterangan:       req.Keterangan,
		NomorReferensi:   nomorReferensi,
		DibuatOleh:       idPengguna,
	}

	// Simpan ke database
	err = s.db.Create(simpanan).Error
	if err != nil {
		return nil, errors.New("gagal mencatat setoran simpanan")
	}

	// Auto-posting ke jurnal akuntansi
	err = s.transaksiService.PostingOtomatisSimpanan(idKoperasi, idPengguna, simpanan.ID)
	if err != nil {
		// Rollback simpanan jika posting gagal
		s.db.Delete(simpanan)
		return nil, fmt.Errorf("gagal posting ke jurnal: %w", err)
	}

	// Reload dengan relasi
	s.db.Preload("Anggota").First(simpanan, simpanan.ID)

	response := simpanan.ToResponse()
	return &response, nil
}

// GenerateNomorReferensi menghasilkan nomor referensi setoran
// Format: SMP-YYYYMMDD-NNNN
// Uses row-level locking to prevent race conditions in concurrent requests
func (s *SimpananService) GenerateNomorReferensi(idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	tanggalStr := tanggal.Format("20060102")
	tanggalDate := tanggal.Format("2006-01-02")
	var nomorReferensi string

	// Use transaction with row-level locking to prevent race conditions
	err := s.db.Transaction(func(tx *gorm.DB) error {
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
			return err
		}

		nomorReferensi = fmt.Sprintf("SMP-%s-%04d", tanggalStr, nomorUrut)
		return nil
	})

	if err != nil {
		return "", errors.New("gagal generate nomor referensi")
	}

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
