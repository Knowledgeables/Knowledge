package users

import (
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/lib/pq"
)

var (
	ErrUsernameTaken = errors.New("username already taken")
	ErrEmailTaken    = errors.New("email already taken")
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
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Constraint, "username") {
				return ErrUsernameTaken
			}
			if strings.Contains(pqErr.Constraint, "email") {
				return ErrEmailTaken
			}
		}
		return err
	}

	slog.Info("user added", "user_id", user.ID) // #nosec G706 -- JSON handler escapes all values
	return nil
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	row := r.db.QueryRow(
		"SELECT id, username, email, password_hash, should_change_password FROM users WHERE username = $1",
		username,
	)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.ShouldChangePassword)

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
		"SELECT id, username, email, password_hash, should_change_password FROM users WHERE id = $1",
		id,
	)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.ShouldChangePassword)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdatePassword(userID int64, newPasswordHash string) error {
	_, err := r.db.Exec(
		"UPDATE users SET password_hash = $1, should_change_password = false WHERE id = $2",
		newPasswordHash, userID,
	)
	return err
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	row := r.db.QueryRow(
		"SELECT id, username, email, password_hash, should_change_password FROM users WHERE email = $1",
		email,
	)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.ShouldChangePassword)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreatePasswordResetToken(userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)",
		userID, tokenHash, expiresAt,
	)
	return err
}

func (r *Repository) FindPasswordResetToken(tokenHash string) (*PasswordResetToken, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, token_hash, expires_at, used_at, created_at FROM password_reset_tokens WHERE token_hash = $1",
		tokenHash,
	)

	var token PasswordResetToken
	var usedAt sql.NullTime
	err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &usedAt, &token.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if usedAt.Valid {
		token.UsedAt = &usedAt.Time
	}

	return &token, nil
}

func (r *Repository) MarkTokenUsed(id int64) error {
	_, err := r.db.Exec(
		"UPDATE password_reset_tokens SET used_at = NOW() WHERE id = $1",
		id,
	)
	return err
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
