package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Register(*User) error
	FindByUsername(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindById(int64) (*User, error)
	FindAll() ([]User, error)
	UpdatePassword(int64, string) error
	CreatePasswordResetToken(userID int64, tokenHash string, expiresAt time.Time) error
	FindPasswordResetToken(tokenHash string) (*PasswordResetToken, error)
	MarkTokenUsed(id int64) error
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(username, email, rawPassword string) (*User, error) {

	if username == "" || email == "" || rawPassword == "" {
		return nil, errors.New("missing fields")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashed),
	}

	if err := s.repo.Register(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByUsername(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("missing username")
	}

	user, err := s.repo.FindByUsername(username)

	return user, err
}

func (s *Service) GetByID(id int64) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	user, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *Service) Login(username, password string) (*User, error) {

	if username == "" || password == "" {
		return nil, errors.New("invalid credentials")
	}

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *Service) ChangePassword(userID int64, newPassword string) error {
	if newPassword == "" {
		return errors.New("missing password")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(userID, string(hashed))
}

func (s *Service) GetByEmail(email string) (*User, error) {
	if email == "" {
		return nil, errors.New("missing email")
	}
	return s.repo.FindByEmail(email)
}

func (s *Service) CreateResetToken(userID int64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	rawToken := hex.EncodeToString(b)

	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])

	if err := s.repo.CreatePasswordResetToken(userID, tokenHash, time.Now().Add(time.Hour)); err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *Service) ConsumeResetToken(rawToken string) (int64, error) {
	if rawToken == "" {
		return 0, errors.New("missing token")
	}

	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])

	token, err := s.repo.FindPasswordResetToken(tokenHash)
	if err != nil {
		return 0, err
	}
	if token == nil {
		return 0, errors.New("invalid token")
	}
	if time.Now().After(token.ExpiresAt) {
		return 0, errors.New("token expired")
	}
	if token.UsedAt != nil {
		return 0, errors.New("token already used")
	}

	if err := s.repo.MarkTokenUsed(token.ID); err != nil {
		return 0, err
	}

	return token.UserID, nil
}

func (s *Service) GetAll() ([]User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Ekstra sikkerhed: nulstil hash hvis structen stadig har feltet
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}
