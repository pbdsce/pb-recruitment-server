package stores

import "app/internal/models"

type ProblemStore struct {
	problems map[string]*models.Problem
}

func NewProblemStore() *ProblemStore {
	return &ProblemStore{}
}
