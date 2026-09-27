package event

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllEventsUser(c *gin.Context) {
	var events []models.Event

	if err := config.DB.Preload("Destination").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua event",
		"data":    events,
	})
}

func GetEventByIDUser(c *gin.Context) {
	id := c.Param("id")
	var event models.Event

	if err := config.DB.Preload("Destination").First(&event, "id_event = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data event",
		"data":    event,
	})
}
