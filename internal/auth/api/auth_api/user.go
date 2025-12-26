package auth_api

import (
    "context"
    
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
)

func (s *AuthServiceAPI) GetUser(ctx context.Context, req *auth_api.GetUserRequest) (*auth_api.GetUserResponse, error) {
    user, err := s.authService.GetUser(ctx, req.UserId)
    if err != nil {
        return nil, err
    }
    
    return &auth_api.GetUserResponse{
        User: &proto_models.User{
            Id:              user.ID,
            Username:        user.Username,
            Email:           user.Email,
            MaxRoomReached:  user.MaxRoomReached,
        },
    }, nil
}