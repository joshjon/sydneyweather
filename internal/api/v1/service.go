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

type Service struct {
	City      string
	primary   WeatherStackClient
	failOver  OpenWeatherClient
	respCache *valueCache[*GetWeatherResponse]
}

type Config struct {
	City               string
	WeatherStackAPIKey string
	OpenWeatherAPIKey  string
	CacheExpiry        time.Duration
}

func NewService(cfg Config) *Service {
	return &Service{
		City:      cfg.City,
		primary:   weather.NewWeatherStackClient(cfg.WeatherStackAPIKey),
		failOver:  weather.NewOpenWeatherClient(cfg.OpenWeatherAPIKey),
		respCache: newValueCache[*GetWeatherResponse](cfg.CacheExpiry),
	}
}

type GetWeatherResponse struct {
	WindSpeed   int `json:"wind_speed"`
	TempDegrees int `json:"temperature_degrees"`
}

func (s *Service) GetWeather(ctx echo.Context) error {
	if strings.ToLower(ctx.QueryParam("city")) != city {
		return echo.NewHTTPError(http.StatusBadRequest, "query param 'city' must have value 'sydney'")
	}

	if resp, ok := s.respCache.get(); ok {
		return ctx.JSON(http.StatusOK, resp)
	}

	var resp *GetWeatherResponse
	defer func() {
		if resp != nil {
			s.respCache.put(&resp)
		}
	}()

	primaryResp, err := s.primary.GetWeather(s.City)
	if err == nil {
		resp = &GetWeatherResponse{
			WindSpeed:   primaryResp.Current.WindSpeed,
			TempDegrees: primaryResp.Current.Temperature,
		}
		return ctx.JSON(http.StatusOK, resp)
	}
	log.Printf("error getting weather from primary source: %v", err)

	failOverResp, err := s.failOver.GetWeather(s.City)
	if err == nil {
		resp = &GetWeatherResponse{
			WindSpeed:   failOverResp.Wind.Speed,
			TempDegrees: failOverResp.Main.Temp,
		}
		return ctx.JSON(http.StatusOK, resp)
	}
	log.Printf("error getting weather from fail over source: %v", err)

	return echo.NewHTTPError(http.StatusServiceUnavailable)
}
