package models

type Contest struct {
	ID                    string `json:"id"` // UUID as string
	Name                  string `json:"name"`
	RegistrationStartTime int64  `json:"registration_start_time"` // Unix timestamp
	RegistrationEndTime   int64  `json:"registration_end_time"`   // Unix timestamp
	StartTime             int64  `json:"start_time"`              // Unix timestamp
	EndTime               int64  `json:"end_time"`                // Unix timestamp
	Status                string `json:"status"`                  // "upcoming", "registration_open", "registration_closed", "active", "ended"
}
