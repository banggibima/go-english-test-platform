package dashboard

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetSummary(ctx context.Context, userID string) (*DashboardResponse, error) {
	var result DashboardResponse

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM tests`).Scan(&result.TotalTests); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM attempts
		WHERE user_id = $1
	`, userID).Scan(&result.TotalAttempts); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM attempts
		WHERE user_id = $1
		  AND status = 'completed'
	`, userID).Scan(&result.CompletedAttempts); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx, `
		SELECT COALESCE(AVG(percentage), 0)
		FROM results
		WHERE user_id = $1
	`, userID).Scan(&result.AverageScore); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			test_id,
			attempt_id,
			score,
			max_score,
			percentage,
			status
		FROM results
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 5
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recentResults := make([]RecentResult, 0)

	for rows.Next() {
		var item RecentResult

		if err := rows.Scan(
			&item.ResultID,
			&item.TestID,
			&item.AttemptID,
			&item.Score,
			&item.MaxScore,
			&item.Percentage,
			&item.Status,
		); err != nil {
			return nil, err
		}

		recentResults = append(recentResults, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result.RecentResults = recentResults

	return &result, nil
}
