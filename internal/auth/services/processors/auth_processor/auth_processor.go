package auth_processor

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

type authService interface {
    UpsertUsers(ctx context.Context, users []*models.User) error
    UpsertPartnerships(ctx context.Context, partnerships []*models.Partnership) error
}

type AuthProcessor struct {
    authService authService
}

func NewAuthProcessor(authService authService) *AuthProcessor {
    return &AuthProcessor{
        authService: authService,
    }
}