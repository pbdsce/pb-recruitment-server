package models

type Problem struct {
	ID          string         `json:"id"`         // UUID as string
	ContestID   string         `json:"contest_id"` // Foreign key reference to Contest
	Name        string         `json:"name"`
	Description string         `json:"description"` // Problem statement/description
	Score       int            `json:"score"`
	Type        SubmissionType `json:"type"`             // "code" or "mcq"
	Answer      []int          `json:"answer,omitempty"` // Only for MCQ; supports multi-choice
}
