package authService

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) GetPartnershipsByIDs(ctx context.Context, ids []uint64) ([]*models.Partnership, error) {
    return s.authStorage.GetPartnershipsByIDs(ctx, ids)
}