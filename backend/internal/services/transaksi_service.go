package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/utils"
	"errors"
	"fmt"
	"time"

	apperrors "cooperative-erp-lite/internal/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TransaksiService menangani logika bisnis transaksi akuntansi
type TransaksiService struct {
	db     *gorm.DB
	logger *utils.Logger
}

// NewTransaksiService membuat instance baru TransaksiService
func NewTransaksiService(db *gorm.DB) *TransaksiService {
	return &TransaksiService{
		db:     db,
		logger: utils.NewLogger("TransaksiService"),
	}
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
	const method = "BuatTransaksi"

	// Validasi baris transaksi (debit = kredit)
	if err := s.ValidasiTransaksi(req.BarisTransaksi); err != nil {
		s.logger.Error(method, "Validasi transaksi gagal", err, map[string]interface{}{
			"koperasi_id":  idKoperasi.String(),
			"jumlah_baris": len(req.BarisTransaksi),
		})
		return nil, err
	}

	// Generate nomor jurnal
	nomorJurnal, err := s.GenerateNomorJurnal(idKoperasi, req.TanggalTransaksi)
	if err != nil {
		s.logger.Error(method, "Gagal generate nomor jurnal", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return nil, err
	}

	// Hitung total debit dan kredit
	var totalDebit, totalKredit float64
	for _, baris := range req.BarisTransaksi {
		totalDebit += baris.JumlahDebit
		totalKredit += baris.JumlahKredit
	}

	// Buat transaksi dengan baris-barisnya dalam satu transaction
	var transaksi models.Transaksi

	err = s.db.Transaction(func(tx *gorm.DB) error {
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
			s.logger.Error(method, "Gagal membuat header transaksi di database", err, map[string]interface{}{
				"koperasi_id":  idKoperasi.String(),
				"nomor_jurnal": nomorJurnal,
				"total_debit":  totalDebit,
				"total_kredit": totalKredit,
			})
			return utils.WrapDatabaseError(err, "Gagal membuat transaksi")
		}

		// Buat baris transaksi
		for _, barisReq := range req.BarisTransaksi {
			// Validasi akun exists
			var akun models.Akun
			if err := tx.Where("id = ? AND id_koperasi = ?", barisReq.IDAkun, idKoperasi).First(&akun).Error; err != nil {
				s.logger.Error(method, "Akun tidak ditemukan saat membuat transaksi", err, map[string]interface{}{
					"koperasi_id": idKoperasi.String(),
					"akun_id":     barisReq.IDAkun.String(),
				})
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.WrapDatabaseError(err, "Akun")
				}
				return utils.WrapDatabaseError(err, "Gagal mengambil data akun")
			}

			baris := models.BarisTransaksi{
				IDTransaksi:  transaksi.ID,
				IDAkun:       barisReq.IDAkun,
				JumlahDebit:  barisReq.JumlahDebit,
				JumlahKredit: barisReq.JumlahKredit,
				Keterangan:   barisReq.Keterangan,
			}

			if err := tx.Create(&baris).Error; err != nil {
				s.logger.Error(method, "Gagal membuat baris transaksi", err, map[string]interface{}{
					"transaksi_id": transaksi.ID.String(),
					"akun_id":      barisReq.IDAkun.String(),
				})
				return utils.WrapDatabaseError(err, "Gagal membuat baris transaksi")
			}
		}

		return nil
	})

	if err != nil {
		s.logger.Error(method, "Transaksi database gagal", err, map[string]interface{}{
			"koperasi_id":  idKoperasi.String(),
			"nomor_jurnal": nomorJurnal,
		})
		return nil, err
	}

	// Reload dengan baris transaksi
	if err := s.db.Preload("BarisTransaksi.Akun").First(&transaksi, transaksi.ID).Error; err != nil {
		s.logger.Error(method, "Gagal reload data transaksi setelah berhasil", err, map[string]interface{}{
			"transaksi_id": transaksi.ID.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data transaksi")
	}

	s.logger.Info(method, "Berhasil membuat transaksi", map[string]interface{}{
		"transaksi_id": transaksi.ID.String(),
		"nomor_jurnal": nomorJurnal,
		"total_debit":  totalDebit,
		"total_kredit": totalKredit,
		"jumlah_baris": len(req.BarisTransaksi),
		"koperasi_id":  idKoperasi.String(),
	})

	respons := transaksi.ToResponse()
	return &respons, nil
}

// ValidasiTransaksi memvalidasi bahwa total debit = total kredit
func (s *TransaksiService) ValidasiTransaksi(barisTransaksi []BuatBarisTransaksiRequest) error {
	if len(barisTransaksi) < 2 {
		return errors.New("transaksi harus memiliki minimal 2 baris (debit dan kredit)")
	}

	var totalDebit, totalKredit float64
	hasDebit := false
	hasKredit := false

	for _, baris := range barisTransaksi {
		// Validasi tidak boleh debit dan kredit bersamaan
		if baris.JumlahDebit > 0 && baris.JumlahKredit > 0 {
			return apperrors.ErrDebitKreditKeduanya
		}

		// Validasi minimal salah satu harus ada
		if baris.JumlahDebit == 0 && baris.JumlahKredit == 0 {
			return errors.New("setiap baris harus memiliki nilai debit atau kredit")
		}

		totalDebit += baris.JumlahDebit
		totalKredit += baris.JumlahKredit

		if baris.JumlahDebit > 0 {
			hasDebit = true
		}
		if baris.JumlahKredit > 0 {
			hasKredit = true
		}
	}

	// Validasi ada debit dan kredit
	if !hasDebit || !hasKredit {
		return errors.New("transaksi harus memiliki minimal satu baris debit dan satu baris kredit")
	}

	// Validasi balanced (debit = kredit)
	if totalDebit != totalKredit {
		return apperrors.ErrDebitKreditTidakBalance
	}

	return nil
}

// GenerateNomorJurnal menghasilkan nomor jurnal otomatis
// Format: JRN-YYYYMMDD-NNNN (contoh: JRN-20250116-0001)
// Menggunakan row-level locking untuk mencegah race condition pada concurrent requests
func (s *TransaksiService) GenerateNomorJurnal(idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	const method = "GenerateNomorJurnal"

	tanggalStr := tanggal.Format("20060102")
	tanggalDate := tanggal.Format("2006-01-02")
	var nomorJurnal string

	// Gunakan transaction dengan row-level locking untuk mencegah race condition
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Lock dan ambil nomor jurnal terakhir untuk tanggal ini
		var lastTransaksi models.Transaksi
		err := tx.Where("id_koperasi = ? AND DATE(tanggal_transaksi) = ?", idKoperasi, tanggalDate).
			Order("nomor_jurnal DESC").
			Limit(1).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastTransaksi).Error

		nomorUrut := 1

		// Jika ada transaksi sebelumnya, parse dan increment
		if err == nil && lastTransaksi.NomorJurnal != "" {
			// Extract number dari JRN-20250116-0001
			var parsedTanggal string
			var parsedUrut int
			_, scanErr := fmt.Sscanf(lastTransaksi.NomorJurnal, "JRN-%s-%04d", &parsedTanggal, &parsedUrut)
			if scanErr == nil && parsedTanggal == tanggalStr {
				nomorUrut = parsedUrut + 1
			}
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Gagal mengambil nomor jurnal terakhir", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"tanggal":     tanggalDate,
			})
			return utils.WrapDatabaseError(err, "Gagal mengambil nomor jurnal terakhir")
		}

		nomorJurnal = fmt.Sprintf("JRN-%s-%04d", tanggalStr, nomorUrut)
		return nil
	})

	if err != nil {
		s.logger.Error(method, "Transaksi generate nomor jurnal gagal", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"tanggal":     tanggalDate,
		})
		return "", err
	}

	s.logger.Debug(method, "Berhasil generate nomor jurnal", map[string]interface{}{
		"nomor_jurnal": nomorJurnal,
		"koperasi_id":  idKoperasi.String(),
	})

	return nomorJurnal, nil
}

