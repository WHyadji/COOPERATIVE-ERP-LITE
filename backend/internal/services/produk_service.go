package services

import (
	"cooperative-erp-lite/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProdukService menangani logika bisnis produk
type ProdukService struct {
	db *gorm.DB
}

// NewProdukService membuat instance baru ProdukService
func NewProdukService(db *gorm.DB) *ProdukService {
	return &ProdukService{db: db}
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
	// Validasi kode produk unique
	var count int64
	s.db.Model(&models.Produk{}).
		Where("id_koperasi = ? AND kode_produk = ?", idKoperasi, req.KodeProduk).
		Count(&count)

	if count > 0 {
		return nil, errors.New("kode produk sudah digunakan")
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

	err := s.db.Create(produk).Error
	if err != nil {
		return nil, errors.New("gagal membuat produk")
	}

	response := produk.ToResponse()
	return &response, nil
}

// DapatkanSemuaProduk mengambil daftar produk dengan filter
func (s *ProdukService) DapatkanSemuaProduk(idKoperasi uuid.UUID, kategori, search string, statusAktif *bool, page, pageSize int) ([]models.ProdukResponse, int64, error) {
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
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("nama_produk ASC").Find(&produkList).Error

	if err != nil {
		return nil, 0, errors.New("gagal mengambil daftar produk")
	}

	// Convert to response
	responses := make([]models.ProdukResponse, len(produkList))
	for i, produk := range produkList {
		responses[i] = produk.ToResponse()
	}

	return responses, total, nil
}

// DapatkanProduk mengambil produk berdasarkan ID
func (s *ProdukService) DapatkanProduk(id uuid.UUID) (*models.ProdukResponse, error) {
	var produk models.Produk
	err := s.db.Where("id = ?", id).First(&produk).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	response := produk.ToResponse()
	return &response, nil
}

// DapatkanProdukByKode mengambil produk berdasarkan kode
func (s *ProdukService) DapatkanProdukByKode(idKoperasi uuid.UUID, kodeProduk string) (*models.ProdukResponse, error) {
	var produk models.Produk
	err := s.db.Where("id_koperasi = ? AND kode_produk = ?", idKoperasi, kodeProduk).First(&produk).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	response := produk.ToResponse()
	return &response, nil
}

// DapatkanProdukByBarcode mengambil produk berdasarkan barcode
func (s *ProdukService) DapatkanProdukByBarcode(idKoperasi uuid.UUID, barcode string) (*models.ProdukResponse, error) {
	var produk models.Produk
	err := s.db.Where("id_koperasi = ? AND barcode = ?", idKoperasi, barcode).First(&produk).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	response := produk.ToResponse()
	return &response, nil
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

// PerbaruiProduk mengupdate data produk
func (s *ProdukService) PerbaruiProduk(idKoperasi, id uuid.UUID, req *PerbaruiProdukRequest) (*models.ProdukResponse, error) {
	// Cek apakah produk ada DAN milik koperasi yang benar (multi-tenant validation)
	var produk models.Produk
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan atau tidak memiliki akses")
		}
		return nil, err
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
		return nil, errors.New("gagal memperbarui produk")
	}

	response := produk.ToResponse()
	return &response, nil
}

// HapusProduk menghapus produk (dengan validasi)
func (s *ProdukService) HapusProduk(idKoperasi, id uuid.UUID) error {
	// Cek apakah produk ada DAN milik koperasi yang benar (multi-tenant validation)
	var produk models.Produk
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&produk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("produk tidak ditemukan atau tidak memiliki akses")
		}
		return err
	}

	// Cek apakah ada di item penjualan
	var countPenjualan int64
	s.db.Model(&models.ItemPenjualan{}).Where("id_produk = ?", id).Count(&countPenjualan)

	if countPenjualan > 0 {
		return errors.New("tidak dapat menghapus produk yang sudah pernah dijual")
	}

	// Soft delete
	err = s.db.Delete(&produk).Error
	if err != nil {
		return errors.New("gagal menghapus produk")
	}

	return nil
}

// KurangiStok mengurangi stok produk
func (s *ProdukService) KurangiStok(id uuid.UUID, jumlah int) error {
	return s.KurangiStokWithTx(s.db, id, jumlah)
}

// KurangiStokWithTx mengurangi stok produk within an existing transaction
func (s *ProdukService) KurangiStokWithTx(tx *gorm.DB, id uuid.UUID, jumlah int) error {
	var produk models.Produk
	err := tx.Where("id = ?", id).First(&produk).Error
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	// Validasi stok cukup
	if produk.Stok < jumlah {
		return errors.New("stok tidak mencukupi")
	}

	// Kurangi stok
	produk.Stok -= jumlah
	err = tx.Save(&produk).Error
	if err != nil {
		return errors.New("gagal mengurangi stok")
	}

	return nil
}

// TambahStok menambah stok produk
func (s *ProdukService) TambahStok(id uuid.UUID, jumlah int) error {
	var produk models.Produk
	err := s.db.Where("id = ?", id).First(&produk).Error
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	// Tambah stok
	produk.Stok += jumlah
	err = s.db.Save(&produk).Error
	if err != nil {
		return errors.New("gagal menambah stok")
	}

	return nil
}

// CekStokTersedia mengecek apakah stok tersedia
func (s *ProdukService) CekStokTersedia(id uuid.UUID, jumlah int) (bool, error) {
	var produk models.Produk
	err := s.db.Where("id = ?", id).First(&produk).Error
	if err != nil {
		return false, errors.New("produk tidak ditemukan")
	}

	return produk.Stok >= jumlah, nil
}

// DapatkanProdukStokRendah mengambil produk dengan stok rendah
func (s *ProdukService) DapatkanProdukStokRendah(idKoperasi uuid.UUID) ([]models.ProdukResponse, error) {
	var produkList []models.Produk

	// Produk dengan stok <= stok minimum
	err := s.db.Where("id_koperasi = ? AND status_aktif = ? AND stok <= stok_minimum", idKoperasi, true).
		Order("stok ASC").
		Find(&produkList).Error

	if err != nil {
		return nil, errors.New("gagal mengambil daftar produk stok rendah")
	}

	// Convert to response
	responses := make([]models.ProdukResponse, len(produkList))
	for i, produk := range produkList {
		responses[i] = produk.ToResponse()
	}

	return responses, nil
}
