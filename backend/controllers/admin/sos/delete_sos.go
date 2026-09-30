package sos

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteSOS(c *gin.Context) {
	id := c.Param("id")

	var sos models.SOS

	if err := config.DB.First(&sos, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SOS tidak ditemukan"})
		return
	}

	// SOS is referenced by villa.sosId. Deleting it while it is still
	// in use would violate the database foreign-key constraint.
	var villaCount int64
	if err := config.DB.Model(&models.Villa{}).
		Where("sosId = ?", sos.ID).
		Count(&villaCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memeriksa penggunaan SOS",
		})
		return
	}

	if villaCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":       "SOS tidak dapat dihapus karena masih digunakan oleh villa",
			"villa_count": villaCount,
		})
		return
	}

	if err := config.DB.Delete(&sos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus SOS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SOS berhasil dihapus",
		"data":    sos,
	})
}
