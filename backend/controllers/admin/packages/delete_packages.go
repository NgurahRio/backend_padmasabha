package packagess

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeletePackages(c *gin.Context) {
	destIDStr := c.Param("villaId")
	if destIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "villaId wajib diisi"})
		return
	}

	destID, convErr := strconv.Atoi(destIDStr)
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "villaId harus berupa angka"})
		return
	}

	var pkg models.Packages
	if err := config.DB.First(&pkg, "villaId = ?", destID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Package untuk villaId ini tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&pkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus package",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Package berhasil dihapus",
	})
}
