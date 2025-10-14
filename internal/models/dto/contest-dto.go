package dto

import "app/internal/models"

// GetContestResponse represents the response for getting contest details
type GetContestResponse struct {
	models.Contest
	IsRegistered bool `json:"is_registered,omitempty"` // Only included when user is authenticated
}

// ModifyRegistrationRequest represents the request for registering/unregistering
type ModifyRegistrationRequest struct {
	Action string `json:"action" validate:"required,oneof=register unregister"`
}

// GetContestProblemsResponse represents the response for getting contest problems list
type GetContestProblemsResponse struct {
	ContestID string            `json:"contest_id"`
	Problems  []ProblemOverview `json:"problems"`
}

type ProblemOverview struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
	Type   string `json:"type"`   // "code" or "mcq"
	Status string `json:"status"` // "solved", "attempted", "not_attempted"
}

// GetProblemStatementResponse represents the response for getting problem statement
type GetProblemStatementResponse struct {
	ProblemID   string `json:"problem_id"`
	ContestID   string `json:"contest_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Score       int    `json:"score"`
	Type        string `json:"type"` // "code" or "mcq"
	// // For code problems
	// TimeLimit   int `json:"time_limit"`   // milliseconds
	// MemoryLimit int `json:"memory_limit"` // MB
	// // For MCQ problems
	// Options  []string `json:"options"`  // Available options
	// Multiple bool     `json:"multiple"` // Single or multiple choice
}
