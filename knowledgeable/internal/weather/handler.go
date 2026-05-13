package weather

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler exposes the weather endpoints.
type Handler struct {
	service    *Service
	logger     *slog.Logger
	tmplLoader func() *template.Template
}

// NewHandler creates a Handler with the given service, logger, and template loader.
func NewHandler(service *Service, logger *slog.Logger, tmplLoader func() *template.Template) *Handler {
	return &Handler{
		service:    service,
		logger:     logger,
		tmplLoader: tmplLoader,
	}
}

// GetWeatherPage handles GET /weather — renders the weather page for users.
// @Summary Serve Weather page
// @Description Render the weather page
// @Tags weather
// @Produce html
// @Success 200 {string} string "Weather page"
// @Router /weather [get]
func (h *Handler) GetWeatherPage(w http.ResponseWriter, r *http.Request) {
	tmpl := h.tmplLoader()
	if err := tmpl.ExecuteTemplate(w, "weather.html", nil); err != nil {
		h.logger.Error("GetWeatherPage: template execution failed", "error", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// GetWeather handles GET /api/weather?lat=55.68&lon=12.57
// @Summary Get weather data
// @Description Returns current weather for the given coordinates
// @Tags weather
// @Produce json
// @Param lat query number true "Latitude (-90 to 90)"
// @Param lon query number true "Longitude (-180 to 180)"
// @Success 200 {object} weather.WeatherData
// @Failure 400 {string} string "missing or invalid coordinates"
// @Failure 502 {string} string "failed to retrieve weather data"
// @Router /api/weather [get]
func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	lat, lon, err := parseCoordinates(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.service.GetWeather(r.Context(), lat, lon)
	if err != nil {
		h.logger.Error("failed to fetch weather", "error", err, "lat", lat, "lon", lon)
		http.Error(w, "failed to retrieve weather data", http.StatusBadGateway)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		h.logger.Error("failed to encode weather response", "error", err)
		http.Error(w, "failed to encode weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = buf.WriteTo(w)
}

func parseCoordinates(r *http.Request) (float64, float64, error) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		return 0, 0, errMissingCoordinates
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || lat < -90 || lat > 90 {
		return 0, 0, errInvalidLatitude
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil || lon < -180 || lon > 180 {
		return 0, 0, errInvalidLongitude
	}

	return lat, lon, nil
}

type weatherError string

func (e weatherError) Error() string { return string(e) }

const (
	errMissingCoordinates weatherError = "both 'lat' and 'lon' query parameters are required"
	errInvalidLatitude    weatherError = "'lat' must be a number between -90 and 90"
	errInvalidLongitude   weatherError = "'lon' must be a number between -180 and 180"
)
