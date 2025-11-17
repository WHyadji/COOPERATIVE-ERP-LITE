package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/utils"
	"errors"

	"gorm.io/gorm"
)

// AuthService menangani logika bisnis autentikasi
type AuthService struct {
	db      *gorm.DB
	jwtUtil *utils.JWTUtil
}

// NewAuthService membuat instance baru AuthService
func NewAuthService(db *gorm.DB, jwtUtil *utils.JWTUtil) *AuthService {
	return &AuthService{
		db:      db,
		jwtUtil: jwtUtil,
	}
}

// LoginRequest adalah struktur request untuk login
type LoginRequest struct {
	NamaPengguna string `json:"namaPengguna" binding:"required"`
	KataSandi    string `json:"kataSandi" binding:"required"`
}

// LoginResponse adalah struktur response untuk login
type LoginResponse struct {
	Token    string                    `json:"token"`
	Pengguna models.PenggunaResponse   `json:"pengguna"`
}

// Login melakukan autentikasi pengguna dan menghasilkan JWT token
func (s *AuthService) Login(namaPengguna, kataSandi string) (*LoginResponse, error) {
	// Cari pengguna berdasarkan username
	var pengguna models.Pengguna
	err := s.db.Where("nama_pengguna = ? AND status_aktif = ?", namaPengguna, true).
		First(&pengguna).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("nama pengguna atau kata sandi salah")
		}
		return nil, err
	}

	// Verifikasi password
	if !pengguna.CekKataSandi(kataSandi) {
		return nil, errors.New("nama pengguna atau kata sandi salah")
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(&pengguna)
	if err != nil {
		return nil, errors.New("gagal membuat token autentikasi")
	}

	// Buat response
	response := &LoginResponse{
		Token:    token,
		Pengguna: pengguna.ToResponse(),
	}

	return response, nil
}

// ValidasiToken memvalidasi JWT token dan mengembalikan claims
func (s *AuthService) ValidasiToken(tokenString string) (*utils.JWTClaims, error) {
	claims, err := s.jwtUtil.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("token tidak valid atau sudah kadaluarsa")
	}

	// Optional: Cek apakah pengguna masih aktif di database
	var pengguna models.Pengguna
	err = s.db.Where("id = ? AND status_aktif = ?", claims.IDPengguna, true).
		First(&pengguna).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan atau tidak aktif")
		}
		return nil, err
	}

	return claims, nil
}

// RefreshToken membuat token baru dari token yang sudah ada
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	// Validasi token lama terlebih dahulu
	claims, err := s.ValidasiToken(tokenString)
	if err != nil {
		return "", err
	}

	// Dapatkan data pengguna terbaru dari database
	var pengguna models.Pengguna
	err = s.db.Where("id = ?", claims.IDPengguna).First(&pengguna).Error
	if err != nil {
		return "", errors.New("pengguna tidak ditemukan")
	}

	// Generate token baru
	newToken, err := s.jwtUtil.GenerateToken(&pengguna)
	if err != nil {
		return "", errors.New("gagal membuat token baru")
	}

	return newToken, nil
}

// DapatkanProfilPengguna mengambil profil pengguna berdasarkan ID dari token
func (s *AuthService) DapatkanProfilPengguna(idPengguna string) (*models.PenggunaResponse, error) {
	var pengguna models.Pengguna
	err := s.db.Where("id = ?", idPengguna).First(&pengguna).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, err
	}

	response := pengguna.ToResponse()
	return &response, nil
}

// UbahKataSandi mengubah password pengguna
func (s *AuthService) UbahKataSandi(idPengguna, kataSandiLama, kataSandiBaru string) error {
	// Dapatkan pengguna
	var pengguna models.Pengguna
	err := s.db.Where("id = ?", idPengguna).First(&pengguna).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("pengguna tidak ditemukan")
		}
		return err
	}

	// Verifikasi password lama
	if !pengguna.CekKataSandi(kataSandiLama) {
		return errors.New("kata sandi lama tidak sesuai")
	}

	// Validasi password baru
	if err := utils.ValidasiKataSandi(kataSandiBaru); err != nil {
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
