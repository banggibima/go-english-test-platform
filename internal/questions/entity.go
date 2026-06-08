package questions

import (
	"time"
)

type Question struct {
	ID            string
	SectionID     string
	QuestionType  string
	QuestionText  string
	QuestionData  []byte
	CorrectAnswer *string
	Explanation   *string
	Points        int
	DisplayOrder  int
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
