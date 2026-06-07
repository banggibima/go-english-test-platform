package questions

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

func (r *Repository) FindAll(ctx context.Context) ([]Question, error) {
	query := `
		SELECT
			id, section_id, question_type, question_text, question_data,
			correct_answer, explanation, points, display_order, is_active,
			created_at, updated_at
		FROM questions
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question

	for rows.Next() {
		var question Question

		if err := rows.Scan(
			&question.ID,
			&question.SectionID,
			&question.QuestionType,
			&question.QuestionText,
			&question.QuestionData,
			&question.CorrectAnswer,
			&question.Explanation,
			&question.Points,
			&question.DisplayOrder,
			&question.IsActive,
			&question.CreatedAt,
			&question.UpdatedAt,
		); err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	return questions, rows.Err()
}

func (r *Repository) FindBySectionID(ctx context.Context, sectionID string) ([]Question, error) {
	query := `
		SELECT
			id, section_id, question_type, question_text, question_data,
			correct_answer, explanation, points, display_order, is_active,
			created_at, updated_at
		FROM questions
		WHERE section_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(ctx, query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question

	for rows.Next() {
		var question Question

		if err := rows.Scan(
			&question.ID,
			&question.SectionID,
			&question.QuestionType,
			&question.QuestionText,
			&question.QuestionData,
			&question.CorrectAnswer,
			&question.Explanation,
			&question.Points,
			&question.DisplayOrder,
			&question.IsActive,
			&question.CreatedAt,
			&question.UpdatedAt,
		); err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	return questions, rows.Err()
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Question, error) {
	query := `
		SELECT
			id, section_id, question_type, question_text, question_data,
			correct_answer, explanation, points, display_order, is_active,
			created_at, updated_at
		FROM questions
		WHERE id = $1
	`

	var question Question

	err := r.db.QueryRow(ctx, query, id).Scan(
		&question.ID,
		&question.SectionID,
		&question.QuestionType,
		&question.QuestionText,
		&question.QuestionData,
		&question.CorrectAnswer,
		&question.Explanation,
		&question.Points,
		&question.DisplayOrder,
		&question.IsActive,
		&question.CreatedAt,
		&question.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &question, nil
}

func (r *Repository) Create(ctx context.Context, question *Question) error {
	query := `
		INSERT INTO questions (
			id, section_id, question_type, question_text, question_data,
			correct_answer, explanation, points, display_order, is_active
		)
		VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, TRUE
		)
		RETURNING id, is_active, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		question.SectionID,
		question.QuestionType,
		question.QuestionText,
		question.QuestionData,
		question.CorrectAnswer,
		question.Explanation,
		question.Points,
		question.DisplayOrder,
	).Scan(
		&question.ID,
		&question.IsActive,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
}

func (r *Repository) Update(ctx context.Context, id string, question *Question) (*Question, error) {
	query := `
		UPDATE questions
		SET
			question_text = $1,
			question_data = $2,
			correct_answer = $3,
			explanation = $4,
			points = $5,
			display_order = $6,
			is_active = $7,
			updated_at = NOW()
		WHERE id = $8
		RETURNING
			id, section_id, question_type, question_text, question_data,
			correct_answer, explanation, points, display_order, is_active,
			created_at, updated_at
	`

	var updated Question

	err := r.db.QueryRow(
		ctx,
		query,
		question.QuestionText,
		question.QuestionData,
		question.CorrectAnswer,
		question.Explanation,
		question.Points,
		question.DisplayOrder,
		question.IsActive,
		id,
	).Scan(
		&updated.ID,
		&updated.SectionID,
		&updated.QuestionType,
		&updated.QuestionText,
		&updated.QuestionData,
		&updated.CorrectAnswer,
		&updated.Explanation,
		&updated.Points,
		&updated.DisplayOrder,
		&updated.IsActive,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &updated, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM questions
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
