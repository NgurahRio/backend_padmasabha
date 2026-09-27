package facility

import (
	"backend/config"
	"backend/models"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllFacility(c *gin.Context) {
	var facilities []models.Facility
	if err := config.DB.Find(&facilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data fasilitas"})
		return
	}

	var result []gin.H
	for _, f := range facilities {

		iconBase64 := ""
		if len(f.Icon) > 0 {
			iconBase64 = base64.StdEncoding.EncodeToString(f.Icon)
		}

		result = append(result, gin.H{
			"id_facility":  f.IDFacility,
			"namefacility": f.NameFacility,
			"icon":         iconBase64,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data fasilitas berhasil diambil",
		"data":    result,
	})
}

func GetFacilityByID(c *gin.Context) {
	id := c.Param("id")

	facilityID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	var facility models.Facility
	if err := config.DB.First(&facility, "id_facility = ?", facilityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Fasilitas tidak ditemukan"})
		return
	}

	iconBase64 := ""
	if len(facility.Icon) > 0 {
		iconBase64 = base64.StdEncoding.EncodeToString(facility.Icon)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail fasilitas berhasil diambil",
		"data": gin.H{
			"id_facility":  facility.IDFacility,
			"namefacility": facility.NameFacility,
			"icon":         iconBase64,
		},
	})
}
