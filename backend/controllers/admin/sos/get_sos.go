package sos

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllSOS(c *gin.Context) {
	var sos []models.SOS

	if err := config.DB.Find(&sos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data SOS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua data SOS",
		"data":    sos,
	})
}

func GetSOSByID(c *gin.Context) {
	id := c.Param("id")
	var data models.SOS

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SOS tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data SOS",
		"data":    data,
	})
}