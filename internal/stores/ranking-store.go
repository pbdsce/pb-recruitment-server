package stores

import (
	"database/sql"
)

type RankingStore struct {
	db *sql.DB
}

func NewRankingStore(db *sql.DB) *RankingStore {
	return &RankingStore{
		db: db,
	}
}
