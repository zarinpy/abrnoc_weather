package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zarinpy/abrnoc_weather/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Try to load .env file, but don't fail if it doesn't exist (useful for tests)
	_ = godotenv.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	if err := DB.AutoMigrate(&models.Weather{}, &models.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}
