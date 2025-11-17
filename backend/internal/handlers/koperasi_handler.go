package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// KoperasiHandler menangani endpoint manajemen koperasi
type KoperasiHandler struct {
	koperasiService *services.KoperasiService
}

// NewKoperasiHandler membuat instance baru KoperasiHandler
func NewKoperasiHandler(koperasiService *services.KoperasiService) *KoperasiHandler {
	return &KoperasiHandler{
		koperasiService: koperasiService,
	}
}

// Create handles POST /api/v1/koperasi
// @Summary Buat koperasi baru
// @Description Membuat koperasi baru (hanya Admin)
// @Tags Koperasi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.BuatKoperasiRequest true "Data koperasi"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Router /api/v1/koperasi [post]
func (h *KoperasiHandler) Create(c *gin.Context) {
	var req services.BuatKoperasiRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Buat koperasi
	koperasi, err := h.koperasiService.BuatKoperasi(&req)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Koperasi berhasil dibuat", koperasi)
}

// List handles GET /api/v1/koperasi
// @Summary List semua koperasi
// @Description Mendapatkan daftar semua koperasi dengan pagination
// @Tags Koperasi
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman" default(1)
// @Param pageSize query int false "Jumlah per halaman" default(20)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/koperasi [get]
func (h *KoperasiHandler) List(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// Get list koperasi
	koperasiList, total, err := h.koperasiService.GetSemuaKoperasi(page, pageSize)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	// Calculate pagination metadata
	pagination := utils.CalculatePaginationMeta(page, pageSize, total)

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Data koperasi berhasil diambil", koperasiList, pagination)
}

// GetByID handles GET /api/v1/koperasi/:id
// @Summary Get detail koperasi
// @Description Mendapatkan detail koperasi berdasarkan ID
// @Tags Koperasi
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Koperasi (UUID)"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/v1/koperasi/{id} [get]
func (h *KoperasiHandler) GetByID(c *gin.Context) {
	// Parse UUID dari parameter
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID koperasi tidak valid")
		return
	}

	// Get koperasi
	koperasi, err := h.koperasiService.GetKoperasiByID(id)
	if err != nil {
		utils.NotFoundResponse(c, "Koperasi tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data koperasi berhasil diambil", koperasi)
}

// Update handles PUT /api/v1/koperasi/:id
// @Summary Update koperasi
// @Description Mengupdate data koperasi
// @Tags Koperasi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Koperasi (UUID)"
// @Param body body services.PerbaruiKoperasiRequest true "Data koperasi"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/v1/koperasi/{id} [put]
func (h *KoperasiHandler) Update(c *gin.Context) {
	// Parse UUID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID koperasi tidak valid")
		return
	}

	var req services.PerbaruiKoperasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Update koperasi
	koperasi, err := h.koperasiService.PerbaruiKoperasi(id, &req)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Koperasi berhasil diupdate", koperasi)
}

// Delete handles DELETE /api/v1/koperasi/:id
// @Summary Hapus koperasi
// @Description Menghapus koperasi (soft delete)
// @Tags Koperasi
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Koperasi (UUID)"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/v1/koperasi/{id} [delete]
func (h *KoperasiHandler) Delete(c *gin.Context) {
	// Parse UUID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID koperasi tidak valid")
		return
	}

	// Delete koperasi
	if err := h.koperasiService.HapusKoperasi(id); err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Koperasi berhasil dihapus", nil)
}

// GetStatistik handles GET /api/v1/koperasi/:id/statistik
// @Summary Get statistik koperasi
// @Description Mendapatkan statistik koperasi (jumlah anggota, total simpanan, dll)
// @Tags Koperasi
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Koperasi (UUID)"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/v1/koperasi/{id}/statistik [get]
func (h *KoperasiHandler) GetStatistik(c *gin.Context) {
	// Parse UUID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID koperasi tidak valid")
		return
	}

	// Get statistik
	statistik, err := h.koperasiService.GetStatistikKoperasi(id)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistik koperasi berhasil diambil", statistik)
}
