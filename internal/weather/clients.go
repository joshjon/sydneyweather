package weather

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	weatherStackBaseURL = "http://api.weatherstack.com"
	openWeatherBaseURL  = "https://api.openweathermap.org"
)

// WeatherStackClient is a simple client for retrieving basic weather data from
// the weatherstack API. A valid API key must be provided in order to successfully
// authenticate on each request.
type WeatherStackClient struct {
	http   *resty.Client
	apiKey string
}

func NewWeatherStackClient(apiKey string) *WeatherStackClient {
	client := newRestyClient(weatherStackBaseURL)
	return &WeatherStackClient{
		http:   client,
		apiKey: apiKey,
	}
}

// GetWeather returns the temperature and wind speed for the specified city.
func (c *WeatherStackClient) GetWeather(city string) (*WeatherStackResponse, error) {
	req := c.http.R().
		SetQueryParam("access_key", c.apiKey).
		SetQueryParam("units", "m"). // Celsius
		SetQueryParam("query", city).
		SetResult(WeatherStackResponse{})
	return get[WeatherStackResponse, WeatherStackErrorResponse](req, "/current")
}

// OpenWeatherClient is a simple client for retrieving basic weather data from
// the OpenWeather API. A valid API key must be provided in order to successfully
// authenticate on each request.
type OpenWeatherClient struct {
	http   *resty.Client
	apiKey string
}

func NewOpenWeatherClient(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		http:   newRestyClient(openWeatherBaseURL),
		apiKey: apiKey,
	}
}

// GetWeather returns the temperature and wind speed for the specified city.
func (c *OpenWeatherClient) GetWeather(city string) (*OpenWeatherResponse, error) {
	req := c.http.R().
		SetQueryParam("appid", c.apiKey).
		SetQueryParam("units", "metric"). // Celsius
		SetQueryParam("q", city).
		SetResult(&OpenWeatherResponse{})
	return get[OpenWeatherResponse, OpenWeatherErrorResponse](req, "/data/2.5/weather")
}

// get performs a GET request to the specified URL and returns the response.
// The R and E type parameters are used when unmarshalling any response or error
// body.
func get[R any, E any](req *resty.Request, url string) (*R, error) {
	httpResp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	if httpResp.IsSuccess() {
		return httpResp.Result().(*R), nil
	}

	var errResp E
	if err = json.Unmarshal(httpResp.Body(), &errResp); err != nil {
		return nil, newHTTPError(httpResp.StatusCode(), nil)
	}

	return nil, newHTTPError(httpResp.StatusCode(), errResp)

}

func newRestyClient(baseURL string) *resty.Client {
	return resty.New().
		SetBaseURL(baseURL).
		SetHeader("Content-Type", "application/json")
}

func newHTTPError(code int, err any) error {
	if err != nil {
		return fmt.Errorf("http error; status code: %d; error: %+v",
			code,
			err,
		)
	}
	return fmt.Errorf("http error; status code: %d", code)
}
