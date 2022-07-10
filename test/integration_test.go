//go:build integration

package test

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/joshjon/sydneyweather/internal/api/v1"
)

const baseURL = "http://localhost:8080"

func TestService_GetWeather(t *testing.T) {
	var resp api.GetWeatherResponse
	httpResp, err := resty.New().R().
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetQueryParam("city", "Sydney").
		SetResult(&resp).
		Get(baseURL + "/v1/weather")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, httpResp.StatusCode())
	require.NotEmpty(t, resp.WindSpeed)
	require.NotEmpty(t, resp.TempDegrees)
	t.Logf("response: %+v", resp)
}
