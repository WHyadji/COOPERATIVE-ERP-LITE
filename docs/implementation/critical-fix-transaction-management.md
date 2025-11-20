# CRITICAL FIX: Transaction Management Issues

**Status:** 🔴 CRITICAL - Must Fix Before Production
**Priority:** P0 (Highest)
**Estimated Effort:** 2-3 hours
**Impact:** Data Consistency & Accounting Integrity
**Assignee:** Backend Developer Team
**Due Date:** ASAP

---

## Executive Summary

Ditemukan **2 critical bugs** pada transaction management di modul **POS (Penjualan)** dan **Share Capital (Simpanan)** yang dapat menyebabkan:

1. **Data Inconsistency**: Transaksi tersimpan tanpa journal entry
2. **Accounting Errors**: Laporan keuangan tidak balance
3. **Orphan Records**: Data POS/Simpanan tanpa pasangan di accounting
4. **Manual Cleanup Required**: Butuh intervensi manual untuk perbaiki data

**Root Cause:**
Posting ke jurnal akuntansi dilakukan **DILUAR database transaction**, sehingga tidak atomic.

**Solution:**
Pindahkan semua operasi accounting posting **KE DALAM** transaction yang sama.

---

## Table of Contents

1. [Problem Analysis](#problem-analysis)
2. [Technical Details](#technical-details)
3. [Implementation Guide](#implementation-guide)
4. [Code Changes Required](#code-changes-required)
5. [Testing Checklist](#testing-checklist)
6. [Rollback Plan](#rollback-plan)
7. [Acceptance Criteria](#acceptance-criteria)

---

## Problem Analysis

### Current Architecture (BUGGY)

```
┌─────────────────────────────────────────────┐
│  POS/Simpanan Transaction Flow (CURRENT)   │
└─────────────────────────────────────────────┘

Step 1: Begin DB Transaction
  ├─ Create Sale/Simpanan Record     ✓
  ├─ Create Items                    ✓
  ├─ Update Stock                    ✓
  └─ COMMIT Transaction              ✓

Step 2: Post to Accounting (OUTSIDE Transaction) ❌
  ├─ Create Journal Entry
  └─ If FAILED:
      ├─ Sale/Simpanan already saved! 💥
      ├─ Data inconsistency!
      └─ Manual rollback unreliable
```

### Target Architecture (CORRECT)

```
┌─────────────────────────────────────────────┐
│  POS/Simpanan Transaction Flow (FIXED)     │
└─────────────────────────────────────────────┘

Step 1: Begin DB Transaction
  ├─ Create Sale/Simpanan Record     ✓
  ├─ Create Items                    ✓
  ├─ Update Stock                    ✓
  ├─ Create Journal Entry            ✓ (MOVED HERE!)
  └─ COMMIT ALL or ROLLBACK ALL      ✓

Result: ATOMIC - All or Nothing!
```

---

## Technical Details

### Bug #1: PenjualanService (Line 88-148)

**File:** `backend/internal/services/penjualan_service.go`

**Current Code (BUGGY):**

```go
// Line 88-136: Transaction for sale only
err = s.db.Transaction(func(tx *gorm.DB) error {
    // 1. Create sale record
    if err := tx.Create(&penjualan).Error; err != nil {
        return errors.New("gagal membuat penjualan")
    }

    // 2. Create items
    for _, itemReq := range req.Items {
        // ... create items ...
    }

    return nil  // ← Transaction ENDS here
})

// Line 142-147: Posting OUTSIDE transaction ❌
err = s.transaksiService.PostingOtomatisPenjualan(idKoperasi, idKasir, penjualan.ID)
if err != nil {
    // ⚠️ Sale already saved, but journal entry failed!
    // This creates orphan data!
    return nil, fmt.Errorf("penjualan berhasil, tetapi posting gagal: %w", err)
}
```

**Problem:**
- If `PostingOtomatisPenjualan()` fails due to:
  - Network timeout
  - Database connection lost
  - Accounting service error
  - Validation error in journal entry
- Then: Sale is saved, but journal is NOT created
- Result: POS shows sale, but accounting balance is wrong!

**Impact:**
- Trial Balance won't balance
- Income Statement incomplete
- Cash account incorrect
- Requires manual journal entry correction

---

### Bug #2: SimpananService (Line 95-106)

**File:** `backend/internal/services/simpanan_service.go`

**Current Code (BUGGY):**

```go
// Line 95-98: Create simpanan (separate transaction)
err = s.db.Create(simpanan).Error
if err != nil {
    return nil, errors.New("gagal mencatat setoran simpanan")
}

// Line 101-106: Posting OUTSIDE, with manual rollback ❌
err = s.transaksiService.PostingOtomatisSimpanan(idKoperasi, idPengguna, simpanan.ID)
if err != nil {
    // Manual rollback - UNRELIABLE!
    s.db.Delete(simpanan)  // ← This can ALSO fail!
    return nil, fmt.Errorf("gagal posting ke jurnal: %w", err)
}
```

**Problem:**
- Manual rollback with `s.db.Delete()` is NOT atomic
- If `Delete()` fails, simpanan stays in database without journal
- Race condition: Another request might read simpanan before delete
- No transaction = no guarantee

**Impact:**
- Share Capital balance incorrect
- Member balance report wrong
- Liability account (Simpanan Anggota) incorrect

---

## Implementation Guide

### Step 1: Create New Methods in TransaksiService

**File to Edit:** `backend/internal/services/transaksi_service.go`

**What to Add:** 2 new methods that accept existing transaction

```go
// Add these methods to TransaksiService

// PostingOtomatisPenjualanWithTx posts sale to journal using existing transaction
// This ensures atomicity with the sale creation
func (s *TransaksiService) PostingOtomatisPenjualanWithTx(
    tx *gorm.DB,
    idKoperasi uuid.UUID,
    idKasir uuid.UUID,
    idPenjualan uuid.UUID,
) error {
    // 1. Get sale data
    var penjualan models.Penjualan
    if err := tx.Where("id = ?", idPenjualan).First(&penjualan).Error; err != nil {
        return errors.New("penjualan tidak ditemukan")
    }

    // 2. Get account IDs (Kas & Pendapatan Penjualan)
    var akunKas, akunPendapatan models.Akun

    // Kas (Cash) - usually code 1-1-01 or similar
    if err := tx.Where("id_koperasi = ? AND kode_akun LIKE ?", idKoperasi, "1-1%").
        Where("nama_akun ILIKE ?", "%kas%").
        First(&akunKas).Error; err != nil {
        return errors.New("akun kas tidak ditemukan")
    }

    // Pendapatan Penjualan - usually code 4-1-01 or similar
    if err := tx.Where("id_koperasi = ? AND kode_akun LIKE ?", idKoperasi, "4-1%").
        Where("nama_akun ILIKE ?", "%penjualan%").
        First(&akunPendapatan).Error; err != nil {
        return errors.New("akun pendapatan penjualan tidak ditemukan")
    }

    // 3. Create journal entry
    transaksi := &models.Transaksi{
        IDKoperasi:  idKoperasi,
        Tanggal:     penjualan.TanggalPenjualan,
        Deskripsi:   fmt.Sprintf("Penjualan %s", penjualan.NomorPenjualan),
        TotalDebit:  penjualan.TotalBelanja,
        TotalKredit: penjualan.TotalBelanja,
        DibuatOleh:  idKasir,
    }

    if err := tx.Create(transaksi).Error; err != nil {
        return errors.New("gagal membuat jurnal penjualan")
    }

    // 4. Create debit line (Kas)
    debitLine := &models.TransaksiLine{
        IDTransaksi: transaksi.ID,
        IDAkun:      akunKas.ID,
        Debit:       penjualan.TotalBelanja,
        Kredit:      0,
        Deskripsi:   fmt.Sprintf("Penerimaan kas dari %s", penjualan.NomorPenjualan),
    }

    if err := tx.Create(debitLine).Error; err != nil {
        return errors.New("gagal membuat debit line")
    }

    // 5. Create credit line (Pendapatan)
    kreditLine := &models.TransaksiLine{
        IDTransaksi: transaksi.ID,
        IDAkun:      akunPendapatan.ID,
        Debit:       0,
        Kredit:      penjualan.TotalBelanja,
        Deskripsi:   fmt.Sprintf("Pendapatan dari %s", penjualan.NomorPenjualan),
    }

    if err := tx.Create(kreditLine).Error; err != nil {
        return errors.New("gagal membuat kredit line")
    }

    return nil
}

// PostingOtomatisSimpananWithTx posts share capital to journal using existing transaction
func (s *TransaksiService) PostingOtomatisSimpananWithTx(
    tx *gorm.DB,
    idKoperasi uuid.UUID,
    idPengguna uuid.UUID,
    idSimpanan uuid.UUID,
) error {
    // 1. Get simpanan data
    var simpanan models.Simpanan
    if err := tx.Where("id = ?", idSimpanan).Preload("Anggota").First(&simpanan).Error; err != nil {
        return errors.New("simpanan tidak ditemukan")
    }

    // 2. Get account IDs (Kas & Simpanan Anggota)
    var akunKas, akunSimpanan models.Akun

    // Kas (Cash)
    if err := tx.Where("id_koperasi = ? AND kode_akun LIKE ?", idKoperasi, "1-1%").
        Where("nama_akun ILIKE ?", "%kas%").
        First(&akunKas).Error; err != nil {
        return errors.New("akun kas tidak ditemukan")
    }

    // Simpanan Anggota (Liability) - code depends on type
    var kodeAkunPattern string
    switch simpanan.TipeSimpanan {
    case models.SimpananPokok:
        kodeAkunPattern = "%pokok%"
    case models.SimpananWajib:
        kodeAkunPattern = "%wajib%"
    case models.SimpananSukarela:
        kodeAkunPattern = "%sukarela%"
    default:
        return errors.New("tipe simpanan tidak valid")
    }

    if err := tx.Where("id_koperasi = ? AND kode_akun LIKE ?", idKoperasi, "2-%").
        Where("nama_akun ILIKE ?", kodeAkunPattern).
        First(&akunSimpanan).Error; err != nil {
        return errors.New("akun simpanan tidak ditemukan")
    }

    // 3. Create journal entry
    transaksi := &models.Transaksi{
        IDKoperasi:  idKoperasi,
        Tanggal:     simpanan.TanggalTransaksi,
        Deskripsi:   fmt.Sprintf("Setoran %s - %s (%s)", simpanan.TipeSimpanan, simpanan.Anggota.NamaLengkap, simpanan.NomorReferensi),
        TotalDebit:  simpanan.JumlahSetoran,
        TotalKredit: simpanan.JumlahSetoran,
        DibuatOleh:  idPengguna,
    }

    if err := tx.Create(transaksi).Error; err != nil {
        return errors.New("gagal membuat jurnal simpanan")
    }

    // 4. Create debit line (Kas)
    debitLine := &models.TransaksiLine{
        IDTransaksi: transaksi.ID,
        IDAkun:      akunKas.ID,
        Debit:       simpanan.JumlahSetoran,
        Kredit:      0,
        Deskripsi:   fmt.Sprintf("Penerimaan setoran %s", simpanan.NomorReferensi),
    }

    if err := tx.Create(debitLine).Error; err != nil {
        return errors.New("gagal membuat debit line")
    }

    // 5. Create credit line (Simpanan Anggota)
    kreditLine := &models.TransaksiLine{
        IDTransaksi: transaksi.ID,
        IDAkun:      akunSimpanan.ID,
        Debit:       0,
        Kredit:      simpanan.JumlahSetoran,
        Deskripsi:   fmt.Sprintf("Simpanan a.n. %s", simpanan.Anggota.NamaLengkap),
    }

    if err := tx.Create(kreditLine).Error; err != nil {
        return errors.New("gagal membuat kredit line")
    }

    return nil
}
```

**Important Notes:**
- These methods use `tx *gorm.DB` parameter instead of `s.db`
- This allows them to participate in existing transaction
- Account lookup uses pattern matching - adjust based on your Chart of Accounts
- Error messages are descriptive for debugging

---

### Step 2: Update PenjualanService

**File to Edit:** `backend/internal/services/penjualan_service.go`

**Method to Fix:** `ProsesPenjualan()` (Line 47-154)

**Find this code block (Line 88-148):**

```go
// OLD CODE - DELETE THIS
err = s.db.Transaction(func(tx *gorm.DB) error {
    // 1. Buat record penjualan
    penjualan = models.Penjualan{
        IDKoperasi:       idKoperasi,
        NomorPenjualan:   nomorPenjualan,
        TanggalPenjualan: time.Now(),
        IDAnggota:        req.IDAnggota,
        TotalBelanja:     totalBelanja,
        MetodePembayaran: models.PembayaranTunai,
        JumlahBayar:      req.JumlahBayar,
        Kembalian:        kembalian,
        IDKasir:          idKasir,
        Catatan:          req.Catatan,
    }

    if err := tx.Create(&penjualan).Error; err != nil {
        return errors.New("gagal membuat penjualan")
    }

    // 2. Buat item penjualan dan kurangi stok
    for _, itemReq := range req.Items {
        // ... existing code ...
    }

    return nil
})

if err != nil {
    return nil, err
}

// 3. Auto-posting ke jurnal akuntansi
err = s.transaksiService.PostingOtomatisPenjualan(idKoperasi, idKasir, penjualan.ID)
if err != nil {
    return nil, fmt.Errorf("penjualan berhasil, tetapi posting gagal: %w", err)
}
```

**Replace with (NEW CODE):**

```go
// NEW CODE - REPLACE ENTIRE TRANSACTION BLOCK
err = s.db.Transaction(func(tx *gorm.DB) error {
    // 1. Buat record penjualan
    penjualan = models.Penjualan{
        IDKoperasi:       idKoperasi,
        NomorPenjualan:   nomorPenjualan,
        TanggalPenjualan: time.Now(),
        IDAnggota:        req.IDAnggota,
        TotalBelanja:     totalBelanja,
        MetodePembayaran: models.PembayaranTunai,
        JumlahBayar:      req.JumlahBayar,
        Kembalian:        kembalian,
        IDKasir:          idKasir,
        Catatan:          req.Catatan,
    }

    if err := tx.Create(&penjualan).Error; err != nil {
        return errors.New("gagal membuat penjualan")
    }

    // 2. Buat item penjualan dan kurangi stok
    for _, itemReq := range req.Items {
        // Dapatkan produk untuk nama
        var produk models.Produk
        if err := tx.Where("id = ?", itemReq.IDProduk).First(&produk).Error; err != nil {
            return fmt.Errorf("produk %s tidak ditemukan", itemReq.IDProduk)
        }

        // Buat item penjualan
        item := models.ItemPenjualan{
            IDPenjualan: penjualan.ID,
            IDProduk:    itemReq.IDProduk,
            NamaProduk:  produk.NamaProduk,
            Kuantitas:   itemReq.Kuantitas,
            HargaSatuan: itemReq.HargaSatuan,
        }

        if err := tx.Create(&item).Error; err != nil {
            return errors.New("gagal membuat item penjualan")
        }

        // Kurangi stok produk
        if err := s.produkService.KurangiStok(itemReq.IDProduk, itemReq.Kuantitas); err != nil {
            return fmt.Errorf("gagal mengurangi stok: %w", err)
        }
    }

    // 3. Auto-posting ke jurnal akuntansi (MOVED INSIDE TRANSACTION!)
    if err := s.transaksiService.PostingOtomatisPenjualanWithTx(tx, idKoperasi, idKasir, penjualan.ID); err != nil {
        // This will rollback EVERYTHING: sale + items + stock updates
        return fmt.Errorf("gagal posting ke jurnal: %w", err)
    }

    return nil  // Commit ALL: sale + items + journal entry
})

if err != nil {
    return nil, err
}

// No separate posting needed - it's already done inside transaction!
```

**Key Changes:**
1. ✅ Moved `PostingOtomatisPenjualan` call INSIDE transaction
2. ✅ Changed to `PostingOtomatisPenjualanWithTx(tx, ...)` - passes transaction object
3. ✅ Removed separate posting call after transaction
4. ✅ Error in posting = auto rollback ALL changes

---

### Step 3: Update SimpananService

**File to Edit:** `backend/internal/services/simpanan_service.go`

**Method to Fix:** `CatatSetoran()` (Line 39-113)

**Find this code block (Line 83-106):**

```go
// OLD CODE - DELETE THIS
simpanan := &models.Simpanan{
    IDKoperasi:       idKoperasi,
    IDAnggota:        req.IDAnggota,
    TipeSimpanan:     req.TipeSimpanan,
    TanggalTransaksi: req.TanggalTransaksi,
    JumlahSetoran:    req.JumlahSetoran,
    Keterangan:       req.Keterangan,
    NomorReferensi:   nomorReferensi,
    DibuatOleh:       idPengguna,
}

// Simpan ke database
err = s.db.Create(simpanan).Error
if err != nil {
    return nil, errors.New("gagal mencatat setoran simpanan")
}

// Auto-posting ke jurnal akuntansi
err = s.transaksiService.PostingOtomatisSimpanan(idKoperasi, idPengguna, simpanan.ID)
if err != nil {
    // Rollback simpanan jika posting gagal
    s.db.Delete(simpanan)
    return nil, fmt.Errorf("gagal posting ke jurnal: %w", err)
}
```

**Replace with (NEW CODE):**

```go
// NEW CODE - REPLACE WITH TRANSACTION BLOCK
simpanan := &models.Simpanan{
    IDKoperasi:       idKoperasi,
    IDAnggota:        req.IDAnggota,
    TipeSimpanan:     req.TipeSimpanan,
    TanggalTransaksi: req.TanggalTransaksi,
    JumlahSetoran:    req.JumlahSetoran,
    Keterangan:       req.Keterangan,
    NomorReferensi:   nomorReferensi,
    DibuatOleh:       idPengguna,
}

// Simpan dalam transaction bersama posting
err = s.db.Transaction(func(tx *gorm.DB) error {
    // 1. Create simpanan record
    if err := tx.Create(simpanan).Error; err != nil {
        return errors.New("gagal mencatat setoran simpanan")
    }

    // 2. Auto-posting ke jurnal akuntansi (INSIDE TRANSACTION!)
    if err := s.transaksiService.PostingOtomatisSimpananWithTx(tx, idKoperasi, idPengguna, simpanan.ID); err != nil {
        // This will auto-rollback simpanan creation
        return fmt.Errorf("gagal posting ke jurnal: %w", err)
    }

    return nil  // Commit both: simpanan + journal entry
})

if err != nil {
    return nil, err
}

// No manual rollback needed - GORM handles it automatically!
```

**Key Changes:**
1. ✅ Wrapped both operations in `s.db.Transaction()`
2. ✅ Changed to `PostingOtomatisSimpananWithTx(tx, ...)` - passes transaction
3. ✅ Removed manual `s.db.Delete(simpanan)` - no longer needed!
4. ✅ GORM automatically rolls back on error

---

## Code Changes Required

### Summary Checklist

- [ ] **File 1:** `backend/internal/services/transaksi_service.go`
  - [ ] Add method: `PostingOtomatisPenjualanWithTx(tx *gorm.DB, ...) error`
  - [ ] Add method: `PostingOtomatisSimpananWithTx(tx *gorm.DB, ...) error`
  - [ ] Estimated lines: ~200 lines added

- [ ] **File 2:** `backend/internal/services/penjualan_service.go`
  - [ ] Modify method: `ProsesPenjualan()` line 88-148
  - [ ] Move posting inside transaction
  - [ ] Remove separate posting call
  - [ ] Estimated lines: ~5 lines changed

- [ ] **File 3:** `backend/internal/services/simpanan_service.go`
  - [ ] Modify method: `CatatSetoran()` line 83-106
  - [ ] Wrap in transaction
  - [ ] Remove manual rollback
  - [ ] Estimated lines: ~10 lines changed

**Total Changes:** 3 files, ~215 lines added/modified

---

## Testing Checklist

### Unit Tests

Create test files to verify transaction behavior:

#### Test 1: PenjualanService - Posting Failure Rolls Back Sale

**File:** `backend/internal/services/penjualan_service_transaction_test.go`

```go
package services

import (
    "testing"
    "cooperative-erp-lite/internal/models"
    "github.com/stretchr/testify/assert"
)

func TestProsesPenjualan_PostingFails_RollbackSale(t *testing.T) {
    // Setup test database
    db := setupTestDB()

    // Initialize services
    transaksiService := NewTransaksiService(db)
    produkService := NewProdukService(db)
    penjualanService := NewPenjualanService(db, produkService, transaksiService)

    // Setup test data: cooperative, user, product
    koperasi := createTestKoperasi(db)
    kasir := createTestUser(db, koperasi.ID)
    produk := createTestProduk(db, koperasi.ID, 100) // stock: 100

    // IMPORTANT: Delete Kas account to force posting failure
    db.Where("id_koperasi = ? AND nama_akun ILIKE ?", koperasi.ID, "%kas%").Delete(&models.Akun{})

    // Create sale request
    req := &ProsesPenjualanRequest{
        Items: []ItemPenjualanRequest{
            {
                IDProduk:    produk.ID,
                Kuantitas:   10,
                HargaSatuan: 10000,
            },
        },
        JumlahBayar: 100000,
    }

    // Execute
    _, err := penjualanService.ProsesPenjualan(koperasi.ID, kasir.ID, req)

    // Assertions
    assert.Error(t, err, "Should return error when posting fails")
    assert.Contains(t, err.Error(), "akun kas tidak ditemukan")

    // CRITICAL: Verify sale was NOT created (rollback successful)
    var count int64
    db.Model(&models.Penjualan{}).Where("id_koperasi = ?", koperasi.ID).Count(&count)
    assert.Equal(t, int64(0), count, "Sale should NOT be saved when posting fails")

    // Verify stock was NOT reduced (rollback successful)
    var updatedProduk models.Produk
    db.First(&updatedProduk, produk.ID)
    assert.Equal(t, 100, updatedProduk.Stok, "Stock should NOT be reduced when posting fails")
}

func TestProsesPenjualan_Success_AllCommitted(t *testing.T) {
    // Setup
    db := setupTestDB()
    transaksiService := NewTransaksiService(db)
    produkService := NewProdukService(db)
    penjualanService := NewPenjualanService(db, produkService, transaksiService)

    // Setup test data WITH proper accounts
    koperasi := createTestKoperasi(db)
    kasir := createTestUser(db, koperasi.ID)
    produk := createTestProduk(db, koperasi.ID, 100)

    // Create required accounts
    akunKas := createTestAkun(db, koperasi.ID, "1-1-01", "Kas")
    akunPendapatan := createTestAkun(db, koperasi.ID, "4-1-01", "Pendapatan Penjualan")

    // Create sale
    req := &ProsesPenjualanRequest{
        Items: []ItemPenjualanRequest{
            {IDProduk: produk.ID, Kuantitas: 10, HargaSatuan: 10000},
        },
        JumlahBayar: 100000,
    }

    // Execute
    result, err := penjualanService.ProsesPenjualan(koperasi.ID, kasir.ID, req)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)

    // Verify sale created
    var sale models.Penjualan
    db.First(&sale, result.ID)
    assert.Equal(t, float64(100000), sale.TotalBelanja)

    // Verify journal entry created
    var transaksi models.Transaksi
    err = db.Where("deskripsi LIKE ?", "%Penjualan%").First(&transaksi).Error
    assert.NoError(t, err)
    assert.Equal(t, float64(100000), transaksi.TotalDebit)

    // Verify journal lines created (2 lines: debit + credit)
    var lines []models.TransaksiLine
    db.Where("id_transaksi = ?", transaksi.ID).Find(&lines)
    assert.Equal(t, 2, len(lines), "Should have 2 journal lines")

    // Verify stock reduced
    var updatedProduk models.Produk
    db.First(&updatedProduk, produk.ID)
    assert.Equal(t, 90, updatedProduk.Stok, "Stock should be reduced by 10")
}
```

#### Test 2: SimpananService - Transaction Atomicity

**File:** `backend/internal/services/simpanan_service_transaction_test.go`

```go
package services

import (
    "testing"
    "time"
    "cooperative-erp-lite/internal/models"
    "github.com/stretchr/testify/assert"
)

func TestCatatSetoran_PostingFails_RollbackSimpanan(t *testing.T) {
    // Setup
    db := setupTestDB()
    transaksiService := NewTransaksiService(db)
    simpananService := NewSimpananService(db, transaksiService)

    // Setup test data
    koperasi := createTestKoperasi(db)
    pengguna := createTestUser(db, koperasi.ID)
    anggota := createTestAnggota(db, koperasi.ID)

    // IMPORTANT: Delete Kas account to force posting failure
    db.Where("id_koperasi = ? AND nama_akun ILIKE ?", koperasi.ID, "%kas%").Delete(&models.Akun{})

    // Create simpanan request
    req := &CatatSetoranRequest{
        IDAnggota:        anggota.ID,
        TipeSimpanan:     models.SimpananPokok,
        TanggalTransaksi: time.Now(),
        JumlahSetoran:    100000,
        Keterangan:       "Simpanan Pokok",
    }

    // Execute
    _, err := simpananService.CatatSetoran(koperasi.ID, pengguna.ID, req)

    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "akun kas tidak ditemukan")

    // CRITICAL: Verify simpanan was NOT created
    var count int64
    db.Model(&models.Simpanan{}).Where("id_anggota = ?", anggota.ID).Count(&count)
    assert.Equal(t, int64(0), count, "Simpanan should NOT be saved when posting fails")
}

func TestCatatSetoran_Success_AllCommitted(t *testing.T) {
    // Setup
    db := setupTestDB()
    transaksiService := NewTransaksiService(db)
    simpananService := NewSimpananService(db, transaksiService)

    // Setup test data WITH accounts
    koperasi := createTestKoperasi(db)
    pengguna := createTestUser(db, koperasi.ID)
    anggota := createTestAnggota(db, koperasi.ID)

    // Create accounts
    akunKas := createTestAkun(db, koperasi.ID, "1-1-01", "Kas")
    akunSimpananPokok := createTestAkun(db, koperasi.ID, "2-1-01", "Simpanan Pokok")

    // Create simpanan
    req := &CatatSetoranRequest{
        IDAnggota:        anggota.ID,
        TipeSimpanan:     models.SimpananPokok,
        TanggalTransaksi: time.Now(),
        JumlahSetoran:    100000,
    }

    // Execute
    result, err := simpananService.CatatSetoran(koperasi.ID, pengguna.ID, req)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)

    // Verify simpanan created
    var simpanan models.Simpanan
    db.First(&simpanan, result.ID)
    assert.Equal(t, float64(100000), simpanan.JumlahSetoran)

    // Verify journal created
    var transaksi models.Transaksi
    err = db.Where("deskripsi LIKE ?", "%Simpanan%").First(&transaksi).Error
    assert.NoError(t, err)

    // Verify journal lines
    var lines []models.TransaksiLine
    db.Where("id_transaksi = ?", transaksi.ID).Find(&lines)
    assert.Equal(t, 2, len(lines))
}
```

### Integration Tests

#### Test 3: End-to-End POS Flow

```bash
# Manual test in Postman/curl

# 1. Create product
POST /api/v1/produk
{
  "namaProduk": "Test Product",
  "hargaJual": 10000,
  "stok": 100
}

# 2. Verify accounts exist
GET /api/v1/akun
# Should return Kas and Pendapatan accounts

# 3. Create sale
POST /api/v1/penjualan
{
  "items": [
    {
      "idProduk": "{{productId}}",
      "kuantitas": 10,
      "hargaSatuan": 10000
    }
  ],
  "jumlahBayar": 100000
}

# 4. Verify sale created
GET /api/v1/penjualan

# 5. Verify journal created
GET /api/v1/transaksi
# Should show journal entry for the sale

# 6. Verify stock reduced
GET /api/v1/produk/{{productId}}
# Stock should be 90 (100 - 10)
```

#### Test 4: Failure Scenario Test

```bash
# Simulate posting failure

# 1. Delete Kas account (to force error)
DELETE /api/v1/akun/{{kasAccountId}}

# 2. Try to create sale
POST /api/v1/penjualan
{
  "items": [...],
  "jumlahBayar": 100000
}

# Expected: 500 error with message "akun kas tidak ditemukan"

# 3. Verify NO sale was created (rollback worked)
GET /api/v1/penjualan
# Should return empty or previous sales only

# 4. Verify stock NOT reduced
GET /api/v1/produk/{{productId}}
# Stock should still be 100
```

### Performance Tests

```go
func BenchmarkProsesPenjualan(b *testing.B) {
    // Setup
    db := setupTestDB()
    // ... setup services and test data ...

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := &ProsesPenjualanRequest{...}
        _, err := penjualanService.ProsesPenjualan(koperasi.ID, kasir.ID, req)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// Target: < 100ms per transaction (including journal posting)
```

---

## Rollback Plan

If issues arise after deployment:

### Quick Rollback (Emergency)

```bash
# 1. Revert to previous git commit
git revert <commit-hash>
git push origin main

# 2. Redeploy
fly deploy  # or your deployment command

# Time: ~5 minutes
```

### Data Cleanup (If orphan records exist)

```sql
-- Find sales without journal entries
SELECT p.id, p.nomor_penjualan, p.total_belanja
FROM penjualan p
LEFT JOIN transaksi t ON t.deskripsi LIKE '%' || p.nomor_penjualan || '%'
WHERE t.id IS NULL;

-- Find simpanan without journal entries
SELECT s.id, s.nomor_referensi, s.jumlah_setoran
FROM simpanan s
LEFT JOIN transaksi t ON t.deskripsi LIKE '%' || s.nomor_referensi || '%'
WHERE t.id IS NULL;

-- Manually create missing journal entries or delete orphan records
-- (Consult with accounting team before executing!)
```

---

## Acceptance Criteria

### Definition of Done

- [ ] All 3 files modified successfully
- [ ] Code compiles without errors: `go build ./...`
- [ ] All existing tests pass: `go test ./...`
- [ ] New transaction tests added and passing
- [ ] Manual testing completed (checklist below)
- [ ] Code reviewed by senior developer
- [ ] Deployed to staging environment
- [ ] Smoke tests passed on staging
- [ ] Approved by QA team

### Manual Testing Checklist

**Scenario 1: POS - Successful Sale**
- [ ] Create product with stock > 0
- [ ] Create sale via API
- [ ] Verify sale record created
- [ ] Verify journal entry created
- [ ] Verify stock reduced
- [ ] Verify trial balance still balances

**Scenario 2: POS - Account Missing (Failure)**
- [ ] Delete Kas or Pendapatan account
- [ ] Try to create sale
- [ ] Verify error returned
- [ ] Verify NO sale record created
- [ ] Verify stock NOT reduced
- [ ] Restore accounts for next tests

**Scenario 3: Simpanan - Successful Deposit**
- [ ] Create member
- [ ] Create simpanan via API
- [ ] Verify simpanan record created
- [ ] Verify journal entry created
- [ ] Verify member balance updated

**Scenario 4: Simpanan - Posting Failure**
- [ ] Delete Simpanan account
- [ ] Try to create simpanan
- [ ] Verify error returned
- [ ] Verify NO simpanan record created
- [ ] Restore accounts

**Scenario 5: Concurrent Transactions**
- [ ] Run 10 simultaneous sale requests
- [ ] Verify all 10 completed successfully OR all failed
- [ ] Verify NO partial transactions
- [ ] Verify trial balance still balances

### Success Metrics

- ✅ 100% of sales have matching journal entries
- ✅ 100% of simpanan have matching journal entries
- ✅ 0 orphan records in database
- ✅ Trial balance balances after all transactions
- ✅ Transaction time < 100ms (p95)
- ✅ Zero data inconsistency errors in logs

---

## Additional Notes

### Common Pitfalls to Avoid

1. **Don't forget to pass `tx` parameter**
   - Wrong: `s.transaksiService.PostingOtomatisPenjualan(...)`
   - Correct: `s.transaksiService.PostingOtomatisPenjualanWithTx(tx, ...)`

2. **Don't use `s.db` inside transaction**
   - Wrong: `if err := s.db.Create(...)`
   - Correct: `if err := tx.Create(...)`

3. **Don't call methods that start their own transaction**
   - Check if `KurangiStok()` uses `s.db.Transaction()` internally
   - If yes, refactor to accept `tx` parameter

4. **Account lookup patterns**
   - Adjust SQL patterns (`LIKE '1-1%'`, `ILIKE '%kas%'`) to match your Chart of Accounts
   - Coordinate with accounting team for exact account codes

### Questions & Answers

**Q: Will this slow down transactions?**
A: No, minimal impact. All operations already execute sequentially. We're just wrapping them in a transaction block.

**Q: What if `PostingOtomatisPenjualan` already exists?**
A: Keep the old method for backward compatibility. Add new `WithTx` variant. You can deprecate old method later.

**Q: How to handle existing orphan records?**
A: Run data cleanup script (see Rollback Plan section) or create compensating journal entries.

**Q: Need to update ProdukService.KurangiStok()?**
A: Check the implementation. If it starts its own transaction, create `KurangiStokWithTx(tx, ...)` variant.

---

## Timeline & Milestones

| Day | Task | Owner | Status |
|-----|------|-------|--------|
| Day 1 AM | Create `WithTx` methods in TransaksiService | Backend Dev | ⏳ Pending |
| Day 1 PM | Update PenjualanService | Backend Dev | ⏳ Pending |
| Day 2 AM | Update SimpananService | Backend Dev | ⏳ Pending |
| Day 2 PM | Write unit tests | Backend Dev | ⏳ Pending |
| Day 3 AM | Integration testing | QA | ⏳ Pending |
| Day 3 PM | Deploy to staging | DevOps | ⏳ Pending |
| Day 4 | Final review & production deploy | Tech Lead | ⏳ Pending |

**Target Completion:** 4 working days from start

---

## Support & Escalation

**If you encounter issues:**

1. **Technical questions**: Contact Backend Tech Lead
2. **Accounting logic questions**: Contact Finance/Accounting Team
3. **Urgent blockers**: Escalate to Project Manager
4. **Database issues**: Contact DBA/DevOps

**Documentation references:**
- GORM Transactions: https://gorm.io/docs/transactions.html
- Go UUID package: https://github.com/google/uuid
- Project Chart of Accounts: `docs/chart-of-accounts.md`

---

**Document Version:** 1.0
**Last Updated:** 2025-01-19
**Author:** Technical Architect
**Reviewers:** Backend Lead, QA Lead, Product Manager
