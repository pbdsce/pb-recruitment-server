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
		IsRegistered(context.Context, string, string) (bool, error)
		GetContest(context.Context, string) (*dto.GetContestResponse, error)
		RegisterUser(context.Context, string, string) error
		UnregisterUser(context.Context, string, string) error
	}
	Users interface {
		CreateUser(context.Context, *auth.UserRecord, *dto.CreateUserRequest) error
		GetUserProfile(context.Context, string) (*models.User, error)
		UpdateUserProfile(context.Context, string, *dto.UpdateUserProfileRequest) error
	}
	Submissions interface {
		GetSubmissionStatusByID(context.Context, string) (*models.Submission, error)
		GetSubmissionDetailsByID(context.Context, string) (*models.Submission, error)
		GetTestCaseResultsBySubmissionID(context.Context, string) ([]models.TestCaseResult, error)
		ListUserSubmissionsByProblemID(context.Context, string, string, int) ([]models.Submission, error)
	}
	Rankings interface {
		// todo: add ranking store
	}
	Problems interface {
		GetProblemList(context.Context, string) ([]dto.ProblemOverview, error)
		GetProblem(context.Context, string, string) (*dto.GetProblemStatementResponse, error)
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
