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

// PenjualanService menangani logika bisnis penjualan (POS)
type PenjualanService struct {
	db               *gorm.DB
	produkService    *ProdukService
	transaksiService *TransaksiService
}

// NewPenjualanService membuat instance baru PenjualanService
func NewPenjualanService(db *gorm.DB, produkService *ProdukService, transaksiService *TransaksiService) *PenjualanService {
	return &PenjualanService{
		db:               db,
		produkService:    produkService,
		transaksiService: transaksiService,
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
	// Validasi items (stok tersedia)
	if err := s.ValidasiItemPenjualan(req.Items); err != nil {
		return nil, err
	}

	// Hitung total belanja
	var totalBelanja float64
	for _, item := range req.Items {
		totalBelanja += item.HargaSatuan * float64(item.Kuantitas)
	}

	// Validasi pembayaran
	if err := s.ValidasiPembayaran(totalBelanja, req.JumlahBayar); err != nil {
		return nil, err
	}

	// Generate nomor penjualan
	nomorPenjualan, err := s.GenerateNomorPenjualan(idKoperasi, time.Now())
	if err != nil {
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
			return errors.New("gagal membuat penjualan")
		}

		// 2. Buat item penjualan dan kurangi stok
		for _, itemReq := range req.Items {
			// Dapatkan produk untuk nama
			var produk models.Produk
			if err := tx.Where("id = ?", itemReq.IDProduk).First(&produk).Error; err != nil {
				return fmt.Errorf("produk %s tidak ditemukan", itemReq.IDProduk)
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
				return errors.New("gagal membuat item penjualan")
			}

			// Kurangi stok produk
			if err := s.produkService.KurangiStok(itemReq.IDProduk, itemReq.Kuantitas); err != nil {
				return fmt.Errorf("gagal mengurangi stok: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 3. Auto-posting ke jurnal akuntansi
	err = s.transaksiService.PostingOtomatisPenjualan(idKoperasi, idKasir, penjualan.ID)
	if err != nil {
		// Warning: penjualan sudah tersimpan tapi posting gagal
		// Bisa di-handle dengan background job untuk retry
		return nil, fmt.Errorf("penjualan berhasil, tetapi posting gagal: %w", err)
	}

	// Reload dengan relasi
	s.db.Preload("ItemPenjualan.Produk").Preload("Kasir").Preload("Anggota").First(&penjualan, penjualan.ID)

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
// Uses row-level locking to prevent race conditions in concurrent requests
func (s *PenjualanService) GenerateNomorPenjualan(idKoperasi uuid.UUID, tanggal time.Time) (string, error) {
	tanggalStr := tanggal.Format("20060102")
	tanggalDate := tanggal.Format("2006-01-02")
	var nomorPenjualan string

	// Use transaction with row-level locking to prevent race conditions
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Lock and get the last sales number for this date
		var lastPenjualan models.Penjualan
		err := tx.Where("id_koperasi = ? AND DATE(tanggal_penjualan) = ?", idKoperasi, tanggalDate).
			Order("nomor_penjualan DESC").
			Limit(1).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastPenjualan).Error

		nomorUrut := 1

		// If there's a previous sale, parse and increment
		if err == nil && lastPenjualan.NomorPenjualan != "" {
			// Extract number from POS-20250116-0001
			var parsedTanggal string
			var parsedUrut int
			_, scanErr := fmt.Sscanf(lastPenjualan.NomorPenjualan, "POS-%s-%04d", &parsedTanggal, &parsedUrut)
			if scanErr == nil && parsedTanggal == tanggalStr {
				nomorUrut = parsedUrut + 1
			}
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		nomorPenjualan = fmt.Sprintf("POS-%s-%04d", tanggalStr, nomorUrut)
		return nil
	})

	if err != nil {
		return "", errors.New("gagal generate nomor penjualan")
	}

	return nomorPenjualan, nil
}

// DapatkanSemuaPenjualan mengambil daftar penjualan dengan filter
func (s *PenjualanService) DapatkanSemuaPenjualan(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string, idKasir *uuid.UUID, page, pageSize int) ([]models.PenjualanResponse, int64, error) {
	var penjualanList []models.Penjualan
	var total int64

	query := s.db.Model(&models.Penjualan{}).Where("id_koperasi = ?", idKoperasi)

	// Apply filters
	if tanggalMulai != "" {
		query = query.Where("tanggal_penjualan >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("tanggal_penjualan <= ?", tanggalAkhir)
	}
	if idKasir != nil {
		query = query.Where("id_kasir = ?", *idKasir)
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("tanggal_penjualan DESC").
		Preload("ItemPenjualan.Produk").
		Preload("Kasir").
		Preload("Anggota").
		Find(&penjualanList).Error

	if err != nil {
		return nil, 0, errors.New("gagal mengambil daftar penjualan")
	}

	// Convert to response
	responses := make([]models.PenjualanResponse, len(penjualanList))
	for i, penjualan := range penjualanList {
		responses[i] = penjualan.ToResponse()
	}

	return responses, total, nil
}

// DapatkanPenjualan mengambil penjualan berdasarkan ID
func (s *PenjualanService) DapatkanPenjualan(id uuid.UUID) (*models.PenjualanResponse, error) {
	var penjualan models.Penjualan
	err := s.db.Preload("ItemPenjualan.Produk").
		Preload("Kasir").
		Preload("Anggota").
		Where("id = ?", id).
		First(&penjualan).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("penjualan tidak ditemukan")
		}
		return nil, err
	}

	response := penjualan.ToResponse()
	return &response, nil
}

// DapatkanStruk mengambil data struk digital
func (s *PenjualanService) DapatkanStruk(id uuid.UUID) (*models.PenjualanResponse, error) {
	return s.DapatkanPenjualan(id)
}

// HitungTotalPenjualan menghitung total penjualan dalam periode
func (s *PenjualanService) HitungTotalPenjualan(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string) (map[string]interface{}, error) {
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
		return nil, errors.New("gagal menghitung total penjualan")
	}

	summary := map[string]interface{}{
		"totalPenjualan":  result.TotalPenjualan,
		"jumlahTransaksi": result.JumlahTransaksi,
		"rataRata":        float64(0),
	}

	if result.JumlahTransaksi > 0 {
		summary["rataRata"] = result.TotalPenjualan / float64(result.JumlahTransaksi)
	}

	return summary, nil
}

// DapatkanPenjualanHariIni mengambil penjualan hari ini
func (s *PenjualanService) DapatkanPenjualanHariIni(idKoperasi uuid.UUID) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")
	return s.HitungTotalPenjualan(idKoperasi, today, today)
}

// DapatkanTopProduk mengambil produk terlaris
func (s *PenjualanService) DapatkanTopProduk(idKoperasi uuid.UUID, limit int) ([]map[string]interface{}, error) {
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
		return nil, errors.New("gagal mengambil top produk")
	}

	// Convert to map
	topProduk := make([]map[string]interface{}, len(results))
	for i, result := range results {
		topProduk[i] = map[string]interface{}{
			"idProduk":     result.IDProduk,
			"namaProduk":   result.NamaProduk,
			"totalTerjual": result.TotalTerjual,
			"totalNilai":   result.TotalNilai,
		}
	}

	return topProduk, nil
}
