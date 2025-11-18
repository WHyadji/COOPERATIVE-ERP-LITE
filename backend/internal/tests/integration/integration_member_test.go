package integration

import (
	"fmt"
	"testing"

	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"github.com/google/uuid"
)

// TestMemberIntegration_CompleteCRUDCycle tests complete member CRUD operations
func TestMemberIntegration_CompleteCRUDCycle(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	// Setup koperasi
	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi Member",
	}
	db.Create(koperasi)

	anggotaService := services.NewAnggotaService(db)

	// Step 1: CREATE - Add new member
	createReq := &services.BuatAnggotaRequest{
		NamaLengkap:   "John Doe",
		NIK:           "1234567890123456",
		TempatLahir:   "Jakarta",
		JenisKelamin:  "L",
		Alamat:        "Jl. Test No. 123",
		RT:            "001",
		RW:            "002",
		Kelurahan:     "Test Kelurahan",
		Kecamatan:     "Test Kecamatan",
		KotaKabupaten: "Jakarta Selatan",
		Provinsi:      "DKI Jakarta",
		KodePos:       "12345",
		NoTelepon:     "08123456789",
		Email:         "john@test.com",
	}

	createdMember, err := anggotaService.BuatAnggota(koperasiID, createReq)
	if err != nil {
		t.Fatalf("Failed to create member: %v", err)
	}

	if createdMember.ID == uuid.Nil {
		t.Error("Created member should have valid ID")
	}
	if createdMember.NamaLengkap != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", createdMember.NamaLengkap)
	}

	t.Logf("✓ Member created successfully with ID: %s", createdMember.ID)

	// Step 2: READ - Get member by ID
	retrievedMember, err := anggotaService.DapatkanAnggota(createdMember.ID)
	if err != nil {
		t.Fatalf("Failed to get member by ID: %v", err)
	}

	if retrievedMember.ID != createdMember.ID {
		t.Error("Retrieved member should have same ID")
	}
	if retrievedMember.NomorAnggota != createdMember.NomorAnggota {
		t.Error("Retrieved member should have correct member number")
	}

	t.Logf("✓ Member retrieved successfully: %s", retrievedMember.NamaLengkap)

	// Step 3: UPDATE - Modify member data
	updateReq := &services.PerbaruiAnggotaRequest{
		NamaLengkap: "John Doe Updated",
		NoTelepon:   "08987654321",
		Email:       "john.updated@test.com",
		Alamat:      "Jl. Updated No. 456",
	}

	updatedMember, err := anggotaService.PerbaruiAnggota(koperasiID, createdMember.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update member: %v", err)
	}

	if updatedMember.NamaLengkap != "John Doe Updated" {
		t.Error("Member name should be updated")
	}
	if updatedMember.NoTelepon != "08987654321" {
		t.Error("Member phone should be updated")
	}

	t.Logf("✓ Member updated successfully: %s", updatedMember.NamaLengkap)

	// Step 4: LIST - Get all members
	daftarAnggota, _, err := anggotaService.DapatkanSemuaAnggota(koperasiID, "", "", 1, 10)
	if err != nil {
		t.Fatalf("Failed to get all members: %v", err)
	}

	if len(daftarAnggota) != 1 {
		t.Errorf("Expected 1 member, got %d", len(daftarAnggota))
	}

	t.Logf("✓ Member list retrieved: %d members", len(daftarAnggota))

	// Step 5: DELETE - Remove member
	err = anggotaService.HapusAnggota(koperasiID, createdMember.ID)
	if err != nil {
		t.Fatalf("Failed to delete member: %v", err)
	}

	// Verify deletion
	_, err = anggotaService.DapatkanAnggota(createdMember.ID)
	if err == nil {
		t.Error("Deleted member should not be retrievable")
	}

	t.Logf("✓ Member deleted successfully")

	t.Log("✅ Complete member CRUD integration test passed")
}

