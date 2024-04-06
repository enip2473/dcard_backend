package main

import (
	"log"

	"dcard_backend/api"
	"dcard_backend/internal/config"
	"dcard_backend/internal/db"
	"dcard_backend/pkg/ads"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database connection
	database, err := db.ConnectGorm(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = database.AutoMigrate(
		&ads.Ad{},
		&ads.Country{},
		&ads.Platform{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = db.InitRedis(cfg.RedisAddr, cfg.RedisPass)

	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	app := fiber.New()

	adRepo := ads.NewAdRepository(database)
	app.Post("/api/v1/ad", api.CreateAd(adRepo))
	app.Get("/api/v1/ad", api.ListAds(adRepo))
	log.Fatal(app.Listen(":" + cfg.Port))
}
