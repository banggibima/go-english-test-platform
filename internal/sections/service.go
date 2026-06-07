package sections

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

func (s *Service) FindAll(ctx context.Context) ([]SectionResponse, error) {
	items, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]SectionResponse, 0, len(items))

	for _, item := range items {
		responses = append(responses, toSectionResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindByTestID(ctx context.Context, testID string) ([]SectionResponse, error) {
	items, err := s.repository.FindByTestID(ctx, testID)
	if err != nil {
		return nil, err
	}

	responses := make([]SectionResponse, 0, len(items))

	for _, item := range items {
		responses = append(responses, toSectionResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*SectionResponse, error) {
	section, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if section == nil {
		return nil, errors.New("section not found")
	}

	response := toSectionResponse(section)

	return &response, nil
}

func (s *Service) Create(ctx context.Context, req CreateSectionRequest) (*SectionResponse, error) {
	section := &Section{
		TestID:          req.TestID,
		Title:           req.Title,
		SectionCode:     req.SectionCode,
		DisplayOrder:    req.DisplayOrder,
		DurationMinutes: req.DurationMinutes,
	}

	if err := s.repository.Create(ctx, section); err != nil {
		return nil, err
	}

	response := toSectionResponse(section)

	return &response, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateSectionRequest) (*SectionResponse, error) {
	section := &Section{
		Title:           req.Title,
		DisplayOrder:    req.DisplayOrder,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
	}

	updated, err := s.repository.Update(ctx, id, section)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		return nil, errors.New("section not found")
	}

	response := toSectionResponse(updated)

	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	section, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if section == nil {
		return errors.New("section not found")
	}

	return s.repository.Delete(ctx, id)
}

func toSectionResponse(section *Section) SectionResponse {
	return SectionResponse{
		ID:              section.ID,
		TestID:          section.TestID,
		Title:           section.Title,
		SectionCode:     section.SectionCode,
		DisplayOrder:    section.DisplayOrder,
		DurationMinutes: section.DurationMinutes,
		TotalQuestions:  section.TotalQuestions,
		IsActive:        section.IsActive,
	}
}
