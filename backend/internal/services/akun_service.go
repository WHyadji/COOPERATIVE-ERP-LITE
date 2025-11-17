package services

import (
	"cooperative-erp-lite/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AkunService menangani logika bisnis Chart of Accounts
type AkunService struct {
	db *gorm.DB
}

// NewAkunService membuat instance baru AkunService
func NewAkunService(db *gorm.DB) *AkunService {
	return &AkunService{db: db}
}

// BuatAkunRequest adalah struktur request untuk membuat akun
type BuatAkunRequest struct {
	KodeAkun    string          `json:"kodeAkun" binding:"required"`
	NamaAkun    string          `json:"namaAkun" binding:"required"`
	TipeAkun    models.TipeAkun `json:"tipeAkun" binding:"required"`
	IDInduk     *uuid.UUID      `json:"idInduk"`
	NormalSaldo string          `json:"normalSaldo"` // Optional, auto-detect dari tipe
	Deskripsi   string          `json:"deskripsi"`
}

// BuatAkun membuat akun baru
func (s *AkunService) BuatAkun(idKoperasi uuid.UUID, req *BuatAkunRequest) (*models.AkunResponse, error) {
	// Cek apakah kode akun sudah ada
	var count int64
	s.db.Model(&models.Akun{}).
		Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, req.KodeAkun).
		Count(&count)

	if count > 0 {
		return nil, errors.New("kode akun sudah digunakan")
	}

	// Validasi parent jika ada
	if req.IDInduk != nil {
		var parentAkun models.Akun
		err := s.db.Where("id = ? AND id_koperasi = ?", req.IDInduk, idKoperasi).First(&parentAkun).Error
		if err != nil {
			return nil, errors.New("akun induk tidak ditemukan")
		}
	}

	// Buat akun baru
	akun := &models.Akun{
		IDKoperasi:  idKoperasi,
		KodeAkun:    req.KodeAkun,
		NamaAkun:    req.NamaAkun,
		TipeAkun:    req.TipeAkun,
		IDInduk:     req.IDInduk,
		NormalSaldo: req.NormalSaldo,
		Deskripsi:   req.Deskripsi,
		StatusAktif: true,
	}

	// BeforeCreate hook akan set normal saldo jika kosong
	err := s.db.Create(akun).Error
	if err != nil {
		return nil, errors.New("gagal membuat akun")
	}

	response := akun.ToResponse()
	return &response, nil
}

// DapatkanSemuaAkun mengambil daftar akun dengan filter
func (s *AkunService) DapatkanSemuaAkun(idKoperasi uuid.UUID, tipeAkun string, statusAktif *bool) ([]models.AkunResponse, error) {
	var akunList []models.Akun

	query := s.db.Where("id_koperasi = ?", idKoperasi).Preload("AkunInduk")

	// Apply filters
	if tipeAkun != "" {
		query = query.Where("tipe_akun = ?", tipeAkun)
	}
	if statusAktif != nil {
		query = query.Where("status_aktif = ?", *statusAktif)
	}

	err := query.Order("kode_akun ASC").Find(&akunList).Error
	if err != nil {
		return nil, errors.New("gagal mengambil daftar akun")
	}

	// Convert to response
	responses := make([]models.AkunResponse, len(akunList))
	for i, akun := range akunList {
		responses[i] = akun.ToResponse()
	}

	return responses, nil
}

// DapatkanAkun mengambil akun berdasarkan ID
func (s *AkunService) DapatkanAkun(id uuid.UUID) (*models.AkunResponse, error) {
	var akun models.Akun
	err := s.db.Preload("AkunInduk").Where("id = ?", id).First(&akun).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akun tidak ditemukan")
		}
		return nil, err
	}

	response := akun.ToResponse()
	return &response, nil
}

// DapatkanAkunByKode mengambil akun berdasarkan kode
func (s *AkunService) DapatkanAkunByKode(idKoperasi uuid.UUID, kodeAkun string) (*models.AkunResponse, error) {
	var akun models.Akun
	err := s.db.Preload("AkunInduk").
		Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, kodeAkun).
		First(&akun).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akun tidak ditemukan")
		}
		return nil, err
	}

	response := akun.ToResponse()
	return &response, nil
}

