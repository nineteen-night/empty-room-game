package auth_api

import (
    "context"
    "log"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/auth/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    "github.com/samber/lo"
)

func (s *AuthServiceAPI) GetPartnershipsByIDs(ctx context.Context, req *auth_api.GetPartnershipsByIDsRequest) (*auth_api.GetPartnershipsByIDsResponse, error) {
    log.Printf("Received GetPartnershipsByIDs request with IDs: %v", req.Ids)

    partnerships, err := s.authService.GetPartnershipsByIDs(ctx, req.Ids)
    if err != nil {
        return &auth_api.GetPartnershipsByIDsResponse{}, err
    }

    return &auth_api.GetPartnershipsByIDsResponse{
        Partnerships: mapPartnershipsToProto(partnerships),
    }, nil
}

func mapPartnershipsToProto(partnerships []*models.Partnership) []*proto_models.PartnershipModel {
    return lo.Map(partnerships, func(p *models.Partnership, _ int) *proto_models.PartnershipModel {
        return &proto_models.PartnershipModel{
            Id:        p.ID,
            Player1Id: p.Player1ID,
            Player2Id: p.Player2ID,
            Status:    p.Status,
        }
    })
}