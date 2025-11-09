package services

import (
	"app/internal/stores"
	"context"
)

type AdminService struct {
	stores *stores.Storage
}

func NewAdminService(stores *stores.Storage) *AdminService {
	return &AdminService{stores: stores}
}

// IsAdmin checks if a user is an admin.
func (s *AdminService) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return s.stores.Admins.IsAdmin(ctx, userID)
}
