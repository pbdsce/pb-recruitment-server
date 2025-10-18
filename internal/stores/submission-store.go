package stores

import (
	"app/internal/models"
	"app/internal/common"
	"context"
	"database/sql"
	"fmt"
	"log"
	"github.com/lib/pq"
	"encoding/json"
)

type SubmissionStore struct {
	db *sql.DB
}

func NewSubmissionStore(db *sql.DB) *SubmissionStore {
	return &SubmissionStore{
		db: db,
	}
}

func (s *SubmissionStore) GetSubmissionStatusByID(ctx context.Context, id string) (*models.Submission, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
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
			return nil, common.ErrNotFound
		}
		log.Printf("submission-store: row scan failed for ID %s: %v", id, err)
		return nil, fmt.Errorf("scan submission: %w", err)
	}

	return &sub, nil
}

func (s *SubmissionStore) GetSubmissionDetailsByID(ctx context.Context, id string) (*models.Submission, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
	}

	const q = `
		SELECT user_id, contest_id, problem_id, type, language, code, choices, status, created_at, runtime, memory, test_case_results
		FROM submissions
		WHERE id = $1
	`
	var sub models.Submission
	sub.ID = id

	var rawTestCaseResults string
	row := s.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(
		&sub.UserID,
		&sub.ContestID,
		&sub.ProblemID,
		&sub.Type,
		&sub.Language,
		&sub.Code,
		pq.Array(&sub.Option),
		&sub.Status,
		&sub.CreatedAt,
		&sub.Runtime,
		&sub.Memory,
		&rawTestCaseResults,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrNotFound
		}
		log.Printf("submission-store: row scan failed for ID %s: %v", id, err)
		return nil, fmt.Errorf("scan submission: %w", err)
	}

	if err := json.Unmarshal([]byte(rawTestCaseResults), &sub.TestCaseResults); err != nil {
		log.Printf("submission-store: failed to unmarshal test_case_results for submission ID %s: %v", id, err)
		return nil, fmt.Errorf("unmarshal test_case_results: %w", err)
	}
	
	return &sub, nil
}