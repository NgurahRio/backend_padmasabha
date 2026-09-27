package subcategory

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSubcategory(c *gin.Context) {
	var input struct {
		CategoriesID      uint   `json:"categoriesId" binding:"required"`
		NameSubcategories string `json:"namesubcategories" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subcategory := models.Subcategory{
		CategoryID: input.CategoriesID,
		Name:       input.NameSubcategories,
	}

	if err := config.DB.Create(&subcategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan subcategory"})
		return
	}

	if err := config.DB.Preload("Category").First(&subcategory, subcategory.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data relasi kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subcategory berhasil ditambahkan",
		"data":    subcategory,
	})
}
