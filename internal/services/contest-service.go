package services

import (
	"app/internal/common"
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/stores"
	"context"
)

type ContestService struct {
	stores *stores.Storage
}

func NewContestService(stores *stores.Storage) *ContestService {
	return &ContestService{stores: stores}
}

func (cs *ContestService) RegisterParticipant(contestID string, userID string) error {
	// Registration logic would go here
	return nil
}

func (cs *ContestService) ListContests(ctx context.Context, page int) ([]models.Contest, error) {
	return cs.stores.Contests.ListContests(ctx, page)
}

func (cs *ContestService) GetProblemVisibility(ctx context.Context, contestID string, userID string) error {

	contest, err := cs.GetContest(ctx, contestID, userID)
	if err != nil {
		return err
	}

	if contest.IsRegistered == nil || !*contest.IsRegistered {
		return common.UserNotRegisteredError
	}

	if contest.GetRunningStatus() == models.ContestRunningUpcoming {
		return common.ContestNotRunningError
	}

	return nil
}

func (cs *ContestService) GetContestProblemsList(ctx context.Context, contestID string) ([]dto.ProblemOverview, error) {
	return cs.stores.Problems.GetProblemList(ctx, contestID)
}

func (cs *ContestService) GetContestProblem(ctx context.Context, contestID string, problemID string) (*dto.GetProblemStatementResponse, error) {
	return cs.stores.Problems.GetProblem(ctx, problemID, contestID)
}

func (cs *ContestService) GetContest(ctx context.Context, contestID string, userID string) (*dto.GetContestResponse, error) {
	contest_response, err := cs.stores.Contests.GetContest(ctx, contestID)
	if err != nil {
		return nil, err
	}

	if userID == "" {
		return contest_response, nil
	}

	r, err := cs.stores.Contests.IsRegistered(ctx, contestID, userID)
	if err != nil {
		return nil, err
	}

	contest_response.IsRegistered = &r
	return contest_response, nil
}
