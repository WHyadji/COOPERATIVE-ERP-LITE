package services

import (
	"cooperative-erp-lite/internal/models"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestPerbaruiPengguna_CrossTenantBlocked tests that cross-tenant updates are blocked
func TestPerbaruiPengguna_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	// Setup two separate cooperatives
	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	// Create user in Koperasi A
	penggunaA := &models.Pengguna{
		IDKoperasi:   koperasiA,
		NamaLengkap:  "User A",
		NamaPengguna: "userA",
		Email:        "usera@test.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	penggunaA.SetKataSandi("password123")
	db.Create(penggunaA)

	service := NewPenggunaService(db)

	// Attempt to update User A using Koperasi B's context
	req := &PerbaruiPenggunaRequest{
		NamaLengkap: "Hacked Name",
	}

	_, err := service.PerbaruiPengguna(koperasiB, penggunaA.ID, req)

	// Should fail - not found or no access
	if err == nil {
		t.Fatal("Expected error when updating cross-tenant user, got nil")
	}

	if !strings.Contains(err.Error(), "tidak ditemukan atau tidak memiliki akses") {
		t.Errorf("Expected 'tidak ditemukan atau tidak memiliki akses' error, got: %v", err)
	}

	// Verify data unchanged
	var unchanged models.Pengguna
	db.First(&unchanged, penggunaA.ID)
	if unchanged.NamaLengkap != "User A" {
		t.Errorf("Data was modified! Expected 'User A', got '%s'", unchanged.NamaLengkap)
	}

	t.Logf("✓ Cross-tenant update blocked successfully")
}

// TestPerbaruiPengguna_SameTenantAllowed tests that same-tenant updates work
func TestPerbaruiPengguna_SameTenantAllowed(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koop := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	db.Create(koop)
	defer cleanupTestData(db, koperasiA)

	pengguna := &models.Pengguna{
		IDKoperasi:   koperasiA,
		NamaLengkap:  "User A",
		NamaPengguna: "userA",
		Email:        "usera@test.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	pengguna.SetKataSandi("password123")
	db.Create(pengguna)

	service := NewPenggunaService(db)

	// Update using correct cooperative ID
	req := &PerbaruiPenggunaRequest{
		NamaLengkap: "Updated Name",
	}

	result, err := service.PerbaruiPengguna(koperasiA, pengguna.ID, req)

	if err != nil {
		t.Fatalf("Expected success for same-tenant update, got error: %v", err)
	}

	if result.NamaLengkap != "Updated Name" {
		t.Errorf("Expected 'Updated Name', got '%s'", result.NamaLengkap)
	}

	t.Logf("✓ Same-tenant update allowed successfully")
}

// TestHapusPengguna_CrossTenantBlocked tests that cross-tenant deletes are blocked
func TestHapusPengguna_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	penggunaA := &models.Pengguna{
		IDKoperasi:   koperasiA,
		NamaLengkap:  "User A",
		NamaPengguna: "userA",
		Email:        "usera@test.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	penggunaA.SetKataSandi("password123")
	db.Create(penggunaA)

	service := NewPenggunaService(db)

	// Attempt to delete User A using Koperasi B's context
	err := service.HapusPengguna(koperasiB, penggunaA.ID)

	if err == nil {
		t.Fatal("Expected error when deleting cross-tenant user, got nil")
	}

	if !strings.Contains(err.Error(), "tidak ditemukan atau tidak memiliki akses") {
		t.Errorf("Expected 'tidak ditemukan atau tidak memiliki akses' error, got: %v", err)
	}

	// Verify user still exists
	var stillExists models.Pengguna
	result := db.First(&stillExists, penggunaA.ID)
	if result.Error != nil {
		t.Errorf("User was deleted! Should still exist")
	}

	t.Logf("✓ Cross-tenant delete blocked successfully")
}

// TestPerbaruiAnggota_CrossTenantBlocked tests anggota service multi-tenant validation
func TestPerbaruiAnggota_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	anggotaA := &models.Anggota{
		IDKoperasi:       koperasiA,
		NomorAnggota:     "KOOP-2025-0001",
		NamaLengkap:      "Anggota A",
		TanggalBergabung: time.Now(),
		Status:           models.StatusAktif,
	}
	db.Create(anggotaA)

	service := NewAnggotaService(db)

	req := &PerbaruiAnggotaRequest{
		NamaLengkap: "Hacked Name",
	}

	_, err := service.PerbaruiAnggota(koperasiB, anggotaA.ID, req)

	if err == nil {
		t.Fatal("Expected error when updating cross-tenant anggota, got nil")
	}

	if !strings.Contains(err.Error(), "tidak ditemukan atau tidak memiliki akses") {
		t.Errorf("Expected 'tidak ditemukan atau tidak memiliki akses' error, got: %v", err)
	}

	var unchanged models.Anggota
	db.First(&unchanged, anggotaA.ID)
	if unchanged.NamaLengkap != "Anggota A" {
		t.Errorf("Data was modified! Expected 'Anggota A', got '%s'", unchanged.NamaLengkap)
	}

	t.Logf("✓ Anggota cross-tenant update blocked successfully")
}

// TestHapusAnggota_CrossTenantBlocked tests anggota delete multi-tenant validation
func TestHapusAnggota_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	anggotaA := &models.Anggota{
		IDKoperasi:       koperasiA,
		NomorAnggota:     "KOOP-2025-0001",
		NamaLengkap:      "Anggota A",
		TanggalBergabung: time.Now(),
		Status:           models.StatusAktif,
	}
	db.Create(anggotaA)

	service := NewAnggotaService(db)

	err := service.HapusAnggota(koperasiB, anggotaA.ID)

	if err == nil {
		t.Fatal("Expected error when deleting cross-tenant anggota, got nil")
	}

	var stillExists models.Anggota
	result := db.First(&stillExists, anggotaA.ID)
	if result.Error != nil {
		t.Errorf("Anggota was deleted! Should still exist")
	}

	t.Logf("✓ Anggota cross-tenant delete blocked successfully")
}