// PerbaruiAkunRequest adalah struktur request untuk update akun
type PerbaruiAkunRequest struct {
	NamaAkun    string `json:"namaAkun"`
	Deskripsi   string `json:"deskripsi"`
	StatusAktif *bool  `json:"statusAktif"`
}

// PerbaruiAkun mengupdate data akun
func (s *AkunService) PerbaruiAkun(idKoperasi, id uuid.UUID, req *PerbaruiAkunRequest) (*models.AkunResponse, error) {
	// Cek apakah akun ada DAN milik koperasi yang benar (multi-tenant validation)
	var akun models.Akun
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&akun).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akun tidak ditemukan atau tidak memiliki akses")
		}
		return nil, err
	}

	// Update fields
	if req.NamaAkun != "" {
		akun.NamaAkun = req.NamaAkun
	}
	if req.Deskripsi != "" {
		akun.Deskripsi = req.Deskripsi
	}
	if req.StatusAktif != nil {
		akun.StatusAktif = *req.StatusAktif
	}

	err = s.db.Save(&akun).Error
	if err != nil {
		return nil, errors.New("gagal memperbarui akun")
	}

	response := akun.ToResponse()
	return &response, nil
}

// HapusAkun menghapus akun (dengan validasi)
func (s *AkunService) HapusAkun(idKoperasi, id uuid.UUID) error {
	// Cek apakah akun ada DAN milik koperasi yang benar (multi-tenant validation)
	var akun models.Akun
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&akun).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("akun tidak ditemukan atau tidak memiliki akses")
		}
		return err
	}

	// Cek apakah ada transaksi terkait
	var countTransaksi int64
	s.db.Model(&models.BarisTransaksi{}).Where("id_akun = ?", id).Count(&countTransaksi)

	if countTransaksi > 0 {
		return errors.New("tidak dapat menghapus akun yang sudah memiliki transaksi")
	}

	// Cek apakah ada sub-akun
	var countSubAkun int64
	s.db.Model(&models.Akun{}).Where("id_induk = ?", id).Count(&countSubAkun)

	if countSubAkun > 0 {
		return errors.New("tidak dapat menghapus akun yang memiliki sub-akun")
	}

	// Soft delete
	err = s.db.Delete(&akun).Error
	if err != nil {
		return errors.New("gagal menghapus akun")
	}

	return nil
}

// HitungSaldoAkun menghitung saldo akun sampai tanggal tertentu
func (s *AkunService) HitungSaldoAkun(idAkun uuid.UUID, tanggalAkhir string) (float64, error) {
	// Dapatkan akun untuk mengetahui normal saldo
	var akun models.Akun
	err := s.db.Where("id = ?", idAkun).First(&akun).Error
	if err != nil {
		return 0, errors.New("akun tidak ditemukan")
	}

	// Hitung total debit dan kredit
	type SaldoResult struct {
		TotalDebit  float64
		TotalKredit float64
	}

	var result SaldoResult
	query := s.db.Model(&models.BarisTransaksi{}).
		Select("COALESCE(SUM(jumlah_debit), 0) as total_debit, COALESCE(SUM(jumlah_kredit), 0) as total_kredit").
		Joins("JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi").
		Where("baris_transaksi.id_akun = ?", idAkun)

	if tanggalAkhir != "" {
		query = query.Where("transaksi.tanggal_transaksi <= ?", tanggalAkhir)
	}

	err = query.Scan(&result).Error
	if err != nil {
		return 0, errors.New("gagal menghitung saldo")
	}

	// Hitung saldo berdasarkan normal saldo
	var saldo float64
	if akun.NormalSaldo == "debit" {
		saldo = result.TotalDebit - result.TotalKredit
	} else {
		saldo = result.TotalKredit - result.TotalDebit
	}

	return saldo, nil
}

