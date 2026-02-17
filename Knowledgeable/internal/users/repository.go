package users

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *User) error {
	result, err := r.db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	row := r.db.QueryRow(
		"SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?",
		username,
	)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindById(id int64) (*User, error) {
	row := r.db.QueryRow(
		"SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?",
		id,
	)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}

	return &user, nil
}
