package roles

import (
	"context"
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) FindAll(ctx context.Context) ([]RoleResponse, error) {
	items, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]RoleResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toRoleResponse(&item))
	}

	return responses, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*RoleResponse, error) {
	role, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	response := toRoleResponse(role)
	return &response, nil
}

func toRoleResponse(role *Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}
