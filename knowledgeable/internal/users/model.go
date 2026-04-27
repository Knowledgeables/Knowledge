package users

import "time"

type User struct {
	ID                   int64
	Username             string
	Email                string
	PasswordHash         string
	ShouldChangePassword bool
}

type PasswordResetToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}