package stores

import (
	"app/internal/models"
	"context"
	"database/sql"
	"encoding/json"
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

func (s *ProblemStore) GetContestProblems(ctx context.Context, contestID string) ([]models.Problem, error) {
	const q = `
		SELECT id, contest_id, name, description, score, type, answer
		FROM problems
		WHERE contest_id = $1
		ORDER BY name
	`

	rows, err := s.db.QueryContext(ctx, q, contestID)
	if err != nil {
		log.Printf("problem-store: query failed: %v", err)
		return nil, fmt.Errorf("query contest problems: %w", err)
	}
	defer rows.Close()

	var problems []models.Problem
	for rows.Next() {
		var p models.Problem
		var answer sql.NullString
		var description sql.NullString

		if err := rows.Scan(&p.ID, &p.ContestID, &p.Name, &description, &p.Score, &p.Type, &answer); err != nil {
			log.Printf("problem-store: row scan failed: %v", err)
			return nil, fmt.Errorf("scan problem row: %w", err)
		}

		if description.Valid {
			p.Description = description.String
		}

		if answer.Valid {
			// Parse the JSON array string into []int
			var answerInts []int
			if err := json.Unmarshal([]byte(answer.String), &answerInts); err != nil {
				log.Printf("problem-store: failed to parse answer JSON: %v", err)
				// Continue with empty answer if parsing fails
				p.Answer = []int{}
			} else {
				p.Answer = answerInts
			}
		} else {
			p.Answer = []int{}
		}

		problems = append(problems, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("problem-store: rows error: %v", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return problems, nil
}

func (s *ProblemStore) GetProblem(ctx context.Context, problemID string, contestID string) (*models.Problem, error) {
	const q = `
		SELECT id, contest_id, name, description, score, type, answer
		FROM problems
		WHERE id = $1 AND contest_id = $2
	`

	var p models.Problem
	var answer sql.NullString
	var description sql.NullString

	err := s.db.QueryRowContext(ctx, q, problemID, contestID).Scan(
		&p.ID, &p.ContestID, &p.Name, &description, &p.Score, &p.Type, &answer,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Problem not found
		}
		log.Printf("problem-store: query failed: %v", err)
		return nil, fmt.Errorf("query problem: %w", err)
	}

	if description.Valid {
		p.Description = description.String
	}

	if answer.Valid {
		// Parse the JSON array string into []int
		var answerInts []int
		if err := json.Unmarshal([]byte(answer.String), &answerInts); err != nil {
			log.Printf("problem-store: failed to parse answer JSON: %v", err)
			// Continue with empty answer if parsing fails
			p.Answer = []int{}
		} else {
			p.Answer = answerInts
		}
	} else {
		p.Answer = []int{}
	}

	return &p, nil
}
