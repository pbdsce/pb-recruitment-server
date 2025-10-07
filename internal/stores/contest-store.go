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

// SAMPLE -- fetches all contests from the database.
func (s *ContestStore) ListContests(ctx context.Context) ([]models.Contest, error) {

	if s == nil || s.db == nil {
		return nil, fmt.Errorf("contest store: db is not initialized")
	}

	const q = `SELECT id, name, registration_start_time, registration_end_time, start_time, end_time FROM contests`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		log.Printf("contest-store: query failed: %v", err)
		return nil, fmt.Errorf("query contests: %w", err)
	}
	defer rows.Close()

	var contests []models.Contest
	for rows.Next() {
		var c models.Contest
		var regStart, regEnd, start, end time.Time

		if err := rows.Scan(&c.ID, &c.Name, &regStart, &regEnd, &start, &end); err != nil {
			log.Printf("contest-store: row scan failed: %v", err)
			return nil, fmt.Errorf("scan contest row: %w", err)
		}

		c.RegistrationStartTime = regStart.Unix()
		c.RegistrationEndTime = regEnd.Unix()
		c.StartTime = start.Unix()
		c.EndTime = end.Unix()

		contests = append(contests, c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("contest-store: rows error: %v", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return contests, nil
}
