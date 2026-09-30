package restaurant

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteRestaurant(c *gin.Context) {
	id := c.Param("id")

	var restaurant models.Restaurant

	if err := config.DB.First(&restaurant, "id_restaurant = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Restaurant tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant berhasil dihapus"})
}
