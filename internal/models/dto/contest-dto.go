package dto

import "time"

type CreateContestRequest struct {
	Name                  string `json:"name" validate:"required"`
	RegistrationStartTime int64  `json:"registration_start_time" validate:"required,min=0"`
	RegistrationEndTime   int64  `json:"registration_end_time" validate:"required,min=0"`
	StartTime             int64  `json:"start_time" validate:"required,min=0"`
	EndTime               int64  `json:"end_time" validate:"required,min=0"`
}

type UpdateContestRequest struct {
	Name                  *string `json:"name,omitempty" validate:"omitempty"`
	RegistrationStartTime *int64  `json:"registration_start_time,omitempty" validate:"omitempty,min=0"`
	RegistrationEndTime   *int64  `json:"registration_end_time,omitempty" validate:"omitempty,min=0"`
	StartTime             *int64  `json:"start_time,omitempty" validate:"omitempty,min=0"`
	EndTime               *int64  `json:"end_time,omitempty" validate:"omitempty,min=0"`
}

type ContestResponse struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	RegistrationStartTime time.Time `json:"registration_start_time"`
	RegistrationEndTime   time.Time `json:"registration_end_time"`
	StartTime             time.Time `json:"start_time"`
	EndTime               time.Time `json:"end_time"`
	Status                string    `json:"status"` // "upcoming", "active", "ended", "registration_open", etc.
}

type ListContestsResponse struct {
	Contests   []ContestResponse `json:"contests"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	HasMore    bool              `json:"has_more"`
}

type RegistrationRequest struct {
	Action string `json:"action" validate:"required,oneof=register unregister"`
}

type RegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ContestProblemRow struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
	Type  string `json:"type"` // MCQ or Code
}

type ContestProblemsResponse struct {
	ContestID string              `json:"contest_id"`
	Problems  []ContestProblemRow `json:"problems"`
}

type LeaderboardRow struct {
	Rank     int    `json:"rank"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	USN      string `json:"usn"`
	Score    int    `json:"score"`
	Solved   int    `json:"solved"` // problems solved
}

type ContestLeaderboardResponse struct {
	ContestID   string           `json:"contest_id"`
	ContestName string           `json:"contest_name"`
	Entries     []LeaderboardRow `json:"entries"`
	Page        int              `json:"page"`
	PageSize    int              `json:"page_size"`
	TotalCount  int              `json:"total_count"`
	HasMore     bool             `json:"has_more"`
}
