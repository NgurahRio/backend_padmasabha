package sos

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSOS(c *gin.Context) {
	var input struct {
		Name    string `json:"name_sos" binding:"required"`
		Alamat  string `json:"alamat_sos" binding:"required"`
		Telepon string `json:"telepon" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sos := models.SOS{
		Name:    input.Name,
		Alamat:  input.Alamat,
		Telepon: input.Telepon,
	}

	if err := config.DB.Create(&sos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat SOS"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "SOS berhasil dibuat",
		"data":    sos,
	})
}
