package common

type UserAlreadyExistsError struct {}

func (e UserAlreadyExistsError) Error() string {
	return "user already exists"
}

type UserNotFoundError struct {}

func (e UserNotFoundError) Error() string {
	return "user not found"
}

