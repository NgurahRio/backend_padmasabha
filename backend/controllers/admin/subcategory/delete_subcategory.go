package subcategory

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteSubcategory(c *gin.Context) {
	id := c.Param("id")
	var subcategory models.Subcategory

	if err := config.DB.First(&subcategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&subcategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus subcategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus subcategory",
		"data":    subcategory,
	})
}