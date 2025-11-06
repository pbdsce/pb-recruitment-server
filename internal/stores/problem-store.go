package stores

import (
	"app/internal/common"
	"app/internal/models/dto"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type ProblemStore struct {
	db *sql.DB
}

func NewProblemStore(db *sql.DB) *ProblemStore {
	return &ProblemStore{
		db: db,
	}
}

func (s *ProblemStore) GetProblemList(ctx context.Context, contestID string) ([]dto.ProblemOverview, error) {
	const q = `
		SELECT id, name, score, type
		FROM problems
		WHERE contest_id = $1
	`

	rows, err := s.db.QueryContext(ctx, q, contestID)
	defer rows.Close()
	if err != nil {
		log.Printf("problem-store: query failed: %v", err)
		return nil, fmt.Errorf("query contest problems: %w", err)
	}

	var problems []dto.ProblemOverview
	for rows.Next() {
		var p dto.ProblemOverview

		if err := rows.Scan(&p.ID, &p.Name, &p.Score, &p.Type); err != nil {
			log.Printf("problem-store: row scan failed: %v", err)
			return nil, fmt.Errorf("scan problem row: %w", err)
		}

		problems = append(problems, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("problem-store: rows error: %v", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(problems) == 0 {
		log.Printf("Failed to find contest problems for contest %s", contestID)
		return nil, common.ContestNotFoundError
	}

	return problems, nil
}

func (s *ProblemStore) GetProblem(ctx context.Context, problemID string, contestID string) (*dto.GetProblemStatementResponse, error) {
	const q = `
		SELECT id, contest_id, name, description, score, type
		FROM problems
		WHERE id = $1 AND contest_id = $2
	`

	var p dto.GetProblemStatementResponse

	err := s.db.QueryRowContext(ctx, q, problemID, contestID).Scan(
		&p.ProblemID, &p.ContestID, &p.Name, &p.Description, &p.Score, &p.Type,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Failed to find problem for contest %s and problem %s", contestID, problemID)
			return nil, common.ContestNotFoundError
		}
		log.Printf("problem-store: query failed: %v", err)
		return nil, fmt.Errorf("query problem: %w", err)
	}

	return &p, nil
}
