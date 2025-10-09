package stores

import (
	"app/internal/models"
	"context"
	"database/sql"
)

type TestCaseResult struct {
	ID       int
	Status   string
	Duration int64
}

type JudgeResult struct {
	SubmissionID string
	Status       models.SubmissionStatus
	RuntimeMs    int64
	MemoryKB     int64
	TestCases    []TestCaseResult
}

type Submissions interface {
	CreateSubmission(ctx context.Context, submission *models.Submission) (*models.Submission, error)
	GetSubmissionByID(ctx context.Context, id string) (*models.Submission, error)
	ListSubmissionsByProblem(ctx context.Context, userID string, contestID string, problemID string, limit int) ([]models.Submission, error)
	GetJudgeResultBySubmissionID(ctx context.Context, submissionID string) (*JudgeResult, error)
}

type Storage struct {
	// Declarations of method extensions for each store go here
	Contests interface {
		ListContests(context.Context, int) ([]models.Contest, error)
	}
	Users interface {
		// todo: add user store
	}
	Submissions Submissions
	Rankings    interface {
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

func ProvideSubmissions(s *Storage) Submissions {
	return s.Submissions
}
