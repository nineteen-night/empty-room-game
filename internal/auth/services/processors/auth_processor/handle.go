package auth_processor

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (p *AuthProcessor) HandleUsers(ctx context.Context, user *models.User) error {
    return p.authService.UpsertUsers(ctx, []*models.User{user})
}

func (p *AuthProcessor) HandlePartnerships(ctx context.Context, partnership *models.Partnership) error {
    return p.authService.UpsertPartnerships(ctx, []*models.Partnership{partnership})
}