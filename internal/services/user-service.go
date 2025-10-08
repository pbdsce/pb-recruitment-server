package services

import (
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/stores"
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/v4/auth"
)

type UserService struct {
	stores        *stores.Storage
	authClient *auth.Client
}

func NewUserService(stores *stores.Storage, authClient *auth.Client) *UserService {
	return &UserService{stores, authClient}
}

func (us *UserService) CreateUser(ctx context.Context, userID string, req *dto.CreateUserRequest) (bool, error) {

	user, err := us.authClient.GetUser(ctx, userID)
	if err != nil {
		log.Printf("user-service: error fetching user: %v", err)
		return false, fmt.Errorf("error fetching user: %v", err)
	}

	return us.stores.Users.CreateUser(ctx, user, req)
}

func (us *UserService) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	return us.stores.Users.GetUserProfile(ctx, userID)
}

func (us *UserService) UpdateUserProfile(ctx context.Context, userID string, req *dto.UpdateUserProfileRequest) (bool, error) {
	return us.stores.Users.UpdateUserProfile(ctx, userID, req)
}
