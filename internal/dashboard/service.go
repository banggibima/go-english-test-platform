package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RepositoryInterface interface {
	GetSummary(ctx context.Context, userID string) (*DashboardResponse, error)
}

type CacheInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type Service struct {
	repository RepositoryInterface
	cache      CacheInterface
}

func NewService(repository RepositoryInterface, cache CacheInterface) *Service {
	return &Service{
		repository: repository,
		cache:      cache,
	}
}

func (s *Service) GetSummary(ctx context.Context, userID string) (*DashboardResponse, error) {
	key := fmt.Sprintf("dashboard:user:%s", userID)

	cached, err := s.cache.Get(ctx, key).Result()
	if err == nil {
		var result DashboardResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	result, err := s.repository.GetSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(result)
	_ = s.cache.Set(ctx, key, bytes, 5*time.Minute).Err()

	return result, nil
}