// InisialisasiCOADefault membuat Chart of Accounts default untuk koperasi baru
func (s *AkunService) InisialisasiCOADefault(idKoperasi uuid.UUID) error {
	// COA Default untuk Koperasi Indonesia
	coaDefault := []models.Akun{
		// ASET
		{IDKoperasi: idKoperasi, KodeAkun: "1000", NamaAkun: "ASET", TipeAkun: models.AkunAset, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "1100", NamaAkun: "Aset Lancar", TipeAkun: models.AkunAset, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "1101", NamaAkun: "Kas", TipeAkun: models.AkunAset, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "1102", NamaAkun: "Bank", TipeAkun: models.AkunAset, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "1200", NamaAkun: "Piutang", TipeAkun: models.AkunAset, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "1201", NamaAkun: "Piutang Anggota", TipeAkun: models.AkunAset, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "1300", NamaAkun: "Persediaan", TipeAkun: models.AkunAset, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "1301", NamaAkun: "Persediaan Barang Dagangan", TipeAkun: models.AkunAset, NormalSaldo: "debit"},

		// KEWAJIBAN
		{IDKoperasi: idKoperasi, KodeAkun: "2000", NamaAkun: "KEWAJIBAN", TipeAkun: models.AkunKewajiban, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "2100", NamaAkun: "Kewajiban Jangka Pendek", TipeAkun: models.AkunKewajiban, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "2101", NamaAkun: "Hutang Usaha", TipeAkun: models.AkunKewajiban, NormalSaldo: "kredit"},

		// MODAL
		{IDKoperasi: idKoperasi, KodeAkun: "3000", NamaAkun: "MODAL", TipeAkun: models.AkunModal, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "3100", NamaAkun: "Modal Koperasi", TipeAkun: models.AkunModal, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "3101", NamaAkun: "Simpanan Pokok", TipeAkun: models.AkunModal, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "3102", NamaAkun: "Simpanan Wajib", TipeAkun: models.AkunModal, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "3103", NamaAkun: "Simpanan Sukarela", TipeAkun: models.AkunModal, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "3200", NamaAkun: "Sisa Hasil Usaha (SHU)", TipeAkun: models.AkunModal, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "3201", NamaAkun: "SHU Tahun Berjalan", TipeAkun: models.AkunModal, NormalSaldo: "kredit"},

		// PENDAPATAN
		{IDKoperasi: idKoperasi, KodeAkun: "4000", NamaAkun: "PENDAPATAN", TipeAkun: models.AkunPendapatan, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "4100", NamaAkun: "Pendapatan Usaha", TipeAkun: models.AkunPendapatan, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "4101", NamaAkun: "Penjualan", TipeAkun: models.AkunPendapatan, NormalSaldo: "kredit"},
		{IDKoperasi: idKoperasi, KodeAkun: "4200", NamaAkun: "Pendapatan Lain-lain", TipeAkun: models.AkunPendapatan, NormalSaldo: "kredit"},

		// BEBAN
		{IDKoperasi: idKoperasi, KodeAkun: "5000", NamaAkun: "BEBAN", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "5100", NamaAkun: "Beban Operasional", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "5101", NamaAkun: "Beban Gaji", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "5102", NamaAkun: "Beban Listrik", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "5103", NamaAkun: "Beban Air", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "5104", NamaAkun: "Beban Telepon & Internet", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "5200", NamaAkun: "Harga Pokok Penjualan", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
		{IDKoperasi: idKoperasi, KodeAkun: "5201", NamaAkun: "HPP", TipeAkun: models.AkunBeban, NormalSaldo: "debit"},
	}

	// Insert semua akun dalam satu transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, akun := range coaDefault {
			akun.StatusAktif = true
			if err := tx.Create(&akun).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("gagal inisialisasi COA: %w", err)
	}

	return nil
}

