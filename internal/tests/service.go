package tests

import (
	"context"
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

func (s *Service) FindAll(ctx context.Context) ([]TestResponse, error) {
	items, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]TestResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toTestResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*TestResponse, error) {
	test, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if test == nil {
		return nil, errors.New("test not found")
	}

	response := toTestResponse(test)
	return &response, nil
}

func (s *Service) Create(ctx context.Context, userID string, req CreateTestRequest) (*TestResponse, error) {
	test := &Test{
		Title:           req.Title,
		Description:     req.Description,
		TestCode:        req.TestCode,
		DurationMinutes: req.DurationMinutes,
		PassingScore:    req.PassingScore,
		CreatedBy:       &userID,
	}

	if err := s.repository.Create(ctx, test); err != nil {
		return nil, err
	}

	response := toTestResponse(test)
	return &response, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateTestRequest) (*TestResponse, error) {
	test := &Test{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		PassingScore:    req.PassingScore,
		IsActive:        req.IsActive,
	}

	updated, err := s.repository.Update(ctx, id, test)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		return nil, errors.New("test not found")
	}

	response := toTestResponse(updated)
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	test, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if test == nil {
		return errors.New("test not found")
	}

	return s.repository.Delete(ctx, id)
}

func toTestResponse(test *Test) TestResponse {
	return TestResponse{
		ID:              test.ID,
		Title:           test.Title,
		Description:     test.Description,
		TestCode:        test.TestCode,
		DurationMinutes: test.DurationMinutes,
		TotalQuestions:  test.TotalQuestions,
		PassingScore:    test.PassingScore,
		IsActive:        test.IsActive,
		CreatedBy:       test.CreatedBy,
	}
}
