package router

import (
	"github.com/gin-gonic/gin"
	"hos_schedule/internal/pkg/response"
)

func Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			response.Success(c, gin.H{"status": "ok"})
		})
	}
}
