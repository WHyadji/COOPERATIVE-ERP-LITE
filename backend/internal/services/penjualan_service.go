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

	// Proses dalam transaction - includes auto-posting for data consistency
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

			// Kurangi stok produk within transaction
			if err := s.produkService.KurangiStokWithTx(tx, itemReq.IDProduk, itemReq.Kuantitas); err != nil {
				return fmt.Errorf("gagal mengurangi stok: %w", err)
			}
		}

		// 3. Auto-posting ke jurnal akuntansi within same transaction
		if err := s.postingPenjualanWithTx(tx, idKoperasi, idKasir, penjualan.ID); err != nil {
			return fmt.Errorf("gagal posting ke jurnal: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err // Automatic rollback on any error
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

// postingPenjualanWithTx creates journal entry for penjualan within an existing transaction
// This ensures atomicity - if posting fails, penjualan and stock changes are also rolled back
func (s *PenjualanService) postingPenjualanWithTx(tx *gorm.DB, idKoperasi, idPengguna, idPenjualan uuid.UUID) error {
	// Get penjualan data with items
	var penjualan models.Penjualan
	if err := tx.Preload("ItemPenjualan.Produk").Where("id = ?", idPenjualan).First(&penjualan).Error; err != nil {
		return errors.New("penjualan tidak ditemukan")
	}

	// Get required accounts
	var akunKas, akunPenjualan, akunHPP, akunPersediaan models.Akun
	if err := tx.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1101").First(&akunKas).Error; err != nil {
		return errors.New("akun kas tidak ditemukan")
	}
	if err := tx.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "4101").First(&akunPenjualan).Error; err != nil {
		return errors.New("akun penjualan tidak ditemukan")
	}
	if err := tx.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "5201").First(&akunHPP).Error; err != nil {
		return errors.New("akun HPP tidak ditemukan")
	}
	if err := tx.Where("id_koperasi = ? AND kode_akun = ?", idKoperasi, "1301").First(&akunPersediaan).Error; err != nil {
		return errors.New("akun persediaan tidak ditemukan")
	}

	// Calculate total HPP
	var totalHPP float64
	for _, item := range penjualan.ItemPenjualan {
		totalHPP += item.Produk.HargaBeli * float64(item.Kuantitas)
	}

	// Generate journal number within the same transaction
	tanggalStr := penjualan.TanggalPenjualan.Format("20060102")
	tanggalDate := penjualan.TanggalPenjualan.Format("2006-01-02")

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

	// Calculate totals for journal
	totalDebit := penjualan.TotalBelanja
	totalKredit := penjualan.TotalBelanja
	if totalHPP > 0 {
		totalDebit += totalHPP
		totalKredit += totalHPP
	}

	// Create journal entry
	transaksi := &models.Transaksi{
		IDKoperasi:       idKoperasi,
		NomorJurnal:      nomorJurnal,
		TanggalTransaksi: penjualan.TanggalPenjualan,
		Deskripsi:        fmt.Sprintf("Penjualan %s", penjualan.NomorPenjualan),
		NomorReferensi:   penjualan.NomorPenjualan,
		TipeTransaksi:    "penjualan",
		TotalDebit:       totalDebit,
		TotalKredit:      totalKredit,
		StatusBalanced:   true,
		DibuatOleh:       idPengguna,
	}

	if err := tx.Create(transaksi).Error; err != nil {
		return fmt.Errorf("gagal membuat jurnal: %w", err)
	}

	// Create journal lines
	barisTransaksi := []models.BarisTransaksi{
		// Kas bertambah (debit)
		{
			IDTransaksi:  transaksi.ID,
			IDAkun:       akunKas.ID,
			JumlahDebit:  penjualan.TotalBelanja,
			JumlahKredit: 0,
			Keterangan:   "Penerimaan kas dari penjualan",
		},
		// Penjualan bertambah (kredit)
		{
			IDTransaksi:  transaksi.ID,
			IDAkun:       akunPenjualan.ID,
			JumlahDebit:  0,
			JumlahKredit: penjualan.TotalBelanja,
			Keterangan:   "Pendapatan penjualan",
		},
	}

	// Add HPP entries if applicable
	if totalHPP > 0 {
		barisTransaksi = append(barisTransaksi,
			models.BarisTransaksi{
				IDTransaksi:  transaksi.ID,
				IDAkun:       akunHPP.ID,
				JumlahDebit:  totalHPP,
				JumlahKredit: 0,
				Keterangan:   "Harga Pokok Penjualan",
			},
			models.BarisTransaksi{
				IDTransaksi:  transaksi.ID,
				IDAkun:       akunPersediaan.ID,
				JumlahDebit:  0,
				JumlahKredit: totalHPP,
				Keterangan:   "Pengurangan persediaan",
			},
		)
	}

	for _, baris := range barisTransaksi {
		if err := tx.Create(&baris).Error; err != nil {
			return fmt.Errorf("gagal membuat baris jurnal: %w", err)
		}
	}

	// Update penjualan with transaction ID
	penjualan.IDTransaksi = &transaksi.ID
	if err := tx.Save(&penjualan).Error; err != nil {
		return fmt.Errorf("gagal update penjualan dengan ID transaksi: %w", err)
	}

	return nil
}
