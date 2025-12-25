package auth_api

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    "github.com/samber/lo"
)

func (s *AuthServiceAPI) UpsertPartnerships(ctx context.Context, req *auth_api.UpsertPartnershipsRequest) (*auth_api.UpsertPartnershipsResponse, error) {
    partnerships := mapPartnershipsFromProto(req.Partnerships)
    err := s.authService.UpsertPartnerships(ctx, partnerships)
    if err != nil {
        return &auth_api.UpsertPartnershipsResponse{}, err
    }
    return &auth_api.UpsertPartnershipsResponse{}, nil
}

func mapPartnershipsFromProto(partnerships []*proto_models.PartnershipUpsertModel) []*models.Partnership {
    return lo.Map(partnerships, func(p *proto_models.PartnershipUpsertModel, _ int) *models.Partnership {
        return &models.Partnership{
            Player1ID: p.Player1Id,
            Player2ID: p.Player2Id,
            Status:    "pending",
        }
    })
}