package weather

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const weatherStackAPIKey = "_"

func TestWeatherStackClient_GetWeather(t *testing.T) {
	client := NewWeatherStackClient(weatherStackAPIKey)
	_, err := client.GetWeather("Sydney")
	require.NoError(t, err)
}
