package auth

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentSession(c *gin.Context) {
	adminValue, ok := c.Get("admin")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	admin := adminValue.(models.User)
	c.JSON(http.StatusOK, gin.H{"user": adminResponse(admin)})
}

func LogoutAdmin(c *gin.Context) {
	if sessionValue, ok := c.Get("admin_session"); ok {
		session := sessionValue.(models.AdminSession)
		config.DB.Delete(&session)
	}
	middleware.ClearAdminSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func LogoutAllAdminSessions(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	if err := config.DB.Where("user_id = ?", userID).Delete(&models.AdminSession{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengakhiri sesi"})
		return
	}
	middleware.ClearAdminSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Semua perangkat berhasil dikeluarkan"})
}
