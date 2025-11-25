package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ServerConfig represents HTTP server configuration.
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig holds database connection parameters.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// AuthConfig stores authentication related values.
type AuthConfig struct {
	JWTSecret string
}

// WeatherConfig contains external API configuration.
type WeatherConfig struct {
	OpenWeatherAPIKey string
}

// Config aggregates all configuration values.
type Config struct {
	Server  ServerConfig
	DB      DatabaseConfig
	Auth    AuthConfig
	Weather WeatherConfig
}

// Load reads environment variables, validates them, and returns a Config.
func Load() (*Config, error) {
	// Load .env silently (useful locally); explicit env vars override.
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host: envWithDefault("SERVER_HOST", "0.0.0.0"),
			Port: intWithDefault("SERVER_PORT", 8080),
		},
		DB: DatabaseConfig{
			Host:     envWithDefault("DB_HOST", "localhost"),
			Port:     intWithDefault("DB_PORT", 5432),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     envWithDefault("DB_NAME", "abrnoc_weather"),
			SSLMode:  envWithDefault("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret: os.Getenv("JWT_SECRET"),
		},
		Weather: WeatherConfig{
			OpenWeatherAPIKey: os.Getenv("OPENWEATHER_API_KEY"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string

	if c.DB.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DB.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Auth.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if c.Weather.OpenWeatherAPIKey == "" {
		missing = append(missing, "OPENWEATHER_API_KEY")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	return nil
}

func envWithDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func intWithDefault(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	if parsed, err := strconv.Atoi(val); err == nil {
		return parsed
	}
	return def
}
