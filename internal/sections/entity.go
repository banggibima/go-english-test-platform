package sections

import (
	"time"
)

type Section struct {
	ID              string
	TestID          string
	Title           string
	SectionCode     string
	DisplayOrder    int
	DurationMinutes int
	TotalQuestions  int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
