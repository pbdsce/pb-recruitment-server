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