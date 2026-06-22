package main

import (
	"context"
	"fmt"
	"log"

	"hos_schedule/internal/config"
	"hos_schedule/internal/middleware"
	"hos_schedule/internal/model"
	"hos_schedule/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := cfg.Database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Hospital{},
		&model.Campus{},
		&model.Department{},
		&model.Doctor{},
		&model.Schedule{},
		&model.Patient{},
		&model.Appointment{},
		&model.Notification{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	rdb := cfg.Redis.Connect()
	if err := cfg.Redis.Ping(context.Background(), rdb); err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()
	r.Use(middleware.CORS())

	router.Register(r, db, rdb, cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
