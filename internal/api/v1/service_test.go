package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRegisterService(t *testing.T) {
	e := echo.New()
	RegisterService(e, NewService())
	routes := e.Routes()
	require.Len(t, routes, 1)
	require.Equal(t, routes[0].Method, http.MethodGet)
	require.Equal(t, routes[0].Path, "/v1/weather")
}

func TestService_GetWeather(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("city")
	ctx.SetParamValues("sydney")

	s := NewService()
	err := s.GetWeather(ctx)
	require.NoError(t, err)

	var resp GetWeatherResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, 10, resp.WindSpeed)
	require.Equal(t, 10, resp.TempDegrees)
}
