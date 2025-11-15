package stores

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type AdminStore struct {
	db *sql.DB
}

func NewAdminStore(db *sql.DB) *AdminStore {
	return &AdminStore{db: db}
}

func (s *AdminStore) IsAdmin(ctx context.Context, userID string) (bool, error) {
	if s == nil || s.db == nil {
		return false, fmt.Errorf("admin store: db is not initialized")
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM admin WHERE user_id = $1)"

	row := s.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&exists)
	if err != nil {
		log.Printf("admin-store: query failed: %v", err)
		return false, fmt.Errorf("query admin status: %w", err)
	}

	return exists, nil
}
