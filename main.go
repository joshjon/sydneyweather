package main

import (
	"flag"
	"log"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joshjon/sydneyweather/internal/api/v1"
)

func main() {
	var port = flag.Int("p", 8080, "Port for HTTP server")
	flag.Parse()
	e := echo.New()
	e.Use(middleware.Logger())
	api.RegisterService(e, api.NewService())

	addr := &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: *port,
	}

	if err := e.Start(addr.String()); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