// DapatkanSemuaTransaksi mengambil daftar transaksi dengan filter
func (s *TransaksiService) DapatkanSemuaTransaksi(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir, tipeTransaksi string, page, pageSize int) ([]models.TransaksiResponse, int64, error) {
	const method = "DapatkanSemuaTransaksi"

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
	if err := query.Count(&total).Error; err != nil {
		s.logger.Error(method, "Gagal menghitung total transaksi", err, map[string]interface{}{
			"koperasi_id":    idKoperasi.String(),
			"tanggal_mulai":  tanggalMulai,
			"tanggal_akhir":  tanggalAkhir,
			"tipe_transaksi": tipeTransaksi,
		})
		return nil, 0, utils.WrapDatabaseError(err, "Gagal menghitung total transaksi")
	}

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("tanggal_transaksi DESC, nomor_jurnal DESC").
		Preload("BarisTransaksi.Akun").
		Find(&transaksiList).Error

	if err != nil {
		s.logger.Error(method, "Gagal mengambil daftar transaksi", err, map[string]interface{}{
			"koperasi_id":    idKoperasi.String(),
			"tanggal_mulai":  tanggalMulai,
			"tanggal_akhir":  tanggalAkhir,
			"tipe_transaksi": tipeTransaksi,
			"page":           page,
			"page_size":      pageSize,
		})
		return nil, 0, utils.WrapDatabaseError(err, "Gagal mengambil daftar transaksi")
	}

	// Convert to response
	responseDaftar := make([]models.TransaksiResponse, len(transaksiList))
	for i, transaksi := range transaksiList {
		responseDaftar[i] = transaksi.ToResponse()
	}

	s.logger.Debug(method, "Berhasil mengambil daftar transaksi", map[string]interface{}{
		"koperasi_id": idKoperasi.String(),
		"total":       total,
		"jumlah":      len(transaksiList),
		"page":        page,
	})

	return responseDaftar, total, nil
}

