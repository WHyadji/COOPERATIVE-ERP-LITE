package services

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/pkg/validasi"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AnggotaService menangani logika bisnis anggota koperasi
type AnggotaService struct {
	db *gorm.DB
}

// NewAnggotaService membuat instance baru AnggotaService
func NewAnggotaService(db *gorm.DB) *AnggotaService {
	return &AnggotaService{db: db}
}

// BuatAnggotaRequest adalah struktur request untuk membuat anggota
type BuatAnggotaRequest struct {
	NamaLengkap      string     `json:"namaLengkap" binding:"required"`
	NIK              string     `json:"nik"`
	TanggalLahir     *time.Time `json:"tanggalLahir"`
	TempatLahir      string     `json:"tempatLahir"`
	JenisKelamin     string     `json:"jenisKelamin"`
	Alamat           string     `json:"alamat"`
	RT               string     `json:"rt"`
	RW               string     `json:"rw"`
	Kelurahan        string     `json:"kelurahan"`
	Kecamatan        string     `json:"kecamatan"`
	KotaKabupaten    string     `json:"kotaKabupaten"`
	Provinsi         string     `json:"provinsi"`
	KodePos          string     `json:"kodePos"`
	NoTelepon        string     `json:"noTelepon"`
	Email            string     `json:"email"`
	Pekerjaan        string     `json:"pekerjaan"`
	TanggalBergabung time.Time  `json:"tanggalBergabung"`
}

// BuatAnggota membuat anggota baru dengan auto-generate nomor anggota
func (s *AnggotaService) BuatAnggota(idKoperasi uuid.UUID, req *BuatAnggotaRequest) (*models.AnggotaResponse, error) {
	// Initialize validator
	validator := validasi.Baru()

	// Validasi business logic
	if err := validator.TeksWajib(req.NamaLengkap, "nama lengkap", 3, 255); err != nil {
		return nil, err
	}

	if err := validator.Email(req.Email); err != nil {
		return nil, err
	}

	if err := validator.NomorHP(req.NoTelepon); err != nil {
		return nil, err
	}

	if err := validator.JenisKelamin(req.JenisKelamin); err != nil {
		return nil, err
	}

	// Validasi tanggal lahir jika ada
	if req.TanggalLahir != nil {
		if err := validator.TanggalLahir(*req.TanggalLahir); err != nil {
			return nil, err
		}
	}

	// Validasi field opsional
	if err := validator.TeksOpsional(req.NIK, "NIK", 16); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.TempatLahir, "tempat lahir", 100); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.Alamat, "alamat", 500); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.Pekerjaan, "pekerjaan", 100); err != nil {
		return nil, err
	}

	// Generate nomor anggota otomatis
	nomorAnggota, err := s.GenerateNomorAnggota(idKoperasi)
	if err != nil {
		return nil, err
	}

	// Set tanggal bergabung ke hari ini jika tidak diisi
	tanggalBergabung := req.TanggalBergabung
	if tanggalBergabung.IsZero() {
		tanggalBergabung = time.Now()
	}

	// Buat anggota baru
	anggota := &models.Anggota{
		IDKoperasi:       idKoperasi,
		NomorAnggota:     nomorAnggota,
		NamaLengkap:      req.NamaLengkap,
		NIK:              req.NIK,
		TanggalLahir:     req.TanggalLahir,
		TempatLahir:      req.TempatLahir,
		JenisKelamin:     req.JenisKelamin,
		Alamat:           req.Alamat,
		RT:               req.RT,
		RW:               req.RW,
		Kelurahan:        req.Kelurahan,
		Kecamatan:        req.Kecamatan,
		KotaKabupaten:    req.KotaKabupaten,
		Provinsi:         req.Provinsi,
		KodePos:          req.KodePos,
		NoTelepon:        req.NoTelepon,
		Email:            req.Email,
		Pekerjaan:        req.Pekerjaan,
		TanggalBergabung: tanggalBergabung,
		Status:           models.StatusAktif,
	}

	// Simpan ke database
	err = s.db.Create(anggota).Error
	if err != nil {
		return nil, errors.New("gagal membuat anggota")
	}

	response := anggota.ToResponse()
	return &response, nil
}

