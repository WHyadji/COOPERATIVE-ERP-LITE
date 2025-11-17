package utils

import (
	"strings"
	"testing"
)

// TestValidasiKataSandi_Valid tests that valid strong passwords pass validation
func TestValidasiKataSandi_Valid(t *testing.T) {
	testCases := []struct {
		nama      string
		kataSandi string
	}{
		{"Password dengan semua persyaratan", "Koperasi@2025"},
		{"Password panjang dengan simbol", "Bendahara#12345"},
		{"Password dengan multiple simbol", "MyK0p3r@si!2025"},
		{"Password 10 karakter tepat", "Admin#2025"},
		{"Password dengan underscore", "Simpanan_Aman99!"},
		{"Password dengan berbagai simbol", "Test!@#$%123Abc"},
	}

	for _, tc := range testCases {
		t.Run(tc.nama, func(t *testing.T) {
			err := ValidasiKataSandi(tc.kataSandi)
			if err != nil {
				t.Errorf("Password '%s' seharusnya valid, tapi mendapat error: %v", tc.kataSandi, err)
			}
		})
	}

	t.Logf("✓ Semua password valid lolos validasi")
}

// TestValidasiKataSandi_TooShort tests that passwords shorter than 10 chars are rejected
func TestValidasiKataSandi_TooShort(t *testing.T) {
	testCases := []struct {
		nama      string
		kataSandi string
	}{
		{"6 karakter", "Abc@12"},
		{"8 karakter", "Abcd@123"},
		{"9 karakter", "Abcd@1234"},
		{"Kosong", ""},
		{"5 karakter", "Ab@12"},
	}

	for _, tc := range testCases {
		t.Run(tc.nama, func(t *testing.T) {
			err := ValidasiKataSandi(tc.kataSandi)
			if err == nil {
				t.Errorf("Password '%s' terlalu pendek, seharusnya ditolak", tc.kataSandi)
			}
			if !strings.Contains(err.Error(), "minimal 10 karakter") {
				t.Errorf("Expected error 'minimal 10 karakter', got: %v", err)
			}
		})
	}

	t.Logf("✓ Password terlalu pendek ditolak dengan benar")
}

// TestValidasiKataSandi_NoUppercase tests that passwords without uppercase are rejected
func TestValidasiKataSandi_NoUppercase(t *testing.T) {
	testCases := []string{
		"koperasi@2025",
		"bendahara#12345",
		"simpanan!aman99",
		"admin@123456",
	}

	for _, kataSandi := range testCases {
		t.Run(kataSandi, func(t *testing.T) {
			err := ValidasiKataSandi(kataSandi)
			if err == nil {
				t.Errorf("Password '%s' tanpa huruf besar, seharusnya ditolak", kataSandi)
			}
			if !strings.Contains(err.Error(), "huruf besar") {
				t.Errorf("Expected error tentang 'huruf besar', got: %v", err)
			}
		})
	}

	t.Logf("✓ Password tanpa huruf besar ditolak dengan benar")
}

// TestValidasiKataSandi_NoLowercase tests that passwords without lowercase are rejected
func TestValidasiKataSandi_NoLowercase(t *testing.T) {
	testCases := []string{
		"KOPERASI@2025",
		"BENDAHARA#12345",
		"SIMPANAN!AMAN99",
		"ADMIN@123456",
	}

	for _, kataSandi := range testCases {
		t.Run(kataSandi, func(t *testing.T) {
			err := ValidasiKataSandi(kataSandi)
			if err == nil {
				t.Errorf("Password '%s' tanpa huruf kecil, seharusnya ditolak", kataSandi)
			}
			if !strings.Contains(err.Error(), "huruf kecil") {
				t.Errorf("Expected error tentang 'huruf kecil', got: %v", err)
			}
		})
	}

	t.Logf("✓ Password tanpa huruf kecil ditolak dengan benar")
}

// TestValidasiKataSandi_NoDigit tests that passwords without numbers are rejected
func TestValidasiKataSandi_NoDigit(t *testing.T) {
	testCases := []string{
		"Koperasi@Aman",
		"Bendahara#Sukses",
		"Simpanan!Aman",
		"Admin@Koperasi",
	}

	for _, kataSandi := range testCases {
		t.Run(kataSandi, func(t *testing.T) {
			err := ValidasiKataSandi(kataSandi)
			if err == nil {
				t.Errorf("Password '%s' tanpa angka, seharusnya ditolak", kataSandi)
			}
			if !strings.Contains(err.Error(), "angka") {
				t.Errorf("Expected error tentang 'angka', got: %v", err)
			}
		})
	}

	t.Logf("✓ Password tanpa angka ditolak dengan benar")
}

