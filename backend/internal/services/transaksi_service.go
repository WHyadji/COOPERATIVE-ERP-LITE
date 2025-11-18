package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/pkg/validasi"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// EpsilonTolerance adalah toleransi untuk perbandingan floating-point dalam rupiah
	// 0.01 = 1 sen toleransi untuk menangani masalah floating-point precision
	EpsilonTolerance = 0.01
)

// TransaksiService menangani logika bisnis transaksi akuntansi
type TransaksiService struct {
	db *gorm.DB
}

// NewTransaksiService membuat instance baru TransaksiService
func NewTransaksiService(db *gorm.DB) *TransaksiService {
	return &TransaksiService{db: db}
}

// BuatTransaksiRequest adalah struktur request untuk membuat transaksi
type BuatTransaksiRequest struct {
	TanggalTransaksi time.Time                   `json:"tanggalTransaksi" binding:"required"`
	Deskripsi        string                      `json:"deskripsi" binding:"required"`
	NomorReferensi   string                      `json:"nomorReferensi"`
	TipeTransaksi    string                      `json:"tipeTransaksi"`
	BarisTransaksi   []BuatBarisTransaksiRequest `json:"barisTransaksi" binding:"required,min=2"`
}

// BuatBarisTransaksiRequest adalah struktur untuk baris transaksi
type BuatBarisTransaksiRequest struct {
	IDAkun       uuid.UUID `json:"idAkun" binding:"required"`
	JumlahDebit  float64   `json:"jumlahDebit"`
	JumlahKredit float64   `json:"jumlahKredit"`
	Keterangan   string    `json:"keterangan"`
}

