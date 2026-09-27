package review

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteReview menghapus review milik user.
func DeleteReview(c *gin.Context) {
	reviewIDStr := c.Param("id")
	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID review tidak valid"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var review models.Review
	if err := config.DB.Where("id_review = ? AND userId = ?", reviewID, userID).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review tidak ditemukan atau bukan milik anda"})
		return
	}

	if err := config.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review berhasil dihapus",
	})
}
