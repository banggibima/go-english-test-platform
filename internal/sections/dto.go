package sections

type CreateSectionRequest struct {
	TestID          string `json:"test_id" validate:"required"`
	Title           string `json:"title" validate:"required,min=3,max=255"`
	SectionCode     string `json:"section_code" validate:"required,min=3,max=50"`
	DisplayOrder    int    `json:"display_order" validate:"required,min=1"`
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=1"`
}

type UpdateSectionRequest struct {
	Title           string `json:"title" validate:"required,min=3,max=255"`
	DisplayOrder    int    `json:"display_order" validate:"required,min=1"`
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=1"`
	IsActive        bool   `json:"is_active"`
}

type SectionResponse struct {
	ID              string `json:"id"`
	TestID          string `json:"test_id"`
	Title           string `json:"title"`
	SectionCode     string `json:"section_code"`
	DisplayOrder    int    `json:"display_order"`
	DurationMinutes int    `json:"duration_minutes"`
	TotalQuestions  int    `json:"total_questions"`
	IsActive        bool   `json:"is_active"`
}
