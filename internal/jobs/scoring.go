package jobs

import (
	"math"
	"strings"
)

type ScoreSummary struct {
	Score             int
	MaxScore          int
	TotalQuestions    int
	AnsweredQuestions int
	Percentage        float64
}

func CalculateScore(answers []AnswerWithQuestion) ScoreSummary {
	score := 0
	maxScore := 0
	totalQuestions := len(answers)
	answeredQuestions := len(answers)

	for _, answer := range answers {
		maxScore += answer.Points

		if answer.AnswerText == nil || answer.CorrectAnswer == nil {
			continue
		}

		userAnswer := strings.TrimSpace(strings.ToLower(*answer.AnswerText))
		correctAnswer := strings.TrimSpace(strings.ToLower(*answer.CorrectAnswer))

		if userAnswer == correctAnswer {
			score += answer.Points
		}
	}

	percentage := 0.0
	if maxScore > 0 {
		percentage = math.Round((float64(score)/float64(maxScore))*10000) / 100
	}

	return ScoreSummary{
		Score:             score,
		MaxScore:          maxScore,
		TotalQuestions:    totalQuestions,
		AnsweredQuestions: answeredQuestions,
		Percentage:        percentage,
	}
}
