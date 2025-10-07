package models

type Ranking struct {
	ContestID    string `json:"contest_id"` // Primary key
	UserID       string `json:"user_id"`    // Primary key
	Score        int    `json:"score"`
	Hidden       bool   `json:"hidden"`
	Disqualified bool   `json:"disqualified"`
	Shortlisted  bool   `json:"shortlisted"`
}