// DapatkanTransaksi mengambil transaksi berdasarkan ID
func (s *TransaksiService) DapatkanTransaksi(idKoperasi, id uuid.UUID) (*models.TransaksiResponse, error) {
	const method = "DapatkanTransaksi"

	var transaksi models.Transaksi
	err := s.db.Preload("BarisTransaksi.Akun").Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&transaksi).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Transaksi tidak ditemukan atau tidak memiliki akses", err, map[string]interface{}{
				"transaksi_id": id.String(),
				"koperasi_id":  idKoperasi.String(),
			})
			return nil, apperrors.ErrTransaksiTidakDitemukan
		}
		s.logger.Error(method, "Gagal mengambil data transaksi", err, map[string]interface{}{
			"transaksi_id": id.String(),
			"koperasi_id":  idKoperasi.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data transaksi")
	}

	s.logger.Debug(method, "Berhasil mengambil data transaksi", map[string]interface{}{
		"transaksi_id": id.String(),
		"koperasi_id":  idKoperasi.String(),
		"nomor_jurnal": transaksi.NomorJurnal,
		"jumlah_baris": len(transaksi.BarisTransaksi),
	})

	respons := transaksi.ToResponse()
	return &respons, nil
}

// PerbaruiTransaksi memperbarui transaksi yang sudah ada
func (s *TransaksiService) PerbaruiTransaksi(idKoperasi, id uuid.UUID, req *BuatTransaksiRequest) (*models.TransaksiResponse, error) {
	const method = "PerbaruiTransaksi"

	// Validasi baris transaksi (debit = kredit)
	if err := s.ValidasiTransaksi(req.BarisTransaksi); err != nil {
		s.logger.Error(method, "Validasi transaksi gagal", err, map[string]interface{}{
			"koperasi_id":  idKoperasi.String(),
			"transaksi_id": id.String(),
			"jumlah_baris": len(req.BarisTransaksi),
		})
		return nil, err
	}

	// Ambil transaksi yang akan diupdate
	var existingTransaksi models.Transaksi
	err := s.db.Where("id = ?", id).First(&existingTransaksi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Transaksi tidak ditemukan", err, map[string]interface{}{
				"koperasi_id":  idKoperasi.String(),
				"transaksi_id": id.String(),
			})
			return nil, apperrors.ErrTransaksiTidakDitemukan
		}
		s.logger.Error(method, "Gagal mengambil data transaksi", err, map[string]interface{}{
			"transaksi_id": id.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data transaksi")
	}

	// SECURITY: Validate multi-tenant access - ensure transaction belongs to the cooperative
	if existingTransaksi.IDKoperasi != idKoperasi {
		s.logger.Error(method, "Akses ditolak: transaksi bukan milik koperasi", nil, map[string]interface{}{
			"koperasi_id":           idKoperasi.String(),
			"transaksi_id":          id.String(),
			"koperasi_id_transaksi": existingTransaksi.IDKoperasi.String(),
		})
		return nil, apperrors.ErrTransaksiTidakDitemukan
	}

	// TODO: Add check for "posted" status when the field is added to the model
	// For now, we allow updates to all transactions

	// Hitung total debit dan kredit
	var totalDebit, totalKredit float64
	for _, baris := range req.BarisTransaksi {
		totalDebit += baris.JumlahDebit
		totalKredit += baris.JumlahKredit
	}

	// Update transaksi dalam transaction
	var transaksi models.Transaksi

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Hapus baris transaksi yang lama
		if err := tx.Where("id_transaksi = ?", id).Delete(&models.BarisTransaksi{}).Error; err != nil {
			s.logger.Error(method, "Gagal menghapus baris transaksi lama", err, map[string]interface{}{
				"transaksi_id": id.String(),
			})
			return utils.WrapDatabaseError(err, "Gagal menghapus baris transaksi lama")
		}

		// Update header transaksi
		updates := map[string]interface{}{
			"tanggal_transaksi": req.TanggalTransaksi,
			"deskripsi":         req.Deskripsi,
			"nomor_referensi":   req.NomorReferensi,
			"tipe_transaksi":    req.TipeTransaksi,
			"total_debit":       totalDebit,
			"total_kredit":      totalKredit,
			"status_balanced":   true,
		}

		if err := tx.Model(&models.Transaksi{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			s.logger.Error(method, "Gagal memperbarui header transaksi", err, map[string]interface{}{
				"transaksi_id": id.String(),
			})
			return utils.WrapDatabaseError(err, "Gagal memperbarui transaksi")
		}

		// Buat baris transaksi yang baru
		for _, barisReq := range req.BarisTransaksi {
			baris := models.BarisTransaksi{
				IDTransaksi:  id,
				IDAkun:       barisReq.IDAkun,
				JumlahDebit:  barisReq.JumlahDebit,
				JumlahKredit: barisReq.JumlahKredit,
				Keterangan:   barisReq.Keterangan,
			}

			if err := tx.Create(&baris).Error; err != nil {
				s.logger.Error(method, "Gagal membuat baris transaksi baru", err, map[string]interface{}{
					"transaksi_id": id.String(),
					"akun_id":      barisReq.IDAkun.String(),
				})
				return utils.WrapDatabaseError(err, "Gagal membuat baris transaksi baru")
			}
		}

		// Load transaksi yang sudah diupdate dengan semua relasi
		if err := tx.Preload("BarisTransaksi.Akun").Where("id = ?", id).First(&transaksi).Error; err != nil {
			s.logger.Error(method, "Gagal mengambil transaksi yang diupdate", err, map[string]interface{}{
				"transaksi_id": id.String(),
			})
			return utils.WrapDatabaseError(err, "Gagal mengambil transaksi yang diupdate")
		}

		return nil
	})

	if err != nil {
		s.logger.Error(method, "Transaksi database gagal", err, map[string]interface{}{
			"koperasi_id":  idKoperasi.String(),
			"transaksi_id": id.String(),
		})
		return nil, err
	}

	s.logger.Info(method, "Berhasil memperbarui transaksi", map[string]interface{}{
		"transaksi_id": id.String(),
		"koperasi_id":  idKoperasi.String(),
		"total_debit":  totalDebit,
		"total_kredit": totalKredit,
		"jumlah_baris": len(req.BarisTransaksi),
	})

	respons := transaksi.ToResponse()
	return &respons, nil
}

