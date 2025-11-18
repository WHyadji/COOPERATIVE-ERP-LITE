package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Koperasi merepresentasikan entitas koperasi dalam sistem
// Setiap koperasi adalah tenant terpisah dalam multi-tenant architecture
type Koperasi struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	NamaKoperasi      string         `gorm:"type:varchar(255);not null" json:"namaKoperasi" validate:"required"`
	Alamat            string         `gorm:"type:text" json:"alamat"`
	NoTelepon         string         `gorm:"type:varchar(20)" json:"noTelepon"`
	Email             string         `gorm:"type:varchar(100)" json:"email" validate:"omitempty,email"`
	LogoURL           string         `gorm:"type:varchar(500)" json:"logoUrl"`
	TahunBukuMulai    int            `gorm:"type:int;default:1" json:"tahunBukuMulai" validate:"min=1,max=12"` // Bulan mulai tahun buku (1-12)
	Pengaturan        string         `gorm:"type:jsonb" json:"pengaturan"`                                     // JSON untuk pengaturan koperasi
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relasi
	Pengguna     []Pengguna     `gorm:"foreignKey:IDKoperasi" json:"-"`
	Anggota      []Anggota      `gorm:"foreignKey:IDKoperasi" json:"-"`
	Akun         []Akun         `gorm:"foreignKey:IDKoperasi" json:"-"`
	Produk       []Produk       `gorm:"foreignKey:IDKoperasi" json:"-"`
	Simpanan     []Simpanan     `gorm:"foreignKey:IDKoperasi" json:"-"`
	Transaksi    []Transaksi    `gorm:"foreignKey:IDKoperasi" json:"-"`
	Penjualan    []Penjualan    `gorm:"foreignKey:IDKoperasi" json:"-"`
}

// BeforeCreate hook untuk generate UUID sebelum create
func (k *Koperasi) BeforeCreate(tx *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	// Set default tahun buku mulai Januari jika belum diset
	if k.TahunBukuMulai == 0 {
		k.TahunBukuMulai = 1
	}
	// Set default Pengaturan ke empty JSON object jika kosong
	// PostgreSQL jsonb tidak menerima empty string
	if k.Pengaturan == "" {
		k.Pengaturan = "{}"
	}
	return nil
}

// TableName menentukan nama tabel di database
func (Koperasi) TableName() string {
	return "koperasi"
}
