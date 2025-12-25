package auth_api

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    "github.com/samber/lo"
)

func (s *AuthServiceAPI) UpsertUsers(ctx context.Context, req *auth_api.UpsertUsersRequest) (*auth_api.UpsertUsersResponse, error) {
    users := mapUsersFromProto(req.Users)
    err := s.authService.UpsertUsers(ctx, users)
    if err != nil {
        return &auth_api.UpsertUsersResponse{}, err
    }
    return &auth_api.UpsertUsersResponse{}, nil
}

func mapUsersFromProto(users []*proto_models.UserUpsertModel) []*models.User {
    return lo.Map(users, func(u *proto_models.UserUpsertModel, _ int) *models.User {
        return &models.User{
            Username: u.Username,
            Email:    u.Email,
            Password: u.Password,
        }
    })
}