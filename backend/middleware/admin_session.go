package middleware

import (
	"backend/config"
	"backend/models"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const AdminSessionCookie = "admin_session"

func HashSessionToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func AdminSessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(AdminSessionCookie)
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak ditemukan"})
			c.Abort()
			return
		}

		var session models.AdminSession
		err = config.DB.Where("token_hash = ? AND revoked_at IS NULL", HashSessionToken(token)).First(&session).Error
		if err != nil || time.Now().After(session.ExpiresAt) {
			if err == nil {
				config.DB.Delete(&session)
			}
			ClearAdminSessionCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi telah berakhir"})
			c.Abort()
			return
		}

		var admin models.User
		if err := config.DB.Where("id_users = ? AND roleId = 2", session.UserID).First(&admin).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memvalidasi sesi"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
			}
			c.Abort()
			return
		}

		if time.Since(session.LastUsedAt) >= time.Minute {
			config.DB.Model(&session).Update("last_used_at", time.Now())
		}
		c.Set("user_id", admin.ID)
		c.Set("role_id", admin.RoleID)
		c.Set("admin", admin)
		c.Set("admin_session", session)
		c.Next()
	}
}

// StartAdminSessionCleanup removes expired/revoked rows even when the browser
// never sends another request after the session expires.
func StartAdminSessionCleanup() {
	cleanup := func() {
		config.DB.Where("expires_at <= ? OR revoked_at IS NOT NULL", time.Now()).Delete(&models.AdminSession{})
		config.DB.Where("otp_expires_at <= ? OR used_at IS NOT NULL", time.Now()).Delete(&models.PasswordReset{})
	}
	cleanup()
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			cleanup()
		}
	}()
}

func ClearAdminSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(AdminSessionCookie, "", -1, "/", "", IsProduction(), true)
}

func IsProduction() bool {
	return gin.Mode() == gin.ReleaseMode
}
