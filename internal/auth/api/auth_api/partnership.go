package auth_api

import (
    "context"
    
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
)

func (s *AuthServiceAPI) CreatePartnership(ctx context.Context, req *auth_api.CreatePartnershipRequest) (*auth_api.CreatePartnershipResponse, error) {
    partnership, err := s.authService.CreatePartnership(ctx, req.UserId, req.PartnerUsername)
    if err != nil {
        return nil, err
    }
    
    return &auth_api.CreatePartnershipResponse{
        Partnership: &proto_models.Partnership{
            Id:       partnership.ID,
            User1Id:  partnership.User1ID,
            User2Id:  partnership.User2ID,
        },
    }, nil
}

func (s *AuthServiceAPI) TerminatePartnership(ctx context.Context, req *auth_api.TerminatePartnershipRequest) (*auth_api.TerminatePartnershipResponse, error) {
    err := s.authService.TerminatePartnership(ctx, req.PartnershipId)
    if err != nil {
        return nil, err
    }
    
    return &auth_api.TerminatePartnershipResponse{
        Success: true,
    }, nil
}