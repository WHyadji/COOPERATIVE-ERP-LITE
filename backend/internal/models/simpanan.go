package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TipeSimpanan mendefinisikan jenis simpanan koperasi
type TipeSimpanan string

const (
	SimpananPokok     TipeSimpanan = "POKOK"     // Simpanan pokok - dibayar sekali saat bergabung
	SimpananWajib     TipeSimpanan = "WAJIB"     // Simpanan wajib - dibayar rutin (bulanan)
	SimpananSukarela  TipeSimpanan = "SUKARELA"  // Simpanan sukarela - opsional
)

// Simpanan merepresentasikan transaksi simpanan anggota
type Simpanan struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDKoperasi        uuid.UUID      `gorm:"type:uuid;not null;index" json:"idKoperasi" validate:"required"`
	IDAnggota         uuid.UUID      `gorm:"type:uuid;not null;index" json:"idAnggota" validate:"required"`
	TipeSimpanan      TipeSimpanan   `gorm:"type:varchar(20);not null" json:"tipeSimpanan" validate:"required,oneof=POKOK WAJIB SUKARELA"`
	TanggalTransaksi  time.Time      `gorm:"type:date;not null;index" json:"tanggalTransaksi" validate:"required"`
	JumlahSetoran     float64        `gorm:"type:decimal(15,2);not null" json:"jumlahSetoran" validate:"required,gt=0"`
	Keterangan        string         `gorm:"type:text" json:"keterangan"`
	NomorReferensi    string         `gorm:"type:varchar(50)" json:"nomorReferensi"` // Nomor bukti transaksi
	IDTransaksi       *uuid.UUID     `gorm:"type:uuid;index" json:"idTransaksi"`     // Link ke jurnal akuntansi
	DibuatOleh        uuid.UUID      `gorm:"type:uuid" json:"dibuatOleh"`            // ID pengguna yang membuat transaksi
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Koperasi  Koperasi   `gorm:"foreignKey:IDKoperasi;constraint:OnDelete:CASCADE" json:"-"`
	Anggota   Anggota    `gorm:"foreignKey:IDAnggota;constraint:OnDelete:CASCADE" json:"-"`
	Transaksi *Transaksi `gorm:"foreignKey:IDTransaksi" json:"-"`
}

// BeforeCreate hook untuk generate UUID
func (s *Simpanan) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	// Set tanggal transaksi ke hari ini jika belum diset
	if s.TanggalTransaksi.IsZero() {
		s.TanggalTransaksi = time.Now()
	}

	return nil
}

// TableName menentukan nama tabel di database
func (Simpanan) TableName() string {
	return "simpanan"
}

// SimpananResponse adalah response untuk API
type SimpananResponse struct {
	ID               uuid.UUID    `json:"id"`
	IDAnggota        uuid.UUID    `json:"idAnggota"`
	NamaAnggota      string       `json:"namaAnggota"`
	NomorAnggota     string       `json:"nomorAnggota"`
	TipeSimpanan     TipeSimpanan `json:"tipeSimpanan"`
	TanggalTransaksi time.Time    `json:"tanggalTransaksi"`
	JumlahSetoran    float64      `json:"jumlahSetoran"`
	Keterangan       string       `json:"keterangan"`
	NomorReferensi   string       `json:"nomorReferensi"`
}

// ToResponse mengkonversi Simpanan ke SimpananResponse
func (s *Simpanan) ToResponse() SimpananResponse {
	resp := SimpananResponse{
		ID:               s.ID,
		IDAnggota:        s.IDAnggota,
		TipeSimpanan:     s.TipeSimpanan,
		TanggalTransaksi: s.TanggalTransaksi,
		JumlahSetoran:    s.JumlahSetoran,
		Keterangan:       s.Keterangan,
		NomorReferensi:   s.NomorReferensi,
	}

	// Populate nama anggota jika relasi sudah di-load
	if s.Anggota.ID != uuid.Nil {
		resp.NamaAnggota = s.Anggota.NamaLengkap
		resp.NomorAnggota = s.Anggota.NomorAnggota
	}

	return resp
}

// RingkasanSimpanan adalah struktur untuk summary simpanan koperasi
type RingkasanSimpanan struct {
	TotalSimpananPokok    float64 `json:"totalSimpananPokok"`
	TotalSimpananWajib    float64 `json:"totalSimpananWajib"`
	TotalSimpananSukarela float64 `json:"totalSimpananSukarela"`
	TotalSemuaSimpanan    float64 `json:"totalSemuaSimpanan"`
	JumlahAnggota         int64   `json:"jumlahAnggota"`
}

// SaldoSimpananAnggota adalah struktur untuk saldo simpanan per anggota
type SaldoSimpananAnggota struct {
	IDAnggota        uuid.UUID `json:"idAnggota"`
	NomorAnggota     string    `json:"nomorAnggota"`
	NamaAnggota      string    `json:"namaAnggota"`
	SimpananPokok    float64   `json:"simpananPokok"`
	SimpananWajib    float64   `json:"simpananWajib"`
	SimpananSukarela float64   `json:"simpananSukarela"`
	TotalSimpanan    float64   `json:"totalSimpanan"`
}
