package middleware

import (
	"strings"

	"hos_schedule/internal/config"
	"hos_schedule/internal/pkg/jwt"
	"hos_schedule/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(cfg.JWT.Secret, tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
