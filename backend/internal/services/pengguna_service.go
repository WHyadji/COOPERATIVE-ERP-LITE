package services

import (
	"cooperative-erp-lite/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PenggunaService menangani logika bisnis pengguna
type PenggunaService struct {
	db *gorm.DB
}

// NewPenggunaService membuat instance baru PenggunaService
func NewPenggunaService(db *gorm.DB) *PenggunaService {
	return &PenggunaService{db: db}
}

// BuatPenggunaRequest adalah struktur request untuk membuat pengguna
type BuatPenggunaRequest struct {
	NamaLengkap  string               `json:"namaLengkap" binding:"required"`
	NamaPengguna string               `json:"namaPengguna" binding:"required,min=3"`
	Email        string               `json:"email" binding:"required,email"`
	KataSandi    string               `json:"kataSandi" binding:"required,min=6"`
	Peran        models.PeranPengguna `json:"peran" binding:"required"`
}

// BuatPengguna membuat pengguna baru
func (s *PenggunaService) BuatPengguna(idKoperasi uuid.UUID, req *BuatPenggunaRequest) (*models.PenggunaResponse, error) {
	// Cek apakah username sudah ada di koperasi yang sama
	var count int64
	s.db.Model(&models.Pengguna{}).
		Where("id_koperasi = ? AND nama_pengguna = ?", idKoperasi, req.NamaPengguna).
		Count(&count)

	if count > 0 {
		return nil, errors.New("nama pengguna sudah digunakan")
	}

	// Buat pengguna baru
	pengguna := &models.Pengguna{
		IDKoperasi:   idKoperasi,
		NamaLengkap:  req.NamaLengkap,
		NamaPengguna: req.NamaPengguna,
		Email:        req.Email,
		Peran:        req.Peran,
		StatusAktif:  true,
	}

	// Hash password
	err := pengguna.SetKataSandi(req.KataSandi)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi kata sandi")
	}

	// Simpan ke database
	err = s.db.Create(pengguna).Error
	if err != nil {
		return nil, errors.New("gagal membuat pengguna")
	}

	response := pengguna.ToResponse()
	return &response, nil
}

// DapatkanSemuaPengguna mengambil daftar pengguna dengan filter
func (s *PenggunaService) DapatkanSemuaPengguna(idKoperasi uuid.UUID, peran string, statusAktif *bool, page, pageSize int) ([]models.PenggunaResponse, int64, error) {
	var penggunaList []models.Pengguna
	var total int64

	query := s.db.Model(&models.Pengguna{}).Where("id_koperasi = ?", idKoperasi)

	// Apply filters
	if peran != "" {
		query = query.Where("peran = ?", peran)
	}
	if statusAktif != nil {
		query = query.Where("status_aktif = ?", *statusAktif)
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("tanggal_dibuat DESC").Find(&penggunaList).Error

	if err != nil {
		return nil, 0, errors.New("gagal mengambil daftar pengguna")
	}

	// Convert to response
	responses := make([]models.PenggunaResponse, len(penggunaList))
	for i, pengguna := range penggunaList {
		responses[i] = pengguna.ToResponse()
	}

	return responses, total, nil
}

// DapatkanPengguna mengambil data pengguna berdasarkan ID
func (s *PenggunaService) DapatkanPengguna(id uuid.UUID) (*models.PenggunaResponse, error) {
	var pengguna models.Pengguna
	err := s.db.Where("id = ?", id).First(&pengguna).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, err
	}

	response := pengguna.ToResponse()
	return &response, nil
}

// PerbaruiPenggunaRequest adalah struktur request untuk update pengguna
type PerbaruiPenggunaRequest struct {
	NamaLengkap string               `json:"namaLengkap"`
	Email       string               `json:"email"`
	Peran       models.PeranPengguna `json:"peran"`
	StatusAktif *bool                `json:"statusAktif"`
}

