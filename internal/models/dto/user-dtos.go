package dto

type CreateUserRequest struct {
	Name         string `json:"name" validate:"required"`
	USN          string `json:"usn" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	CurrentYear  int    `json:"current_year" validate:"required,min=1,max=3"`
	Department   string `json:"department" validate:"required"`
}

type UpdateUserProfileRequest struct {
	Name         string `json:"name" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	Department   string `json:"department" validate:"required"`
}

type SignupRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	USN          string `json:"usn" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	CurrentYear  int    `json:"current_year" validate:"required,min=1,max=3"`
	Department   string `json:"department" validate:"required"`
}

type SignupResponse struct {
	UserID string `json:"user_id"`
}
