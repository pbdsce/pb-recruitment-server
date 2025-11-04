package stores

import (
	"app/internal/common"
	"app/internal/models"
	"app/internal/models/dto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type ContestStore struct {
	db *sql.DB
}

func NewContestStore(db *sql.DB) *ContestStore {
	return &ContestStore{
		db: db,
	}
}

func (s *ContestStore) ListContests(ctx context.Context, page int) ([]models.Contest, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("contest store: db is not initialized")
	}

	const pageSize = 20
	page = max(0, page)
	offset := page * pageSize

	const q = `
		SELECT id, name, registration_start_time, registration_end_time, start_time, end_time, eligible_to
		FROM contests
		ORDER BY start_time DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, q, pageSize, offset)
	if err != nil {
		log.Printf("contest-store: query failed: %v", err)
		return nil, fmt.Errorf("query contests: %w", err)
	}
	defer rows.Close()

	contests := make([]models.Contest, 0)
	for rows.Next() {
		var c models.Contest

		if err := rows.Scan(&c.ID, &c.Name, &c.RegistrationStartTime, &c.RegistrationEndTime, &c.StartTime, &c.EndTime, &c.EligibleTo); err != nil {
			log.Printf("contest-store: row scan failed: %v", err)
			return nil, fmt.Errorf("scan contest row: %w", err)
		}

		contests = append(contests, c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("contest-store: rows error: %v", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return contests, nil
}

func (s *ContestStore) IsRegistered(ctx context.Context, contestID string, userID string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM contest_registrations
			WHERE contest_id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := s.db.QueryRowContext(ctx, q, contestID, userID).Scan(&exists)
	if err != nil {
		log.Printf("contest-store: query failed: %v", err)
		return false, fmt.Errorf("query contest registration: %w", err)
	}

	return exists, nil
}

func (s *ContestStore) GetContest(ctx context.Context, contestID string) (*dto.GetContestResponse, error) {
	const q = `
		SELECT id, name, registration_start_time, registration_end_time, start_time, end_time, eligible_to
		FROM contests
		WHERE id = $1
	`

	var c dto.GetContestResponse
	err := s.db.QueryRowContext(ctx, q, contestID).Scan(
		&c.ID, &c.Name, &c.RegistrationStartTime, &c.RegistrationEndTime, &c.StartTime, &c.EndTime, &c.EligibleTo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ContestNotFoundError
		}
		log.Printf("contest-store: query failed: %v", err)
		return nil, fmt.Errorf("query contest: %w", err)
	}

	return &c, nil
}

func (s *ContestStore) RegisterUser(ctx context.Context, contestID string, userID string) error {
	const q = `
		INSERT INTO contest_registrations (contest_id, user_id, registered_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (contest_id, user_id) DO NOTHING
		`

	res, err := s.db.ExecContext(ctx, q, contestID, userID, time.Now().Unix())
	if err != nil {
		log.Printf("contest-store: query failed: %v", err)
		return fmt.Errorf("query contest registration: %w", err)
	}

	// If rows affected is 0, then the user already registered
	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("user-store: rows error %v", err)
		return fmt.Errorf("rows error: %w", err)
	}

	if affected == 0 {
		log.Printf("user-store: user %s already registered", userID)
		return common.UserAlreadyExistsError
	}

	return nil
}

func (s *ContestStore) UnregisterUser(ctx context.Context, contestID string, userID string) error {
	const q = `
		DELETE FROM contest_registrations
		WHERE contest_id = $1 AND user_id = $2
		`

	res, err := s.db.ExecContext(ctx, q, contestID, userID)
	if err != nil {
		log.Printf("contest-store: query failed: %v", err)
		return fmt.Errorf("query contest registration: %w", err)
	}

	// If rows affected is 0, then the user never registered
	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("user-store: rows error %v", err)
		return fmt.Errorf("rows error: %w", err)
	}

	if affected == 0 {
		log.Printf("user-store: user %s not registered", userID)
		return common.UserNotFoundError
	}

	return nil
}
