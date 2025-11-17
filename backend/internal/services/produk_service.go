package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProdukService menangani logika bisnis produk
type ProdukService struct {
	db     *gorm.DB
	logger *utils.Logger
}

// NewProdukService membuat instance baru ProdukService
func NewProdukService(db *gorm.DB) *ProdukService {
	return &ProdukService{
		db:     db,
		logger: utils.NewLogger("ProdukService"),
	}
}

// BuatProdukRequest adalah struktur request untuk membuat produk
type BuatProdukRequest struct {
	KodeProduk  string  `json:"kodeProduk" binding:"required"`
	NamaProduk  string  `json:"namaProduk" binding:"required"`
	Kategori    string  `json:"kategori"`
	Deskripsi   string  `json:"deskripsi"`
	Harga       float64 `json:"harga" binding:"required,gte=0"`
	HargaBeli   float64 `json:"hargaBeli" binding:"gte=0"`
	Stok        int     `json:"stok"`
	StokMinimum int     `json:"stokMinimum"`
	Satuan      string  `json:"satuan"`
	Barcode     string  `json:"barcode"`
	GambarURL   string  `json:"gambarUrl"`
}

// BuatProduk membuat produk baru
func (s *ProdukService) BuatProduk(idKoperasi uuid.UUID, req *BuatProdukRequest) (*models.ProdukResponse, error) {
	const method = "BuatProduk"

	// Validasi kode produk unique
	var jumlah int64
	err := s.db.Model(&models.Produk{}).
		Where("id_koperasi = ? AND kode_produk = ?", idKoperasi, req.KodeProduk).
		Count(&jumlah).Error

	if err != nil {
		s.logger.Error(method, "Gagal mengecek kode produk", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"kode_produk": req.KodeProduk,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengecek kode produk")
	}

	if jumlah > 0 {
		s.logger.Error(method, "Kode produk sudah digunakan", nil, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"kode_produk": req.KodeProduk,
		})
		return nil, utils.NewValidationError("Kode produk sudah digunakan")
	}

	// Buat produk
	produk := &models.Produk{
		IDKoperasi:  idKoperasi,
		KodeProduk:  req.KodeProduk,
		NamaProduk:  req.NamaProduk,
		Kategori:    req.Kategori,
		Deskripsi:   req.Deskripsi,
		Harga:       req.Harga,
		HargaBeli:   req.HargaBeli,
		Stok:        req.Stok,
		StokMinimum: req.StokMinimum,
		Satuan:      req.Satuan,
		Barcode:     req.Barcode,
		GambarURL:   req.GambarURL,
		StatusAktif: true,
	}

	err = s.db.Create(produk).Error
	if err != nil {
		s.logger.Error(method, "Gagal membuat produk", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"kode_produk": req.KodeProduk,
			"nama_produk": req.NamaProduk,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal membuat produk")
	}

	s.logger.Info(method, "Berhasil membuat produk", map[string]interface{}{
		"koperasi_id": idKoperasi.String(),
		"produk_id":   produk.ID.String(),
		"kode_produk": produk.KodeProduk,
		"nama_produk": produk.NamaProduk,
	})

	respons := produk.ToResponse()
	return &respons, nil
}

// DapatkanSemuaProduk mengambil daftar produk dengan filter
func (s *ProdukService) DapatkanSemuaProduk(idKoperasi uuid.UUID, kategori, search string, statusAktif *bool, page, pageSize int) ([]models.ProdukResponse, int64, error) {
	const method = "DapatkanSemuaProduk"

	var produkList []models.Produk
	var total int64

	query := s.db.Model(&models.Produk{}).Where("id_koperasi = ?", idKoperasi)

	// Apply filters
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	if search != "" {
		query = query.Where("nama_produk ILIKE ? OR kode_produk ILIKE ? OR barcode ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if statusAktif != nil {
		query = query.Where("status_aktif = ?", *statusAktif)
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		s.logger.Error(method, "Gagal menghitung total produk", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"kategori":    kategori,
			"search":      search,
		})
		return nil, 0, utils.WrapDatabaseError(err, "Gagal menghitung total produk")
	}

	// Pagination
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("nama_produk ASC").Find(&produkList).Error

	if err != nil {
		s.logger.Error(method, "Gagal mengambil daftar produk", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"kategori":    kategori,
			"search":      search,
			"page":        page,
			"page_size":   pageSize,
		})
		return nil, 0, utils.WrapDatabaseError(err, "Gagal mengambil daftar produk")
	}

	s.logger.Debug(method, "Berhasil mengambil daftar produk", map[string]interface{}{
		"koperasi_id":   idKoperasi.String(),
		"total":         total,
		"jumlah_result": len(produkList),
		"page":          page,
	})

	// Convert to response
	responseDaftar := make([]models.ProdukResponse, len(produkList))
	for i, produk := range produkList {
		responseDaftar[i] = produk.ToResponse()
	}

	return responseDaftar, total, nil
}

