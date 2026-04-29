package pages

import (
	"encoding/json"
	"errors"
	"html/template"
	"knowledgeable/internal/auth"
	"knowledgeable/internal/middleware"
	"knowledgeable/internal/observability"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type SearchResponse struct {
	Results []Page `json:"results"`
	Count   int    `json:"count"`
}

type CrawlerTargetsResponse struct {
	Targets []CrawlSignal `json:"targets"`
	Count   int           `json:"count"`
}

type CrawlerIngestResponse struct {
	Upserted int `json:"upserted"`
	Received int `json:"received"`
}

type CrawlerExistingURLsResponse struct {
	URLs  []string `json:"urls"`
	Count int      `json:"count"`
}

type Handler struct {
	service  *Service
	loadTmpl func() *template.Template
}

func NewHandler(service *Service, load func() *template.Template) *Handler {
	return &Handler{
		service:  service,
		loadTmpl: load,
	}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {

	trackingID := middleware.GetTrackingID(r)

	var userID *int64
	if cookie, err := r.Cookie("session_id"); err == nil {
		if id, ok := auth.Get(cookie.Value); ok {
			userID = &id
		}
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	lang := r.URL.Query().Get("language")

	if lang == "" {
		lang = "en"
	}

	var results []Page
	var count int
	var err error

	if query == "" {

		slog.Info("search_empty",
			observability.LogAttrs("search_empty", trackingID, userID)...,
		)
		results = []Page{}
		count = 0
	} else {

		// the intent of the search request
		slog.Info("search",
			observability.LogAttrs("search", trackingID, userID,
				"query", query,
				"language", lang,
			)...,
		)

		results, count, err = h.service.Search(query, Language(lang))
		if err != nil {
			slog.Error("search_failed",
				observability.LogAttrs("search_failed", trackingID, userID,
					"query", query,
					"language", lang,
					"error", err.Error(),
				)...,
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// the outcome of the search query
		slog.Info("search_result",
			observability.LogAttrs("search_result", trackingID, userID,
				"query", query,
				"language", lang,
				"results", count,
			)...,
		)

		if err := h.service.RecordSignal(query, Language(lang)); err != nil {
			slog.Warn("record search signal failed",
				observability.LogAttrs("search_signal_failed", trackingID, userID,
					"query", query,
					"language", lang,
					"error", err.Error(),
				)...,
			)
		}
	}

	data := struct {
		Query   string
		Results []Page
		Count   int
	}{
		Query:   query,
		Results: results,
		Count:   count,
	}

	tmpl := h.loadTmpl()

	if err := tmpl.ExecuteTemplate(w, "search.html", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SearchAPI(w http.ResponseWriter, r *http.Request) {

	trackingID := middleware.GetTrackingID(r)

	var userID *int64
	if cookie, err := r.Cookie("session_id"); err == nil {
		if id, ok := auth.Get(cookie.Value); ok {
			userID = &id
		}
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	lang := r.URL.Query().Get("language")

	if lang == "" {
		lang = "en"
	}

	var results []Page
	var count int
	var err error

	if query == "" {
		slog.Info("search_empty",
			observability.LogAttrs("search_empty", trackingID, userID)...,
		)
		results = []Page{}
		count = 0
	} else {
		// the intend of the search requets
		slog.Info("search",
			observability.LogAttrs("search", trackingID, userID,
				"query", query,
				"language", lang,
			)...,
		)

		results, count, err = h.service.Search(query, Language(lang))
		if err != nil {
			slog.Error("search_failed",
				observability.LogAttrs("search_failed", trackingID, userID,
					"query", query,
					"language", lang,
					"error", err.Error(),
				)...,
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// the outcome of the search query
		slog.Info("search_result",
			observability.LogAttrs("search_result", trackingID, userID,
				"query", query,
				"language", lang,
				"results", count,
			)...,
		)

		if err := h.service.RecordSignal(query, Language(lang)); err != nil {
			slog.Warn("record search_api signal failed",
				observability.LogAttrs("search_signal_failed", trackingID, userID,
					"query", query,
					"language", lang,
					"error", err.Error(),
				)...,
			)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(SearchResponse{
		Results: results,
		Count:   count,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CrawlerTargetsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expectedKey := os.Getenv("CRAWLER_KEY")
	if expectedKey == "" {
		expectedKey = "dev-key"
	}
	if r.Header.Get("X-Crawler-Key") != expectedKey {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit := 20
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	targets, err := h.service.GetTopSignals(limit)
	if err != nil {
		slog.Error("crawler targets query failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(CrawlerTargetsResponse{
		Targets: targets,
		Count:   len(targets),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CrawlerExistingURLsAPI godoc
// @Summary Get existing crawler URLs
// @Description Return list of all URLs already in database
// @Tags crawler
// @Produce json
// @Param X-Crawler-Key header string true "Crawler authentication key"
// @Success 200 {object} pages.CrawlerExistingURLsResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 405 {string} string "method not allowed"
// @Failure 500 {string} string "internal error"
// @Router /api/crawler/existing-urls [get]
func (h *Handler) CrawlerExistingURLsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expectedKey := os.Getenv("CRAWLER_KEY")
	if expectedKey == "" {
		expectedKey = "dev-key"
	}
	if r.Header.Get("X-Crawler-Key") != expectedKey {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := h.service.GetExistingURLs()
	if err != nil {
		slog.Error("crawler existing urls query failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(CrawlerExistingURLsResponse{
		URLs:  urls,
		Count: len(urls),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CrawlerIngestAPI godoc
// @Summary Ingest crawled pages
// @Description Upsert crawler results into pages table
// @Tags crawler
// @Accept json
// @Produce json
// @Param X-Crawler-Key header string true "Crawler authentication key"
// @Param body body []pages.CrawlerIngestItem true "Crawler pages payload"
// @Success 200 {object} pages.CrawlerIngestResponse
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 405 {string} string "method not allowed"
// @Failure 500 {string} string "internal error"
// @Router /api/crawler/ingest [post]
func (h *Handler) CrawlerIngestAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expectedKey := os.Getenv("CRAWLER_KEY")
	if expectedKey == "" {
		expectedKey = "dev-key"
	}
	if r.Header.Get("X-Crawler-Key") != expectedKey {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	const maxBodySize = 10 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var items []CrawlerIngestItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	upserted, err := h.service.IngestCrawlerPages(items)
	if err != nil {
		slog.Error("crawler ingest failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(CrawlerIngestResponse{
		Upserted: upserted,
		Received: len(items),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	trackingID := middleware.GetTrackingID(r)

	query := r.URL.Query().Get("q")
	url := r.URL.Query().Get("url")
	posStr := r.URL.Query().Get("pos")

	pos, _ := strconv.Atoi(posStr)

	// ❗ VALIDATION (vigtigt)
	if url == "" || !strings.HasPrefix(url, "http") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	// log click
	slog.Info("search_click",
		observability.LogAttrs("search_click", trackingID, nil,
			"query", query,
			"url", url,
			"position", pos,
		)...,
	)

	// redirect
	http.Redirect(w, r, url, http.StatusFound)
}