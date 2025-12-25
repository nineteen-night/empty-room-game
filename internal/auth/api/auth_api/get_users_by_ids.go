package auth_api

import (
    "context"
    "log"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    "github.com/samber/lo"
)

func (s *AuthServiceAPI) GetUsersByIDs(ctx context.Context, req *auth_api.GetUsersByIDsRequest) (*auth_api.GetUsersByIDsResponse, error) {
    log.Printf("Received GetUsersByIDs request with IDs: %v", req.Ids)

    users, err := s.authService.GetUsersByIDs(ctx, req.Ids)
    if err != nil {
        return &auth_api.GetUsersByIDsResponse{}, err
    }

    return &auth_api.GetUsersByIDsResponse{
        Users: mapUsersToProto(users),
    }, nil
}

func mapUsersToProto(users []*models.User) []*proto_models.UserModel {
    return lo.Map(users, func(user *models.User, _ int) *proto_models.UserModel {
        return &proto_models.UserModel{
            Id:       user.ID,
            Username: user.Username,
            Email:    user.Email,
        }
    })
}