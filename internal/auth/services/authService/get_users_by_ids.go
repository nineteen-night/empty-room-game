package authService

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) GetUsersByIDs(ctx context.Context, ids []uint64) ([]*models.User, error) {
    return s.authStorage.GetUsersByIDs(ctx, ids)
}