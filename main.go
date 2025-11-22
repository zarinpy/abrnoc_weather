package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swag "github.com/swaggo/gin-swagger"
	_ "github.com/zarinpy/abrnoc_weather/docs"
	"github.com/zarinpy/abrnoc_weather/internals/db"
	"github.com/zarinpy/abrnoc_weather/internals/handlers"
	"github.com/zarinpy/abrnoc_weather/internals/middleware"
)

// @title Cloudzy Weather API
// @version 1.0
// @description A RESTful weather API service with authentication
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	db.Connect()

	r := gin.Default()

	// Enable CORS
	r.Use(middleware.CORSMiddleware())

	r.GET("/swagger/*any", swag.WrapHandler(swaggerFiles.Handler))

	// Auth routes (public)
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	// Public routes
	r.GET("/weather", handlers.GetAllWeather)
	r.GET("/weather/:id", handlers.GetWeatherByID)
	r.GET("/weather/latest/:cityName", handlers.GetLatestByCity)

	// Protected routes
	protected := r.Group("/weather")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("", handlers.CreateWeather)
		protected.PUT("/:id", handlers.UpdateWeather)
		protected.DELETE("/:id", handlers.DeleteWeather)
	}

	r.Run("0.0.0.0:8080")
}
