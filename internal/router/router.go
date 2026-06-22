package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hos_schedule/internal/config"
	"hos_schedule/internal/handler"
	"hos_schedule/internal/middleware"
	"hos_schedule/internal/pkg/response"
	"hos_schedule/internal/repository"
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

		hospitalRepo := repository.NewHospitalRepo(db)
		hospitalService := service.NewHospitalService(hospitalRepo)
		hospitalHandler := handler.NewHospitalHandler(hospitalService)

		departmentRepo := repository.NewDepartmentRepo(db)
		departmentService := service.NewDepartmentService(departmentRepo)
		departmentHandler := handler.NewDepartmentHandler(departmentService)

		api.GET("/hospitals", hospitalHandler.List)
		api.GET("/hospitals/:id", hospitalHandler.GetByID)
		api.GET("/hospitals/:id/campuses", hospitalHandler.GetCampuses)

		api.GET("/departments", departmentHandler.List)
		api.GET("/departments/:id", departmentHandler.GetByID)
	}
}
