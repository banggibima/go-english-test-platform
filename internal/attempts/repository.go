package attempts

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

func (r *Repository) Create(ctx context.Context, attempt *Attempt) error {
	query := `
		INSERT INTO attempts (
			id, user_id, test_id, status, total_questions, max_score
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			'in_progress',
			(SELECT total_questions FROM tests WHERE id = $2),
			(SELECT total_questions FROM tests WHERE id = $2)
		)
		RETURNING
			id, status, started_at, total_questions, answered_questions,
			score, max_score, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		attempt.UserID,
		attempt.TestID,
	).Scan(
		&attempt.ID,
		&attempt.Status,
		&attempt.StartedAt,
		&attempt.TotalQuestions,
		&attempt.AnsweredQuestions,
		&attempt.Score,
		&attempt.MaxScore,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
	)
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID string) ([]Attempt, error) {
	query := `
		SELECT
			id, user_id, test_id, status, started_at, submitted_at,
			completed_at, total_questions, answered_questions, score,
			max_score, created_at, updated_at
		FROM attempts
		WHERE user_id = $1
		ORDER BY started_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []Attempt

	for rows.Next() {
		var attempt Attempt

		if err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.TestID,
			&attempt.Status,
			&attempt.StartedAt,
			&attempt.SubmittedAt,
			&attempt.CompletedAt,
			&attempt.TotalQuestions,
			&attempt.AnsweredQuestions,
			&attempt.Score,
			&attempt.MaxScore,
			&attempt.CreatedAt,
			&attempt.UpdatedAt,
		); err != nil {
			return nil, err
		}

		attempts = append(attempts, attempt)
	}

	return attempts, rows.Err()
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Attempt, error) {
	query := `
		SELECT
			id, user_id, test_id, status, started_at, submitted_at,
			completed_at, total_questions, answered_questions, score,
			max_score, created_at, updated_at
		FROM attempts
		WHERE id = $1
	`

	var attempt Attempt

	err := r.db.QueryRow(ctx, query, id).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.TestID,
		&attempt.Status,
		&attempt.StartedAt,
		&attempt.SubmittedAt,
		&attempt.CompletedAt,
		&attempt.TotalQuestions,
		&attempt.AnsweredQuestions,
		&attempt.Score,
		&attempt.MaxScore,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &attempt, nil
}

func (r *Repository) Submit(ctx context.Context, id string, userID string) (*Attempt, error) {
	query := `
		UPDATE attempts
		SET
			status = 'submitted',
			submitted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		  AND user_id = $2
		  AND status = 'in_progress'
		RETURNING
			id, user_id, test_id, status, started_at, submitted_at,
			completed_at, total_questions, answered_questions, score,
			max_score, created_at, updated_at
	`

	var attempt Attempt

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.TestID,
		&attempt.Status,
		&attempt.StartedAt,
		&attempt.SubmittedAt,
		&attempt.CompletedAt,
		&attempt.TotalQuestions,
		&attempt.AnsweredQuestions,
		&attempt.Score,
		&attempt.MaxScore,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &attempt, nil
}
