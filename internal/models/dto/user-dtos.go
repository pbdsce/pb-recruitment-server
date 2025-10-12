package dto

type CreateUserRequest struct {
	Name         string `json:"name" validate:"required"`
	USN          string `json:"usn" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	CurrentYear  int    `json:"current_year" validate:"required,min=2000,max=2100"`
	Department   string `json:"department" validate:"required"`
}

type UpdateUserProfileRequest struct {
	Name         string `json:"name" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	Department   string `json:"department" validate:"required"`
}
