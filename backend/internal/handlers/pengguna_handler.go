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

// PenggunaHandler menangani endpoint manajemen pengguna
type PenggunaHandler struct {
	penggunaService *services.PenggunaService
}

// NewPenggunaHandler membuat instance baru PenggunaHandler
func NewPenggunaHandler(penggunaService *services.PenggunaService) *PenggunaHandler {
	return &PenggunaHandler{
		penggunaService: penggunaService,
	}
}

// Create handles POST /api/v1/pengguna
func (h *PenggunaHandler) Create(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	var req services.BuatPenggunaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	pengguna, err := h.penggunaService.BuatPengguna(koperasiUUID, &req)
	if err != nil {
		if err.Error() == "nama pengguna sudah ada" {
			utils.ConflictResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Pengguna berhasil dibuat", pengguna)
}

// List handles GET /api/v1/pengguna
func (h *PenggunaHandler) List(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	peran := c.Query("peran")               // optional filter
	statusAktif := c.Query("statusAktif")    // optional filter

	// Parse status aktif
	var statusAktifPtr *bool
	if statusAktif != "" {
		aktif := statusAktif == "true"
		statusAktifPtr = &aktif
	}

	// Parse peran
	var peranPtr *models.PeranPengguna
	if peran != "" {
		p := models.PeranPengguna(peran)
		peranPtr = &p
	}

	penggunaList, total, err := h.penggunaService.GetSemuaPengguna(koperasiUUID, peranPtr, statusAktifPtr, page, pageSize)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	pagination := utils.CalculatePaginationMeta(page, pageSize, total)
	utils.PaginatedSuccessResponse(c, http.StatusOK, "Data pengguna berhasil diambil", penggunaList, pagination)
}

// GetByID handles GET /api/v1/pengguna/:id
func (h *PenggunaHandler) GetByID(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID pengguna tidak valid")
		return
	}

	pengguna, err := h.penggunaService.GetPenggunaByID(koperasiUUID, id)
	if err != nil {
		utils.NotFoundResponse(c, "Pengguna tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data pengguna berhasil diambil", pengguna)
}

// Update handles PUT /api/v1/pengguna/:id
func (h *PenggunaHandler) Update(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID pengguna tidak valid")
		return
	}

	var req services.PerbaruiPenggunaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	pengguna, err := h.penggunaService.PerbaruiPengguna(koperasiUUID, id, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pengguna berhasil diupdate", pengguna)
}

// Delete handles DELETE /api/v1/pengguna/:id
func (h *PenggunaHandler) Delete(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID pengguna tidak valid")
		return
	}

	if err := h.penggunaService.HapusPengguna(koperasiUUID, id); err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pengguna berhasil dihapus", nil)
}

// ResetPassword handles POST /api/v1/pengguna/:id/reset-password
func (h *PenggunaHandler) ResetPassword(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID pengguna tidak valid")
		return
	}

	var req struct {
		KataSandiBaru string `json:"kataSandiBaru" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.penggunaService.ResetKataSandi(koperasiUUID, id, req.KataSandiBaru); err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Kata sandi berhasil direset", nil)
}

// ChangePassword handles PUT /api/v1/pengguna/:id/change-password
func (h *PenggunaHandler) ChangePassword(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID pengguna tidak valid")
		return
	}

	var req struct {
		KataSandiLama string `json:"kataSandiLama" binding:"required"`
		KataSandiBaru string `json:"kataSandiBaru" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.penggunaService.UbahKataSandiAdmin(koperasiUUID, id, req.KataSandiLama, req.KataSandiBaru); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Kata sandi berhasil diubah", nil)
}
