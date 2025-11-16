package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MetodePembayaran mendefinisikan metode pembayaran
type MetodePembayaran string

const (
	PembayaranTunai MetodePembayaran = "tunai" // Cash only untuk MVP
)

// Penjualan merepresentasikan transaksi penjualan di POS
type Penjualan struct {
	ID                uuid.UUID        `gorm:"type:uuid;primary_key" json:"id"`
	IDKoperasi        uuid.UUID        `gorm:"type:uuid;not null;index" json:"idKoperasi" validate:"required"`
	NomorPenjualan    string           `gorm:"type:varchar(50);not null;uniqueIndex:idx_koperasi_nomor_penjualan" json:"nomorPenjualan" validate:"required"`
	TanggalPenjualan  time.Time        `gorm:"type:timestamp;not null;index" json:"tanggalPenjualan" validate:"required"`
	IDAnggota         *uuid.UUID       `gorm:"type:uuid;index" json:"idAnggota"` // Opsional, bisa non-member
	TotalBelanja      float64          `gorm:"type:decimal(15,2);not null" json:"totalBelanja" validate:"required,gt=0"`
	MetodePembayaran  MetodePembayaran `gorm:"type:varchar(20);not null;default:'tunai'" json:"metodePembayaran"`
	JumlahBayar       float64          `gorm:"type:decimal(15,2);not null" json:"jumlahBayar" validate:"required,gte=0"`
	Kembalian         float64          `gorm:"type:decimal(15,2);not null;default:0" json:"kembalian"`
	IDKasir           uuid.UUID        `gorm:"type:uuid;not null" json:"idKasir" validate:"required"`
	IDTransaksi       *uuid.UUID       `gorm:"type:uuid;index" json:"idTransaksi"` // Link ke jurnal akuntansi
	Catatan           string           `gorm:"type:text" json:"catatan"`
	TanggalDibuat     time.Time        `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time        `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt   `gorm:"index" json:"-"`

	// Relasi
	Koperasi      Koperasi        `gorm:"foreignKey:IDKoperasi;constraint:OnDelete:CASCADE" json:"-"`
	Anggota       *Anggota        `gorm:"foreignKey:IDAnggota" json:"anggota,omitempty"`
	Kasir         Pengguna        `gorm:"foreignKey:IDKasir" json:"kasir,omitempty"`
	Transaksi     *Transaksi      `gorm:"foreignKey:IDTransaksi" json:"-"`
	ItemPenjualan []ItemPenjualan `gorm:"foreignKey:IDPenjualan;constraint:OnDelete:CASCADE" json:"itemPenjualan,omitempty"`
}

// BeforeCreate hook untuk generate UUID dan nomor penjualan
func (p *Penjualan) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Set tanggal penjualan ke sekarang jika belum diset
	if p.TanggalPenjualan.IsZero() {
		p.TanggalPenjualan = time.Now()
	}

	// Set metode pembayaran default
	if p.MetodePembayaran == "" {
		p.MetodePembayaran = PembayaranTunai
	}

	// Hitung kembalian
	p.Kembalian = p.JumlahBayar - p.TotalBelanja
	if p.Kembalian < 0 {
		p.Kembalian = 0
	}

	return nil
}

// TableName menentukan nama tabel di database
func (Penjualan) TableName() string {
	return "penjualan"
}

