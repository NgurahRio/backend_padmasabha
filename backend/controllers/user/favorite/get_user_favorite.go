package favorite

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserFavorites(c *gin.Context) {
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

	favorites := make([]models.Favorite, 0)
	if err := config.DB.Preload("User").Preload("Destination").Preload("Event").Where("userId = ?", userID).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil favorites"})
		return
	}

	// Sanitasi data user agar password tidak ikut terbawa ke response
	for i := range favorites {
		if favorites[i].User != nil {
			orig := favorites[i].User
			// Buat salinan tanpa field sensitif
			favorites[i].User = &models.User{
				ID:       orig.ID,
				Username: orig.Username,
				RoleID:   orig.RoleID,
				Image:    orig.Image,
			}
		}
	}

	destSet := make(map[uint]struct{})
	eventSet := make(map[uint]struct{})

	for _, f := range favorites {
		if f.DestinationID != nil {
			destSet[*f.DestinationID] = struct{}{}
		}
		if f.EventID != nil {
			eventSet[*f.EventID] = struct{}{}
		}
	}

	destIDs := make([]uint, 0, len(destSet))
	for id := range destSet {
		destIDs = append(destIDs, id)
	}

	eventIDs := make([]uint, 0, len(eventSet))
	for id := range eventSet {
		eventIDs = append(eventIDs, id)
	}

	destinations := make([]models.Destination, 0)
	if len(destIDs) > 0 {
		if err := config.DB.Where("id_destination IN ?", destIDs).Find(&destinations).Error; err != nil {
			fmt.Println("warning: gagal mengambil detail destinasi:", err)
		}
	}

	events := make([]models.Event, 0)
	if len(eventIDs) > 0 {
		if err := config.DB.Preload("Destination").Where("id_event IN ?", eventIDs).Find(&events).Error; err != nil {
			fmt.Println("warning: gagal mengambil detail event:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"favorites":     favorites,
		"destinations":  destinations,
		"events":        events,
		"totalFavorite": len(favorites),
	})
}
