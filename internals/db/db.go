package db

import (
	"fmt"

	"github.com/zarinpy/abrnoc_weather/internals/config"
	"github.com/zarinpy/abrnoc_weather/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates and configures a new database connection.
func New(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("enable uuid extension: %w", err)
	}

	if err := db.AutoMigrate(&models.Weather{}, &models.User{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil
}
