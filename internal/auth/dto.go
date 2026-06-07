package auth

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type UserResponse struct {
	ID         string  `json:"id"`
	RoleID     string  `json:"role_id"`
	FullName   string  `json:"full_name"`
	Email      string  `json:"email"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	IsActive   bool    `json:"is_active"`
	IsVerified bool    `json:"is_verified"`
}
