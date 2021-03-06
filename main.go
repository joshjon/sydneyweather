package main

import (
	"log"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joshjon/sydneyweather/internal/api/v1"
	"github.com/joshjon/sydneyweather/internal/config"
)

// main starts a new echo server registered with the sydney weather service.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v\n", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	serviceCfg := api.Config{
		WeatherStackAPIKey: cfg.WeatherStackAPIKey,
		OpenWeatherAPIKey:  cfg.OpenWeatherAPIKey,
		CacheExpiry:        cfg.CacheExpiry,
	}

	service := api.NewService(serviceCfg)
	registerService(e, service)

	addr := &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: cfg.ServerPort,
	}

	if err = e.Start(addr.String()); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}

func registerService(e *echo.Echo, s *api.Service) {
	v1 := e.Group("/v1")
	v1.GET("/weather", s.GetWeather)
}
