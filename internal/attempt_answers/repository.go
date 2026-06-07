package attemptanswers

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

func (r *Repository) SaveOrUpdate(ctx context.Context, answer *AttemptAnswer) error {
	query := `
		INSERT INTO attempt_answers (
			id,
			attempt_id,
			question_id,
			answer_text,
			answer_data
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			$4
		)
		ON CONFLICT (attempt_id, question_id)
		DO UPDATE SET
			answer_text = EXCLUDED.answer_text,
			answer_data = EXCLUDED.answer_data,
			answered_at = NOW(),
			updated_at = NOW()
		RETURNING
			id,
			is_correct,
			points_awarded,
			answered_at,
			created_at,
			updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		answer.AttemptID,
		answer.QuestionID,
		answer.AnswerText,
		answer.AnswerData,
	).Scan(
		&answer.ID,
		&answer.IsCorrect,
		&answer.PointsAwarded,
		&answer.AnsweredAt,
		&answer.CreatedAt,
		&answer.UpdatedAt,
	)
}

func (r *Repository) FindByAttemptID(ctx context.Context, attemptID string) ([]AttemptAnswer, error) {
	query := `
		SELECT
			id,
			attempt_id,
			question_id,
			answer_text,
			answer_data,
			is_correct,
			points_awarded,
			answered_at,
			created_at,
			updated_at
		FROM attempt_answers
		WHERE attempt_id = $1
		ORDER BY answered_at ASC
	`

	rows, err := r.db.Query(ctx, query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []AttemptAnswer

	for rows.Next() {
		var answer AttemptAnswer

		err := rows.Scan(
			&answer.ID,
			&answer.AttemptID,
			&answer.QuestionID,
			&answer.AnswerText,
			&answer.AnswerData,
			&answer.IsCorrect,
			&answer.PointsAwarded,
			&answer.AnsweredAt,
			&answer.CreatedAt,
			&answer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, rows.Err()
}

func (r *Repository) FindByAttemptAndQuestion(
	ctx context.Context,
	attemptID string,
	questionID string,
) (*AttemptAnswer, error) {
	query := `
		SELECT
			id,
			attempt_id,
			question_id,
			answer_text,
			answer_data,
			is_correct,
			points_awarded,
			answered_at,
			created_at,
			updated_at
		FROM attempt_answers
		WHERE attempt_id = $1
		  AND question_id = $2
	`

	var answer AttemptAnswer

	err := r.db.QueryRow(
		ctx,
		query,
		attemptID,
		questionID,
	).Scan(
		&answer.ID,
		&answer.AttemptID,
		&answer.QuestionID,
		&answer.AnswerText,
		&answer.AnswerData,
		&answer.IsCorrect,
		&answer.PointsAwarded,
		&answer.AnsweredAt,
		&answer.CreatedAt,
		&answer.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &answer, nil
}