// TestSetPINPortal_CrossTenantBlocked tests PIN setting multi-tenant validation
func TestSetPINPortal_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	anggotaA := &models.Anggota{
		IDKoperasi:       koperasiA,
		NomorAnggota:     "KOOP-2025-0001",
		NamaLengkap:      "Anggota A",
		TanggalBergabung: time.Now(),
		Status:           models.StatusAktif,
	}
	db.Create(anggotaA)

	service := NewAnggotaService(db)

	err := service.SetPINPortal(koperasiB, anggotaA.ID, "1234")

	if err == nil {
		t.Fatal("Expected error when setting PIN for cross-tenant anggota, got nil")
	}

	var unchanged models.Anggota
	db.First(&unchanged, anggotaA.ID)
	if unchanged.PINPortal != "" {
		t.Errorf("PIN was set! Should remain empty")
	}

	t.Logf("✓ SetPINPortal cross-tenant blocked successfully")
}

// TestPerbaruiAkun_CrossTenantBlocked tests akun update multi-tenant validation
func TestPerbaruiAkun_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	akunA := &models.Akun{
		IDKoperasi:  koperasiA,
		KodeAkun:    "1101",
		NamaAkun:    "Kas",
		TipeAkun:    models.AkunAset,
		NormalSaldo: "debit",
		StatusAktif: true,
	}
	db.Create(akunA)

	service := NewAkunService(db)

	req := &PerbaruiAkunRequest{
		NamaAkun: "Hacked Account",
	}

	_, err := service.PerbaruiAkun(koperasiB, akunA.ID, req)

	if err == nil {
		t.Fatal("Expected error when updating cross-tenant akun, got nil")
	}

	var unchanged models.Akun
	db.First(&unchanged, akunA.ID)
	if unchanged.NamaAkun != "Kas" {
		t.Errorf("Data was modified! Expected 'Kas', got '%s'", unchanged.NamaAkun)
	}

	t.Logf("✓ Akun cross-tenant update blocked successfully")
}

// TestHapusAkun_CrossTenantBlocked tests akun delete multi-tenant validation
func TestHapusAkun_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	akunA := &models.Akun{
		IDKoperasi:  koperasiA,
		KodeAkun:    "1101",
		NamaAkun:    "Kas",
		TipeAkun:    models.AkunAset,
		NormalSaldo: "debit",
		StatusAktif: true,
	}
	db.Create(akunA)

	service := NewAkunService(db)

	err := service.HapusAkun(koperasiB, akunA.ID)

	if err == nil {
		t.Fatal("Expected error when deleting cross-tenant akun, got nil")
	}

	var stillExists models.Akun
	result := db.Unscoped().First(&stillExists, akunA.ID)
	if result.Error != nil {
		t.Errorf("Akun was deleted! Should still exist")
	}
	if stillExists.TanggalDihapus.Valid {
		t.Errorf("Akun was soft-deleted! Should not be deleted")
	}

	t.Logf("✓ Akun cross-tenant delete blocked successfully")
}

