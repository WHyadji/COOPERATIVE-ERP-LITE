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

// PenjualanService menangani logika bisnis penjualan (POS)
type PenjualanService struct {
	db               *gorm.DB
	produkService    *ProdukService
	transaksiService *TransaksiService
	logger           *utils.Logger
}

// NewPenjualanService membuat instance baru PenjualanService
func NewPenjualanService(db *gorm.DB, produkService *ProdukService, transaksiService *TransaksiService) *PenjualanService {
	return &PenjualanService{
		db:               db,
		produkService:    produkService,
		transaksiService: transaksiService,
		logger:           utils.NewLogger("PenjualanService"),
	}
}

// ItemPenjualanRequest adalah struktur untuk item dalam penjualan
type ItemPenjualanRequest struct {
	IDProduk    uuid.UUID `json:"idProduk" binding:"required"`
	Kuantitas   int       `json:"kuantitas" binding:"required,gt=0"`
	HargaSatuan float64   `json:"hargaSatuan" binding:"required,gt=0"`
}

// ProsesPenjualanRequest adalah struktur request untuk proses penjualan
type ProsesPenjualanRequest struct {
	IDAnggota   *uuid.UUID             `json:"idAnggota"` // Optional
	Items       []ItemPenjualanRequest `json:"items" binding:"required,min=1"`
	JumlahBayar float64                `json:"jumlahBayar" binding:"required,gt=0"`
	Catatan     string                 `json:"catatan"`
}

