package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ServiceServer interface {
	GetWeather(ctx echo.Context) error
}

func RegisterService(e *echo.Echo, s ServiceServer) {
	v1 := e.Group("/v1")
	v1.GET("/weather", s.GetWeather)
}

type WeatherClient interface {
	GetWeather(city string)
}

type Service struct {
	WeatherStackClient WeatherClient
	OpenWeatherClient  WeatherClient
}

func NewService() *Service {
	return &Service{}
}

type GetWeatherResponse struct {
	WindSpeed   int `json:"wind_speed"`
	TempDegrees int `json:"temperature_degrees"`
}

func (s *Service) GetWeather(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, GetWeatherResponse{
		WindSpeed:   10,
		TempDegrees: 10,
	})
}
