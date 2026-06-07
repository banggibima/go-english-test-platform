package results

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
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, result *Result) error {
	query := `
		INSERT INTO results (
			id,
			attempt_id,
			user_id,
			test_id,
			score,
			max_score,
			percentage,
			status,
			feedback,
			result_data,
			graded_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
		RETURNING
			id,
			created_at,
			updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		result.AttemptID,
		result.UserID,
		result.TestID,
		result.Score,
		result.MaxScore,
		result.Percentage,
		result.Status,
		result.Feedback,
		result.ResultData,
		result.GradedAt,
	).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Result, error) {
	query := `
		SELECT
			id,
			attempt_id,
			user_id,
			test_id,
			score,
			max_score,
			percentage,
			status,
			feedback,
			result_data,
			graded_at,
			created_at,
			updated_at
		FROM results
		WHERE id = $1
	`

	var result Result

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.AttemptID,
		&result.UserID,
		&result.TestID,
		&result.Score,
		&result.MaxScore,
		&result.Percentage,
		&result.Status,
		&result.Feedback,
		&result.ResultData,
		&result.GradedAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *Repository) FindByAttemptID(ctx context.Context, attemptID string) (*Result, error) {
	query := `
		SELECT
			id,
			attempt_id,
			user_id,
			test_id,
			score,
			max_score,
			percentage,
			status,
			feedback,
			result_data,
			graded_at,
			created_at,
			updated_at
		FROM results
		WHERE attempt_id = $1
	`

	var result Result

	err := r.db.QueryRow(ctx, query, attemptID).Scan(
		&result.ID,
		&result.AttemptID,
		&result.UserID,
		&result.TestID,
		&result.Score,
		&result.MaxScore,
		&result.Percentage,
		&result.Status,
		&result.Feedback,
		&result.ResultData,
		&result.GradedAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID string) ([]Result, error) {
	query := `
		SELECT
			id,
			attempt_id,
			user_id,
			test_id,
			score,
			max_score,
			percentage,
			status,
			feedback,
			result_data,
			graded_at,
			created_at,
			updated_at
		FROM results
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result

	for rows.Next() {
		var result Result

		err := rows.Scan(
			&result.ID,
			&result.AttemptID,
			&result.UserID,
			&result.TestID,
			&result.Score,
			&result.MaxScore,
			&result.Percentage,
			&result.Status,
			&result.Feedback,
			&result.ResultData,
			&result.GradedAt,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, rows.Err()
}
