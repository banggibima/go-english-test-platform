package results

import (
	"encoding/json"
)

type ResultResponse struct {
	ID         string          `json:"id"`
	AttemptID  string          `json:"attempt_id"`
	UserID     string          `json:"user_id"`
	TestID     string          `json:"test_id"`
	Score      int             `json:"score"`
	MaxScore   int             `json:"max_score"`
	Percentage float64         `json:"percentage"`
	Status     string          `json:"status"`
	Feedback   *string         `json:"feedback,omitempty"`
	ResultData json.RawMessage `json:"result_data"`
	GradedAt   *string         `json:"graded_at,omitempty"`
}
