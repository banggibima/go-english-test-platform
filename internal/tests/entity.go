package tests

import (
	"time"
)

type Test struct {
	ID              string
	Title           string
	Description     *string
	TestCode        string
	DurationMinutes int
	TotalQuestions  int
	PassingScore    int
	IsActive        bool
	CreatedBy       *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
