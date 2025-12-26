package auth_api

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
)

type authService interface {
    Register(ctx context.Context, username, email, password string) (*models.User, error)
    Login(ctx context.Context, username, password string) (*models.User, error)
    CreatePartnership(ctx context.Context, currentUserID, partnerUsername string) (*models.Partnership, error)
    TerminatePartnership(ctx context.Context, partnershipID string) error
    GetUser(ctx context.Context, userID string) (*models.User, error)
    HandleRoomCompleted(ctx context.Context, userID string, roomNumber int32) error
}

type AuthServiceAPI struct {
    auth_api.UnimplementedAuthServiceServer
    authService authService
}

func NewAuthServiceAPI(authService authService) *AuthServiceAPI {
    return &AuthServiceAPI{
        authService: authService,
    }
}