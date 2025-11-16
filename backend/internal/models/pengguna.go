package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PeranPengguna mendefinisikan role/peran pengguna dalam sistem
type PeranPengguna string

const (
	PeranAdmin     PeranPengguna = "admin"     // Admin koperasi - akses penuh
	PeranBendahara PeranPengguna = "bendahara" // Bendahara - akses keuangan
	PeranKasir     PeranPengguna = "kasir"     // Kasir - akses POS
	PeranAnggota   PeranPengguna = "anggota"   // Anggota - akses portal (read-only)
)

// Pengguna merepresentasikan user/pengguna sistem
type Pengguna struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDKoperasi        uuid.UUID      `gorm:"type:uuid;not null;index" json:"idKoperasi" validate:"required"`
	NamaLengkap       string         `gorm:"type:varchar(255);not null" json:"namaLengkap" validate:"required"`
	NamaPengguna      string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_koperasi_username" json:"namaPengguna" validate:"required,min=3"`
	Email             string         `gorm:"type:varchar(100);not null" json:"email" validate:"required,email"`
	KataSandiHash     string         `gorm:"type:varchar(255);not null" json:"-"` // Password hash, tidak di-export ke JSON
	Peran             PeranPengguna  `gorm:"type:varchar(20);not null" json:"peran" validate:"required,oneof=admin bendahara kasir anggota"`
	StatusAktif       bool           `gorm:"type:boolean;default:true" json:"statusAktif"`
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Koperasi Koperasi `gorm:"foreignKey:IDKoperasi;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate hook untuk generate UUID dan hash password
func (p *Pengguna) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// SetKataSandi meng-hash password dan menyimpannya
func (p *Pengguna) SetKataSandi(kataSandi string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(kataSandi), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.KataSandiHash = string(hashedPassword)
	return nil
}

// CekKataSandi memverifikasi password yang diberikan dengan hash yang tersimpan
func (p *Pengguna) CekKataSandi(kataSandi string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.KataSandiHash), []byte(kataSandi))
	return err == nil
}

// TableName menentukan nama tabel di database
func (Pengguna) TableName() string {
	return "pengguna"
}

// PenggunaResponse adalah response untuk API (tanpa password hash)
type PenggunaResponse struct {
	ID           uuid.UUID     `json:"id"`
	IDKoperasi   uuid.UUID     `json:"idKoperasi"`
	NamaLengkap  string        `json:"namaLengkap"`
	NamaPengguna string        `json:"namaPengguna"`
	Email        string        `json:"email"`
	Peran        PeranPengguna `json:"peran"`
	StatusAktif  bool          `json:"statusAktif"`
}

// ToResponse mengkonversi Pengguna ke PenggunaResponse
func (p *Pengguna) ToResponse() PenggunaResponse {
	return PenggunaResponse{
		ID:           p.ID,
		IDKoperasi:   p.IDKoperasi,
		NamaLengkap:  p.NamaLengkap,
		NamaPengguna: p.NamaPengguna,
		Email:        p.Email,
		Peran:        p.Peran,
		StatusAktif:  p.StatusAktif,
	}
}