// GenerateNomorAnggota menghasilkan nomor anggota otomatis
// Format: KOOP-YYYY-NNNN (contoh: KOOP-2025-0001)
// Uses row-level locking to prevent race conditions in concurrent requests
func (s *AnggotaService) GenerateNomorAnggota(idKoperasi uuid.UUID) (string, error) {
	tahun := time.Now().Year()
	var nomorAnggota string

	// Use transaction with row-level locking to prevent race conditions
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Lock and get the last member number for this year
		var lastAnggota models.Anggota
		err := tx.Where("id_koperasi = ? AND EXTRACT(YEAR FROM tanggal_bergabung) = ?", idKoperasi, tahun).
			Order("nomor_anggota DESC").
			Limit(1).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastAnggota).Error

		nomorUrut := 1

		// If there's a previous member, parse and increment
		if err == nil && lastAnggota.NomorAnggota != "" {
			// Extract number from KOOP-2025-0001
			var parsedTahun, parsedUrut int
			_, scanErr := fmt.Sscanf(lastAnggota.NomorAnggota, "KOOP-%d-%04d", &parsedTahun, &parsedUrut)
			if scanErr == nil && parsedTahun == tahun {
				nomorUrut = parsedUrut + 1
			}
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		nomorAnggota = fmt.Sprintf("KOOP-%d-%04d", tahun, nomorUrut)
		return nil
	})

	if err != nil {
		return "", errors.New("gagal generate nomor anggota")
	}

	return nomorAnggota, nil
}

// DapatkanSemuaAnggota mengambil daftar anggota dengan pagination dan filter
func (s *AnggotaService) DapatkanSemuaAnggota(idKoperasi uuid.UUID, status string, search string, page, pageSize int) ([]models.AnggotaResponse, int64, error) {
	var anggotaList []models.Anggota
	var total int64

	query := s.db.Model(&models.Anggota{}).Where("id_koperasi = ?", idKoperasi)

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Search by nama atau nomor anggota
	if search != "" {
		query = query.Where("nama_lengkap ILIKE ? OR nomor_anggota ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("tanggal_bergabung DESC").Find(&anggotaList).Error

	if err != nil {
		return nil, 0, errors.New("gagal mengambil daftar anggota")
	}

	// Convert to response
	responses := make([]models.AnggotaResponse, len(anggotaList))
	for i, anggota := range anggotaList {
		responses[i] = anggota.ToResponse()
	}

	return responses, total, nil
}

// DapatkanAnggota mengambil data anggota berdasarkan ID
func (s *AnggotaService) DapatkanAnggota(id uuid.UUID) (*models.AnggotaResponse, error) {
	var anggota models.Anggota
	err := s.db.Where("id = ?", id).First(&anggota).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("anggota tidak ditemukan")
		}
		return nil, err
	}

	response := anggota.ToResponse()
	return &response, nil
}

// DapatkanAnggotaByNomor mengambil anggota berdasarkan nomor anggota
func (s *AnggotaService) DapatkanAnggotaByNomor(idKoperasi uuid.UUID, nomorAnggota string) (*models.AnggotaResponse, error) {
	var anggota models.Anggota
	err := s.db.Where("id_koperasi = ? AND nomor_anggota = ?", idKoperasi, nomorAnggota).First(&anggota).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("anggota tidak ditemukan")
		}
		return nil, err
	}

	response := anggota.ToResponse()
	return &response, nil
}

// PerbaruiAnggotaRequest adalah struktur request untuk update anggota
type PerbaruiAnggotaRequest struct {
	NamaLengkap   string               `json:"namaLengkap"`
	NIK           string               `json:"nik"`
	TanggalLahir  *time.Time           `json:"tanggalLahir"`
	TempatLahir   string               `json:"tempatLahir"`
	JenisKelamin  string               `json:"jenisKelamin"`
	Alamat        string               `json:"alamat"`
	RT            string               `json:"rt"`
	RW            string               `json:"rw"`
	Kelurahan     string               `json:"kelurahan"`
	Kecamatan     string               `json:"kecamatan"`
	KotaKabupaten string               `json:"kotaKabupaten"`
	Provinsi      string               `json:"provinsi"`
	KodePos       string               `json:"kodePos"`
	NoTelepon     string               `json:"noTelepon"`
	Email         string               `json:"email"`
	Pekerjaan     string               `json:"pekerjaan"`
	Status        models.StatusAnggota `json:"status"`
	FotoURL       string               `json:"fotoUrl"`
	Catatan       string               `json:"catatan"`
}

