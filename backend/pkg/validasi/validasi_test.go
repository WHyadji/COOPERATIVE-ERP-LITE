package validasi

import (
	"testing"
	"time"
)

// TestJumlah menguji validasi amount/jumlah
func TestJumlah(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		jumlah      float64
		namaField   string
		shouldError bool
	}{
		{"Valid amount 100", 100.50, "test", false},
		{"Valid amount 0.01", 0.01, "test", false},
		{"Valid amount max 2 decimal", 1234.56, "test", false},
		{"Invalid amount 0", 0, "test", true},
		{"Invalid negative", -10, "test", true},
		{"Invalid too large", 9999999999, "test", true},
		{"Invalid too many decimals", 10.12345, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Jumlah(tt.jumlah, tt.namaField)
			if (err != nil) != tt.shouldError {
				t.Errorf("Jumlah() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestTanggalTransaksi menguji validasi tanggal transaksi
func TestTanggalTransaksi(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		tanggal     time.Time
		shouldError bool
	}{
		{"Valid today", time.Now(), false},
		{"Valid yesterday", time.Now().AddDate(0, 0, -1), false},
		{"Valid 1 month ago", time.Now().AddDate(0, -1, 0), false},
		{"Invalid future date", time.Now().AddDate(0, 0, 1), true},
		{"Invalid more than 1 year ago", time.Now().AddDate(-1, 0, -1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.TanggalTransaksi(tt.tanggal)
			if (err != nil) != tt.shouldError {
				t.Errorf("TanggalTransaksi() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestTeksWajib menguji validasi teks wajib
func TestTeksWajib(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		teks        string
		namaField   string
		panjangMin  int
		panjangMax  int
		shouldError bool
	}{
		{"Valid text", "Hello World", "test", 3, 255, false},
		{"Valid min length", "abc", "test", 3, 255, false},
		{"Invalid empty", "", "test", 3, 255, true},
		{"Invalid too short", "ab", "test", 3, 255, true},
		{"Invalid too long", "a very long text that exceeds the maximum allowed length", "test", 3, 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.TeksWajib(tt.teks, tt.namaField, tt.panjangMin, tt.panjangMax)
			if (err != nil) != tt.shouldError {
				t.Errorf("TeksWajib() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestTeksOpsional menguji validasi teks opsional
func TestTeksOpsional(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		teks        string
		namaField   string
		panjangMax  int
		shouldError bool
	}{
		{"Valid text", "Hello World", "test", 255, false},
		{"Valid empty (opsional)", "", "test", 255, false},
		{"Invalid too long", "a very long text that exceeds the maximum allowed length", "test", 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.TeksOpsional(tt.teks, tt.namaField, tt.panjangMax)
			if (err != nil) != tt.shouldError {
				t.Errorf("TeksOpsional() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestEmail menguji validasi email
func TestEmail(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		email       string
		shouldError bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid email with subdomain", "user@mail.example.com", false},
		{"Valid empty (opsional)", "", false},
		{"Invalid format", "invalid-email", true},
		{"Invalid missing @", "test.example.com", true},
		{"Invalid missing domain", "test@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Email(tt.email)
			if (err != nil) != tt.shouldError {
				t.Errorf("Email() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestNomorHP menguji validasi nomor HP Indonesia
func TestNomorHP(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		nomorHP     string
		shouldError bool
	}{
		{"Valid with 08", "08123456789", false},
		{"Valid with +62", "+628123456789", false},
		{"Valid with 62", "628123456789", false},
		{"Valid empty (opsional)", "", false},
		{"Invalid too short", "081234", true},
		{"Invalid format", "123456789", true},
		{"Invalid non-numeric", "081234abcde", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.NomorHP(tt.nomorHP)
			if (err != nil) != tt.shouldError {
				t.Errorf("NomorHP() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestJenisKelamin menguji validasi jenis kelamin
func TestJenisKelamin(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		jk          string
		shouldError bool
	}{
		{"Valid L", "L", false},
		{"Valid P", "P", false},
		{"Valid empty (opsional)", "", false},
		{"Invalid value", "M", true},
		{"Invalid value", "X", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.JenisKelamin(tt.jk)
			if (err != nil) != tt.shouldError {
				t.Errorf("JenisKelamin() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestEnum menguji validasi enum
func TestEnum(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name            string
		nilai           string
		namaField       string
		nilaiDiizinkan  []string
		shouldError     bool
	}{
		{"Valid value", "aktif", "status", []string{"aktif", "nonaktif"}, false},
		{"Valid value 2", "nonaktif", "status", []string{"aktif", "nonaktif"}, false},
		{"Invalid value", "pending", "status", []string{"aktif", "nonaktif"}, true},
		{"Invalid empty", "", "status", []string{"aktif", "nonaktif"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Enum(tt.nilai, tt.namaField, tt.nilaiDiizinkan)
			if (err != nil) != tt.shouldError {
				t.Errorf("Enum() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestKuantitasProduk menguji validasi kuantitas produk
func TestKuantitasProduk(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		kuantitas   float64
		namaField   string
		shouldError bool
	}{
		{"Valid quantity 1", 1, "test", false},
		{"Valid quantity 100", 100, "test", false},
		{"Invalid 0", 0, "test", true},
		{"Invalid negative", -10, "test", true},
		{"Invalid too large", 1000001, "test", true},
		{"Invalid fractional 1.5", 1.5, "test", true},
		{"Invalid fractional 10.25", 10.25, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.KuantitasProduk(tt.kuantitas, tt.namaField)
			if (err != nil) != tt.shouldError {
				t.Errorf("KuantitasProduk() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestPersentase menguji validasi persentase
func TestPersentase(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		persentase  float64
		namaField   string
		shouldError bool
	}{
		{"Valid 0%", 0, "test", false},
		{"Valid 50%", 50, "test", false},
		{"Valid 100%", 100, "test", false},
		{"Valid 50.25%", 50.25, "test", false},
		{"Invalid negative", -10, "test", true},
		{"Invalid over 100", 101, "test", true},
		{"Invalid too many decimals", 50.123, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Persentase(tt.persentase, tt.namaField)
			if (err != nil) != tt.shouldError {
				t.Errorf("Persentase() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestKodeAkun menguji validasi kode akun
func TestKodeAkun(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		kode        string
		shouldError bool
	}{
		{"Valid 4 digit", "1101", false},
		{"Valid with sub-account", "1101-01", false},
		{"Invalid empty", "", true},
		{"Invalid too short", "110", true},
		{"Invalid non-numeric", "ABCD", true},
		{"Invalid format", "1101-1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.KodeAkun(tt.kode)
			if (err != nil) != tt.shouldError {
				t.Errorf("KodeAkun() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

// TestTanggalLahir menguji validasi tanggal lahir
func TestTanggalLahir(t *testing.T) {
	validator := Baru()

	tests := []struct {
		name        string
		tanggal     time.Time
		shouldError bool
	}{
		{"Valid 20 years old", time.Now().AddDate(-20, 0, 0), false},
		{"Valid 17 years old (exact)", time.Now().AddDate(-17, 0, 0), false},
		{"Valid 50 years old", time.Now().AddDate(-50, 0, 0), false},
		{"Invalid future date", time.Now().AddDate(0, 0, 1), true},
		{"Invalid too young (16 years)", time.Now().AddDate(-16, 0, 0), true},
		{"Invalid too old (121 years)", time.Now().AddDate(-121, 0, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.TanggalLahir(tt.tanggal)
			if (err != nil) != tt.shouldError {
				t.Errorf("TanggalLahir() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}
