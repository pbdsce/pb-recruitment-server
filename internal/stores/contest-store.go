package stores

import "app/internal/models"

type ContestStore struct {
	contests map[string]*models.Contest
}

func NewContestStore() *ContestStore {
	return &ContestStore{}
}