// PerbaruiAnggota mengupdate data anggota
func (s *AnggotaService) PerbaruiAnggota(idKoperasi, id uuid.UUID, req *PerbaruiAnggotaRequest) (*models.AnggotaResponse, error) {
	// Initialize validator
	validator := validasi.Baru()

	// Validasi business logic untuk field yang akan diupdate
	if req.NamaLengkap != "" {
		if err := validator.TeksWajib(req.NamaLengkap, "nama lengkap", 3, 255); err != nil {
			return nil, err
		}
	}

	if req.Email != "" {
		if err := validator.Email(req.Email); err != nil {
			return nil, err
		}
	}

	if req.NoTelepon != "" {
		if err := validator.NomorHP(req.NoTelepon); err != nil {
			return nil, err
		}
	}

	if req.JenisKelamin != "" {
		if err := validator.JenisKelamin(req.JenisKelamin); err != nil {
			return nil, err
		}
	}

	if req.TanggalLahir != nil {
		if err := validator.TanggalLahir(*req.TanggalLahir); err != nil {
			return nil, err
		}
	}

	// Validasi field opsional
	if err := validator.TeksOpsional(req.NIK, "NIK", 16); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.TempatLahir, "tempat lahir", 100); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.Alamat, "alamat", 500); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.Pekerjaan, "pekerjaan", 100); err != nil {
		return nil, err
	}

	if err := validator.TeksOpsional(req.Catatan, "catatan", 1000); err != nil {
		return nil, err
	}

	// Cek apakah anggota ada DAN milik koperasi yang benar (multi-tenant validation)
	var anggota models.Anggota
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&anggota).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("anggota tidak ditemukan atau tidak memiliki akses")
		}
		return nil, err
	}

	// Update fields (hanya yang tidak kosong)
	if req.NamaLengkap != "" {
		anggota.NamaLengkap = req.NamaLengkap
	}
	if req.NIK != "" {
		anggota.NIK = req.NIK
	}
	if req.TanggalLahir != nil {
		anggota.TanggalLahir = req.TanggalLahir
	}
	if req.TempatLahir != "" {
		anggota.TempatLahir = req.TempatLahir
	}
	if req.JenisKelamin != "" {
		anggota.JenisKelamin = req.JenisKelamin
	}
	if req.Alamat != "" {
		anggota.Alamat = req.Alamat
	}
	if req.RT != "" {
		anggota.RT = req.RT
	}
	if req.RW != "" {
		anggota.RW = req.RW
	}
	if req.Kelurahan != "" {
		anggota.Kelurahan = req.Kelurahan
	}
	if req.Kecamatan != "" {
		anggota.Kecamatan = req.Kecamatan
	}
	if req.KotaKabupaten != "" {
		anggota.KotaKabupaten = req.KotaKabupaten
	}
	if req.Provinsi != "" {
		anggota.Provinsi = req.Provinsi
	}
	if req.KodePos != "" {
		anggota.KodePos = req.KodePos
	}
	if req.NoTelepon != "" {
		anggota.NoTelepon = req.NoTelepon
	}
	if req.Email != "" {
		anggota.Email = req.Email
	}
	if req.Pekerjaan != "" {
		anggota.Pekerjaan = req.Pekerjaan
	}
	if req.Status != "" {
		anggota.Status = req.Status
	}
	if req.FotoURL != "" {
		anggota.FotoURL = req.FotoURL
	}
	if req.Catatan != "" {
		anggota.Catatan = req.Catatan
	}

	// Simpan perubahan
	err = s.db.Save(&anggota).Error
	if err != nil {
		return nil, errors.New("gagal memperbarui anggota")
	}

	response := anggota.ToResponse()
	return &response, nil
}

