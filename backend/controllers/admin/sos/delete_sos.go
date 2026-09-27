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

	// SOS is referenced by destination.sosId. Deleting it while it is still
	// in use would violate the database foreign-key constraint.
	var destinationCount int64
	if err := config.DB.Model(&models.Destination{}).
		Where("sosId = ?", sos.ID).
		Count(&destinationCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memeriksa penggunaan SOS",
		})
		return
	}

	if destinationCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":             "SOS tidak dapat dihapus karena masih digunakan oleh destinasi",
			"destination_count": destinationCount,
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
