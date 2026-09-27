package auth

import (
	"backend/config"
	"backend/models"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mail "github.com/go-mail/mail/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	otpLifetime   = 10 * time.Minute
	tokenLifetime = 10 * time.Minute
	resendDelay   = 60 * time.Second
	maxAttempts   = 5
)

func hashSecret(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func secureDigits(length int) (string, error) {
	result, buffer := make([]byte, length), make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	for i := range result {
		result[i] = '0' + buffer[i]%10
	}
	return string(result), nil
}

func secureToken() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func sendOTPEmail(to, otp string) error {
	host, port := os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT")
	username, password := os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = username
	}
	if host == "" || port == "" || username == "" || password == "" || from == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("SMTP_PORT is invalid")
	}
	body := fmt.Sprintf("Kode OTP Anda adalah %s. Kode berlaku selama 10 menit. Jangan berikan kode ini kepada siapa pun.", otp)
	message := mail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Travora Manager - Kode reset password")
	message.SetBody("text/plain", body)
	return mail.NewDialer(host, portNumber, username, password).DialAndSend(message)
}

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak valid"})
		return
	}
	var admin models.User
	if err := config.DB.Where("LOWER(email) = ? AND roleId = 2", strings.ToLower(strings.TrimSpace(input.Email))).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email admin tidak terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa email"})
		return
	}

	var previous models.PasswordReset
	if err := config.DB.Where("user_id = ? AND used_at IS NULL", admin.ID).Order("id DESC").First(&previous).Error; err == nil {
		if remaining := resendDelay - time.Since(previous.LastSentAt); remaining > 0 {
			c.Header("Retry-After", strconv.Itoa(int(remaining.Seconds())+1))
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Tunggu sebelum meminta OTP baru"})
			return
		}
	}
	otp, err := secureDigits(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat OTP"})
		return
	}
	now := time.Now()
	reset := models.PasswordReset{UserID: admin.ID, OTPHash: hashSecret(otp), OTPExpiresAt: now.Add(otpLifetime), LastSentAt: now}
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", admin.ID).Delete(&models.PasswordReset{}).Error; err != nil {
			return err
		}
		return tx.Create(&reset).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat permintaan reset"})
		return
	}
	if err := sendOTPEmail(admin.Email, otp); err != nil {
		config.DB.Delete(&reset)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Layanan email belum tersedia"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kode OTP berhasil dikirim ke email admin."})
}

func VerifyResetOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau OTP tidak valid"})
		return
	}
	var admin models.User
	if err := config.DB.Where("LOWER(email) = ? AND roleId = 2", strings.ToLower(strings.TrimSpace(input.Email))).First(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP tidak valid atau kedaluwarsa"})
		return
	}
	var reset models.PasswordReset
	if err := config.DB.Where("user_id = ? AND used_at IS NULL", admin.ID).Order("id DESC").First(&reset).Error; err != nil || time.Now().After(reset.OTPExpiresAt) || reset.Attempts >= maxAttempts {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP tidak valid atau kedaluwarsa"})
		return
	}
	if subtle.ConstantTimeCompare([]byte(reset.OTPHash), []byte(hashSecret(input.OTP))) != 1 {
		config.DB.Model(&reset).UpdateColumn("attempts", gorm.Expr("attempts + 1"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP tidak valid atau kedaluwarsa"})
		return
	}
	token, err := secureToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token reset"})
		return
	}
	tokenHash, expires := hashSecret(token), time.Now().Add(tokenLifetime)
	if err := config.DB.Model(&reset).Updates(map[string]any{"reset_token_hash": tokenHash, "token_expires_at": expires}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi OTP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP berhasil diverifikasi", "resetToken": token})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		ResetToken string `json:"resetToken" binding:"required"`
		Password   string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 8 karakter"})
		return
	}
	var reset models.PasswordReset
	if err := config.DB.Where("reset_token_hash = ? AND used_at IS NULL", hashSecret(input.ResetToken)).First(&reset).Error; err != nil || reset.TokenExpiresAt == nil || time.Now().After(*reset.TokenExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token reset tidak valid atau kedaluwarsa"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id_users = ?", reset.UserID).Update("password", string(hashed)).Error; err != nil {
			return err
		}
		if err := tx.Delete(&reset).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ?", reset.UserID).Delete(&models.AdminSession{}).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengganti password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diganti. Silakan login kembali."})
}
