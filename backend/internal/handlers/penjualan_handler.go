package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PenjualanHandler menangani endpoint POS penjualan
type PenjualanHandler struct {
	penjualanService *services.PenjualanService
}

// NewPenjualanHandler membuat instance baru PenjualanHandler
func NewPenjualanHandler(penjualanService *services.PenjualanService) *PenjualanHandler {
	return &PenjualanHandler{
		penjualanService: penjualanService,
	}
}

// ProsesPenjualan handles POST /api/v1/penjualan
func (h *PenjualanHandler) ProsesPenjualan(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idPengguna, _ := c.Get("idPengguna")
	kasirUUID := idPengguna.(uuid.UUID)

	var req services.ProsesPenjualanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	penjualan, err := h.penjualanService.ProsesPenjualan(koperasiUUID, kasirUUID, &req)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Penjualan berhasil diproses", penjualan)
}

// List handles GET /api/v1/penjualan
func (h *PenjualanHandler) List(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")
	idKasirStr := c.Query("idKasir")

	// Parse ID kasir jika ada
	var idKasirPtr *uuid.UUID
	if idKasirStr != "" {
		id, err := uuid.Parse(idKasirStr)
		if err == nil {
			idKasirPtr = &id
		}
	}

	penjualanList, total, err := h.penjualanService.DapatkanSemuaPenjualan(
		koperasiUUID, tanggalMulai, tanggalAkhir, idKasirPtr, page, pageSize,
	)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	pagination := utils.CalculatePaginationMeta(page, pageSize, total)
	utils.PaginatedSuccessResponse(c, http.StatusOK, "Data penjualan berhasil diambil", penjualanList, pagination)
}

// GetByID handles GET /api/v1/penjualan/:id
func (h *PenjualanHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID penjualan tidak valid")
		return
	}

	penjualan, err := h.penjualanService.DapatkanPenjualan(id)
	if err != nil {
		utils.NotFoundResponse(c, "Penjualan tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data penjualan berhasil diambil", penjualan)
}

// GetStruk handles GET /api/v1/penjualan/:id/struk
func (h *PenjualanHandler) GetStruk(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID penjualan tidak valid")
		return
	}

	struk, err := h.penjualanService.DapatkanStruk(id)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Struk digital berhasil digenerate", struk)
}

// GetHariIni handles GET /api/v1/penjualan/hari-ini
func (h *PenjualanHandler) GetHariIni(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	summary, err := h.penjualanService.DapatkanPenjualanHariIni(koperasiUUID)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Summary penjualan hari ini berhasil diambil", summary)
}

// GetTopProduk handles GET /api/v1/penjualan/top-produk
func (h *PenjualanHandler) GetTopProduk(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	topProduk, err := h.penjualanService.DapatkanTopProduk(koperasiUUID, limit)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Top produk berhasil diambil", topProduk)
}
