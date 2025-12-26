package auth_api

import (
    "context"
    
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
)

func (s *AuthServiceAPI) Login(ctx context.Context, req *auth_api.LoginRequest) (*auth_api.LoginResponse, error) {
    user, err := s.authService.Login(ctx, req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    
    return &auth_api.LoginResponse{
        User: &proto_models.User{
            Id:              user.ID,
            Username:        user.Username,
            Email:           user.Email,
            MaxRoomReached:  user.MaxRoomReached,
        },
    }, nil
}