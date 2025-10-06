package stores

import "app/internal/models"

type SubmissionStore struct {
	submissions map[string]*models.Submission
}

func NewSubmissionStore() *SubmissionStore {
	return &SubmissionStore{}
}
