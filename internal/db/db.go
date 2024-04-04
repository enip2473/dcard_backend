package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectGorm initializes a database connection using GORM
func ConnectGorm(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Optional: Configure GORM settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (adjust these based on your needs)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
