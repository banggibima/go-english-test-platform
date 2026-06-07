package users

type UserResponse struct {
	ID         string  `json:"id"`
	RoleID     string  `json:"role_id"`
	FullName   string  `json:"full_name"`
	Email      string  `json:"email"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	IsActive   bool    `json:"is_active"`
	IsVerified bool    `json:"is_verified"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=255"`
}
