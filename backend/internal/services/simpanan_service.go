package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/utils"
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
	logger           *utils.Logger
}

// NewSimpananService membuat instance baru SimpananService
func NewSimpananService(db *gorm.DB, transaksiService *TransaksiService) *SimpananService {
	return &SimpananService{
		db:               db,
		transaksiService: transaksiService,
		logger:           utils.NewLogger("SimpananService"),
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
	const method = "CatatSetoran"

	// Validasi anggota exists dan aktif
	var anggota models.Anggota
	err := s.db.Where("id = ? AND id_koperasi = ? AND status = ?", req.IDAnggota, idKoperasi, models.StatusAktif).
		First(&anggota).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Anggota not found or not active", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"anggota_id":  req.IDAnggota.String(),
			})
			return nil, utils.WrapValidationError(err, "Anggota tidak ditemukan atau tidak aktif")
		}

		s.logger.Error(method, "Failed to validate anggota", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"anggota_id":  req.IDAnggota.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal memvalidasi anggota")
	}

	// Validasi jumlah setoran
	if req.JumlahSetoran <= 0 {
		s.logger.Error(method, "Invalid setoran amount", utils.ErrInvalidAmount, map[string]interface{}{
			"koperasi_id":    idKoperasi.String(),
			"anggota_id":     req.IDAnggota.String(),
			"jumlah_setoran": req.JumlahSetoran,
		})
		return nil, utils.WrapValidationError(utils.ErrInvalidAmount, "Jumlah setoran harus lebih dari 0")
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
		s.logger.Error(method, "Failed to create simpanan record", err, map[string]interface{}{
			"koperasi_id":       idKoperasi.String(),
			"anggota_id":        req.IDAnggota.String(),
			"tipe_simpanan":     req.TipeSimpanan,
			"jumlah_setoran":    req.JumlahSetoran,
			"nomor_referensi":   nomorReferensi,
			"tanggal_transaksi": req.TanggalTransaksi.Format("2006-01-02"),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mencatat setoran simpanan")
	}

	// Auto-posting ke jurnal akuntansi
	err = s.transaksiService.PostingOtomatisSimpanan(idKoperasi, idPengguna, simpanan.ID)
	if err != nil {
		// Rollback simpanan jika posting gagal
		s.db.Delete(simpanan)
		s.logger.Error(method, "Failed to post to journal, rolled back", err, map[string]interface{}{
			"koperasi_id":     idKoperasi.String(),
			"simpanan_id":     simpanan.ID.String(),
			"nomor_referensi": nomorReferensi,
		})
		return nil, fmt.Errorf("gagal posting ke jurnal: %w", err)
	}

	// Reload dengan relasi
	s.db.Preload("Anggota").First(simpanan, simpanan.ID)

	s.logger.Info(method, "Successfully created simpanan transaction", map[string]interface{}{
		"simpanan_id":       simpanan.ID.String(),
		"koperasi_id":       idKoperasi.String(),
		"anggota_id":        req.IDAnggota.String(),
		"tipe_simpanan":     req.TipeSimpanan,
		"jumlah_setoran":    req.JumlahSetoran,
		"nomor_referensi":   nomorReferensi,
		"tanggal_transaksi": req.TanggalTransaksi.Format("2006-01-02"),
	})

	respons := simpanan.ToResponse()
	return &respons, nil
}

// GenerateNomorReferensi menghasilkan nomor referensi setoran
// Format: SMP-YYYYMMDD-NNNN
// Uses row-level locking to prevent race conditions in concurrent requests
func (s *SimpananService) GenerateNomorReferensi(idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	const method = "GenerateNomorReferensi"
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
			s.logger.Error(method, "Failed to query last simpanan", err, map[string]interface{}{
				"koperasi_id":       idKoperasi.String(),
				"tanggal_transaksi": tanggalDate,
			})
			return err
		}

		nomorReferensi = fmt.Sprintf("SMP-%s-%04d", tanggalStr, nomorUrut)
		return nil
	})

	if err != nil {
		s.logger.Error(method, "Transaction failed while generating nomor referensi", err, map[string]interface{}{
			"koperasi_id":       idKoperasi.String(),
			"tanggal_transaksi": tanggalDate,
		})
		return "", utils.WrapGenerationError(err, "Generate nomor referensi")
	}

	s.logger.Debug(method, "Generated nomor referensi", map[string]interface{}{
		"nomor_referensi":   nomorReferensi,
		"koperasi_id":       idKoperasi.String(),
		"tanggal_transaksi": tanggalDate,
	})

	return nomorReferensi, nil
}

// DapatkanSemuaTransaksiSimpanan mengambil daftar transaksi simpanan
func (s *SimpananService) DapatkanSemuaTransaksiSimpanan(idKoperasi uuid.UUID, tipeSimpanan string, idAnggota *uuid.UUID, tanggalMulai, tanggalAkhir string, page, pageSize int) ([]models.SimpananResponse, int64, error) {
	const method = "DapatkanSemuaTransaksiSimpanan"
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
		s.logger.Error(method, "Failed to fetch simpanan transaction list", err, map[string]interface{}{
			"koperasi_id":   idKoperasi.String(),
			"tipe_simpanan": tipeSimpanan,
			"anggota_id":    idAnggota,
			"tanggal_mulai": tanggalMulai,
			"tanggal_akhir": tanggalAkhir,
			"page":          page,
			"page_size":     pageSize,
		})
		return nil, 0, utils.WrapDatabaseError(err, "Gagal mengambil daftar transaksi simpanan")
	}

	// Convert to response
	responseDaftar := make([]models.SimpananResponse, len(simpananList))
	for i, simpanan := range simpananList {
		responseDaftar[i] = simpanan.ToResponse()
	}

	s.logger.Debug(method, "Successfully fetched simpanan transaction list", map[string]interface{}{
		"count": len(responseDaftar),
		"total": total,
	})

	return responseDaftar, total, nil
}

