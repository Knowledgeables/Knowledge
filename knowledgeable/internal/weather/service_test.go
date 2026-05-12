package weather

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

// roundTripFunc lets tests intercept outgoing HTTP calls without a real server.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// fakeService returns a Service whose HTTP client returns the given response.
func fakeService(resp OpenMeteoResponse) *Service {
	body, _ := json.Marshal(resp)
	return &Service{
		client: &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	}
}

// HAPPY PATH

func TestDescribeWeatherCode(t *testing.T) {
	cases := []struct {
		code int
		want string
	}{
		{0, "Clear sky"},
		{1, "Partly cloudy"},
		{3, "Partly cloudy"},
		{10, "Fog"},
		{51, "Drizzle"},
		{56, "Freezing drizzle"},
		{61, "Rain"},
		{65, "Heavy rain"},
		{66, "Freezing rain"},
		{73, "Snow"},
		{77, "Snow grains"},
		{80, "Rain showers"},
		{85, "Snow showers"},
		{95, "Thunderstorm"},
		{96, "Thunderstorm with hail"},
		{999, "Unknown"},
	}

	for _, tc := range cases {
		got := describeWeatherCode(tc.code)
		if got != tc.want {
			t.Errorf("code %d: want %q, got %q", tc.code, tc.want, got)
		}
	}
}

func TestGetWeather_MapsFields(t *testing.T) {
	svc := fakeService(OpenMeteoResponse{
		CurrentWeather: CurrentWeather{
			Temperature: 18.3,
			WindSpeed:   3.2,
			WeatherCode: 63,
			Time:        "2024-06-15T10:00",
		},
		Current: CurrentData{Humidity: 75},
	})

	data, err := svc.GetWeather(context.Background(), 55.68, 12.57)
	if err != nil {
		t.Fatal(err)
	}

	if data.Temperature != 18.3 {
		t.Errorf("want temperature 18.3, got %v", data.Temperature)
	}
	if data.WindSpeed != 3.2 {
		t.Errorf("want wind speed 3.2, got %v", data.WindSpeed)
	}
	if data.Humidity != 75 {
		t.Errorf("want humidity 75, got %v", data.Humidity)
	}
	if data.Description != "Rain" {
		t.Errorf("want description 'Rain', got %q", data.Description)
	}
}

func TestGetWeather_RetrievedAt_UsesAPITime(t *testing.T) {
	svc := fakeService(OpenMeteoResponse{
		CurrentWeather: CurrentWeather{
			Temperature: 20.5,
			WindSpeed:   5.0,
			WeatherCode: 0,
			Time:        "2024-06-15T14:30",
		},
		Current: CurrentData{Humidity: 60},
	})

	data, err := svc.GetWeather(context.Background(), 55.68, 12.57)
	if err != nil {
		t.Fatal(err)
	}

	want := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	if !data.RetrievedAt.Equal(want) {
		t.Errorf("want %v, got %v", want, data.RetrievedAt)
	}
}

// ERROR PATH

func TestGetWeather_RetrievedAt_FallsBackOnBadTime(t *testing.T) {
	before := time.Now().UTC()
	svc := fakeService(OpenMeteoResponse{
		CurrentWeather: CurrentWeather{Time: "not-a-time"},
	})

	data, err := svc.GetWeather(context.Background(), 55.68, 12.57)
	if err != nil {
		t.Fatal(err)
	}
	after := time.Now().UTC()

	if data.RetrievedAt.Before(before) || data.RetrievedAt.After(after) {
		t.Errorf("fallback time %v not in expected range [%v, %v]", data.RetrievedAt, before, after)
	}
}