// TestValidasiKataSandi_NoSpecialChar tests that passwords without special chars are rejected
func TestValidasiKataSandi_NoSpecialChar(t *testing.T) {
	testCases := []string{
		"Koperasi2025",
		"Bendahara12345",
		"SimpananAman99",
		"Admin123456",
	}

	for _, kataSandi := range testCases {
		t.Run(kataSandi, func(t *testing.T) {
			err := ValidasiKataSandi(kataSandi)
			if err == nil {
				t.Errorf("Password '%s' tanpa karakter spesial, seharusnya ditolak", kataSandi)
			}
			if !strings.Contains(err.Error(), "karakter spesial") {
				t.Errorf("Expected error tentang 'karakter spesial', got: %v", err)
			}
		})
	}

	t.Logf("✓ Password tanpa karakter spesial ditolak dengan benar")
}

// TestValidasiKataSandi_CommonPasswords tests that common weak passwords are rejected
func TestValidasiKataSandi_CommonPasswords(t *testing.T) {
	testCases := []struct {
		nama      string
		kataSandi string
	}{
		{"password generic", "Password123!"},
		{"password dengan angka", "Password1234!"},
		{"12345678 dengan simbol", "12345678!A"},
		{"qwerty dengan simbol", "Qwerty123!"},
		{"admin dengan simbol", "Admin123!@"},
		{"administrator", "Administrator123!"},
		{"koperasi umum", "Koperasi123!"},
		{"bendahara umum", "Bendahara123!"},
		{"kasir umum", "Kasir123!@#"},
		{"simpanan umum", "Simpanan123!"},
		{"welcome", "Welcome123!"},
		{"P@ssw0rd", "P@ssw0rd12"},
	}

	for _, tc := range testCases {
		t.Run(tc.nama, func(t *testing.T) {
			err := ValidasiKataSandi(tc.kataSandi)
			if err == nil {
				t.Errorf("Password '%s' adalah password umum, seharusnya ditolak", tc.kataSandi)
				return
			}
			if !strings.Contains(err.Error(), "terlalu umum") {
				t.Errorf("Expected error tentang password 'terlalu umum', got: %v", err)
			}
		})
	}

	t.Logf("✓ Password umum/lemah ditolak dengan benar")
}

// TestValidasiKataSandi_EdgeCases tests edge cases
func TestValidasiKataSandi_EdgeCases(t *testing.T) {
	testCases := []struct {
		nama            string
		kataSandi       string
		seharusnyaValid bool
		pesanError      string
	}{
		{"Tepat 10 karakter valid", "Admin@2025", true, ""},
		{"Sangat panjang valid", "ThisIsAVeryLongAndSecurePassword123!@#", true, ""},
		{"Dengan spasi valid", "My Pass@2025", true, ""},
		{"Unicode karakter", "Koperasi™2025!", true, ""},
		{"Hanya simbol tanpa huruf", "!@#$%^&*()", false, "huruf"},
		{"Hanya angka", "1234567890", false, "huruf"},
	}

	for _, tc := range testCases {
		t.Run(tc.nama, func(t *testing.T) {
			err := ValidasiKataSandi(tc.kataSandi)
			if tc.seharusnyaValid {
				if err != nil {
					t.Errorf("Password '%s' seharusnya valid, tapi mendapat error: %v", tc.kataSandi, err)
				}
			} else {
				if err == nil {
					t.Errorf("Password '%s' seharusnya tidak valid", tc.kataSandi)
				}
				if tc.pesanError != "" && !strings.Contains(err.Error(), tc.pesanError) {
					t.Errorf("Expected error containing '%s', got: %v", tc.pesanError, err)
				}
			}
		})
	}

	t.Logf("✓ Edge cases ditangani dengan benar")
}

// TestDapatkanPersyaratanKataSandi tests that requirements list is returned correctly
func TestDapatkanPersyaratanKataSandi(t *testing.T) {
	persyaratan := DapatkanPersyaratanKataSandi()

	if len(persyaratan) == 0 {
		t.Error("Persyaratan kata sandi tidak boleh kosong")
	}

	// Verify expected requirements are present
	expectedKeywords := []string{
		"10 karakter",
		"huruf besar",
		"huruf kecil",
		"angka",
		"karakter spesial",
		"umum",
	}

	persyaratanString := strings.Join(persyaratan, " ")
	for _, keyword := range expectedKeywords {
		if !strings.Contains(persyaratanString, keyword) {
			t.Errorf("Persyaratan seharusnya mengandung '%s'", keyword)
		}
	}

	t.Logf("✓ Persyaratan kata sandi: %v", persyaratan)
}

