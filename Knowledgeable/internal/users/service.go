package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
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

	return s.repo.FindByUsername(username)
}

func (s *Service) GetByID(id int64) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	return s.repo.FindById(id)
}

func (s *Service) Login(username, password string) (*User, error) {

	if username == "" || password == "" {
		return nil, errors.New("invalid credentials")
	}

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
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
