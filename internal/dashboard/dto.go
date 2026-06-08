package dashboard

type DashboardResponse struct {
	TotalTests        int64          `json:"total_tests"`
	TotalAttempts     int64          `json:"total_attempts"`
	CompletedAttempts int64          `json:"completed_attempts"`
	AverageScore      float64        `json:"average_score"`
	RecentResults     []RecentResult `json:"recent_results"`
}

type RecentResult struct {
	ResultID   string  `json:"result_id"`
	TestID     string  `json:"test_id"`
	AttemptID  string  `json:"attempt_id"`
	Score      int     `json:"score"`
	MaxScore   int     `json:"max_score"`
	Percentage float64 `json:"percentage"`
	Status     string  `json:"status"`
}
