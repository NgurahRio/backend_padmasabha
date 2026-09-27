package favorite

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteFavorite(c *gin.Context) {
	// Ambil ID dari path; gunakan query ?type=event untuk hapus event bookmark, default: destination
	idStr := c.Param("destinationId")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID bookmark diperlukan"})
		return
	}

	idVal, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	targetID := uint(idVal)

	kind := c.DefaultQuery("type", "destination")

	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	var userID uint
	switch v := uidVal.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe user_id tidak dikenali"})
		return
	}

	// Hapus record favorite sesuai jenis
	db := config.DB.Where("userId = ?", userID)
	if kind == "event" {
		db = db.Where("eventId = ?", targetID)
	} else {
		db = db.Where("destinationId = ?", targetID)
	}

	if err := db.Delete(&models.Favorite{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite dihapus"})
}
