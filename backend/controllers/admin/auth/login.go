package auth

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginAdmin(c *gin.Context) {
	var input struct {
		LoginID        string `json:"loginId" binding:"required"`
		Password       string `json:"password" binding:"required"`
		RememberDevice bool   `json:"rememberDevice"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan password wajib diisi"})
		return
	}

	var admin models.User
	loginID := strings.TrimSpace(input.LoginID)
	if err := config.DB.Where("(username = ? OR LOWER(email) = ?) AND roleId = 2", loginID, strings.ToLower(loginID)).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat sesi"})
		return
	}
	token := hex.EncodeToString(rawToken)
	now := time.Now()
	duration := 24 * time.Hour
	maxAge := 0 // Session cookie: hilang saat browser ditutup.
	if input.RememberDevice {
		duration = 7 * 24 * time.Hour
		maxAge = int(duration.Seconds())
	}
	session := models.AdminSession{
		UserID: admin.ID, TokenHash: middleware.HashSessionToken(token),
		UserAgent: c.Request.UserAgent(), IPAddress: c.ClientIP(),
		ExpiresAt: now.Add(duration), LastUsedAt: now,
	}
	if err := config.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan sesi"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(middleware.AdminSessionCookie, token, maxAge, "/", "", middleware.IsProduction(), true)
	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "user": adminResponse(admin), "expiresAt": session.ExpiresAt})
}

func adminResponse(admin models.User) gin.H {
	return gin.H{"id_users": admin.ID, "username": admin.Username, "email": admin.Email, "roleId": admin.RoleID}
}
