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
	db, err := db.ConnectGorm(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&ads.Ad{},
		&ads.Country{},
		&ads.Platform{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := fiber.New()

	adRepo := ads.NewAdRepository(db)
	app.Post("/api/v1/ad", api.CreateAd(adRepo))
	app.Get("/api/v1/ad", api.ListAds(adRepo))
	log.Fatal(app.Listen(":" + cfg.Port))
}
