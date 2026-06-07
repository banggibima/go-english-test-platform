package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT
			id,
			role_id,
			full_name,
			email,
			password_hash,
			avatar_url,
			is_active,
			is_verified,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.IsActive,
		&user.IsVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id string, fullName string) (*User, error) {
	query := `
		UPDATE users
		SET full_name = $1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING
			id,
			role_id,
			full_name,
			email,
			password_hash,
			avatar_url,
			is_active,
			is_verified,
			last_login_at,
			created_at,
			updated_at
	`

	var user User

	err := r.db.QueryRow(ctx, query, fullName, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.IsActive,
		&user.IsVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
