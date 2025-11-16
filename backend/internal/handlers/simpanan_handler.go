package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SimpananHandler menangani endpoint simpanan anggota
type SimpananHandler struct {
	simpananService *services.SimpananService
}

// NewSimpananHandler membuat instance baru SimpananHandler
func NewSimpananHandler(simpananService *services.SimpananService) *SimpananHandler {
	return &SimpananHandler{
		simpananService: simpananService,
	}
}

// CatatSetoran handles POST /api/v1/simpanan/setor
func (h *SimpananHandler) CatatSetoran(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idPengguna, _ := c.Get("idPengguna")
	penggunaUUID := idPengguna.(uuid.UUID)

	var req services.CatatSetoranRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	simpanan, err := h.simpananService.CatatSetoran(koperasiUUID, penggunaUUID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Setoran simpanan berhasil dicatat", simpanan)
}

// List handles GET /api/v1/simpanan
func (h *SimpananHandler) List(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	tipeSimpanan := c.Query("tipeSimpanan")
	idAnggotaStr := c.Query("idAnggota")
	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")

	// Parse ID anggota jika ada
	var idAnggotaPtr *uuid.UUID
	if idAnggotaStr != "" {
		id, err := uuid.Parse(idAnggotaStr)
		if err == nil {
			idAnggotaPtr = &id
		}
	}

	simpananList, total, err := h.simpananService.GetTransaksiSimpanan(
		koperasiUUID, tipeSimpanan, idAnggotaPtr, tanggalMulai, tanggalAkhir, page, pageSize,
	)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	pagination := utils.CalculatePaginationMeta(page, pageSize, total)
	utils.PaginatedSuccessResponse(c, http.StatusOK, "Data simpanan berhasil diambil", simpananList, pagination)
}

// GetSaldoAnggota handles GET /api/v1/simpanan/anggota/:id
func (h *SimpananHandler) GetSaldoAnggota(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	idAnggota, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID anggota tidak valid")
		return
	}

	saldo, err := h.simpananService.GetSaldoSimpanan(koperasiUUID, idAnggota)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Saldo simpanan berhasil diambil", saldo)
}

// GetRingkasan handles GET /api/v1/simpanan/ringkasan
func (h *SimpananHandler) GetRingkasan(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	ringkasan, err := h.simpananService.GetRingkasanSimpanan(koperasiUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Ringkasan simpanan berhasil diambil", ringkasan)
}

// GetLaporanSaldo handles GET /api/v1/simpanan/laporan-saldo
func (h *SimpananHandler) GetLaporanSaldo(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	laporanSaldo, err := h.simpananService.GetLaporanSaldoSemuaAnggota(koperasiUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Laporan saldo simpanan berhasil diambil", laporanSaldo)
}
