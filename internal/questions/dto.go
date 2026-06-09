package questions

import (
	"encoding/json"
)

type CreateQuestionRequest struct {
	SectionID     string          `json:"section_id" validate:"required"`
	QuestionType  string          `json:"question_type" validate:"required"`
	QuestionText  string          `json:"question_text" validate:"required"`
	QuestionData  json.RawMessage `json:"question_data" swaggertype:"object"`
	CorrectAnswer *string         `json:"correct_answer"`
	Explanation   *string         `json:"explanation"`
	Points        int             `json:"points" validate:"min=1"`
	DisplayOrder  int             `json:"display_order" validate:"required,min=1"`
}

type UpdateQuestionRequest struct {
	QuestionText  string          `json:"question_text" validate:"required"`
	QuestionData  json.RawMessage `json:"question_data" swaggertype:"object"`
	CorrectAnswer *string         `json:"correct_answer"`
	Explanation   *string         `json:"explanation"`
	Points        int             `json:"points" validate:"min=1"`
	DisplayOrder  int             `json:"display_order" validate:"required,min=1"`
	IsActive      bool            `json:"is_active"`
}

type QuestionResponse struct {
	ID            string          `json:"id"`
	SectionID     string          `json:"section_id"`
	QuestionType  string          `json:"question_type"`
	QuestionText  string          `json:"question_text"`
	QuestionData  json.RawMessage `json:"question_data" swaggertype:"object"`
	CorrectAnswer *string         `json:"correct_answer,omitempty"`
	Explanation   *string         `json:"explanation,omitempty"`
	Points        int             `json:"points"`
	DisplayOrder  int             `json:"display_order"`
	IsActive      bool            `json:"is_active"`
}
