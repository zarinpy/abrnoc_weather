package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zarinpy/abrnoc_weather/internals/models"
)

type OpenWeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
	} `json:"sys"`
}

func init() {
	// Try to load .env file, but don't fail if it doesn't exist (useful for tests)
	_ = godotenv.Load()
}

func FetchCurrentWeather(city, country string) (*models.Weather, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s,%s&appid=%s&units=metric",
		city,
		country,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("city not found or API error", err)
	}
	defer resp.Body.Close()

	var data OpenWeatherResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&data)
	if decodeErr != nil {
		return nil, decodeErr
	}

	weather := &models.Weather{
		CityName:    data.Name,
		Country:     data.Sys.Country,
		Temperature: data.Main.Temp,
		Description: data.Weather[0].Description,
		Humidity:    data.Main.Humidity,
		WindSpeed:   data.Wind.Speed,
		FetchedAt:   time.Now(),
	}

	return weather, nil
}
