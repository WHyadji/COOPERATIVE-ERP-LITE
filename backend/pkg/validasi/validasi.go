// Package validasi menyediakan fungsi-fungsi validasi business logic untuk Cooperative ERP Lite.
//
// Package ini berisi helper validasi untuk berbagai kebutuhan business logic meliputi:
//   - Validasi finansial (jumlah uang, kuantitas, persentase)
//   - Validasi tanggal (tanggal transaksi, tanggal lahir)
//   - Validasi string (teks wajib/opsional dengan batas panjang)
//   - Validasi format (email, nomor HP, kode akun)
//   - Validasi enum (jenis kelamin, status, dll)
//
// Semua error message menggunakan Bahasa Indonesia untuk pengalaman pengguna yang lebih baik.
//
// Penggunaan:
//
//	validator := validasi.Baru()
//	if err := validator.Jumlah(amount, "jumlah setoran"); err != nil {
//	    return err
//	}
//
// Validasi yang tersedia:
//   - Jumlah: Validasi nilai uang (max 999 juta, 2 desimal)
//   - TanggalTransaksi: Validasi tanggal transaksi (tidak future, max 1 tahun)
//   - TanggalLahir: Validasi tanggal lahir (min umur 17 tahun)
//   - TeksWajib: Validasi teks dengan min/max panjang
//   - TeksOpsional: Validasi teks opsional dengan max panjang
//   - Email: Validasi format email
//   - NomorHP: Validasi nomor HP Indonesia (08xx/+628xx)
//   - JenisKelamin: Validasi jenis kelamin (L/P)
//   - Enum: Validasi nilai enum terhadap daftar yang diizinkan
//   - KuantitasProduk: Validasi kuantitas produk (bilangan bulat positif)
//   - Persentase: Validasi nilai persentase (0-100, 2 desimal)
//   - KodeAkun: Validasi format kode akun (XXXX atau XXXX-XX)
package validasi

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"
)

// Validasi adalah struct helper untuk validasi business logic
type Validasi struct{}

// Baru membuat instance Validasi baru
func Baru() *Validasi {
	return &Validasi{}
}

// Jumlah memvalidasi nilai uang/amount
// - Harus lebih dari 0
// - Tidak boleh lebih dari 999,999,999 (999 juta - batas wajar)
// - Maksimal 2 angka di belakang koma
func (v *Validasi) Jumlah(jumlah float64, namaField string) error {
	if jumlah <= 0 {
		return fmt.Errorf("%s harus lebih dari 0", namaField)
	}

	// Batas maksimal 999 juta (batas wajar untuk transaksi koperasi)
	if jumlah > 999999999 {
		return fmt.Errorf("%s terlalu besar (maksimal Rp 999.999.999)", namaField)
	}

	// Validasi presisi (maksimal 2 angka di belakang koma)
	dibulatkan := math.Round(jumlah*100) / 100
	if math.Abs(jumlah-dibulatkan) > 0.001 {
		return fmt.Errorf("%s hanya boleh 2 angka di belakang koma", namaField)
	}

	return nil
}

// TanggalTransaksi memvalidasi tanggal transaksi
// - Tidak boleh di masa depan
// - Tidak boleh lebih dari 1 tahun yang lalu
func (v *Validasi) TanggalTransaksi(tanggal time.Time) error {
	sekarang := time.Now()

	// Tidak boleh di masa depan
	if tanggal.After(sekarang) {
		return errors.New("tanggal transaksi tidak boleh di masa depan")
	}

	// Tidak boleh lebih dari 1 tahun yang lalu
	setahunLalu := sekarang.AddDate(-1, 0, 0)
	if tanggal.Before(setahunLalu) {
		return errors.New("tanggal transaksi terlalu lama (maksimal 1 tahun yang lalu)")
	}

	return nil
}

// TeksWajib memvalidasi teks yang wajib diisi dengan panjang minimal dan maksimal
func (v *Validasi) TeksWajib(teks string, namaField string, panjangMin, panjangMax int) error {
	panjang := len(teks)

	if panjang == 0 {
		return fmt.Errorf("%s wajib diisi", namaField)
	}

	if panjang < panjangMin {
		return fmt.Errorf("%s terlalu pendek (minimal %d karakter)", namaField, panjangMin)
	}

	if panjang > panjangMax {
		return fmt.Errorf("%s terlalu panjang (maksimal %d karakter)", namaField, panjangMax)
	}

	return nil
}

// TeksOpsional memvalidasi teks yang opsional dengan panjang maksimal
func (v *Validasi) TeksOpsional(teks string, namaField string, panjangMax int) error {
	// Jika kosong, tidak perlu validasi
	if teks == "" {
		return nil
	}

	panjang := len(teks)
	if panjang > panjangMax {
		return fmt.Errorf("%s terlalu panjang (maksimal %d karakter)", namaField, panjangMax)
	}

	return nil
}

