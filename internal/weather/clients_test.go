package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	wantCity   = "Sydney"
	wantAPIKey = "some-key"
)

func TestNewWeatherStackClient(t *testing.T) {
	client := NewWeatherStackClient(wantAPIKey)
	require.Equal(t, wantAPIKey, client.apiKey)
	require.Equal(t, weatherStackBaseURL, client.http.BaseURL)
}

func TestWeatherStackClient_GetWeather(t *testing.T) {
	wantResp := WeatherStackResponse{
		Current: WeatherStackCurrent{
			WindSpeed:   10,
			Temperature: 20,
		},
	}
	wantURLValues := weatherStackURLValues(wantAPIKey)
	srv := mockServer(t, "/current", wantURLValues, http.StatusOK, wantResp)
	defer srv.Close()

	client := WeatherStackClient{
		http:   newRestyClient(srv.URL),
		apiKey: wantAPIKey,
	}
	resp, err := client.GetWeather(wantCity)
	require.NoError(t, err)
	require.Equal(t, wantResp, *resp)
}

func TestWeatherStackClient_GetWeather_error(t *testing.T) {
	tests := []struct {
		name        string
		wantCode    int
		wantErrResp *WeatherStackErrorResponse
	}{
		{
			name:     "status code and error body",
			wantCode: 444,
			wantErrResp: &WeatherStackErrorResponse{
				Success: false,
				Error: WeatherStackError{
					Code: 444,
					Type: "some-type",
					Info: "some-info",
				},
			},
		},
		{
			name:     "status code without error body",
			wantCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantURLValues := weatherStackURLValues(wantAPIKey)
			srv := mockServer(t, "/current", wantURLValues, tt.wantCode, tt.wantErrResp)
			defer srv.Close()

			client := WeatherStackClient{
				http:   newRestyClient(srv.URL),
				apiKey: wantAPIKey,
			}
			resp, err := client.GetWeather(wantCity)

			if tt.wantErrResp != nil {
				require.EqualError(t, err, newHTTPError(tt.wantCode, *tt.wantErrResp).Error())
			} else {
				require.EqualError(t, err, newHTTPError(tt.wantCode, nil).Error())
			}

			require.Nil(t, resp)
		})
	}
}

func TestNewOpenWeatherClient(t *testing.T) {
	client := NewOpenWeatherClient(wantAPIKey)
	require.Equal(t, wantAPIKey, client.apiKey)
	require.Equal(t, openWeatherBaseURL, client.http.BaseURL)
}

func TestOpenWeatherClient_GetWeather(t *testing.T) {
	wantResp := OpenWeatherResponse{
		Main: OpenWeatherMain{
			Temp: 10,
		},
		Wind: OpenWeatherWind{
			Speed: 20,
		},
	}
	wantURLValues := openWeatherURLValues(wantAPIKey)
	srv := mockServer(t, "/data/2.5/weather", wantURLValues, http.StatusOK, wantResp)
	defer srv.Close()

	client := OpenWeatherClient{
		http:   newRestyClient(srv.URL),
		apiKey: wantAPIKey,
	}
	resp, err := client.GetWeather(wantCity)
	require.NoError(t, err)
	require.Equal(t, wantResp, *resp)
}

func TestOpenWeather_GetWeather_error(t *testing.T) {
	tests := []struct {
		name        string
		wantCode    int
		wantErrResp *OpenWeatherErrorResponse
	}{
		{
			name:     "status code and error body",
			wantCode: 444,
			wantErrResp: &OpenWeatherErrorResponse{
				Code:    444,
				Message: "some-message",
			},
		},
		{
			name:     "status code without error body",
			wantCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantURLValues := openWeatherURLValues(wantAPIKey)
			srv := mockServer(t, "/data/2.5/weather", wantURLValues, tt.wantCode, tt.wantErrResp)
			defer srv.Close()

			client := OpenWeatherClient{
				http:   newRestyClient(srv.URL),
				apiKey: wantAPIKey,
			}
			resp, err := client.GetWeather(wantCity)

			if tt.wantErrResp != nil {
				require.EqualError(t, err, newHTTPError(tt.wantCode, *tt.wantErrResp).Error())
			} else {
				require.EqualError(t, err, newHTTPError(tt.wantCode, nil).Error())
			}

			require.Nil(t, resp)
		})
	}
}

func mockServer(t *testing.T, urlPath string, wantURLValues url.Values, wantCode int, wantResp any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == urlPath {
			require.Equal(t, wantURLValues, r.URL.Query())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(wantCode)
			if !(wantResp == nil || (reflect.ValueOf(wantResp).Kind() == reflect.Ptr && reflect.ValueOf(wantResp).IsNil())) {
				resp, err := json.Marshal(wantResp)
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			}
		} else {
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
}

func weatherStackURLValues(wantAPIKey string) url.Values {
	values := url.Values{}
	values.Set("access_key", wantAPIKey)
	values.Set("units", "m")
	values.Set("query", wantCity)
	return values
}

func openWeatherURLValues(wantAPIKey string) url.Values {
	values := url.Values{}
	values.Set("appid", wantAPIKey)
	values.Set("units", "metric")
	values.Set("q", wantCity)
	return values
}
