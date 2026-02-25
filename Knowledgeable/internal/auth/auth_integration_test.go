package auth

import (
	"html/template"
	"knowledgeable/internal/users"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"errors"
)

type successUserService struct{}

func (s *successUserService) Login(username, password string) (*users.User, error) {
	return &users.User{ID: 1, Username: "rasmus"}, nil
}

// HAPPY PATH
func TestLogin_SetsSessionCookie(t *testing.T) {

	userService := &successUserService{}
	handler := NewHandler(
		userService,
		template.New("dummy"),
	)

	sessions = map[string]int64{}

	form := url.Values{}
	form.Add("username", "rasmus")
	form.Add("password", "secret")

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", res.StatusCode)
	}

	cookies := res.Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	if cookies[0].Name != "session_id" {
		t.Fatal("session cookie not set")
	}

	sessionID := cookies[0].Value

	userID, ok := Get(sessionID)
	if !ok {
		t.Fatal("session not stored")
	}

	if userID != 1 {
		t.Fatal("wrong user id in session")
	}

	if res.Header.Get("Location") != "/dashboard" {
		t.Fatal("expected redirect to /dashboard")
	}

	if !cookies[0].HttpOnly {
		t.Fatal("cookie should be HttpOnly")
	}
}


// ERROR PATH
type errorUserService struct{}

func (e *errorUserService) Login(username, password string) (*users.User, error) {
	return nil, errors.New("invalid credentials")
}
