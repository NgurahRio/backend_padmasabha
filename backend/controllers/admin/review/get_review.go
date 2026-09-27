package review

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllReview(c *gin.Context) {
	var reviews []models.Review

	if err := config.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Destination").
		Preload("Event").
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data review"})
		return
	}

	for i := range reviews {
		reviews[i].User.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua review",
		"data":    reviews,
	})
}

func GetReviewByID(c *gin.Context) {
	id := c.Param("id")
	var review models.Review

	if err := config.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Destination").
		Preload("Event").
		First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review tidak ditemukan"})
		return
	}

	review.User.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data review",
		"data":    review,
	})
}
