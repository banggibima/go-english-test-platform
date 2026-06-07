package auth

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

func (r *Repository) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	var role Role

	err := r.db.QueryRow(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &role, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			id,
			role_id,
			full_name,
			email,
			password_hash,
			is_active,
			is_verified
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			$4,
			TRUE,
			FALSE
		)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		user.RoleID,
		user.FullName,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT
			u.id,
			u.role_id,
			r.name,
			u.full_name,
			u.email,
			u.password_hash,
			u.avatar_url,
			u.is_active,
			u.is_verified,
			u.last_login_at,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.RoleID,
		&user.RoleName,
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

func (r *Repository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *Repository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			id,
			user_id,
			token,
			expires_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3
		)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
	).Scan(
		&token.ID,
		&token.CreatedAt,
	)
}

func (r *Repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `
		SELECT
			id,
			user_id,
			token,
			expires_at,
			revoked_at,
			created_at
		FROM refresh_tokens
		WHERE token = $1
	`

	var refreshToken RefreshToken

	err := r.db.QueryRow(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.RevokedAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &refreshToken, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token = $1
	`

	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT
			u.id,
			u.role_id,
			r.name,
			u.full_name,
			u.email,
			u.password_hash,
			u.avatar_url,
			u.is_active,
			u.is_verified,
			u.last_login_at,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.RoleName,
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
