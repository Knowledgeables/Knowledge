package pages

import (
	"encoding/json"
	"html/template"
	"net/http"
)

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

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	allPages, err := h.service.GetAllPages()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl := h.loadTmpl()

	if err := tmpl.ExecuteTemplate(w, "pages.html", allPages); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

}
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	lang := r.URL.Query().Get("language")

	if lang == "" {
		lang = "en"
	}

	results, err := h.service.Search(query, Language(lang))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := struct {
		Query   string
		Results []Page
	}{
		Query:   query,
		Results: results,
	}

	tmpl := h.loadTmpl()

	if err := tmpl.ExecuteTemplate(w, "search.html", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ViewPage(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")

	page, err := h.service.FindByURL(url)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl := h.loadTmpl()

	if err := tmpl.ExecuteTemplate(w, "page.html", page); err != nil {
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
// @Success 200 {array} object
// @Failure 500 {string} string "internal error"
// @Router /api/search [get]
func (h *Handler) SearchAPI(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	lang := r.URL.Query().Get("language")

	if lang == "" {
		lang = "en"
	}

	results, err := h.service.Search(query, Language(lang))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
