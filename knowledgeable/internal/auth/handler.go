package auth

import (
	"encoding/json"
	"html/template"
	"knowledgeable/internal/middleware"
	"knowledgeable/internal/observability"
	"knowledgeable/internal/users"
	"log/slog"
	"net/http"
	"os"
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

type MeResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserService interface {
	Login(string, string) (*users.User, error)
	GetByID(id int64) (*users.User, error)
	ChangePassword(userID int64, newPassword string) error
}

type Handler struct {
	userService UserService
	loadTmpl    func() *template.Template
}

type LoginPageData struct {
	Error string
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
// @Success 303 {string} string "Redirect to home"
// @Router /login [get]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err == nil {
		if _, ok := Get(cookie.Value); ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	tmpl := h.loadTmpl()

	data := LoginPageData{}
	if r.URL.Query().Get("error") == "invalid_credentials" {
		data.Error = "Invalid username or password"
	}

	if err := tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
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

	trackingID := middleware.GetTrackingID(r)

	userID, _ := ResolveFromRequest(r)

	cookie, err := r.Cookie("session_id")
	if err == nil {
		Delete(cookie.Value)
	}

	slog.Info("logout",
		observability.LogAttrs("auth.logout", trackingID, userID)...,
	)

	http.SetCookie(w, &http.Cookie{  // #nosec G124
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   os.Getenv("APP_ENV") == "production",
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

	
	trackingID := middleware.GetTrackingID(r)

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

		slog.Warn("login failed",
			observability.LogAttrs("auth.login_failed", trackingID, nil,
				"error", err.Error(),
			)...,
		)

		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	sessionID, err := Create(user.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{ // #nosec G124
		Name:     "session_id",
		Value:    sessionID,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Path:     "/",
		Secure:   os.Getenv("APP_ENV") == "production",
	})


	slog.Info("login",
		observability.LogAttrs("auth.login", trackingID, &user.ID)...,
	)

	if user.ShouldChangePassword {
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ChangePassword godoc
// @Summary Change password
// @Description Force-change password for authenticated user
// @Tags auth
// @Produce html
// @Success 303 {string} string "Redirect to home"
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "unauthorized"
// @Router /change-password [get,post]
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	trackingID := middleware.GetTrackingID(r)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Debug("unauthorized change of password",
			observability.LogAttrs("auth.unauthorized", trackingID, nil)...,
		)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := Get(cookie.Value)
	if !ok {
		slog.Debug("unauthorized",
			observability.LogAttrs("auth.unauthorized", trackingID, nil)...,
		)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl := h.loadTmpl()
		if err := tmpl.ExecuteTemplate(w, "change_password.html", nil); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		password := strings.TrimSpace(r.FormValue("password"))
		confirm := strings.TrimSpace(r.FormValue("confirm_password"))

		if password == "" || confirm == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		if password != confirm {
			http.Error(w, "passwords do not match", http.StatusBadRequest)
			return
		}

		if err := h.userService.ChangePassword(userID, password); err != nil {
			slog.Error("password_change_failed",
				observability.LogAttrs("auth.password_change_failed", trackingID, &userID,
					"error", err.Error(),
				)...,
			)

			http.Error(w, "failed to change password", http.StatusInternalServerError)
			return
		}

		slog.Info("password_changed",
			observability.LogAttrs("auth.password_changed", trackingID, &userID)...,
		)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// Me godoc
// @Summary Get current user
// @Description Returns authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} MeResponse
// @Failure 401 {string} string "unauthorized"
// @Router /api/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	trackingID := middleware.GetTrackingID(r)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Debug("unauthorized",
			observability.LogAttrs("auth.unauthorized", trackingID, nil)...,
		)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := Get(cookie.Value)
	if !ok {
		slog.Debug("unauthorized",
			observability.LogAttrs("auth.unauthorized", trackingID, nil)...,
		)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		slog.Warn("me_failed",
			observability.LogAttrs("auth.me_failed", trackingID, &userID,
				"error", err.Error(),
			)...,
		)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	res := MeResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("me_encode_failed",
			observability.LogAttrs("auth.me_encode_failed", trackingID, &userID,
				"error", err.Error(),
			)...,
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