// DapatkanProduk mengambil produk berdasarkan ID dengan validasi multi-tenant
func (s *ProdukService) DapatkanProduk(idKoperasi, id uuid.UUID) (*models.ProdukResponse, error) {
	const method = "DapatkanProduk"

	var produk models.Produk
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&produk).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan atau tidak memiliki akses", err, map[string]interface{}{
				"produk_id":   id.String(),
				"koperasi_id": idKoperasi.String(),
			})
			return nil, utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"produk_id":   id.String(),
			"koperasi_id": idKoperasi.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	s.logger.Debug(method, "Berhasil mengambil data produk", map[string]interface{}{
		"produk_id":   id.String(),
		"koperasi_id": idKoperasi.String(),
		"nama_produk": produk.NamaProduk,
		"kode_produk": produk.KodeProduk,
	})

	respons := produk.ToResponse()
	return &respons, nil
}

// DapatkanProdukByKode mengambil produk berdasarkan kode
func (s *ProdukService) DapatkanProdukByKode(idKoperasi uuid.UUID, kodeProduk string) (*models.ProdukResponse, error) {
	const method = "DapatkanProdukByKode"

	var produk models.Produk
	err := s.db.Where("id_koperasi = ? AND kode_produk = ?", idKoperasi, kodeProduk).First(&produk).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"kode_produk": kodeProduk,
			})
			return nil, utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"kode_produk": kodeProduk,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	s.logger.Debug(method, "Berhasil mengambil data produk", map[string]interface{}{
		"koperasi_id": idKoperasi.String(),
		"produk_id":   produk.ID.String(),
		"kode_produk": kodeProduk,
		"nama_produk": produk.NamaProduk,
	})

	respons := produk.ToResponse()
	return &respons, nil
}

// DapatkanProdukByBarcode mengambil produk berdasarkan barcode
func (s *ProdukService) DapatkanProdukByBarcode(idKoperasi uuid.UUID, barcode string) (*models.ProdukResponse, error) {
	const method = "DapatkanProdukByBarcode"

	var produk models.Produk
	err := s.db.Where("id_koperasi = ? AND barcode = ?", idKoperasi, barcode).First(&produk).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan", err, map[string]interface{}{
				"koperasi_id": idKoperasi.String(),
				"barcode":     barcode,
			})
			return nil, utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
			"barcode":     barcode,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	s.logger.Debug(method, "Berhasil mengambil data produk", map[string]interface{}{
		"koperasi_id": idKoperasi.String(),
		"produk_id":   produk.ID.String(),
		"barcode":     barcode,
		"nama_produk": produk.NamaProduk,
	})

	respons := produk.ToResponse()
	return &respons, nil
}

// PerbaruiProdukRequest adalah struktur request untuk update produk
type PerbaruiProdukRequest struct {
	NamaProduk  string  `json:"namaProduk"`
	Kategori    string  `json:"kategori"`
	Deskripsi   string  `json:"deskripsi"`
	Harga       float64 `json:"harga"`
	HargaBeli   float64 `json:"hargaBeli"`
	StokMinimum int     `json:"stokMinimum"`
	Satuan      string  `json:"satuan"`
	Barcode     string  `json:"barcode"`
	GambarURL   string  `json:"gambarUrl"`
	StatusAktif *bool   `json:"statusAktif"`
}