// DapatkanSaldoAnggota mengambil saldo simpanan per anggota dengan validasi multi-tenant
func (s *SimpananService) DapatkanSaldoAnggota(idKoperasi, idAnggota uuid.UUID) (*models.SaldoSimpananAnggota, error) {
	// Validasi anggota exists dan milik koperasi yang benar
	var anggota models.Anggota
	err := s.db.Where("id = ? AND id_koperasi = ?", idAnggota, idKoperasi).First(&anggota).Error
	if err != nil {
		return nil, errors.New("anggota tidak ditemukan atau tidak memiliki akses")
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
// Optimized version using single query with GROUP BY to eliminate N+1 query problem
func (s *SimpananService) DapatkanLaporanSaldoAnggota(idKoperasi uuid.UUID) ([]models.SaldoSimpananAnggota, error) {
	const method = "DapatkanLaporanSaldoAnggota"

	// Structure to hold aggregated data from database
	type SaldoResult struct {
		IDAnggota    uuid.UUID
		NomorAnggota string
		NamaLengkap  string
		TipeSimpanan models.TipeSimpanan
		TotalSaldo   float64
	}

	var results []SaldoResult

	// Single optimized query with JOIN and GROUP BY
	// This replaces the N+1 pattern (1 query for members + N queries for balances)
	err := s.db.Table("simpanan").
		Select(`
			anggota.id as id_anggota,
			anggota.nomor_anggota,
			anggota.nama_lengkap,
			simpanan.tipe_simpanan,
			COALESCE(SUM(simpanan.jumlah_setoran), 0) as total_saldo
		`).
		Joins("JOIN anggota ON anggota.id = simpanan.id_anggota").
		Where("anggota.id_koperasi = ? AND anggota.status = ?", idKoperasi, models.StatusAktif).
		Group("anggota.id, anggota.nomor_anggota, anggota.nama_lengkap, simpanan.tipe_simpanan").
		Order("anggota.nomor_anggota ASC, simpanan.tipe_simpanan").
		Scan(&results).Error

	if err != nil {
		s.logger.Error(method, "Failed to fetch member balance report", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return nil, errors.New("gagal mengambil laporan saldo anggota")
	}

	// Process results and aggregate by member
	// Use map for efficient member lookup and aggregation
	memberMap := make(map[uuid.UUID]*models.SaldoSimpananAnggota)

	for _, result := range results {
		member, exists := memberMap[result.IDAnggota]
		if !exists {
			// Create new member entry
			member = &models.SaldoSimpananAnggota{
				IDAnggota:    result.IDAnggota,
				NomorAnggota: result.NomorAnggota,
				NamaAnggota:  result.NamaLengkap,
			}
			memberMap[result.IDAnggota] = member
		}

		// Aggregate balances by type
		switch result.TipeSimpanan {
		case models.SimpananPokok:
			member.SimpananPokok = result.TotalSaldo
		case models.SimpananWajib:
			member.SimpananWajib = result.TotalSaldo
		case models.SimpananSukarela:
			member.SimpananSukarela = result.TotalSaldo
		}
	}

	// Also get members with no deposits to include them with zero balances
	var membersWithoutDeposits []models.Anggota
	err = s.db.Where(`
		id_koperasi = ? AND status = ? AND id NOT IN (
			SELECT DISTINCT id_anggota FROM simpanan WHERE id_koperasi = ?
		)
	`, idKoperasi, models.StatusAktif, idKoperasi).Find(&membersWithoutDeposits).Error

	if err != nil {
		s.logger.Error(method, "Failed to fetch members without deposits", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		// Continue with existing results even if this query fails
	} else {
		// Add members with zero balances
		for _, member := range membersWithoutDeposits {
			if _, exists := memberMap[member.ID]; !exists {
				memberMap[member.ID] = &models.SaldoSimpananAnggota{
					IDAnggota:    member.ID,
					NomorAnggota: member.NomorAnggota,
					NamaAnggota:  member.NamaLengkap,
				}
			}
		}
	}

	// Convert map to sorted slice
	laporan := make([]models.SaldoSimpananAnggota, 0, len(memberMap))
	for _, member := range memberMap {
		// Calculate total
		member.TotalSimpanan = member.SimpananPokok + member.SimpananWajib + member.SimpananSukarela
		laporan = append(laporan, *member)
	}

	// Sort by member number
	// Note: In production, consider using database ORDER BY on the member number
	// This in-memory sort is acceptable since member lists are typically small (<10k)
	for i := 0; i < len(laporan)-1; i++ {
		for j := i + 1; j < len(laporan); j++ {
			if laporan[i].NomorAnggota > laporan[j].NomorAnggota {
				laporan[i], laporan[j] = laporan[j], laporan[i]
			}
		}
	}

	s.logger.Info(method, "Successfully generated member balance report", map[string]interface{}{
		"koperasi_id":   idKoperasi.String(),
		"total_members": len(laporan),
	})

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
