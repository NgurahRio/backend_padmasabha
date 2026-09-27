package auth

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginUser menangani autentikasi pengguna (username/email + password)
func LoginUser(c *gin.Context) {
	var input struct {
		Identifier string `json:"identifier" binding:"required"` // username atau email
		Password   string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid. Pastikan identifier dan password terisi."})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ? OR email = ?", input.Identifier, input.Identifier).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa database"})
		return
	}

	// Periksa password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/email atau password salah"})
		return
	}

	// Generate JWT
	token, err := middleware.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"id_users": user.ID,
			"username": user.Username,
			"email":    user.Email,
			"roleId":   user.RoleID,
			"image":    user.Image,
		},
	})
}

