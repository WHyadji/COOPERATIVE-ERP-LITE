package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Produk merepresentasikan produk yang dijual di koperasi
type Produk struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDKoperasi        uuid.UUID      `gorm:"type:uuid;not null;index" json:"idKoperasi" validate:"required"`
	KodeProduk        string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_koperasi_kode_produk" json:"kodeProduk" validate:"required"`
	NamaProduk        string         `gorm:"type:varchar(255);not null" json:"namaProduk" validate:"required"`
	Kategori          string         `gorm:"type:varchar(100)" json:"kategori"`
	Deskripsi         string         `gorm:"type:text" json:"deskripsi"`
	Harga             float64        `gorm:"type:decimal(15,2);not null" json:"harga" validate:"required,gte=0"`
	HargaBeli         float64        `gorm:"type:decimal(15,2)" json:"hargaBeli" validate:"gte=0"` // Harga beli/HPP
	Stok              int            `gorm:"type:int;default:0" json:"stok"`
	StokMinimum       int            `gorm:"type:int;default:0" json:"stokMinimum"`
	Satuan            string         `gorm:"type:varchar(20);default:'pcs'" json:"satuan"` // pcs, kg, liter, dll
	Barcode           string         `gorm:"type:varchar(100)" json:"barcode"`
	GambarURL         string         `gorm:"type:varchar(500)" json:"gambarUrl"`
	StatusAktif       bool           `gorm:"type:boolean;default:true" json:"statusAktif"`
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Koperasi     Koperasi       `gorm:"foreignKey:IDKoperasi;constraint:OnDelete:CASCADE" json:"-"`
	ItemPenjualan []ItemPenjualan `gorm:"foreignKey:IDProduk" json:"-"`
}

// BeforeCreate hook untuk generate UUID
func (p *Produk) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Set satuan default jika belum diset
	if p.Satuan == "" {
		p.Satuan = "pcs"
	}

	return nil
}

// TableName menentukan nama tabel di database
func (Produk) TableName() string {
	return "produk"
}

// ProdukResponse adalah response untuk API
type ProdukResponse struct {
	ID          uuid.UUID `json:"id"`
	KodeProduk  string    `json:"kodeProduk"`
	NamaProduk  string    `json:"namaProduk"`
	Kategori    string    `json:"kategori"`
	Deskripsi   string    `json:"deskripsi"`
	Harga       float64   `json:"harga"`
	HargaBeli   float64   `json:"hargaBeli"`
	Stok        int       `json:"stok"`
	StokMinimum int       `json:"stokMinimum"`
	Satuan      string    `json:"satuan"`
	Barcode     string    `json:"barcode"`
	GambarURL   string    `json:"gambarUrl"`
	StatusAktif bool      `json:"statusAktif"`
}

// ToResponse mengkonversi Produk ke ProdukResponse
func (p *Produk) ToResponse() ProdukResponse {
	return ProdukResponse{
		ID:          p.ID,
		KodeProduk:  p.KodeProduk,
		NamaProduk:  p.NamaProduk,
		Kategori:    p.Kategori,
		Deskripsi:   p.Deskripsi,
		Harga:       p.Harga,
		HargaBeli:   p.HargaBeli,
		Stok:        p.Stok,
		StokMinimum: p.StokMinimum,
		Satuan:      p.Satuan,
		Barcode:     p.Barcode,
		GambarURL:   p.GambarURL,
		StatusAktif: p.StatusAktif,
	}
}
