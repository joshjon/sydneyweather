package weather

type WeatherStackCurrent struct {
	WindSpeed   int `json:"wind_speed"`
	Temperature int `json:"temperature"`
}

type WeatherStackResponse struct {
	Current WeatherStackCurrent `json:"current"`
}

type WeatherStackError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

type WeatherStackErrorResponse struct {
	Success bool              `json:"success"`
	Error   WeatherStackError `json:"error"`
}

type OpenWeatherMain struct {
	Temp int `json:"temp"`
}

type OpenWeatherWind struct {
	Speed int `json:"speed"`
}

type OpenWeatherResponse struct {
	Main OpenWeatherMain
	Wind OpenWeatherWind
}

type OpenWeatherErrorResponse struct {
	Code    int    `json:"cod"` // 'code' misspelled in open weather response
	Message string `json:"message"`
}
