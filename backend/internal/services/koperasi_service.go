package services

import (
	"cooperative-erp-lite/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KoperasiService menangani logika bisnis koperasi
type KoperasiService struct {
	db *gorm.DB
}

// NewKoperasiService membuat instance baru KoperasiService
func NewKoperasiService(db *gorm.DB) *KoperasiService {
	return &KoperasiService{db: db}
}

// BuatKoperasiRequest adalah struktur request untuk membuat koperasi
type BuatKoperasiRequest struct {
	NamaKoperasi   string `json:"namaKoperasi" binding:"required"`
	Alamat         string `json:"alamat"`
	NoTelepon      string `json:"noTelepon"`
	Email          string `json:"email"`
	TahunBukuMulai int    `json:"tahunBukuMulai"`
}

// BuatKoperasi membuat koperasi baru
func (s *KoperasiService) BuatKoperasi(req *BuatKoperasiRequest) (*models.Koperasi, error) {
	koperasi := &models.Koperasi{
		NamaKoperasi:   req.NamaKoperasi,
		Alamat:         req.Alamat,
		NoTelepon:      req.NoTelepon,
		Email:          req.Email,
		TahunBukuMulai: req.TahunBukuMulai,
	}

	// Simpan ke database
	err := s.db.Create(koperasi).Error
	if err != nil {
		return nil, errors.New("gagal membuat koperasi")
	}

	return koperasi, nil
}

// DapatkanKoperasi mengambil data koperasi berdasarkan ID
func (s *KoperasiService) DapatkanKoperasi(idKoperasi uuid.UUID) (*models.Koperasi, error) {
	var koperasi models.Koperasi
	err := s.db.Where("id = ?", idKoperasi).First(&koperasi).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("koperasi tidak ditemukan")
		}
		return nil, err
	}

	return &koperasi, nil
}

// PerbaruiKoperasiRequest adalah struktur request untuk update koperasi
type PerbaruiKoperasiRequest struct {
	NamaKoperasi   string `json:"namaKoperasi"`
	Alamat         string `json:"alamat"`
	NoTelepon      string `json:"noTelepon"`
	Email          string `json:"email"`
	LogoURL        string `json:"logoUrl"`
	TahunBukuMulai int    `json:"tahunBukuMulai"`
}

// PerbaruiKoperasi mengupdate data koperasi
func (s *KoperasiService) PerbaruiKoperasi(idKoperasi uuid.UUID, req *PerbaruiKoperasiRequest) (*models.Koperasi, error) {
	// Cek apakah koperasi ada
	koperasi, err := s.DapatkanKoperasi(idKoperasi)
	if err != nil {
		return nil, err
	}

	// Update fields yang disediakan
	if req.NamaKoperasi != "" {
		koperasi.NamaKoperasi = req.NamaKoperasi
	}
	if req.Alamat != "" {
		koperasi.Alamat = req.Alamat
	}
	if req.NoTelepon != "" {
		koperasi.NoTelepon = req.NoTelepon
	}
	if req.Email != "" {
		koperasi.Email = req.Email
	}
	if req.LogoURL != "" {
		koperasi.LogoURL = req.LogoURL
	}
	if req.TahunBukuMulai > 0 {
		koperasi.TahunBukuMulai = req.TahunBukuMulai
	}

	// Simpan perubahan
	err = s.db.Save(koperasi).Error
	if err != nil {
		return nil, errors.New("gagal memperbarui koperasi")
	}

	return koperasi, nil
}

// DapatkanSemuaKoperasi mengambil daftar semua koperasi
func (s *KoperasiService) DapatkanSemuaKoperasi() ([]models.Koperasi, error) {
	var koperasiList []models.Koperasi
	err := s.db.Find(&koperasiList).Error

	if err != nil {
		return nil, errors.New("gagal mengambil daftar koperasi")
	}

	return koperasiList, nil
}

// HapusKoperasi menghapus koperasi (soft delete)
func (s *KoperasiService) HapusKoperasi(idKoperasi uuid.UUID) error {
	// Cek apakah koperasi ada
	_, err := s.DapatkanKoperasi(idKoperasi)
	if err != nil {
		return err
	}

	// Soft delete
	err = s.db.Delete(&models.Koperasi{}, idKoperasi).Error
	if err != nil {
		return errors.New("gagal menghapus koperasi")
	}

	return nil
}

// DapatkanStatistikKoperasi mengambil statistik koperasi
func (s *KoperasiService) DapatkanStatistikKoperasi(idKoperasi uuid.UUID) (map[string]interface{}, error) {
	// Cek apakah koperasi ada
	_, err := s.DapatkanKoperasi(idKoperasi)
	if err != nil {
		return nil, err
	}

	statistik := make(map[string]interface{})

	// Hitung jumlah anggota
	var jumlahAnggota int64
	s.db.Model(&models.Anggota{}).Where("id_koperasi = ? AND status = ?", idKoperasi, models.StatusAktif).Count(&jumlahAnggota)
	statistik["jumlahAnggota"] = jumlahAnggota

	// Hitung jumlah pengguna
	var jumlahPengguna int64
	s.db.Model(&models.Pengguna{}).Where("id_koperasi = ? AND status_aktif = ?", idKoperasi, true).Count(&jumlahPengguna)
	statistik["jumlahPengguna"] = jumlahPengguna

	// Hitung jumlah produk
	var jumlahProduk int64
	s.db.Model(&models.Produk{}).Where("id_koperasi = ? AND status_aktif = ?", idKoperasi, true).Count(&jumlahProduk)
	statistik["jumlahProduk"] = jumlahProduk

	return statistik, nil
}

// GetSemuaKoperasi is an English wrapper for DapatkanSemuaKoperasi with pagination
func (s *KoperasiService) GetSemuaKoperasi(page, pageSize int) ([]models.Koperasi, int64, error) {
	// Get all koperasi
	koperasiList, err := s.DapatkanSemuaKoperasi()
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]models.Koperasi, len(koperasiList))
	copy(responses, koperasiList)

	total := int64(len(responses))

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(responses) {
		return []models.Koperasi{}, total, nil
	}
	if end > len(responses) {
		end = len(responses)
	}

	return responses[start:end], total, nil
}

// GetKoperasiByID is an English wrapper for DapatkanKoperasi
func (s *KoperasiService) GetKoperasiByID(id uuid.UUID) (*models.Koperasi, error) {
	koperasi, err := s.DapatkanKoperasi(id)
	if err != nil {
		return nil, err
	}
	response := *koperasi
	return &response, nil
}

// GetStatistikKoperasi is an English wrapper for DapatkanStatistikKoperasi
func (s *KoperasiService) GetStatistikKoperasi(idKoperasi uuid.UUID) (map[string]interface{}, error) {
	return s.DapatkanStatistikKoperasi(idKoperasi)
}
