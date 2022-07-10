package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/joshjon/sydneyweather/internal/weather"
)

const city = "sydney"

type WeatherStackClient interface {
	GetWeather(city string) (*weather.WeatherStackResponse, error)
}

type OpenWeatherClient interface {
	GetWeather(city string) (*weather.OpenWeatherResponse, error)
}

// Service is an HTTP api that provides basic weather data for a given city.
// The service is configured to only accept 'sydney' as a city, although, it can
// be easily tweaked to support any city if required.
type Service struct {
	primary   WeatherStackClient
	failOver  OpenWeatherClient
	respCache *valueCache[*GetWeatherResponse]
}

type Config struct {
	WeatherStackAPIKey string
	OpenWeatherAPIKey  string
	CacheExpiry        time.Duration
}

func NewService(cfg Config) *Service {
	return &Service{
		primary:   weather.NewWeatherStackClient(cfg.WeatherStackAPIKey),
		failOver:  weather.NewOpenWeatherClient(cfg.OpenWeatherAPIKey),
		respCache: newValueCache[*GetWeatherResponse](cfg.CacheExpiry),
	}
}

type GetWeatherResponse struct {
	WindSpeed   int `json:"wind_speed"`
	TempDegrees int `json:"temperature_degrees"`
}

// GetWeather returns the temperature and wind speed for the specified city.
// Data retrieval is prioritized in the following order: cache (non expired),
// primary source, fail over source, cache (stale).
func (s *Service) GetWeather(ctx echo.Context) error {
	if strings.ToLower(ctx.QueryParam("city")) != city {
		return echo.NewHTTPError(http.StatusBadRequest, "query param 'city' must have value 'sydney'")
	}

	if !s.respCache.expired() {
		if resp, ok := s.respCache.get(); ok {
			return ctx.JSON(http.StatusOK, resp)
		}
	}

	var resp *GetWeatherResponse
	defer func() {
		if resp != nil {
			s.respCache.put(&resp)
		}
	}()

	primaryResp, err := s.primary.GetWeather(city)
	if err == nil {
		resp = &GetWeatherResponse{
			WindSpeed:   primaryResp.Current.WindSpeed,
			TempDegrees: primaryResp.Current.Temperature,
		}
		return ctx.JSON(http.StatusOK, resp)
	}
	log.Printf("error getting weather from primary source: %v\n", err)

	failOverResp, err := s.failOver.GetWeather(city)
	if err == nil {
		resp = &GetWeatherResponse{
			WindSpeed:   int(failOverResp.Wind.Speed),
			TempDegrees: int(failOverResp.Main.Temp),
		}
		return ctx.JSON(http.StatusOK, resp)
	}
	log.Printf("error getting weather from fail over source: %v\n", err)

	// Serve stale weather data
	if resp, ok := s.respCache.get(); ok {
		return ctx.JSON(http.StatusOK, resp)
	}

	return echo.NewHTTPError(http.StatusServiceUnavailable)
}
