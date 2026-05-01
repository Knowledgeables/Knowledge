//go:build integration
package auth

import (
	"errors"
	"html/template"
	"knowledgeable/internal/users"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type successUserService struct{}

func (s *successUserService) Login(username, password string) (*users.User, error) {
	return &users.User{ID: 1, Username: "rasmus"}, nil
}

func (s *successUserService) GetByID(id int64) (*users.User, error) {
	return &users.User{
		ID:       id,
		Username: "rasmus",
	}, nil
}

func (s *successUserService) ChangePassword(userID int64, newPassword string) error {
	return nil
}

func (s *successUserService) GetByEmail(email string) (*users.User, error) {
	return &users.User{ID: 1, Email: email}, nil
}

func (s *successUserService) CreateResetToken(userID int64) (string, error) {
	return "token", nil
}

func (s *successUserService) ConsumeResetToken(rawToken string) (int64, error) {
	return 1, nil
}

// HAPPY PATH
func TestLoginAPI_SetsSessionCookie(t *testing.T) {

	userService := &successUserService{}
	handler := NewHandler(userService, func() *template.Template {
		return template.New("dummy")
	}, nil, "")

	sessions = map[string]int64{}

	form := url.Values{}
	form.Add("username", "rasmus")
	form.Add("password", "secret")

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(form.Encode()),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler.LoginAPI(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", res.StatusCode)
	}

	if res.Header.Get("Location") != "/" {
		t.Fatal("expected redirect to /")
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

	if !cookies[0].HttpOnly {
		t.Fatal("cookie should be HttpOnly")
	}

}

//Login failure

type failUserService struct{}

func (s *failUserService) Login(username, password string) (*users.User, error) {
	return nil, errors.New("invalid credentials")
}

func (s *failUserService) GetByID(id int64) (*users.User, error) {
	return nil, errors.New("user not found")
}

func (s *failUserService) ChangePassword(userID int64, newPassword string) error {
	return errors.New("change password failed")
}

func (s *failUserService) GetByEmail(email string) (*users.User, error) {
	return nil, errors.New("user not found")
}

func (s *failUserService) CreateResetToken(userID int64) (string, error) {
	return "", errors.New("failed")
}

func (s *failUserService) ConsumeResetToken(rawToken string) (int64, error) {
	return 0, errors.New("invalid token")
}

type mockEmailSender struct {
	sentTo  string
	err     error
}

func (m *mockEmailSender) Send(to, subject, html string) error {
	m.sentTo = to
	return m.err
}

func TestLoginAPI_InvalidCredentials(t *testing.T) {

	userService := &failUserService{}
	handler := NewHandler(userService, func() *template.Template {
		return template.New("dummy")
	}, nil, "")

	sessions = map[string]int64{}

	form := url.Values{}
	form.Add("username", "rasmus")
	form.Add("password", "wrong")

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler.LoginAPI(rr, req)

	res := rr.Result()

	// should redirect back to login with error param
	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", res.StatusCode)
	}

	if res.Header.Get("Location") != "/login?error=invalid_credentials" {
		t.Fatalf("expected redirect to /login?error=invalid_credentials, got %s", res.Header.Get("Location"))
	}

	// ingen cookie
	if len(res.Cookies()) != 0 {
		t.Fatal("expected no cookies on failed login")
	}

	// ingen session
	if len(sessions) != 0 {
		t.Fatal("expected no session to be created")
	}
}

// ForgotPassword

func TestForgotPassword_UnknownEmail_RedirectsToSent(t *testing.T) {
	handler := NewHandler(&failUserService{}, func() *template.Template {
		return template.New("dummy")
	}, nil, "")

	form := url.Values{}
	form.Add("email", "unknown@test.com")

	req := httptest.NewRequest(http.MethodPost, "/forgot-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ForgotPassword(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", res.StatusCode)
	}
	if res.Header.Get("Location") != "/forgot-password?sent=1" {
		t.Fatalf("expected redirect to ?sent=1, got %s", res.Header.Get("Location"))
	}
}

func TestForgotPassword_KnownUser_SendsEmailAndRedirects(t *testing.T) {
	sender := &mockEmailSender{}
	handler := NewHandler(&successUserService{}, func() *template.Template {
		return template.New("dummy")
	}, sender, "http://localhost:8080")

	form := url.Values{}
	form.Add("email", "rasmus@test.com")

	req := httptest.NewRequest(http.MethodPost, "/forgot-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ForgotPassword(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", res.StatusCode)
	}
	if res.Header.Get("Location") != "/forgot-password?sent=1" {
		t.Fatalf("expected redirect to ?sent=1, got %s", res.Header.Get("Location"))
	}
	if sender.sentTo != "rasmus@test.com" {
		t.Fatalf("expected email sent to rasmus@test.com, got %q", sender.sentTo)
	}
}

func TestForgotPassword_MissingEmail_RedirectsWithError(t *testing.T) {
	handler := NewHandler(&successUserService{}, func() *template.Template {
		return template.New("dummy")
	}, nil, "")

	form := url.Values{}
	form.Add("email", "")

	req := httptest.NewRequest(http.MethodPost, "/forgot-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ForgotPassword(rr, req)
	res := rr.Result()

	if res.Header.Get("Location") != "/forgot-password?error=missing_email" {
		t.Fatalf("expected missing_email redirect, got %s", res.Header.Get("Location"))
	}
}

// ResetPassword

func TestResetPassword_POST_Success(t *testing.T) {
	handler := NewHandler(&successUserService{}, func() *template.Template {
		return template.New("dummy")
	}, nil, "")

	form := url.Values{}
	form.Add("token", "validtoken")
	form.Add("password", "newpassword1")
	form.Add("confirm_password", "newpassword1")

	req := httptest.NewRequest(http.MethodPost, "/reset-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ResetPassword(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", res.StatusCode)
	}
	if res.Header.Get("Location") != "/login?reset=1" {
		t.Fatalf("expected redirect to /login?reset=1, got %s", res.Header.Get("Location"))
	}
}

func TestResetPassword_POST_PasswordMismatch_Redirects(t *testing.T) {
	handler := NewHandler(&successUserService{}, func() *template.Template {
		return template.New("dummy")
	}, nil, "")

	form := url.Values{}
	form.Add("token", "sometoken")
	form.Add("password", "password1")
	form.Add("confirm_password", "password2")

	req := httptest.NewRequest(http.MethodPost, "/reset-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ResetPassword(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", res.StatusCode)
	}
	location := res.Header.Get("Location")
	if !strings.Contains(location, "error=passwords_mismatch") {
		t.Fatalf("expected passwords_mismatch error, got %s", location)
	}
}

func TestResetPassword_POST_InvalidToken_RedirectsToForgot(t *testing.T) {
	handler := NewHandler(&failUserService{}, func() *template.Template {
		return template.New("dummy")
	}, nil, "")

	form := url.Values{}
	form.Add("token", "expiredtoken")
	form.Add("password", "newpassword1")
	form.Add("confirm_password", "newpassword1")

	req := httptest.NewRequest(http.MethodPost, "/reset-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ResetPassword(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", res.StatusCode)
	}
	if res.Header.Get("Location") != "/forgot-password?error=invalid_token" {
		t.Fatalf("expected redirect to /forgot-password?error=invalid_token, got %s", res.Header.Get("Location"))
	}
}