// ProsesPenjualan memproses transaksi penjualan lengkap
func (s *PenjualanService) ProsesPenjualan(idKoperasi, idKasir uuid.UUID, req *ProsesPenjualanRequest) (*models.PenjualanResponse, error) {
	const method = "ProsesPenjualan"

	// Validasi items (stok tersedia)
	if err := s.ValidasiItemPenjualan(req.Items); err != nil {
		s.logger.Error(method, "Validasi item penjualan gagal", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"kasir_id":    idKasir.String(),
			"jumlah_item": len(req.Items),
		})
		return nil, err
	}

	// Hitung total belanja
	var totalBelanja float64
	for _, item := range req.Items {
		totalBelanja += item.HargaSatuan * float64(item.Kuantitas)
	}

	// Validasi pembayaran
	if err := s.ValidasiPembayaran(totalBelanja, req.JumlahBayar); err != nil {
		s.logger.Error(method, "Validasi pembayaran gagal", err, map[string]interface{}{
			"koperasi_id":   idKoperasi.String(),
			"total_belanja": totalBelanja,
			"jumlah_bayar":  req.JumlahBayar,
		})
		return nil, err
	}

	// Generate nomor penjualan
	nomorPenjualan, err := s.GenerateNomorPenjualan(idKoperasi, time.Now())
	if err != nil {
		s.logger.Error(method, "Gagal generate nomor penjualan", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return nil, err
	}

	// Hitung kembalian
	kembalian := req.JumlahBayar - totalBelanja

	// Proses dalam transaction
	var penjualan models.Penjualan

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buat record penjualan
		penjualan = models.Penjualan{
			IDKoperasi:       idKoperasi,
			NomorPenjualan:   nomorPenjualan,
			TanggalPenjualan: time.Now(),
			IDAnggota:        req.IDAnggota,
			TotalBelanja:     totalBelanja,
			MetodePembayaran: models.PembayaranTunai,
			JumlahBayar:      req.JumlahBayar,
			Kembalian:        kembalian,
			IDKasir:          idKasir,
			Catatan:          req.Catatan,
		}

		if err := tx.Create(&penjualan).Error; err != nil {
			s.logger.Error(method, "Gagal membuat record penjualan di database", err, map[string]interface{}{
				"koperasi_id":     idKoperasi.String(),
				"nomor_penjualan": nomorPenjualan,
				"total_belanja":   totalBelanja,
			})
			return utils.WrapDatabaseError(err, "Gagal membuat penjualan")
		}

		// 2. Buat item penjualan dan kurangi stok
		for _, itemReq := range req.Items {
			// Dapatkan produk untuk nama
			var produk models.Produk
			if err := tx.Where("id = ?", itemReq.IDProduk).First(&produk).Error; err != nil {
				s.logger.Error(method, "Produk tidak ditemukan saat proses penjualan", err, map[string]interface{}{
					"koperasi_id": idKoperasi.String(),
					"produk_id":   itemReq.IDProduk.String(),
				})
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.WrapDatabaseError(err, "Produk")
				}
				return utils.WrapDatabaseError(err, "Gagal mengambil data produk")
			}

			// Buat item penjualan
			item := models.ItemPenjualan{
				IDPenjualan: penjualan.ID,
				IDProduk:    itemReq.IDProduk,
				NamaProduk:  produk.NamaProduk, // Snapshot nama
				Kuantitas:   itemReq.Kuantitas,
				HargaSatuan: itemReq.HargaSatuan,
			}

			if err := tx.Create(&item).Error; err != nil {
				s.logger.Error(method, "Gagal membuat item penjualan", err, map[string]interface{}{
					"penjualan_id": penjualan.ID.String(),
					"produk_id":    itemReq.IDProduk.String(),
					"nama_produk":  produk.NamaProduk,
				})
				return utils.WrapDatabaseError(err, "Gagal membuat item penjualan")
			}

			// Kurangi stok produk
			if err := s.produkService.KurangiStok(itemReq.IDProduk, itemReq.Kuantitas); err != nil {
				s.logger.Error(method, "Gagal mengurangi stok produk", err, map[string]interface{}{
					"produk_id":   itemReq.IDProduk.String(),
					"nama_produk": produk.NamaProduk,
					"kuantitas":   itemReq.Kuantitas,
				})
				return fmt.Errorf("gagal mengurangi stok: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		s.logger.Error(method, "Transaksi penjualan gagal", err, map[string]interface{}{
			"koperasi_id":     idKoperasi.String(),
			"nomor_penjualan": nomorPenjualan,
		})
		return nil, err
	}

	// 3. Auto-posting ke jurnal akuntansi
	err = s.transaksiService.PostingOtomatisPenjualan(idKoperasi, idKasir, penjualan.ID)
	if err != nil {
		// Warning: penjualan sudah tersimpan tapi posting gagal
		s.logger.Error(method, "Penjualan berhasil tetapi posting ke jurnal gagal", err, map[string]interface{}{
			"penjualan_id":    penjualan.ID.String(),
			"nomor_penjualan": nomorPenjualan,
			"koperasi_id":     idKoperasi.String(),
		})
		return nil, fmt.Errorf("penjualan berhasil, tetapi posting gagal: %w", err)
	}

	// Reload dengan relasi
	if err := s.db.Preload("ItemPenjualan.Produk").Preload("Kasir").Preload("Anggota").First(&penjualan, penjualan.ID).Error; err != nil {
		s.logger.Error(method, "Gagal reload data penjualan setelah berhasil", err, map[string]interface{}{
			"penjualan_id": penjualan.ID.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data penjualan")
	}

	s.logger.Info(method, "Berhasil memproses penjualan", map[string]interface{}{
		"penjualan_id":    penjualan.ID.String(),
		"nomor_penjualan": nomorPenjualan,
		"total_belanja":   totalBelanja,
		"jumlah_item":     len(req.Items),
		"koperasi_id":     idKoperasi.String(),
		"kasir_id":        idKasir.String(),
	})

	response := penjualan.ToResponse()
	return &response, nil
}

// ValidasiItemPenjualan memvalidasi semua item (stok tersedia)
func (s *PenjualanService) ValidasiItemPenjualan(items []ItemPenjualanRequest) error {
	for _, item := range items {
		tersedia, err := s.produkService.CekStokTersedia(item.IDProduk, item.Kuantitas)
		if err != nil {
			return err
		}

		if !tersedia {
			return fmt.Errorf("stok produk tidak mencukupi untuk item %s", item.IDProduk)
		}
	}

	return nil
}

// ValidasiPembayaran memvalidasi jumlah bayar cukup
func (s *PenjualanService) ValidasiPembayaran(totalBelanja, jumlahBayar float64) error {
	if jumlahBayar < totalBelanja {
		return fmt.Errorf("jumlah bayar (%.2f) kurang dari total belanja (%.2f)", jumlahBayar, totalBelanja)
	}
	return nil
}

// GenerateNomorPenjualan menghasilkan nomor penjualan otomatis
// Format: POS-YYYYMMDD-NNNN
// Menggunakan row-level locking untuk mencegah race condition pada concurrent requests
func (s *PenjualanService) GenerateNomorPenjualan(idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	const method = "GenerateNomorPenjualan"

	tanggalStr := tanggal.Format("20060102")
	tanggalDate := tanggal.Format("2006-01-02")
	var nomorPenjualan string

	// Gunakan transaction dengan row-level locking untuk mencegah race condition
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Lock dan ambil nomor penjualan terakhir untuk tanggal ini
		var lastPenjualan models.Penjualan
		err := tx.Where("id_koperasi = ? AND DATE(tanggal_penjualan) = ?", idKoperasi, tanggalDate).
			Order("nomor_penjualan DESC").
			Limit(1).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastPenjualan).Error

		nomorUrut := 1

		// Jika ada penjualan sebelumnya, parse dan increment
		if err == nil && lastPenjualan.NomorPenjualan != "" {
			// Extract number dari POS-20250116-0001
			var parsedTanggal string
			var parsedUrut int
			_, scanErr := fmt.Sscanf(lastPenjualan.NomorPenjualan, "POS-%s-%04d", &parsedTanggal, &parsedUrut)
			if scanErr == nil && parsedTanggal == tanggalStr {
				nomorUrut = parsedUrut + 1
			}
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Gagal mengambil nomor penjualan terakhir", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"tanggal":     tanggalDate,
			})
			return utils.WrapDatabaseError(err, "Gagal mengambil nomor penjualan terakhir")
		}

		nomorPenjualan = fmt.Sprintf("POS-%s-%04d", tanggalStr, nomorUrut)
		return nil
	})

	if err != nil {
		s.logger.Error(method, "Transaksi generate nomor penjualan gagal", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"tanggal":     tanggalDate,
		})
		return "", err
	}

	s.logger.Debug(method, "Berhasil generate nomor penjualan", map[string]interface{}{
		"nomor_penjualan": nomorPenjualan,
		"koperasi_id":     idKoperasi.String(),
	})

	return nomorPenjualan, nil
}

