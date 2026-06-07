package attempts

import "time"

type Attempt struct {
	ID                string
	UserID            string
	TestID            string
	Status            string
	StartedAt         time.Time
	SubmittedAt       *time.Time
	CompletedAt       *time.Time
	TotalQuestions    int
	AnsweredQuestions int
	Score             *int
	MaxScore          int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
