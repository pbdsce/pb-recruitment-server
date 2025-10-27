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