// BuatTransaksi membuat jurnal entry baru dengan validasi double-entry
func (s *TransaksiService) BuatTransaksi(idKoperasi, idPengguna uuid.UUID, req *BuatTransaksiRequest) (*models.TransaksiResponse, error) {
	// Initialize validator
	validator := validasi.Baru()

	// Validasi business logic
	if err := validator.TanggalTransaksi(req.TanggalTransaksi); err != nil {
		return nil, err
	}

	if err := validator.TeksWajib(req.Deskripsi, "deskripsi", 5, 500); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.NomorReferensi, "nomor referensi", 50); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.TipeTransaksi, "tipe transaksi", 50); err != nil {
		return nil, err
	}

	// Validasi baris transaksi (debit = kredit)
	if err := s.ValidasiTransaksi(req.BarisTransaksi); err != nil {
		return nil, err
	}

	// Hitung total debit dan kredit
	var totalDebit, totalKredit float64
	for _, baris := range req.BarisTransaksi {
		totalDebit += baris.JumlahDebit
		totalKredit += baris.JumlahKredit
	}

	// Buat transaksi dengan baris-barisnya dalam satu transaction
	// IMPORTANT: GenerateNomorJurnal is now called INSIDE the transaction to prevent race conditions
	var transaksi models.Transaksi

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Generate nomor jurnal INSIDE transaction with row-level locking
		// This prevents race conditions when multiple requests come simultaneously
		nomorJurnal, err := s.generateNomorJurnalInTx(tx, idKoperasi, req.TanggalTransaksi)
		if err != nil {
			return err
		}

		// Buat header transaksi
		transaksi = models.Transaksi{
			IDKoperasi:       idKoperasi,
			NomorJurnal:      nomorJurnal,
			TanggalTransaksi: req.TanggalTransaksi,
			Deskripsi:        req.Deskripsi,
			NomorReferensi:   req.NomorReferensi,
			TipeTransaksi:    req.TipeTransaksi,
			TotalDebit:       totalDebit,
			TotalKredit:      totalKredit,
			StatusBalanced:   true,
			DibuatOleh:       idPengguna,
		}

		if err := tx.Create(&transaksi).Error; err != nil {
			return errors.New("gagal membuat transaksi")
		}

		// Buat baris transaksi
		for _, barisReq := range req.BarisTransaksi {
			// Validasi akun exists
			var akun models.Akun
			if err := tx.Where("id = ? AND id_koperasi = ?", barisReq.IDAkun, idKoperasi).First(&akun).Error; err != nil {
				return fmt.Errorf("akun %s tidak ditemukan", barisReq.IDAkun)
			}

			baris := models.BarisTransaksi{
				IDTransaksi:  transaksi.ID,
				IDAkun:       barisReq.IDAkun,
				JumlahDebit:  barisReq.JumlahDebit,
				JumlahKredit: barisReq.JumlahKredit,
				Keterangan:   barisReq.Keterangan,
			}

			if err := tx.Create(&baris).Error; err != nil {
				return errors.New("gagal membuat baris transaksi")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload dengan baris transaksi
	s.db.Preload("BarisTransaksi.Akun").First(&transaksi, transaksi.ID)

	response := transaksi.ToResponse()
	return &response, nil
}

// ValidasiTransaksi memvalidasi bahwa total debit = total kredit
func (s *TransaksiService) ValidasiTransaksi(barisTransaksi []BuatBarisTransaksiRequest) error {
	validator := validasi.Baru()

	if len(barisTransaksi) < 2 {
		return errors.New("transaksi harus memiliki minimal 2 baris (debit dan kredit)")
	}

	var totalDebit, totalKredit float64
	hasDebit := false
	hasKredit := false

	for i, baris := range barisTransaksi {
		// Validasi tidak boleh debit dan kredit bersamaan
		if baris.JumlahDebit > 0 && baris.JumlahKredit > 0 {
			return errors.New("satu baris tidak boleh memiliki debit dan kredit sekaligus")
		}

		// Validasi minimal salah satu harus ada
		if baris.JumlahDebit == 0 && baris.JumlahKredit == 0 {
			return errors.New("setiap baris harus memiliki nilai debit atau kredit")
		}

		// Validasi jumlah debit jika ada
		if baris.JumlahDebit > 0 {
			if err := validator.Jumlah(baris.JumlahDebit, fmt.Sprintf("jumlah debit baris ke-%d", i+1)); err != nil {
				return err
			}
			hasDebit = true
		}

		// Validasi jumlah kredit jika ada
		if baris.JumlahKredit > 0 {
			if err := validator.Jumlah(baris.JumlahKredit, fmt.Sprintf("jumlah kredit baris ke-%d", i+1)); err != nil {
				return err
			}
			hasKredit = true
		}

		// Validasi keterangan (opsional)
		if err := validator.TeksOpsional(baris.Keterangan, fmt.Sprintf("keterangan baris ke-%d", i+1), 500); err != nil {
			return err
		}

		totalDebit += baris.JumlahDebit
		totalKredit += baris.JumlahKredit
	}

	// Validasi ada debit dan kredit
	if !hasDebit || !hasKredit {
		return errors.New("transaksi harus memiliki minimal satu baris debit dan satu baris kredit")
	}

	// Validasi total tidak boleh nol (menggunakan epsilon untuk floating-point)
	if totalDebit < EpsilonTolerance {
		return errors.New("total debit dan kredit tidak boleh 0")
	}

	// Validasi balanced (debit = kredit) dengan epsilon tolerance untuk floating-point precision
	// Menggunakan math.Abs untuk menghitung selisih absolut dan membandingkan dengan epsilon
	// Ini mengatasi masalah floating-point arithmetic seperti: 0.1 + 0.1 + 0.1 != 0.3
	if math.Abs(totalDebit-totalKredit) > EpsilonTolerance {
		return fmt.Errorf("total debit (%.2f) tidak sama dengan total kredit (%.2f)", totalDebit, totalKredit)
	}

	return nil
}

// generateNomorJurnalInTx menghasilkan nomor jurnal otomatis dalam transaction yang sudah ada
// Format: JRN-YYYYMMDD-NNNN (contoh: JRN-20250116-0001)
// MUST be called within an existing transaction to maintain row-level locking throughout
// the entire transaction lifecycle, preventing race conditions in concurrent requests
//
// Uses PostgreSQL advisory lock to prevent race conditions even when no rows exist yet.
// This solves the "first transaction of the day" problem where SELECT FOR UPDATE
// cannot lock non-existent rows.
func (s *TransaksiService) generateNomorJurnalInTx(tx *gorm.DB, idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	tanggalStr := tanggal.Format("20060102")
	tanggalDate := tanggal.Format("2006-01-02")

	// CRITICAL: Use PostgreSQL advisory lock to prevent race conditions
	// Advisory lock works even when no rows exist (solves first-transaction-of-day problem)
	// Lock ID is derived from hash of (koperasi_id + date) to ensure uniqueness per koperasi per day
	// Lock is automatically released when transaction commits/rolls back
	lockKey := generateAdvisoryLockKey(idKoperasi, tanggalStr)
	if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
		return "", fmt.Errorf("gagal acquire advisory lock: %w", err)
	}

	// Now safely get the last journal number for this date
	// Even if no row exists, we're protected by the advisory lock above
	var lastTransaksi models.Transaksi
	err := tx.Where("id_koperasi = ? AND DATE(tanggal_transaksi) = ?", idKoperasi, tanggalDate).
		Order("nomor_jurnal DESC").
		Limit(1).
		First(&lastTransaksi).Error

	nomorUrut := 1

	// If there's a previous transaction, parse and increment
	if err == nil && lastTransaksi.NomorJurnal != "" {
		// Extract number from JRN-20250116-0001
		var parsedTanggal string
		var parsedUrut int
		_, scanErr := fmt.Sscanf(lastTransaksi.NomorJurnal, "JRN-%s-%04d", &parsedTanggal, &parsedUrut)
		if scanErr == nil && parsedTanggal == tanggalStr {
			nomorUrut = parsedUrut + 1
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("gagal membaca transaksi terakhir: %w", err)
	}

	nomorJurnal := fmt.Sprintf("JRN-%s-%04d", tanggalStr, nomorUrut)
	return nomorJurnal, nil
}

// generateAdvisoryLockKey generates a consistent int64 lock key from koperasi ID and date
// This ensures the same key is generated for the same koperasi and date across all instances
func generateAdvisoryLockKey(idKoperasi uuid.UUID, tanggalStr string) int64 {
	// Combine koperasi ID and date to create unique lock key
	// Use first 8 bytes of UUID XOR with hash of date string
	h := int64(0)

	// XOR first 8 bytes of UUID
	uuidBytes := [16]byte(idKoperasi)
	for i := 0; i < 8; i++ {
		h = (h << 8) | int64(uuidBytes[i])
	}

	// XOR with simple hash of date string
	for i := 0; i < len(tanggalStr); i++ {
		h ^= int64(tanggalStr[i]) << uint(i%8*8)
	}

	return h
}

// GenerateNomorJurnal is a public wrapper for backward compatibility
// It creates its own transaction if called outside of an existing transaction
// WARNING: This method is deprecated for use in BuatTransaksi to avoid race conditions
// Prefer using generateNomorJurnalInTx within an existing transaction instead
func (s *TransaksiService) GenerateNomorJurnal(idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	var nomorJurnal string
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		nomorJurnal, err = s.generateNomorJurnalInTx(tx, idKoperasi, tanggal)
		return err
	})

	if err != nil {
		return "", errors.New("gagal generate nomor jurnal")
	}

	return nomorJurnal, nil
}

