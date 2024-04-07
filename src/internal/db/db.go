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
		for i := 0; i < 5; i++ {
			fmt.Printf("Error connecting to database, retrying in 10 seconds: %v\n", err)
			time.Sleep(time.Second * 10)
			db, err = ConnectGorm(dbURL)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
