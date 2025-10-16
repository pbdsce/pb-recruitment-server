package stores

import (
	"app/internal/models"
	"app/internal/models/dto"
	"context"
	"database/sql"

	"firebase.google.com/go/v4/auth"
)

type Storage struct {
	// Declarations of method extensions for each store go here
	Contests interface {
		ListContests(context.Context, int) ([]models.Contest, error)
	}
	Users interface {
		CreateUser(context.Context, *auth.UserRecord, *dto.CreateUserRequest) error
		GetUserProfile(context.Context, string) (*models.User, error)
		UpdateUserProfile(context.Context, string, *dto.UpdateUserProfileRequest) error
	}
	Submissions interface {
		GetSubmissionStatusByID(context.Context, string) (*models.Submission, error)
		GetSubmissionDetailsByID(context.Context, string) (*models.Submission, error)
	}
	Rankings interface {
		// todo: add ranking store
	}
	Problems interface {
		// todo: add problem store
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Contests:    NewContestStore(db),
		Users:       NewUserStore(db),
		Submissions: NewSubmissionStore(db),
		Rankings:    NewRankingStore(db),
		Problems:    NewProblemStore(db),
	}
}
