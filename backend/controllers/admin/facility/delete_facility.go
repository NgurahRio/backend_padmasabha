package facility

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteFacility(c *gin.Context) {
	id := c.Param("id")

	var facility models.Facility

	if err := config.DB.First(&facility, "id_facility = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Fasilitas tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&facility).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus fasilitas",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fasilitas berhasil dihapus",
	})
}
