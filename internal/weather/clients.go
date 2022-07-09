package weather

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type WeatherStackClient struct {
	client *resty.Client
	apiKey string
}

func NewWeatherStackClient(apiKey string) *WeatherStackClient {
	return &WeatherStackClient{
		client: newClient("http://api.weatherstack.com/").
			SetQueryParam("access_key", apiKey).
			SetQueryParam("units", "m"), // Celsius
		apiKey: apiKey,
	}
}

type WeatherStackResponse struct {
	Current struct {
		WindSpeed   int `json:"wind_speed"`
		Temperature int `json:"temperature"`
	} `json:"current"`
}

func (c *WeatherStackClient) GetWeather(city string) (*WeatherStackResponse, error) {
	req := c.client.R().
		SetQueryParam("query", city).
		SetResult(&WeatherStackResponse{})

	httpResp, err := req.Get("/current")
	if err != nil {
		return nil, err
	}

	return httpResp.Result().(*WeatherStackResponse), nil
}

func newClient(baseURL string) *resty.Client {
	return resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(200*time.Millisecond).
		SetBaseURL(baseURL).
		SetHeader("Content-Type", "application/json")
}
