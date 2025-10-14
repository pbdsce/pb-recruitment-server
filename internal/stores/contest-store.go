package stores

import (
	"app/internal/models"
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
		SELECT id, name, registration_start_time, registration_end_time, start_time, end_time
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

		if err := rows.Scan(&c.ID, &c.Name, &c.RegistrationStartTime, &c.RegistrationEndTime, &c.StartTime, &c.EndTime); err != nil {
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

func (s *ContestStore) GetContest(ctx context.Context, contestID string) (*models.Contest, error) {

	const q = `
		SELECT id, name, registration_start_time, registration_end_time, start_time, end_time
		FROM contests
		WHERE id = $1
	`

	var c models.Contest
	err := s.db.QueryRowContext(ctx, q, contestID).Scan(
		&c.ID, &c.Name, &c.RegistrationStartTime, &c.RegistrationEndTime, &c.StartTime, &c.EndTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("contest-store: query failed: %v", err)
		return nil, fmt.Errorf("query contest: %w", err)
	}

	return &c, nil
}

func (s *ContestStore) IsUserRegistered(ctx context.Context, contestID, userID string) (bool, error) {

	const q = `
		SELECT EXISTS(
			SELECT 1 FROM contest_registrations 
			WHERE contest_id = $1 AND user_id = $2
		)
	`
	var exists bool
	err := s.db.QueryRowContext(ctx, q, contestID, userID).Scan(&exists)
	if err != nil {
		log.Printf("contest-store: query failed: %v", err)
		return false, fmt.Errorf("query user registration: %w", err)
	}

	return exists, nil
}

func (s *ContestStore) RegisterUser(ctx context.Context, contestID, userID string) error {

	now := time.Now().Unix()
	const q = `
		INSERT INTO contest_registrations (contest_id, user_id, registered_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (contest_id, user_id) DO NOTHING
	`

	_, err := s.db.ExecContext(ctx, q, contestID, userID, now)
	if err != nil {
		log.Printf("contest-store: register user failed: %v", err)
		return fmt.Errorf("register user: %w", err)
	}

	return nil
}

func (s *ContestStore) UnregisterUser(ctx context.Context, contestID, userID string) error {

	const q = `DELETE FROM contest_registrations WHERE contest_id = $1 AND user_id = $2`

	_, err := s.db.ExecContext(ctx, q, contestID, userID)
	if err != nil {
		log.Printf("contest-store: unregister user failed: %v", err)
		return fmt.Errorf("unregister user: %w", err)
	}

	return nil
}