// PerbaruiProduk mengupdate data produk dengan validasi multi-tenant
func (s *ProdukService) PerbaruiProduk(idKoperasi, id uuid.UUID, req *PerbaruiProdukRequest) (*models.ProdukResponse, error) {
	const method = "PerbaruiProduk"

	var produk models.Produk
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan atau tidak memiliki akses", err, map[string]interface{}{
				"produk_id":   id.String(),
				"koperasi_id": idKoperasi.String(),
			})
			return nil, utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"produk_id":   id.String(),
			"koperasi_id": idKoperasi.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	// Update fields
	if req.NamaProduk != "" {
		produk.NamaProduk = req.NamaProduk
	}
	if req.Kategori != "" {
		produk.Kategori = req.Kategori
	}
	if req.Deskripsi != "" {
		produk.Deskripsi = req.Deskripsi
	}
	if req.Harga > 0 {
		produk.Harga = req.Harga
	}
	if req.HargaBeli >= 0 {
		produk.HargaBeli = req.HargaBeli
	}
	if req.StokMinimum >= 0 {
		produk.StokMinimum = req.StokMinimum
	}
	if req.Satuan != "" {
		produk.Satuan = req.Satuan
	}
	if req.Barcode != "" {
		produk.Barcode = req.Barcode
	}
	if req.GambarURL != "" {
		produk.GambarURL = req.GambarURL
	}
	if req.StatusAktif != nil {
		produk.StatusAktif = *req.StatusAktif
	}

	err = s.db.Save(&produk).Error
	if err != nil {
		s.logger.Error(method, "Gagal memperbarui produk", err, map[string]interface{}{
			"produk_id":   id.String(),
			"nama_produk": req.NamaProduk,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal memperbarui produk")
	}

	s.logger.Info(method, "Berhasil memperbarui produk", map[string]interface{}{
		"produk_id":   id.String(),
		"nama_produk": produk.NamaProduk,
		"kode_produk": produk.KodeProduk,
	})

	respons := produk.ToResponse()
	return &respons, nil
}

// HapusProduk menghapus produk (dengan validasi multi-tenant)
func (s *ProdukService) HapusProduk(idKoperasi, id uuid.UUID) error {
	const method = "HapusProduk"

	// Ambil data produk terlebih dahulu untuk logging dengan validasi multi-tenant
	var produk models.Produk
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan atau tidak memiliki akses", err, map[string]interface{}{
				"produk_id":   id.String(),
				"koperasi_id": idKoperasi.String(),
			})
			return utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"produk_id":   id.String(),
			"koperasi_id": idKoperasi.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	// Cek apakah ada di item penjualan
	var jumlahPenjualan int64
	err = s.db.Model(&models.ItemPenjualan{}).Where("id_produk = ?", id).Count(&jumlahPenjualan).Error
	if err != nil {
		s.logger.Error(method, "Gagal memeriksa item penjualan", err, map[string]interface{}{
			"produk_id": id.String(),
		})
		return utils.WrapDatabaseError(err, "Gagal memeriksa item penjualan")
	}

	if jumlahPenjualan > 0 {
		s.logger.Error(method, "Tidak dapat menghapus produk yang sudah pernah dijual", nil, map[string]interface{}{
			"produk_id":       id.String(),
			"nama_produk":     produk.NamaProduk,
			"count_penjualan": jumlahPenjualan,
		})
		return utils.NewValidationError("Tidak dapat menghapus produk yang sudah pernah dijual")
	}

	// Soft delete
	err = s.db.Delete(&models.Produk{}, id).Error
	if err != nil {
		s.logger.Error(method, "Gagal menghapus produk", err, map[string]interface{}{
			"produk_id":   id.String(),
			"nama_produk": produk.NamaProduk,
		})
		return utils.WrapDatabaseError(err, "Gagal menghapus produk")
	}

	s.logger.Info(method, "Berhasil menghapus produk", map[string]interface{}{
		"produk_id":   id.String(),
		"nama_produk": produk.NamaProduk,
		"kode_produk": produk.KodeProduk,
	})

	return nil
}

