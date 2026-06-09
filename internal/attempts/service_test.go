package attempts

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAttemptRepository struct {
	mock.Mock
}

func (m *mockAttemptRepository) Create(ctx context.Context, attempt *Attempt) error {
	args := m.Called(ctx, attempt)

	attempt.ID = "attempt-id"
	attempt.Status = "in_progress"
	attempt.StartedAt = time.Now()

	return args.Error(0)
}

func (m *mockAttemptRepository) FindAllByUserID(ctx context.Context, userID string) ([]Attempt, error) {
	args := m.Called(ctx, userID)

	items, _ := args.Get(0).([]Attempt)
	return items, args.Error(1)
}

func (m *mockAttemptRepository) FindByID(ctx context.Context, id string) (*Attempt, error) {
	args := m.Called(ctx, id)

	attempt, _ := args.Get(0).(*Attempt)
	return attempt, args.Error(1)
}

func (m *mockAttemptRepository) Submit(ctx context.Context, id string, userID string) (*Attempt, error) {
	args := m.Called(ctx, id, userID)

	attempt, _ := args.Get(0).(*Attempt)
	return attempt, args.Error(1)
}

type mockQueue struct {
	mock.Mock
}

func (m *mockQueue) Publish(queueName string, body []byte) error {
	args := m.Called(queueName, body)
	return args.Error(0)
}

func TestStartAttemptSuccess(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	req := StartAttemptRequest{
		TestID: "test-id",
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*attempts.Attempt")).
		Return(nil)

	result, err := service.StartAttempt(context.Background(), "user-id", req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "attempt-id", result.ID)
	assert.Equal(t, "user-id", result.UserID)
	assert.Equal(t, req.TestID, result.TestID)
	assert.Equal(t, "in_progress", result.Status)

	repo.AssertExpectations(t)
}

func TestStartAttemptRepositoryError(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	req := StartAttemptRequest{
		TestID: "test-id",
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*attempts.Attempt")).
		Return(errors.New("database error"))

	result, err := service.StartAttempt(context.Background(), "user-id", req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")

	repo.AssertExpectations(t)
}

func TestFindAllByUserIDSuccess(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	items := []Attempt{
		{
			ID:        "attempt-id-1",
			UserID:    "user-id",
			TestID:    "test-id",
			Status:    "in_progress",
			StartedAt: time.Now(),
		},
		{
			ID:        "attempt-id-2",
			UserID:    "user-id",
			TestID:    "test-id",
			Status:    "completed",
			StartedAt: time.Now(),
		},
	}

	repo.On("FindAllByUserID", mock.Anything, "user-id").
		Return(items, nil)

	result, err := service.FindAllByUserID(context.Background(), "user-id")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "attempt-id-1", result[0].ID)
	assert.Equal(t, "attempt-id-2", result[1].ID)

	repo.AssertExpectations(t)
}

func TestFindByIDSuccess(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	attempt := &Attempt{
		ID:        "attempt-id",
		UserID:    "user-id",
		TestID:    "test-id",
		Status:    "in_progress",
		StartedAt: time.Now(),
	}

	repo.On("FindByID", mock.Anything, "attempt-id").
		Return(attempt, nil)

	result, err := service.FindByID(context.Background(), "attempt-id", "user-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, attempt.ID, result.ID)
	assert.Equal(t, attempt.UserID, result.UserID)

	repo.AssertExpectations(t)
}

func TestFindByIDNotFound(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	repo.On("FindByID", mock.Anything, "attempt-id").
		Return((*Attempt)(nil), nil)

	result, err := service.FindByID(context.Background(), "attempt-id", "user-id")

	assert.Nil(t, result)
	assert.EqualError(t, err, "attempt not found")

	repo.AssertExpectations(t)
}

func TestFindByIDForbidden(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	attempt := &Attempt{
		ID:        "attempt-id",
		UserID:    "other-user-id",
		TestID:    "test-id",
		Status:    "in_progress",
		StartedAt: time.Now(),
	}

	repo.On("FindByID", mock.Anything, "attempt-id").
		Return(attempt, nil)

	result, err := service.FindByID(context.Background(), "attempt-id", "user-id")

	assert.Nil(t, result)
	assert.EqualError(t, err, "forbidden")

	repo.AssertExpectations(t)
}

func TestSubmitSuccess(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	now := time.Now()

	attempt := &Attempt{
		ID:          "attempt-id",
		UserID:      "user-id",
		TestID:      "test-id",
		Status:      "submitted",
		StartedAt:   now,
		SubmittedAt: &now,
	}

	repo.On("Submit", mock.Anything, "attempt-id", "user-id").
		Return(attempt, nil)

	queue.On("Publish", "score-attempt", mock.AnythingOfType("[]uint8")).
		Return(nil)

	result, err := service.Submit(context.Background(), "attempt-id", "user-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, attempt.ID, result.ID)
	assert.Equal(t, "submitted", result.Status)
	assert.NotNil(t, result.SubmittedAt)

	repo.AssertExpectations(t)
	queue.AssertExpectations(t)
}

func TestSubmitNotFound(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	repo.On("Submit", mock.Anything, "attempt-id", "user-id").
		Return((*Attempt)(nil), nil)

	result, err := service.Submit(context.Background(), "attempt-id", "user-id")

	assert.Nil(t, result)
	assert.EqualError(t, err, "attempt not found or already submitted")

	repo.AssertExpectations(t)
	queue.AssertNotCalled(t, "Publish")
}

func TestSubmitRepositoryError(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	repo.On("Submit", mock.Anything, "attempt-id", "user-id").
		Return((*Attempt)(nil), errors.New("database error"))

	result, err := service.Submit(context.Background(), "attempt-id", "user-id")

	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")

	repo.AssertExpectations(t)
	queue.AssertNotCalled(t, "Publish")
}

func TestSubmitPublishError(t *testing.T) {
	repo := new(mockAttemptRepository)
	queue := new(mockQueue)
	service := NewService(repo, queue)

	now := time.Now()

	attempt := &Attempt{
		ID:          "attempt-id",
		UserID:      "user-id",
		TestID:      "test-id",
		Status:      "submitted",
		StartedAt:   now,
		SubmittedAt: &now,
	}

	repo.On("Submit", mock.Anything, "attempt-id", "user-id").
		Return(attempt, nil)

	queue.On("Publish", "score-attempt", mock.AnythingOfType("[]uint8")).
		Return(errors.New("rabbitmq error"))

	result, err := service.Submit(context.Background(), "attempt-id", "user-id")

	assert.Nil(t, result)
	assert.EqualError(t, err, "rabbitmq error")

	repo.AssertExpectations(t)
	queue.AssertExpectations(t)
}
