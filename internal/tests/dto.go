package tests

type CreateTestRequest struct {
	Title           string  `json:"title" validate:"required,min=3,max=255"`
	Description     *string `json:"description"`
	TestCode        string  `json:"test_code" validate:"required,min=3,max=100"`
	DurationMinutes int     `json:"duration_minutes" validate:"required,min=1"`
	PassingScore    int     `json:"passing_score" validate:"min=0,max=100"`
}

type UpdateTestRequest struct {
	Title           string  `json:"title" validate:"required,min=3,max=255"`
	Description     *string `json:"description"`
	DurationMinutes int     `json:"duration_minutes" validate:"required,min=1"`
	PassingScore    int     `json:"passing_score" validate:"min=0,max=100"`
	IsActive        bool    `json:"is_active"`
}

type TestResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     *string `json:"description,omitempty"`
	TestCode        string  `json:"test_code"`
	DurationMinutes int     `json:"duration_minutes"`
	TotalQuestions  int     `json:"total_questions"`
	PassingScore    int     `json:"passing_score"`
	IsActive        bool    `json:"is_active"`
	CreatedBy       *string `json:"created_by,omitempty"`
}
