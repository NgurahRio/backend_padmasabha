package subpackage

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllSubpackagesUser(c *gin.Context) {
	var subpackages []models.SubPackage

	if err := config.DB.Find(&subpackages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data subpackage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua subpackage",
		"data":    subpackages,
	})
}

func GetSubpackageByIDUser(c *gin.Context) {
	id := c.Param("id")
	var subpackage models.SubPackage

	if err := config.DB.First(&subpackage, "id_subpackage = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "ID Subpackage tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data subpackage",
		"data":    subpackage,
	})
}
