package middleware

import (
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah middleware untuk validasi JWT token
func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Token autentikasi tidak ditemukan")
			c.Abort()
			return
		}

		// Cek format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Format token tidak valid")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Token tidak valid atau sudah kadaluarsa")
			c.Abort()
			return
		}

		// Set user info ke context untuk digunakan di handler
		c.Set("idPengguna", claims.IDPengguna)
		c.Set("idKoperasi", claims.IDKoperasi)
		c.Set("namaPengguna", claims.NamaPengguna)
		c.Set("namaLengkap", claims.NamaLengkap)
		c.Set("peran", claims.Peran)

		c.Next()
	}
}

// RequireRole adalah middleware untuk memeriksa role pengguna
func RequireRole(allowedRoles ...models.PeranPengguna) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil peran dari context (sudah di-set oleh AuthMiddleware)
		peranInterface, exists := c.Get("peran")
		if !exists {
			utils.ForbiddenResponse(c, "Informasi peran tidak ditemukan")
			c.Abort()
			return
		}

		peran, ok := peranInterface.(models.PeranPengguna)
		if !ok {
			utils.ForbiddenResponse(c, "Format peran tidak valid")
			c.Abort()
			return
		}

		// Cek apakah peran pengguna termasuk dalam allowed roles
		for _, allowedRole := range allowedRoles {
			if peran == allowedRole {
				c.Next()
				return
			}
		}

		utils.ForbiddenResponse(c, "Anda tidak memiliki akses ke resource ini")
		c.Abort()
	}
}

// OptionalAuth adalah middleware untuk autentikasi opsional
// Jika token ada dan valid, set user info ke context
// Jika tidak ada atau tidak valid, lanjutkan tanpa error
func OptionalAuth(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Set user info ke context
		c.Set("idPengguna", claims.IDPengguna)
		c.Set("idKoperasi", claims.IDKoperasi)
		c.Set("namaPengguna", claims.NamaPengguna)
		c.Set("namaLengkap", claims.NamaLengkap)
		c.Set("peran", claims.Peran)

		c.Next()
	}
}
