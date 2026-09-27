package auth

import (
	"backend/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)
func LogoutUser(c *gin.Context) {
	// Ambil token dari header Authorization
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token tidak ditemukan"})
		return
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format token tidak valid"})
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")

	// Validasi token terlebih dahulu
	_, err := middleware.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	// Tambahkan token ke blacklist (middleware)
	middleware.BlacklistToken(token)

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil. Token telah di-blacklist."})
}
