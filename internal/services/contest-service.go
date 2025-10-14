package services

import (
	"app/internal/models"
	"app/internal/models/dto"
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

func (cs *ContestService) ModifyRegistration(ctx context.Context, contestID string, userID string, action string) error {
	// Validate contest exists
	contest, err := cs.stores.Contests.GetContest(ctx, contestID)
	if err != nil {
		return fmt.Errorf("failed to get contest: %w", err)
	}
	if contest == nil {
		return fmt.Errorf("contest not found")
	}

	// Check if contest is in registration period
	now := time.Now().Unix()
	if now < contest.RegistrationStartTime || now >= contest.RegistrationEndTime {
		return fmt.Errorf("contest registration is not currently open")
	}

	// Check if user is already registered
	isRegistered, err := cs.stores.Contests.IsUserRegistered(ctx, contestID, userID)
	if err != nil {
		return fmt.Errorf("failed to check registration status: %w", err)
	}

	switch action {
	case "register":
		if isRegistered {
			return fmt.Errorf("user is already registered for this contest")
		}
		return cs.stores.Contests.RegisterUser(ctx, contestID, userID)

	case "unregister":
		if !isRegistered {
			return fmt.Errorf("user is not registered for this contest")
		}
		return cs.stores.Contests.UnregisterUser(ctx, contestID, userID)

	default:
		return fmt.Errorf("invalid action: %s", action)
	}
}

func (cs *ContestService) ListContests(ctx context.Context, page int) ([]models.Contest, error) {
	// Validate page parameter
	if page < 0 {
		page = 0
	}

	// Get contests from store
	contests, err := cs.stores.Contests.ListContests(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contests: %w", err)
	}

	// Add status to each contest
	for i := range contests {
		contests[i].Status = cs.getContestStatus(contests[i])
	}

	return contests, nil
}

// getContestStatus determines the current status of a contest based on timestamps
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

func (cs *ContestService) GetContest(ctx context.Context, contestID string, userID string, isAuthenticated bool) (*dto.GetContestResponse, error) {
	// Get contest details
	contest, err := cs.stores.Contests.GetContest(ctx, contestID)
	if err != nil {
		return nil, err
	}

	if contest == nil {
		return nil, nil // Contest not found
	}

	// Add status to contest
	contest.Status = cs.getContestStatus(*contest)

	response := &dto.GetContestResponse{
		Contest: *contest,
	}

	// If user is authenticated, check if they're registered
	if isAuthenticated && userID != "" {
		isRegistered, err := cs.stores.Contests.IsUserRegistered(ctx, contestID, userID)
		if err != nil {
			return nil, err
		}
		response.IsRegistered = isRegistered
	}

	return response, nil
}

func (cs *ContestService) GetContestProblemsList(ctx context.Context, contestID string, userID string) (*dto.GetContestProblemsResponse, error) {
	// Validate contest exists
	contest, err := cs.stores.Contests.GetContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest: %w", err)
	}
	if contest == nil {
		return nil, fmt.Errorf("contest not found")
	}

	// Check if user is registered for the contest
	isRegistered, err := cs.stores.Contests.IsUserRegistered(ctx, contestID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check registration status: %w", err)
	}
	if !isRegistered {
		return nil, fmt.Errorf("user is not registered for this contest")
	}

	// Get problems from store
	problems, err := cs.stores.Problems.GetContestProblems(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest problems: %w", err)
	}

	// Convert to ProblemOverview and add status
	problemOverviews := make([]dto.ProblemOverview, len(problems))
	for i, problem := range problems {
		status := "not_attempted"
		// TODO: Check user's submission status for this problem

		problemOverviews[i] = dto.ProblemOverview{
			ID:     problem.ID,
			Name:   problem.Name,
			Score:  problem.Score,
			Type:   string(problem.Type),
			Status: status,
		}
	}

	return &dto.GetContestProblemsResponse{
		ContestID: contestID,
		Problems:  problemOverviews,
	}, nil
}

func (cs *ContestService) GetContestProblemStatement(ctx context.Context, contestID string, problemID string, userID string) (*dto.GetProblemStatementResponse, error) {
	// Validate contest exists
	contest, err := cs.stores.Contests.GetContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest: %w", err)
	}
	if contest == nil {
		return nil, fmt.Errorf("contest not found")
	}

	// Check if user is registered for the contest
	isRegistered, err := cs.stores.Contests.IsUserRegistered(ctx, contestID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check registration status: %w", err)
	}
	if !isRegistered {
		return nil, fmt.Errorf("user is not registered for this contest")
	}

	// Get problem from store
	problem, err := cs.stores.Problems.GetProblem(ctx, problemID, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}
	if problem == nil {
		return nil, nil // Problem not found
	}

	// Convert to GetProblemStatementResponse
	response := &dto.GetProblemStatementResponse{
		ProblemID:   problem.ID,
		ContestID:   problem.ContestID,
		Name:        problem.Name,
		Description: problem.Description,
		Score:       problem.Score,
		Type:        string(problem.Type),
	}

	// // Add type-specific fields
	// if problem.Type == models.Code {
	// 	// TODO: Add time_limit and memory_limit fields to Problem model if needed
	// 	response.TimeLimit = 1000  // Default 1 second
	// 	response.MemoryLimit = 256 // Default 256 MB
	// } else if problem.Type == models.MCQ {
	// 	// TODO: Add options and multiple fields to Problem model if needed
	// 	response.Options = []string{"Option A", "Option B", "Option C", "Option D"} // Placeholder
	// 	response.Multiple = false                                                   // Default single choice
	// }

	return response, nil
}
