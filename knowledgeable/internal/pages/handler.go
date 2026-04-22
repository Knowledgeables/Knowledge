package pages

import (
	"encoding/json"
	"html/template"
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

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	lang := r.URL.Query().Get("language")

	if lang == "" {
		lang = "en"
	}

	var results []Page
	var count int
	var err error

	if query == "" {
		results = []Page{}
		count = 0
	} else {
		results, count, err = h.service.Search(query, Language(lang))
		if err != nil {
			slog.Error("search failed", "query", query, "language", lang, "error", err) // #nosec G706 -- JSON handler escapes all values
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slog.Info("search", "query", query, "language", lang, "results", count) // #nosec G706 -- JSON handler escapes all values
		if err := h.service.RecordSignal(query, Language(lang)); err != nil {
			slog.Warn("record search signal failed", "query", query, "language", lang, "error", err) // #nosec G706 -- JSON handler escapes all values
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

// SearchAPI godoc
// @Summary Search
// @Description Search pages by query and optional language
// @Tags pages
// @Produce json
// @Param q query string true "Search query"
// @Param language query string false "Language code"
// @Success 200 {object} pages.SearchResponse
// @Failure 500 {string} string "internal error"
// @Router /api/search [get]
func (h *Handler) SearchAPI(w http.ResponseWriter, r *http.Request) {

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	lang := r.URL.Query().Get("language")

	if lang == "" {
		lang = "en"
	}

	var results []Page
	var count int
	var err error

	if query == "" {
		results = []Page{}
		count = 0
	} else {
		results, count, err = h.service.Search(query, Language(lang))
		if err != nil {
			slog.Error("search_api failed", "query", query, "language", lang, "error", err) // #nosec G706 -- JSON handler escapes all values
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slog.Info("search_api", "query", query, "language", lang, "results", count) // #nosec G706 -- JSON handler escapes all values
		if err := h.service.RecordSignal(query, Language(lang)); err != nil {
			slog.Warn("record search_api signal failed", "query", query, "language", lang, "error", err) // #nosec G706 -- JSON handler escapes all values
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

// CrawlerTargetsAPI godoc
// @Summary Get top crawl targets
// @Description Returns top search signals for crawler seed targets
// @Tags crawler
// @Produce json
// @Param X-Crawler-Key header string true "Crawler authentication key"
// @Param limit query int false "Max targets to return (default 20, max 100)"
// @Success 200 {object} pages.CrawlerTargetsResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 405 {string} string "method not allowed"
// @Failure 500 {string} string "internal error"
// @Router /api/crawler/targets [get]
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

	var items []CrawlerIngestItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
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
