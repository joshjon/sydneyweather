//go:build integration

package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshjon/sydneyweather/internal/api/v1"
)

const baseURL = "http://localhost:8080/v1"

func TestService_GetWeather(t *testing.T) {
	url := fmt.Sprintf("%s/weather?city=sydney", baseURL)
	req, err := http.NewRequest(echo.GET, url, nil)
	require.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	httpResp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, "200 OK", httpResp.Status)

	respBody, err := ioutil.ReadAll(httpResp.Body)
	require.NoError(t, err)

	var resp api.GetWeatherResponse
	require.NoError(t, json.Unmarshal(respBody, &resp))
	require.Equal(t, 10, resp.WindSpeed)
	require.Equal(t, 10, resp.TempDegrees)
}
