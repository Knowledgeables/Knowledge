package pages

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	allPages, err := h.service.GetAllPages()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/pages.html"))
	if err := tmpl.Execute(w, allPages); err != nil {
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

	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	if err := tmpl.Execute(w, data); err != nil {
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

	tmpl := template.Must(template.ParseFiles("templates/page.html"))

	if err := tmpl.Execute(w, page); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

}

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
