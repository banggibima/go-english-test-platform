package sections

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

func (r *Repository) FindAll(ctx context.Context) ([]Section, error) {
	query := `
		SELECT
			id, test_id, title, section_code, display_order,
			duration_minutes, total_questions, is_active,
			created_at, updated_at
		FROM sections
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []Section

	for rows.Next() {
		var section Section

		err := rows.Scan(
			&section.ID,
			&section.TestID,
			&section.Title,
			&section.SectionCode,
			&section.DisplayOrder,
			&section.DurationMinutes,
			&section.TotalQuestions,
			&section.IsActive,
			&section.CreatedAt,
			&section.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	return sections, rows.Err()
}

func (r *Repository) FindByTestID(ctx context.Context, testID string) ([]Section, error) {
	query := `
		SELECT
			id, test_id, title, section_code, display_order,
			duration_minutes, total_questions, is_active,
			created_at, updated_at
		FROM sections
		WHERE test_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(ctx, query, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []Section

	for rows.Next() {
		var section Section

		err := rows.Scan(
			&section.ID,
			&section.TestID,
			&section.Title,
			&section.SectionCode,
			&section.DisplayOrder,
			&section.DurationMinutes,
			&section.TotalQuestions,
			&section.IsActive,
			&section.CreatedAt,
			&section.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	return sections, rows.Err()
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Section, error) {
	query := `
		SELECT
			id, test_id, title, section_code, display_order,
			duration_minutes, total_questions, is_active,
			created_at, updated_at
		FROM sections
		WHERE id = $1
	`

	var section Section

	err := r.db.QueryRow(ctx, query, id).Scan(
		&section.ID,
		&section.TestID,
		&section.Title,
		&section.SectionCode,
		&section.DisplayOrder,
		&section.DurationMinutes,
		&section.TotalQuestions,
		&section.IsActive,
		&section.CreatedAt,
		&section.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &section, nil
}

func (r *Repository) Create(ctx context.Context, section *Section) error {
	query := `
		INSERT INTO sections (
			id, test_id, title, section_code, display_order,
			duration_minutes, total_questions, is_active
		)
		VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, 0, TRUE
		)
		RETURNING id, total_questions, is_active, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		section.TestID,
		section.Title,
		section.SectionCode,
		section.DisplayOrder,
		section.DurationMinutes,
	).Scan(
		&section.ID,
		&section.TotalQuestions,
		&section.IsActive,
		&section.CreatedAt,
		&section.UpdatedAt,
	)
}

func (r *Repository) Update(ctx context.Context, id string, section *Section) (*Section, error) {
	query := `
		UPDATE sections
		SET
			title = $1,
			display_order = $2,
			duration_minutes = $3,
			is_active = $4,
			updated_at = NOW()
		WHERE id = $5
		RETURNING
			id, test_id, title, section_code, display_order,
			duration_minutes, total_questions, is_active,
			created_at, updated_at
	`

	var updated Section

	err := r.db.QueryRow(
		ctx,
		query,
		section.Title,
		section.DisplayOrder,
		section.DurationMinutes,
		section.IsActive,
		id,
	).Scan(
		&updated.ID,
		&updated.TestID,
		&updated.Title,
		&updated.SectionCode,
		&updated.DisplayOrder,
		&updated.DurationMinutes,
		&updated.TotalQuestions,
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
		DELETE FROM sections
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
