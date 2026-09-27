package favorite

import (
	"backend/config"
	"backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddFavorite(c *gin.Context) {
	var input struct {
		DestinationID *uint `json:"destinationId"`
		EventID       *uint `json:"eventId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid. Kirim JSON yang benar."})
		return
	}

	if input.DestinationID == nil && input.EventID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "destinationId atau eventId wajib diisi (minimal salah satu)."})
		return
	}

	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	var userID uint
	switch v := uidVal.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe user_id tidak dikenali"})
		return
	}

	if input.DestinationID != nil {
		var dest models.Destination
		if err := config.DB.First(&dest, "id_destination = ?", *input.DestinationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "destinationId tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal validasi destination"})
			return
		}
	}

	if input.EventID != nil {
		var ev models.Event
		if err := config.DB.First(&ev, "id_event = ?", *input.EventID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "eventId tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal validasi event"})
			return
		}
	}

	var existing models.Favorite
	q := config.DB.Where("userId = ?", userID)

	if input.DestinationID == nil {
		q = q.Where("destinationId IS NULL")
	} else {
		q = q.Where("destinationId = ?", *input.DestinationID)
	}

	if input.EventID == nil {
		q = q.Where("eventId IS NULL")
	} else {
		q = q.Where("eventId = ?", *input.EventID)
	}

	err := q.First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Item sudah ada di favorite"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan database saat cek favorite"})
		return
	}

	fav := models.Favorite{
		UserID:        userID,
		DestinationID: input.DestinationID,
		EventID:       input.EventID,
	}

	if err := config.DB.Create(&fav).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan favorite"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Favorite ditambahkan",
		"favorite": fav,
	})
}
