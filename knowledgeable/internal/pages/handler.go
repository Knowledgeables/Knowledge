package pages

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

type SearchResponse struct {
	Results []Page `json:"results"`
	Count   int    `json:"count"`
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