// HapusTransaksi menghapus transaksi (soft delete) dengan validasi multi-tenant
func (s *TransaksiService) HapusTransaksi(idKoperasi, id uuid.UUID) error {
	const method = "HapusTransaksi"

	// Soft delete transaksi dengan validasi multi-tenant (baris transaksi akan cascade delete)
	result := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).Delete(&models.Transaksi{})
	if result.Error != nil {
		s.logger.Error(method, "Gagal menghapus transaksi", result.Error, map[string]interface{}{
			"transaksi_id": id.String(),
			"koperasi_id":  idKoperasi.String(),
		})
		return utils.WrapDatabaseError(result.Error, "Gagal menghapus transaksi")
	}

	if result.RowsAffected == 0 {
		s.logger.Error(method, "Transaksi tidak ditemukan atau tidak memiliki akses", nil, map[string]interface{}{
			"transaksi_id": id.String(),
			"koperasi_id":  idKoperasi.String(),
		})
		return apperrors.ErrTransaksiTidakDitemukan
	}

	s.logger.Info(method, "Berhasil menghapus transaksi", map[string]interface{}{
		"transaksi_id": id.String(),
		"koperasi_id":  idKoperasi.String(),
	})

	return nil
}

// ReverseTransaksi membuat jurnal pembalik (reversing entry) untuk membatalkan transaksi
func (s *TransaksiService) ReverseTransaksi(idKoperasi, idPengguna, id uuid.UUID, keterangan string) (*models.TransaksiResponse, error) {
	const method = "ReverseTransaksi"

	// Ambil transaksi original yang akan di-reverse
	var originalTransaksi models.Transaksi
	err := s.db.Preload("BarisTransaksi").Where("id = ?", id).First(&originalTransaksi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Transaksi tidak ditemukan", err, map[string]interface{}{
				"koperasi_id":  idKoperasi.String(),
				"transaksi_id": id.String(),
			})
			return nil, apperrors.ErrTransaksiTidakDitemukan
		}
		s.logger.Error(method, "Gagal mengambil data transaksi", err, map[string]interface{}{
			"transaksi_id": id.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data transaksi")
	}

	// SECURITY: Validate multi-tenant access - ensure transaction belongs to the cooperative
	if originalTransaksi.IDKoperasi != idKoperasi {
		s.logger.Error(method, "Akses ditolak: transaksi bukan milik koperasi", nil, map[string]interface{}{
			"koperasi_id":           idKoperasi.String(),
			"transaksi_id":          id.String(),
			"koperasi_id_transaksi": originalTransaksi.IDKoperasi.String(),
		})
		return nil, apperrors.ErrTransaksiTidakDitemukan
	}

	// Validate that there are transaction lines to reverse
	if len(originalTransaksi.BarisTransaksi) == 0 {
		s.logger.Error(method, "Tidak ada baris transaksi untuk direverse", nil, map[string]interface{}{
			"transaksi_id": id.String(),
		})
		return nil, apperrors.ErrTidakAdaBarisTransaksi
	}

	// Generate nomor jurnal untuk transaksi pembalik
	tanggalReverse := time.Now()
	nomorJurnal, err := s.GenerateNomorJurnal(idKoperasi, tanggalReverse)
	if err != nil {
		s.logger.Error(method, "Gagal generate nomor jurnal untuk reversal", err, map[string]interface{}{
			"koperasi_id":  idKoperasi.String(),
			"transaksi_id": id.String(),
		})
		return nil, err
	}

	// Buat deskripsi untuk jurnal pembalik
	deskripsiReverse := fmt.Sprintf("REVERSAL: %s", originalTransaksi.Deskripsi)
	if keterangan != "" {
		deskripsiReverse = fmt.Sprintf("%s - %s", deskripsiReverse, keterangan)
	}

	// Buat transaksi pembalik dengan debit/kredit yang dibalik
	var reversalTransaksi models.Transaksi

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Buat header transaksi pembalik
		reversalTransaksi = models.Transaksi{
			IDKoperasi:       idKoperasi,
			NomorJurnal:      nomorJurnal,
			TanggalTransaksi: tanggalReverse,
			Deskripsi:        deskripsiReverse,
			NomorReferensi:   fmt.Sprintf("REV-%s", originalTransaksi.NomorJurnal), // Reference to original
			TipeTransaksi:    "reversal",
			TotalDebit:       originalTransaksi.TotalKredit, // Swap totals
			TotalKredit:      originalTransaksi.TotalDebit,  // Swap totals
			StatusBalanced:   true,
			DibuatOleh:       idPengguna,
		}

		if err := tx.Create(&reversalTransaksi).Error; err != nil {
			s.logger.Error(method, "Gagal membuat header transaksi pembalik", err, map[string]interface{}{
				"koperasi_id":        idKoperasi.String(),
				"nomor_jurnal":       nomorJurnal,
				"transaksi_original": originalTransaksi.NomorJurnal,
			})
			return utils.WrapDatabaseError(err, "Gagal membuat transaksi pembalik")
		}

		// Buat baris transaksi dengan debit/kredit yang dibalik
		for _, barisOriginal := range originalTransaksi.BarisTransaksi {
			barisReversal := models.BarisTransaksi{
				IDTransaksi:  reversalTransaksi.ID,
				IDAkun:       barisOriginal.IDAkun,
				JumlahDebit:  barisOriginal.JumlahKredit, // Swap: kredit jadi debit
				JumlahKredit: barisOriginal.JumlahDebit,  // Swap: debit jadi kredit
				Keterangan:   fmt.Sprintf("Reversal: %s", barisOriginal.Keterangan),
			}

			if err := tx.Create(&barisReversal).Error; err != nil {
				s.logger.Error(method, "Gagal membuat baris transaksi pembalik", err, map[string]interface{}{
					"transaksi_reversal_id": reversalTransaksi.ID.String(),
					"akun_id":               barisOriginal.IDAkun.String(),
				})
				return utils.WrapDatabaseError(err, "Gagal membuat baris transaksi pembalik")
			}
		}

		// Load transaksi pembalik dengan semua relasi
		if err := tx.Preload("BarisTransaksi.Akun").Where("id = ?", reversalTransaksi.ID).First(&reversalTransaksi).Error; err != nil {
			s.logger.Error(method, "Gagal mengambil transaksi pembalik", err, map[string]interface{}{
				"transaksi_reversal_id": reversalTransaksi.ID.String(),
			})
			return utils.WrapDatabaseError(err, "Gagal mengambil transaksi pembalik")
		}

		return nil
	})

	if err != nil {
		s.logger.Error(method, "Transaksi database gagal", err, map[string]interface{}{
			"koperasi_id":  idKoperasi.String(),
			"transaksi_id": id.String(),
			"nomor_jurnal": nomorJurnal,
		})
		return nil, err
	}

	s.logger.Info(method, "Berhasil membuat transaksi reversal", map[string]interface{}{
		"transaksi_original_id": id.String(),
		"transaksi_reversal_id": reversalTransaksi.ID.String(),
		"nomor_jurnal_original": originalTransaksi.NomorJurnal,
		"nomor_jurnal_reversal": nomorJurnal,
		"koperasi_id":           idKoperasi.String(),
	})

	respons := reversalTransaksi.ToResponse()
	return &respons, nil
}

