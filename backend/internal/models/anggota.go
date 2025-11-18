package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StatusAnggota mendefinisikan status keanggotaan
type StatusAnggota string

const (
	StatusAktif     StatusAnggota = "aktif"     // Anggota aktif
	StatusNonAktif  StatusAnggota = "nonaktif"  // Anggota tidak aktif
	StatusDiberhentikan StatusAnggota = "diberhentikan" // Anggota diberhentikan
)

// Anggota merepresentasikan anggota koperasi
type Anggota struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDKoperasi        uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_koperasi_nomor" json:"idKoperasi" validate:"required"`
	NomorAnggota      string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_koperasi_nomor" json:"nomorAnggota" validate:"required"`
	NamaLengkap       string         `gorm:"type:varchar(255);not null" json:"namaLengkap" validate:"required"`
	NIK               string         `gorm:"type:varchar(16)" json:"nik" validate:"omitempty,len=16"` // Nomor Induk Kependudukan
	TanggalLahir      *time.Time     `gorm:"type:date" json:"tanggalLahir"`
	TempatLahir       string         `gorm:"type:varchar(100)" json:"tempatLahir"`
	JenisKelamin      string         `gorm:"type:varchar(10)" json:"jenisKelamin" validate:"omitempty,oneof=L P"` // L=Laki-laki, P=Perempuan
	Alamat            string         `gorm:"type:text" json:"alamat"`
	RT                string         `gorm:"type:varchar(3)" json:"rt"`
	RW                string         `gorm:"type:varchar(3)" json:"rw"`
	Kelurahan         string         `gorm:"type:varchar(100)" json:"kelurahan"`
	Kecamatan         string         `gorm:"type:varchar(100)" json:"kecamatan"`
	KotaKabupaten     string         `gorm:"type:varchar(100)" json:"kotaKabupaten"`
	Provinsi          string         `gorm:"type:varchar(100)" json:"provinsi"`
	KodePos           string         `gorm:"type:varchar(5)" json:"kodePos"`
	NoTelepon         string         `gorm:"type:varchar(20)" json:"noTelepon"`
	Email             string         `gorm:"type:varchar(100)" json:"email" validate:"omitempty,email"`
	Pekerjaan         string         `gorm:"type:varchar(100)" json:"pekerjaan"`
	TanggalBergabung  time.Time      `gorm:"type:date;not null" json:"tanggalBergabung" validate:"required"`
	Status            StatusAnggota  `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	PINPortal         string         `gorm:"type:varchar(255)" json:"-"` // Hash PIN untuk login portal anggota
	FotoURL           string         `gorm:"type:varchar(500)" json:"fotoUrl"`
	Catatan           string         `gorm:"type:text" json:"catatan"`
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Koperasi Koperasi   `gorm:"foreignKey:IDKoperasi;constraint:OnDelete:CASCADE" json:"-"`
	Simpanan []Simpanan `gorm:"foreignKey:IDAnggota" json:"-"`
}

// BeforeCreate hook untuk generate UUID
func (a *Anggota) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	// Set status default jika belum diset
	if a.Status == "" {
		a.Status = StatusAktif
	}

	// Set tanggal bergabung ke hari ini jika belum diset
	if a.TanggalBergabung.IsZero() {
		a.TanggalBergabung = time.Now()
	}

	return nil
}

// TableName menentukan nama tabel di database
func (Anggota) TableName() string {
	return "anggota"
}

// AnggotaResponse adalah response untuk API
type AnggotaResponse struct {
	ID               uuid.UUID     `json:"id"`
	IDKoperasi       uuid.UUID     `json:"idKoperasi"`
	NomorAnggota     string        `json:"nomorAnggota"`
	NamaLengkap      string        `json:"namaLengkap"`
	NIK              string        `json:"nik"`
	TanggalLahir     *time.Time    `json:"tanggalLahir"`
	TempatLahir      string        `json:"tempatLahir"`
	JenisKelamin     string        `json:"jenisKelamin"`
	Alamat           string        `json:"alamat"`
	NoTelepon        string        `json:"noTelepon"`
	Email            string        `json:"email"`
	Pekerjaan        string        `json:"pekerjaan"`
	TanggalBergabung time.Time     `json:"tanggalBergabung"`
	Status           StatusAnggota `json:"status"`
	FotoURL          string        `json:"fotoUrl"`
}

// ToResponse mengkonversi Anggota ke AnggotaResponse
func (a *Anggota) ToResponse() AnggotaResponse {
	return AnggotaResponse{
		ID:               a.ID,
		IDKoperasi:       a.IDKoperasi,
		NomorAnggota:     a.NomorAnggota,
		NamaLengkap:      a.NamaLengkap,
		NIK:              a.NIK,
		TanggalLahir:     a.TanggalLahir,
		TempatLahir:      a.TempatLahir,
		JenisKelamin:     a.JenisKelamin,
		Alamat:           a.Alamat,
		NoTelepon:        a.NoTelepon,
		Email:            a.Email,
		Pekerjaan:        a.Pekerjaan,
		TanggalBergabung: a.TanggalBergabung,
		Status:           a.Status,
		FotoURL:          a.FotoURL,
	}
}
