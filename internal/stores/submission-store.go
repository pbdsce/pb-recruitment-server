package stores

import (
	"app/internal/models"
	"app/internal/common"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/google/uuid"
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
		SELECT user_id, contest_id, problem_id, type, language, code, choices, status, created_at, runtime, memory
		FROM submissions
		WHERE id = $1
	`
	var sub models.Submission
	sub.ID = id

	var rawChoices sql.NullString

	row := s.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(
		&sub.UserID,
		&sub.ContestID,
		&sub.ProblemID,
		&sub.Type,
		&sub.Language,
		&sub.Code,
		&rawChoices,
		&sub.Status,
		&sub.CreatedAt,
		&sub.Runtime,
		&sub.Memory,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrNotFound
		}
		log.Printf("submission-store: row scan failed for ID %s: %v", id, err)
		return nil, fmt.Errorf("scan submission: %w", err)
	}

	sub.Option = []int{}
	if rawChoices.Valid && rawChoices.String != "" && rawChoices.String != "{}" {
		choiceStr := strings.TrimSpace(strings.Trim(rawChoices.String, "{}"))

		if choiceStr != "" {
			parts := strings.Split(choiceStr, ",")

			for _, part := range parts {
				val, err := strconv.Atoi(strings.TrimSpace(part))
				if err != nil {
					log.Printf("submission-store: failed to parse choice value '%s': %v", part, err)
					continue
				}
				sub.Option = append(sub.Option, val)
			}
		}
	}

	testCaseResults, err := s.GetTestCaseResultsBySubmissionID(ctx, id)
	if err != nil {
		log.Printf("submission-store: failed to get test case results for submission ID %s: %v", id, err)
		sub.TestCaseResults = []models.TestCaseResult{}
	} else {
		sub.TestCaseResults = testCaseResults
	}
	
	return &sub, nil
}

func (s *SubmissionStore) GetTestCaseResultsBySubmissionID(ctx context.Context, submissionID string) ([]models.TestCaseResult, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
	}

	const q = `
		SELECT id, submission_id, test_case_id, status, runtime, memory, created_at
		FROM test_case_results
		WHERE submission_id = $1
		ORDER BY created_at ASC
	`
	rows, err := s.db.QueryContext(ctx, q, submissionID)
	if err != nil {
		return nil, fmt.Errorf("query test case results: %w", err)
	}
	defer rows.Close()

	var results []models.TestCaseResult
	for rows.Next() {
		var res models.TestCaseResult
		if err := rows.Scan(
			&res.ID,
			&res.SubmissionID,
			&res.TestCaseID,
			&res.Status,
			&res.Runtime,
			&res.Memory,
			&res.CreatedAt,
		); err != nil {
			log.Printf("submission-store: failed to scan test case result row for submission  %s: %v", submissionID, err)
			continue
		}
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (s *SubmissionStore) ListUserSubmissionsByProblemID(ctx context.Context, userID, problemID string, page int) ([]models.Submission, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
	}

	const pageSize = 20
	page = max(0, page)
	offset := page * pageSize

	const q = `
		SELECT id, contest_id, problem_id, type, language, status, created_at, runtime, memory
		FROM submissions
		WHERE user_id = $1 AND problem_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	
	rows, err := s.db.QueryContext(ctx, q, userID, problemID, pageSize, offset)
	if err != nil {
		log.Printf("submission-store: query failed: %v", err)
		return nil, fmt.Errorf("query user submissions: %w", err)
	}
	defer rows.Close()

	submissions := make([]models.Submission, 0)
	for rows.Next() {
		var sub models.Submission

		if err := rows.Scan(
			&sub.ID,
			&sub.ContestID,
			&sub.ProblemID,
			&sub.Type,
			&sub.Language,
			&sub.Status,
			&sub.CreatedAt,
			&sub.Runtime,
			&sub.Memory,
		); err != nil {
			log.Printf("submission-store: failed to scan submission row: %v", err)
			continue
		}
		submissions = append(submissions, sub)
	}

	if err := rows.Err(); err != nil {
		log.Printf("submission-store: rows error: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return submissions, nil
}

func (s *SubmissionStore) CreateSubmission(ctx context.Context, sub *models.Submission) (string, error) {
	if s == nil || s.db == nil {
		return "", fmt.Errorf("submission store: db is not initialized")
	}

	sub.ID = uuid.NewString()
	sub.CreatedAt = time.Now().Unix()

	choiceStrings := make([]string, len(sub.Option))
	for i, choice := range sub.Option {
		choiceStrings[i] = strconv.Itoa(choice)
	}
	mcqChoices := fmt.Sprintf("{%s}", strings.Join(choiceStrings, ","))

	const q = `
		INSERT INTO 
		submissions (id, user_id, contest_id, problem_id, type, language, code, choices, status, created_at, runtime, memory)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	var submissionID string
	err := s.db.QueryRowContext(ctx, q,
		sub.ID,
		sub.UserID,
		sub.ContestID,
		sub.ProblemID,
		sub.Type,
		sub.Language,
		sub.Code,
		mcqChoices,
		sub.Status,
		sub.CreatedAt,
		sub.Runtime,
		sub.Memory,
	).Scan(&submissionID)

	if err != nil {
		log.Printf("submission-store: failed to insert submission: %v", err)
		return "", fmt.Errorf("insert submission: %w", err)
	}

	return submissionID, nil
}