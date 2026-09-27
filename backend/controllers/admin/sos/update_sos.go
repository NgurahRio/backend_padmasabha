package sos

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateSOS(c *gin.Context) {
	id := c.Param("id")

	var sos models.SOS
	if err := config.DB.First(&sos, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SOS tidak ditemukan"})
		return
	}

	var input struct {
		Name    string `json:"name_sos"`
		Alamat  string `json:"alamat_sos"`
		Telepon string `json:"telepon"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := models.SOS{
		Name:    input.Name,
		Alamat:  input.Alamat,
		Telepon: input.Telepon,
	}

	if err := config.DB.Model(&sos).Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update SOS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SOS berhasil diupdate",
		"data":    sos,
	})
}
