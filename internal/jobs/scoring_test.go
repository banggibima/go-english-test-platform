package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func stringPtr(value string) *string {
	return &value
}

func TestCalculateScoreCorrectAnswer(t *testing.T) {
	answers := []AnswerWithQuestion{
		{
			AnswerText:    stringPtr("A"),
			CorrectAnswer: stringPtr("A"),
			Points:        1,
		},
	}

	result := CalculateScore(answers)

	assert.Equal(t, 1, result.Score)
	assert.Equal(t, 1, result.MaxScore)
	assert.Equal(t, 1, result.TotalQuestions)
	assert.Equal(t, 1, result.AnsweredQuestions)
	assert.Equal(t, 100.0, result.Percentage)
}

func TestCalculateScoreWrongAnswer(t *testing.T) {
	answers := []AnswerWithQuestion{
		{
			AnswerText:    stringPtr("B"),
			CorrectAnswer: stringPtr("A"),
			Points:        1,
		},
	}

	result := CalculateScore(answers)

	assert.Equal(t, 0, result.Score)
	assert.Equal(t, 1, result.MaxScore)
	assert.Equal(t, 1, result.TotalQuestions)
	assert.Equal(t, 1, result.AnsweredQuestions)
	assert.Equal(t, 0.0, result.Percentage)
}

func TestCalculateScoreCaseInsensitiveAndTrimmed(t *testing.T) {
	answers := []AnswerWithQuestion{
		{
			AnswerText:    stringPtr(" a "),
			CorrectAnswer: stringPtr("A"),
			Points:        2,
		},
	}

	result := CalculateScore(answers)

	assert.Equal(t, 2, result.Score)
	assert.Equal(t, 2, result.MaxScore)
	assert.Equal(t, 100.0, result.Percentage)
}

func TestCalculateScoreMultipleAnswers(t *testing.T) {
	answers := []AnswerWithQuestion{
		{
			AnswerText:    stringPtr("A"),
			CorrectAnswer: stringPtr("A"),
			Points:        1,
		},
		{
			AnswerText:    stringPtr("B"),
			CorrectAnswer: stringPtr("C"),
			Points:        1,
		},
		{
			AnswerText:    stringPtr("D"),
			CorrectAnswer: stringPtr("D"),
			Points:        2,
		},
	}

	result := CalculateScore(answers)

	assert.Equal(t, 3, result.Score)
	assert.Equal(t, 4, result.MaxScore)
	assert.Equal(t, 3, result.TotalQuestions)
	assert.Equal(t, 3, result.AnsweredQuestions)
	assert.Equal(t, 75.0, result.Percentage)
}

func TestCalculateScoreNilAnswer(t *testing.T) {
	answers := []AnswerWithQuestion{
		{
			AnswerText:    nil,
			CorrectAnswer: stringPtr("A"),
			Points:        1,
		},
	}

	result := CalculateScore(answers)

	assert.Equal(t, 0, result.Score)
	assert.Equal(t, 1, result.MaxScore)
	assert.Equal(t, 0.0, result.Percentage)
}

func TestCalculateScoreEmptyAnswers(t *testing.T) {
	result := CalculateScore([]AnswerWithQuestion{})

	assert.Equal(t, 0, result.Score)
	assert.Equal(t, 0, result.MaxScore)
	assert.Equal(t, 0, result.TotalQuestions)
	assert.Equal(t, 0, result.AnsweredQuestions)
	assert.Equal(t, 0.0, result.Percentage)
}
