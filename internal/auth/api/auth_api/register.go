package auth_api

import (
    "context"
    
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
)

func (s *AuthServiceAPI) Register(ctx context.Context, req *auth_api.RegisterRequest) (*auth_api.RegisterResponse, error) {
    user, err := s.authService.Register(ctx, req.Username, req.Email, req.Password)
    if err != nil {
        return nil, err
    }
    
    return &auth_api.RegisterResponse{
        User: &proto_models.User{
            Id:              user.ID,
            Username:        user.Username,
            Email:           user.Email,
            MaxRoomReached:  user.MaxRoomReached,
        },
    }, nil
}