package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler menangani endpoint otentikasi
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/v1/auth/login
// @Summary Login pengguna
// @Description Login dengan username dan password, mengembalikan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body services.LoginRequest true "Login credentials"
// @Success 200 {object} utils.APIResponse{data=services.LoginResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest

	// Bind dan validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Panggil service login
	response, err := h.authService.Login(req.NamaPengguna, req.KataSandi)
	if err != nil {
		utils.UnauthorizedResponse(c, "Nama pengguna atau kata sandi salah")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", response)
}

// GetProfile handles GET /api/v1/auth/profile
// @Summary Get user profile
// @Description Mendapatkan profil pengguna yang sedang login
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=models.PenggunaResponse}
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Ambil ID pengguna dari context (di-set oleh auth middleware)
	idPengguna, exists := c.Get("idPengguna")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	penggunaUUID := idPengguna.(uuid.UUID)

	// Get profil pengguna
	profil, err := h.authService.GetProfilPengguna(penggunaUUID)
	if err != nil {
		utils.NotFoundResponse(c, "Pengguna tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profil berhasil diambil", profil)
}

// ChangePassword handles PUT /api/v1/auth/change-password
// @Summary Ubah password
// @Description Ubah password pengguna yang sedang login
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.UbahKataSandiRequest true "Change password request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Ambil ID pengguna dari context
	idPengguna, exists := c.Get("idPengguna")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	penggunaUUID := idPengguna.(uuid.UUID)

	var req services.UbahKataSandiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Ubah password
	err := h.authService.UbahKataSandiByUUID(penggunaUUID, req.KataSandiLama, req.KataSandiBaru)
	if err != nil {
		// Check jenis error
		if err.Error() == "kata sandi lama tidak sesuai" {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.SafeInternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Kata sandi berhasil diubah", nil)
}

// RefreshToken handles POST /api/v1/auth/refresh
// @Summary Refresh JWT token
// @Description Generate token baru dengan memperpanjang expiration
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=services.LoginResponse}
// @Failure 401 {object} utils.APIResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Ambil ID pengguna dari context
	idPengguna, exists := c.Get("idPengguna")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return
	}

	penggunaUUID := idPengguna.(uuid.UUID)

	// Generate token baru
	response, err := h.authService.RefreshTokenByUUID(penggunaUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Gagal generate token baru", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token berhasil di-refresh", response)
}

// Logout handles POST /api/v1/auth/logout
// @Summary Logout pengguna
// @Description Logout pengguna (client-side token removal)
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Note: Karena menggunakan JWT stateless, logout dilakukan di client-side
	// dengan menghapus token dari storage. Handler ini hanya konfirmasi.

	utils.SuccessResponse(c, http.StatusOK, "Logout berhasil. Silakan hapus token dari client", nil)
}
