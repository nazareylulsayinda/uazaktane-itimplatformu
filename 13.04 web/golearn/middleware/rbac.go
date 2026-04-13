package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TeacherOnly middleware ensures the user is a teacher
func TeacherOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: teacher access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
