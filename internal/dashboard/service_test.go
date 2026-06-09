package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDashboardRepository struct {
	mock.Mock
}

func (m *mockDashboardRepository) GetSummary(ctx context.Context, userID string) (*DashboardResponse, error) {
	args := m.Called(ctx, userID)

	result, _ := args.Get(0).(*DashboardResponse)
	return result, args.Error(1)
}

type mockDashboardCache struct {
	mock.Mock
}

func (m *mockDashboardCache) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)

	cmd := redis.NewStringCmd(ctx)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
		return cmd
	}

	cmd.SetVal(args.String(0))
	return cmd
}

func (m *mockDashboardCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)

	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
		return cmd
	}

	cmd.SetVal("OK")
	return cmd
}

func TestGetSummaryCacheHit(t *testing.T) {
	repo := new(mockDashboardRepository)
	cache := new(mockDashboardCache)
	service := NewService(repo, cache)

	userID := "user-id"

	expected := DashboardResponse{
		TotalTests:        3,
		TotalAttempts:     5,
		CompletedAttempts: 2,
		AverageScore:      80,
		RecentResults: []RecentResult{
			{
				ResultID:   "result-id",
				TestID:     "test-id",
				AttemptID:  "attempt-id",
				Score:      8,
				MaxScore:   10,
				Percentage: 80,
				Status:     "completed",
			},
		},
	}

	bytes, _ := json.Marshal(expected)

	cache.On("Get", mock.Anything, "dashboard:user:"+userID).
		Return(string(bytes), nil)

	result, err := service.GetSummary(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected.TotalTests, result.TotalTests)
	assert.Equal(t, expected.TotalAttempts, result.TotalAttempts)
	assert.Equal(t, expected.CompletedAttempts, result.CompletedAttempts)
	assert.Equal(t, expected.AverageScore, result.AverageScore)
	assert.Len(t, result.RecentResults, 1)

	repo.AssertNotCalled(t, "GetSummary")
	cache.AssertExpectations(t)
}

func TestGetSummaryCacheMiss(t *testing.T) {
	repo := new(mockDashboardRepository)
	cache := new(mockDashboardCache)
	service := NewService(repo, cache)

	userID := "user-id"

	expected := &DashboardResponse{
		TotalTests:        3,
		TotalAttempts:     5,
		CompletedAttempts: 2,
		AverageScore:      80,
		RecentResults:     []RecentResult{},
	}

	cache.On("Get", mock.Anything, "dashboard:user:"+userID).
		Return("", redis.Nil)

	repo.On("GetSummary", mock.Anything, userID).
		Return(expected, nil)

	cache.On("Set", mock.Anything, "dashboard:user:"+userID, mock.Anything, 5*time.Minute).
		Return(nil)

	result, err := service.GetSummary(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected.TotalTests, result.TotalTests)
	assert.Equal(t, expected.TotalAttempts, result.TotalAttempts)
	assert.Equal(t, expected.CompletedAttempts, result.CompletedAttempts)
	assert.Equal(t, expected.AverageScore, result.AverageScore)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestGetSummaryInvalidCacheJSON(t *testing.T) {
	repo := new(mockDashboardRepository)
	cache := new(mockDashboardCache)
	service := NewService(repo, cache)

	userID := "user-id"

	expected := &DashboardResponse{
		TotalTests:        1,
		TotalAttempts:     2,
		CompletedAttempts: 1,
		AverageScore:      100,
		RecentResults:     []RecentResult{},
	}

	cache.On("Get", mock.Anything, "dashboard:user:"+userID).
		Return("invalid-json", nil)

	repo.On("GetSummary", mock.Anything, userID).
		Return(expected, nil)

	cache.On("Set", mock.Anything, "dashboard:user:"+userID, mock.Anything, 5*time.Minute).
		Return(nil)

	result, err := service.GetSummary(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected.TotalTests, result.TotalTests)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestGetSummaryRepositoryError(t *testing.T) {
	repo := new(mockDashboardRepository)
	cache := new(mockDashboardCache)
	service := NewService(repo, cache)

	userID := "user-id"

	cache.On("Get", mock.Anything, "dashboard:user:"+userID).
		Return("", redis.Nil)

	repo.On("GetSummary", mock.Anything, userID).
		Return((*DashboardResponse)(nil), errors.New("database error"))

	result, err := service.GetSummary(context.Background(), userID)

	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}