// TestContohKataSandiKuat tests that password examples are valid
func TestContohKataSandiKuat(t *testing.T) {
	contohList := ContohKataSandiKuat()

	if len(contohList) == 0 {
		t.Error("Contoh kata sandi tidak boleh kosong")
	}

	// Verify all examples pass validation
	for _, contoh := range contohList {
		err := ValidasiKataSandi(contoh)
		if err != nil {
			t.Errorf("Contoh password '%s' seharusnya valid, tapi mendapat error: %v", contoh, err)
		}
	}

	t.Logf("✓ Semua contoh kata sandi valid: %v", contohList)
}

// TestHitungKekuatanKataSandi tests password strength calculation
func TestHitungKekuatanKataSandi(t *testing.T) {
	testCases := []struct {
		nama            string
		kataSandi       string
		levelDiharapkan string
		skorMinimal     int
	}{
		{"Password lemah", "abc123", "lemah", 0},
		{"Password lemah 2", "password", "lemah", 0},
		{"Password sedang", "Password123", "sedang", 3},
		{"Password sedang 2", "Koperasi2025", "sedang", 3},
		{"Password kuat", "Koperasi@2025", "kuat", 5},
		{"Password kuat 2", "MyK0p3r@si!2025", "kuat", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.nama, func(t *testing.T) {
			kekuatan := HitungKekuatanKataSandi(tc.kataSandi)

			if kekuatan.Level != tc.levelDiharapkan {
				t.Errorf("Password '%s': diharapkan level '%s', mendapat '%s'",
					tc.kataSandi, tc.levelDiharapkan, kekuatan.Level)
			}

			if kekuatan.Skor < tc.skorMinimal {
				t.Errorf("Password '%s': diharapkan skor minimal %d, mendapat %d",
					tc.kataSandi, tc.skorMinimal, kekuatan.Skor)
			}

			if kekuatan.Keterangan == "" {
				t.Error("Keterangan tidak boleh kosong")
			}

			t.Logf("Password '%s': Level=%s, Skor=%d, Keterangan=%s, Saran=%v",
				tc.kataSandi, kekuatan.Level, kekuatan.Skor, kekuatan.Keterangan, kekuatan.Saran)
		})
	}

	t.Logf("✓ Penghitungan kekuatan kata sandi bekerja dengan benar")
}

// TestHitungKekuatanKataSandi_Saran tests that suggestions are provided for weak passwords
func TestHitungKekuatanKataSandi_Saran(t *testing.T) {
	testCases := []struct {
		nama               string
		kataSandi          string
		seharusnyaAdaSaran bool
	}{
		{"Password sangat lemah", "abc", true},
		{"Password tanpa angka", "abcdefghij", true},
		{"Password tanpa huruf besar", "abcdefgh123", true},
		{"Password kuat", "Koperasi@2025", false},
	}

	for _, tc := range testCases {
		t.Run(tc.nama, func(t *testing.T) {
			kekuatan := HitungKekuatanKataSandi(tc.kataSandi)

			if tc.seharusnyaAdaSaran {
				if len(kekuatan.Saran) == 0 {
					t.Errorf("Password lemah '%s' seharusnya mendapat saran", tc.kataSandi)
				}
			}

			if len(kekuatan.Saran) > 0 {
				t.Logf("Saran untuk '%s': %v", tc.kataSandi, kekuatan.Saran)
			}
		})
	}

	t.Logf("✓ Saran untuk password lemah diberikan dengan benar")
}

// BenchmarkValidasiKataSandi benchmarks password validation performance
func BenchmarkValidasiKataSandi(b *testing.B) {
	kataSandi := "Koperasi@2025"
	for i := 0; i < b.N; i++ {
		ValidasiKataSandi(kataSandi)
	}
}

// BenchmarkHitungKekuatanKataSandi benchmarks password strength calculation performance
func BenchmarkHitungKekuatanKataSandi(b *testing.B) {
	kataSandi := "Koperasi@2025"
	for i := 0; i < b.N; i++ {
		HitungKekuatanKataSandi(kataSandi)
	}
}
