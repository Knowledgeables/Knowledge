package pages

import (
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

	tmpl.Execute(w, allPages)
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

	tmpl := template.Must(template.ParseFiles("templates/search.html"))

	data := struct {
		Query   string
		Results []Page
	}{
		Query:   query,
		Results: results,
	}

	tmpl.Execute(w, data)
}
func (h *Handler) ViewPage(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")

	page, err := h.service.FindByURL(url)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/page.html"))
	tmpl.Execute(w, page)
}
