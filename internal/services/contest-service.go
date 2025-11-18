package services

import (
	"app/internal/common"
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/stores"
	"context"
	"slices"

	"fmt"

	"github.com/google/uuid"

	"github.com/labstack/gommon/log"
)

type ContestService struct {
	stores *stores.Storage
}

func NewContestService(stores *stores.Storage) *ContestService {
	return &ContestService{stores: stores}
}

func (cs *ContestService) CreateContest(ctx context.Context, contest *models.Contest) (*models.Contest, error) {
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

func (cs *ContestService) ModifyRegistration(ctx context.Context, contestID string, userID string, action dto.RegisterationAction) error {
	contest, err := cs.stores.Contests.GetContest(ctx, contestID)
	if err != nil {
		return err
	}

	if contest.GetRegistrationStatus() != models.ContestRegistrationOpen {
		log.Errorf("contest %s is not open for registration", contestID)
		return common.ContestRegistrationClosedError
	}

	switch action {
	case dto.RegisterAction:
		user, err := cs.stores.Users.GetUserProfile(ctx, userID)
		if err != nil {
			log.Errorf("failed to get user profile for user %s: %v", userID, err)
			return err
		}

		if !slices.Contains(contest.EligibleTo, user.CurrentYear) {
			log.Errorf("user %s is not eligible to contest %s", userID, contestID)
			return common.InvalidYearError
		}

		return cs.stores.Contests.RegisterUser(ctx, contestID, userID)

	case dto.UnregisterAction:
		return cs.stores.Contests.UnregisterUser(ctx, contestID, userID)

	default:
		return fmt.Errorf("invalid action: %s", action)
	}
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

func (cs *ContestService) GetContestRegistrations(ctx context.Context, contestID string) ([]dto.ContestRegistration, error) {
	return cs.stores.Contests.GetContestRegistrations(ctx, contestID)
}