// DapatkanSemuaPenjualan mengambil daftar penjualan dengan filter
func (s *PenjualanService) DapatkanSemuaPenjualan(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string, idKasir *uuid.UUID, page, pageSize int) ([]models.PenjualanResponse, int64, error) {
	const method = "DapatkanSemuaPenjualan"

	// Validate and normalize pagination parameters to prevent DoS attacks
	validPage, validPageSize := utils.ValidatePagination(page, pageSize)

	var penjualanList []models.Penjualan
	var total int64

	query := s.db.Model(&models.Penjualan{}).Where("id_koperasi = ?", idKoperasi)

	// Terapkan filter
	if tanggalMulai != "" {
		query = query.Where("tanggal_penjualan >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("tanggal_penjualan <= ?", tanggalAkhir)
	}
	if idKasir != nil {
		query = query.Where("id_kasir = ?", *idKasir)
	}

	// Hitung total
	query.Count(&total)

	// Pagination with validated parameters
	offset := utils.CalculateOffset(validPage, validPageSize)

	// Create context with timeout to prevent long-running queries
	ctx, cancel := utils.CreateQueryContext()
	defer cancel()

	err := query.WithContext(ctx).Offset(offset).Limit(validPageSize).
		Order("tanggal_penjualan DESC").
		Preload("ItemPenjualan.Produk").
		Preload("Kasir").
		Preload("Anggota").
		Find(&penjualanList).Error

	if err != nil {
		s.logger.Error(method, "Gagal mengambil daftar penjualan dari database", err, map[string]interface{}{
			"koperasi_id":   idKoperasi.String(),
			"tanggal_mulai": tanggalMulai,
			"tanggal_akhir": tanggalAkhir,
			"page":          validPage,
			"page_size":     validPageSize,
		})
		return nil, 0, utils.WrapDatabaseError(err, "Gagal mengambil daftar penjualan")
	}

	// Convert ke response
	responses := make([]models.PenjualanResponse, len(penjualanList))
	for i, penjualan := range penjualanList {
		responses[i] = penjualan.ToResponse()
	}

	s.logger.Debug(method, "Berhasil mengambil daftar penjualan", map[string]interface{}{
		"koperasi_id": idKoperasi.String(),
		"total":       total,
		"count":       len(responses),
		"page":        validPage,
	})

	return responses, total, nil
}

// DapatkanPenjualan mengambil penjualan berdasarkan ID
func (s *PenjualanService) DapatkanPenjualan(idKoperasi, id uuid.UUID) (*models.PenjualanResponse, error) {
	const method = "DapatkanPenjualan"

	var penjualan models.Penjualan
	err := s.db.Preload("ItemPenjualan.Produk").
		Preload("Kasir").
		Preload("Anggota").
		Where("id = ? AND id_koperasi = ?", id, idKoperasi).
		First(&penjualan).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Penjualan tidak ditemukan atau tidak memiliki akses", err, map[string]interface{}{
				"penjualan_id": id.String(),
				"koperasi_id":  idKoperasi.String(),
			})
			return nil, utils.WrapDatabaseError(err, "Penjualan")
		}
		s.logger.Error(method, "Gagal mengambil data penjualan", err, map[string]interface{}{
			"penjualan_id": id.String(),
			"koperasi_id":  idKoperasi.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data penjualan")
	}

	s.logger.Debug(method, "Berhasil mengambil data penjualan", map[string]interface{}{
		"penjualan_id":    id.String(),
		"koperasi_id":     idKoperasi.String(),
		"nomor_penjualan": penjualan.NomorPenjualan,
	})

	response := penjualan.ToResponse()
	return &response, nil
}

