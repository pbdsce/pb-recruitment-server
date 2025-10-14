package dto

type ProblemOverview struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
	Type   string `json:"type"`   // "code" or "mcq"
	Status string `json:"status"` // "solved", "attempted", "not_attempted"
}

type GetProblemStatementResponse struct {
	ProblemID   string `json:"problem_id"`
	ContestID   string `json:"contest_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Score       int    `json:"score"`
	Type        string `json:"type"` // "code" or "mcq"
	// For code problems (future enhancement)
	// TimeLimit   int `json:"time_limit"`   // milliseconds
	// MemoryLimit int `json:"memory_limit"` // MB
	// For MCQ problems (future enhancement)
	// Options  []string `json:"options"`  // Available options
	// Multiple bool     `json:"multiple"` // Single or multiple choice
}
