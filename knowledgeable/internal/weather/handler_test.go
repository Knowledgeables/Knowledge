package weather

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeHandler(resp OpenMeteoResponse) *Handler {
	return NewHandler(fakeService(resp), slog.Default(), nil)
}

// parseCoordinates — happy path

func TestParseCoordinates_Valid(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?lat=55.68&lon=12.57", nil)
	lat, lon, err := parseCoordinates(r)
	if err != nil {
		t.Fatal(err)
	}
	if lat != 55.68 || lon != 12.57 {
		t.Fatalf("got lat=%v lon=%v", lat, lon)
	}
}

func TestParseCoordinates_Boundaries(t *testing.T) {
	cases := []struct{ lat, lon string }{
		{"-90", "-180"},
		{"90", "180"},
		{"0", "0"},
	}
	for _, tc := range cases {
		r := httptest.NewRequest(http.MethodGet, "/?lat="+tc.lat+"&lon="+tc.lon, nil)
		if _, _, err := parseCoordinates(r); err != nil {
			t.Errorf("lat=%s lon=%s: unexpected error: %v", tc.lat, tc.lon, err)
		}
	}
}

// parseCoordinates — error path

func TestParseCoordinates_Missing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	_, _, err := parseCoordinates(r)
	if err == nil {
		t.Fatal("expected error for missing params")
	}
}

func TestParseCoordinates_InvalidLat(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?lat=999&lon=12.57", nil)
	_, _, err := parseCoordinates(r)
	if err == nil {
		t.Fatal("expected error for out-of-range lat")
	}
}

func TestParseCoordinates_InvalidLon(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?lat=55.68&lon=999", nil)
	_, _, err := parseCoordinates(r)
	if err == nil {
		t.Fatal("expected error for out-of-range lon")
	}
}

func TestParseCoordinates_NonNumeric(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?lat=abc&lon=12.57", nil)
	_, _, err := parseCoordinates(r)
	if err == nil {
		t.Fatal("expected error for non-numeric lat")
	}
}

// GetWeather handler — happy path

func TestGetWeather_Success(t *testing.T) {
	h := makeHandler(OpenMeteoResponse{
		CurrentWeather: CurrentWeather{
			Temperature: 20.0,
			WindSpeed:   4.0,
			WeatherCode: 0,
			Time:        "2024-06-15T12:00",
		},
		Current: CurrentData{Humidity: 55},
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/weather?lat=55.68&lon=12.57", nil)
	h.GetWeather(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	var data WeatherData
	if err := json.NewDecoder(w.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.Temperature != 20.0 {
		t.Errorf("want temperature 20.0, got %v", data.Temperature)
	}
	if data.Description != "Clear sky" {
		t.Errorf("want 'Clear sky', got %q", data.Description)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("want Content-Type application/json, got %q", ct)
	}
}

// GetWeather handler — error path

func TestGetWeather_MethodNotAllowed(t *testing.T) {
	h := makeHandler(OpenMeteoResponse{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/weather?lat=55.68&lon=12.57", nil)

	h.GetWeather(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d", w.Code)
	}
}

func TestGetWeather_MissingParams(t *testing.T) {
	h := makeHandler(OpenMeteoResponse{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/weather", nil)

	h.GetWeather(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestGetWeather_InvalidCoords(t *testing.T) {
	h := makeHandler(OpenMeteoResponse{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/weather?lat=999&lon=12.57", nil)

	h.GetWeather(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}