// DapatkanSemuaTransaksi mengambil daftar transaksi dengan filter
func (s *TransaksiService) DapatkanSemuaTransaksi(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir, tipeTransaksi string, page, pageSize int) ([]models.TransaksiResponse, int64, error) {
	var transaksiList []models.Transaksi
	var total int64

	query := s.db.Model(&models.Transaksi{}).Where("id_koperasi = ?", idKoperasi)

	// Apply filters
	if tanggalMulai != "" {
		query = query.Where("tanggal_transaksi >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("tanggal_transaksi <= ?", tanggalAkhir)
	}
	if tipeTransaksi != "" {
		query = query.Where("tipe_transaksi = ?", tipeTransaksi)
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("tanggal_transaksi DESC, nomor_jurnal DESC").
		Preload("BarisTransaksi.Akun").
		Find(&transaksiList).Error

	if err != nil {
		return nil, 0, errors.New("gagal mengambil daftar transaksi")
	}

	// Convert to response
	responses := make([]models.TransaksiResponse, len(transaksiList))
	for i, transaksi := range transaksiList {
		responses[i] = transaksi.ToResponse()
	}

	return responses, total, nil
}

// DapatkanTransaksi mengambil transaksi berdasarkan ID
func (s *TransaksiService) DapatkanTransaksi(id uuid.UUID) (*models.TransaksiResponse, error) {
	var transaksi models.Transaksi
	err := s.db.Preload("BarisTransaksi.Akun").Where("id = ?", id).First(&transaksi).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaksi tidak ditemukan")
		}
		return nil, err
	}

	response := transaksi.ToResponse()
	return &response, nil
}

