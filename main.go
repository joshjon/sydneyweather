package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joshjon/sydneyweather/internal/api/v1"
)

var cfg = api.Config{
	City:               "Sydney",
	WeatherStackAPIKey: "",
	OpenWeatherAPIKey:  "",
	CacheExpiry:        3 * time.Second,
}

func main() {
	var port = flag.Int("p", 8080, "Port for HTTP server")
	flag.Parse()
	e := echo.New()
	e.Use(middleware.Logger())

	service := api.NewService(cfg)
	registerService(e, service)

	addr := &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: *port,
	}

	if err := e.Start(addr.String()); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

type serviceServer interface {
	GetWeather(ctx echo.Context) error
}

func registerService(e *echo.Echo, s serviceServer) {
	v1 := e.Group("/v1")
	v1.GET("/weather", s.GetWeather)
}
