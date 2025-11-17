# Package Validasi

Package `validasi` menyediakan helper functions untuk validasi business logic di layer service. Package ini dirancang khusus untuk Cooperative ERP Lite dengan naming convention bahasa Indonesia yang konsisten.

## 📋 Daftar Isi

- [Instalasi](#instalasi)
- [Penggunaan Dasar](#penggunaan-dasar)
- [Fungsi-Fungsi Validasi](#fungsi-fungsi-validasi)
- [Contoh Penggunaan](#contoh-penggunaan)
- [Testing](#testing)
- [Best Practices](#best-practices)

## 🚀 Instalasi

Package ini sudah tersedia di `backend/pkg/validasi`. Untuk menggunakannya, import di service Anda:

```go
import "cooperative-erp-lite/pkg/validasi"
```

## 💡 Penggunaan Dasar

### Inisialisasi Validator

```go
validator := validasi.Baru()
```

### Contoh Validasi Sederhana

```go
// Validasi jumlah uang
if err := validator.Jumlah(100000, "harga produk"); err != nil {
    return nil, err
}

// Validasi email
if err := validator.Email("user@example.com"); err != nil {
    return nil, err
}

// Validasi tanggal transaksi
if err := validator.TanggalTransaksi(time.Now()); err != nil {
    return nil, err
}
```

## 📚 Fungsi-Fungsi Validasi

### 1. Validasi Jumlah/Amount

```go
func (v *Validasi) Jumlah(jumlah float64, namaField string) error
```

**Validasi:**
- Harus lebih dari 0
- Maksimal 999,999,999 (999 juta rupiah)
- Maksimal 2 angka di belakang koma

**Contoh:**
```go
// ✅ Valid
validator.Jumlah(100.50, "total bayar")      // OK
validator.Jumlah(1234.56, "harga")          // OK

// ❌ Invalid
validator.Jumlah(0, "total")                 // Error: harus lebih dari 0
validator.Jumlah(-100, "total")              // Error: harus lebih dari 0
validator.Jumlah(999999999999, "total")      // Error: terlalu besar
validator.Jumlah(10.12345, "total")          // Error: max 2 desimal
```

---

### 2. Validasi Tanggal Transaksi

```go
func (v *Validasi) TanggalTransaksi(tanggal time.Time) error
```

**Validasi:**
- Tidak boleh di masa depan
- Tidak boleh lebih dari 1 tahun yang lalu

**Contoh:**
```go
// ✅ Valid
validator.TanggalTransaksi(time.Now())                    // OK
validator.TanggalTransaksi(time.Now().AddDate(0, -6, 0))  // OK - 6 bulan lalu

// ❌ Invalid
validator.TanggalTransaksi(time.Now().AddDate(0, 0, 1))   // Error: masa depan
validator.TanggalTransaksi(time.Now().AddDate(-2, 0, 0))  // Error: terlalu lama
```

---

### 3. Validasi Teks Wajib

```go
func (v *Validasi) TeksWajib(teks string, namaField string, panjangMin, panjangMax int) error
```

**Validasi:**
- Tidak boleh kosong
- Panjang minimal sesuai parameter
- Panjang maksimal sesuai parameter

**Contoh:**
```go
// ✅ Valid
validator.TeksWajib("John Doe", "nama", 3, 255)          // OK
validator.TeksWajib("abc", "username", 3, 50)            // OK

// ❌ Invalid
validator.TeksWajib("", "nama", 3, 255)                  // Error: wajib diisi
validator.TeksWajib("ab", "nama", 3, 255)                // Error: terlalu pendek
validator.TeksWajib("nama yang sangat panjang...", "nama", 3, 10)  // Error: terlalu panjang
```

---

### 4. Validasi Teks Opsional

```go
func (v *Validasi) TeksOpsional(teks string, namaField string, panjangMax int) error
```

**Validasi:**
- Boleh kosong (opsional)
- Jika diisi, maksimal sesuai parameter

**Contoh:**
```go
// ✅ Valid
validator.TeksOpsional("", "keterangan", 500)            // OK - boleh kosong
validator.TeksOpsional("Catatan singkat", "keterangan", 500)  // OK

// ❌ Invalid
validator.TeksOpsional("teks yang sangat panjang...", "keterangan", 10)  // Error: terlalu panjang
```

---

### 5. Validasi Email

```go
func (v *Validasi) Email(email string) error
```

**Validasi:**
- Boleh kosong (opsional)
- Format email yang valid
- Maksimal 255 karakter

**Contoh:**
```go
// ✅ Valid
validator.Email("")                           // OK - opsional
validator.Email("user@example.com")           // OK
validator.Email("admin@mail.example.com")     // OK

// ❌ Invalid
validator.Email("invalid-email")              // Error: format tidak valid
validator.Email("test@")                      // Error: format tidak valid
validator.Email("@example.com")               // Error: format tidak valid
```

---

### 6. Validasi Nomor HP (Indonesia)

```go
func (v *Validasi) NomorHP(nomorHP string) error
```

**Validasi:**
- Boleh kosong (opsional)
- Format Indonesia: 08xx, +62xx, atau 62xx
- Minimal 10 digit, maksimal 14 digit

**Contoh:**
```go
// ✅ Valid
validator.NomorHP("")                         // OK - opsional
validator.NomorHP("08123456789")              // OK
validator.NomorHP("+628123456789")            // OK
validator.NomorHP("628123456789")             // OK

// ❌ Invalid
validator.NomorHP("081234")                   // Error: terlalu pendek
validator.NomorHP("123456789")                // Error: format tidak valid
```

---

### 7. Validasi Jenis Kelamin

```go
func (v *Validasi) JenisKelamin(jenisKelamin string) error
```

**Validasi:**
- Boleh kosong (opsional)
- Harus "L" (Laki-laki) atau "P" (Perempuan)

**Contoh:**
```go
// ✅ Valid
validator.JenisKelamin("")       // OK - opsional
validator.JenisKelamin("L")      // OK
validator.JenisKelamin("P")      // OK

// ❌ Invalid
validator.JenisKelamin("M")      // Error: harus L atau P
validator.JenisKelamin("X")      // Error: harus L atau P
```

---

### 8. Validasi Enum

```go
func (v *Validasi) Enum(nilai string, namaField string, nilaiDiizinkan []string) error
```

**Validasi:**
- Tidak boleh kosong
- Harus ada dalam daftar nilai yang diizinkan

**Contoh:**
```go
// ✅ Valid
validator.Enum("aktif", "status", []string{"aktif", "nonaktif"})     // OK
validator.Enum("nonaktif", "status", []string{"aktif", "nonaktif"})  // OK

// ❌ Invalid
validator.Enum("pending", "status", []string{"aktif", "nonaktif"})   // Error: tidak valid
validator.Enum("", "status", []string{"aktif", "nonaktif"})          // Error: wajib diisi
```

---

### 9. Validasi Kuantitas Produk

```go
func (v *Validasi) KuantitasProduk(kuantitas float64, namaField string) error
```

**Validasi:**
- Harus lebih dari 0
- Maksimal 1,000,000 unit

**Contoh:**
```go
// ✅ Valid
validator.KuantitasProduk(1, "jumlah")        // OK
validator.KuantitasProduk(100, "jumlah")      // OK

// ❌ Invalid
validator.KuantitasProduk(0, "jumlah")        // Error: harus lebih dari 0
validator.KuantitasProduk(-10, "jumlah")      // Error: harus lebih dari 0
validator.KuantitasProduk(1000001, "jumlah")  // Error: terlalu besar
```

---

### 10. Validasi Persentase

```go
func (v *Validasi) Persentase(persentase float64, namaField string) error
```

**Validasi:**
- Harus antara 0 dan 100
- Maksimal 2 angka di belakang koma

**Contoh:**
```go
// ✅ Valid
validator.Persentase(0, "diskon")         // OK
validator.Persentase(50, "diskon")        // OK
validator.Persentase(100, "diskon")       // OK
validator.Persentase(50.25, "diskon")     // OK

// ❌ Invalid
validator.Persentase(-10, "diskon")       // Error: tidak boleh negatif
validator.Persentase(101, "diskon")       // Error: tidak boleh lebih dari 100
validator.Persentase(50.123, "diskon")    // Error: max 2 desimal
```

---

### 11. Validasi Kode Akun

```go
func (v *Validasi) KodeAkun(kode string) error
```

**Validasi:**
- Tidak boleh kosong
- Format: XXXX atau XXXX-XX (4 digit atau 4 digit + 2 digit)

**Contoh:**
```go
// ✅ Valid
validator.KodeAkun("1101")        // OK
validator.KodeAkun("1101-01")     // OK
validator.KodeAkun("5201-05")     // OK

// ❌ Invalid
validator.KodeAkun("")            // Error: wajib diisi
validator.KodeAkun("110")         // Error: format tidak valid
validator.KodeAkun("ABCD")        // Error: format tidak valid
validator.KodeAkun("1101-1")      // Error: format tidak valid
```

---

### 12. Validasi Tanggal Lahir

```go
func (v *Validasi) TanggalLahir(tanggalLahir time.Time) error
```

**Validasi:**
- Tidak boleh di masa depan
- Tidak boleh lebih dari 120 tahun yang lalu
- Minimal usia 17 tahun (usia minimal anggota koperasi)

**Contoh:**
```go
// ✅ Valid
validator.TanggalLahir(time.Now().AddDate(-20, 0, 0))   // OK - 20 tahun
validator.TanggalLahir(time.Now().AddDate(-17, 0, 0))   // OK - 17 tahun
validator.TanggalLahir(time.Now().AddDate(-50, 0, 0))   // OK - 50 tahun

// ❌ Invalid
validator.TanggalLahir(time.Now().AddDate(0, 0, 1))     // Error: masa depan
validator.TanggalLahir(time.Now().AddDate(-16, 0, 0))   // Error: terlalu muda
validator.TanggalLahir(time.Now().AddDate(-121, 0, 0))  // Error: terlalu tua
```

---

## 🔍 Contoh Penggunaan di Service

### Contoh 1: Validasi di Simpanan Service

```go
func (s *SimpananService) CatatSetoran(idKoperasi, idPengguna uuid.UUID, req *CatatSetoranRequest) (*models.SimpananResponse, error) {
    // Initialize validator
    validator := validasi.Baru()

    // Validasi business logic
    if err := validator.Jumlah(req.JumlahSetoran, "jumlah setoran"); err != nil {
        return nil, err
    }

    if err := validator.TanggalTransaksi(req.TanggalTransaksi); err != nil {
        return nil, err
    }

    if err := validator.TeksOpsional(req.Keterangan, "keterangan", 500); err != nil {
        return nil, err
    }

    // Lanjutkan dengan business logic...
}
```

### Contoh 2: Validasi di Anggota Service

```go
func (s *AnggotaService) BuatAnggota(idKoperasi uuid.UUID, req *BuatAnggotaRequest) (*models.AnggotaResponse, error) {
    validator := validasi.Baru()

    // Validasi field wajib
    if err := validator.TeksWajib(req.NamaLengkap, "nama lengkap", 3, 255); err != nil {
        return nil, err
    }

    if err := validator.Email(req.Email); err != nil {
        return nil, err
    }

    if err := validator.NomorHP(req.NoTelepon); err != nil {
        return nil, err
    }

    if err := validator.JenisKelamin(req.JenisKelamin); err != nil {
        return nil, err
    }

    // Validasi tanggal lahir jika ada
    if req.TanggalLahir != nil {
        if err := validator.TanggalLahir(*req.TanggalLahir); err != nil {
            return nil, err
        }
    }

    // Lanjutkan dengan business logic...
}
```

### Contoh 3: Validasi di Transaksi Service

```go
func (s *TransaksiService) BuatTransaksi(idKoperasi, idPengguna uuid.UUID, req *BuatTransaksiRequest) (*models.TransaksiResponse, error) {
    validator := validasi.Baru()

    // Validasi header transaksi
    if err := validator.TanggalTransaksi(req.TanggalTransaksi); err != nil {
        return nil, err
    }

    if err := validator.TeksWajib(req.Deskripsi, "deskripsi", 5, 500); err != nil {
        return nil, err
    }

    // Validasi baris transaksi
    for i, baris := range req.BarisTransaksi {
        if baris.JumlahDebit > 0 {
            if err := validator.Jumlah(baris.JumlahDebit, fmt.Sprintf("jumlah debit baris ke-%d", i+1)); err != nil {
                return nil, err
            }
        }

        if baris.JumlahKredit > 0 {
            if err := validator.Jumlah(baris.JumlahKredit, fmt.Sprintf("jumlah kredit baris ke-%d", i+1)); err != nil {
                return nil, err
            }
        }
    }

    // Lanjutkan dengan business logic...
}
```

## 🧪 Testing

Package ini dilengkapi dengan unit tests komprehensif. Untuk menjalankan tests:

```bash
# Run tests
cd backend
go test ./pkg/validasi/... -v

# Run tests with coverage
go test ./pkg/validasi/... -cover

# Run specific test
go test ./pkg/validasi/... -v -run TestJumlah
```

**Test Coverage:**
- 12 test suites
- 60+ test scenarios
- Coverage: 100% untuk semua fungsi validasi

## ✅ Best Practices

### 1. Selalu Inisialisasi Validator di Awal Function

```go
// ✅ Good
func (s *Service) Create(req *Request) error {
    validator := validasi.Baru()

    // Validasi dulu sebelum business logic
    if err := validator.Jumlah(req.Amount, "amount"); err != nil {
        return err
    }

    // Business logic...
}

// ❌ Bad - tidak ada validasi
func (s *Service) Create(req *Request) error {
    // Langsung business logic tanpa validasi
    s.db.Create(req)
}
```

### 2. Validasi Sesuai Urutan Logis

```go
// ✅ Good - validasi terstruktur
validator := validasi.Baru()

// 1. Validasi field wajib dulu
if err := validator.TeksWajib(req.Nama, "nama", 3, 255); err != nil {
    return nil, err
}

// 2. Validasi format
if err := validator.Email(req.Email); err != nil {
    return nil, err
}

// 3. Validasi field opsional
if err := validator.TeksOpsional(req.Catatan, "catatan", 500); err != nil {
    return nil, err
}
```

### 3. Gunakan Nama Field yang Jelas

```go
// ✅ Good - nama field deskriptif
validator.Jumlah(req.HargaBeli, "harga beli")
validator.Jumlah(req.HargaJual, "harga jual")

// ❌ Bad - nama field tidak jelas
validator.Jumlah(req.HargaBeli, "harga")  // Harga apa?
validator.Jumlah(req.HargaJual, "value")  // Tidak deskriptif
```

### 4. Validasi Conditional Fields

```go
// ✅ Good - validasi hanya jika field ada
if req.TanggalLahir != nil {
    if err := validator.TanggalLahir(*req.TanggalLahir); err != nil {
        return nil, err
    }
}

if req.Email != "" {
    if err := validator.Email(req.Email); err != nil {
        return nil, err
    }
}
```

### 5. Error Handling yang Konsisten

```go
// ✅ Good - return error langsung
if err := validator.Jumlah(amount, "total"); err != nil {
    return nil, err  // Error message sudah jelas dari validator
}

// ❌ Bad - wrap error tanpa info tambahan
if err := validator.Jumlah(amount, "total"); err != nil {
    return nil, fmt.Errorf("validation failed: %w", err)  // Redundant
}
```

## 📝 Error Messages

Semua error messages menggunakan bahasa Indonesia yang jelas dan user-friendly:

- ✅ "jumlah setoran harus lebih dari 0"
- ✅ "tanggal transaksi tidak boleh di masa depan"
- ✅ "email tidak valid"
- ✅ "nomor HP tidak valid (contoh: 08123456789)"
- ✅ "nama lengkap terlalu pendek (minimal 3 karakter)"

Error messages dirancang untuk:
1. Mudah dipahami user
2. Memberikan petunjuk cara perbaikan
3. Konsisten dengan terminologi bisnis
4. Sesuai konteks koperasi Indonesia

## 🔒 Security Benefits

Package validasi ini membantu mencegah:

1. **SQL Injection** - Validasi string length mencegah payload berlebihan
2. **Data Integrity Issues** - Validasi amount & date mencegah data tidak valid
3. **Business Logic Bypass** - Validasi enum & foreign keys mencegah manipulasi
4. **DoS Attacks** - Validasi length mencegah input berlebihan

## 📊 Performance

- **Minimal Overhead**: Validasi in-memory, tidak ada database call
- **No Dependencies**: Pure Go, tidak ada external library
- **Fast Execution**: Rata-rata < 1ms per validasi

## 🤝 Contributing

Untuk menambah fungsi validasi baru:

1. Tambahkan function di `validasi.go`
2. Gunakan naming bahasa Indonesia yang konsisten
3. Tambahkan unit tests di `validasi_test.go`
4. Update dokumentasi ini
5. Pastikan semua tests pass

## 📄 License

Bagian dari Cooperative ERP Lite project.

## 📞 Support

Untuk pertanyaan atau issue, silakan buat issue di GitHub repository.
