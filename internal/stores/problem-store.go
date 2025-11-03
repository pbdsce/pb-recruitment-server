package stores

import (
	"app/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type ProblemStore struct {
	db *sql.DB
}

func NewProblemStore(db *sql.DB) *ProblemStore {
	return &ProblemStore{
		db: db,
	}
}

func (s *ProblemStore) CreateProblem(ctx context.Context, p *models.Problem) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("problem store: db is not initialized")
	}

	const q = `
        INSERT INTO problems (id, contest_id, name, score, type, answer)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := s.db.ExecContext(ctx, q,
		p.ID,
		p.ContestID,
		p.Name,
		p.Score,
		p.Type,
		pq.Array(p.Answer),
	)

	if err != nil {
		log.Printf("problem-store: insert failed: %v", err)
		return fmt.Errorf("insert problem: %w", err)
	}

	return nil
}

func (s *ProblemStore) UpdateProblem(ctx context.Context, p *models.Problem) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("problem store: db is not initialized")
	}

	const q = `
        UPDATE problems
        SET name = $3,
            score = $4,
            type = $5,
            answer = $6
        WHERE id = $1 AND contest_id = $2
    `

	_, err := s.db.ExecContext(ctx, q,
		p.ID,
		p.ContestID,
		p.Name,
		p.Score,
		p.Type,
		pq.Array(p.Answer),
	)

	if err != nil {
		log.Printf("problem-store: update failed: %v", err)
		return fmt.Errorf("update problem: %w", err)
	}

	return nil
}

func (s *ProblemStore) DeleteProblem(ctx context.Context, contestID string, problemID string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("problem store: db is not initialized")
	}

	const q = `DELETE FROM problems WHERE id = $1 AND contest_id = $2`

	_, err := s.db.ExecContext(ctx, q, problemID, contestID)

	if err != nil {
		log.Printf("problem-store: delete failed: %v", err)
		return fmt.Errorf("delete problem: %w", err)
	}

	return nil
}
