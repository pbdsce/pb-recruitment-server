package stores

import (
	"app/internal/models/dto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type RankingStore struct {
	db *sql.DB
}

func NewRankingStore(db *sql.DB) *RankingStore {
	return &RankingStore{
		db: db,
	}
}

func (s *RankingStore) UpdateLeaderboardUser(ctx context.Context, contestID string, userID string, req *dto.UpdateLeaderboardUserRequest) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("ranking store: db is not initialized")
	}

	query := "UPDATE rankings SET "
	args := []interface{}{}
	argID := 1

	if req.Hidden != nil {
		query += fmt.Sprintf("hidden = $%d, ", argID)
		args = append(args, *req.Hidden)
		argID++
	}
	if req.Disqualified != nil {
		query += fmt.Sprintf("disqualified = $%d, ", argID)
		args = append(args, *req.Disqualified)
		argID++
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query = strings.TrimSuffix(query, ", ")

	query += fmt.Sprintf(" WHERE contest_id = $%d AND user_id = $%d", argID, argID+1)
	args = append(args, contestID, userID)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("ranking-store: update failed: %v", err)
		return fmt.Errorf("update ranking: %w", err)
	}

	return nil
}
