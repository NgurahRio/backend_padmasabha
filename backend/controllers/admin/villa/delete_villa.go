package villa

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteVilla(c *gin.Context) {
	id := c.Param("id")

	var villa models.Villa

	if err := config.DB.First(&villa, "id_villa = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Villa tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&villa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus villa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Villa berhasil dihapus"})
}