// TestPerbaruiProduk_CrossTenantBlocked tests produk update multi-tenant validation
func TestPerbaruiProduk_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	produkA := &models.Produk{
		IDKoperasi:  koperasiA,
		KodeProduk:  "P001",
		NamaProduk:  "Product A",
		Harga:       10000,
		Stok:        100,
		StatusAktif: true,
	}
	db.Create(produkA)

	service := NewProdukService(db)

	req := &PerbaruiProdukRequest{
		NamaProduk: "Hacked Product",
	}

	_, err := service.PerbaruiProduk(koperasiB, produkA.ID, req)

	if err == nil {
		t.Fatal("Expected error when updating cross-tenant produk, got nil")
	}

	var unchanged models.Produk
	db.First(&unchanged, produkA.ID)
	if unchanged.NamaProduk != "Product A" {
		t.Errorf("Data was modified! Expected 'Product A', got '%s'", unchanged.NamaProduk)
	}

	t.Logf("✓ Produk cross-tenant update blocked successfully")
}

// TestHapusProduk_CrossTenantBlocked tests produk delete multi-tenant validation
func TestHapusProduk_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	produkA := &models.Produk{
		IDKoperasi:  koperasiA,
		KodeProduk:  "P001",
		NamaProduk:  "Product A",
		Harga:       10000,
		Stok:        100,
		StatusAktif: true,
	}
	db.Create(produkA)

	service := NewProdukService(db)

	err := service.HapusProduk(koperasiB, produkA.ID)

	if err == nil {
		t.Fatal("Expected error when deleting cross-tenant produk, got nil")
	}

	var stillExists models.Produk
	result := db.First(&stillExists, produkA.ID)
	if result.Error != nil {
		t.Errorf("Produk was deleted! Should still exist")
	}

	t.Logf("✓ Produk cross-tenant delete blocked successfully")
}

// TestUbahKataSandiPengguna_CrossTenantBlocked tests password change multi-tenant validation
func TestUbahKataSandiPengguna_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	penggunaA := &models.Pengguna{
		IDKoperasi:   koperasiA,
		NamaLengkap:  "User A",
		NamaPengguna: "userA",
		Email:        "usera@test.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	penggunaA.SetKataSandi("oldpassword123")
	originalHash := penggunaA.KataSandiHash
	db.Create(penggunaA)

	service := NewPenggunaService(db)

	err := service.UbahKataSandiPengguna(koperasiB, penggunaA.ID, "newpassword123")

	if err == nil {
		t.Fatal("Expected error when changing password for cross-tenant user, got nil")
	}

	var unchanged models.Pengguna
	db.First(&unchanged, penggunaA.ID)
	if unchanged.KataSandiHash != originalHash {
		t.Errorf("Password was changed! Hash should remain unchanged")
	}

	t.Logf("✓ Password change cross-tenant blocked successfully")
}

// TestResetKataSandi_CrossTenantBlocked tests password reset multi-tenant validation
func TestResetKataSandi_CrossTenantBlocked(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()

	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	penggunaA := &models.Pengguna{
		IDKoperasi:   koperasiA,
		NamaLengkap:  "User A",
		NamaPengguna: "userA",
		Email:        "usera@test.com",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}
	penggunaA.SetKataSandi("oldpassword123")
	originalHash := penggunaA.KataSandiHash
	db.Create(penggunaA)

	service := NewPenggunaService(db)

	_, err := service.ResetKataSandi(koperasiB, penggunaA.ID)

	if err == nil {
		t.Fatal("Expected error when resetting password for cross-tenant user, got nil")
	}

	var unchanged models.Pengguna
	db.First(&unchanged, penggunaA.ID)
	if unchanged.KataSandiHash != originalHash {
		t.Errorf("Password was reset! Hash should remain unchanged")
	}

	t.Logf("✓ Password reset cross-tenant blocked successfully")
}

