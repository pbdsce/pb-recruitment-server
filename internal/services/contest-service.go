package services

import (
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/stores"
	"context"

	"github.com/google/uuid"
)

type ContestService struct {
	stores *stores.Storage
}

func NewContestService(stores *stores.Storage) *ContestService {
	return &ContestService{stores: stores}
}

func (cs *ContestService) CreateContest(ctx context.Context, contest *models.Contest) (*models.Contest, error) {
	contest.ID = uuid.NewString()
	if err := cs.stores.Contests.CreateContest(ctx, contest); err != nil {
		return nil, err
	}
	return contest, nil
}

func (cs *ContestService) UpdateContest(ctx context.Context, contest *models.Contest) (*models.Contest, error) {
	if err := cs.stores.Contests.UpdateContest(ctx, contest); err != nil {
		return nil, err
	}
	return contest, nil
}

func (cs *ContestService) DeleteContest(ctx context.Context, contestID string) error {
	return cs.stores.Contests.DeleteContest(ctx, contestID)
}

func (cs *ContestService) RegisterParticipant(contestID string, userID string) error {
	// Registration logic would go here
	return nil
}

func (cs *ContestService) ListContests(ctx context.Context, page int) ([]models.Contest, error) {
	return cs.stores.Contests.ListContests(ctx, page)
}

//Problem Reated Services

func (cs *ContestService) CreateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error) {

	problem.ID = uuid.NewString()

	if err := cs.stores.Problems.CreateProblem(ctx, problem); err != nil {
		return nil, err
	}

	return problem, nil
}

func (cs *ContestService) UpdateProblem(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	if err := cs.stores.Problems.UpdateProblem(ctx, problem); err != nil {
		return nil, err
	}
	return problem, nil
}

func (cs *ContestService) DeleteProblem(ctx context.Context, contestID string, problemID string) error {
	return cs.stores.Problems.DeleteProblem(ctx, contestID, problemID)
}

//Leaderboard related services

func (cs *ContestService) UpdateLeaderboardUser(ctx context.Context, contestID string, userID string, req *dto.UpdateLeaderboardUserRequest) error {
	return cs.stores.Rankings.UpdateLeaderboardUser(ctx, contestID, userID, req)
}
