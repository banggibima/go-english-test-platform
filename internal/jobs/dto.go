package jobs

type ScoreAttemptJob struct {
	AttemptID string `json:"attempt_id"`
	UserID    string `json:"user_id"`
	TestID    string `json:"test_id"`
}
