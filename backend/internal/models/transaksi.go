package models

import (
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// EpsilonTolerance adalah toleransi untuk perbandingan floating-point dalam rupiah
	// 0.01 = 1 sen toleransi untuk menangani masalah floating-point precision
	EpsilonTolerance = 0.01
)

// Transaksi merepresentasikan jurnal transaksi akuntansi (header)
type Transaksi struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDKoperasi        uuid.UUID      `gorm:"type:uuid;not null;index" json:"idKoperasi" validate:"required"`
	NomorJurnal       string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_koperasi_nomor_jurnal" json:"nomorJurnal" validate:"required"`
	TanggalTransaksi  time.Time      `gorm:"type:date;not null;index" json:"tanggalTransaksi" validate:"required"`
	Deskripsi         string         `gorm:"type:text;not null" json:"deskripsi" validate:"required"`
	NomorReferensi    string         `gorm:"type:varchar(50)" json:"nomorReferensi"` // Nomor bukti transaksi eksternal
	TipeTransaksi     string         `gorm:"type:varchar(50)" json:"tipeTransaksi"`  // manual, penjualan, simpanan, dll
	TotalDebit        float64        `gorm:"type:decimal(15,2);not null;default:0" json:"totalDebit"`
	TotalKredit       float64        `gorm:"type:decimal(15,2);not null;default:0" json:"totalKredit"`
	StatusBalanced    bool           `gorm:"type:boolean;default:false" json:"statusBalanced"` // Apakah debit = kredit
	DibuatOleh        uuid.UUID      `gorm:"type:uuid" json:"dibuatOleh"`
	DiperbaruiOleh    uuid.UUID      `gorm:"type:uuid" json:"diperbaruiOleh"`
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Koperasi       Koperasi         `gorm:"foreignKey:IDKoperasi;constraint:OnDelete:CASCADE" json:"-"`
	BarisTransaksi []BarisTransaksi `gorm:"foreignKey:IDTransaksi;constraint:OnDelete:CASCADE" json:"barisTransaksi,omitempty"`
}

// BeforeCreate hook untuk generate UUID
func (t *Transaksi) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	// Set tanggal transaksi ke hari ini jika belum diset
	if t.TanggalTransaksi.IsZero() {
		t.TanggalTransaksi = time.Now()
	}

	return nil
}

// BeforeSave hook untuk menghitung total dan validasi balance
func (t *Transaksi) BeforeSave(tx *gorm.DB) error {
	// Hitung total debit dan kredit dari baris transaksi
	var totalDebit, totalKredit float64
	for _, baris := range t.BarisTransaksi {
		totalDebit += baris.JumlahDebit
		totalKredit += baris.JumlahKredit
	}

	t.TotalDebit = totalDebit
	t.TotalKredit = totalKredit

	// Cek apakah transaksi balanced (debit = kredit) dengan epsilon tolerance
	// Menggunakan math.Abs untuk mengatasi floating-point precision issues
	isBalanced := math.Abs(totalDebit-totalKredit) <= EpsilonTolerance
	hasValue := totalDebit >= EpsilonTolerance
	t.StatusBalanced = isBalanced && hasValue

	return nil
}

// TableName menentukan nama tabel di database
func (Transaksi) TableName() string {
	return "transaksi"
}

// BarisTransaksi merepresentasikan baris/detail jurnal transaksi
type BarisTransaksi struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IDTransaksi       uuid.UUID      `gorm:"type:uuid;not null;index" json:"idTransaksi" validate:"required"`
	IDAkun            uuid.UUID      `gorm:"type:uuid;not null;index" json:"idAkun" validate:"required"`
	JumlahDebit       float64        `gorm:"type:decimal(15,2);not null;default:0" json:"jumlahDebit" validate:"gte=0"`
	JumlahKredit      float64        `gorm:"type:decimal(15,2);not null;default:0" json:"jumlahKredit" validate:"gte=0"`
	Keterangan        string         `gorm:"type:text" json:"keterangan"`
	TanggalDibuat     time.Time      `gorm:"autoCreateTime" json:"tanggalDibuat"`
	TanggalDiperbarui time.Time      `gorm:"autoUpdateTime" json:"tanggalDiperbarui"`
	TanggalDihapus    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Transaksi Transaksi `gorm:"foreignKey:IDTransaksi;constraint:OnDelete:CASCADE" json:"-"`
	Akun      Akun      `gorm:"foreignKey:IDAkun;constraint:OnDelete:RESTRICT" json:"akun,omitempty"`
}