// KurangiStok mengurangi stok produk
func (s *ProdukService) KurangiStok(id uuid.UUID, jumlah int) error {
	const method = "KurangiStok"

	var produk models.Produk
	err := s.db.Where("id = ?", id).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan", err, map[string]interface{}{
				"produk_id": id.String(),
				"kuantitas": jumlah,
			})
			return utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"produk_id": id.String(),
			"kuantitas": jumlah,
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	// Validasi stok cukup
	if produk.Stok < jumlah {
		s.logger.Error(method, "Stok tidak mencukupi", nil, map[string]interface{}{
			"produk_id":         id.String(),
			"nama_produk":       produk.NamaProduk,
			"stok_saat_ini":     produk.Stok,
			"kuantitas_diminta": jumlah,
		})
		return utils.NewValidationError("Stok tidak mencukupi")
	}

	// Kurangi stok
	stokLama := produk.Stok
	produk.Stok -= jumlah
	err = s.db.Save(&produk).Error
	if err != nil {
		s.logger.Error(method, "Gagal mengurangi stok", err, map[string]interface{}{
			"produk_id":   id.String(),
			"nama_produk": produk.NamaProduk,
			"kuantitas":   jumlah,
		})
		return utils.WrapDatabaseError(err, "Gagal mengurangi stok")
	}

	s.logger.Info(method, "Berhasil mengurangi stok", map[string]interface{}{
		"produk_id":   id.String(),
		"nama_produk": produk.NamaProduk,
		"stok_lama":   stokLama,
		"stok_baru":   produk.Stok,
		"kuantitas":   jumlah,
	})

	return nil
}

// TambahStok menambah stok produk
func (s *ProdukService) TambahStok(id uuid.UUID, jumlah int) error {
	const method = "TambahStok"

	var produk models.Produk
	err := s.db.Where("id = ?", id).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan", err, map[string]interface{}{
				"produk_id": id.String(),
				"kuantitas": jumlah,
			})
			return utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"produk_id": id.String(),
			"kuantitas": jumlah,
		})
		return utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	// Tambah stok
	stokLama := produk.Stok
	produk.Stok += jumlah
	err = s.db.Save(&produk).Error
	if err != nil {
		s.logger.Error(method, "Gagal menambah stok", err, map[string]interface{}{
			"produk_id":   id.String(),
			"nama_produk": produk.NamaProduk,
			"kuantitas":   jumlah,
		})
		return utils.WrapDatabaseError(err, "Gagal menambah stok")
	}

	s.logger.Info(method, "Berhasil menambah stok", map[string]interface{}{
		"produk_id":   id.String(),
		"nama_produk": produk.NamaProduk,
		"stok_lama":   stokLama,
		"stok_baru":   produk.Stok,
		"kuantitas":   jumlah,
	})

	return nil
}

// CekStokTersedia mengecek apakah stok tersedia
func (s *ProdukService) CekStokTersedia(id uuid.UUID, jumlah int) (bool, error) {
	const method = "CekStokTersedia"

	var produk models.Produk
	err := s.db.Where("id = ?", id).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan", err, map[string]interface{}{
				"produk_id": id.String(),
				"kuantitas": jumlah,
			})
			return false, utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"produk_id": id.String(),
			"kuantitas": jumlah,
		})
		return false, utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	tersedia := produk.Stok >= jumlah

	s.logger.Debug(method, "Pengecekan stok selesai", map[string]interface{}{
		"produk_id":         id.String(),
		"nama_produk":       produk.NamaProduk,
		"stok_saat_ini":     produk.Stok,
		"kuantitas_diminta": jumlah,
		"tersedia":          tersedia,
	})

	return tersedia, nil
}

