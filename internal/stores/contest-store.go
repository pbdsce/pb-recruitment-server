package stores

import (
	"app/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
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

func (s *ContestStore) CreateContest(ctx context.Context, c *models.Contest) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("contest store: db is not initialized")
	}

	const q = `
        INSERT INTO contests (id, name, registration_start_time, registration_end_time, start_time, end_time)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := s.db.ExecContext(ctx, q,
		c.ID,
		c.Name,
		c.RegistrationStartTime,
		c.RegistrationEndTime,
		c.StartTime,
		c.EndTime,
	)

	if err != nil {
		log.Printf("contest-store: insert failed: %v", err)
		return fmt.Errorf("insert contest: %w", err)
	}

	return nil
}

func (s *ContestStore) UpdateContest(ctx context.Context, c *models.Contest) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("contest store: db is not initialized")
	}

	const q = `
        UPDATE contests
        SET name = $2,
            registration_start_time = $3,
            registration_end_time = $4,
            start_time = $5,
            end_time = $6
        WHERE id = $1
    `
	_, err := s.db.ExecContext(ctx, q,
		c.ID,
		c.Name,
		c.RegistrationStartTime,
		c.RegistrationEndTime,
		c.StartTime,
		c.EndTime,
	)

	if err != nil {
		log.Printf("contest-store: update failed: %v", err)
		return fmt.Errorf("update contest: %w", err)
	}
	return nil
}

func (s *ContestStore) DeleteContest(ctx context.Context, contestID string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("contest store: db is not initialized")
	}

	const q = `DELETE FROM contests WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, contestID)
	if err != nil {
		log.Printf("contest-store: delete failed: %v", err)
		return fmt.Errorf("delete contest: %w", err)
	}
	return nil
}