// Email memvalidasi format email
func (v *Validasi) Email(email string) error {
	// Email opsional
	if email == "" {
		return nil
	}

	// Regex untuk validasi email yang lebih strict
	// Mencegah consecutive dots dan trailing dots
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("format email tidak valid")
	}

	// Validasi panjang email
	if len(email) > 255 {
		return errors.New("email terlalu panjang (maksimal 255 karakter)")
	}

	return nil
}

// NomorHP memvalidasi format nomor HP Indonesia
func (v *Validasi) NomorHP(nomorHP string) error {
	// Nomor HP opsional
	if nomorHP == "" {
		return nil
	}

	// Format Indonesia: 08xx-xxxx-xxxx atau +628xx-xxxx-xxxx
	// Minimal 10 digit, maksimal 14 digit
	hpRegex := regexp.MustCompile(`^(\+62|62|0)[0-9]{9,13}$`)
	if !hpRegex.MatchString(nomorHP) {
		return errors.New("format nomor HP tidak valid (contoh: 08123456789 atau +628123456789)")
	}

	return nil
}

// JenisKelamin memvalidasi jenis kelamin (L atau P)
func (v *Validasi) JenisKelamin(jenisKelamin string) error {
	// Opsional
	if jenisKelamin == "" {
		return nil
	}

	if jenisKelamin != "L" && jenisKelamin != "P" {
		return errors.New("jenis kelamin harus L (Laki-laki) atau P (Perempuan)")
	}

	return nil
}

// Enum memvalidasi nilai enum terhadap daftar nilai yang diizinkan
func (v *Validasi) Enum(nilai string, namaField string, nilaiDiizinkan []string) error {
	if nilai == "" {
		return fmt.Errorf("%s wajib diisi", namaField)
	}

	for _, validNilai := range nilaiDiizinkan {
		if nilai == validNilai {
			return nil
		}
	}

	return fmt.Errorf("%s tidak valid (nilai yang diizinkan: %v)", namaField, nilaiDiizinkan)
}

// KuantitasProduk memvalidasi kuantitas/quantity produk
// - Harus lebih dari 0
// - Harus bilangan bulat (untuk produk yang tidak bisa dipecah)
func (v *Validasi) KuantitasProduk(kuantitas float64, namaField string) error {
	if kuantitas <= 0 {
		return fmt.Errorf("%s harus lebih dari 0", namaField)
	}

	// Maksimal 1 juta unit (batas wajar)
	if kuantitas > 1000000 {
		return fmt.Errorf("%s terlalu besar (maksimal 1.000.000)", namaField)
	}

	// Validasi bilangan bulat (produk tidak boleh pecahan)
	if kuantitas != math.Floor(kuantitas) {
		return fmt.Errorf("%s harus bilangan bulat (tidak boleh ada pecahan)", namaField)
	}

	return nil
}

// Persentase memvalidasi nilai persentase
// - Harus antara 0 dan 100
// - Maksimal 2 angka di belakang koma
func (v *Validasi) Persentase(persentase float64, namaField string) error {
	if persentase < 0 {
		return fmt.Errorf("%s tidak boleh negatif", namaField)
	}

	if persentase > 100 {
		return fmt.Errorf("%s tidak boleh lebih dari 100", namaField)
	}

	// Validasi presisi (maksimal 2 angka di belakang koma)
	dibulatkan := math.Round(persentase*100) / 100
	if math.Abs(persentase-dibulatkan) > 0.001 {
		return fmt.Errorf("%s hanya boleh 2 angka di belakang koma", namaField)
	}

	return nil
}

// KodeAkun memvalidasi format kode akun (Chart of Accounts)
// Format: XXXX atau XXXX-XX (4 digit atau 4 digit + 2 digit)
func (v *Validasi) KodeAkun(kode string) error {
	if kode == "" {
		return errors.New("kode akun wajib diisi")
	}

	// Validasi format kode akun
	kodeRegex := regexp.MustCompile(`^[0-9]{4}(-[0-9]{2})?$`)
	if !kodeRegex.MatchString(kode) {
		return errors.New("format kode akun tidak valid (contoh: 1101 atau 1101-01)")
	}

	return nil
}

// TanggalLahir memvalidasi tanggal lahir
// - Tidak boleh di masa depan
// - Tidak boleh lebih dari 120 tahun yang lalu
// - Harus minimal 17 tahun (usia minimal anggota koperasi)
func (v *Validasi) TanggalLahir(tanggalLahir time.Time) error {
	sekarang := time.Now()

	// Tidak boleh di masa depan
	if tanggalLahir.After(sekarang) {
		return errors.New("tanggal lahir tidak boleh di masa depan")
	}

	// Tidak boleh lebih dari 120 tahun yang lalu
	seratus20TahunLalu := sekarang.AddDate(-120, 0, 0)
	if tanggalLahir.Before(seratus20TahunLalu) {
		return errors.New("tanggal lahir tidak valid (lebih dari 120 tahun yang lalu)")
	}

	// Minimal 17 tahun (usia minimal anggota koperasi)
	tujuhbelasTahunLalu := sekarang.AddDate(-17, 0, 0)
	if tanggalLahir.After(tujuhbelasTahunLalu) {
		return errors.New("usia minimal anggota koperasi adalah 17 tahun")
	}

	return nil
}
