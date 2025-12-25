package authService

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

type AuthStorage interface {
    GetUsersByIDs(ctx context.Context, ids []uint64) ([]*models.User, error)
    UpsertUsers(ctx context.Context, users []*models.User) error
    GetPartnershipsByIDs(ctx context.Context, ids []uint64) ([]*models.Partnership, error)
    UpsertPartnerships(ctx context.Context, partnerships []*models.Partnership) error
    GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type AuthService struct {
    authStorage AuthStorage
    minNameLen  int
    maxNameLen  int
}

func NewAuthService(ctx context.Context, authStorage AuthStorage, minNameLen, maxNameLen int) *AuthService {
    return &AuthService{
        authStorage: authStorage,
        minNameLen:  minNameLen,
        maxNameLen:  maxNameLen,
    }
}