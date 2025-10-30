package dto

import "app/internal/models"

type ProblemOverview struct {
	ID    string                `json:"id"`
	Name  string                `json:"name"`
	Score int                   `json:"score"`
	Type  models.SubmissionType `json:"type"`
}

type GetProblemStatementResponse struct {
	ProblemID   string                `json:"problem_id"`
	ContestID   string                `json:"contest_id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Score       int                   `json:"score"`
	Type        models.SubmissionType `json:"type"`
}
