package review

import (
	"backend/config"
	"backend/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getUserID extracts user_id from context and normalizes to uint.
func getUserID(c *gin.Context) (uint, bool) {
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return 0, false
	}

	switch v := uidVal.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe user_id tidak dikenali"})
		return 0, false
	}
}

// AddReview membuat review baru untuk destinasi atau event.
func AddReview(c *gin.Context) {
	var input struct {
		DestinationID *uint  `json:"destinationId"`
		EventID       *uint  `json:"eventId"`
		Rating        int    `json:"rating" binding:"required"`
		Comment       string `json:"comment" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid"})
		return
	}

	// Validasi target review (wajib pilih salah satu).
	if (input.DestinationID == nil && input.EventID == nil) || (input.DestinationID != nil && input.EventID != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pilih salah satu: destinationId atau eventId"})
		return
	}

	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating harus antara 1 sampai 5"})
		return
	}

	input.Comment = strings.TrimSpace(input.Comment)
	if input.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Komentar tidak boleh kosong"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	// Pastikan target review ada agar tidak melanggar foreign key.
	if input.DestinationID != nil {
		var dest models.Destination
		if err := config.DB.First(&dest, "id_destination = ?", *input.DestinationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Destinasi tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa destinasi"})
			return
		}
	} else if input.EventID != nil {
		var event models.Event
		if err := config.DB.First(&event, "id_event = ?", *input.EventID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Event tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa event"})
			return
		}
	}

	// Cegah duplikat review oleh user pada target yang sama.
	var count int64
	check := config.DB.Model(&models.Review{}).Where("userId = ?", userID)
	if input.DestinationID != nil {
		check = check.Where("destinationId = ?", *input.DestinationID)
	} else {
		check = check.Where("eventId = ?", *input.EventID)
	}
	if err := check.Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa review"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah memberikan review untuk item ini"})
		return
	}

	review := models.Review{
		UserID:  userID,
		Rating:  input.Rating,
		Comment: input.Comment,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	if input.DestinationID != nil {
		review.DestinationID = input.DestinationID
	}
	if input.EventID != nil {
		review.EventID = input.EventID
	}

	if err := config.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review berhasil ditambahkan",
		"data":    review,
	})
}
