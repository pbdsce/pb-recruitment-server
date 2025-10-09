package services

import (
	"app/internal/models"
	"app/internal/models/dto"
	"app/internal/stores"
	"context"
	"errors"
	"fmt"
)

type SubmissionService struct {
	SubmissionStore stores.Submissions
}

func NewSubmissionService(submissionStore stores.Submissions) *SubmissionService {
	return &SubmissionService{
		SubmissionStore: submissionStore,
	}
}

func (ss *SubmissionService) SubmitSubmission(ctx context.Context, userID string, req *dto.SubmitSubmissionRequest) (*dto.SubmitSubmissionResponse, error) {
	var submissionType models.SubmissionType

	if req.Code != "" {
		submissionType = models.Code
		if req.Language == "" {
			return nil, errors.New("language is required for code submissions")
		}
	} else if len(req.Option) > 0 {
		submissionType = models.MCQ
	} else {
		return nil, errors.New("submission must contain either code or selected options")
	}

	submission := &models.Submission{
		UserID:    userID,
		ContestID: req.ContestID,
		ProblemID: req.ProblemID,
		Type:      submissionType,
		Language:  req.Language,
		Code:      req.Code,
		Option:    req.Option,
		Status:    models.Pending,
	}

	createdSubmission, err := ss.SubmissionStore.CreateSubmission(ctx, submission)
	if err != nil {
		return nil, fmt.Errorf("service: failed to save submission: %w", err)
	}

	return &dto.SubmitSubmissionResponse{
		SubmissionID: createdSubmission.ID,
	}, nil
}

func (ss *SubmissionService) GetSubmission(ctx context.Context, userID string, submissionID string) (*models.Submission, error) {
	submission, err := ss.SubmissionStore.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}

	if submission.UserID != userID {
		return nil, errors.New("unauthorized access: submission does not belong to user")
	}

	if submission.Status == models.Pending {
		submission.Code = "Submission is pending judgment."
	}

	return submission, nil
}

func (ss *SubmissionService) ListProblemSubmissions(ctx context.Context, userID string, contestID string, problemID string, limit int) (*dto.ListProblemSubmissionsResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	submissions, err := ss.SubmissionStore.ListSubmissionsByProblem(ctx, userID, contestID, problemID, limit)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list submissions: %w", err)
	}

	for i := range submissions {
		submissions[i].Code = ""
		submissions[i].Option = nil
		submissions[i].Language = "N/A"
	}

	return &dto.ListProblemSubmissionsResponse{
		Submissions: submissions,
	}, nil
}

func (ss *SubmissionService) GetSubmissionStatus(ctx context.Context, userID string, submissionID string) (models.SubmissionStatus, error) {
	submission, err := ss.GetSubmission(ctx, userID, submissionID)
	if err != nil {
		return "", err
	}
	return submission.Status, nil
}

func (ss *SubmissionService) GetSubmissionDetails(ctx context.Context, userID string, submissionID string) (*dto.SubmissionDetailsResponse, error) {
	submission, err := ss.SubmissionStore.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}

	if submission.UserID != userID {
		return nil, errors.New("unauthorized access: submission details do not belong to user")
	}

	if submission.Status == models.Pending {
		return nil, errors.New("submission result is pending and not yet available")
	}

	judgeResult, err := ss.SubmissionStore.GetJudgeResultBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to fetch judge results: %w", err)
	}

	response := &dto.SubmissionDetailsResponse{
		SubmissionID: judgeResult.SubmissionID,
		Status:       judgeResult.Status,
		RuntimeMs:    judgeResult.RuntimeMs,
		MemoryKB:     judgeResult.MemoryKB,
		TestCases:    make([]dto.TestCaseResultDTO, len(judgeResult.TestCases)),
	}

	for i, tc := range judgeResult.TestCases {
		response.TestCases[i] = dto.TestCaseResultDTO{
			ID:       tc.ID,
			Status:   tc.Status,
			Duration: tc.Duration,
		}
	}

	return response, nil
}
