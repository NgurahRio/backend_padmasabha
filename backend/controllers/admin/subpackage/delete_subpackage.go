package subpackage

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteSubpackage(c *gin.Context) {
	id := c.Param("id")
	var subpackage models.SubPackage

	if err := config.DB.First(&subpackage, "id_subpackage = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Id Subpackage tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&subpackage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus subpackage",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subpackage berhasil dihapus",
	})
}
