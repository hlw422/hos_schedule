package main

import (
	"fmt"
	"log"

	"hos_schedule/internal/config"
	"hos_schedule/internal/middleware"
	"hos_schedule/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()
	r.Use(middleware.CORS())

	router.Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
