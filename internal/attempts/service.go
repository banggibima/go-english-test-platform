package attempts

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/banggibima/go-english-test-platform/internal/jobs"
	"github.com/banggibima/go-english-test-platform/pkg/queue"
)

type Service struct {
	repository *Repository
	rabbit     *queue.RabbitMQ
}

func NewService(repository *Repository, rabbit *queue.RabbitMQ) *Service {
	return &Service{
		repository: repository,
		rabbit:     rabbit,
	}
}

func (s *Service) StartAttempt(ctx context.Context, userID string, req StartAttemptRequest) (*AttemptResponse, error) {
	attempt := &Attempt{
		UserID: userID,
		TestID: req.TestID,
	}

	if err := s.repository.Create(ctx, attempt); err != nil {
		return nil, err
	}

	response := toAttemptResponse(attempt)
	return &response, nil
}

func (s *Service) FindAllByUserID(ctx context.Context, userID string) ([]AttemptResponse, error) {
	items, err := s.repository.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]AttemptResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toAttemptResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindByID(ctx context.Context, id string, userID string) (*AttemptResponse, error) {
	attempt, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if attempt == nil {
		return nil, errors.New("attempt not found")
	}

	if attempt.UserID != userID {
		return nil, errors.New("forbidden")
	}

	response := toAttemptResponse(attempt)
	return &response, nil
}

func (s *Service) Submit(ctx context.Context, id string, userID string) (*AttemptResponse, error) {
	attempt, err := s.repository.Submit(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if attempt == nil {
		return nil, errors.New("attempt not found or already submitted")
	}

	job := jobs.ScoreAttemptJob{
		AttemptID: attempt.ID,
		UserID:    attempt.UserID,
		TestID:    attempt.TestID,
	}

	body, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	if err := s.rabbit.Publish("score-attempt", body); err != nil {
		return nil, err
	}

	response := toAttemptResponse(attempt)
	return &response, nil
}

func toAttemptResponse(attempt *Attempt) AttemptResponse {
	return AttemptResponse{
		ID:                attempt.ID,
		UserID:            attempt.UserID,
		TestID:            attempt.TestID,
		Status:            attempt.Status,
		TotalQuestions:    attempt.TotalQuestions,
		AnsweredQuestions: attempt.AnsweredQuestions,
		Score:             attempt.Score,
		MaxScore:          attempt.MaxScore,
		StartedAt:         formatTime(attempt.StartedAt),
		SubmittedAt:       formatTimePtr(attempt.SubmittedAt),
		CompletedAt:       formatTimePtr(attempt.CompletedAt),
	}
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	formatted := t.Format(time.RFC3339)
	return &formatted
}
