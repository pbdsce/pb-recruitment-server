package stores

import "app/internal/models"

type UserStore struct {
	users map[string]*models.User
}

func NewUserStore() *UserStore {
	return &UserStore{}
}