// AdjustStok menyesuaikan stok produk (tambah atau kurang) dengan keterangan
func (s *ProdukService) AdjustStok(idKoperasi, id uuid.UUID, jumlah int, keterangan string) (*models.ProdukResponse, error) {
	const method = "AdjustStok"

	// Validasi produk exists dan milik koperasi yang benar (multi-tenant)
	var produk models.Produk
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error(method, "Produk tidak ditemukan atau tidak memiliki akses", err, map[string]interface{}{
				"produk_id":   id.String(),
				"koperasi_id": idKoperasi.String(),
			})
			return nil, utils.WrapDatabaseError(err, "Produk")
		}
		s.logger.Error(method, "Gagal mengambil data produk", err, map[string]interface{}{
			"produk_id":   id.String(),
			"koperasi_id": idKoperasi.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil data produk")
	}

	// Validasi jumlah tidak boleh 0
	if jumlah == 0 {
		s.logger.Error(method, "Jumlah adjustment tidak boleh 0", nil, map[string]interface{}{
			"produk_id": id.String(),
			"jumlah":    jumlah,
		})
		return nil, utils.NewValidationError("Jumlah adjustment tidak boleh 0")
	}

	// Jika pengurangan, validasi stok cukup
	if jumlah < 0 && produk.Stok < -jumlah {
		s.logger.Error(method, "Stok tidak mencukupi untuk pengurangan", nil, map[string]interface{}{
			"produk_id":          id.String(),
			"nama_produk":        produk.NamaProduk,
			"stok_saat_ini":      produk.Stok,
			"jumlah_pengurangan": -jumlah,
		})
		return nil, utils.NewValidationError("Stok tidak mencukupi")
	}

	// Adjust stok
	stokLama := produk.Stok
	produk.Stok += jumlah

	err = s.db.Save(&produk).Error
	if err != nil {
		s.logger.Error(method, "Gagal menyesuaikan stok", err, map[string]interface{}{
			"produk_id":   id.String(),
			"nama_produk": produk.NamaProduk,
			"jumlah":      jumlah,
		})
		return nil, utils.WrapDatabaseError(err, "Gagal menyesuaikan stok")
	}

	jenisAdjustment := "penambahan"
	if jumlah < 0 {
		jenisAdjustment = "pengurangan"
	}

	s.logger.Info(method, "Berhasil menyesuaikan stok", map[string]interface{}{
		"produk_id":   id.String(),
		"nama_produk": produk.NamaProduk,
		"stok_lama":   stokLama,
		"stok_baru":   produk.Stok,
		"jumlah":      jumlah,
		"jenis":       jenisAdjustment,
		"keterangan":  keterangan,
	})

	// Return response
	return &models.ProdukResponse{
		ID:          produk.ID,
		KodeProduk:  produk.KodeProduk,
		NamaProduk:  produk.NamaProduk,
		Kategori:    produk.Kategori,
		Deskripsi:   produk.Deskripsi,
		Harga:       produk.Harga,
		HargaBeli:   produk.HargaBeli,
		Stok:        produk.Stok,
		StokMinimum: produk.StokMinimum,
		Satuan:      produk.Satuan,
		Barcode:     produk.Barcode,
		GambarURL:   produk.GambarURL,
		StatusAktif: produk.StatusAktif,
	}, nil
}

// DapatkanProdukStokRendah mengambil produk dengan stok rendah
func (s *ProdukService) DapatkanProdukStokRendah(idKoperasi uuid.UUID) ([]models.ProdukResponse, error) {
	const method = "DapatkanProdukStokRendah"

	var produkList []models.Produk

	// Produk dengan stok <= stok minimum
	err := s.db.Where("id_koperasi = ? AND status_aktif = ? AND stok <= stok_minimum", idKoperasi, true).
		Order("stok ASC").
		Find(&produkList).Error

	if err != nil {
		s.logger.Error(method, "Gagal mengambil daftar produk stok rendah", err, map[string]interface{}{
			"koperasi_id": idKoperasi.String(),
		})
		return nil, utils.WrapDatabaseError(err, "Gagal mengambil daftar produk stok rendah")
	}

	s.logger.Debug(method, "Berhasil mengambil daftar produk stok rendah", map[string]interface{}{
		"koperasi_id": idKoperasi.String(),
		"jumlah":      len(produkList),
	})

	// Convert to response
	responseDaftar := make([]models.ProdukResponse, len(produkList))
	for i, produk := range produkList {
		responseDaftar[i] = produk.ToResponse()
	}

	return responseDaftar, nil
}
