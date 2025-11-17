package utils

import (
	"errors"
	"strings"
	"unicode"
)

// ValidasiKataSandi memvalidasi kekuatan kata sandi dengan persyaratan keamanan yang kuat
// Persyaratan:
// - Minimal 10 karakter
// - Minimal satu huruf besar (uppercase)
// - Minimal satu huruf kecil (lowercase)
// - Minimal satu angka
// - Minimal satu karakter spesial (!@#$%^&*()_+-=[]{}|;:,.<>?)
// - Tidak boleh menggunakan kata sandi yang umum/lemah
func ValidasiKataSandi(kataSandi string) error {
	// Validasi panjang minimal
	if len(kataSandi) < 10 {
		return errors.New("kata sandi minimal 10 karakter")
	}

	var (
		adaHurufBesar      = false
		adaHurufKecil      = false
		adaAngka           = false
		adaKarakterSpesial = false
	)

	// Periksa setiap karakter
	for _, karakter := range kataSandi {
		switch {
		case unicode.IsUpper(karakter):
			adaHurufBesar = true
		case unicode.IsLower(karakter):
			adaHurufKecil = true
		case unicode.IsDigit(karakter):
			adaAngka = true
		case unicode.IsPunct(karakter) || unicode.IsSymbol(karakter):
			adaKarakterSpesial = true
		}
	}

	// Validasi huruf besar
	if !adaHurufBesar {
		return errors.New("kata sandi harus mengandung minimal satu huruf besar")
	}

	// Validasi huruf kecil
	if !adaHurufKecil {
		return errors.New("kata sandi harus mengandung minimal satu huruf kecil")
	}

	// Validasi angka
	if !adaAngka {
		return errors.New("kata sandi harus mengandung minimal satu angka")
	}

	// Validasi karakter spesial
	if !adaKarakterSpesial {
		return errors.New("kata sandi harus mengandung minimal satu karakter spesial (!@#$%^&*)")
	}

	// Cek terhadap kata sandi yang umum/lemah
	kataSandiUmum := []string{
		"password", "password12", "password123", "password1234",
		"12345678", "12345678aa", "12345678a", "123456789", "1234567890",
		"qwerty", "qwerty123", "qwertyuiop",
		"admin", "admin123", "admin1234", "administrator", "administrator123",
		"koperasi", "koperasi123",
		"bendahara", "bendahara123",
		"kasir", "kasir123",
		"simpanan", "simpanan123",
		"welcome", "welcome123",
		"letmein", "letmein123",
		"abc123", "abc12345",
		"p@ssw0rd", "p@ssw0rd12", "pssw0rd12", "p@ssword", "passw0rd", "password12",
	}

	// Hapus karakter spesial untuk deteksi variasi password umum
	// Contoh: "Password123!" akan menjadi "password123"
	kataSandiTanpaSimbol := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, kataSandi)

	kataSandiLower := strings.ToLower(kataSandi)
	for _, kataSandiLemah := range kataSandiUmum {
		kataSandiLemahLower := strings.ToLower(kataSandiLemah)
		// Cek exact match (case-insensitive)
		if kataSandiLower == kataSandiLemahLower {
			return errors.New("kata sandi terlalu umum, gunakan kombinasi yang lebih unik")
		}
		// Cek match tanpa simbol (untuk deteksi "Password123!" == "password123")
		if kataSandiTanpaSimbol == kataSandiLemahLower {
			return errors.New("kata sandi terlalu umum, gunakan kombinasi yang lebih unik")
		}
	}

	return nil
}

// DapatkanPersyaratanKataSandi mengembalikan daftar persyaratan kata sandi
// Digunakan untuk menampilkan panduan kepada pengguna
func DapatkanPersyaratanKataSandi() []string {
	return []string{
		"Minimal 10 karakter",
		"Minimal satu huruf besar (A-Z)",
		"Minimal satu huruf kecil (a-z)",
		"Minimal satu angka (0-9)",
		"Minimal satu karakter spesial (!@#$%^&*)",
		"Tidak boleh menggunakan kata sandi yang umum",
	}
}

// ContohKataSandiKuat mengembalikan contoh kata sandi yang kuat
// Digunakan untuk memberikan panduan kepada pengguna
func ContohKataSandiKuat() []string {
	return []string{
		"Koperasi@2025",
		"SimpananSaya#456",
		"Anggota!Kuat99",
		"MyK0p3r@si!",
	}
}

// KekuatanKataSandi menghitung skor kekuatan kata sandi (0-5)
// 0-2: Lemah, 3-4: Sedang, 5-6: Kuat
type KekuatanKataSandi struct {
	Skor       int      `json:"skor"`
	Level      string   `json:"level"`
	Keterangan string   `json:"keterangan"`
	Saran      []string `json:"saran"`
}

// HitungKekuatanKataSandi menghitung kekuatan kata sandi dan memberikan saran
func HitungKekuatanKataSandi(kataSandi string) *KekuatanKataSandi {
	skor := 0
	saran := []string{}

	// Panjang kata sandi
	if len(kataSandi) >= 8 {
		skor++
	} else {
		saran = append(saran, "Tambahkan lebih banyak karakter (minimal 10)")
	}

	if len(kataSandi) >= 12 {
		skor++
	} else if len(kataSandi) >= 8 {
		saran = append(saran, "Gunakan minimal 12 karakter untuk keamanan maksimal")
	}

	// Periksa kompleksitas
	var (
		adaHurufBesar      = false
		adaHurufKecil      = false
		adaAngka           = false
		adaKarakterSpesial = false
	)

	for _, karakter := range kataSandi {
		switch {
		case unicode.IsUpper(karakter):
			adaHurufBesar = true
		case unicode.IsLower(karakter):
			adaHurufKecil = true
		case unicode.IsDigit(karakter):
			adaAngka = true
		case unicode.IsPunct(karakter) || unicode.IsSymbol(karakter):
			adaKarakterSpesial = true
		}
	}

	if adaHurufKecil {
		skor++
	} else {
		saran = append(saran, "Tambahkan huruf kecil (a-z)")
	}

	if adaHurufBesar {
		skor++
	} else {
		saran = append(saran, "Tambahkan huruf besar (A-Z)")
	}

	if adaAngka {
		skor++
	} else {
		saran = append(saran, "Tambahkan angka (0-9)")
	}

	if adaKarakterSpesial {
		skor++
	} else {
		saran = append(saran, "Tambahkan karakter spesial (!@#$%^&*)")
	}

	// Tentukan level
	var level, keterangan string
	if skor <= 2 {
		level = "lemah"
		keterangan = "Kata sandi Anda sangat lemah dan mudah ditebak"
	} else if skor <= 4 {
		level = "sedang"
		keterangan = "Kata sandi Anda cukup baik, tetapi bisa ditingkatkan"
	} else {
		level = "kuat"
		keterangan = "Kata sandi Anda kuat dan aman"
	}

	return &KekuatanKataSandi{
		Skor:       skor,
		Level:      level,
		Keterangan: keterangan,
		Saran:      saran,
	}
}
