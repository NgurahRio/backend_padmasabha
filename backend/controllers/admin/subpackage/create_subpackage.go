package subpackage

import (
	"backend/config"
	"backend/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSubpackage(c *gin.Context) {
	jenisPackage := c.PostForm("jenispackage")

	if jenisPackage == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jenispackage wajib diisi"})
		return
	}

	file, _ := c.FormFile("image")

	var imageBytes []byte

	if file != nil {
		opened, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka file"})
			return
		}
		defer opened.Close()

		imageBytes, err = io.ReadAll(opened)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file"})
			return
		}
	}

	subpackage := models.SubPackage{
		Packagetype: jenisPackage,
		Image:       imageBytes, 
	}

	if err := config.DB.Create(&subpackage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan subpackage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subpackage berhasil ditambahkan",
		"data":    subpackage,
	})
}
