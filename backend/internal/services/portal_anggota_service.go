package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/utils"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PortalAnggotaService menangani logika bisnis portal anggota
type PortalAnggotaService struct {
	db      *gorm.DB
	jwtUtil *utils.JWTUtil
}

// NewPortalAnggotaService membuat instance baru PortalAnggotaService
func NewPortalAnggotaService(db *gorm.DB, jwtUtil *utils.JWTUtil) *PortalAnggotaService {
	return &PortalAnggotaService{
		db:      db,
		jwtUtil: jwtUtil,
	}
}

// LoginAnggotaRequest adalah struktur request untuk login portal anggota
type LoginAnggotaRequest struct {
	NomorAnggota string `json:"nomorAnggota" binding:"required"`
	PIN          string `json:"pin" binding:"required,len=6"`
}

// LoginAnggotaResponse adalah struktur response untuk login portal anggota
type LoginAnggotaResponse struct {
	Token   string                  `json:"token"`
	Anggota models.AnggotaResponse `json:"anggota"`
}

// LoginAnggota melakukan autentikasi anggota dengan nomor anggota dan PIN
func (s *PortalAnggotaService) LoginAnggota(idKoperasi uuid.UUID, nomorAnggota, pin string) (*LoginAnggotaResponse, error) {
	// Cari anggota berdasarkan nomor anggota
	var anggota models.Anggota
	err := s.db.Where("id_koperasi = ? AND nomor_anggota = ? AND status = ?",
		idKoperasi, nomorAnggota, models.StatusAktif).
		First(&anggota).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("nomor anggota tidak ditemukan atau tidak aktif")
		}
		return nil, err
	}

	// Cek apakah PIN sudah diset
	if anggota.PINPortal == "" {
		return nil, errors.New("PIN portal belum diset. Silakan hubungi admin koperasi")
	}

	// Verifikasi PIN
	err = bcrypt.CompareHashAndPassword([]byte(anggota.PINPortal), []byte(pin))
	if err != nil {
		return nil, errors.New("PIN tidak valid")
	}

	// Generate JWT token untuk anggota
	token, err := s.jwtUtil.GenerateTokenAnggota(&anggota)
	if err != nil {
		return nil, errors.New("gagal membuat token autentikasi")
	}

	// Buat response
	response := &LoginAnggotaResponse{
		Token:   token,
		Anggota: anggota.ToResponse(),
	}

	return response, nil
}

// GetInfoAnggota mengambil informasi anggota berdasarkan ID
func (s *PortalAnggotaService) GetInfoAnggota(idKoperasi, idAnggota uuid.UUID) (*models.AnggotaResponse, error) {
	var anggota models.Anggota
	err := s.db.Where("id_koperasi = ? AND id = ?", idKoperasi, idAnggota).
		First(&anggota).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("anggota tidak ditemukan")
		}
		return nil, err
	}

	response := anggota.ToResponse()
	return &response, nil
}

// GetSaldoAnggota mengambil saldo simpanan anggota
func (s *PortalAnggotaService) GetSaldoAnggota(idKoperasi, idAnggota uuid.UUID) (*models.SaldoSimpananAnggota, error) {
	var saldo models.SaldoSimpananAnggota

	// Query untuk menghitung total simpanan per jenis
	err := s.db.Table("simpanan").
		Select(`
			? as id_anggota,
			COALESCE(SUM(CASE WHEN tipe_simpanan = 'pokok' THEN jumlah_setoran ELSE 0 END), 0) as simpanan_pokok,
			COALESCE(SUM(CASE WHEN tipe_simpanan = 'wajib' THEN jumlah_setoran ELSE 0 END), 0) as simpanan_wajib,
			COALESCE(SUM(CASE WHEN tipe_simpanan = 'sukarela' THEN jumlah_setoran ELSE 0 END), 0) as simpanan_sukarela,
			COALESCE(SUM(jumlah_setoran), 0) as total_simpanan
		`, idAnggota).
		Where("id_koperasi = ? AND id_anggota = ? AND deleted_at IS NULL", idKoperasi, idAnggota).
		Scan(&saldo).Error

	if err != nil {
		return nil, err
	}

	// Ambil info anggota untuk melengkapi response
	var anggota models.Anggota
	err = s.db.Where("id = ?", idAnggota).First(&anggota).Error
	if err == nil {
		saldo.NomorAnggota = anggota.NomorAnggota
		saldo.NamaAnggota = anggota.NamaLengkap
	}

	saldo.IDAnggota = idAnggota

	return &saldo, nil
}

// RiwayatTransaksiAnggota adalah struktur untuk riwayat transaksi anggota
type RiwayatTransaksiAnggota struct {
	ID               uuid.UUID           `json:"id"`
	TanggalTransaksi string              `json:"tanggalTransaksi"`
	TipeSimpanan     models.TipeSimpanan `json:"tipeSimpanan"`
	Jumlah           float64             `json:"jumlah"`
	Keterangan       string              `json:"keterangan"`
	NomorReferensi   string              `json:"nomorReferensi"`
}

// GetRiwayatTransaksi mengambil riwayat transaksi simpanan anggota
func (s *PortalAnggotaService) GetRiwayatTransaksi(idKoperasi, idAnggota uuid.UUID, limit, offset int) ([]RiwayatTransaksiAnggota, int64, error) {
	var riwayat []RiwayatTransaksiAnggota
	var total int64

	// Query count
	err := s.db.Table("simpanan").
		Where("id_koperasi = ? AND id_anggota = ? AND deleted_at IS NULL", idKoperasi, idAnggota).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	// Query data dengan pagination
	err = s.db.Table("simpanan").
		Select(`
			id,
			TO_CHAR(tanggal_transaksi, 'YYYY-MM-DD') as tanggal_transaksi,
			tipe_simpanan,
			jumlah_setoran as jumlah,
			keterangan,
			nomor_referensi
		`).
		Where("id_koperasi = ? AND id_anggota = ? AND deleted_at IS NULL", idKoperasi, idAnggota).
		Order("tanggal_transaksi DESC, tanggal_dibuat DESC").
		Limit(limit).
		Offset(offset).
		Scan(&riwayat).Error

	if err != nil {
		return nil, 0, err
	}

	return riwayat, total, nil
}

// UbahPIN mengubah PIN portal anggota
func (s *PortalAnggotaService) UbahPIN(idKoperasi, idAnggota uuid.UUID, pinLama, pinBaru string) error {
	// Ambil anggota
	var anggota models.Anggota
	err := s.db.Where("id_koperasi = ? AND id = ?", idKoperasi, idAnggota).
		First(&anggota).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("anggota tidak ditemukan")
		}
		return err
	}

	// Verifikasi PIN lama
	if anggota.PINPortal != "" {
		err = bcrypt.CompareHashAndPassword([]byte(anggota.PINPortal), []byte(pinLama))
		if err != nil {
			return errors.New("PIN lama tidak sesuai")
		}
	}

	// Validasi PIN baru
	if len(pinBaru) != 6 {
		return errors.New("PIN baru harus 6 digit")
	}

	// Hash PIN baru
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(pinBaru), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi PIN baru")
	}

	// Update PIN
	err = s.db.Model(&anggota).Update("pin_portal", string(hashedPIN)).Error
	if err != nil {
		return errors.New("gagal menyimpan PIN baru")
	}

	return nil
}
