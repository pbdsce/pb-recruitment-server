package services

import (
	"app/internal/models"
	"app/internal/stores"
	"context"
)

type SubmissionService struct {
	stores *stores.Storage
}

func NewSubmissionService(stores *stores.Storage) *SubmissionService {
	return &SubmissionService{stores: stores}
}

func (ss *SubmissionService) GetSubmissionStatusByID(ctx context.Context, id string) (*models.Submission, error) {
	sub, err := ss.stores.Submissions.GetSubmissionStatusByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (ss *SubmissionService) GetSubmissionDetailsByID(ctx context.Context, id string) (*models.Submission, error) {
	sub, err := ss.stores.Submissions.GetSubmissionDetailsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (ss *SubmissionService) ListUserSubmissionsByProblemID(ctx context.Context, userID, problemID string, page int) ([]models.Submission, error) {
	sub, err := ss.stores.Submissions.ListUserSubmissionsByProblemID(ctx, userID, problemID, page)
	if err != nil {
		return nil, err
	}
	return sub, nil
}