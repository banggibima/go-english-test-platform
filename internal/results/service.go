package results

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindByID(ctx context.Context, id string, userID string) (*ResultResponse, error) {
	result, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("result not found")
	}

	if result.UserID != userID {
		return nil, errors.New("forbidden")
	}

	response := toResultResponse(result)
	return &response, nil
}

func (s *Service) FindByAttemptID(ctx context.Context, attemptID string, userID string) (*ResultResponse, error) {
	result, err := s.repository.FindByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("result not found")
	}

	if result.UserID != userID {
		return nil, errors.New("forbidden")
	}

	response := toResultResponse(result)
	return &response, nil
}

func (s *Service) FindByUserID(ctx context.Context, userID string) ([]ResultResponse, error) {
	items, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]ResultResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toResultResponse(&item))
	}

	return responses, nil
}

func toResultResponse(result *Result) ResultResponse {
	return ResultResponse{
		ID:         result.ID,
		AttemptID:  result.AttemptID,
		UserID:     result.UserID,
		TestID:     result.TestID,
		Score:      result.Score,
		MaxScore:   result.MaxScore,
		Percentage: result.Percentage,
		Status:     result.Status,
		Feedback:   result.Feedback,
		ResultData: json.RawMessage(result.ResultData),
		GradedAt:   formatTimePtr(result.GradedAt),
	}
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	formatted := t.Format(time.RFC3339)
	return &formatted
}
