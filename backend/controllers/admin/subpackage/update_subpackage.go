package subpackage

import (
	"backend/config"
	"backend/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateSubpackage(c *gin.Context) {
	id := c.Param("id")
	var subpackage models.SubPackage

	if err := config.DB.First(&subpackage, "id_subpackage = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Subpackage tidak ditemukan",
		})
		return
	}

	updatedData := gin.H{}

	jenisPackage := c.PostForm("jenispackage")
	if jenisPackage != "" {
		subpackage.Packagetype = jenisPackage
		updatedData["jenispackage"] = jenisPackage
	}

	file, _ := c.FormFile("image")
	if file != nil {
		opened, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Gagal membuka file",
			})
			return
		}
		defer opened.Close()

		imageBytes, err := io.ReadAll(opened)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Gagal membaca file",
			})
			return
		}

		subpackage.Image = imageBytes
		updatedData["image"] = "gambar diperbarui"
	}

	if err := config.DB.Save(&subpackage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengupdate subpackage",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subpackage berhasil diperbarui",
		"data":    updatedData,
	})
}
