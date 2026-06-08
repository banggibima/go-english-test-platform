package users

import (
	"time"
)

type User struct {
	ID           string
	RoleID       string
	FullName     string
	Email        string
	PasswordHash string
	AvatarURL    *string
	IsActive     bool
	IsVerified   bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
