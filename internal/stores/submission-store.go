package stores

import (
	"app/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type SubmissionStore struct {
	db *sql.DB
}

func NewSubmissionStore(db *sql.DB) *SubmissionStore {
	return &SubmissionStore{
		db: db,
	}
}

func (s *SubmissionStore) CreateSubmission(ctx context.Context, submission *models.Submission) (*models.Submission, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
	}

	const q = `
		INSERT INTO submissions (user_id, contest_id, problem_id, type, language, code, option, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at
	`
	if submission.Status == "" {
		submission.Status = models.Pending
	}
	submission.CreatedAt = time.Now().Unix()
	selectedOptionValue := submission.Option

	err := s.db.QueryRowContext(
		ctx,
		q,
		submission.UserID,
		submission.ContestID,
		submission.ProblemID,
		submission.Type,
		submission.Language,
		submission.Code,
		selectedOptionValue,
		submission.Status,
		submission.CreatedAt,
	).Scan(&submission.ID, &submission.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("submission-store: failed to create submission: %w", err)
	}
	return submission, nil
}

func (s *SubmissionStore) GetSubmissionByID(ctx context.Context, id string) (*models.Submission, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
	}

	const q = `
		SELECT id, user_id, contest_id, problem_id, type, language, code, option, status, created_at
		FROM submissions
		WHERE id = $1
	`
	var submission models.Submission
	var selectedOption []int

	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ContestID,
		&submission.ProblemID,
		&submission.Type,
		&submission.Language,
		&submission.Code,
		&selectedOption,
		&submission.Status,
		&submission.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("submission not found")
	}
	if err != nil {
		return nil, fmt.Errorf("submission-store: failed to get submission by ID: %w", err)
	}

	submission.Option = selectedOption
	return &submission, nil
}

func (s *SubmissionStore) ListSubmissionsByProblem(ctx context.Context, userID string, contestID string, problemID string, limit int) ([]models.Submission, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
	}

	const q = `
		SELECT id, user_id, contest_id, problem_id, type, language, code, option, status, created_at
		FROM submissions
		WHERE user_id = $1 AND contest_id = $2 AND problem_id = $3
		ORDER BY created_at DESC
		LIMIT $4
	`

	rows, err := s.db.QueryContext(ctx, q, userID, contestID, problemID, limit)
	if err != nil {
		return nil, fmt.Errorf("submission-store: failed to list submissions by problem: %w", err)
	}
	defer rows.Close()

	submissions := make([]models.Submission, 0)
	for rows.Next() {
		var submission models.Submission
		var selectedOption []int

		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ContestID,
			&submission.ProblemID,
			&submission.Type,
			&submission.Language,
			&submission.Code,
			&selectedOption,
			&submission.Status,
			&submission.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("submission-store: failed to scan row: %w", err)
		}

		submission.Option = selectedOption
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("submission-store: iteration error: %w", err)
	}

	return submissions, nil
}

func (s *SubmissionStore) GetJudgeResultBySubmissionID(ctx context.Context, submissionID string) (*JudgeResult, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("submission store: db is not initialized")
	}

	const resultQuery = `
		SELECT 
			submission_id, 
			status, 
			runtime_ms, 
			memory_kb
		FROM judge_results
		WHERE submission_id = $1
	` //judge_results name can be varied once the actual table is created

	judgeResult := &JudgeResult{SubmissionID: submissionID}

	err := s.db.QueryRowContext(ctx, resultQuery, submissionID).Scan(
		&judgeResult.SubmissionID,
		&judgeResult.Status,
		&judgeResult.RuntimeMs,
		&judgeResult.MemoryKB,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("judge result data not found for this submission")
		}
		return nil, fmt.Errorf("submission store: failed to get judge result header: %w", err)
	}

	const testCaseQuery = `
		SELECT 
			id, 
			status, 
			duration_ms
		FROM test_case_results
		WHERE submission_id = $1
		ORDER BY id ASC
	`
	rows, err := s.db.QueryContext(ctx, testCaseQuery, submissionID)
	if err != nil {
		return nil, fmt.Errorf("submission store: failed to query test case results: %w", err)
	}
	defer rows.Close()

	judgeResult.TestCases = make([]TestCaseResult, 0)

	for rows.Next() {
		var tc TestCaseResult
		err := rows.Scan(
			&tc.ID,
			&tc.Status,
			&tc.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("submission store: failed to scan test case result: %w", err)
		}
		judgeResult.TestCases = append(judgeResult.TestCases, tc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("submission store: error iterating through test case results: %w", err)
	}

	return judgeResult, nil
}