// DapatkanHierarkiAkun mengambil struktur hierarki COA
func (s *AkunService) DapatkanHierarkiAkun(idKoperasi uuid.UUID) ([]models.AkunResponse, error) {
	var akunList []models.Akun

	// Ambil semua akun dengan preload
	err := s.db.Where("id_koperasi = ? AND status_aktif = ?", idKoperasi, true).
		Preload("AkunInduk").
		Preload("SubAkun").
		Order("kode_akun ASC").
		Find(&akunList).Error

	if err != nil {
		return nil, errors.New("gagal mengambil hierarki akun")
	}

	// Convert to response
	responses := make([]models.AkunResponse, len(akunList))
	for i, akun := range akunList {
		responses[i] = akun.ToResponse()
	}

	return responses, nil
}

// GetSemuaAkun is a wrapper for DapatkanSemuaAkun with pagination support
func (s *AkunService) GetSemuaAkun(idKoperasi uuid.UUID, tipeAkun *models.TipeAkun, statusAktif *bool) ([]models.AkunResponse, error) {
	// Convert TipeAkun pointer to string
	tipeAkunStr := ""
	if tipeAkun != nil {
		tipeAkunStr = string(*tipeAkun)
	}

	return s.DapatkanSemuaAkun(idKoperasi, tipeAkunStr, statusAktif)
}

// GetAkunByID is a wrapper for DapatkanAkun with multi-tenant validation
func (s *AkunService) GetAkunByID(idKoperasi, id uuid.UUID) (*models.AkunResponse, error) {
	// Get akun
	akun, err := s.DapatkanAkun(id)
	if err != nil {
		return nil, err
	}

	// Validate multi-tenancy - ensure akun belongs to the correct cooperative
	var akunModel models.Akun
	err = s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&akunModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akun tidak ditemukan atau tidak memiliki akses")
		}
		return nil, err
	}

	return akun, nil
}

// GetBukuBesar mengambil buku besar (ledger) untuk akun tertentu
func (s *AkunService) GetBukuBesar(idKoperasi, idAkun uuid.UUID, tanggalMulai, tanggalAkhir string) (interface{}, error) {
	// Validate akun exists and belongs to cooperative
	var akun models.Akun
	err := s.db.Where("id = ? AND id_koperasi = ?", idAkun, idKoperasi).First(&akun).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akun tidak ditemukan atau tidak memiliki akses")
		}
		return nil, err
	}

	// Query baris transaksi untuk akun ini
	type BukuBesarEntry struct {
		TanggalTransaksi time.Time `json:"tanggalTransaksi"`
		NomorJurnal      string    `json:"nomorJurnal"`
		Deskripsi        string    `json:"deskripsi"`
		JumlahDebit      float64   `json:"jumlahDebit"`
		JumlahKredit     float64   `json:"jumlahKredit"`
		Saldo            float64   `json:"saldo"`
	}

	var entries []BukuBesarEntry
	query := s.db.Table("baris_transaksi").
		Select("transaksi.tanggal_transaksi, transaksi.nomor_jurnal, transaksi.deskripsi, baris_transaksi.jumlah_debit, baris_transaksi.jumlah_kredit").
		Joins("JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi").
		Where("baris_transaksi.id_akun = ? AND transaksi.id_koperasi = ?", idAkun, idKoperasi)

	if tanggalMulai != "" {
		query = query.Where("transaksi.tanggal_transaksi >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("transaksi.tanggal_transaksi <= ?", tanggalAkhir)
	}

	err = query.Order("transaksi.tanggal_transaksi ASC, transaksi.nomor_jurnal ASC").Scan(&entries).Error
	if err != nil {
		return nil, errors.New("gagal mengambil buku besar")
	}

	// Calculate running balance
	saldo := 0.0
	for i := range entries {
		if akun.NormalSaldo == "debit" {
			saldo += entries[i].JumlahDebit - entries[i].JumlahKredit
		} else {
			saldo += entries[i].JumlahKredit - entries[i].JumlahDebit
		}
		entries[i].Saldo = saldo
	}

	return map[string]interface{}{
		"akun":    akun.ToResponse(),
		"entries": entries,
	}, nil
}
