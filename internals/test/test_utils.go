package test

import (
	"os"

	"github.com/zarinpy/abrnoc_weather/internals/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	testDB.AutoMigrate(&models.User{}, &models.Weather{})

	os.Setenv("JWT_SECRET", "test_secret_key_for_testing")

	return testDB
}

// CleanupTestDB closes the test database connection
func CleanupTestDB(testDB *gorm.DB) {
	sqlDB, _ := testDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// ResetDB clears all tables
func ResetDB(testDB *gorm.DB) {
	testDB.Exec("DELETE FROM users")
	testDB.Exec("DELETE FROM weathers")
}
