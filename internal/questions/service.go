package questions

import (
	"context"
	"encoding/json"
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindAll(ctx context.Context) ([]QuestionResponse, error) {
	items, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]QuestionResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toQuestionResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindBySectionID(ctx context.Context, sectionID string) ([]QuestionResponse, error) {
	items, err := s.repository.FindBySectionID(ctx, sectionID)
	if err != nil {
		return nil, err
	}

	responses := make([]QuestionResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toQuestionResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*QuestionResponse, error) {
	question, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if question == nil {
		return nil, errors.New("question not found")
	}

	response := toQuestionResponse(question)
	return &response, nil
}

func (s *Service) Create(ctx context.Context, req CreateQuestionRequest) (*QuestionResponse, error) {
	questionData := req.QuestionData
	if len(questionData) == 0 {
		questionData = json.RawMessage(`{}`)
	}

	question := &Question{
		SectionID:     req.SectionID,
		QuestionType:  req.QuestionType,
		QuestionText:  req.QuestionText,
		QuestionData:  questionData,
		CorrectAnswer: req.CorrectAnswer,
		Explanation:   req.Explanation,
		Points:        req.Points,
		DisplayOrder:  req.DisplayOrder,
	}

	if question.Points == 0 {
		question.Points = 1
	}

	if err := s.repository.Create(ctx, question); err != nil {
		return nil, err
	}

	response := toQuestionResponse(question)
	return &response, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateQuestionRequest) (*QuestionResponse, error) {
	questionData := req.QuestionData
	if len(questionData) == 0 {
		questionData = json.RawMessage(`{}`)
	}

	question := &Question{
		QuestionText:  req.QuestionText,
		QuestionData:  questionData,
		CorrectAnswer: req.CorrectAnswer,
		Explanation:   req.Explanation,
		Points:        req.Points,
		DisplayOrder:  req.DisplayOrder,
		IsActive:      req.IsActive,
	}

	if question.Points == 0 {
		question.Points = 1
	}

	updated, err := s.repository.Update(ctx, id, question)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		return nil, errors.New("question not found")
	}

	response := toQuestionResponse(updated)
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	question, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if question == nil {
		return errors.New("question not found")
	}

	return s.repository.Delete(ctx, id)
}

func toQuestionResponse(question *Question) QuestionResponse {
	return QuestionResponse{
		ID:            question.ID,
		SectionID:     question.SectionID,
		QuestionType:  question.QuestionType,
		QuestionText:  question.QuestionText,
		QuestionData:  json.RawMessage(question.QuestionData),
		CorrectAnswer: question.CorrectAnswer,
		Explanation:   question.Explanation,
		Points:        question.Points,
		DisplayOrder:  question.DisplayOrder,
		IsActive:      question.IsActive,
	}
}
