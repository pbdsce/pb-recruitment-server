package dto

type UpdateLeaderboardUserRequest struct {
	Hidden       *bool `json:"hidden"`
	Disqualified *bool `json:"disqualified"`
}
