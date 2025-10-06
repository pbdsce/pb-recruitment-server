package stores

import "app/internal/models"

type RankingStore struct {
	rankings map[string]*models.Ranking
}

func NewRankingStore() *RankingStore {
	return &RankingStore{}
}
