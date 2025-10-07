package stores

import (
	"app/internal/models"
	"context"
	"database/sql"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	//example
	//interaction with db
	user := &models.User{}
	return user, nil
}