// TestMultiTenant_ComprehensiveSecurity runs all security tests together
func TestMultiTenant_ComprehensiveSecurity(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	// Create 3 cooperatives to test isolation
	koop1 := uuid.New()
	koop2 := uuid.New()
	koop3 := uuid.New()

	for _, id := range []uuid.UUID{koop1, koop2, koop3} {
		db.Create(&models.Koperasi{ID: id, NamaKoperasi: "Koperasi " + id.String()[:8]})
		defer cleanupTestData(db, id)
	}

	// Create resources in each cooperative
	for i, koopID := range []uuid.UUID{koop1, koop2, koop3} {
		// Create pengguna
		pengguna := &models.Pengguna{
			IDKoperasi:   koopID,
			NamaLengkap:  fmt.Sprintf("User %d", i+1),
			NamaPengguna: fmt.Sprintf("user%d", i+1),
			Email:        fmt.Sprintf("user%d@test.com", i+1),
			Peran:        models.PeranAdmin,
			StatusAktif:  true,
		}
		pengguna.SetKataSandi("password")
		db.Create(pengguna)

		// Create anggota
		anggota := &models.Anggota{
			IDKoperasi:       koopID,
			NomorAnggota:     fmt.Sprintf("KOOP-%d-0001", 2025+i),
			NamaLengkap:      fmt.Sprintf("Member %d", i+1),
			TanggalBergabung: time.Now(),
			Status:           models.StatusAktif,
		}
		db.Create(anggota)

		// Create akun
		akun := &models.Akun{
			IDKoperasi:  koopID,
			KodeAkun:    fmt.Sprintf("110%d", i+1),
			NamaAkun:    fmt.Sprintf("Cash %d", i+1),
			TipeAkun:    models.AkunAset,
			NormalSaldo: "debit",
			StatusAktif: true,
		}
		db.Create(akun)

		// Create produk
		produk := &models.Produk{
			IDKoperasi:  koopID,
			KodeProduk:  fmt.Sprintf("P00%d", i+1),
			NamaProduk:  fmt.Sprintf("Product %d", i+1),
			Harga:       10000 * float64(i+1),
			Stok:        100,
			StatusAktif: true,
		}
		db.Create(produk)
	}

	// Now attempt cross-tenant operations
	// Koop1 trying to access Koop2's resources
	var koop2Pengguna models.Pengguna
	var koop2Anggota models.Anggota
	var koop2Akun models.Akun
	var koop2Produk models.Produk

	db.Where("id_koperasi = ?", koop2).First(&koop2Pengguna)
	db.Where("id_koperasi = ?", koop2).First(&koop2Anggota)
	db.Where("id_koperasi = ?", koop2).First(&koop2Akun)
	db.Where("id_koperasi = ?", koop2).First(&koop2Produk)

	penggunaService := NewPenggunaService(db)
	anggotaService := NewAnggotaService(db)
	akunService := NewAkunService(db)
	produkService := NewProdukService(db)

	// All these operations should fail
	testCases := []struct {
		name string
		fn   func() error
	}{
		{
			name: "PerbaruiPengguna cross-tenant",
			fn: func() error {
				_, err := penggunaService.PerbaruiPengguna(koop1, koop2Pengguna.ID, &PerbaruiPenggunaRequest{NamaLengkap: "Hack"})
				return err
			},
		},
		{
			name: "HapusPengguna cross-tenant",
			fn: func() error {
				return penggunaService.HapusPengguna(koop1, koop2Pengguna.ID)
			},
		},
		{
			name: "PerbaruiAnggota cross-tenant",
			fn: func() error {
				_, err := anggotaService.PerbaruiAnggota(koop1, koop2Anggota.ID, &PerbaruiAnggotaRequest{NamaLengkap: "Hack"})
				return err
			},
		},
		{
			name: "HapusAnggota cross-tenant",
			fn: func() error {
				return anggotaService.HapusAnggota(koop1, koop2Anggota.ID)
			},
		},
		{
			name: "SetPINPortal cross-tenant",
			fn: func() error {
				return anggotaService.SetPINPortal(koop1, koop2Anggota.ID, "1234")
			},
		},
		{
			name: "PerbaruiAkun cross-tenant",
			fn: func() error {
				_, err := akunService.PerbaruiAkun(koop1, koop2Akun.ID, &PerbaruiAkunRequest{NamaAkun: "Hack"})
				return err
			},
		},
		{
			name: "HapusAkun cross-tenant",
			fn: func() error {
				return akunService.HapusAkun(koop1, koop2Akun.ID)
			},
		},
		{
			name: "PerbaruiProduk cross-tenant",
			fn: func() error {
				_, err := produkService.PerbaruiProduk(koop1, koop2Produk.ID, &PerbaruiProdukRequest{NamaProduk: "Hack"})
				return err
			},
		},
		{
			name: "HapusProduk cross-tenant",
			fn: func() error {
				return produkService.HapusProduk(koop1, koop2Produk.ID)
			},
		},
	}

	passCount := 0
	for _, tc := range testCases {
		err := tc.fn()
		if err == nil {
			t.Errorf("%s: Expected error, got nil - SECURITY VULNERABILITY!", tc.name)
		} else if !strings.Contains(err.Error(), "tidak ditemukan atau tidak memiliki akses") {
			t.Errorf("%s: Expected access denied error, got: %v", tc.name, err)
		} else {
			passCount++
			t.Logf("✓ %s: Blocked correctly", tc.name)
		}
	}

	if passCount == len(testCases) {
		t.Logf("✓✓✓ All %d cross-tenant operations blocked successfully!", passCount)
		t.Logf("✓✓✓ Multi-tenant isolation is SECURE!")
	} else {
		t.Errorf("Security check failed: Only %d/%d operations blocked", passCount, len(testCases))
	}
}
