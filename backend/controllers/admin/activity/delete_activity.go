package activity

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteActivity(c *gin.Context) {
	id := c.Param("id")

	var activity models.Activity

	if err := config.DB.First(&activity, "id_activity = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Activity tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity berhasil dihapus"})
}
