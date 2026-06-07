package jobs

type Attempt struct {
	ID     string
	UserID string
	TestID string
}

type AnswerWithQuestion struct {
	AnswerID      string
	AnswerText    *string
	CorrectAnswer *string
	Points        int
}
