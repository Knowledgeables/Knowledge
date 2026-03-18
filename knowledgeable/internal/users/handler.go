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

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Username string `json:"username"`
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


// RegisterPage godoc
// @Summary Server Register page
// @Description Shows the register form
// @Tags pages
// @Produce html
// @Success 200 {string} string "Register page"
// @Router /register [get]
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

// RegisterAPI godoc
// @Summary Register user
// @Description Create a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Register user"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal error"
// @Router /api/register [post]
func (h *Handler) RegisterAPI(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		log.Println("RegisterAPI error:", err)
		http.Error(w, "registration failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := RegisterResponse{
		Status:   "ok",
		Message:  "user registered",
		Username: user.Username,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}
