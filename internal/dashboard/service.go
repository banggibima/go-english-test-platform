package dashboard

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetSummary(ctx context.Context, userID string) (*DashboardResponse, error) {
	return s.repository.GetSummary(ctx, userID)
}
