package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zarinpy/abrnoc_weather/internals/models"
	"github.com/zarinpy/abrnoc_weather/internals/test"
	"gorm.io/gorm"
)

type stubWeatherClient struct{}

func (stubWeatherClient) FetchCurrentWeather(_ context.Context, city, country string) (*models.Weather, error) {
	return &models.Weather{
		CityName:    city,
		Country:     country,
		Temperature: 20,
		Description: "stub",
		Humidity:    50,
		WindSpeed:   1.0,
		FetchedAt:   time.Now(),
	}, nil
}

func setupWeatherTestRouter(t *testing.T) (*gin.Engine, *WeatherHandler) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	testDB := test.SetupTestDB()
	handler := NewWeatherHandler(testDB, stubWeatherClient{})

	router := gin.New()
	router.GET("/weather", handler.GetAllWeather)

	return router, handler
}

func createTestWeatherRecords(dbInstance *gorm.DB, count int) {
	for i := 0; i < count; i++ {
		record := models.Weather{
			ID:          uuid.New(),
			CityName:    "TestCity",
			Country:     "TC",
			Temperature: 20.5 + float64(i),
			Description: "test description",
			Humidity:    65,
			WindSpeed:   3.2,
			FetchedAt:   time.Now(),
			CreatedAt:   time.Now().Add(time.Duration(i) * time.Second),
		}
		dbInstance.Create(&record)
	}
}

func TestGetAllWeather_DefaultPagination(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	createTestWeatherRecords(handler.db, 15)

	req, _ := http.NewRequest("GET", "/weather", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedWeatherResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, int64(15), response.Total)
	assert.Equal(t, 2, response.TotalPages)
	assert.Len(t, response.Data, 10) // Should return 10 records (default limit)
}

func TestGetAllWeather_WithPageAndLimit(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	createTestWeatherRecords(handler.db, 25)

	req, _ := http.NewRequest("GET", "/weather?page=2&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedWeatherResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, int64(25), response.Total)
	assert.Equal(t, 3, response.TotalPages)
	assert.Len(t, response.Data, 10) // Should return 10 records
}

func TestGetAllWeather_LastPage(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	createTestWeatherRecords(handler.db, 15)

	req, _ := http.NewRequest("GET", "/weather?page=2&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedWeatherResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, int64(15), response.Total)
	assert.Equal(t, 2, response.TotalPages)
	assert.Len(t, response.Data, 5) // Should return remaining 5 records
}

func TestGetAllWeather_EmptyResult(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	req, _ := http.NewRequest("GET", "/weather", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedWeatherResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, int64(0), response.Total)
	assert.Equal(t, 1, response.TotalPages) // Should be 1 even with 0 records
	assert.Len(t, response.Data, 0)
}

func TestGetAllWeather_MaxLimit(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	createTestWeatherRecords(handler.db, 150)

	req, _ := http.NewRequest("GET", "/weather?page=1&limit=200", nil) // Requesting more than max
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedWeatherResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 100, response.Limit) // Should be capped at 100
	assert.Len(t, response.Data, 100)
}

func TestGetAllWeather_InvalidPage(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	req, _ := http.NewRequest("GET", "/weather?page=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllWeather_InvalidLimit(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	req, _ := http.NewRequest("GET", "/weather?limit=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllWeather_NegativePage(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	req, _ := http.NewRequest("GET", "/weather?page=-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllWeather_NegativeLimit(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	req, _ := http.NewRequest("GET", "/weather?limit=-5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllWeather_Ordering(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	createTestWeatherRecords(handler.db, 5)

	req, _ := http.NewRequest("GET", "/weather?limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedWeatherResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Len(t, response.Data, 5)

	// Verify records are ordered by created_at DESC (newest first)
	for i := 0; i < len(response.Data)-1; i++ {
		assert.True(t, response.Data[i].CreatedAt.After(response.Data[i+1].CreatedAt) ||
			response.Data[i].CreatedAt.Equal(response.Data[i+1].CreatedAt))
	}
}

func TestGetAllWeather_CustomLimit(t *testing.T) {
	router, handler := setupWeatherTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	createTestWeatherRecords(handler.db, 20)

	req, _ := http.NewRequest("GET", "/weather?page=1&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedWeatherResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 5, response.Limit)
	assert.Equal(t, int64(20), response.Total)
	assert.Equal(t, 4, response.TotalPages) // 20 records / 5 per page = 4 pages
	assert.Len(t, response.Data, 5)
}
