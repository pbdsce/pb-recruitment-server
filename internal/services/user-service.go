package services

import (
	"app/internal/common"
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/stores"
	"context"
	"database/sql"
	"fmt"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

type UserService struct {
	stores     *stores.Storage
	authClient *auth.Client
}

func NewUserService(stores *stores.Storage, authClient *auth.Client) *UserService {
	return &UserService{stores, authClient}
}

func (us *UserService) CreateUser(ctx context.Context, userID string, req *dto.CreateUserRequest) error {
	user, err := us.authClient.GetUser(ctx, userID)
	if err != nil {
		log.Printf("user-service: error fetching user: %v", err)
		return fmt.Errorf("error fetching user: %v", err)
	}

	return us.stores.Users.CreateUser(ctx, user, req)
}

func (us *UserService) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	return us.stores.Users.GetUserProfile(ctx, userID)
}

func (us *UserService) UpdateUserProfile(ctx context.Context, userID string, req *dto.UpdateUserProfileRequest) error {
	return us.stores.Users.UpdateUserProfile(ctx, userID, req)
}

func (us *UserService) Signup(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error) {
	if us.stores == nil || us.stores.DB == nil {
		return nil, fmt.Errorf("user service: storage is not initialized")
	}

	tx, err := us.stores.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("user service: failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
				log.Printf("user-service: rollback failed: %v", rbErr)
			}
		}
	}()

	tempID := uuid.NewString()

	if err := us.stores.Users.CreatePendingUserTx(ctx, tx, tempID, req); err != nil {
		return nil, err
	}

	userRecord, err := us.authClient.CreateUser(ctx, (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		DisplayName(req.Name))
	if err != nil {
		if auth.IsEmailAlreadyExists(err) || auth.IsUIDAlreadyExists(err) {
			return nil, common.UserAlreadyExistsError
		}
		return nil, fmt.Errorf("firebase create user: %w", err)
	}

	if err := us.stores.Users.UpdateUserIDTx(ctx, tx, tempID, userRecord.UID); err != nil {
		if delErr := us.authClient.DeleteUser(ctx, userRecord.UID); delErr != nil {
			log.Printf("user-service: cleanup firebase user %s failed: %v", userRecord.UID, delErr)
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if delErr := us.authClient.DeleteUser(ctx, userRecord.UID); delErr != nil {
			log.Printf("user-service: cleanup firebase user %s failed after commit error: %v", userRecord.UID, delErr)
		}
		return nil, fmt.Errorf("user service: commit failed: %w", err)
	}
	committed = true

	customToken, err := us.authClient.CustomToken(ctx, userRecord.UID)
	if err != nil {
		return nil, fmt.Errorf("firebase custom token: %w", err)
	}

	return &dto.SignupResponse{
		UserID:      userRecord.UID,
		CustomToken: customToken,
	}, nil
}
