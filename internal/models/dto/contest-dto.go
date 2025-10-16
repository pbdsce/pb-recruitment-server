package dto

import "app/internal/models"

// GetContestResponse represents the response for getting contest details
type GetContestResponse struct {
	models.Contest
	IsRegistered bool `json:"is_registered"`
}

type ModifyRegistrationRequest struct {
	Action string `json:"action" validate:"required,oneof=register unregister"`
}

type GetContestProblemsResponse struct {
	ContestID string            `json:"contest_id"`
	Problems  []ProblemOverview `json:"problems"`
}
