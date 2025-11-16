package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TipeAkun mendefinisikan tipe akun dalam Chart of Accounts
type TipeAkun string

const (
	AkunAset      TipeAkun = "aset"      // Aset/Harta
	AkunKewajiban TipeAkun = "kewajiban" // Kewajiban/Hutang
	AkunModal     TipeAkun = "modal"     // Modal/Ekuitas
	AkunPendapatan TipeAkun = "pendapatan" // Pendapatan
	AkunBeban     TipeAkun = "beban"     // Beban/Biaya
)

// Akun merepresentasikan Chart of Accounts (Bagan Akun)
type Akun struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDKoperasi        uuid.UUID      `gorm:"type:uuid;not null;index" json:"idKoperasi" validate:"required"`
	KodeAkun          string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_koperasi_kode_akun" json:"kodeAkun" validate:"required"`
	NamaAkun          string         `gorm:"type:varchar(255);not null" json:"namaAkun" validate:"required"`
	TipeAkun          TipeAkun       `gorm:"type:varchar(20);not null" json:"tipeAkun" validate:"required,oneof=aset kewajiban modal pendapatan beban"`
	IDInduk           *uuid.UUID     `gorm:"type:uuid;index" json:"idInduk"` // Parent account untuk hierarchical COA
	NormalSaldo       string         `gorm:"type:varchar(6);not null" json:"normalSaldo" validate:"oneof=debit kredit"` // debit atau kredit
	Deskripsi         string         `gorm:"type:text" json:"deskripsi"`
	StatusAktif       bool           `gorm:"type:boolean;default:true" json:"statusAktif"`
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Koperasi        Koperasi          `gorm:"foreignKey:IDKoperasi;constraint:OnDelete:CASCADE" json:"-"`
	AkunInduk       *Akun             `gorm:"foreignKey:IDInduk" json:"-"`
	SubAkun         []Akun            `gorm:"foreignKey:IDInduk" json:"-"`
	BarisTransaksi  []BarisTransaksi  `gorm:"foreignKey:IDAkun" json:"-"`
}

// BeforeCreate hook untuk generate UUID dan set normal saldo
func (a *Akun) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	// Set normal saldo berdasarkan tipe akun jika belum diset
	if a.NormalSaldo == "" {
		switch a.TipeAkun {
		case AkunAset, AkunBeban:
			a.NormalSaldo = "debit"
		case AkunKewajiban, AkunModal, AkunPendapatan:
			a.NormalSaldo = "kredit"
		}
	}

	return nil
}

// TableName menentukan nama tabel di database
func (Akun) TableName() string {
	return "akun"
}

// AkunResponse adalah response untuk API
type AkunResponse struct {
	ID          uuid.UUID `json:"id"`
	KodeAkun    string    `json:"kodeAkun"`
	NamaAkun    string    `json:"namaAkun"`
	TipeAkun    TipeAkun  `json:"tipeAkun"`
	IDInduk     *uuid.UUID `json:"idInduk"`
	NamaInduk   string    `json:"namaInduk,omitempty"`
	NormalSaldo string    `json:"normalSaldo"`
	Deskripsi   string    `json:"deskripsi"`
	StatusAktif bool      `json:"statusAktif"`
	Saldo       float64   `json:"saldo,omitempty"` // Computed field
}

// ToResponse mengkonversi Akun ke AkunResponse
func (a *Akun) ToResponse() AkunResponse {
	resp := AkunResponse{
		ID:          a.ID,
		KodeAkun:    a.KodeAkun,
		NamaAkun:    a.NamaAkun,
		TipeAkun:    a.TipeAkun,
		IDInduk:     a.IDInduk,
		NormalSaldo: a.NormalSaldo,
		Deskripsi:   a.Deskripsi,
		StatusAktif: a.StatusAktif,
	}

	// Populate nama induk jika relasi sudah di-load
	if a.AkunInduk != nil && a.AkunInduk.ID != uuid.Nil {
		resp.NamaInduk = a.AkunInduk.NamaAkun
	}

	return resp
}
