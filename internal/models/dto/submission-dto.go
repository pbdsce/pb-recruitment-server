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

type ListProblemSubmissionsResponse struct {
	Submissions []models.Submission `json:"submissions"`
}

type TestCaseResultDTO struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Duration int64  `json:"duration_ms"`
}

type SubmissionDetailsResponse struct {
	SubmissionID string                  `json:"submission_id"`
	Status       models.SubmissionStatus `json:"status"`
	RuntimeMs    int64                   `json:"runtime_ms"`
	MemoryKB     int64                   `json:"memory_kb"`
	TestCases    []TestCaseResultDTO     `json:"test_cases"`
}
