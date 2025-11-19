package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PortalAnggotaHandler menangani endpoint portal anggota
type PortalAnggotaHandler struct {
	portalService *services.PortalAnggotaService
}

// NewPortalAnggotaHandler membuat instance baru PortalAnggotaHandler
func NewPortalAnggotaHandler(portalService *services.PortalAnggotaService) *PortalAnggotaHandler {
	return &PortalAnggotaHandler{
		portalService: portalService,
	}
}

// Login handles POST /api/v1/portal/login
// @Summary Login portal anggota
// @Description Login anggota dengan nomor anggota dan PIN
// @Tags Portal Anggota
// @Accept json
// @Produce json
// @Param idKoperasi query string true "ID Koperasi"
// @Param body body services.LoginAnggotaRequest true "Login credentials"
// @Success 200 {object} utils.APIResponse{data=services.LoginAnggotaResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/portal/login [post]
func (h *PortalAnggotaHandler) Login(c *gin.Context) {
	// Get ID koperasi from query parameter
	idKoperasiStr := c.Query("idKoperasi")
	if idKoperasiStr == "" {
		utils.BadRequestResponse(c, "ID koperasi diperlukan")
		return
	}

	idKoperasi, err := uuid.Parse(idKoperasiStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID koperasi tidak valid")
		return
	}

	var req services.LoginAnggotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Login anggota
	response, err := h.portalService.LoginAnggota(idKoperasi, req.NomorAnggota, req.PIN)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", response)
}

// GetProfile handles GET /api/v1/portal/profile
// @Summary Get member profile
// @Description Mendapatkan profil anggota yang sedang login
// @Tags Portal Anggota
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=models.AnggotaResponse}
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/portal/profile [get]
func (h *PortalAnggotaHandler) GetProfile(c *gin.Context) {
	// Get ID anggota from context (set by auth middleware)
	idAnggota, exists := c.Get("idAnggota")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	idKoperasi, exists := c.Get("idKoperasi")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	anggotaUUID := idAnggota.(uuid.UUID)
	koperasiUUID := idKoperasi.(uuid.UUID)

	// Get profile
	profile, err := h.portalService.GetInfoAnggota(koperasiUUID, anggotaUUID)
	if err != nil {
		utils.NotFoundResponse(c, "Profil tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profil berhasil diambil", profile)
}

// GetSaldo handles GET /api/v1/portal/saldo
// @Summary Get member balance
// @Description Mendapatkan saldo simpanan anggota
// @Tags Portal Anggota
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=models.SaldoSimpananAnggota}
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/portal/saldo [get]
func (h *PortalAnggotaHandler) GetSaldo(c *gin.Context) {
	// Get ID anggota from context
	idAnggota, exists := c.Get("idAnggota")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	idKoperasi, exists := c.Get("idKoperasi")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	anggotaUUID := idAnggota.(uuid.UUID)
	koperasiUUID := idKoperasi.(uuid.UUID)

	// Get balance
	saldo, err := h.portalService.GetSaldoAnggota(koperasiUUID, anggotaUUID)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Saldo berhasil diambil", saldo)
}

// GetRiwayat handles GET /api/v1/portal/riwayat
// @Summary Get transaction history
// @Description Mendapatkan riwayat transaksi simpanan anggota
// @Tags Portal Anggota
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} utils.APIResponse{data=[]services.RiwayatTransaksiAnggota}
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/portal/riwayat [get]
func (h *PortalAnggotaHandler) GetRiwayat(c *gin.Context) {
	// Get ID anggota from context
	idAnggota, exists := c.Get("idAnggota")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	idKoperasi, exists := c.Get("idKoperasi")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	anggotaUUID := idAnggota.(uuid.UUID)
	koperasiUUID := idKoperasi.(uuid.UUID)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get transaction history
	riwayat, total, err := h.portalService.GetRiwayatTransaksi(koperasiUUID, anggotaUUID, pageSize, offset)
	if err != nil {
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	// Calculate pagination metadata
	pagination := utils.CalculatePaginationMeta(page, pageSize, total)

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Riwayat transaksi berhasil diambil", riwayat, pagination)
}

// UbahPIN handles PUT /api/v1/portal/ubah-pin
// @Summary Change member PIN
// @Description Ubah PIN portal anggota
// @Tags Portal Anggota
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body struct{PINLama string; PINBaru string} true "Change PIN request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/portal/ubah-pin [put]
func (h *PortalAnggotaHandler) UbahPIN(c *gin.Context) {
	// Get ID anggota from context
	idAnggota, exists := c.Get("idAnggota")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	idKoperasi, exists := c.Get("idKoperasi")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	anggotaUUID := idAnggota.(uuid.UUID)
	koperasiUUID := idKoperasi.(uuid.UUID)

	var req struct {
		PINLama string `json:"pinLama" binding:"required,len=6"`
		PINBaru string `json:"pinBaru" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "PIN harus 6 digit")
		return
	}

	// Change PIN
	err := h.portalService.UbahPIN(koperasiUUID, anggotaUUID, req.PINLama, req.PINBaru)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "PIN berhasil diubah", nil)
}