// PerbaruiPengguna mengupdate data pengguna
func (s *PenggunaService) PerbaruiPengguna(idKoperasi, id uuid.UUID, req *PerbaruiPenggunaRequest) (*models.PenggunaResponse, error) {
	// Cek apakah pengguna ada DAN milik koperasi yang benar (multi-tenant validation)
	var pengguna models.Pengguna
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&pengguna).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan atau tidak memiliki akses")
		}
		return nil, err
	}

	// Update fields yang disediakan
	if req.NamaLengkap != "" {
		pengguna.NamaLengkap = req.NamaLengkap
	}
	if req.Email != "" {
		pengguna.Email = req.Email
	}
	if req.Peran != "" {
		pengguna.Peran = req.Peran
	}
	if req.StatusAktif != nil {
		pengguna.StatusAktif = *req.StatusAktif
	}

	// Simpan perubahan
	err = s.db.Save(&pengguna).Error
	if err != nil {
		return nil, errors.New("gagal memperbarui pengguna")
	}

	response := pengguna.ToResponse()
	return &response, nil
}

// HapusPengguna menghapus pengguna (soft delete)
func (s *PenggunaService) HapusPengguna(idKoperasi, id uuid.UUID) error {
	// Cek apakah pengguna ada DAN milik koperasi yang benar (multi-tenant validation)
	var pengguna models.Pengguna
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&pengguna).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("pengguna tidak ditemukan atau tidak memiliki akses")
		}
		return err
	}

	// Soft delete
	err = s.db.Delete(&pengguna).Error
	if err != nil {
		return errors.New("gagal menghapus pengguna")
	}

	return nil
}

// UbahKataSandiPengguna mengubah password pengguna (oleh admin)
func (s *PenggunaService) UbahKataSandiPengguna(idKoperasi, id uuid.UUID, kataSandiBaru string) error {
	// Validasi password
	if len(kataSandiBaru) < 6 {
		return errors.New("kata sandi minimal 6 karakter")
	}

	// Cek apakah pengguna ada DAN milik koperasi yang benar (multi-tenant validation)
	var pengguna models.Pengguna
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&pengguna).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("pengguna tidak ditemukan atau tidak memiliki akses")
		}
		return err
	}

	// Set password baru
	err = pengguna.SetKataSandi(kataSandiBaru)
	if err != nil {
		return errors.New("gagal mengenkripsi kata sandi baru")
	}

	// Simpan ke database
	err = s.db.Save(&pengguna).Error
	if err != nil {
		return errors.New("gagal menyimpan kata sandi baru")
	}

	return nil
}

// ResetKataSandi mereset password pengguna ke default (admin only)
func (s *PenggunaService) ResetKataSandi(idKoperasi, id uuid.UUID) (string, error) {
	// Generate password default (menggunakan username)
	var pengguna models.Pengguna
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&pengguna).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("pengguna tidak ditemukan atau tidak memiliki akses")
		}
		return "", err
	}

	// Password default: username + "123"
	passwordDefault := fmt.Sprintf("%s123", pengguna.NamaPengguna)

	// Set password
	err = pengguna.SetKataSandi(passwordDefault)
	if err != nil {
		return "", errors.New("gagal mengenkripsi kata sandi")
	}

	// Simpan
	err = s.db.Save(&pengguna).Error
	if err != nil {
		return "", errors.New("gagal mereset kata sandi")
	}

	return passwordDefault, nil
}

// DapatkanPenggunaByUsername mengambil pengguna berdasarkan username
func (s *PenggunaService) DapatkanPenggunaByUsername(idKoperasi uuid.UUID, namaPengguna string) (*models.PenggunaResponse, error) {
	var pengguna models.Pengguna
	err := s.db.Where("id_koperasi = ? AND nama_pengguna = ?", idKoperasi, namaPengguna).First(&pengguna).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, err
	}

	response := pengguna.ToResponse()
	return &response, nil
}
