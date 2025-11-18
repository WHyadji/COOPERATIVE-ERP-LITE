# Handler Implementation Guide - Multi-tenant Security

## CRITICAL: Service Method Signature Changes

**Date**: 2025-11-16
**Issue**: Fixed critical multi-tenant validation vulnerability (CVE-INTERNAL-2025-001)

All update and delete service methods now require `idKoperasi` as the first parameter to enforce multi-tenant isolation.

## Updated Service Method Signatures

### PenggunaService

```go
// BEFORE (VULNERABLE)
func (s *PenggunaService) PerbaruiPengguna(id uuid.UUID, req *PerbaruiPenggunaRequest) (*models.PenggunaResponse, error)
func (s *PenggunaService) HapusPengguna(id uuid.UUID) error
func (s *PenggunaService) UbahKataSandiPengguna(id uuid.UUID, kataSandiBaru string) error
func (s *PenggunaService) ResetKataSandi(id uuid.UUID) (string, error)

// AFTER (SECURE)
func (s *PenggunaService) PerbaruiPengguna(idKoperasi, id uuid.UUID, req *PerbaruiPenggunaRequest) (*models.PenggunaResponse, error)
func (s *PenggunaService) HapusPengguna(idKoperasi, id uuid.UUID) error
func (s *PenggunaService) UbahKataSandiPengguna(idKoperasi, id uuid.UUID, kataSandiBaru string) error
func (s *PenggunaService) ResetKataSandi(idKoperasi, id uuid.UUID) (string, error)
```

### AnggotaService

```go
// BEFORE (VULNERABLE)
func (s *AnggotaService) PerbaruiAnggota(id uuid.UUID, req *PerbaruiAnggotaRequest) (*models.AnggotaResponse, error)
func (s *AnggotaService) HapusAnggota(id uuid.UUID) error
func (s *AnggotaService) SetPINPortal(id uuid.UUID, pin string) error

// AFTER (SECURE)
func (s *AnggotaService) PerbaruiAnggota(idKoperasi, id uuid.UUID, req *PerbaruiAnggotaRequest) (*models.AnggotaResponse, error)
func (s *AnggotaService) HapusAnggota(idKoperasi, id uuid.UUID) error
func (s *AnggotaService) SetPINPortal(idKoperasi, id uuid.UUID, pin string) error
```

### AkunService

```go
// BEFORE (VULNERABLE)
func (s *AkunService) PerbaruiAkun(id uuid.UUID, req *PerbaruiAkunRequest) (*models.AkunResponse, error)
func (s *AkunService) HapusAkun(id uuid.UUID) error

// AFTER (SECURE)
func (s *AkunService) PerbaruiAkun(idKoperasi, id uuid.UUID, req *PerbaruiAkunRequest) (*models.AkunResponse, error)
func (s *AkunService) HapusAkun(idKoperasi, id uuid.UUID) error
```

### ProdukService

```go
// BEFORE (VULNERABLE)
func (s *ProdukService) PerbaruiProduk(id uuid.UUID, req *PerbaruiProdukRequest) (*models.ProdukResponse, error)
func (s *ProdukService) HapusProduk(id uuid.UUID) error

// AFTER (SECURE)
func (s *ProdukService) PerbaruiProduk(idKoperasi, id uuid.UUID, req *PerbaruiProdukRequest) (*models.ProdukResponse, error)
func (s *ProdukService) HapusProduk(idKoperasi, id uuid.UUID) error
```

## Handler Implementation Pattern

When implementing handlers, **ALWAYS** extract `idKoperasi` from JWT token and pass it to service methods:

```go
package handlers

import (
	"cooperative-erp-lite/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PenggunaHandler struct {
	service *services.PenggunaService
}

func NewPenggunaHandler(service *services.PenggunaService) *PenggunaHandler {
	return &PenggunaHandler{service: service}
}

// PerbaruiPengguna - CORRECT IMPLEMENTATION
func (h *PenggunaHandler) PerbaruiPengguna(c *gin.Context) {
	// CRITICAL: Extract idKoperasi from JWT middleware
	idKoperasiStr := c.GetString("cooperative_id")
	if idKoperasiStr == "" {
		c.JSON(401, gin.H{"error": "ID koperasi tidak ditemukan dalam token"})
		return
	}

	idKoperasi, err := uuid.Parse(idKoperasiStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID koperasi tidak valid"})
		return
	}

	// Parse resource ID from URL parameter
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	// Parse request body
	var req services.PerbaruiPenggunaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// CRITICAL: Pass idKoperasi to service method
	result, err := h.service.PerbaruiPengguna(idKoperasi, id, &req)
	if err != nil {
		// Error message already handles "tidak ditemukan atau tidak memiliki akses"
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// HapusPengguna - CORRECT IMPLEMENTATION
func (h *PenggunaHandler) HapusPengguna(c *gin.Context) {
	// Extract idKoperasi from JWT
	idKoperasiStr := c.GetString("cooperative_id")
	idKoperasi, err := uuid.Parse(idKoperasiStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Parse resource ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	// CRITICAL: Pass idKoperasi to service method
	err = h.service.HapusPengguna(idKoperasi, id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Pengguna berhasil dihapus"})
}

// UbahKataSandiPengguna - CORRECT IMPLEMENTATION
func (h *PenggunaHandler) UbahKataSandiPengguna(c *gin.Context) {
	// Extract idKoperasi from JWT
	idKoperasiStr := c.GetString("cooperative_id")
	idKoperasi, err := uuid.Parse(idKoperasiStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Parse resource ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	// Parse request
	var req struct {
		KataSandiBaru string `json:"kataSandiBaru" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// CRITICAL: Pass idKoperasi to service method
	err = h.service.UbahKataSandiPengguna(idKoperasi, id, req.KataSandiBaru)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Kata sandi berhasil diubah"})
}