// DapatkanStruk mengambil data struk digital
func (s *PenjualanService) DapatkanStruk(idKoperasi, id uuid.UUID) (*models.PenjualanResponse, error) {
	return s.DapatkanPenjualan(idKoperasi, id)
}

// HitungTotalPenjualan menghitung total penjualan dalam periode
func (s *PenjualanService) HitungTotalPenjualan(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string) (map[string]interface{}, error) {
	const method = "HitungTotalPenjualan"

	type SalesResult struct {
		TotalPenjualan  float64
		JumlahTransaksi int64
	}

	var result SalesResult
	query := s.db.Model(&models.Penjualan{}).
		Select("COALESCE(SUM(total_belanja), 0) as total_penjualan, COUNT(*) as jumlah_transaksi").
		Where("id_koperasi = ?", idKoperasi)

	if tanggalMulai != "" {
		query = query.Where("tanggal_penjualan >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("tanggal_penjualan <= ?", tanggalAkhir)
	}

	err := query.Scan(&result).Error
	if err != nil {
		s.logger.Error(method, "Gagal menghitung total penjualan", err, map[string]interface{}{
			"koperasi_id":   idKoperasi.String(),
			"tanggal_mulai": tanggalMulai,
			"tanggal_akhir": tanggalAkhir,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal menghitung total penjualan")
	}

	summary := map[string]interface{}{
		"totalPenjualan":  result.TotalPenjualan,
		"jumlahTransaksi": result.JumlahTransaksi,
		"rataRata":        float64(0),
	}

	if result.JumlahTransaksi > 0 {
		summary["rataRata"] = result.TotalPenjualan / float64(result.JumlahTransaksi)
	}

	s.logger.Debug(method, "Berhasil menghitung total penjualan", map[string]interface{}{
		"koperasi_id":      idKoperasi.String(),
		"total_penjualan":  result.TotalPenjualan,
		"jumlah_transaksi": result.JumlahTransaksi,
	})

	return summary, nil
}

// DapatkanPenjualanHariIni mengambil penjualan hari ini
func (s *PenjualanService) DapatkanPenjualanHariIni(idKoperasi uuid.UUID) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")
	return s.HitungTotalPenjualan(idKoperasi, today, today)
}

// DapatkanTopProduk mengambil produk terlaris
func (s *PenjualanService) DapatkanTopProduk(idKoperasi uuid.UUID, limit int) ([]map[string]interface{}, error) {
	const method = "DapatkanTopProduk"

	type TopProduk struct {
		IDProduk     uuid.UUID
		NamaProduk   string
		TotalTerjual int
		TotalNilai   float64
	}

	var results []TopProduk
	err := s.db.Model(&models.ItemPenjualan{}).
		Select("item_penjualan.id_produk, item_penjualan.nama_produk, SUM(item_penjualan.kuantitas) as total_terjual, SUM(item_penjualan.subtotal) as total_nilai").
		Joins("JOIN penjualan ON penjualan.id = item_penjualan.id_penjualan").
		Where("penjualan.id_koperasi = ?", idKoperasi).
		Group("item_penjualan.id_produk, item_penjualan.nama_produk").
		Order("total_terjual DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		s.logger.Error(method, "Gagal mengambil data top produk", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"limit":       limit,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil top produk")
	}

	// Convert ke map
	topProduk := make([]map[string]interface{}, len(results))
	for i, result := range results {
		topProduk[i] = map[string]interface{}{
			"idProduk":     result.IDProduk,
			"namaProduk":   result.NamaProduk,
			"totalTerjual": result.TotalTerjual,
			"totalNilai":   result.TotalNilai,
		}
	}

	s.logger.Debug(method, "Berhasil mengambil top produk", map[string]interface{}{
		"koperasi_id": idKoperasi.String(),
		"count":       len(topProduk),
	})

	return topProduk, nil
}
