package weather

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler exposes the weather endpoint.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a Handler with the given service and logger.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes mounts the weather endpoint onto the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/weather", h.GetWeather)
}

// GetWeather handles GET /api/weather?lat=55.68&lon=12.57
func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode weather response", "error", err)
	}
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
