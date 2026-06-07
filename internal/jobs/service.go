package jobs

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/banggibima/go-english-test-platform/internal/results"
)

type Service struct {
	jobRepository    *Repository
	resultRepository *results.Repository
}

func NewService(jobRepository *Repository, resultRepository *results.Repository) *Service {
	return &Service{
		jobRepository:    jobRepository,
		resultRepository: resultRepository,
	}
}

func (s *Service) HandleScoreAttempt(ctx context.Context, body []byte) error {
	var job ScoreAttemptJob

	if err := json.Unmarshal(body, &job); err != nil {
		return err
	}

	attempt, err := s.jobRepository.FindAttempt(ctx, job.AttemptID)
	if err != nil {
		return err
	}

	answers, err := s.jobRepository.FindAnswers(ctx, job.AttemptID)
	if err != nil {
		return err
	}

	score := 0
	maxScore := 0
	totalQuestions := len(answers)
	answeredQuestions := len(answers)

	for _, answer := range answers {
		maxScore += answer.Points

		isCorrect := false
		pointsAwarded := 0

		if answer.AnswerText != nil && answer.CorrectAnswer != nil {
			userAnswer := strings.TrimSpace(strings.ToLower(*answer.AnswerText))
			correctAnswer := strings.TrimSpace(strings.ToLower(*answer.CorrectAnswer))

			if userAnswer == correctAnswer {
				isCorrect = true
				pointsAwarded = answer.Points
				score += answer.Points
			}
		}

		if err := s.jobRepository.UpdateAnswerScore(ctx, answer.AnswerID, isCorrect, pointsAwarded); err != nil {
			return err
		}
	}

	percentage := 0.0
	if maxScore > 0 {
		percentage = math.Round((float64(score)/float64(maxScore))*10000) / 100
	}

	resultData, err := BuildResultData(score, maxScore, answeredQuestions)
	if err != nil {
		return err
	}

	now := time.Now()

	result := &results.Result{
		AttemptID:  attempt.ID,
		UserID:     attempt.UserID,
		TestID:     attempt.TestID,
		Score:      score,
		MaxScore:   maxScore,
		Percentage: percentage,
		Status:     "completed",
		ResultData: resultData,
		GradedAt:   &now,
	}

	if err := s.resultRepository.Create(ctx, result); err != nil {
		return err
	}

	return s.jobRepository.CompleteAttempt(ctx, attempt.ID, score, maxScore, totalQuestions, answeredQuestions)
}
