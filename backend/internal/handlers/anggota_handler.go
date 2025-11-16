package handlers

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AnggotaHandler menangani endpoint manajemen anggota
type AnggotaHandler struct {
	anggotaService *services.AnggotaService
}

// NewAnggotaHandler membuat instance baru AnggotaHandler
func NewAnggotaHandler(anggotaService *services.AnggotaService) *AnggotaHandler {
	return &AnggotaHandler{
		anggotaService: anggotaService,
	}
}

// Create handles POST /api/v1/anggota
func (h *AnggotaHandler) Create(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	var req services.BuatAnggotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	anggota, err := h.anggotaService.BuatAnggota(koperasiUUID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Anggota berhasil dibuat", anggota)
}

// List handles GET /api/v1/anggota
func (h *AnggotaHandler) List(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	search := c.Query("search")

	// Parse status
	var statusPtr *models.StatusAnggota
	if status != "" {
		s := models.StatusAnggota(status)
		statusPtr = &s
	}

	anggotaList, total, err := h.anggotaService.GetSemuaAnggota(koperasiUUID, statusPtr, search, page, pageSize)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	pagination := utils.CalculatePaginationMeta(page, pageSize, total)
	utils.PaginatedSuccessResponse(c, http.StatusOK, "Data anggota berhasil diambil", anggotaList, pagination)
}

// GetByID handles GET /api/v1/anggota/:id
func (h *AnggotaHandler) GetByID(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID anggota tidak valid")
		return
	}

	anggota, err := h.anggotaService.GetAnggotaByID(koperasiUUID, id)
	if err != nil {
		utils.NotFoundResponse(c, "Anggota tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data anggota berhasil diambil", anggota)
}

// GetByNomor handles GET /api/v1/anggota/nomor/:nomor
func (h *AnggotaHandler) GetByNomor(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	nomorAnggota := c.Param("nomor")

	anggota, err := h.anggotaService.GetAnggotaByNomor(koperasiUUID, nomorAnggota)
	if err != nil {
		utils.NotFoundResponse(c, "Anggota tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data anggota berhasil diambil", anggota)
}

// Update handles PUT /api/v1/anggota/:id
func (h *AnggotaHandler) Update(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID anggota tidak valid")
		return
	}

	var req services.PerbaruiAnggotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	anggota, err := h.anggotaService.PerbaruiAnggota(koperasiUUID, id, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Anggota berhasil diupdate", anggota)
}

// Delete handles DELETE /api/v1/anggota/:id
func (h *AnggotaHandler) Delete(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID anggota tidak valid")
		return
	}

	if err := h.anggotaService.HapusAnggota(koperasiUUID, id); err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Anggota berhasil dihapus", nil)
}

// SetPIN handles POST /api/v1/anggota/:id/set-pin
func (h *AnggotaHandler) SetPIN(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID anggota tidak valid")
		return
	}

	var req struct {
		PIN string `json:"pin" binding:"required,len=6,numeric"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "PIN harus 6 digit angka")
		return
	}

	if err := h.anggotaService.SetPINPortal(koperasiUUID, id, req.PIN); err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "PIN portal berhasil diset", nil)
}

// ValidatePIN handles POST /api/v1/anggota/validate-pin
func (h *AnggotaHandler) ValidatePIN(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	var req struct {
		NomorAnggota string `json:"nomorAnggota" binding:"required"`
		PIN          string `json:"pin" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	valid, err := h.anggotaService.ValidasiPINPortal(koperasiUUID, req.NomorAnggota, req.PIN)
	if err != nil || !valid {
		utils.UnauthorizedResponse(c, "Nomor anggota atau PIN tidak valid")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "PIN valid", gin.H{"valid": true})
}

// GetStatistik handles GET /api/v1/anggota/statistik
func (h *AnggotaHandler) GetStatistik(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	jumlah, err := h.anggotaService.HitungJumlahAnggota(koperasiUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistik anggota berhasil diambil", gin.H{
		"jumlahAnggota": jumlah,
	})
}
