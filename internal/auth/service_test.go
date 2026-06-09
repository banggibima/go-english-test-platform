package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/banggibima/go-english-test-platform/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)

	user, _ := args.Get(0).(*User)
	return user, args.Error(1)
}

func (m *mockAuthRepository) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	args := m.Called(ctx, name)

	role, _ := args.Get(0).(*Role)
	return role, args.Error(1)
}

func (m *mockAuthRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)

	user.ID = "user-id"
	user.RoleName = "user"
	user.IsActive = true

	return args.Error(0)
}

func (m *mockAuthRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockAuthRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	args := m.Called(ctx, token)

	token.ID = "refresh-token-id"
	token.CreatedAt = time.Now()

	return args.Error(0)
}

func (m *mockAuthRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	args := m.Called(ctx, token)

	refreshToken, _ := args.Get(0).(*RefreshToken)
	return refreshToken, args.Error(1)
}

func (m *mockAuthRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockAuthRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)

	user, _ := args.Get(0).(*User)
	return user, args.Error(1)
}

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret",
	}
}

func TestRegisterSuccess(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := RegisterRequest{
		FullName: "Bima",
		Email:    "bima@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return((*User)(nil), nil)

	repo.On("GetRoleByName", mock.Anything, "user").
		Return(&Role{
			ID:   "role-id",
			Name: "user",
		}, nil)

	repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*auth.User")).
		Return(nil)

	repo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).
		Return(nil)

	result, err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.User.Email)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	repo.AssertExpectations(t)
}

func TestRegisterEmailAlreadyExists(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := RegisterRequest{
		FullName: "Bima",
		Email:    "bima@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(&User{
			ID:    "user-id",
			Email: req.Email,
		}, nil)

	result, err := service.Register(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "email already registered")

	repo.AssertExpectations(t)
}

func TestRegisterDefaultRoleNotFound(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := RegisterRequest{
		FullName: "Bima",
		Email:    "bima@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return((*User)(nil), nil)

	repo.On("GetRoleByName", mock.Anything, "user").
		Return((*Role)(nil), nil)

	result, err := service.Register(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "default role not found")

	repo.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	passwordHash, _ := hashPassword("password123")

	req := LoginRequest{
		Email:    "bima@example.com",
		Password: "password123",
	}

	user := &User{
		ID:           "user-id",
		RoleID:       "role-id",
		RoleName:     "user",
		FullName:     "Bima",
		Email:        req.Email,
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(user, nil)

	repo.On("UpdateLastLogin", mock.Anything, user.ID).
		Return(nil)

	repo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).
		Return(nil)

	result, err := service.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.User.Email)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	repo.AssertExpectations(t)
}

func TestLoginInvalidEmail(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := LoginRequest{
		Email:    "bima@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return((*User)(nil), nil)

	result, err := service.Login(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid email or password")

	repo.AssertExpectations(t)
}

func TestLoginInvalidPassword(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	passwordHash, _ := hashPassword("password123")

	req := LoginRequest{
		Email:    "bima@example.com",
		Password: "wrong-password",
	}

	user := &User{
		ID:           "user-id",
		RoleID:       "role-id",
		RoleName:     "user",
		Email:        req.Email,
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(user, nil)

	result, err := service.Login(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid email or password")

	repo.AssertExpectations(t)
}

func TestLoginInactiveUser(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	passwordHash, _ := hashPassword("password123")

	req := LoginRequest{
		Email:    "bima@example.com",
		Password: "password123",
	}

	user := &User{
		ID:           "user-id",
		Email:        req.Email,
		PasswordHash: passwordHash,
		IsActive:     false,
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(user, nil)

	result, err := service.Login(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "user is inactive")

	repo.AssertExpectations(t)
}

func TestRefreshTokenSuccess(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := RefreshTokenRequest{
		RefreshToken: "refresh-token",
	}

	refreshToken := &RefreshToken{
		ID:        "refresh-token-id",
		UserID:    "user-id",
		Token:     req.RefreshToken,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	user := &User{
		ID:       "user-id",
		RoleID:   "role-id",
		RoleName: "user",
		FullName: "Bima",
		Email:    "bima@example.com",
		IsActive: true,
	}

	repo.On("GetRefreshToken", mock.Anything, req.RefreshToken).
		Return(refreshToken, nil)

	repo.On("GetUserByID", mock.Anything, refreshToken.UserID).
		Return(user, nil)

	repo.On("RevokeRefreshToken", mock.Anything, req.RefreshToken).
		Return(nil)

	repo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).
		Return(nil)

	result, err := service.RefreshToken(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Email, result.User.Email)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	repo.AssertExpectations(t)
}

func TestRefreshTokenInvalid(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := RefreshTokenRequest{
		RefreshToken: "invalid-token",
	}

	repo.On("GetRefreshToken", mock.Anything, req.RefreshToken).
		Return((*RefreshToken)(nil), nil)

	result, err := service.RefreshToken(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid refresh token")

	repo.AssertExpectations(t)
}

func TestRefreshTokenExpired(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := RefreshTokenRequest{
		RefreshToken: "expired-token",
	}

	refreshToken := &RefreshToken{
		ID:        "refresh-token-id",
		UserID:    "user-id",
		Token:     req.RefreshToken,
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	repo.On("GetRefreshToken", mock.Anything, req.RefreshToken).
		Return(refreshToken, nil)

	result, err := service.RefreshToken(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "refresh token expired")

	repo.AssertExpectations(t)
}

func TestRepositoryError(t *testing.T) {
	repo := new(mockAuthRepository)
	service := NewService(repo, testConfig())

	req := LoginRequest{
		Email:    "bima@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", mock.Anything, req.Email).
		Return((*User)(nil), errors.New("database error"))

	result, err := service.Login(context.Background(), req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")

	repo.AssertExpectations(t)
}
