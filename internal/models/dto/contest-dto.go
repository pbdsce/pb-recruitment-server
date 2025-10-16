package dto

import "app/internal/models"

// GetContestResponse represents the response for getting contest details
type GetContestResponse struct {
	models.Contest
	IsRegistered bool `json:"is_registered"` // Only included when user is authenticated
}