// HapusAnggota menghapus anggota (soft delete)
func (s *AnggotaService) HapusAnggota(idKoperasi, id uuid.UUID) error {
	// Cek apakah anggota ada DAN milik koperasi yang benar (multi-tenant validation)
	var anggota models.Anggota
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&anggota).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("anggota tidak ditemukan atau tidak memiliki akses")
		}
		return err
	}

	// Validasi: tidak boleh menghapus anggota yang memiliki transaksi simpanan
	var jumlahSimpanan int64
	s.db.Model(&models.Simpanan{}).Where("id_anggota = ?", id).Count(&jumlahSimpanan)
	if jumlahSimpanan > 0 {
		return errors.New("tidak dapat menghapus anggota yang memiliki transaksi simpanan")
	}

	// Validasi: tidak boleh menghapus anggota yang memiliki transaksi penjualan
	var jumlahPenjualan int64
	s.db.Model(&models.Penjualan{}).Where("id_anggota = ?", id).Count(&jumlahPenjualan)
	if jumlahPenjualan > 0 {
		return errors.New("tidak dapat menghapus anggota yang memiliki transaksi penjualan")
	}

	// Soft delete
	err = s.db.Delete(&anggota).Error
	if err != nil {
		return errors.New("gagal menghapus anggota")
	}

	return nil
}

// SetPINPortal mengatur PIN untuk login portal anggota
func (s *AnggotaService) SetPINPortal(idKoperasi, id uuid.UUID, pin string) error {
	// Validasi PIN (4-6 digit)
	if len(pin) < 4 || len(pin) > 6 {
		return errors.New("PIN harus 4-6 digit")
	}

	// Cek apakah anggota ada DAN milik koperasi yang benar (multi-tenant validation)
	var anggota models.Anggota
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&anggota).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("anggota tidak ditemukan atau tidak memiliki akses")
		}
		return err
	}

	// Hash PIN menggunakan bcrypt
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi PIN")
	}

	anggota.PINPortal = string(hashedPIN)

	// Simpan
	err = s.db.Save(&anggota).Error
	if err != nil {
		return errors.New("gagal menyimpan PIN")
	}

	return nil
}

// ValidasiPINPortal memvalidasi PIN untuk login portal anggota
func (s *AnggotaService) ValidasiPINPortal(idKoperasi uuid.UUID, nomorAnggota, pin string) (*models.AnggotaResponse, error) {
	// Cari anggota berdasarkan nomor anggota
	var anggota models.Anggota
	err := s.db.Where("id_koperasi = ? AND nomor_anggota = ? AND status = ?", idKoperasi, nomorAnggota, models.StatusAktif).
		First(&anggota).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("nomor anggota atau PIN salah")
		}
		return nil, err
	}

	// Cek apakah PIN sudah di-set
	if anggota.PINPortal == "" {
		return nil, errors.New("PIN belum di-set, silakan hubungi admin")
	}

	// Verifikasi PIN
	err = bcrypt.CompareHashAndPassword([]byte(anggota.PINPortal), []byte(pin))
	if err != nil {
		return nil, errors.New("nomor anggota atau PIN salah")
	}

	response := anggota.ToResponse()
	return &response, nil
}

// HitungJumlahAnggota menghitung jumlah anggota berdasarkan status
func (s *AnggotaService) HitungJumlahAnggota(idKoperasi uuid.UUID, status string) (int64, error) {
	var count int64
	query := s.db.Model(&models.Anggota{}).Where("id_koperasi = ?", idKoperasi)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, errors.New("gagal menghitung jumlah anggota")
	}

	return count, nil
}

// GetSemuaAnggota is an English wrapper for DapatkanSemuaAnggota
func (s *AnggotaService) GetSemuaAnggota(idKoperasi uuid.UUID, status *models.StatusAnggota, search string, page, pageSize int) ([]models.AnggotaResponse, int64, error) {
	statusStr := ""
	if status != nil {
		statusStr = string(*status)
	}
	return s.DapatkanSemuaAnggota(idKoperasi, statusStr, search, page, pageSize)
}

// GetAnggotaByID is an English wrapper for DapatkanAnggota with multi-tenant isolation
func (s *AnggotaService) GetAnggotaByID(idKoperasi, id uuid.UUID) (*models.AnggotaResponse, error) {
	var anggota models.Anggota
	err := s.db.Where("id = ? AND id_koperasi = ?", id, idKoperasi).First(&anggota).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("anggota tidak ditemukan")
		}
		return nil, err
	}

	response := anggota.ToResponse()
	return &response, nil
}

// GetAnggotaByNomor is an English wrapper for DapatkanAnggotaByNomor
func (s *AnggotaService) GetAnggotaByNomor(idKoperasi uuid.UUID, nomorAnggota string) (*models.AnggotaResponse, error) {
	return s.DapatkanAnggotaByNomor(idKoperasi, nomorAnggota)
}
