package services

import (
	"app/internal/common"
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/stores"
	"context"
	"fmt"
	"log"
	"strconv"
)

type ContestService struct {
	stores *stores.Storage
}

func NewContestService(stores *stores.Storage) *ContestService {
	return &ContestService{stores: stores}
}

func (cs *ContestService) ModifyRegistration(ctx context.Context, contestID string, userID string, action dto.RegisterationAction) error {
	contest, err := cs.stores.Contests.GetContest(ctx, contestID)
	if err != nil {
		return err
	}

	if contest.GetRegistrationStatus() != models.ContestRegistrationOpen {
		log.Printf("contest %s is not open for registration", contestID)
		return common.ContestRegistrationClosedError
	}

	switch action {
	case dto.RegisterAction:
		user, err := cs.stores.Users.GetUserProfile(ctx, userID)
		if err != nil {
			return err
		}

		if contest.EligibleTo != strconv.Itoa(user.CurrentYear) {
			log.Printf("user %s is not eligible to contest %s", userID, contestID)
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
