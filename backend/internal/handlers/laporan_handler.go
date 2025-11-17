package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LaporanHandler menangani endpoint laporan keuangan
type LaporanHandler struct {
	laporanService *services.LaporanService
}

// NewLaporanHandler membuat instance baru LaporanHandler
func NewLaporanHandler(laporanService *services.LaporanService) *LaporanHandler {
	return &LaporanHandler{
		laporanService: laporanService,
	}
}

// GetNeraca handles GET /api/v1/laporan/neraca
func (h *LaporanHandler) GetNeraca(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tanggalPer := c.DefaultQuery("tanggalPer", "") // Format: YYYY-MM-DD

	neraca, err := h.laporanService.GenerateLaporanPosisiKeuangan(koperasiUUID, tanggalPer)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Laporan neraca berhasil digenerate", neraca)
}

// GetLabaRugi handles GET /api/v1/laporan/laba-rugi
func (h *LaporanHandler) GetLabaRugi(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")

	if tanggalMulai == "" || tanggalAkhir == "" {
		utils.BadRequestResponse(c, "Parameter tanggalMulai dan tanggalAkhir wajib diisi")
		return
	}

	labaRugi, err := h.laporanService.GenerateLaporanLabaRugi(koperasiUUID, tanggalMulai, tanggalAkhir)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Laporan laba rugi berhasil digenerate", labaRugi)
}

// GetPerubahanModal handles GET /api/v1/laporan/perubahan-modal
func (h *LaporanHandler) GetPerubahanModal(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")

	if tanggalMulai == "" || tanggalAkhir == "" {
		utils.BadRequestResponse(c, "Parameter tanggalMulai dan tanggalAkhir wajib diisi")
		return
	}

	perubahanModal, err := h.laporanService.GenerateLaporanPerubahanModal(koperasiUUID, tanggalMulai, tanggalAkhir)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Laporan perubahan modal berhasil digenerate", perubahanModal)
}

// GetArusKas handles GET /api/v1/laporan/arus-kas
func (h *LaporanHandler) GetArusKas(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")

	if tanggalMulai == "" || tanggalAkhir == "" {
		utils.BadRequestResponse(c, "Parameter tanggalMulai dan tanggalAkhir wajib diisi")
		return
	}

	arusKas, err := h.laporanService.GenerateLaporanArusKas(koperasiUUID, tanggalMulai, tanggalAkhir)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Laporan arus kas berhasil digenerate", arusKas)
}

// GetBukuBesar handles GET /api/v1/laporan/buku-besar
func (h *LaporanHandler) GetBukuBesar(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")
	idAkunStr := c.Query("idAkun") // Optional - jika kosong, semua akun

	var idAkun uuid.UUID
	if idAkunStr != "" {
		id, err := uuid.Parse(idAkunStr)
		if err != nil {
			utils.BadRequestResponse(c, "ID akun tidak valid")
			return
		}
		idAkun = id
	} else {
		// If no account specified, return error
		utils.BadRequestResponse(c, "ID akun diperlukan")
		return
	}

	bukuBesar, err := h.laporanService.GenerateBukuBesar(koperasiUUID, idAkun, tanggalMulai, tanggalAkhir)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Buku besar berhasil digenerate", bukuBesar)
}

// GetNeracaSaldo handles GET /api/v1/laporan/neraca-saldo
func (h *LaporanHandler) GetNeracaSaldo(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tanggalPer := c.DefaultQuery("tanggalPer", "")

	neracaSaldo, err := h.laporanService.GenerateNeracaSaldo(koperasiUUID, tanggalPer)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Neraca saldo berhasil digenerate", neracaSaldo)
}

// GetTransaksiHarian handles GET /api/v1/laporan/transaksi-harian
func (h *LaporanHandler) GetTransaksiHarian(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	tanggal := c.DefaultQuery("tanggal", "") // YYYY-MM-DD, default hari ini

	laporan, err := h.laporanService.GenerateLaporanTransaksiHarian(koperasiUUID, tanggal)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Laporan transaksi harian berhasil digenerate", laporan)
}