// ItemPenjualan merepresentasikan item/produk dalam penjualan
type ItemPenjualan struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDPenjualan  uuid.UUID      `gorm:"type:uuid;not null;index" json:"idPenjualan" validate:"required"`
	IDProduk     uuid.UUID      `gorm:"type:uuid;not null;index" json:"idProduk" validate:"required"`
	NamaProduk   string         `gorm:"type:varchar(255);not null" json:"namaProduk"` // Snapshot nama produk saat transaksi
	Kuantitas    int            `gorm:"type:int;not null" json:"kuantitas" validate:"required,gt=0"`
	HargaSatuan  float64        `gorm:"type:decimal(15,2);not null" json:"hargaSatuan" validate:"required,gt=0"`
	Subtotal     float64        `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Penjualan Penjualan `gorm:"foreignKey:IDPenjualan;constraint:OnDelete:CASCADE" json:"-"`
	Produk    Produk    `gorm:"foreignKey:IDProduk;constraint:OnDelete:RESTRICT" json:"produk,omitempty"`
}

// BeforeCreate hook untuk generate UUID dan hitung subtotal
func (i *ItemPenjualan) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}

	// Hitung subtotal
	i.Subtotal = float64(i.Kuantitas) * i.HargaSatuan

	return nil
}

// BeforeSave hook untuk hitung subtotal
func (i *ItemPenjualan) BeforeSave(tx *gorm.DB) error {
	i.Subtotal = float64(i.Kuantitas) * i.HargaSatuan
	return nil
}

// TableName menentukan nama tabel di database
func (ItemPenjualan) TableName() string {
	return "item_penjualan"
}

// PenjualanResponse adalah response untuk API
type PenjualanResponse struct {
	ID               uuid.UUID             `json:"id"`
	NomorPenjualan   string                `json:"nomorPenjualan"`
	TanggalPenjualan time.Time             `json:"tanggalPenjualan"`
	IDAnggota        *uuid.UUID            `json:"idAnggota"`
	NamaAnggota      string                `json:"namaAnggota,omitempty"`
	NomorAnggota     string                `json:"nomorAnggota,omitempty"`
	TotalBelanja     float64               `json:"totalBelanja"`
	MetodePembayaran MetodePembayaran      `json:"metodePembayaran"`
	JumlahBayar      float64               `json:"jumlahBayar"`
	Kembalian        float64               `json:"kembalian"`
	NamaKasir        string                `json:"namaKasir"`
	Catatan          string                `json:"catatan"`
	ItemPenjualan    []ItemPenjualanResponse `json:"itemPenjualan,omitempty"`
}

// ItemPenjualanResponse adalah response untuk item penjualan
type ItemPenjualanResponse struct {
	ID          uuid.UUID `json:"id"`
	IDProduk    uuid.UUID `json:"idProduk"`
	KodeProduk  string    `json:"kodeProduk,omitempty"`
	NamaProduk  string    `json:"namaProduk"`
	Kuantitas   int       `json:"kuantitas"`
	HargaSatuan float64   `json:"hargaSatuan"`
	Subtotal    float64   `json:"subtotal"`
}

// ToResponse mengkonversi Penjualan ke PenjualanResponse
func (p *Penjualan) ToResponse() PenjualanResponse {
	resp := PenjualanResponse{
		ID:               p.ID,
		NomorPenjualan:   p.NomorPenjualan,
		TanggalPenjualan: p.TanggalPenjualan,
		IDAnggota:        p.IDAnggota,
		TotalBelanja:     p.TotalBelanja,
		MetodePembayaran: p.MetodePembayaran,
		JumlahBayar:      p.JumlahBayar,
		Kembalian:        p.Kembalian,
		Catatan:          p.Catatan,
	}

	// Populate info anggota jika ada dan relasi sudah di-load
	if p.Anggota != nil && p.Anggota.ID != uuid.Nil {
		resp.NamaAnggota = p.Anggota.NamaLengkap
		resp.NomorAnggota = p.Anggota.NomorAnggota
	}

	// Populate info kasir jika relasi sudah di-load
	if p.Kasir.ID != uuid.Nil {
		resp.NamaKasir = p.Kasir.NamaLengkap
	}

	// Convert item penjualan jika ada
	if len(p.ItemPenjualan) > 0 {
		resp.ItemPenjualan = make([]ItemPenjualanResponse, len(p.ItemPenjualan))
		for i, item := range p.ItemPenjualan {
			resp.ItemPenjualan[i] = ItemPenjualanResponse{
				ID:          item.ID,
				IDProduk:    item.IDProduk,
				NamaProduk:  item.NamaProduk,
				Kuantitas:   item.Kuantitas,
				HargaSatuan: item.HargaSatuan,
				Subtotal:    item.Subtotal,
			}

			// Populate kode produk jika relasi sudah di-load
			if item.Produk.ID != uuid.Nil {
				resp.ItemPenjualan[i].KodeProduk = item.Produk.KodeProduk
			}
		}
	}

	return resp
}
