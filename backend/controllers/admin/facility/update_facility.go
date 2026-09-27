package facility

import (
	"backend/config"
	"backend/models"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateFacility(c *gin.Context) {
	id := c.Param("id")

	var facility models.Facility
	if err := config.DB.First(&facility, "id_facility = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Fasilitas tidak ditemukan"})
		return
	}

	nameFacility := c.PostForm("namefacility")

	if nameFacility != "" {
		facility.NameFacility = strings.TrimSpace(nameFacility)
	}

	file, err := c.FormFile("icon")
	if err == nil {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
		if !isValidImage(ext) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Tipe file tidak valid (hanya png, jpg, jpeg, webp, svg)"})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Gagal membuka file"})
			return
		}
		defer openedFile.Close()

		data, err := io.ReadAll(openedFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Gagal membaca file"})
			return
		}

		facility.Icon = data
	} else if !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Gagal mengambil file icon"})
		return
	}

	if err := config.DB.Save(&facility).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal update fasilitas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fasilitas berhasil diperbarui",
		"data": gin.H{
			"id_facility":  facility.IDFacility,
			"namefacility": facility.NameFacility,
			"hasIcon":      len(facility.Icon) > 0,
		},
	})
}
