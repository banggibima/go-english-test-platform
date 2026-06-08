package files

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

func (r *Repository) Create(ctx context.Context, file *File) error {
	query := `
		INSERT INTO files (
			id, user_id, bucket, object_name, original_name,
			content_type, size_bytes, file_url
		)
		VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		file.UserID,
		file.Bucket,
		file.ObjectName,
		file.OriginalName,
		file.ContentType,
		file.SizeBytes,
		file.FileURL,
	).Scan(
		&file.ID,
		&file.CreatedAt,
		&file.UpdatedAt,
	)
}

func (r *Repository) FindByID(ctx context.Context, id string) (*File, error) {
	query := `
		SELECT
			id, user_id, bucket, object_name, original_name,
			content_type, size_bytes, file_url, created_at, updated_at
		FROM files
		WHERE id = $1
	`

	var file File

	err := r.db.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.UserID,
		&file.Bucket,
		&file.ObjectName,
		&file.OriginalName,
		&file.ContentType,
		&file.SizeBytes,
		&file.FileURL,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &file, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID string) ([]File, error) {
	query := `
		SELECT
			id, user_id, bucket, object_name, original_name,
			content_type, size_bytes, file_url, created_at, updated_at
		FROM files
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File

		if err := rows.Scan(
			&file.ID,
			&file.UserID,
			&file.Bucket,
			&file.ObjectName,
			&file.OriginalName,
			&file.ContentType,
			&file.SizeBytes,
			&file.FileURL,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, rows.Err()
}
