package game_api

import (
    "context"
    "log"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/game/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
    "github.com/samber/lo"
)

func (s *GameServiceAPI) GetGameSessionsByIDs(ctx context.Context, req *game_api.GetGameSessionsByIDsRequest) (*game_api.GetGameSessionsByIDsResponse, error) {
    log.Printf("Received GetGameSessionsByIDs request with IDs: %v", req.Ids)

    sessions, err := s.gameService.GetGameSessionsByIDs(ctx, req.Ids)
    if err != nil {
        return &game_api.GetGameSessionsByIDsResponse{}, err
    }

    return &game_api.GetGameSessionsByIDsResponse{
        GameSessions: mapGameSessionsToProto(sessions),
    }, nil
}

func mapGameSessionsToProto(sessions []*models.GameSession) []*proto_models.GameSessionModel {
    return lo.Map(sessions, func(s *models.GameSession, _ int) *proto_models.GameSessionModel {
        return &proto_models.GameSessionModel{
            Id:              s.ID,
            PartnershipId:   s.PartnershipID,
            CurrentRoom:     s.CurrentRoom,
            Status:          s.Status,
            CurrentPlayerId: s.CurrentPlayerID,
        }
    })
}