package review

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ✅ DELETE review berdasarkan ID
func DeleteReview(c *gin.Context) {
	id := c.Param("id")
	var review models.Review

	// Cek apakah review ada
	if err := config.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review tidak ditemukan"})
		return
	}

	// Hapus review
	if err := config.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review berhasil dihapus",
		"data":    review,
	})
}
