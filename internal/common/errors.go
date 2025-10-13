package common

import "errors"

var (
	ErrNotFound                   = errors.New("Resource not found")
	UserAlreadyExistsError        = errors.New("user already exists")
	UserNotFoundError             = errors.New("user not found")
	CreateUserFailedError         = errors.New("failed to create user")
	FetchUserFailedError          = errors.New("failed to fetch user profile")
	UpdateUserFailedError         = errors.New("failed to update user profile")
	ProfileUpdateFailedError      = errors.New("user profile update failed")
	InvalidApplicationNumberError = errors.New("first-year students must provide a valid Application Number (e.g., 25UGDS1234)")
	InvalidUSNError               = errors.New("provide a valid USN (e.g., 1DS24IC015)")
	InvalidMobileNumberError      = errors.New("invalid mobile number format")
	USNYearImmutableError         = errors.New("USN and Year cannot be changed")
)
