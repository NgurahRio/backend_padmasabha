package subcategory

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ✅ GET semua subcategory
func GetAllSubcategoriesUser(c *gin.Context) {
	var subcategories []models.Subcategory

	if err := config.DB.Preload("Category").Find(&subcategories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data subcategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua subcategory",
		"data":    subcategories,
	})
}

// ✅ GET subcategory by ID
func GetSubcategoryByIDUser(c *gin.Context) {
	id := c.Param("id")
	var subcategory models.Subcategory

	if err := config.DB.Preload("Category").First(&subcategory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data subcategory",
		"data":    subcategory,
	})
}
