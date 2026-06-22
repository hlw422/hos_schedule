package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hos_schedule/internal/config"
	"hos_schedule/internal/handler"
	"hos_schedule/internal/middleware"
	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/service"
)

func Register(r *gin.Engine, db *gorm.DB, rdb *redis.Client, cfg *config.Config) {
	authService := service.NewAuthService(db, cfg)
	authHandler := handler.NewAuthHandler(authService)

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			response.Success(c, gin.H{"status": "ok"})
		})

		api.POST("/auth/login", authHandler.Login)

		auth := api.Group("")
		auth.Use(middleware.Auth(cfg))
		{
			auth.GET("/me", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				response.Success(c, gin.H{"user_id": userID})
			})
		}
	}
}
