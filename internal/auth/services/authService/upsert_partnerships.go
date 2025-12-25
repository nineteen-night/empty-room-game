package authService

import (
    "context"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) UpsertPartnerships(ctx context.Context, partnerships []*models.Partnership) error {
    if err := s.validatePartnerships(partnerships); err != nil {
        return err
    }

    for _, partnership := range partnerships {
        users, err := s.GetUsersByIDs(ctx, []uint64{partnership.Player1ID, partnership.Player2ID})
        if err != nil {
            return fmt.Errorf("failed to check users: %w", err)
        }
        
        if len(users) != 2 {
            return fmt.Errorf("one or both users not found")
        }
    }

    return s.authStorage.UpsertPartnerships(ctx, partnerships)
}