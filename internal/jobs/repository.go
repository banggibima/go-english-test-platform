package jobs

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAttempt(ctx context.Context, attemptID string) (*Attempt, error) {
	query := `
		SELECT id, user_id, test_id
		FROM attempts
		WHERE id = $1
	`

	var attempt Attempt

	err := r.db.QueryRow(ctx, query, attemptID).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.TestID,
	)

	if err != nil {
		return nil, err
	}

	return &attempt, nil
}

func (r *Repository) FindAnswers(ctx context.Context, attemptID string) ([]AnswerWithQuestion, error) {
	query := `
		SELECT
			aa.id,
			aa.answer_text,
			q.correct_answer,
			q.points
		FROM attempt_answers aa
		JOIN questions q ON q.id = aa.question_id
		WHERE aa.attempt_id = $1
	`

	rows, err := r.db.Query(ctx, query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []AnswerWithQuestion

	for rows.Next() {
		var item AnswerWithQuestion

		if err := rows.Scan(
			&item.AnswerID,
			&item.AnswerText,
			&item.CorrectAnswer,
			&item.Points,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Repository) UpdateAnswerScore(ctx context.Context, answerID string, isCorrect bool, pointsAwarded int) error {
	query := `
		UPDATE attempt_answers
		SET
			is_correct = $1,
			points_awarded = $2,
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, isCorrect, pointsAwarded, answerID)
	return err
}

func (r *Repository) CompleteAttempt(ctx context.Context, attemptID string, score int, maxScore int, totalQuestions int, answeredQuestions int) error {
	query := `
		UPDATE attempts
		SET
			status = 'completed',
			completed_at = NOW(),
			score = $1,
			max_score = $2,
			total_questions = $3,
			answered_questions = $4,
			updated_at = NOW()
		WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query, score, maxScore, totalQuestions, answeredQuestions, attemptID)
	return err
}

func BuildResultData(score int, maxScore int, answeredQuestions int) ([]byte, error) {
	percentage := 0.0
	if maxScore > 0 {
		percentage = math.Round((float64(score)/float64(maxScore))*10000) / 100
	}

	return json.Marshal(map[string]any{
		"score":              score,
		"max_score":          maxScore,
		"percentage":         percentage,
		"answered_questions": answeredQuestions,
		"graded_at":          time.Now(),
	})
}
