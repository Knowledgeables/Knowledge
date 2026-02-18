package auth

import (
	"net/http"
	"knowledgeable/internal/users"
	"html/template"
)

var loginTmpl = template.Must(template.ParseFiles("templates/login.html"))


type Handler struct {
	userService *users.Service
}

func NewHandler(us *users.Service) *Handler {
	return &Handler{userService: us}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {


	if r.Method == http.MethodGet {
	loginTmpl.Execute(w, nil)
	return
}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.userService.Login(username, password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := Create(user.ID)

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


