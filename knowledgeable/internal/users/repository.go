package users

import (
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Register(user *User) error {
	err := r.db.QueryRow(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	log.Println("User added: ", user.Username)
	return nil
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	row := r.db.QueryRow(
		"SELECT id, username, email, password_hash FROM users WHERE username = $1",
		username,
	)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

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
		"SELECT id, username, email, password_hash FROM users WHERE id = $1",
		id,
	)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindAll() ([]User, error) {
	rows, err := r.db.Query(
		"SELECT id, username, email FROM users",
	)
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
