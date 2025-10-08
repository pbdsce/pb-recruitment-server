package dto

type CreateProblemRequest struct {
	ContestID string `json:"contest_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Score     int    `json:"score" validate:"required,min=0"`
	Type      string `json:"type" validate:"required,oneof=MCQ Code"`
	Answer    []int  `json:"answer"` // Only for MCQ
}

type UpdateProblemRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty"`
	Score  *int    `json:"score,omitempty" validate:"omitempty,min=0"`
	Type   *string `json:"type,omitempty" validate:"omitempty,oneof=MCQ Code"`
	Answer []int   `json:"answer"`
}
type ProblemResponse struct {
	ID        string `json:"id"`
	ContestID string `json:"contest_id"`
	Name      string `json:"name"`
	Score     int    `json:"score"`
	Type      string `json:"type"`
}

type ListProblemsResponse struct {
	ContestID string            `json:"contest_id"`
	Problems  []ProblemResponse `json:"problems"`
}