// TestMemberIntegration_SearchAndFilter tests member search and filtering
func TestMemberIntegration_SearchAndFilter(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)

	anggotaService := services.NewAnggotaService(db)

	// Create multiple members (service will auto-generate member numbers)
	members := []struct {
		nama string
	}{
		{"Ahmad Santoso"},
		{"Budi Hartono"},
		{"Citra Dewi"},
		{"Ahmad Rizki"},
	}

	for i, m := range members {
		req := &services.BuatAnggotaRequest{
			NamaLengkap:   m.nama,
			NIK:           "1234567890123456",
			TempatLahir:   "Jakarta",
			JenisKelamin:  "L",
			Alamat:        "Test",
			RT:            "001",
			RW:            "002",
			Kelurahan:     "Test",
			Kecamatan:     "Test",
			KotaKabupaten: "Jakarta",
			Provinsi:      "DKI",
			KodePos:       "12345",
			NoTelepon:     "08123456789",
			Email:         fmt.Sprintf("member%d@test.com", i+1), // Valid email without spaces
		}

		_, err := anggotaService.BuatAnggota(koperasiID, req)
		if err != nil {
			t.Fatalf("Failed to create member %s: %v", m.nama, err)
		}
	}

	// Test 1: Search by name
	searchResults, _, err := anggotaService.DapatkanSemuaAnggota(koperasiID, "", "Ahmad", 1, 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(searchResults) != 2 {
		t.Errorf("Expected 2 members with name 'Ahmad', got %d", len(searchResults))
	}

	t.Logf("✓ Search by name returned %d results", len(searchResults))

	// Test 2: Filter by status (aktif is default)
	activeMembers, _, err := anggotaService.DapatkanSemuaAnggota(koperasiID, string(models.StatusAktif), "", 1, 10)
	if err != nil {
		t.Fatalf("Filter failed: %v", err)
	}

	if len(activeMembers) != 4 {
		t.Errorf("Expected 4 active members, got %d", len(activeMembers))
	}

	t.Logf("✓ Filter by status returned %d active members", len(activeMembers))

	// Test 3: Pagination
	page1, total, err := anggotaService.DapatkanSemuaAnggota(koperasiID, "", "", 1, 2)
	if err != nil {
		t.Fatalf("Pagination failed: %v", err)
	}

	if len(page1) != 2 {
		t.Errorf("Expected 2 members on page 1, got %d", len(page1))
	}
	if total != 4 {
		t.Errorf("Expected total 4 members, got %d", total)
	}

	page2, _, err := anggotaService.DapatkanSemuaAnggota(koperasiID, "", "", 2, 2)
	if err != nil {
		t.Fatalf("Pagination page 2 failed: %v", err)
	}

	if len(page2) != 2 {
		t.Errorf("Expected 2 members on page 2, got %d", len(page2))
	}

	t.Logf("✓ Pagination working correctly: page1=%d, page2=%d, total=%d", len(page1), len(page2), total)

	t.Log("✅ Member search and filter integration test passed")
}

// TestMemberIntegration_DuplicateValidation tests that auto-generated member numbers are unique
func TestMemberIntegration_DuplicateValidation(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiID := uuid.New()
	defer cleanupTestData(db, koperasiID)

	koperasi := &models.Koperasi{
		ID:           koperasiID,
		NamaKoperasi: "Test Koperasi",
	}
	db.Create(koperasi)

	anggotaService := services.NewAnggotaService(db)

	// Create first member
	req1 := &services.BuatAnggotaRequest{
		NamaLengkap:  "First Member",
		NIK:          "1234567890123456",
		TempatLahir:  "Jakarta",
		JenisKelamin: "L",
		Alamat:       "Test",
		RT:           "001",
		RW:           "002",
		Kelurahan:    "Test",
		Kecamatan:    "Test",
		KotaKabupaten: "Jakarta",
		Provinsi:     "DKI",
		KodePos:      "12345",
		NoTelepon:    "08123456789",
		Email:        "first@test.com",
	}

	member1, err := anggotaService.BuatAnggota(koperasiID, req1)
	if err != nil {
		t.Fatalf("Failed to create first member: %v", err)
	}

	// Create second member
	req2 := &services.BuatAnggotaRequest{
		NamaLengkap:  "Second Member",
		NIK:          "9876543210987654",
		TempatLahir:  "Bandung",
		JenisKelamin: "P",
		Alamat:       "Test 2",
		RT:           "001",
		RW:           "002",
		Kelurahan:    "Test",
		Kecamatan:    "Test",
		KotaKabupaten: "Bandung",
		Provinsi:     "Jabar",
		KodePos:      "54321",
		NoTelepon:    "08987654321",
		Email:        "second@test.com",
	}

	member2, err := anggotaService.BuatAnggota(koperasiID, req2)
	if err != nil {
		t.Fatalf("Failed to create second member: %v", err)
	}

	// Verify member numbers are unique (auto-generated)
	if member1.NomorAnggota == member2.NomorAnggota {
		t.Error("Auto-generated member numbers should be unique")
	}

	t.Logf("✓ Member numbers are unique: %s != %s", member1.NomorAnggota, member2.NomorAnggota)
	t.Log("✅ Auto-generated member number uniqueness test passed")
}

