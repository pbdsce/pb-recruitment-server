package models

import "time"

type Contest struct {
	ID                    string `json:"id"` // UUID as string
	Name                  string `json:"name"`
	RegistrationStartTime int64  `json:"registration_start_time"` // Unix timestamp
	RegistrationEndTime   int64  `json:"registration_end_time"`   // Unix timestamp
	StartTime             int64  `json:"start_time"`              // Unix timestamp
	EndTime               int64  `json:"end_time"`                // Unix timestamp
	EligibleTo            string `json:"eligible_to"`             // Student year restriction
}

type ContestRegistrationStatus string

const (
	ContestRegistrationUpcoming ContestRegistrationStatus = "upcoming"
	ContestRegistrationOpen                               = "open"
	ContestRegistrationClosed                             = "closed"
)

type ContestRunningStatus string

const (
	ContestRunningUpcoming ContestRunningStatus = "upcoming"
	ContestRunningOpen                          = "open"
	ContestRunningClosed                        = "closed"
)

func (c *Contest) GetRegistrationStatus() ContestRegistrationStatus {
	now := time.Now().Unix()
	if c.RegistrationStartTime > now {
		return ContestRegistrationUpcoming
	} else if c.RegistrationStartTime <= now && c.RegistrationEndTime >= now {
		return ContestRegistrationOpen
	}
	return ContestRegistrationClosed
}

func (c *Contest) GetRunningStatus() ContestRunningStatus {
	now := time.Now().Unix()
	if c.StartTime > now {
		return ContestRunningUpcoming
	} else if c.StartTime <= now && c.EndTime >= now {
		return ContestRunningOpen
	}
	return ContestRunningClosed
}
