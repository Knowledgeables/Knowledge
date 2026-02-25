package auth

import (
	"encoding/json"
	"html/template"
	"knowledgeable/internal/users"
	"net/http"
)

type UserService interface {
	Login(string, string) (*users.User, error)
}

type Handler struct {
	userService UserService
	loginTmpl   *template.Template
}

func NewHandler(us UserService, tmpl *template.Template) *Handler {
	return &Handler{
		userService: us,
		loginTmpl:   tmpl,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		cookie, err := r.Cookie("session_id")
		if err == nil {
			if _, ok := Get(cookie.Value); ok {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}
		}

		if err := h.loginTmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "missing credentials", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(username, password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID, err := Create(user.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // slet cookie
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func (h *Handler) LoginAPI(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "missing credentials", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID, err := Create(user.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "login successful",
	}); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}
