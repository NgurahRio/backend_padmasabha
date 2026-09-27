package user

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID       uint   `json:"id_users"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   uint   `json:"roleId"`
	RoleName string `json:"roleName"`
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			RoleID:   u.RoleID,
			RoleName: u.Role.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua pengguna",
		"data":    userResponses,
	})
}


func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		RoleName: user.Role.Name,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data pengguna",
		"data":    response,
	})
}

