package auth_api

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
)

type authService interface {
    GetUsersByIDs(ctx context.Context, ids []uint64) ([]*models.User, error)
    UpsertUsers(ctx context.Context, users []*models.User) error
    GetPartnershipsByIDs(ctx context.Context, ids []uint64) ([]*models.Partnership, error)
    UpsertPartnerships(ctx context.Context, partnerships []*models.Partnership) error
}

// AuthServiceAPI реализует grpc AuthServiceServer
type AuthServiceAPI struct {
    auth_api.UnimplementedAuthServiceServer
    authService authService
}

func NewAuthServiceAPI(authService authService) *AuthServiceAPI {
    return &AuthServiceAPI{
        authService: authService,
    }
}