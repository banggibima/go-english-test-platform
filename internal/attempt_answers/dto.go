package attempt_answers

import (
	"encoding/json"
)

type SaveAnswerRequest struct {
	AttemptID  string          `json:"attempt_id" validate:"required"`
	QuestionID string          `json:"question_id" validate:"required"`
	AnswerText *string         `json:"answer_text"`
	AnswerData json.RawMessage `json:"answer_data,omitempty" swaggertype:"object"`
}

type AttemptAnswerResponse struct {
	ID            string          `json:"id"`
	AttemptID     string          `json:"attempt_id"`
	QuestionID    string          `json:"question_id"`
	AnswerText    *string         `json:"answer_text,omitempty"`
	AnswerData    json.RawMessage `json:"answer_data" swaggertype:"object"`
	IsCorrect     *bool           `json:"is_correct,omitempty"`
	PointsAwarded int             `json:"points_awarded"`
	AnsweredAt    string          `json:"answered_at"`
}
