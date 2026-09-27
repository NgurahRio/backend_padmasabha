package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleOnly(roles ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleVal, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: no role info"})
			c.Abort()
			return
		}

		userRole := roleVal.(uint)

		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role"})
		c.Abort()
	}
}
