package users

import (
	"encoding/json"
	"html/template"
	"knowledgeable/internal/middleware"
	"knowledgeable/internal/observability"
	"log/slog"
	"net/http"
	"os"
)

type Handler struct {
	service       UserService
	loadTmpl      func() *template.Template
	createSession func(int64) (string, error)
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

type RegisterPageData struct {
	Error string
}

type UserService interface {
	Register(string, string, string) (*User, error)
	GetAll() ([]User, error)
}

func NewHandler(s UserService, load func() *template.Template, createSession func(int64) (string, error)) *Handler {
	return &Handler{
		service:       s,
		loadTmpl:      load,
		createSession: createSession,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackingID := middleware.GetTrackingID(r)

	users, err := h.service.GetAll()
	if err != nil {
		slog.Error("get_users_failed",
			"event", "get_users_failed",
			"tracking_id", trackingID,
			"error", err.Error(),
		)
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
		tmpl := h.loadTmpl()

		data := RegisterPageData{}
		switch r.URL.Query().Get("error") {
		case "username_taken":
			data.Error = "Username is already taken"
		case "email_taken":
			data.Error = "Email is already registered"
		}

		if err := tmpl.ExecuteTemplate(w, "register.html", data); err != nil {
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

		trackingID := middleware.GetTrackingID(r)

		if err != nil {
			slog.Warn("register_failed",
				observability.LogAttrs("register_failed", trackingID, nil,
					"error", err.Error(),
				)...,
			)

			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		slog.Info("register",
			observability.LogAttrs("register", trackingID, &user.ID)...,
		) // #nosec G706 -- JSON handler escapes all values
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	ct := r.Header.Get("Content-Type")

	if ct == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		req.Username = r.FormValue("username")
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(req.Username, req.Email, req.Password)

	trackingID := middleware.GetTrackingID(r)

	if err != nil {
		slog.Warn("register_failed",
			observability.LogAttrs("register_failed", trackingID, nil,
				"error", err.Error(),
			)...,
		)

		switch err {
		case ErrUsernameTaken:
			http.Redirect(w, r, "/register?error=username_taken", http.StatusSeeOther)
		case ErrEmailTaken:
			http.Redirect(w, r, "/register?error=email_taken", http.StatusSeeOther)
		default:
			http.Error(w, "registration failed", http.StatusBadRequest)
		}
		return
	}

	sessionID, err := h.createSession(user.ID)
	if err != nil {
		slog.Error("session_failed",
			observability.LogAttrs("session_failed", trackingID, &user.ID,
				"error", err.Error(),
			)...,
		)
	}

	http.SetCookie(w, &http.Cookie{ // #nosec G124 -- Secure is set dynamically based on APP_ENV
		Name:     "session_id",
		Value:    sessionID,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Path:     "/",
		Secure:   os.Getenv("APP_ENV") == "production",
	})

	slog.Info("register",
		observability.LogAttrs("register", trackingID, &user.ID)...,
	) // #nosec G706 -- JSON handler escapes all values
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
