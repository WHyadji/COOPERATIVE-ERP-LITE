package handlers

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AkunHandler menangani endpoint chart of accounts
type AkunHandler struct {
	akunService *services.AkunService
}

// NewAkunHandler membuat instance baru AkunHandler
func NewAkunHandler(akunService *services.AkunService) *AkunHandler {
	return &AkunHandler{
		akunService: akunService,
	}
}

// Create handles POST /api/v1/akun
func (h *AkunHandler) Create(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	var req services.BuatAkunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	akun, err := h.akunService.BuatAkun(koperasiUUID, &req)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Akun berhasil dibuat", akun)
}

// List handles GET /api/v1/akun
func (h *AkunHandler) List(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tipeAkun := c.Query("tipeAkun")
	statusAktif := c.Query("statusAktif")

	// Parse tipe akun
	var tipeAkunPtr *models.TipeAkun
	if tipeAkun != "" {
		t := models.TipeAkun(tipeAkun)
		tipeAkunPtr = &t
	}

	// Parse status aktif
	var statusAktifPtr *bool
	if statusAktif != "" {
		aktif := statusAktif == "true"
		statusAktifPtr = &aktif
	}

	akunList, err := h.akunService.GetSemuaAkun(koperasiUUID, tipeAkunPtr, statusAktifPtr)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data akun berhasil diambil", akunList)
}

// GetByID handles GET /api/v1/akun/:id
func (h *AkunHandler) GetByID(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID akun tidak valid")
		return
	}

	akun, err := h.akunService.GetAkunByID(koperasiUUID, id)
	if err != nil {
		utils.NotFoundResponse(c, "Akun tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data akun berhasil diambil", akun)
}

// Update handles PUT /api/v1/akun/:id
func (h *AkunHandler) Update(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID akun tidak valid")
		return
	}

	var req services.PerbaruiAkunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	akun, err := h.akunService.PerbaruiAkun(koperasiUUID, id, &req)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Akun berhasil diupdate", akun)
}

// Delete handles DELETE /api/v1/akun/:id
func (h *AkunHandler) Delete(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID akun tidak valid")
		return
	}

	if err := h.akunService.HapusAkun(koperasiUUID, id); err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Akun berhasil dihapus", nil)
}

// GetSaldo handles GET /api/v1/akun/:id/saldo
func (h *AkunHandler) GetSaldo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID akun tidak valid")
		return
	}

	tanggalPer := c.Query("tanggalPer") // Optional, format: YYYY-MM-DD

	saldo, err := h.akunService.HitungSaldoAkun(id, tanggalPer)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Saldo akun berhasil dihitung", gin.H{
		"saldo": saldo,
	})
}

// SeedCOA handles POST /api/v1/akun/seed-coa
func (h *AkunHandler) SeedCOA(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	if err := h.akunService.InisialisasiCOADefault(koperasiUUID); err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Chart of Accounts default berhasil di-seed", nil)
}
