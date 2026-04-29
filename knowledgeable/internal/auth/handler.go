package auth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"knowledgeable/internal/middleware"
	"knowledgeable/internal/observability"
	"knowledgeable/internal/ratelimit"
	"knowledgeable/internal/users"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
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
	GetByEmail(email string) (*users.User, error)
	CreateResetToken(userID int64) (string, error)
	ConsumeResetToken(rawToken string) (int64, error)
}

type EmailSender interface {
	Send(to, subject, html string) error
}

type Handler struct {
	userService    UserService
	loadTmpl       func() *template.Template
	emailSender    EmailSender
	baseURL        string
	forgotLimiter  *ratelimit.IPLimiter
	resetLimiter   *ratelimit.IPLimiter
}

type LoginPageData struct {
	Error string
}

type ForgotPasswordPageData struct {
	Sent  bool
	Error string
}

type ResetPasswordPageData struct {
	Token string
	Error string
}

func NewHandler(us UserService, load func() *template.Template, emailSender EmailSender, baseURL string) *Handler {
	return &Handler{
		userService:   us,
		loadTmpl:      load,
		emailSender:   emailSender,
		baseURL:       baseURL,
		forgotLimiter: ratelimit.NewIPLimiter(5, 15*time.Minute),
		resetLimiter:  ratelimit.NewIPLimiter(10, 15*time.Minute),
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

// ForgotPassword godoc
// @Summary Request password reset
// @Description Sends a password reset email if the address is registered
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param email formData string true "Email address"
// @Success 303 {string} string "Redirect to /forgot-password?sent=1"
// @Router /forgot-password [get,post]
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	trackingID := middleware.GetTrackingID(r)

	if r.Method == http.MethodGet {
		tmpl := h.loadTmpl()
		data := ForgotPasswordPageData{}
		switch r.URL.Query().Get("error") {
		case "missing_email":
			data.Error = "Please enter an email address"
		case "invalid_token":
			data.Error = "This reset link is invalid or has expired"
		}
		if r.URL.Query().Get("sent") == "1" {
			data.Sent = true
		}
		if err := tmpl.ExecuteTemplate(w, "forgot_password.html", data); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := ratelimit.RealIP(r)

	if !h.forgotLimiter.Allow(ip) {
		slog.Warn("forgot_password_rate_limited",
			observability.LogAttrs("auth.forgot_password", trackingID, nil, "ip", ip)...,
		)
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	emailAddr := strings.TrimSpace(r.FormValue("email"))
	if emailAddr == "" {
		http.Redirect(w, r, "/forgot-password?error=missing_email", http.StatusSeeOther)
		return
	}

	emailDomain := emailDomainOf(emailAddr)

	slog.Info("forgot_password_requested",
		observability.LogAttrs("auth.forgot_password", trackingID, nil, "email_domain", emailDomain, "ip", ip)...,
	)

	user, err := h.userService.GetByEmail(emailAddr)
	if err != nil || user == nil {
		slog.Info("forgot_password_unknown_email",
			observability.LogAttrs("auth.forgot_password", trackingID, nil, "email_domain", emailDomain)...,
		)
		http.Redirect(w, r, "/forgot-password?sent=1", http.StatusSeeOther)
		return
	}

	rawToken, err := h.userService.CreateResetToken(user.ID)
	if err != nil {
		slog.Error("forgot_password_token_failed",
			observability.LogAttrs("auth.forgot_password", trackingID, &user.ID, "error", err.Error())...,
		)
		http.Redirect(w, r, "/forgot-password?sent=1", http.StatusSeeOther)
		return
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", h.baseURL, rawToken)
	html := fmt.Sprintf(
		`<p>Click <a href="%s">here</a> to reset your password. This link expires in 1 hour.</p>`,
		resetURL,
	)

	if err := h.emailSender.Send(user.Email, "Reset your password", html); err != nil {
		slog.Error("forgot_password_email_failed",
			observability.LogAttrs("auth.forgot_password", trackingID, &user.ID, "error", err.Error())...,
		)
		http.Redirect(w, r, "/forgot-password?sent=1", http.StatusSeeOther)
		return
	}

	slog.Info("forgot_password_sent",
		observability.LogAttrs("auth.forgot_password", trackingID, &user.ID, "email_domain", emailDomain)...,
	)

	http.Redirect(w, r, "/forgot-password?sent=1", http.StatusSeeOther)
}

// ResetPassword godoc
// @Summary Reset password via token
// @Description Validates the reset token and updates the user's password
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param token query string true "Reset token"
// @Success 303 {string} string "Redirect to /login"
// @Router /reset-password [get,post]
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	trackingID := middleware.GetTrackingID(r)

	if r.Method == http.MethodGet {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
			return
		}
		tmpl := h.loadTmpl()
		data := ResetPasswordPageData{Token: token}
		if r.URL.Query().Get("error") == "passwords_mismatch" {
			data.Error = "Passwords do not match"
		}
		if err := tmpl.ExecuteTemplate(w, "reset_password.html", data); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := ratelimit.RealIP(r)

	if !h.resetLimiter.Allow(ip) {
		slog.Warn("reset_password_rate_limited",
			observability.LogAttrs("auth.reset_password", trackingID, nil, "ip", ip)...,
		)
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	token := strings.TrimSpace(r.FormValue("token"))
	password := strings.TrimSpace(r.FormValue("password"))
	confirm := strings.TrimSpace(r.FormValue("confirm_password"))

	if token == "" || password == "" || confirm == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	if password != confirm {
		slog.Debug("reset_password_mismatch",
			observability.LogAttrs("auth.reset_password", trackingID, nil)...,
		)
		http.Redirect(w, r, "/reset-password?token="+token+"&error=passwords_mismatch", http.StatusSeeOther)
		return
	}

	slog.Info("reset_password_attempted",
		observability.LogAttrs("auth.reset_password", trackingID, nil)...,
	)

	userID, err := h.userService.ConsumeResetToken(token)
	if err != nil {
		slog.Warn("reset_password_invalid_token",
			observability.LogAttrs("auth.reset_password", trackingID, nil, "error", err.Error())...,
		)
		http.Redirect(w, r, "/forgot-password?error=invalid_token", http.StatusSeeOther)
		return
	}

	if err := h.userService.ChangePassword(userID, password); err != nil {
		slog.Error("reset_password_change_failed",
			observability.LogAttrs("auth.reset_password", trackingID, &userID, "error", err.Error())...,
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	slog.Info("reset_password_success",
		observability.LogAttrs("auth.reset_password", trackingID, &userID)...,
	)

	http.Redirect(w, r, "/login?reset=1", http.StatusSeeOther)
}

func emailDomainOf(email string) string {
	if i := strings.LastIndex(email, "@"); i >= 0 {
		return email[i+1:]
	}
	return "unknown"
}
