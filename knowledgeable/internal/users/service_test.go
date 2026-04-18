package users

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	user *User
	err  error
}

func (m *mockRepo) Register(u *User) error {
	m.user = u
	return m.err
}

func (m *mockRepo) FindByUsername(username string) (*User, error) {
	return m.user, m.err
}

func (m *mockRepo) FindById(id int64) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.user != nil && m.user.ID == id {
		return m.user, nil
	}
	return nil, nil
}

func (m *mockRepo) UpdatePassword(userID int64, newPasswordHash string) error {
	return m.err
}

func (m *mockRepo) FindAll() ([]User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []User{*m.user}, nil
}

// HAPPY PATH

func TestRegister_HashesPassword(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	_, err := service.Register("rasmus", "mail@test.com", "secret")
	if err != nil {
		t.Fatal(err)
	}

	// repo.user
	if repo.user == nil {
		t.Fatal("repo not called")
	}

	// bcrypt test
	err = bcrypt.CompareHashAndPassword(
		[]byte(repo.user.PasswordHash),
		[]byte("secret"),
	)

	if err != nil {
		t.Fatal("password was not hashed correctly")
	}
}

func TestLogin_Success(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	repo := &mockRepo{
		user: &User{
			Username:     "rasmus",
			PasswordHash: string(hashed),
		},
	}

	service := NewService(repo)

	user, err := service.Login("rasmus", "secret")
	if err != nil {
		t.Fatal(err)
	}

	if user.Username != "rasmus" {
		t.Fatal("wrong user returned")
	}
}

func TestGetByID_Success(t *testing.T) {
	repo := &mockRepo{
		user: &User{ID: 1, Username: "rasmus"},
	}

	service := NewService(repo)

	user, err := service.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 {
		t.Fatal("wrong user returned")
	}
}

// ERROR PATH
func TestRegister_MissingFields(t *testing.T) {
	service := NewService(&mockRepo{})

	_, err := service.Register("", "", "")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRegister_RepoError(t *testing.T) {
	repo := &mockRepo{err: errors.New("db failed")}
	service := NewService(repo)

	_, err := service.Register("rasmus", "mail@test.com", "secret")
	if err == nil {
		t.Fatal("expected repo error")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	repo := &mockRepo{
		user: &User{
			Username:     "rasmus",
			PasswordHash: string(hashed),
		},
	}

	service := NewService(repo)

	_, err := service.Login("rasmus", "wrong")
	if err == nil {
		t.Fatal("expected invalid credentials")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockRepo{user: nil}
	service := NewService(repo)

	_, err := service.Login("rasmus", "secret")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockRepo{user: nil}
	service := NewService(repo)

	_, err := service.GetByID(1)
	if err == nil {
		t.Fatal("expected error")
	}
}
