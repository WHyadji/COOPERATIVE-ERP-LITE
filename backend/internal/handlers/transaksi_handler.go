package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransaksiHandler menangani endpoint transaksi akuntansi
type TransaksiHandler struct {
	transaksiService *services.TransaksiService
}

// NewTransaksiHandler membuat instance baru TransaksiHandler
func NewTransaksiHandler(transaksiService *services.TransaksiService) *TransaksiHandler {
	return &TransaksiHandler{
		transaksiService: transaksiService,
	}
}

// Create handles POST /api/v1/transaksi
func (h *TransaksiHandler) Create(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idPengguna, _ := c.Get("idPengguna")
	penggunaUUID := idPengguna.(uuid.UUID)

	var req services.BuatTransaksiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	transaksi, err := h.transaksiService.BuatTransaksi(koperasiUUID, penggunaUUID, &req)
	if err != nil {
		// Check jika error validasi double-entry
		if err.Error() == "total debit harus sama dengan total kredit" ||
		   err.Error() == "satu baris tidak boleh memiliki debit dan kredit sekaligus" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Transaksi berhasil dibuat", transaksi)
}

// List handles GET /api/v1/transaksi
func (h *TransaksiHandler) List(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")
	tipeTransaksi := c.Query("tipeTransaksi")

	transaksiList, total, err := h.transaksiService.GetSemuaTransaksi(
		koperasiUUID, tanggalMulai, tanggalAkhir, tipeTransaksi, page, pageSize,
	)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	pagination := utils.CalculatePaginationMeta(page, pageSize, total)
	utils.PaginatedSuccessResponse(c, http.StatusOK, "Data transaksi berhasil diambil", transaksiList, pagination)
}

// GetByID handles GET /api/v1/transaksi/:id
func (h *TransaksiHandler) GetByID(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID transaksi tidak valid")
		return
	}

	transaksi, err := h.transaksiService.GetTransaksiByID(koperasiUUID, id)
	if err != nil {
		utils.NotFoundResponse(c, "Transaksi tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data transaksi berhasil diambil", transaksi)
}

// Update handles PUT /api/v1/transaksi/:id
func (h *TransaksiHandler) Update(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID transaksi tidak valid")
		return
	}

	var req services.BuatTransaksiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	transaksi, err := h.transaksiService.PerbaruiTransaksi(koperasiUUID, id, &req)
	if err != nil {
		if err.Error() == "transaksi sudah di-post, tidak dapat diubah" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaksi berhasil diupdate", transaksi)
}

// Delete handles DELETE /api/v1/transaksi/:id
func (h *TransaksiHandler) Delete(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID transaksi tidak valid")
		return
	}

	if err := h.transaksiService.HapusTransaksi(koperasiUUID, id); err != nil {
		if err.Error() == "transaksi sudah di-post, tidak dapat dihapus" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaksi berhasil dihapus", nil)
}

// Reverse handles POST /api/v1/transaksi/:id/reverse
func (h *TransaksiHandler) Reverse(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idPengguna, _ := c.Get("idPengguna")
	penggunaUUID := idPengguna.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID transaksi tidak valid")
		return
	}

	var req struct {
		Keterangan string `json:"keterangan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	reversedTransaksi, err := h.transaksiService.ReverseTransaksi(koperasiUUID, penggunaUUID, id, req.Keterangan)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaksi berhasil di-reverse", reversedTransaksi)
}