// HapusTransaksi menghapus transaksi (soft delete)
func (s *TransaksiService) HapusTransaksi(id uuid.UUID) error {
	// Soft delete transaksi (baris transaksi akan cascade delete)
	err := s.db.Delete(&models.Transaksi{}, id).Error
	if err != nil {
		return errors.New("gagal menghapus transaksi")
	}

	return nil
}

// DapatkanBukuBesar mengambil buku besar (ledger) untuk akun tertentu
func (s *TransaksiService) DapatkanBukuBesar(idAkun uuid.UUID, tanggalMulai, tanggalAkhir string) ([]map[string]interface{}, error) {
	// Validasi akun exists
	var akun models.Akun
	err := s.db.Where("id = ?", idAkun).First(&akun).Error
	if err != nil {
		return nil, errors.New("akun tidak ditemukan")
	}

	// Query baris transaksi untuk akun ini
	var barisTransaksiList []models.BarisTransaksi
	query := s.db.Preload("Transaksi").Where("id_akun = ?", idAkun).
		Joins("JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi")

	if tanggalMulai != "" {
		query = query.Where("transaksi.tanggal_transaksi >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("transaksi.tanggal_transaksi <= ?", tanggalAkhir)
	}

	err = query.Order("transaksi.tanggal_transaksi ASC").Find(&barisTransaksiList).Error
	if err != nil {
		return nil, errors.New("gagal mengambil buku besar")
	}

	// Format response dengan running balance
	var saldo float64
	ledger := make([]map[string]interface{}, 0)

	for _, baris := range barisTransaksiList {
		// Update saldo
		if akun.NormalSaldo == "debit" {
			saldo += baris.JumlahDebit - baris.JumlahKredit
		} else {
			saldo += baris.JumlahKredit - baris.JumlahDebit
		}

		entry := map[string]interface{}{
			"tanggal":     baris.Transaksi.TanggalTransaksi,
			"nomorJurnal": baris.Transaksi.NomorJurnal,
			"deskripsi":   baris.Transaksi.Deskripsi,
			"keterangan":  baris.Keterangan,
			"debit":       baris.JumlahDebit,
			"kredit":      baris.JumlahKredit,
			"saldo":       saldo,
		}
		ledger = append(ledger, entry)
	}

	return ledger, nil
}

// PostingOtomatisSimpanan membuat jurnal otomatis untuk setoran simpanan
func (s *TransaksiService) PostingOtomatisSimpanan(idKoperasi, idPengguna, idSimpanan uuid.UUID) error {
	// Ambil data simpanan
	var simpanan models.Simpanan
	err := s.db.Where("id = ?", idSimpanan).First(&simpanan).Error
	if err != nil {
		return errors.New("simpanan tidak ditemukan")
	}

	// Tentukan akun modal berdasarkan tipe simpanan
	var kodeAkunModal string
	switch simpanan.TipeSimpanan {
	case models.SimpananPokok:
		kodeAkunModal = "3101" // Simpanan Pokok
	case models.SimpananWajib:
		kodeAkunModal = "3102" // Simpanan Wajib
	case models.SimpananSukarela:
		kodeAkunModal = "3103" // Simpanan Sukarela
	}

	// Dapatkan akun kas dan akun modal
	var akunKas, akunModal models.Akun
	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1101").First(&akunKas).Error; err != nil {
		return errors.New("akun kas tidak ditemukan")
	}
	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, kodeAkunModal).First(&akunModal).Error; err != nil {
		return errors.New("akun modal tidak ditemukan")
	}

	// Buat jurnal entry
	req := &BuatTransaksiRequest{
		TanggalTransaksi: simpanan.TanggalTransaksi,
		Deskripsi:        fmt.Sprintf("Setoran %s", simpanan.TipeSimpanan),
		NomorReferensi:   simpanan.NomorReferensi,
		TipeTransaksi:    "simpanan",
		BarisTransaksi: []BuatBarisTransaksiRequest{
			{
				IDAkun:      akunKas.ID,
				JumlahDebit: simpanan.JumlahSetoran,
				Keterangan:  fmt.Sprintf("Setoran %s", simpanan.TipeSimpanan),
			},
			{
				IDAkun:       akunModal.ID,
				JumlahKredit: simpanan.JumlahSetoran,
				Keterangan:   fmt.Sprintf("Setoran %s", simpanan.TipeSimpanan),
			},
		},
	}

	transaksi, err := s.BuatTransaksi(idKoperasi, idPengguna, req)
	if err != nil {
		return fmt.Errorf("gagal posting simpanan: %w", err)
	}

	// Update simpanan dengan ID transaksi
	simpanan.IDTransaksi = &transaksi.ID
	s.db.Save(&simpanan)

	return nil
}

