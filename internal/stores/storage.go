package stores

import (
	"app/internal/models"
	"context"
	"database/sql"
)

type Storage struct {
	// Declarations of method extensions for each store go here
	Contests interface {
		ListContests(context.Context, int) ([]models.Contest, error)
	}
	Users interface {
		// todo: add user store
	}
	Submissions interface {
		GetSubmissionStatusByID(context.Context, string) (models.Submission, error)
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
