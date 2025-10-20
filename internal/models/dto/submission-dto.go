package dto

import "app/internal/models"

type SubmitSubmissionRequest struct {
	ContestID string `json:"contest_id" validate:"required"`
	ProblemID string `json:"problem_id" validate:"required"`
	Language  string `json:"language"`
	Code      string `json:"code"`   // Base64 encoded code
	Option    []int  `json:"option"` // For MCQ type questions
}

type SubmitSubmissionResponse struct {
	SubmissionID string `json:"submission_id"`
}

type ListProblemSubmissionsRequest struct {
	ProblemID string `query:"problem_id" validate:"required"`
	Page      int    `query:"page" validate:"min=0"`
}

type ListProblemSubmissionsResponse struct {
	Submissions []models.Submission `json:"submissions"`
}
