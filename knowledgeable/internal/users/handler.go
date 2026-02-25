package users

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	service UserService
	tmpl    *template.Template
}

type UserService interface {
	Register(string, string, string) (*User, error)
	GetAll() ([]User, error)
}

func NewHandler(s UserService, tmpl *template.Template) *Handler {
	return &Handler{
		service: s,
		tmpl:    tmpl,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.service.GetAll()
	if err != nil {
		log.Println("GetAll error:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := h.tmpl.ExecuteTemplate(w, "register.html", nil); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.service.Register(username, email, password)

		if err != nil {
			log.Println("Register error:", err)
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		log.Println("User created: ", user.Username)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