// PostingOtomatisPenjualan membuat jurnal otomatis untuk penjualan
func (s *TransaksiService) PostingOtomatisPenjualan(idKoperasi, idPengguna, idPenjualan uuid.UUID) error {
	// Ambil data penjualan dengan items
	var penjualan models.Penjualan
	err := s.db.Preload("ItemPenjualan.Produk").Where("id = ?", idPenjualan).First(&penjualan).Error
	if err != nil {
		return errors.New("penjualan tidak ditemukan")
	}

	// Dapatkan akun-akun yang diperlukan
	var akunKas, akunPenjualan, akunHPP, akunPersediaan models.Akun
	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1101").First(&akunKas).Error; err != nil {
		return errors.New("akun kas tidak ditemukan")
	}
	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "4101").First(&akunPenjualan).Error; err != nil {
		return errors.New("akun penjualan tidak ditemukan")
	}
	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "5201").First(&akunHPP).Error; err != nil {
		return errors.New("akun HPP tidak ditemukan")
	}
	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1301").First(&akunPersediaan).Error; err != nil {
		return errors.New("akun persediaan tidak ditemukan")
	}

	// Hitung total HPP
	var totalHPP float64
	for _, item := range penjualan.ItemPenjualan {
		totalHPP += item.Produk.HargaBeli * float64(item.Kuantitas)
	}

	// Buat baris transaksi
	barisTransaksi := []BuatBarisTransaksiRequest{
		// Kas bertambah (debit)
		{
			IDAkun:      akunKas.ID,
			JumlahDebit: penjualan.TotalBelanja,
			Keterangan:  "Penerimaan kas dari penjualan",
		},
		// Penjualan bertambah (kredit)
		{
			IDAkun:       akunPenjualan.ID,
			JumlahKredit: penjualan.TotalBelanja,
			Keterangan:   "Pendapatan penjualan",
		},
	}

	// Jika ada HPP, tambahkan jurnal HPP
	if totalHPP > 0 {
		barisTransaksi = append(barisTransaksi,
			BuatBarisTransaksiRequest{
				IDAkun:      akunHPP.ID,
				JumlahDebit: totalHPP,
				Keterangan:  "Harga Pokok Penjualan",
			},
			BuatBarisTransaksiRequest{
				IDAkun:       akunPersediaan.ID,
				JumlahKredit: totalHPP,
				Keterangan:   "Pengurangan persediaan",
			},
		)
	}

	// Buat jurnal entry
	req := &BuatTransaksiRequest{
		TanggalTransaksi: penjualan.TanggalPenjualan,
		Deskripsi:        fmt.Sprintf("Penjualan %s", penjualan.NomorPenjualan),
		NomorReferensi:   penjualan.NomorPenjualan,
		TipeTransaksi:    "penjualan",
		BarisTransaksi:   barisTransaksi,
	}

	transaksi, err := s.BuatTransaksi(idKoperasi, idPengguna, req)
	if err != nil {
		return fmt.Errorf("gagal posting penjualan: %w", err)
	}

	// Update penjualan dengan ID transaksi
	penjualan.IDTransaksi = &transaksi.ID
	s.db.Save(&penjualan)

	return nil
}
