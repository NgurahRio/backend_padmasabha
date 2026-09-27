package review

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetReviews mengambil review, bisa difilter berdasarkan destinationId atau eventId.
func GetReviews(c *gin.Context) {
	destIDStr := c.Query("destinationId")
	eventIDStr := c.Query("eventId")

	var (
		destID uint64
		evID   uint64
		err    error
	)

	if destIDStr != "" {
		destID, err = strconv.ParseUint(destIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "destinationId tidak valid"})
			return
		}
	}

	if eventIDStr != "" {
		evID, err = strconv.ParseUint(eventIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventId tidak valid"})
			return
		}
	}

	// jangan izinkan filter ganda
	if destIDStr != "" && eventIDStr != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gunakan salah satu filter saja: destinationId atau eventId"})
		return
	}

	db := config.DB.Model(&models.Review{})
	if destIDStr != "" {
		db = db.Where("destinationId = ?", destID)
	}
	if eventIDStr != "" {
		db = db.Where("eventId = ?", evID)
	}

	var reviews []models.Review
	if err := db.
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Destination").
		Preload("Event").
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data review",
		"data":    reviews,
	})
}