package main

import (
	"fmt"
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

	err = db.AutoMigrate(&ads.Ad{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Connected to database", db)

	// // Create Fiber app
	app := fiber.New()

	// // Set up routes using your existing handlers
	adRepo := ads.NewAdRepository(db)
	app.Post("/api/v1/ad", api.CreateAd(adRepo))
	app.Get("/api/v1/ad", api.ListAds(adRepo)) // // Start the server
	log.Fatal(app.Listen(":" + cfg.Port))      // Or ":8080" for a default port
}
