package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/banggibima/go-english-test-platform/config"
	jwtpkg "github.com/banggibima/go-english-test-platform/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateLastLogin(ctx context.Context, userID string) error
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type Service struct {
	repository RepositoryInterface
	cfg        *config.Config
}

func NewService(repository RepositoryInterface, cfg *config.Config) *Service {
	return &Service{
		repository: repository,
		cfg:        cfg,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	existingUser, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	role, err := s.repository.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("default role not found")
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		RoleID:       role.ID,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	if err := comparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := s.repository.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	refreshToken, err := s.repository.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if refreshToken == nil {
		return nil, errors.New("invalid refresh token")
	}

	if refreshToken.RevokedAt != nil {
		return nil, errors.New("refresh token revoked")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.repository.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := s.repository.RevokeRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *Service) generateAuthResponse(ctx context.Context, user *User) (*AuthResponse, error) {
	accessToken, err := jwtpkg.GenerateToken(
		user.ID,
		user.RoleName,
		s.cfg.JWTSecret,
		15*time.Minute,
	)
	if err != nil {
		return nil, err
	}

	refreshTokenValue, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	refreshToken := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenValue,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.repository.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &AuthResponse{
		User: UserResponse{
			ID:         user.ID,
			RoleID:     user.RoleID,
			FullName:   user.FullName,
			Email:      user.Email,
			AvatarURL:  user.AvatarURL,
			IsActive:   user.IsActive,
			IsVerified: user.IsVerified,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshTokenValue,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
