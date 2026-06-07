package attemptanswers

import (
	"context"
	"encoding/json"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) SaveOrUpdate(ctx context.Context, req SaveAnswerRequest) (*AttemptAnswerResponse, error) {
	answerData := req.AnswerData
	if len(answerData) == 0 {
		answerData = json.RawMessage(`{}`)
	}

	answer := &AttemptAnswer{
		AttemptID:  req.AttemptID,
		QuestionID: req.QuestionID,
		AnswerText: req.AnswerText,
		AnswerData: answerData,
	}

	if err := s.repository.SaveOrUpdate(ctx, answer); err != nil {
		return nil, err
	}

	response := toAttemptAnswerResponse(answer)
	return &response, nil
}

func (s *Service) FindByAttemptID(ctx context.Context, attemptID string) ([]AttemptAnswerResponse, error) {
	items, err := s.repository.FindByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	responses := make([]AttemptAnswerResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toAttemptAnswerResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindByAttemptAndQuestion(ctx context.Context, attemptID string, questionID string) (*AttemptAnswerResponse, error) {
	answer, err := s.repository.FindByAttemptAndQuestion(ctx, attemptID, questionID)
	if err != nil {
		return nil, err
	}

	if answer == nil {
		return nil, nil
	}

	response := toAttemptAnswerResponse(answer)
	return &response, nil
}

func toAttemptAnswerResponse(answer *AttemptAnswer) AttemptAnswerResponse {
	return AttemptAnswerResponse{
		ID:            answer.ID,
		AttemptID:     answer.AttemptID,
		QuestionID:    answer.QuestionID,
		AnswerText:    answer.AnswerText,
		AnswerData:    json.RawMessage(answer.AnswerData),
		IsCorrect:     answer.IsCorrect,
		PointsAwarded: answer.PointsAwarded,
		AnsweredAt:    answer.AnsweredAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
