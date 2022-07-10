package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/joshjon/sydneyweather/internal/weather"
)

const (
	wantTemp  = 10
	wantSpeed = 20
	city      = "Sydney"
)

func TestService_GetWeather(t *testing.T) {
	tests := []struct {
		name            string
		primaryEnabled  bool
		failOverEnabled bool
	}{
		{
			name:            "get weather from primary source",
			primaryEnabled:  true,
			failOverEnabled: false,
		},
		{
			name:            "get weather from fail over source",
			primaryEnabled:  false,
			failOverEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/v1/weather?city="+city, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			primary := &mockWeatherStackClient{wantErr: true}
			failOver := &mockOpenWeatherClient{wantErr: true}

			if tt.primaryEnabled {
				primary.wantErr = false
			}
			if tt.failOverEnabled {
				failOver.wantErr = false
			}

			s := Service{
				City:      city,
				primary:   primary,
				failOver:  failOver,
				respCache: newValueCache[*GetWeatherResponse](100 * time.Millisecond),
			}

			err := s.GetWeather(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, rec.Code)

			var resp GetWeatherResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Equal(t, wantSpeed, resp.WindSpeed)
			require.Equal(t, wantTemp, resp.TempDegrees)
		})
	}
}

func TestService_GetWeather_useCache(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/weather?city=Sydney", nil)
	rec1, rec2 := httptest.NewRecorder(), httptest.NewRecorder()
	ctx1, ctx2 := e.NewContext(req, rec1), e.NewContext(req, rec2)

	s := Service{
		City:      "Sydney",
		primary:   &mockWeatherStackClient{},
		failOver:  &mockOpenWeatherClient{},
		respCache: newValueCache[*GetWeatherResponse](100 * time.Millisecond),
	}

	err := s.GetWeather(ctx1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec1.Code)

	var resp GetWeatherResponse
	require.NoError(t, json.Unmarshal(rec1.Body.Bytes(), &resp))
	require.Equal(t, wantSpeed, resp.WindSpeed)
	require.Equal(t, wantTemp, resp.TempDegrees)

	err = s.GetWeather(ctx2)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec2.Code)

	// Return cached response without experiencing client errors
	s.primary = &mockWeatherStackClient{wantErr: true}
	s.failOver = &mockOpenWeatherClient{wantErr: true}
	var resp2 GetWeatherResponse
	require.NoError(t, json.Unmarshal(rec2.Body.Bytes(), &resp2))
	require.Equal(t, wantSpeed, resp2.WindSpeed)
	require.Equal(t, wantTemp, resp2.TempDegrees)
}

func TestService_GetWeather_unavailableError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/weather?city=Sydney", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	s := Service{
		City:      "Sydney",
		primary:   &mockWeatherStackClient{wantErr: true},
		failOver:  &mockOpenWeatherClient{wantErr: true},
		respCache: newValueCache[*GetWeatherResponse](100 * time.Millisecond),
	}

	err := s.GetWeather(ctx)
	require.EqualError(t, err, "code=503, message=Service Unavailable")
}

type mockWeatherStackClient struct {
	wantErr bool
}

func (c *mockWeatherStackClient) GetWeather(_ string) (*weather.WeatherStackResponse, error) {
	if c.wantErr {
		return nil, errors.New("some-error")
	}
	return &weather.WeatherStackResponse{
		Current: weather.WeatherStackCurrent{
			WindSpeed:   wantSpeed,
			Temperature: wantTemp,
		},
	}, nil
}

type mockOpenWeatherClient struct {
	wantErr bool
}

func (c *mockOpenWeatherClient) GetWeather(_ string) (*weather.OpenWeatherResponse, error) {
	if c.wantErr {
		return nil, errors.New("some-error")
	}
	return &weather.OpenWeatherResponse{
		Wind: weather.OpenWeatherWind{
			Speed: wantSpeed,
		},
		Main: weather.OpenWeatherMain{
			Temp: wantTemp,
		},
	}, nil
}