// BeforeCreate hook untuk generate UUID
func (b *BarisTransaksi) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// TableName menentukan nama tabel di database
func (BarisTransaksi) TableName() string {
	return "baris_transaksi"
}

// TransaksiResponse adalah response untuk API
type TransaksiResponse struct {
	ID                 uuid.UUID                `json:"id"`
	NomorJurnal        string                   `json:"nomorJurnal"`
	TanggalTransaksi   time.Time                `json:"tanggalTransaksi"`
	Deskripsi          string                   `json:"deskripsi"`
	NomorReferensi     string                   `json:"nomorReferensi"`
	TipeTransaksi      string                   `json:"tipeTransaksi"`
	TotalDebit         float64                  `json:"totalDebit"`
	TotalKredit        float64                  `json:"totalKredit"`
	StatusBalanced     bool                     `json:"statusBalanced"`
	DibuatOleh         uuid.UUID                `json:"dibuatOleh,omitempty"`
	NamaDibuatOleh     string                   `json:"namaDibuatOleh,omitempty"`
	DiperbaruiOleh     uuid.UUID                `json:"diperbaruiOleh,omitempty"`
	NamaDiperbaruiOleh string                   `json:"namaDiperbaruiOleh,omitempty"`
	TanggalDibuat      time.Time                `json:"tanggalDibuat,omitempty"`
	TanggalDiperbarui  time.Time                `json:"tanggalDiperbarui,omitempty"`
	BarisTransaksi     []BarisTransaksiResponse `json:"barisTransaksi,omitempty"`
}

// BarisTransaksiResponse adalah response untuk baris transaksi
type BarisTransaksiResponse struct {
	ID           uuid.UUID `json:"id"`
	IDAkun       uuid.UUID `json:"idAkun"`
	KodeAkun     string    `json:"kodeAkun"`
	NamaAkun     string    `json:"namaAkun"`
	JumlahDebit  float64   `json:"jumlahDebit"`
	JumlahKredit float64   `json:"jumlahKredit"`
	Keterangan   string    `json:"keterangan"`
}

// ToResponse mengkonversi Transaksi ke TransaksiResponse
func (t *Transaksi) ToResponse() TransaksiResponse {
	resp := TransaksiResponse{
		ID:                t.ID,
		NomorJurnal:       t.NomorJurnal,
		TanggalTransaksi:  t.TanggalTransaksi,
		Deskripsi:         t.Deskripsi,
		NomorReferensi:    t.NomorReferensi,
		TipeTransaksi:     t.TipeTransaksi,
		TotalDebit:        t.TotalDebit,
		TotalKredit:       t.TotalKredit,
		StatusBalanced:    t.StatusBalanced,
		DibuatOleh:        t.DibuatOleh,
		DiperbaruiOleh:    t.DiperbaruiOleh,
		TanggalDibuat:     t.TanggalDibuat,
		TanggalDiperbarui: t.TanggalDiperbarui,
	}

	// Convert baris transaksi jika ada
	if len(t.BarisTransaksi) > 0 {
		resp.BarisTransaksi = make([]BarisTransaksiResponse, len(t.BarisTransaksi))
		for i, baris := range t.BarisTransaksi {
			resp.BarisTransaksi[i] = BarisTransaksiResponse{
				ID:           baris.ID,
				IDAkun:       baris.IDAkun,
				JumlahDebit:  baris.JumlahDebit,
				JumlahKredit: baris.JumlahKredit,
				Keterangan:   baris.Keterangan,
			}

			// Populate info akun jika relasi sudah di-load
			if baris.Akun.ID != uuid.Nil {
				resp.BarisTransaksi[i].KodeAkun = baris.Akun.KodeAkun
				resp.BarisTransaksi[i].NamaAkun = baris.Akun.NamaAkun
			}
		}
	}

	return resp
}
