package event

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	var event models.Event

	if err := config.DB.First(&event, "id_event = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event tidak ditemukan",
		})
		return
	}

	if err := config.DB.Delete(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event berhasil dihapus",
	})
}
