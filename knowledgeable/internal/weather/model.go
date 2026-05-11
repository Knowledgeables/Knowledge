package weather

import "time"

// OpenMeteoResponse represents the relevant fields from the Open-Meteo API.
type OpenMeteoResponse struct {
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	CurrentWeather CurrentWeather `json:"current_weather"`
	Current        CurrentData    `json:"current"`
}

// CurrentWeather holds basic real-time conditions.
type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windspeed"`
	WeatherCode int     `json:"weathercode"`
	Time        string  `json:"time"`
}

// CurrentData holds extra current fields like humidity.
type CurrentData struct {
	Humidity int `json:"relative_humidity_2m"`
}

// WeatherData is the domain model returned to the client.
type WeatherData struct {
	Temperature float64   `json:"temperature_c"`
	WindSpeed   float64   `json:"wind_speed_ms"`
	Humidity    int       `json:"humidity_pct"`
	Description string    `json:"description"`
	RetrievedAt time.Time `json:"retrieved_at"`
}
