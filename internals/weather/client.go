package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zarinpy/abrnoc_weather/internals/models"
)

// Client defines behavior for fetching weather data.
type Client interface {
	FetchCurrentWeather(ctx context.Context, city, country string) (*models.Weather, error)
}

type openWeatherClient struct {
	apiKey string
	client *http.Client
}

// NewClient returns a weather client backed by OpenWeather API.
func NewClient(apiKey string) Client {
	return &openWeatherClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type openWeatherResponse struct {
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

func (c *openWeatherClient) FetchCurrentWeather(ctx context.Context, city, country string) (*models.Weather, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s,%s&appid=%s&units=metric",
		city,
		country,
		c.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build weather request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call weather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API responded with status %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode weather response: %w", err)
	}

	if len(data.Weather) == 0 {
		return nil, fmt.Errorf("weather data missing description")
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
