package auth

import (
	"html/template"
	"knowledgeable/internal/users"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserService interface {
	Login(string, string) (*users.User, error)
}

type Handler struct {
	userService UserService
	loadTmpl    func() *template.Template
}

func NewHandler(us UserService, load func() *template.Template) *Handler {
	return &Handler{
		userService: us,
		loadTmpl:    load,
	}
}

// LoginPage godoc
// @Summary Serve Login page
// @Description Render login page
// @Tags pages
// @Produce html
// @Success 200 {string} string "Login page"
// @Success 303 {string} string "Redirect to dashboard"
// @Router /login [get]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err == nil {
		if _, ok := Get(cookie.Value); ok {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	tmpl := h.loadTmpl()

	if err := tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// Logout godoc
// @Summary Logout user
// @Description Deletes the current session and redirects to login page
// @Tags auth
// @Produce html
// @Success 303 {string} string "Redirect to login page"
// @Header 303 {string} Location "/login"
// @Failure 405 {string} string "method not allowed"
// @Router /logout [post]
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

// LoginAPI godoc
// @Summary Login
// @Description Authenticate user and create session
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string "missing credentials or bad form data"
// @Failure 401 {string} string "invalid credentials"
// @Failure 500 {string} string "internal error"
// @Router /api/login [post]
func (h *Handler) LoginAPI(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	req := LoginRequest{
		Username: strings.TrimSpace(r.FormValue("username")),
		Password: strings.TrimSpace(r.FormValue("password")),
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

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
