package results

import (
	"time"
)

type Result struct {
	ID         string
	AttemptID  string
	UserID     string
	TestID     string
	Score      int
	MaxScore   int
	Percentage float64
	Status     string
	Feedback   *string
	ResultData []byte
	GradedAt   *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
