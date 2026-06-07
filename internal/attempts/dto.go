package attempts

type StartAttemptRequest struct {
	TestID string `json:"test_id" validate:"required"`
}

type AttemptResponse struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	TestID            string  `json:"test_id"`
	Status            string  `json:"status"`
	TotalQuestions    int     `json:"total_questions"`
	AnsweredQuestions int     `json:"answered_questions"`
	Score             *int    `json:"score,omitempty"`
	MaxScore          int     `json:"max_score"`
	StartedAt         string  `json:"started_at"`
	SubmittedAt       *string `json:"submitted_at,omitempty"`
	CompletedAt       *string `json:"completed_at,omitempty"`
}
