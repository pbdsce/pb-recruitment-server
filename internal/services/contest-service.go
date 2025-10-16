package services

import (
	"app/internal/models"
	"app/internal/stores"
	"context"
	"fmt"
	"time"
)

type ContestService struct {
	stores *stores.Storage
}

func NewContestService(stores *stores.Storage) *ContestService {
	return &ContestService{stores: stores}
}

func (cs *ContestService) ListContests(ctx context.Context, page int) ([]models.Contest, error) {
	if page < 0 {
		page = 0
	}

	contests, err := cs.stores.Contests.ListContests(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contests: %w", err)
	}

	for i := range contests {
		contests[i].Status = cs.getContestStatus(contests[i])
	}

	return contests, nil
}

func (cs *ContestService) getContestStatus(contest models.Contest) string {
	now := time.Now().Unix()

	if now < contest.RegistrationStartTime {
		return "upcoming"
	} else if now >= contest.RegistrationStartTime && now < contest.RegistrationEndTime {
		return "registration_open"
	} else if now >= contest.RegistrationEndTime && now < contest.StartTime {
		return "registration_closed"
	} else if now >= contest.StartTime && now < contest.EndTime {
		return "active"
	} else {
		return "ended"
	}
}
