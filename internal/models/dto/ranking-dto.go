package dto

type UpdateRankingRequest struct {
	Score        *int  `json:"score,omitempty" validate:"omitempty,min=0"`
	Hidden       *bool `json:"hidden,omitempty"`
	Disqualified *bool `json:"disqualified,omitempty"`
	Shortlisted  *bool `json:"shortlisted,omitempty"`
}

type RankingResponse struct {
	ContestID    string `json:"contest_id"`
	UserID       string `json:"user_id"`
	Score        int    `json:"score"`
	Hidden       bool   `json:"hidden"`
	Disqualified bool   `json:"disqualified"`
	Shortlisted  bool   `json:"shortlisted"`
}

type ListRankingResponse struct {
	ContestID  string            `json:"contest_id"`
	Rankings   []RankingResponse `json:"rankings"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalCount int               `json:"total_count"`
	HasMore    bool              `json:"has_more"`
}
