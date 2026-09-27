package destination

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteDestination(c *gin.Context) {
	id := c.Param("id")

	var destination models.Destination

	if err := config.DB.First(&destination, "id_destination = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Destinasi tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&destination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus destinasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destinasi berhasil dihapus"})
}
