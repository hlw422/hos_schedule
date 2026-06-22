package middleware

import (
	"hos_schedule/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

var adminRoles = map[string]bool{
	"SCHEDULER":      true,
	"HOSPITAL_ADMIN": true,
	"SUPER_ADMIN":    true,
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "No role found")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || !adminRoles[roleStr] {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
