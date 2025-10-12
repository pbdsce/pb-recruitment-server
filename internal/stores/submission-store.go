package stores

import (
	"app/internal/models"
	"app/internal/common"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type SubmissionStore struct {
	db *sql.DB
}

func NewSubmissionStore(db *sql.DB) *SubmissionStore {
	return &SubmissionStore{
		db: db,
	}
}

func (s *SubmissionStore) GetSubmissionStatusByID(ctx context.Context, id string) (models.Submission, error) {
	if s == nil || s.db == nil {
		return models.Submission{}, fmt.Errorf("submission store: db is not initialized")
	}

	const q = `
		SELECT status, user_id
		FROM submissions
		WHERE id = $1
	`
	var sub models.Submission
	sub.ID = id

	row := s.db.QueryRowContext(ctx, q, id)

	if err := row.Scan(&sub.Status, &sub.UserID); err != nil {
		if err == sql.ErrNoRows {
			return models.Submission{}, common.ErrNotFound
		}
		log.Printf("submission-store: row scan failed for ID %s: %v", id, err)
		return models.Submission{}, fmt.Errorf("scan submission: %w", err)
	}

	return sub, nil
}