// DapatkanBukuBesar mengambil buku besar (ledger) untuk akun tertentu
func (s *TransaksiService) DapatkanBukuBesar(idAkun uuid.UUID, tanggalMulai, tanggalAkhir string) ([]map[string]interface{}, error) {
	const method = "DapatkanBukuBesar"

	// Validasi akun exists
	var akun models.Akun
	err := s.db.Where("id = ?", idAkun).First(&akun).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Akun tidak ditemukan", err, map[string]interface{}{
				"akun_id": idAkun.String(),
			})
			return nil, utils.WrapDatabaseError(err, "Akun")
		}
		s.logger.Error(method, "Gagal mengambil data akun", err, map[string]interface{}{
			"akun_id": idAkun.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data akun")
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
		s.logger.Error(method, "Gagal mengambil data buku besar", err, map[string]interface{}{
			"akun_id":       idAkun.String(),
			"tanggal_mulai": tanggalMulai,
			"tanggal_akhir": tanggalAkhir,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data buku besar")
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

	s.logger.Debug(method, "Berhasil mengambil buku besar", map[string]interface{}{
		"akun_id":      idAkun.String(),
		"kode_akun":    akun.KodeAkun,
		"nama_akun":    akun.NamaAkun,
		"jumlah_baris": len(ledger),
		"saldo_akhir":  saldo,
	})

	return ledger, nil
}

// PostingOtomatisSimpanan membuat jurnal otomatis untuk setoran simpanan
func (s *TransaksiService) PostingOtomatisSimpanan(idKoperasi, idPengguna, idSimpanan uuid.UUID) error {
	const method = "PostingOtomatisSimpanan"

	// Ambil data simpanan
	var simpanan models.Simpanan
	err := s.db.Where("id = ?", idSimpanan).First(&simpanan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Simpanan tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"simpanan_id": idSimpanan.String(),
			})
			return utils.WrapDatabaseError(err, "Simpanan")
		}
		s.logger.Error(method, "Gagal mengambil data simpanan", err, map[string]interface{}{
			"simpanan_id": idSimpanan.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil data simpanan")
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Akun kas tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"kode_akun":   "1101",
			})
			return utils.WrapDatabaseError(err, "Akun kas")
		}
		s.logger.Error(method, "Gagal mengambil akun kas", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil akun kas")
	}

	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, kodeAkunModal).First(&akunModal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Akun modal tidak ditemukan", err, map[string]interface{}{
				"koperasi_id":     idKoperasi.String(),
				"kode_akun_modal": kodeAkunModal,
				"tipe_simpanan":   simpanan.TipeSimpanan,
			})
			return utils.WrapDatabaseError(err, "Akun modal")
		}
		s.logger.Error(method, "Gagal mengambil akun modal", err, map[string]interface{}{
			"koperasi_id":     idKoperasi.String(),
			"kode_akun_modal": kodeAkunModal,
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil akun modal")
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
		s.logger.Error(method, "Gagal membuat transaksi posting simpanan", err, map[string]interface{}{
			"koperasi_id":    idKoperasi.String(),
			"simpanan_id":    idSimpanan.String(),
			"tipe_simpanan":  simpanan.TipeSimpanan,
			"jumlah_setoran": simpanan.JumlahSetoran,
		})
		return fmt.Errorf("gagal posting simpanan: %w", err)
	}

	// Update simpanan dengan ID transaksi
	simpanan.IDTransaksi = &transaksi.ID
	if err := s.db.Save(&simpanan).Error; err != nil {
		s.logger.Error(method, "Gagal update simpanan dengan ID transaksi", err, map[string]interface{}{
			"simpanan_id":  idSimpanan.String(),
			"transaksi_id": transaksi.ID.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal update simpanan dengan ID transaksi")
	}

	s.logger.Info(method, "Berhasil posting otomatis simpanan", map[string]interface{}{
		"simpanan_id":    idSimpanan.String(),
		"transaksi_id":   transaksi.ID.String(),
		"tipe_simpanan":  simpanan.TipeSimpanan,
		"jumlah_setoran": simpanan.JumlahSetoran,
		"koperasi_id":    idKoperasi.String(),
	})

	return nil
}

// PostingOtomatisPenjualan membuat jurnal otomatis untuk penjualan
func (s *TransaksiService) PostingOtomatisPenjualan(idKoperasi, idPengguna, idPenjualan uuid.UUID) error {
	const method = "PostingOtomatisPenjualan"

	// Ambil data penjualan dengan items
	var penjualan models.Penjualan
	err := s.db.Preload("ItemPenjualan.Produk").Where("id = ?", idPenjualan).First(&penjualan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Penjualan tidak ditemukan", err, map[string]interface{}{
				"koperasi_id":  idKoperasi.String(),
				"penjualan_id": idPenjualan.String(),
			})
			return utils.WrapDatabaseError(err, "Penjualan")
		}
		s.logger.Error(method, "Gagal mengambil data penjualan", err, map[string]interface{}{
			"penjualan_id": idPenjualan.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil data penjualan")
	}

	// Dapatkan akun-akun yang diperlukan
	var akunKas, akunPenjualan, akunHPP, akunPersediaan models.Akun
	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1101").First(&akunKas).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Akun kas tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"kode_akun":   "1101",
			})
			return utils.WrapDatabaseError(err, "Akun kas")
		}
		s.logger.Error(method, "Gagal mengambil akun kas", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil akun kas")
	}

	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "4101").First(&akunPenjualan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Akun penjualan tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"kode_akun":   "4101",
			})
			return utils.WrapDatabaseError(err, "Akun penjualan")
		}
		s.logger.Error(method, "Gagal mengambil akun penjualan", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil akun penjualan")
	}

	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "5201").First(&akunHPP).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Akun HPP tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"kode_akun":   "5201",
			})
			return utils.WrapDatabaseError(err, "Akun HPP")
		}
		s.logger.Error(method, "Gagal mengambil akun HPP", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil akun HPP")
	}

	if err := s.db.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1301").First(&akunPersediaan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Akun persediaan tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"kode_akun":   "1301",
			})
			return utils.WrapDatabaseError(err, "Akun persediaan")
		}
		s.logger.Error(method, "Gagal mengambil akun persediaan", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil akun persediaan")
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
		s.logger.Error(method, "Gagal membuat transaksi posting penjualan", err, map[string]interface{}{
			"koperasi_id":     idKoperasi.String(),
			"penjualan_id":    idPenjualan.String(),
			"nomor_penjualan": penjualan.NomorPenjualan,
			"total_belanja":   penjualan.TotalBelanja,
			"total_hpp":       totalHPP,
		})
		return fmt.Errorf("gagal posting penjualan: %w", err)
	}

	// Update penjualan dengan ID transaksi
	penjualan.IDTransaksi = &transaksi.ID
	if err := s.db.Save(&penjualan).Error; err != nil {
		s.logger.Error(method, "Gagal update penjualan dengan ID transaksi", err, map[string]interface{}{
			"penjualan_id": idPenjualan.String(),
			"transaksi_id": transaksi.ID.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal update penjualan dengan ID transaksi")
	}

	s.logger.Info(method, "Berhasil posting otomatis penjualan", map[string]interface{}{
		"penjualan_id":    idPenjualan.String(),
		"transaksi_id":    transaksi.ID.String(),
		"nomor_penjualan": penjualan.NomorPenjualan,
		"total_belanja":   penjualan.TotalBelanja,
		"total_hpp":       totalHPP,
		"jumlah_item":     len(penjualan.ItemPenjualan),
		"koperasi_id":     idKoperasi.String(),
	})

	return nil
}
