package users

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

func (s *Service) Me(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return toUserResponse(user), nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*UserResponse, error) {
	user, err := s.repository.UpdateProfile(ctx, userID, req.FullName)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return toUserResponse(user), nil
}

func toUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		RoleID:     user.RoleID,
		FullName:   user.FullName,
		Email:      user.Email,
		AvatarURL:  user.AvatarURL,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}
}
