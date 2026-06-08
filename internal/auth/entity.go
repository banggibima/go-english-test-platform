package auth

import (
	"time"
)

type Role struct {
	ID          string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID           string
	RoleID       string
	RoleName     string
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

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
