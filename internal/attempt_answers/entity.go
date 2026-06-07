package attemptanswers

import "time"

type AttemptAnswer struct {
	ID            string
	AttemptID     string
	QuestionID    string
	AnswerText    *string
	AnswerData    []byte
	IsCorrect     *bool
	PointsAwarded int
	AnsweredAt    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