// TestMemberIntegration_CrossTenantIsolation tests that members are isolated between cooperatives
func TestMemberIntegration_CrossTenantIsolation(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	koperasiA := uuid.New()
	koperasiB := uuid.New()
	defer cleanupTestData(db, koperasiA)
	defer cleanupTestData(db, koperasiB)

	// Create two cooperatives
	koopA := &models.Koperasi{ID: koperasiA, NamaKoperasi: "Koperasi A"}
	koopB := &models.Koperasi{ID: koperasiB, NamaKoperasi: "Koperasi B"}
	db.Create(koopA)
	db.Create(koopB)

	anggotaService := services.NewAnggotaService(db)

	// Create member in Koperasi A
	reqA := &services.BuatAnggotaRequest{
		NamaLengkap:  "Member A",
		NIK:          "1111111111111111",
		TempatLahir:  "Jakarta",
		JenisKelamin: "L",
		Alamat:       "Test",
		RT:           "001",
		RW:           "002",
		Kelurahan:    "Test",
		Kecamatan:    "Test",
		KotaKabupaten: "Jakarta",
		Provinsi:     "DKI",
		KodePos:      "12345",
		NoTelepon:    "08111111111",
		Email:        "memberA@test.com",
	}

	memberA, err := anggotaService.BuatAnggota(koperasiA, reqA)
	if err != nil {
		t.Fatalf("Failed to create member in Koperasi A: %v", err)
	}

	// Create member in Koperasi B
	reqB := &services.BuatAnggotaRequest{
		NamaLengkap:  "Member B",
		NIK:          "2222222222222222",
		TempatLahir:  "Bandung",
		JenisKelamin: "P",
		Alamat:       "Test",
		RT:           "001",
		RW:           "002",
		Kelurahan:    "Test",
		Kecamatan:    "Test",
		KotaKabupaten: "Bandung",
		Provinsi:     "Jabar",
		KodePos:      "54321",
		NoTelepon:    "08222222222",
		Email:        "memberB@test.com",
	}

	_, err = anggotaService.BuatAnggota(koperasiB, reqB)
	if err != nil {
		t.Fatalf("Failed to create member in Koperasi B: %v", err)
	}

	// Koperasi A should only see 1 member
	membersA, totalA, err := anggotaService.DapatkanSemuaAnggota(koperasiA, "", "", 1, 10)
	if err != nil {
		t.Fatalf("Failed to get members for Koperasi A: %v", err)
	}

	if totalA != 1 {
		t.Errorf("Koperasi A should have 1 member, got %d", totalA)
	}

	// Koperasi B should only see 1 member
	membersB, totalB, err := anggotaService.DapatkanSemuaAnggota(koperasiB, "", "", 1, 10)
	if err != nil {
		t.Fatalf("Failed to get members for Koperasi B: %v", err)
	}

	if totalB != 1 {
		t.Errorf("Koperasi B should have 1 member, got %d", totalB)
	}

	t.Logf("✓ Koperasi A has %d member(s), Koperasi B has %d member(s)", len(membersA), len(membersB))

	// Try to get Koperasi A's member by ID (no cooperative filter)
	// This should succeed as DapatkanAnggota doesn't filter by cooperative
	retrievedMember, err := anggotaService.DapatkanAnggota(memberA.ID)
	if err != nil {
		t.Errorf("Should be able to get member by ID: %v", err)
	}

	// Verify the member belongs to Koperasi A
	if retrievedMember != nil && retrievedMember.IDKoperasi != koperasiA {
		t.Error("Retrieved member should belong to Koperasi A")
	}

	t.Logf("✓ Member retrieved successfully with correct cooperative ID")
	t.Log("✅ Cross-tenant isolation integration test passed")
}
