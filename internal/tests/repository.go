package tests

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

func (r *Repository) FindAll(ctx context.Context) ([]Test, error) {
	query := `
		SELECT
			id, title, description, test_code, duration_minutes,
			total_questions, passing_score, is_active, created_by,
			created_at, updated_at
		FROM tests
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []Test

	for rows.Next() {
		var test Test

		err := rows.Scan(
			&test.ID,
			&test.Title,
			&test.Description,
			&test.TestCode,
			&test.DurationMinutes,
			&test.TotalQuestions,
			&test.PassingScore,
			&test.IsActive,
			&test.CreatedBy,
			&test.CreatedAt,
			&test.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tests = append(tests, test)
	}

	return tests, rows.Err()
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Test, error) {
	query := `
		SELECT
			id, title, description, test_code, duration_minutes,
			total_questions, passing_score, is_active, created_by,
			created_at, updated_at
		FROM tests
		WHERE id = $1
	`

	var test Test

	err := r.db.QueryRow(ctx, query, id).Scan(
		&test.ID,
		&test.Title,
		&test.Description,
		&test.TestCode,
		&test.DurationMinutes,
		&test.TotalQuestions,
		&test.PassingScore,
		&test.IsActive,
		&test.CreatedBy,
		&test.CreatedAt,
		&test.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &test, nil
}

func (r *Repository) Create(ctx context.Context, test *Test) error {
	query := `
		INSERT INTO tests (
			id, title, description, test_code, duration_minutes,
			total_questions, passing_score, is_active, created_by
		)
		VALUES (
			gen_random_uuid(), $1, $2, $3, $4, 0, $5, TRUE, $6
		)
		RETURNING id, total_questions, is_active, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		test.Title,
		test.Description,
		test.TestCode,
		test.DurationMinutes,
		test.PassingScore,
		test.CreatedBy,
	).Scan(
		&test.ID,
		&test.TotalQuestions,
		&test.IsActive,
		&test.CreatedAt,
		&test.UpdatedAt,
	)
}

func (r *Repository) Update(ctx context.Context, id string, test *Test) (*Test, error) {
	query := `
		UPDATE tests
		SET
			title = $1,
			description = $2,
			duration_minutes = $3,
			passing_score = $4,
			is_active = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING
			id, title, description, test_code, duration_minutes,
			total_questions, passing_score, is_active, created_by,
			created_at, updated_at
	`

	var updated Test

	err := r.db.QueryRow(
		ctx,
		query,
		test.Title,
		test.Description,
		test.DurationMinutes,
		test.PassingScore,
		test.IsActive,
		id,
	).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.TestCode,
		&updated.DurationMinutes,
		&updated.TotalQuestions,
		&updated.PassingScore,
		&updated.IsActive,
		&updated.CreatedBy,
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
		DELETE FROM tests
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