// ResetKataSandi - CORRECT IMPLEMENTATION
func (h *PenggunaHandler) ResetKataSandi(c *gin.Context) {
	// Extract idKoperasi from JWT
	idKoperasiStr := c.GetString("cooperative_id")
	idKoperasi, err := uuid.Parse(idKoperasiStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Parse resource ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	// CRITICAL: Pass idKoperasi to service method
	passwordDefault, err := h.service.ResetKataSandi(idKoperasi, id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":        "Kata sandi berhasil direset",
		"passwordDefault": passwordDefault,
	})
}
```

## Security Checklist for Handler Implementation

When implementing ANY handler that calls update/delete service methods:

- [ ] Extract `idKoperasi` from JWT token using `c.GetString("cooperative_id")`
- [ ] Validate that `idKoperasi` is not empty
- [ ] Parse `idKoperasi` to UUID format
- [ ] Pass `idKoperasi` as the **first parameter** to service methods
- [ ] Return 401 if `idKoperasi` is missing (authentication issue)
- [ ] Return 404 if service returns "tidak ditemukan atau tidak memiliki akses" (authorization issue)
- [ ] **NEVER** call update/delete service methods without passing `idKoperasi`

## Testing Multi-tenant Isolation

Every handler must have tests that verify:

1. **Cross-tenant access is blocked**: User from Koperasi A cannot modify resources from Koperasi B
2. **Same-tenant access works**: User from Koperasi A can modify resources from Koperasi A
3. **Missing token is rejected**: Requests without JWT token return 401
4. **Invalid resource ID returns 404**: Non-existent resource returns proper error

Example test:

```go
func TestPerbaruiPengguna_CrossTenantBlocked(t *testing.T) {
	// Setup
	koperasiA := uuid.New()
	koperasiB := uuid.New()

	// Create user in Koperasi A
	penggunaA := createTestPengguna(koperasiA)

	// Generate JWT token for Koperasi B
	tokenB := generateJWTToken(koperasiB)

	// Attempt to update Koperasi A's user with Koperasi B's token
	req := gin.H{"namaLengkap": "Hacked Name"}
	w := performRequest("PUT", "/api/v1/pengguna/"+penggunaA.ID.String(), req, tokenB)

	// Should return 404 (not found or no access)
	assert.Equal(t, 404, w.Code)

	// Verify data unchanged
	var unchanged models.Pengguna
	db.First(&unchanged, penggunaA.ID)
	assert.NotEqual(t, "Hacked Name", unchanged.NamaLengkap)
}
```

## Routes Configuration

Example routes setup with authentication middleware:

```go
// In cmd/api/main.go or routes setup
func setupRoutes(r *gin.Engine, services *Services) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware()) // CRITICAL: Ensures JWT is validated

	// Pengguna routes
	pengguna := api.Group("/pengguna")
	{
		penggunaHandler := handlers.NewPenggunaHandler(services.PenggunaService)

		pengguna.GET("", penggunaHandler.DapatkanSemuaPengguna)
		pengguna.GET("/:id", penggunaHandler.DapatkanPengguna)
		pengguna.POST("", penggunaHandler.BuatPengguna)
		pengguna.PUT("/:id", penggunaHandler.PerbaruiPengguna)         // Requires idKoperasi
		pengguna.DELETE("/:id", penggunaHandler.HapusPengguna)         // Requires idKoperasi
		pengguna.PUT("/:id/password", penggunaHandler.UbahKataSandiPengguna) // Requires idKoperasi
		pengguna.POST("/:id/reset-password", penggunaHandler.ResetKataSandi)  // Requires idKoperasi
	}

	// Apply same pattern for Anggota, Akun, Produk handlers
}
```

## What NOT to Do

### ❌ WRONG - Missing idKoperasi validation
```go
func (h *PenggunaHandler) PerbaruiPengguna(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req services.PerbaruiPenggunaRequest
	c.ShouldBindJSON(&req)

	// SECURITY VULNERABILITY: Not passing idKoperasi!
	result, err := h.service.PerbaruiPengguna(id, &req) // ❌ WRONG!
	c.JSON(200, result)
}
```

### ❌ WRONG - Trusting client-provided idKoperasi
```go
func (h *PenggunaHandler) PerbaruiPengguna(c *gin.Context) {
	// NEVER trust idKoperasi from request body/query params!
	idKoperasi, _ := uuid.Parse(c.Query("idKoperasi")) // ❌ WRONG!

	id, _ := uuid.Parse(c.Param("id"))
	var req services.PerbaruiPenggunaRequest
	c.ShouldBindJSON(&req)

	result, err := h.service.PerbaruiPengguna(idKoperasi, id, &req)
	c.JSON(200, result)
}
```

### ✅ CORRECT - Extract idKoperasi from JWT
```go
func (h *PenggunaHandler) PerbaruiPengguna(c *gin.Context) {
	// ALWAYS extract idKoperasi from JWT middleware
	idKoperasi, _ := uuid.Parse(c.GetString("cooperative_id")) // ✅ CORRECT!

	id, _ := uuid.Parse(c.Param("id"))
	var req services.PerbaruiPenggunaRequest
	c.ShouldBindJSON(&req)

	result, err := h.service.PerbaruiPengguna(idKoperasi, id, &req)
	c.JSON(200, result)
}
```

## References

- OWASP: Broken Access Control (A01:2021)
- Multi-tenant Security Best Practices: https://cheatsheetseries.owasp.org/cheatsheets/Multitenant_Security_Cheat_Sheet.html
- Internal QA Report: `docs/qa-report/services-qa-report-2025-11-16.md`

## Next Steps

1. When implementing handlers, follow this guide exactly
2. Add integration tests for each handler
3. Perform security audit before MVP launch (Week 12)
4. Consider adding automated security scanning in CI/CD pipeline
