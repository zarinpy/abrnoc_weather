package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zarinpy/abrnoc_weather/internals/db"
	"github.com/zarinpy/abrnoc_weather/internals/models"
	"github.com/zarinpy/abrnoc_weather/internals/weather"
)

// PaginatedWeatherResponse represents a paginated response
type PaginatedWeatherResponse struct {
	Data       []models.Weather `json:"data"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"totalPages"`
}

// @Summary Get all weather records
// @Description Retrieve paginated weather records from the database
// @Tags weather
// @Produce json
// @Param page query int false "Page number (default: 1)" default(1) minimum(1)
// @Param limit query int false "Number of records per page (default: 10, max: 100)" default(10) minimum(1) maximum(100)
// @Success 200 {object} PaginatedWeatherResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /weather [get]
func GetAllWeather(c *gin.Context) {
	// Get pagination parameters from query string
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(400, gin.H{"error": "Invalid limit"})
		return
	}

	// Maximum limit to prevent excessive queries
	if limit > 100 {
		limit = 100
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	var total int64
	db.DB.Model(&models.Weather{}).Count(&total)

	// Get paginated results
	var weathers []models.Weather
	result := db.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&weathers)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch weather records"})
		return
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	response := PaginatedWeatherResponse{
		Data:       weathers,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(200, response)
}

// @Summary Get latest weather for city
// @Description Get the most recent weather record for a specific city
// @Tags weather
// @Produce json
// @Param cityName path string true "City Name" example:"London"
// @Success 200 {object} models.Weather
// @Failure 404 {object} map[string]string "No record found for this city"
// @Router /weather/latest/{cityName} [get]
func GetLatestByCity(c *gin.Context) {
	city := c.Param("cityName")
	var weather models.Weather
	result := db.DB.Order("fetched_at desc").Where("city_name = ?", city).First(&weather)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "No record found for this city"})
		return
	}
	c.JSON(200, weather)
}

// CreateWeatherRequest represents the request body for creating weather
type CreateWeatherRequest struct {
	CityName string `json:"cityName" example:"London" binding:"required"`
	Country  string `json:"country" example:"GB" binding:"required"`
}

// @Summary Create new weather entry
// @Description Fetch current weather from OpenWeatherMap API and store it in the database
// @Tags weather
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param weather body CreateWeatherRequest true "Weather creation request"
// @Success 201 {object} models.Weather
// @Failure 400 {object} map[string]string "Bad request or failed to fetch weather"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /weather [post]
func CreateWeather(c *gin.Context) {
	var input CreateWeatherRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	weather, err := weather.FetchCurrentWeather(input.CityName, input.Country)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to fetch weather: " + err.Error()})
		return
	}

	db.DB.Create(weather)
	c.JSON(201, weather)
}

// @Summary Get weather by ID
// @Description Retrieve a specific weather record by its ID
// @Tags weather
// @Produce json
// @Param id path string true "Weather ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} models.Weather
// @Failure 404 {object} map[string]string "Weather record not found"
// @Router /weather/{id} [get]
func GetWeatherByID(c *gin.Context) {
	id := c.Param("id")
	var weather models.Weather
	result := db.DB.Where("id = ?", id).First(&weather)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Weather record not found"})
		return
	}
	c.JSON(200, weather)
}

// UpdateWeatherRequest represents the request body for updating weather
type UpdateWeatherRequest struct {
	CityName    string  `json:"cityName" example:"London"`
	Country     string  `json:"country" example:"GB"`
	Temperature float64 `json:"temperature" example:"15.5"`
	Description string  `json:"description" example:"clear sky"`
	Humidity    int     `json:"humidity" example:"65"`
	WindSpeed   float64 `json:"windSpeed" example:"3.2"`
}

// @Summary Update weather entry
// @Description Update an existing weather record
// @Tags weather
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Weather ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param weather body UpdateWeatherRequest true "Weather update request"
// @Success 200 {object} models.Weather
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Weather record not found"
// @Router /weather/{id} [put]
func UpdateWeather(c *gin.Context) {
	id := c.Param("id")
	var weather models.Weather
	result := db.DB.Where("id = ?", id).First(&weather)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Weather record not found"})
		return
	}

	var input UpdateWeatherRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.CityName != "" {
		weather.CityName = input.CityName
	}
	if input.Country != "" {
		weather.Country = input.Country
	}
	if input.Temperature != 0 {
		weather.Temperature = input.Temperature
	}
	if input.Description != "" {
		weather.Description = input.Description
	}
	if input.Humidity != 0 {
		weather.Humidity = input.Humidity
	}
	if input.WindSpeed != 0 {
		weather.WindSpeed = input.WindSpeed
	}

	db.DB.Save(&weather)
	c.JSON(200, weather)
}

// DeleteWeatherResponse represents the delete response
type DeleteWeatherResponse struct {
	Message string `json:"message" example:"Weather record deleted successfully"`
}

// @Summary Delete weather entry
// @Description Delete a weather record by ID
// @Tags weather
// @Security BearerAuth
// @Produce json
// @Param id path string true "Weather ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} DeleteWeatherResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Weather record not found"
// @Router /weather/{id} [delete]
func DeleteWeather(c *gin.Context) {
	id := c.Param("id")
	result := db.DB.Where("id = ?", id).Delete(&models.Weather{})
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Weather record not found"})
		return
	}
	c.JSON(200, DeleteWeatherResponse{Message: "Weather record deleted successfully"})
}
