package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.open-meteo.com/v1/forecast"

// Service handles weather data retrieval.
type Service struct {
	client *http.Client
}

// NewService creates a Service with a sensible timeout.
func NewService() *Service {
	return &Service{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetWeather fetches the current weather for the given coordinates.
func (s *Service) GetWeather(ctx context.Context, lat, lon float64) (*WeatherData, error) {
	url := fmt.Sprintf(
		"%s?latitude=%.4f&longitude=%.4f&current_weather=true&current=relative_humidity_2m&wind_speed_unit=ms&timezone=auto",
		baseURL, lat, lon,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // #nosec G704 -- URL is baseURL const + range-validated float64 coords; no SSRF possible
	if err != nil {
		return nil, fmt.Errorf("weather: build request: %w", err)
	}

	resp, err := s.client.Do(req) // #nosec G704 -- same as above
	if err != nil {
		return nil, fmt.Errorf("weather: fetch: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather: unexpected status %d", resp.StatusCode)
	}

	var raw OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("weather: decode response: %w", err)
	}

	observedAt, err := time.Parse("2006-01-02T15:04", raw.CurrentWeather.Time)
	if err != nil {
		observedAt = time.Now()
	}

	return &WeatherData{
		Temperature: raw.CurrentWeather.Temperature,
		WindSpeed:   raw.CurrentWeather.WindSpeed,
		Humidity:    raw.Current.Humidity,
		Description: describeWeatherCode(raw.CurrentWeather.WeatherCode),
		RetrievedAt: observedAt.UTC(),
	}, nil
}

// describeWeatherCode maps WMO weather codes to readable strings.
func describeWeatherCode(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code <= 3:
		return "Partly cloudy"
	case code <= 49:
		return "Fog"
	case code <= 55:
		return "Drizzle"
	case code <= 57:
		return "Freezing drizzle"
	case code <= 63:
		return "Rain"
	case code <= 65:
		return "Heavy rain"
	case code <= 67:
		return "Freezing rain"
	case code <= 75:
		return "Snow"
	case code <= 77:
		return "Snow grains"
	case code <= 82:
		return "Rain showers"
	case code <= 86:
		return "Snow showers"
	case code == 95:
		return "Thunderstorm"
	case code <= 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown"
	}
}
