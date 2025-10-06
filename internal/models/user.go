package models

type User struct {
	ID           string `json:"-"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	USN          string `json:"usn"`
	MobileNumber string `json:"mobile_number"`
	JoiningYear  int    `json:"joining_year"`
	Department   string `json:"department"`
}
