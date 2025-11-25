package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swag "github.com/swaggo/gin-swagger"
	_ "github.com/zarinpy/abrnoc_weather/docs"
	"github.com/zarinpy/abrnoc_weather/internals/config"
	"github.com/zarinpy/abrnoc_weather/internals/db"
	"github.com/zarinpy/abrnoc_weather/internals/handlers"
	"github.com/zarinpy/abrnoc_weather/internals/middleware"
	"github.com/zarinpy/abrnoc_weather/internals/weather"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	weatherClient := weather.NewClient(cfg.Weather.OpenWeatherAPIKey)
	authHandler := handlers.NewAuthHandler(database, cfg.Auth.JWTSecret)
	weatherHandler := handlers.NewWeatherHandler(database, weatherClient)
	authMiddleware := middleware.AuthRequired(cfg.Auth.JWTSecret)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware())

	r.GET("/swagger/*any", swag.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	r.GET("/weather", weatherHandler.GetAllWeather)
	r.GET("/weather/:id", weatherHandler.GetWeatherByID)
	r.GET("/weather/latest/:cityName", weatherHandler.GetLatestByCity)

	protected := r.Group("/weather")
	protected.Use(authMiddleware)
	{
		protected.POST("", weatherHandler.CreateWeather)
		protected.PUT("/:id", weatherHandler.UpdateWeather)
		protected.DELETE("/:id", weatherHandler.DeleteWeather)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
