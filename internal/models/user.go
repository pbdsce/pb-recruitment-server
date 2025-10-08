package models

type User struct {
	ID           string `json:"-"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	USN          string `json:"usn"`
	MobileNumber string `json:"mobile_number"`
	CurrentYear  int    `json:"current_year"`
	Department   string `json:"department"`
}
