package dto

type UpdateUserProfileRequest struct {
	Name         string `json:"name" validate:"required"`
	USN          string `json:"usn" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	JoiningYear  int    `json:"joining_year" validate:"required,min=2000,max=2100"`
	Department   string `json:"department" validate:"required"`
}
